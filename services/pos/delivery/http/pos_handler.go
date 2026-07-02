package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"bisnis-rinzi/packages/backend/response"
	"bisnis-rinzi/services/pos/dto"
	"bisnis-rinzi/services/pos/usecase"
)

type POSHandler struct {
	useCase usecase.POSUseCase
}

func NewPOSHandler(uc usecase.POSUseCase) *POSHandler {
	return &POSHandler{useCase: uc}
}

// CheckoutHandler menangani POST /api/pos/sales
func (h *POSHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	// Ekstrak ID Kasir dari Header yang disuntikkan API Gateway
	cashierID := r.Header.Get("X-User-Id")

	var req dto.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Struktur data JSON kasir tidak valid")
		return
	}

	sale, err := h.useCase.Checkout(r.Context(), cashierID, req)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Transaksi checkout kasir berhasil dicatat", map[string]string{
		"id":             sale.ID,
		"invoice_number": sale.InvoiceNumber,
	})
}

// HistoryHandler menangani GET /api/pos/sales
func (h *POSHandler) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP salah")
		return
	}
	sales, err := h.useCase.GetSalesHistory(r.Context())
	if err != nil {
		println("Error GetSalesHistory:", err.Error())
	}
	response.WriteSuccess(w, http.StatusOK, "Daftar riwayat transaksi kasir", sales)
}

// InvoiceHandler menangani GET /api/pos/sales/invoice/{invoice_number}
func (h *POSHandler) InvoiceHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	invoiceNum := parts[5] // Mengambil segmen {invoice_number}

	sale, _ := h.useCase.GetInvoiceData(r.Context(), invoiceNum)
	response.WriteSuccess(w, http.StatusOK, "Data Invoice", sale)
}

// SyncHandler menangani POST /api/pos/sync
func (h *POSHandler) SyncHandler(w http.ResponseWriter, r *http.Request) {
	cashierID := r.Header.Get("X-User-Id")
	var req dto.OfflineSyncPayload
	_ = json.NewDecoder(r.Body).Decode(&req)

	_ = h.useCase.SyncOfflineTransactions(r.Context(), cashierID, req)
	response.WriteSuccess(w, http.StatusOK, "Proses sinkronisasi data transaksi offline selesai diproses", nil)
}

// RetrySyncHandler menangani POST /api/pos/sync/retry
func (h *POSHandler) RetrySyncHandler(w http.ResponseWriter, r *http.Request) {
	cashierID := r.Header.Get("X-User-Id")
	var req dto.OfflineSyncPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload sinkronisasi retry tidak valid")
		return
	}

	_ = h.useCase.RetryFailedSync(r.Context(), cashierID, req)
	response.WriteSuccess(w, http.StatusOK, "Proses percobaan ulang sinkronisasi selesai diproses", nil)
}

// DynamicRouteDispatcher mengurai rute dinamis terproteksi /api/pos/sales/{id}/*
func (h *POSHandler) DynamicRouteDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3] // Menangkap {id}

	if len(parts) == 4 {
		sale, _ := h.useCase.GetSaleDetail(r.Context(), id)
		response.WriteSuccess(w, http.StatusOK, "Detail data transaksi", sale)
		return
	}

	subResource := parts[4]
	if r.Method == http.MethodGet && subResource == "items" { // /api/pos/sales/{id}/items
		items, _ := h.useCase.GetSaleItems(r.Context(), id)
		response.WriteSuccess(w, http.StatusOK, "Daftar barang transaksi", items)
	} else if r.Method == http.MethodGet && subResource == "receipt" { // /api/pos/sales/{id}/receipt
		cashierName := r.Header.Get("X-User-Name")
		receipt, _ := h.useCase.GetReceipt(r.Context(), id, cashierName)
		response.WriteSuccess(w, http.StatusOK, "Struktur data cetak struk belanja kasir", receipt)
	} else {
		http.NotFound(w, r)
	}
}

func (h *POSHandler) SyncStatusHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "Status Jaringan Server POS: ONLINE", nil)
}

func (h *POSHandler) FailedSyncLogsHandler(w http.ResponseWriter, r *http.Request) {
	logs, _ := h.useCase.GetFailedSyncLogs(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar kegagalan sinkronisasi transaksi offline", logs)
}

func (h *POSHandler) TopProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP salah")
		return
	}
	
	sessionID := r.URL.Query().Get("session_id")
	var sid *string
	if sessionID != "" {
		sid = &sessionID
	}

	// Default limit 5
	tops, err := h.useCase.GetTopProducts(r.Context(), 5, sid)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Gagal mengambil data produk terlaris")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Data produk terlaris berhasil diambil", tops)
}

func (h *POSHandler) ProductSalesHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "ID Produk tidak valid")
		return
	}
	productID := parts[3] // /api/pos/products/{id}/sales

	histories, err := h.useCase.GetProductSalesHistory(r.Context(), productID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Riwayat penjualan produk berhasil ditarik", histories)
}

// UploadInvoiceHandler menangani POST /api/pos/sales/{id}/invoice
func (h *POSHandler) UploadInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	// Ekstrak ID dari path
	parts := strings.Split(r.URL.Path, "/")
	// /api/pos/sales/{id}/invoice
	// parts = ["", "api", "pos", "sales", "{id}", "invoice"]
	if len(parts) < 6 {
		response.WriteError(w, http.StatusBadRequest, "Format URL tidak valid")
		return
	}
	saleID := parts[4]

	// Parse multipart form dengan batas ukuran memori (misal: 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Gagal memproses file upload")
		return
	}

	// Dapatkan file dari form key "invoice"
	file, handler, err := r.FormFile("invoice")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File 'invoice' tidak ditemukan")
		return
	}
	defer file.Close()

	// Baca ke dalam byte slice
	fileBytes := make([]byte, handler.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Gagal membaca file upload")
		return
	}

	// Tentukan content type (PDF)
	contentType := handler.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/pdf"
	}

	// Dapatkan nama kasir dari Header
	cashierName := r.Header.Get("X-User-Name")
	if cashierName == "" {
		cashierName = "Unknown_Cashier"
	}

	// Panggil usecase
	url, err := h.useCase.UploadInvoice(r.Context(), saleID, cashierName, fileBytes, contentType)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Invoice berhasil diunggah", map[string]string{
		"invoice_url": url,
	})
}
