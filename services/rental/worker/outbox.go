package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/rental/entity"
)

type RentalOutboxWorker struct {
	dbClient *postgres.DBClient
}

func NewRentalOutboxWorker(db *postgres.DBClient) *RentalOutboxWorker {
	return &RentalOutboxWorker{dbClient: db}
}

func (w *RentalOutboxWorker) Start(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	logger.Info("Background Worker Outbox Event RENTAL Service aktif memantau...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Menghentikan Background Worker RENTAL Outbox...")
			return
		case <-ticker.C:
			w.processPendingRentalEvents(ctx)
		}
	}
}

func (w *RentalOutboxWorker) processPendingRentalEvents(ctx context.Context) {
	tx, err := w.dbClient.Pool.Begin(ctx)
	if err != nil {
		logger.Error("RENTAL Outbox: Gagal membuka transaksi database: %v", err)
		return
	}
	defer tx.Rollback(ctx)

	// Fetch 10 antrean event rental PENDING secara FIFO
	events, err := outbox.FetchPendingEvents(ctx, tx, 10)
	if err != nil {
		logger.Error("RENTAL Outbox: Gagal memuat antrean outbox: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	logger.Info("RENTAL Outbox Worker menemukan %d sirkulasi sewa PENDING, mengeksekusi kliring...", len(events))

	for _, event := range events {
		errDispatch := w.dispatchToFinanceService(ctx, event)

		if errDispatch != nil {
			logger.Error("-> [FAILED] Kontrak Sewa %s ID %s gagal dikirim: %v", event.AggregateID, event.ID, errDispatch)
			errStr := errDispatch.Error()
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusFailed, &errStr)
		} else {
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusSuccess, nil)
			logger.Info("-> [SUCCESS] Berita acara rental %s ID %s sukses terposting ke keuangan", event.EventType, event.ID)
		}
	}

	_ = tx.Commit(ctx)
}

func (w *RentalOutboxWorker) dispatchToFinanceService(_ context.Context, e *outbox.Event) error {
	logger.Info("-> [DISPATCH RENTAL] Mengirimkan data %s ID %s menuju pusat akuntansi", e.EventType, e.ID)

	financeURL := os.Getenv("FINANCE_SERVICE_INTERNAL_URL")
	if financeURL == "" {
		financeURL = "http://localhost:8086"
	}

	if e.AggregateType == "RENTAL_RESERVATION" || e.AggregateType == "RESERVATION" {
		var binaryPayload []byte
		var err error
		var targetEventType string

		// Evaluasi Berdasarkan Spesifikasi Tipe Event (Event Driven Routing)
		switch e.EventType {
		case "PRODUCT_RETURN_PROCESSED":
			localReturn, errParse := entity.UnmarshalProductReturnPayload(e.Payload)
			if errParse != nil {
				return fmt.Errorf("gagal unmarshal payload return lokal: %w", errParse)
			}

			// Mapping data normalisasi untuk pengembalian produk
			normalizedReturn := map[string]interface{}{
				"return_id":         localReturn.ReturnID,
				"reservation_id":    localReturn.ReservationID,
				"late_days":         localReturn.LateDays,
				"total_late_fees":   localReturn.TotalLateFees,
				"total_damage_fees": localReturn.TotalDamageFees,
				"remaining_payment": localReturn.RemainingPayment,
				"grand_total_paid":  localReturn.GrandTotalPaid,
				"timestamp":         localReturn.Timestamp,
			}
			binaryPayload, err = json.Marshal(normalizedReturn)
			targetEventType = "PRODUCT_RETURN_PROCESSED"

		default:
			// Fallback ke pemesanan / reservasi awal boks seserahan
			localReserve, errParse := entity.UnmarshalReservationPayload(e.Payload)
			if errParse != nil {
				return fmt.Errorf("gagal unmarshal payload reservasi lokal: %w", errParse)
			}

			normalizedReserve := map[string]interface{}{
				"reservation_id": localReserve.ReservationID,
				"customer_id":    localReserve.CustomerID,
				"box_code":       localReserve.BoxCode,
				"rental_price":   localReserve.RentalPrice,
				"deposit_amount": localReserve.DepositAmount,
				"penalty_amount": localReserve.PenaltyAmount,
				"status_booking": localReserve.StatusBooking,
				"deadline_date":  localReserve.DeadlineDate,
			}
			binaryPayload, err = json.Marshal(normalizedReserve)
			targetEventType = "RENTAL_CONFIRMED"
		}

		if err != nil {
			return fmt.Errorf("gagal mengemas biner sub-payload operasional sewa: %w", err)
		}

		// Enkapsulasi biner menjadi teks string JSON agar sesuai dengan skema JournalIncomingRequest
		inboxPayload := map[string]interface{}{
			"id":             e.ID,
			"aggregate_type": "RENTAL_RESERVATION", // Sesuai dengan switch-case pada worker finance
			"aggregate_id":   e.AggregateID,
			"event_type":     targetEventType,
			"payload":        string(binaryPayload), // Mengirim raw tekstual string JSON
			"created_at":     e.CreatedAt,
		}

		body, _ := json.Marshal(inboxPayload)
		targetURL := fmt.Sprintf("%s/api/finance/inbox-events", financeURL)

		logger.Info("Mentransfer data berita acara sewa hantaran ke: %s", targetURL)
		if err := sendHTTPRequest("POST", targetURL, body); err != nil {
			return fmt.Errorf("gagal transfer jurnal sewa ke Finance Service: %w", err)
		}
	}

	return nil
}

func sendHTTPRequest(method, url string, payload []byte) error {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d error dari %s", resp.StatusCode, url)
	}
	return nil
}
