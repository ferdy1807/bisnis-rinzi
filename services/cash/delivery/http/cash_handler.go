package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"bisnis-rinzi/packages/backend/response"
	"bisnis-rinzi/services/cash/dto"
	"bisnis-rinzi/services/cash/usecase"
)

type CashHandler struct {
	useCase usecase.CashUseCase
}

func NewCashHandler(uc usecase.CashUseCase) *CashHandler {
	return &CashHandler{useCase: uc}
}

func (h *CashHandler) OpenSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP ditolak")
		return
	}
	cashierID := r.Header.Get("X-User-Id")
	var req dto.OpenSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
		return
	}
	sessionID, err := h.useCase.OpenSession(r.Context(), cashierID, req)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusCreated, "Sesi shift kasir berhasil dibuka", map[string]string{"cashier_session_id": sessionID})
}

func (h *CashHandler) CloseSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP wajib POST")
		return
	}
	cashierID := r.Header.Get("X-User-Id")
	var req dto.CloseSessionRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.useCase.CloseSession(r.Context(), cashierID, req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Sesi shift kasir berhasil ditutup & data ter-outbox aman", nil)
}

func (h *CashHandler) CurrentSessionHandler(w http.ResponseWriter, r *http.Request) {
	cashierID := r.Header.Get("X-User-Id")
	session, _ := h.useCase.GetCurrentSession(r.Context(), cashierID)
	if session == nil {
		response.WriteSuccess(w, http.StatusOK, "Anda tidak memiliki shift aktif saat ini", nil)
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Data status sesi kasir berjalan", session)
}

func (h *CashHandler) ShiftsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		list, _ := h.useCase.GetAllShifts(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar riwayat shift kasir", list)
	} else {
		response.WriteError(w, http.StatusMethodNotAllowed, "Hanya metode GET yang didukung")
	}
}

func (h *CashHandler) ShiftItemHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]
	session, _ := h.useCase.GetShiftByID(r.Context(), id)
	if session == nil {
		response.WriteError(w, http.StatusNotFound, "Data sesi kasir tidak ditemukan")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Detail sesi kasir", session)
}

func (h *CashHandler) ShiftSummaryHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]
	summary, err := h.useCase.GetShiftSummary(r.Context(), id)
	if err != nil || summary == nil {
		response.WriteError(w, http.StatusNotFound, "Gagal mendapatkan ringkasan sesi kasir")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Ringkasan sesi kasir", summary)
}

func (h *CashHandler) UploadShiftReportHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		response.WriteError(w, http.StatusBadRequest, "Gagal memproses form file: "+err.Error())
		return
	}

	file, header, err := r.FormFile("report")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File laporan (report) wajib diunggah")
		return
	}
	defer file.Close()

	session, err := h.useCase.UploadShiftReport(r.Context(), id, file, header.Size, header.Filename)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Laporan shift berhasil diunggah", session)
}

func (h *CashHandler) TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, _ := h.useCase.GetAllTransactions(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Histori mutasi laci kasir", list)
	case http.MethodPost:
		cashierID := r.Header.Get("X-User-Id")
		var req dto.CashTransactionRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.CreateCashTransaction(r.Context(), cashierID, req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Mutasi uang laci kasir berhasil dicatat", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CashHandler) TransactionItemHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]
	trx, _ := h.useCase.GetTransactionByID(r.Context(), id)
	if trx == nil {
		response.WriteError(w, http.StatusNotFound, "Data mutasi kas tidak ditemukan")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Detail mutasi kas", trx)
}

func (h *CashHandler) ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, _ := h.useCase.GetAllExpenses(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar seluruh pengeluaran operasional toko", list)
	case http.MethodPost:
		cashierID := r.Header.Get("X-User-Id")
		var req dto.CreateExpenseRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.AddExpense(r.Context(), cashierID, req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Biaya pengeluaran kas berhasil dibayarkan", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CashHandler) ExpenseItemHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		exp, _ := h.useCase.GetExpenseByID(r.Context(), id)
		if exp == nil {
			response.WriteError(w, http.StatusNotFound, "Data biaya tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail pengeluaran kas", exp)
	case http.MethodPut:
		var req dto.CreateExpenseRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.UpdateExpense(r.Context(), id, req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Pengeluaran operasional berhasil diubah", nil)
	case http.MethodDelete:
		if err := h.useCase.DeleteExpense(r.Context(), id); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Pengeluaran operasional berhasil dihapus", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CashHandler) ExpenseCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cats, _ := h.useCase.GetCategories(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar kategori biaya", cats)
	case http.MethodPost:
		var req dto.CreateExpenseCategoryRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = h.useCase.CreateCategory(r.Context(), req)
		response.WriteSuccess(w, http.StatusCreated, "Kategori pengeluaran baru berhasil ditambahkan", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *CashHandler) ExpenseCategoryIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		cat, _ := h.useCase.GetCategoryByID(r.Context(), id)
		if cat == nil {
			response.WriteError(w, http.StatusNotFound, "Kategori tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail kategori biaya", cat)
	case http.MethodPut:
		var req dto.CreateExpenseCategoryRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = h.useCase.UpdateCategory(r.Context(), id, req)
		response.WriteSuccess(w, http.StatusOK, "Kategori biaya berhasil diubah", nil)
	case http.MethodDelete:
		_ = h.useCase.DeleteCategory(r.Context(), id)
		response.WriteSuccess(w, http.StatusOK, "Kategori biaya berhasil dihapus permanen", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// InternalIncomeHandler menangani POST /internal/cash/income dari POS Outbox
func (h *CashHandler) InternalIncomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req dto.InternalIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Format payload tidak valid")
		return
	}

	if err := h.useCase.RecordInternalIncome(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Berhasil mencatat kas masuk tambahan", nil)
}

func (h *CashHandler) InternalIncomesQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	list, err := h.useCase.GetInternalIncomes(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Daftar pendapatan internal toko", list)
}
