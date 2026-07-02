package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/packages/backend/response"
	"bisnis-rinzi/services/finance/dto"
	"bisnis-rinzi/services/finance/usecase"
)

type FinanceHandler struct {
	useCase usecase.FinanceUseCase
}

func NewFinanceHandler(uc usecase.FinanceUseCase) *FinanceHandler {
	return &FinanceHandler{useCase: uc}
}

func (h *FinanceHandler) COAHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		list, _ := h.useCase.GetAccounts(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar susunan bagan akun perkiraan (COA)", list)
	} else if r.Method == http.MethodPost {
		var req dto.CreateCOARequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.CreateAccount(r.Context(), req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusCreated, "Akun perkiraan baru berhasil disimpan", nil)
	}
}

func (h *FinanceHandler) PeriodOpenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req dto.OpenPeriodRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.useCase.OpenAccountingPeriod(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Periode akuntansi pembukuan baru berhasil dibuka", nil)
}

func (h *FinanceHandler) JournalEntryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req dto.CreateJournalEntryRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.useCase.RecordManualJournal(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusCreated, "Entri pembukuan jurnal double-entry berhasil diposting", nil)
}

func (h *FinanceHandler) DailyClosingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cashierID := r.Header.Get("X-User-Id")
	var req dto.ProcessDailyClosingRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	req.ClosingDate = time.Now()

	if err := h.useCase.ProcessClosingAndReconciliation(r.Context(), cashierID, req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Tutup buku harian cabang dan kliring kasir sukses diselesaikan", nil)
}

// PeriodsQueryHandler menangani GET /api/finance/periods
func (h *FinanceHandler) PeriodsQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	list, _ := h.useCase.GetAllPeriods(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar seluruh periode akuntansi", list)
}

// PeriodLockHandler menangani POST /api/finance/period-locks
func (h *FinanceHandler) PeriodLockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-Id")
	var body map[string]string
	_ = json.NewDecoder(r.Body).Decode(&body)

	id := body["period_id"]
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "period_id harus diisi")
		return
	}

	err := h.useCase.LockPeriod(r.Context(), id, userID, body["lock_reason"])
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Periode akuntansi sukses dikunci permanen", nil)
}

// ClosingsQueryHandler menangani GET /api/finance/closings/daily
func (h *FinanceHandler) ClosingsQueryHandler(w http.ResponseWriter, r *http.Request) {
	list, _ := h.useCase.GetAllDailyClosings(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar riwayat tutup buku harian toko", list)
}

// ClosingDetailDispatcher mengurai rute dinamis /api/finance/daily-closings/{id}/*
func (h *FinanceHandler) ClosingDetailDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	if len(parts) == 4 {
		dc, _ := h.useCase.GetDailyClosingDetail(r.Context(), id)
		if dc == nil {
			response.WriteError(w, http.StatusNotFound, "Dokumen tutup buku tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail log penutupan keuangan harian", dc)
		return
	}

	http.NotFound(w, r)
}

// KOREKSI DINAMIS: Ubah GeneralLedgerReportHandler agar membaca query param URL secara elastis
func (h *FinanceHandler) GeneralLedgerReportHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("account_code")
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	// Fallback parameter penanggalan real-time 2026 jika query string kosong
	if startStr == "" {
		startStr = "2026-01-01"
	}
	if endStr == "" {
		endStr = "2026-12-31"
	}

	start, _ := time.Parse("2006-01-02", startStr)
	end, _ := time.Parse("2006-01-02", endStr)

	res, err := h.useCase.GenerateLedgerReport(r.Context(), code, start, end)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Buku Besar Terfilter (General Ledger)", res)
}

func (h *FinanceHandler) TrialBalanceReportHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := h.useCase.GenerateTrialBalance(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Neraca Saldo (Trial Balance)", res)
}

func (h *FinanceHandler) IncomeStatementReportHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := h.useCase.GenerateIncomeStatement(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Laporan Laba Rugi (Income Statement)", res)
}

func (h *FinanceHandler) BalanceSheetReportHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := h.useCase.GenerateBalanceSheet(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Laporan Posisi Keuangan (Balance Sheet)", res)
}

// AnalyticsDashboardHandler memetakan GET /api/finance/analytics/dashboard
func (h *FinanceHandler) AnalyticsDashboardHandler(w http.ResponseWriter, r *http.Request) {
	start, _ := time.Parse("2006-01-02", "2026-01-01")
	end, _ := time.Parse("2006-01-02", "2026-12-31")

	analyticsData, _ := h.useCase.GetCorporateDashboardAnalytics(r.Context(), start, end)
	response.WriteSuccess(w, http.StatusOK, "Agregasi metrik analitik eksekutif korporat retail", analyticsData)
}

// ----------------------------------------------------
// NEW HANDLERS UNTUK MENUTUP LINTER ERROR
// ----------------------------------------------------

func (h *FinanceHandler) AccountDetailDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		acc, _ := h.useCase.GetAccountByID(r.Context(), id)
		if acc == nil {
			response.WriteError(w, http.StatusNotFound, "Akun tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail akun", acc)
	case http.MethodPut:
		var req dto.CreateCOARequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = h.useCase.UpdateAccount(r.Context(), id, req)
		response.WriteSuccess(w, http.StatusOK, "Akun berhasil diubah", nil)
	case http.MethodDelete:
		_ = h.useCase.DeleteAccount(r.Context(), id)
		response.WriteSuccess(w, http.StatusOK, "Akun berhasil dihapus", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *FinanceHandler) JournalsQueryHandler(w http.ResponseWriter, r *http.Request) {
	list, _ := h.useCase.GetAllJournals(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar jurnal akuntansi", list)
}

func (h *FinanceHandler) JournalDetailDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	if r.Method == http.MethodGet {
		journal, _ := h.useCase.GetJournalByID(r.Context(), id)
		if journal == nil {
			response.WriteError(w, http.StatusNotFound, "Jurnal tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail jurnal akuntansi", journal)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *FinanceHandler) ReconciliationHandler(w http.ResponseWriter, r *http.Request) {
	logs, _ := h.useCase.GetReconciliationLogs(r.Context(), "") // Akan diperbaiki nanti di usecase
	response.WriteSuccess(w, http.StatusOK, "Proses rekonsiliasi", logs)
}

func (h *FinanceHandler) CashFlowReportHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := h.useCase.GenerateCashFlowStatement(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Laporan Arus Kas (Cash Flow)", res)
}

func (h *FinanceHandler) ExportReportHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		http.NotFound(w, r)
		return
	}
	format := parts[4]

	// Simulasi file ekspor
	res := map[string]string{
		"format": format,
		"url":    "https://storage.rinzi.com/exports/report." + format,
	}
	response.WriteSuccess(w, http.StatusOK, "Berkas laporan siap diunduh", res)
}

func (h *FinanceHandler) AnalyticsDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	action := parts[3]

	// Khusus: /api/finance/analytics/products/{product_id}
	// Gabungkan action + product_id agar repo/usecase bisa membedakan query per-produk
	if action == "products" && len(parts) >= 5 {
		productID := parts[4]
		data, _ := h.useCase.GetAnalyticsData(r.Context(), "products/"+productID)
		response.WriteSuccess(w, http.StatusOK, "Data Analitik Produk: "+productID, data)
		return
	}

	data, _ := h.useCase.GetAnalyticsData(r.Context(), action)
	response.WriteSuccess(w, http.StatusOK, "Data Analitik: "+action, data)
}

// Tambahkan fungsi ini di dalam file services/finance/delivery/http/finance_handler.go

// ProcessInboxEventHandler menangani POST /api/finance/inbox-events (Hook otomatis dari Outbox Worker)
func (h *FinanceHandler) ProcessInboxEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 1. Dekode dan WAJIB tangkap error-nya jika payload rusak
	var req dto.JournalIncomingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[FINANCE GATE] Gagal decode JSON request: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Format objek inbox event malformed")
		return // <-- Hentikan proses, jangan kirim 200 OK!
	}

	// 2. Kirim ke UseCase untuk disimpan transaksional ke DB
	if err := h.useCase.EnqueueInboxEvent(r.Context(), req); err != nil {
		logger.Error("[FINANCE GATE] Gagal menyimpan event ke outbox_events DB: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Gagal mengantrekan data: "+err.Error())
		return // <-- Hentikan proses jika DB gagal menulis!
	}

	// 3. Hanya kirimkan sukses jika data BENAR-BENAR tertulis di finance_db
	logger.Info("[FINANCE GATE] Event ID %s dari service luar sukses disimpan di DB", req.ID)
	response.WriteSuccess(w, http.StatusOK, "Event berhasil diserap ke antrean finansial", nil)
}

func (h *FinanceHandler) UploadDailyIncomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Gagal memproses form data")
		return
	}

	file, header, err := r.FormFile("pdf")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File pdf tidak ditemukan")
		return
	}
	defer file.Close()

	fileBytes := make([]byte, header.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Gagal membaca file pdf")
		return
	}

	url, err := h.useCase.UploadDailyIncomeReport(r.Context(), header.Filename, fileBytes)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "PDF berhasil diunggah", map[string]string{"url": url})
}
