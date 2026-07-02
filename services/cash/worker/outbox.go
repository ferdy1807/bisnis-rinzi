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
	"bisnis-rinzi/services/cash/entity"
)

type CashOutboxWorker struct {
	dbClient *postgres.DBClient
}

// NewCashOutboxWorker menginstansiasi worker outbox penjamin eventual consistency untuk sistem keuangan kas.
func NewCashOutboxWorker(db *postgres.DBClient) *CashOutboxWorker {
	return &CashOutboxWorker{dbClient: db}
}

// Start memicu loop berkala ticker pengawasan antrean outbox cash_db.
func (w *CashOutboxWorker) Start(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	logger.Info("Background Worker Outbox Event CASH Service aktif memantau...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Menghentikan Background Worker CASH Outbox...")
			return
		case <-ticker.C:
			w.processPendingCashEvents(ctx)
		}
	}
}

// processPendingCashEvents mengisolasi penarikan baris data PENDING dalam ruang transaksi database tertutup.
func (w *CashOutboxWorker) processPendingCashEvents(ctx context.Context) {
	tx, err := w.dbClient.Pool.Begin(ctx)
	if err != nil {
		logger.Error("CASH Outbox: Gagal menginisialisasi transaksi database: %v", err)
		return
	}
	defer tx.Rollback(ctx)

	// Mengambil antrean event berstatus PENDING dari cash_db secara FIFO menggunakan shared package backend
	events, err := outbox.FetchPendingEvents(ctx, tx, 10)
	if err != nil {
		logger.Error("CASH Outbox: Gagal memuat baris antrean outbox: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	logger.Info("CASH Outbox Worker menemukan %d data mutasi operasional PENDING, memproses kliring...", len(events))

	for _, event := range events {
		errDispatch := w.dispatchToFinanceService(ctx, event)

		if errDispatch != nil {
			logger.Error("-> [FAILED] Mutasi operasional %s ID %s gagal tersinkronisasi: %v", event.EventType, event.ID, errDispatch)
			errStr := errDispatch.Error()
			// Perbarui status lokal cash_db menjadi FAILED via shared library agar bisa ditinjau ulang
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusFailed, &errStr)
		} else {
			// Perbarui status lokal cash_db menjadi SUCCESS (Kliring Berhasil)
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusSuccess, nil)
			logger.Info("-> [SUCCESS] Event %s dengan ID %s sukses terposting menuju Finance Ingestion Gate", event.EventType, event.ID)
		}
	}

	_ = tx.Commit(ctx)
}

// dispatchToFinanceService mengeksekusi normalisasi payload data dan mentransfernya via HTTP REST-API.
func (w *CashOutboxWorker) dispatchToFinanceService(_ context.Context, e *outbox.Event) error {
	financeURL := os.Getenv("FINANCE_SERVICE_INTERNAL_URL")
	if financeURL == "" {
		financeURL = "http://localhost:8086"
	}

	var normalizedPayload map[string]interface{}
	var targetAggregateType string
	var targetEventType string

	// Kasus Bisnis A: Rekonsiliasi Hasil Tutup Shift Kasir Lapangan
	if e.AggregateType == "CASH_SESSION" || e.AggregateType == "SHIFT_SESSION" {
		localShift, err := entity.UnmarshalLocalShift(e.Payload)
		if err != nil {
			return fmt.Errorf("gagal mengurai sub-payload shift lokal: %w", err)
		}

		// Transformasikan skema data lokal agar sesuai dengan kebutuhan entitas StoreShiftPayload di finance pusat
		normalizedPayload = map[string]interface{}{
			"session_id":     localShift.SessionID,
			"user_id":        localShift.CashierID,
			"opening_cash":   localShift.OpeningCash,
			"total_income":   localShift.TotalIncome,
			"manual_deposit": localShift.ManualDeposit,
			"total_expense":  localShift.TotalExpense,
			"actual_cash":    localShift.ActualCash,
			"expected_cash":  localShift.ExpectedCash,
			"discrepancy":    localShift.DifferenceCash,
			"status_shift":   "CLOSED",
			"timestamp":      localShift.ClosedTimestamp,
		}

		targetAggregateType = "SHIFT_SESSION" // Dipaksa menjadi SHIFT_SESSION agar lolos switch-case worker finance
		targetEventType = "SHIFT_CLOSED"

		// Kasus Bisnis B: Sinkronisasi Pengeluaran Kas Tunai / Pembayaran Biaya Operasional Toko
	} else if e.AggregateType == "EXPENSE" {
		localExpense, err := entity.UnmarshalLocalExpense(e.Payload)
		if err != nil {
			return fmt.Errorf("gagal mengurai sub-payload pengeluaran lokal: %w", err)
		}

		// Transformasikan skema beban operasional
		normalizedPayload = map[string]interface{}{
			"reference_id": localExpense.ExpenseID,
			"amount":       localExpense.Amount,
			"description":  fmt.Sprintf("Pengeluaran Kas Toko (%s) - %s", localExpense.Category, localExpense.Description),
			"category":     localExpense.Category,
			"timestamp":    localExpense.Timestamp,
		}

		targetAggregateType = "RETAIL_ORDER" // Diarahkan ke rumpun retail order/cash mutasi finansial
		targetEventType = "EXPENSE_POSTED"
	} else {
		// Abaikan tipe agregat yang tidak ditujukan untuk pembukuan akuntansi umum
		return nil
	}

	// Langkah 1: Ubah objek operasional terurai menjadi biner string teks JSON murni
	binaryPayload, err := json.Marshal(normalizedPayload)
	if err != nil {
		return fmt.Errorf("gagal melakukan enkapsulasi sub-payload operasional: %w", err)
	}

	// Langkah 2: Bungkus ke dalam skema ketat kontrak JournalIncomingRequest
	// Kolom 'payload' WAJIB diisi berupa STRING tekstual agar lolos dari validasi decoder HTTP finance pusat
	inboxPayload := map[string]interface{}{
		"id":             e.ID,
		"aggregate_type": targetAggregateType,
		"aggregate_id":   e.AggregateID,
		"event_type":     targetEventType,
		"payload":        string(binaryPayload), // Konversi biner slice byte menjadi raw tekstual string
		"created_at":     e.CreatedAt,
	}

	body, _ := json.Marshal(inboxPayload)
	targetURL := fmt.Sprintf("%s/api/finance/inbox-events", financeURL) // Ditembak lurus menuju gerbang penerima terpadu baru

	logger.Info("Mentransfer enkapsulasi teks berita acara dari CASH ke Finance Ingestion Gate: %s", targetURL)
	if err := sendHTTPRequest("POST", targetURL, body); err != nil {
		return fmt.Errorf("gagal mentransfer antrean event cash ke Finance Ingestion Gate: %w", err)
	}

	return nil
}

// sendHTTPRequest melakukan dial koneksi HTTP client dengan pembatasan waktu timeout demi performa sistem.
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
