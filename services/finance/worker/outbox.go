package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/finance/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FinanceJournalRepository mendefinisikan kontrak operasi ke buku besar akuntansi finance_db
type FinanceJournalRepository interface {
	InsertJournalEntry(ctx context.Context, tx pgx.Tx, aggregateID, narasi string, debetAcc, kreditAcc string, total float64) error
	IsEventProcessed(ctx context.Context, tx pgx.Tx, eventID string) (bool, error)
	MarkEventAsProcessed(ctx context.Context, tx pgx.Tx, eventID string) error
}

// FinanceOutboxWorker mengelola background proses penyerapan event bisnis global
type FinanceOutboxWorker struct {
	dbPool       *pgxpool.Pool
	journalRepo  FinanceJournalRepository
	pollInterval time.Duration
	batchSize    int
}

// NewFinanceOutboxWorker menginisialisasi instansiasi worker akuntansi otomatis
func NewFinanceOutboxWorker(pool *pgxpool.Pool, repo FinanceJournalRepository, interval time.Duration, batchSize int) *FinanceOutboxWorker {
	return &FinanceOutboxWorker{
		dbPool:       pool,
		journalRepo:  repo,
		pollInterval: interval,
		batchSize:    batchSize,
	}
}

// Start menjalankan background loop ticker untuk memproses event secara berkala
func (w *FinanceOutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	log.Printf("[FINANCE WORKER] Mengaktifkan mesin audit neraca otomatis (Interval: %v)...", w.pollInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("[FINANCE WORKER] Menghentikan pemrosesan worker secara aman (Graceful Shutdown)...")
			return
		case <-ticker.C:
			if err := w.ProcessNextBatch(ctx); err != nil {
				log.Printf("[FINANCE WORKER ERROR] Gagal memproses batch finansial: %v", err)
			}
		}
	}
}

// File: services/finance/worker/outbox.go

// File: services/finance/worker/outbox.go

func (w *FinanceOutboxWorker) ProcessNextBatch(ctx context.Context) error {
	// 1. Ambil data PENDING tanpa membuka transaksi panjang di awal
	// Kita buka koneksi pendek hanya untuk fetch antrean
	fetchTx, err := w.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("worker: gagal membuka transaksi fetch: %w", err)
	}
	events, err := outbox.FetchPendingEvents(ctx, fetchTx, w.batchSize)
	_ = fetchTx.Rollback(ctx) // Langsung tutup setelah fetch selesai

	if err != nil {
		return fmt.Errorf("worker: gagal mengambil antrean outbox: %w", err)
	}
	if len(events) == 0 {
		return nil
	}

	log.Printf("[FINANCE WORKER] Menemukan %d mutasi baru untuk diverifikasi...", len(events))

	// 2. ITERASI: Setiap event memiliki ruang transaksinya sendiri (Terisolasi!)
	for _, e := range events {
		err := w.processSingleEventWithTx(ctx, *e)
		if err != nil {
			log.Printf("[FINANCE WORKER CRITICAL] ID Event %s Gagal Terjurnal: %v", e.ID, err)
		} else {
			log.Printf("[FINANCE WORKER SUCCESS] ID Event %s berhasil diposting ke Jurnal Umum", e.ID)
		}
	}

	return nil
}

// Fungsi pembantu baru untuk mengisolasi transaksi per satu baris event
func (w *FinanceOutboxWorker) processSingleEventWithTx(ctx context.Context, e outbox.Event) error {
	tx, err := w.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("gagal membuka transaksi event: %w", err)
	}
	defer tx.Rollback(ctx)

	// A. Cek Idempotensi secara aman di dalam transaksi bersih
	isProcessed, err := w.journalRepo.IsEventProcessed(ctx, tx, e.ID)
	if err != nil {
		return fmt.Errorf("gagal memeriksa idempotensi: %w", err)
	}
	if isProcessed {
		_ = outbox.UpdateEventStatusTx(ctx, tx, e.ID, outbox.StatusSuccess, nil)
		return tx.Commit(ctx)
	}

	inboxEvent := &entity.FinanceInboxEvent{
		ID:            e.ID,
		AggregateType: e.AggregateType,
		AggregateID:   e.AggregateID,
		EventType:     e.EventType,
		Payload:       e.Payload,
		CreatedAt:     e.CreatedAt,
	}

	// B. Routing Logic
	var processErr error
	switch inboxEvent.AggregateType {
	case "RETAIL_ORDER":
		processErr = w.handleRetailOrder(ctx, tx, inboxEvent)
	case "SHIFT_SESSION":
		processErr = w.handleStoreShift(ctx, tx, inboxEvent)
	case "RENTAL_RESERVATION", "RESERVATION":
		if inboxEvent.EventType == "PRODUCT_RETURN_PROCESSED" {
			processErr = w.handleRentalReturn(ctx, tx, inboxEvent)
		} else {
			processErr = w.handleRentalBooking(ctx, tx, inboxEvent)
		}
	default:
		processErr = fmt.Errorf("tipe agregat '%s' tidak didukung", e.AggregateType)
	}

	// C. Evaluasi Hasil Akhir untuk Event Ini
	if processErr != nil {
		errStr := processErr.Error()
		_ = outbox.UpdateEventStatusTx(ctx, tx, e.ID, outbox.StatusFailed, &errStr)
		_ = tx.Commit(ctx) // Commit status FAILED agar tersimpan di DB dan tidak loop terus menerus
		return processErr
	}

	// D. Tandai Sukses dan Idempotensi jika semua kueri akuntansi lolos
	if err := w.journalRepo.MarkEventAsProcessed(ctx, tx, e.ID); err != nil {
		return fmt.Errorf("gagal menandai idempotensi: %w", err)
	}
	_ = outbox.UpdateEventStatusTx(ctx, tx, e.ID, outbox.StatusSuccess, nil)

	return tx.Commit(ctx)
}

// Handler 1: Penjurnalan POS Kasir Otomatis (Portal Toko)
func (w *FinanceOutboxWorker) handleRetailOrder(ctx context.Context, tx pgx.Tx, ev *entity.FinanceInboxEvent) error {
	payload, err := ev.ToStoreOrder()
	if err != nil {
		return fmt.Errorf("retail_handler: gagal mengurai payload: %w", err)
	}

	narasi := fmt.Sprintf("Pendapatan POS Kasir Ritel Rapi - Invoice No: %s", payload.InvoiceNumber)
	netSales := payload.TotalAmount - payload.TaxAmount

	// Rule Akuntansi: Debet Kas/QRIS (110000), Kredit Pendapatan Toko (410000)
	err = w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID, narasi, "110000", "410000", netSales)
	if err != nil {
		return err
	}

	// Post Jurnal PPN jika ada
	if payload.TaxAmount > 0 {
		err = w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID, narasi+" (PPN)", "110000", "210000", payload.TaxAmount)
		if err != nil {
			return err
		}
	}

	return nil
}

// Handler 2: Penjurnalan Selisih Laci Kasir / Sesi Shift (Portal Toko)
func (w *FinanceOutboxWorker) handleStoreShift(ctx context.Context, tx pgx.Tx, ev *entity.FinanceInboxEvent) error {
	payload, err := ev.ToStoreShift()
	if err != nil {
		return fmt.Errorf("shift_handler: gagal mengurai payload: %w", err)
	}

	if payload.StatusShift != "CLOSED" {
		return nil
	}

	log.Printf("[SHIFT CLOSED] Sesi: %s | Omzet: %.2f | Expense: %.2f | Manual: %.2f | Selisih: %.2f", 
		payload.SessionID, payload.TotalIncome, payload.TotalExpense, payload.ManualDeposit, payload.Discrepancy)

	// Jurnal Omzet Shift (Total Penerimaan POS) - DIHAPUS
	// Penjurnalan pendapatan POS secara spesifik sudah dilakukan per-transaksi (RETAIL_ORDER)
	// sehingga tidak perlu di-jurnal ulang saat penutupan shift agar tidak terjadi double-posting.

	// Jurnal Kas Masuk Manual (Ekstra)
	if payload.ManualDeposit > 0 {
		narasi := fmt.Sprintf("Setoran Kas Manual Tambahan - Sesi: %s", payload.SessionID)
		if err := w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID+"-MAN", narasi, "110000", "421000", payload.ManualDeposit); err != nil {
			return err
		}
	}

	// Jurnal Pengeluaran Operasional Shift
	if payload.TotalExpense > 0 {
		narasi := fmt.Sprintf("Pengeluaran Kas/Biaya Operasional - Sesi: %s", payload.SessionID)
		if err := w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID+"-EXP", narasi, "720000", "110000", payload.TotalExpense); err != nil {
			return err
		}
	}

	// Jurnal Selisih Kasir Aktual vs Sistem
	if payload.Discrepancy != 0 {
		narasi := fmt.Sprintf("Rekonsiliasi Selisih Fisik Laci Shift Kasir - Sesi: %s", payload.SessionID)
		if payload.Discrepancy > 0 {
			// Surplus: Kas bertambah di laci (Debet Kas, Kredit Pendapatan Lain-lain 421000)
			if err := w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID+"-DIS", narasi, "110000", "421000", payload.Discrepancy); err != nil {
				return err
			}
		} else {
			// Minus: Kas tekor (Debet Beban Operasional 720000, Kredit Kas berkurang)
			if err := w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID+"-DIS", narasi, "720000", "110000", MathAbs(payload.Discrepancy)); err != nil {
				return err
			}
		}
	}

	return nil
}

// Handler 3: Penjurnalan Reservasi Awal Hantaran Pernikahan (Portal Sewa)
func (w *FinanceOutboxWorker) handleRentalBooking(ctx context.Context, tx pgx.Tx, ev *entity.FinanceInboxEvent) error {
	payload, err := ev.ToRentalBooking()
	if err != nil {
		return fmt.Errorf("rental_handler: gagal mengurai payload: %w", err)
	}

	// Jurnal Omzet Pokok Sewa HANYA sebesar uang muka (Cash Basis)
	if payload.DepositAmount > 0 {
		narasi := fmt.Sprintf("Pembayaran Awal Sewa - Reservasi: %s", payload.ReservationID)
		err = w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID, narasi, "110000", "420000", payload.DepositAmount)
		if err != nil {
			return err
		}
	}

	return nil
}

// Handler 4: Penjurnalan Pengembalian Box, Denda Kerusakan, Keterlambatan, & Pelunasan Sisa Kontrak (Portal Sewa)
func (w *FinanceOutboxWorker) handleRentalReturn(ctx context.Context, tx pgx.Tx, ev *entity.FinanceInboxEvent) error {
	payload, err := ev.ToRentalReturn()
	if err != nil {
		return fmt.Errorf("rental_return_handler: gagal mengurai payload fdc3ce58: %w", err)
	}

	narasiPusat := fmt.Sprintf("Penyelesaian Kontrak Sewa Box - Return ID: %s, Res: %s", payload.ReturnID, payload.ReservationID)

	// 1. Jurnal Pembayaran Sisa Kontrak Sewa (Jika ada sisa piutang/pelunasan)
	if payload.RemainingPayment > 0 {
		narasiPelunasan := fmt.Sprintf("%s (Pelunasan Sisa Kontrak)", narasiPusat)
		err = w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID, narasiPelunasan, "110000", "420000", payload.RemainingPayment)
		if err != nil {
			return fmt.Errorf("gagal menjurnal pelunasan sewa: %w", err)
		}
	}

	// 2. Jurnal Klaim Pendapatan Denda Keterlambatan (Late Fees)
	if payload.TotalLateFees > 0 {
		narasiKeterlambatan := fmt.Sprintf("%s (Denda Keterlambatan %d Hari)", narasiPusat, payload.LateDays)
		err = w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID, narasiKeterlambatan, "110000", "421000", payload.TotalLateFees)
		if err != nil {
			return fmt.Errorf("gagal menjurnal denda keterlambatan: %w", err)
		}
	}

	// 3. Jurnal Klaim Pendapatan Denda Kerusakan Mika/Akrilik (Damage Fees)
	if payload.TotalDamageFees > 0 {
		narasiKerusakan := fmt.Sprintf("%s (Denda Kerusakan Komponen Fisik)", narasiPusat)
		err = w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID, narasiKerusakan, "110000", "421000", payload.TotalDamageFees)
		if err != nil {
			return fmt.Errorf("gagal menjurnal denda kerusakan mika: %w", err)
		}
	}

	return nil
}

// Handler 5: Penjurnalan Valuasi Persediaan Nilai Barang Gudang / Restock Supplier (Portal Logistik)
func (w *FinanceOutboxWorker) handleInventoryStock(ctx context.Context, tx pgx.Tx, ev *entity.FinanceInboxEvent) error {
	payload, err := ev.ToInventoryStock()
	if err != nil {
		return fmt.Errorf("inventory_handler: gagal mengurai payload logistik: %w", err)
	}

	narasiGudang := fmt.Sprintf("Restock Aset Inventaris Barang Gudang - PO: %s, Item: %s", payload.ReferenceCode, payload.ProductID)

	// Rule Akuntansi: Debet Persediaan Barang Dagang (1105), Kredit Kas (1101) atau Utang Dagang (2102)
	if payload.MutationType == "INCOMING_SUPPLIER" {
		return w.journalRepo.InsertJournalEntry(ctx, tx, ev.AggregateID, narasiGudang, "1105-PERSEDIAAN-BARANG", "1101-KAS", payload.TotalCost)
	}

	return nil
}

// MathAbs helper lokal pengganti math.Abs bertipe float64
func MathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
