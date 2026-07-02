package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"bisnis-rinzi/packages/backend/response"
	"bisnis-rinzi/services/rental/dto"
	"bisnis-rinzi/services/rental/usecase"
)

type RentalHandler struct {
	useCase usecase.RentalUseCase
}

func NewRentalHandler(uc usecase.RentalUseCase) *RentalHandler {
	return &RentalHandler{useCase: uc}
}

func (h *RentalHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, _ := h.useCase.GetCategories(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar kategori objek sewa", list)
	case http.MethodPost:
		var req dto.CreateCategoryRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.CreateCategory(r.Context(), req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusCreated, "Kategori sewa berhasil ditambahkan", nil)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RentalHandler) ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := h.useCase.GetProducts(r.Context())
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Katalog produk sewa/rental retail", list)
	case http.MethodPost:
		var req dto.CreateProductRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.CreateProduct(r.Context(), req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusCreated, "Produk sewa baru berhasil didaftarkan", nil)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RentalHandler) ReservationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cashierID := r.Header.Get("X-User-Id")

	var req dto.CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Struktur berkas payload reservasi salah")
		return
	}

	resID, invoiceNum, err := h.useCase.CreateReservation(r.Context(), cashierID, req)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Reservasi sewa berhasil disimpan", map[string]string{
		"id":             resID,
		"invoice_number": invoiceNum,
	})
}

func (h *RentalHandler) ReturnsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cashierID := r.Header.Get("X-User-Id")

	var req dto.ProcessReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload pengembalian unit sewa tidak valid")
		return
	}

	ret, err := h.useCase.ProcessRentalReturn(r.Context(), cashierID, req)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Proses pengembalian unit sewa dan kliring dana jaminan selesai", ret)
}

func (h *RentalHandler) DynamicRouteDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	if len(parts) == 4 { // /api/rental/reservations/{id}
		res, _ := h.useCase.GetReservationDetail(r.Context(), id)
		if res != nil {
			items, _ := h.useCase.GetReservationItemsDetail(r.Context(), id)
			res.Items = items
		}
		response.WriteSuccess(w, http.StatusOK, "Detail berkas transaksi sewa", res)
		return
	}

	sub := parts[4]
	switch sub {
	case "items": // /api/rental/reservations/{id}/items
		items, _ := h.useCase.GetReservationItemsDetail(r.Context(), id)
		response.WriteSuccess(w, http.StatusOK, "Detail item daftar barang sewa", items)
	case "return": // /api/rental/reservations/{id}/return
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ret, err := h.useCase.GetReturnByReservationID(r.Context(), id)
		if err != nil || ret == nil {
			response.WriteError(w, http.StatusNotFound, "Dokumen pengembalian tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail dokumen pengembalian sewa", ret)
	case "contents_received":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req dto.ReservationContentPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, "Payload salah")
			return
		}
		if err := h.useCase.SaveDepositItems(r.Context(), id, req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Berhasil mencatat barang", nil)
	case "invoice": // /api/rental/reservations/{id}/invoice
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.UploadReservationInvoiceHandler(w, r, id)
	default:
		http.NotFound(w, r)
	}
}

func (h *RentalHandler) ProductsQueryHandler(w http.ResponseWriter, r *http.Request) {
	list, _ := h.useCase.GetProducts(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Katalog produk sewa/rental retail", list)
}

func (h *RentalHandler) UploadProductMediaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	productID := parts[3]

	// Batasi ukuran gambar aset rental maksimal 5MB
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Ukuran file foto aset sewa terlalu besar (Maks 5MB)")
		return
	}

	file, header, err := r.FormFile("media")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File gambar sewa tidak ditemukan dalam form data")
		return
	}
	defer file.Close()

	err = h.useCase.UploadProductMedia(r.Context(), productID, file, header.Filename, header.Header.Get("Content-Type"), header.Size)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Foto aset kelengkapan produk rental sukses disimpan ke MinIO", nil)
}


// HistoryReservationsHandler menangani GET /api/rental/reservations
func (h *RentalHandler) HistoryReservationsHandler(w http.ResponseWriter, r *http.Request) {
	list, _ := h.useCase.GetAllReservations(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar seluruh riwayat pesanan sewa pelanggan", list)
}

// PickupHandler menangani POST /api/rental/reservations/{id}/pickup
func (h *RentalHandler) PickupHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[3]

	if err := h.useCase.PickupRentalItems(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Unit sewa berhasil diambil oleh pelanggan (Status: PICKED_UP)", nil)
}

// UndoPickupHandler menangani POST /api/rental/reservations/{id}/undo_pickup
func (h *RentalHandler) UndoPickupHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[3]
	
	if err := h.useCase.UndoPickupRentalItems(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Berhasil membatalkan pengeluaran unit sewa (Dikembalikan ke READY_FOR_PICKUP)", nil)
}

// ReadyForPickupHandler menangani POST /api/rental/reservations/{id}/ready
func (h *RentalHandler) ReadyForPickupHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[3]
	
	if err := h.useCase.MarkAsReadyForPickup(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Berhasil memindahkan status reservasi ke tahap siap diambil", nil)
}

// UndoReadyForPickupHandler menangani POST /api/rental/reservations/{id}/undo_ready
func (h *RentalHandler) UndoReadyForPickupHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[3]
	
	if err := h.useCase.UndoReadyForPickup(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Berhasil menarik kembali reservasi ke tahap antrean baru (BOOKED)", nil)
}

// CancelHandler menangani POST /api/rental/reservations/{id}/cancel
func (h *RentalHandler) CancelHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[3]

	var req dto.CancelReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Jika body kosong atau invalid, tetap inisialisasi dengan ID dari param
		req = dto.CancelReservationRequest{ReservationID: id, PenaltyFee: 0}
	} else {
		// Pastikan ID diset dari URL param meskipun body diisi
		req.ReservationID = id
	}

	if err := h.useCase.CancelReservation(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Pemesanan sewa dibatalkan secara transaksional (Status: CANCELLED)", nil)
}

// ReturnsQueryHandler menangani GET /api/rental/returns (Audit Jejak Denda Pengembalian)
func (h *RentalHandler) ReturnsQueryHandler(w http.ResponseWriter, r *http.Request) {
	list, err := h.useCase.GetAllReturns(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Riwayat log dokumen pengembalian sewa toko", list)
}

// ReturnDetailDispatcher mengurai rute /api/rental/returns/{id} dan sub-item denda miliknya
func (h *RentalHandler) ReturnDetailDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	if id == "damages" {
		if len(parts) == 4 {
			damages, err := h.useCase.GetDamagedItems(r.Context())
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Daftar kerusakan unit sewa", damages)
			return
		}
		if len(parts) == 6 && parts[5] == "settle" {
			damageID := parts[4]
			var req dto.SettleDamageRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				response.WriteError(w, http.StatusBadRequest, "Invalid payload")
				return
			}
			if err := h.useCase.SettleDamage(r.Context(), damageID, req); err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Status denda kerusakan berhasil diselesaikan", nil)
			return
		}
		http.NotFound(w, r)
		return
	}

	if len(parts) == 4 {
		ret, _ := h.useCase.GetReturnDetail(r.Context(), id)
		response.WriteSuccess(w, http.StatusOK, "Detail dokumen denda pengembalian sewa", ret)
		return
	}

	if len(parts) >= 5 {
		subAction := parts[4]
		if subAction == "items" {
			items, _ := h.useCase.GetReturnItems(r.Context(), id)
			response.WriteSuccess(w, http.StatusOK, "Rincian denda unit barang sewa", items)
			return
		}
		if subAction == "penalty" {
			h.ReturnPenaltyHandler(w, r)
			return
		}
		if subAction == "photos" {
			if len(parts) == 5 {
				if r.Method == http.MethodGet {
					h.GetReturnPhotosHandler(w, r)
					return
				}
				if r.Method == http.MethodPost {
					h.UploadReturnPhotoHandler(w, r)
					return
				}
			} else if len(parts) == 6 && r.Method == http.MethodDelete {
				h.DeleteReturnPhotoHandler(w, r)
				return
			}
		}
		if subAction == "receipt" && r.Method == http.MethodPost {
			h.UploadReturnReceiptHandler(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (h *RentalHandler) CategoryItemHandler(w http.ResponseWriter, r *http.Request) {
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
		response.WriteSuccess(w, http.StatusOK, "Detail kategori sewa", cat)
	case http.MethodPut:
		var req dto.CreateCategoryRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.UpdateCategory(r.Context(), id, req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Kategori sewa berhasil diubah", nil)
	case http.MethodDelete:
		if err := h.useCase.DeleteCategory(r.Context(), id); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Kategori sewa berhasil dihapus", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *RentalHandler) UploadReturnReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		response.WriteError(w, http.StatusBadRequest, "ID return tidak valid")
		return
	}
	returnID := parts[3]

	if err := r.ParseMultipartForm(5 << 20); err != nil { // Max 5 MB
		response.WriteError(w, http.StatusBadRequest, "Ukuran file invoice terlalu besar (Maks 5MB)")
		return
	}

	file, header, err := r.FormFile("receipt")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File receipt/invoice tidak ditemukan dalam form data")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/pdf"
	}

	if err := h.useCase.UploadReturnReceipt(r.Context(), returnID, file, header.Filename, contentType, header.Size); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Invoice/Receipt berhasil diunggah", nil)
}

func (h *RentalHandler) ProductItemHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		prod, _ := h.useCase.GetProductByID(r.Context(), id)
		if prod == nil {
			response.WriteError(w, http.StatusNotFound, "Produk sewa tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail produk sewa", prod)
	case http.MethodPut:
		var req dto.CreateProductRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := h.useCase.UpdateProduct(r.Context(), id, req); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Produk sewa berhasil diubah", nil)
	case http.MethodDelete:
		if err := h.useCase.DeleteProduct(r.Context(), id); err != nil {
			response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Produk sewa berhasil dihapus", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *RentalHandler) ReservationStatusQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	status := parts[3] // active, upcoming, overdue

	var list interface{}
	var msg string

	switch status {
	case "active":
		list, _ = h.useCase.GetActiveReservations(r.Context())
		msg = "Daftar reservasi aktif berjalan"
	case "upcoming":
		list, _ = h.useCase.GetUpcomingReservations(r.Context())
		msg = "Daftar reservasi akan datang"
	case "overdue":
		list, _ = h.useCase.GetOverdueReservations(r.Context())
		msg = "Daftar reservasi terlambat dikembalikan"
	default:
		http.NotFound(w, r)
		return
	}

	response.WriteSuccess(w, http.StatusOK, msg, list)
}

func (h *RentalHandler) AvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	productID := r.URL.Query().Get("product_id")
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")
	qtyStr := r.URL.Query().Get("qty")

	if productID == "" || startStr == "" || endStr == "" {
		response.WriteError(w, http.StatusBadRequest, "product_id, start_date, end_date wajib diisi")
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Format start_date harus RFC3339")
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Format end_date harus RFC3339")
		return
	}

	qty := 1
	if qtyStr != "" {
		fmt.Sscanf(qtyStr, "%d", &qty)
	}

	req := dto.CheckAvailabilityRequest{
		RentalProductID: productID,
		StartDate:       start,
		EndDate:         end,
		QtyRequested:    qty,
	}

	isAvailable, err := h.useCase.CheckStockAvailability(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resData := map[string]interface{}{
		"is_available": isAvailable,
	}
	response.WriteSuccess(w, http.StatusOK, "Status ketersediaan stok rental", resData)
}

func (h *RentalHandler) ProductCalendarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	start := time.Now()
	end := start.AddDate(0, 1, 0) // Default 1 bulan

	if startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			start = t
		}
	}
	if endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			end = t
		}
	}

	list, err := h.useCase.GetProductCalendar(r.Context(), id, start, end)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Kalender reservasi produk", list)
}

func (h *RentalHandler) ReservationsCalendarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	start := time.Now()
	end := start.AddDate(0, 1, 0) // Default 1 bulan

	if startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			start = t
		}
	}
	if endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			end = t
		}
	}

	list, err := h.useCase.GetReservationsCalendar(r.Context(), start, end)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Kalender seluruh reservasi sewa", list)
}

func (h *RentalHandler) DeleteProductMediaHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	productID := parts[3]

	if err := h.useCase.DeleteProductPhoto(r.Context(), productID); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Berhasil menghapus media gambar dari katalog", nil)
}

func (h *RentalHandler) ReturnPenaltyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	summary, err := h.useCase.GetReturnPenaltySummary(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Rincian denda pengembalian unit sewa", summary)
}

func (h *RentalHandler) UploadReturnPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	returnID := parts[3]

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Gagal memproses form data: "+err.Error())
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File foto tidak ditemukan dalam form data")
		return
	}
	defer file.Close()

	err = h.useCase.UploadReturnPhoto(r.Context(), returnID, file, header.Filename, header.Header.Get("Content-Type"), header.Size)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Foto kondisi pengembalian berhasil diunggah", nil)
}

func (h *RentalHandler) GetReturnPhotosHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	returnID := parts[3]

	list, err := h.useCase.GetReturnPhotos(r.Context(), returnID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Daftar foto inspeksi barang saat dikembalikan", list)
}

func (h *RentalHandler) DeleteReturnPhotoHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 6 {
		http.NotFound(w, r)
		return
	}
	photoID := parts[5]

	if err := h.useCase.DeleteReturnPhoto(r.Context(), photoID); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Berhasil menghapus foto inspeksi pengembalian", nil)
}

func (h *RentalHandler) UploadReservationInvoiceHandler(w http.ResponseWriter, r *http.Request, resID string) {
	if resID == "" {
		response.WriteError(w, http.StatusBadRequest, "ID reservasi tidak valid")
		return
	}

	err := r.ParseMultipartForm(5 << 20) // Maksimal 5MB
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Ukuran file PDF terlalu besar, maksimal 5MB")
		return
	}

	file, header, err := r.FormFile("invoice")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File PDF tidak ditemukan dalam request (key: 'invoice')")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if contentType != "application/pdf" && ext != ".pdf" {
		response.WriteError(w, http.StatusBadRequest, "Format file tidak didukung, wajib berupa dokumen PDF")
		return
	}
    // Override if it's a PDF but browser didn't set content type correctly
    if ext == ".pdf" {
        contentType = "application/pdf"
    }

	if err := h.useCase.UploadReservationInvoice(r.Context(), resID, file, header.Filename, contentType, header.Size); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Sistem gagal memproses unggahan PDF dokumen penyewaan")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Berhasil mengunggah nota reservasi ke dalam sistem berkas", nil)
}

// GetPickupReportHandler menangani GET /api/rental/reservations/{id}/pickup
func (h *RentalHandler) GetPickupReportHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id := parts[3]

	res, err := h.useCase.GetReservationDetail(r.Context(), id)
	if err != nil || res == nil {
		response.WriteError(w, http.StatusNotFound, "Laporan penjemputan tidak ditemukan")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Laporan penjemputan unit sewa", res)
}


