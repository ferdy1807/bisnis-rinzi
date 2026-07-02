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
	"bisnis-rinzi/services/inventory/entity"
)

type InventoryOutboxWorker struct {
	dbClient *postgres.DBClient
}

// NewInventoryOutboxWorker menginstansiasi worker outbox penjamin eventual consistency untuk sistem inventory gudang.
func NewInventoryOutboxWorker(db *postgres.DBClient) *InventoryOutboxWorker {
	return &InventoryOutboxWorker{dbClient: db}
}

// Start memicu loop berkala ticker pengawasan antrean outbox inventory_db.
func (w *InventoryOutboxWorker) Start(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	logger.Info("Background Worker Outbox Event INVENTORY Service aktif memantau...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Menghentikan Background Worker INVENTORY Outbox...")
			return
		case <-ticker.C:
			w.processPendingInventoryEvents(ctx)
		}
	}
}

// processPendingInventoryEvents mengisolasi penarikan baris data PENDING dalam ruang transaksi database tertutup.
func (w *InventoryOutboxWorker) processPendingInventoryEvents(ctx context.Context) {
	tx, err := w.dbClient.Pool.Begin(ctx)
	if err != nil {
		logger.Error("INVENTORY Outbox: Gagal menginisialisasi transaksi database: %v", err)
		return
	}
	defer tx.Rollback(ctx)

	// Mengambil antrean event berstatus PENDING dari inventory_db secara FIFO menggunakan shared package backend
	events, err := outbox.FetchPendingEvents(ctx, tx, 10)
	if err != nil {
		logger.Error("INVENTORY Outbox: Gagal memuat baris antrean outbox: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	logger.Info("INVENTORY Outbox Worker menemukan %d data mutasi stok PENDING, menyinkronkan...", len(events))

	for _, event := range events {
		errDispatch := w.dispatchToFinanceService(ctx, event)

		if errDispatch != nil {
			logger.Error("-> [FAILED] Mutasi stok %s ID %s gagal tersinkronisasi: %v", event.EventType, event.ID, errDispatch)
			errStr := errDispatch.Error()
			// Perbarui status lokal inventory_db menjadi FAILED via shared library agar bisa ditinjau ulang
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusFailed, &errStr)
		} else {
			// Perbarui status lokal inventory_db menjadi SUCCESS (Kliring Berhasil)
			_ = outbox.UpdateEventStatusTx(ctx, tx, event.ID, outbox.StatusSuccess, nil)
			logger.Info("-> [SUCCESS] Event %s dengan ID %s sukses terposting menuju Finance Ingestion Gate", event.EventType, event.ID)
		}
	}

	_ = tx.Commit(ctx)
}

// dispatchToFinanceService mengeksekusi normalisasi payload data logistik dan mentransfernya via HTTP REST-API.
func (w *InventoryOutboxWorker) dispatchToFinanceService(_ context.Context, e *outbox.Event) error {
	financeURL := os.Getenv("FINANCE_SERVICE_INTERNAL_URL")
	if financeURL == "" {
		financeURL = "http://localhost:8086"
	}

	// Memproses data mutasi stok masuk gudang untuk kebutuhan sinkronisasi akuntansi nilai persediaan barang dagangan
	if e.AggregateType == "INVENTORY_STOCK" || e.AggregateType == "STOCK_MOVEMENT" {
		localStock, err := entity.UnmarshalLocalStock(e.Payload)
		if err != nil {
			return fmt.Errorf("gagal mengurai sub-payload stok lokal: %w", err)
		}

		// Transformasikan skema data lokal agar sesuai dengan kebutuhan entitas InventoryStockPayload di finance pusat
		normalizedPayload := map[string]interface{}{
			"product_id":     localStock.ProductID,
			"mutation_type":  localStock.MutationType,
			"qty":            localStock.Quantity,
			"unit_cost":      localStock.UnitCost,
			"total_cost":     localStock.TotalCost, // Penting untuk alokasi Debet Akun Persediaan dan Kredit Akun Kas/Utang
			"reference_code": localStock.ReferenceCode,
			"timestamp":      localStock.Timestamp,
		}

		// Langkah A: Ubah objek operasional logistik menjadi biner string teks JSON murni
		binaryPayload, err := json.Marshal(normalizedPayload)
		if err != nil {
			return fmt.Errorf("gagal melakukan enkapsulasi sub-payload operasional stok: %w", err)
		}

		// Langkah B: Bungkus ke dalam skema ketat kontrak JournalIncomingRequest
		// Properti 'payload' WAJIB diisi berupa STRING tekstual agar lolos dari validasi decoder HTTP finance pusat
		inboxPayload := map[string]interface{}{
			"id":             e.ID,
			"aggregate_type": "INVENTORY_STOCK", // Dipaksa menjadi INVENTORY_STOCK agar lolos switch-case worker finance
			"aggregate_id":   e.AggregateID,
			"event_type":     e.EventType,           // Meneruskan "INCOMING_SUPPLIER" atau "STOCK_ADJUSTED"
			"payload":        string(binaryPayload), // Konversi biner slice byte menjadi raw tekstual string JSON
			"created_at":     e.CreatedAt,
		}

		body, _ := json.Marshal(inboxPayload)
		targetURL := fmt.Sprintf("%s/api/finance/inbox-events", financeURL) // Diarahkan langsung menuju gerbang penampung internal terpadu baru

		logger.Info("Mentransfer enkapsulasi teks mutasi logistik dari INVENTORY ke Finance Ingestion Gate: %s", targetURL)
		if err := sendHTTPRequest("POST", targetURL, body); err != nil {
			return fmt.Errorf("gagal mentransfer antrean event stok ke Finance Ingestion Gate: %w", err)
		}
	}

	return nil
}

// sendHTTPRequest melakukan dial koneksi HTTP client dengan pembatasan waktu timeout demi menjaga performa subsistem.
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
