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
	"bisnis-rinzi/services/pos/entity"
)

type OutboxWorker struct {
	dbClient *postgres.DBClient
}

func NewOutboxWorker(db *postgres.DBClient) *OutboxWorker {
	return &OutboxWorker{dbClient: db}
}

// Start memicu perputaran ticker monitoring outbox POS
func (w *OutboxWorker) Start(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	logger.Info("Background Worker Outbox Event POS Service aktif memantau...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Menghentikan Background Worker POS Outbox...")
			return
		case <-ticker.C:
			w.processPendingEvents(ctx)
		}
	}
}

func (w *OutboxWorker) processPendingEvents(ctx context.Context) {
	tx, err := w.dbClient.Pool.Begin(ctx)
	if err != nil {
		logger.Error("POS Outbox: Gagal menginisialisasi transaksi: %v", err)
		return
	}
	defer tx.Rollback(ctx)

	// SINKRONISASI: Menarik batch 10 event PENDING menggunakan package outbox bersama
	events, err := outbox.FetchPendingEvents(ctx, tx, 10)
	if err != nil {
		logger.Error("POS Outbox: Gagal menarik data antrean outbox: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	logger.Info("POS Outbox Worker menemukan %d transaksi kasir berstatus PENDING, menyinkronkan...", len(events))

	for _, event := range events {
		errDispatch := w.dispatchToOtherServices(ctx, event)

		if errDispatch != nil {
			logger.Error("-> [FAILED] Event %s dengan ID %s gagal didistribusikan: %v", event.EventType, event.ID, errDispatch)
			errStr := errDispatch.Error()
			// Perbarui status lokal pos_db menjadi FAILED via shared library
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusFailed, &errStr)
		} else {
			// Perbarui status lokal pos_db menjadi SUCCESS via shared library
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusSuccess, nil)
			logger.Info("-> [SUCCESS] Event %s dengan ID %s sukses terkirim lintas layanan", event.EventType, event.ID)
		}
	}

	_ = tx.Commit(ctx)
}

func (w *OutboxWorker) dispatchToOtherServices(_ context.Context, e *outbox.Event) error {
	logger.Info("-> [DISPATCH] Memproses penyebaran event %s ID %s", e.EventType, e.ID)

	inventoryURL := os.Getenv("INVENTORY_SERVICE_INTERNAL_URL")
	if inventoryURL == "" {
		inventoryURL = "http://localhost:8082"
	}

	cashURL := os.Getenv("CASH_SERVICE_INTERNAL_URL")
	if cashURL == "" {
		cashURL = "http://localhost:8083"
	}

	financeURL := os.Getenv("FINANCE_SERVICE_INTERNAL_URL")
	if financeURL == "" {
		financeURL = "http://localhost:8086"
	}

	if e.EventType == "SALE_COMPLETED" {
		// Bongkar biner payload menggunakan filter entity POS yang bertipe ketat
		payload, err := entity.UnmarshalSalePayload(e.Payload)
		if err != nil {
			return fmt.Errorf("gagal unmarshal payload penjualan pos: %w", err)
		}

		// 1. DISTRIBUSI KE INVENTORY: POST /internal/inventory/stock/deduct
		for _, item := range payload.Items {
			deductPayload := map[string]interface{}{
				"product_id": item.ProductID,
				"qty":        item.Quantity,
				"reference":  fmt.Sprintf("SALE-%s", e.AggregateID),
			}
			body, _ := json.Marshal(deductPayload)
			if err := sendHTTPRequest("POST", inventoryURL+"/internal/inventory/stock/deduct", body); err != nil {
				return fmt.Errorf("gagal deduct stok (Inventory Service): %w", err)
			}
		}

		// 2. DISTRIBUSI KE CASH: POST /internal/cash/income
		cashPayload := map[string]interface{}{
			"amount":      payload.TotalTransaction,
			"description": fmt.Sprintf("Pendapatan penjualan Retail %s", payload.InvoiceNumber),
			"source":      "POS_SALE",
			"reference":   e.AggregateID,
		}
		cashBody, _ := json.Marshal(cashPayload)
		if err := sendHTTPRequest("POST", cashURL+"/internal/cash/income", cashBody); err != nil {
			return fmt.Errorf("gagal catat kas masuk (Cash Service): %w", err)
		}

		// 3. DISTRIBUSI KE FINANCE GATE: POST /api/finance/inbox-events
		// Memasukkan ke format penampung string tekstual agar match dengan JournalIncomingRequest
		binaryOperational, _ := json.Marshal(payload)
		inboxPayload := map[string]interface{}{
			"id":             e.ID,
			"aggregate_type": "RETAIL_ORDER",
			"aggregate_id":   e.AggregateID,
			"event_type":     "ORDER_COMPLETED",
			"payload":        string(binaryOperational), // Dikirim berupa teks string JSON murni
			"created_at":     e.CreatedAt,
		}

		financeBody, _ := json.Marshal(inboxPayload)
		targetFinanceURL := fmt.Sprintf("%s/api/finance/inbox-events", financeURL)

		logger.Info("Mengirimkan enkapsulasi teks omzet kasir ke Finance Gate: %s", targetFinanceURL)
		if err := sendHTTPRequest("POST", targetFinanceURL, financeBody); err != nil {
			return fmt.Errorf("gagal mentransfer antrean event omzet ritel ke Finance Service: %w", err)
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
