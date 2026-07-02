package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/packages/backend/response"
	"bisnis-rinzi/services/inventory/dto"
	"bisnis-rinzi/services/inventory/entity"
	"bisnis-rinzi/services/inventory/usecase"
)

type InventoryHandler struct {
	useCase usecase.InventoryUseCase
}

func NewInventoryHandler(uc usecase.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{useCase: uc}
}

// CreateProductHandler menangani POST /api/inventory/products [cite: 24]
func (h *InventoryHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req dto.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
		return
	}

	productID, err := h.useCase.CreateProduct(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Produk baru dan stok awal berhasil disimpan", map[string]string{"id": productID})
}

// AdjustStockHandler menangani POST /api/inventory/stocks/adjust [cite: 4]
func (h *InventoryHandler) AdjustStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req dto.StockAdjustRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
		return
	}

	if err := h.useCase.AdjustStock(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Penyesuaian stok (opname) berhasil dicatat", nil)
}

// UploadMediaHandler menangani POST /api/inventory/products/{id}/media [cite: 4]
func (h *InventoryHandler) UploadMediaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	// Ambil ID produk dari URL Path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "ID Produk tidak valid")
		return
	}
	productID := pathParts[4]

	// Batasi ukuran upload gambar maks 5MB
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Ukuran file terlalu besar (Maks 5MB)")
		return
	}

	file, header, err := r.FormFile("media")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File media tidak ditemukan dalam form data")
		return
	}
	defer file.Close()

	err = h.useCase.UploadProductMedia(r.Context(), productID, "IMAGE", file, header.Filename, header.Header.Get("Content-Type"), header.Size)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Media gambar produk berhasil diunggah ke MinIO", nil)
}

// SyncChangesHandler menangani GET /api/inventory/sync/changes
func (h *InventoryHandler) SyncChangesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	// Membaca versi terakhir yang dimiliki oleh PWA Client dari query param
	versionStr := r.URL.Query().Get("version")
	clientVersion, _ := strconv.ParseInt(versionStr, 10, 64)

	res, err := h.useCase.GetSyncData(r.Context(), clientVersion)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Berhasil mengambil delta data sinkronisasi PWA", res)
}

// ImportCSVHandler menangani POST /api/inventory/products/import [cite: 3]
func (h *InventoryHandler) ImportCSVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP harus POST")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "File CSV tidak ditemukan")
		return
	}
	defer file.Close()

	if err := h.useCase.ImportProductsFromCSV(r.Context(), file); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Import massal data produk via CSV berhasil", nil)
}

// ExportCSVHandler menangani GET /api/inventory/products/export [cite: 3]
func (h *InventoryHandler) ExportCSVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=katalog_produk_retail.csv")

	if err := h.useCase.ExportProductsToCSV(r.Context(), w); err != nil {
		logger.Error("Gagal ekspor data CSV: %v", err)
	}
}

// ProductsQueryHandler menangani GET /api/inventory/products dan pencarian custom
func (h *InventoryHandler) ProductsQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	products, _ := h.useCase.GetAllProducts(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar produk retail", products)
}

func (h *InventoryHandler) LowStockHandler(w http.ResponseWriter, r *http.Request) {
	products, _ := h.useCase.GetLowStockProducts(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar produk low-stock", products)
}

// --- CATEGORIES HANDLER ---
func (h *InventoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c, _ := h.useCase.GetAllCategories(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar kategori", c)
	case http.MethodPost:
		var b map[string]string
		json.NewDecoder(r.Body).Decode(&b)
		h.useCase.CreateCategory(r.Context(), b["code"], b["name"])
		response.WriteSuccess(w, http.StatusCreated, "Kategori berhasil dibuat", nil)
	}
}

// --- BRANDS HANDLER ---
func (h *InventoryHandler) BrandsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, _ := h.useCase.GetAllBrands(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar brand", b)
	case http.MethodPost:
		var b map[string]string
		json.NewDecoder(r.Body).Decode(&b)
		h.useCase.CreateBrand(r.Context(), b["code"], b["name"])
		response.WriteSuccess(w, http.StatusCreated, "Brand berhasil dibuat", nil)
	}
}

func (h *InventoryHandler) BrandIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		response.WriteError(w, http.StatusBadRequest, "ID Brand tidak valid")
		return
	}
	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		b, _ := h.useCase.GetBrandByID(r.Context(), id)
		if b == nil {
			response.WriteError(w, http.StatusNotFound, "Brand tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail data brand", b)

	case http.MethodPut:
		var b map[string]string
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
			return
		}
		err := h.useCase.UpdateBrand(r.Context(), id, b["code"], b["name"])
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Brand berhasil diperbarui", nil)

	case http.MethodDelete:
		err := h.useCase.DeleteBrand(r.Context(), id)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Brand berhasil dihapus (soft-delete)", nil)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) BrandRouterDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Jika panjang parts == 3, polanya adalah /api/inventory/brands (Koleksi / Buat Baru)
	if len(parts) == 3 {
		switch r.Method {
		case http.MethodGet:
			b, _ := h.useCase.GetAllBrands(r.Context())
			response.WriteSuccess(w, http.StatusOK, "Daftar brand berhasil diambil", b)
		case http.MethodPost:
			var b map[string]string
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
				return
			}
			err := h.useCase.CreateBrand(r.Context(), b["code"], b["name"])
			if err != nil {
				response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusCreated, "Brand berhasil dibuat", nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	// Jika panjang parts == 4, polanya adalah /api/inventory/brands/{id} (Detail / Update / Delete)
	if len(parts) == 4 {
		id := parts[3]
		switch r.Method {
		case http.MethodGet:
			b, _ := h.useCase.GetBrandByID(r.Context(), id)
			if b == nil {
				response.WriteError(w, http.StatusNotFound, "Brand tidak ditemukan")
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Detail data brand", b)
		case http.MethodPut:
			var b map[string]string
			json.NewDecoder(r.Body).Decode(&b)
			h.useCase.UpdateBrand(r.Context(), id, b["code"], b["name"])
			response.WriteSuccess(w, http.StatusOK, "Brand berhasil diperbarui", nil)
		case http.MethodDelete:
			h.useCase.DeleteBrand(r.Context(), id)
			response.WriteSuccess(w, http.StatusOK, "Brand berhasil dihapus", nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

// --- UNITS HANDLER ---
func (h *InventoryHandler) UnitsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u, _ := h.useCase.GetAllUnits(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar satuan unit", u)
	case http.MethodPost:
		var b map[string]string
		json.NewDecoder(r.Body).Decode(&b)
		h.useCase.CreateUnit(r.Context(), b["code"], b["name"])
		response.WriteSuccess(w, http.StatusCreated, "Satuan unit berhasil dibuat", nil)
	}
}

func (h *InventoryHandler) UnitIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		response.WriteError(w, http.StatusBadRequest, "ID Unit tidak valid")
		return
	}
	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		u, _ := h.useCase.GetUnitByID(r.Context(), id)
		if u == nil {
			response.WriteError(w, http.StatusNotFound, "Satuan unit tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail data satuan unit", u)

	case http.MethodDelete:
		err := h.useCase.DeleteUnit(r.Context(), id)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Satuan unit berhasil dihapus permanen", nil)

	case http.MethodPut:
		var b map[string]string
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
			return
		}
		err := h.useCase.UpdateUnit(r.Context(), id, b["code"], b["name"])
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Satuan unit berhasil diperbarui", nil)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) ProductItemHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		p, _ := h.useCase.GetProductByID(r.Context(), id)
		if p == nil {
			response.WriteError(w, http.StatusNotFound, "Produk tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail data produk retail", p)
	case http.MethodPut:
		var p entity.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
			return
		}
		if err := h.useCase.UpdateProduct(r.Context(), id, &p); err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Data produk berhasil diubah", nil)
	case http.MethodDelete:
		h.useCase.DeleteProduct(r.Context(), id)
		response.WriteSuccess(w, http.StatusOK, "Produk berhasil dihapus (soft-delete)", nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handler Baru: GET /api/inventory/products/{id}/stock
func (h *InventoryHandler) ProductStockHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := parts[3] // id produk
	st, _ := h.useCase.GetStock(r.Context(), id)
	response.WriteSuccess(w, http.StatusOK, "Informasi stok produk", st)
}

// GET /api/inventory/products/{id}/stock-history
// =====================================================================================
func (h *InventoryHandler) ProductStockHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	productID := parts[3] // Mengambil segmen {id}

	// Memanggil UseCase untuk mengambil riwayat pergerakan stok khusus produk ini
	history, err := h.useCase.GetStockCard(r.Context(), productID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Histori mutasi stok produk berhasil ditarik", history)
}

// Penyempurnaan Logika: POST /api/inventory/products/{id}/cost-histories
// =====================================================================================
func (h *InventoryHandler) ProductCostHistoriesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	productID := parts[3]

	switch r.Method {
	case http.MethodGet:
		histories, err := h.useCase.GetCostHistories(r.Context(), productID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Data riwayat modal HPP produk berhasil ditarik", histories)

	case http.MethodPost:
		var req struct {
			AverageCost float64 `json:"average_cost"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
			return
		}

		if req.AverageCost <= 0 {
			response.WriteError(w, http.StatusUnprocessableEntity, "Nilai harga modal pokok (HPP) tidak boleh kurang dari atau sama dengan nol")
			return
		}

		if err := h.useCase.AddCostHistory(r.Context(), productID, req.AverageCost, ""); err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		logger.Info("Mengunci HPP Baru untuk produk %s senilai Rp%.2f via %s", productID, req.AverageCost, "SISTEM")
		response.WriteSuccess(w, http.StatusCreated, "Riwayat harga modal pokok (HPP) baru berhasil dikunci ke sistem", nil)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// PublicMediaHandler melayani streaming file secara publik tanpa autentikasi (untuk tag <img>)
func (h *InventoryHandler) PublicMediaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "Media ID tidak valid")
		return
	}
	mediaID := parts[4] // /api/public/inventory/media/{media_id}

	stream, media, err := h.useCase.GetMediaStream(r.Context(), mediaID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", media.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(media.FileSizeValues, 10))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache 1 tahun untuk gambar
	w.WriteHeader(http.StatusOK)
	io.Copy(w, stream)
}

// GET & DELETE /api/inventory/products/{id}/media/{media_id}
// =====================================================================================
func (h *InventoryHandler) ProductMediaItemHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	mediaID := parts[5] // Mengambil segmen {media_id}

	switch r.Method {
	case http.MethodGet:
		action := r.URL.Query().Get("action")
		if action == "view" {
			stream, media, err := h.useCase.GetMediaStream(r.Context(), mediaID)
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			defer stream.Close()

			w.Header().Set("Content-Type", media.MimeType)
			w.Header().Set("Content-Length", strconv.FormatInt(media.FileSizeValues, 10))
			w.WriteHeader(http.StatusOK)
			io.Copy(w, stream)
			return
		}

		media, err := h.useCase.GetMediaItem(r.Context(), mediaID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if media == nil {
			response.WriteError(w, http.StatusNotFound, "Data media tidak ditemukan")
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Detail data berkas gambar produk", media)

	case http.MethodDelete:
		err := h.useCase.DeleteProductMedia(r.Context(), mediaID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Berkas gambar produk berhasil dihapus permanen dari Object Storage MinIO dan Database", nil)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Tambahan Baru: Menangani GET /api/inventory/products/{id}/media
// =====================================================================================
func (h *InventoryHandler) ProductMediaListHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	productID := parts[3]

	mediaList, err := h.useCase.GetProductMedia(r.Context(), productID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Seluruh daftar berkas gambar produk berhasil ditarik", mediaList)
}

// Handler Baru: GET /api/inventory/categories/{id} (Koreksi Respon Kosong)
func (h *InventoryHandler) CategoryIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 4 {
		id := parts[3]
		switch r.Method {
		case http.MethodGet:
			c, _ := h.useCase.GetCategoryByID(r.Context(), id)
			if c == nil {
				response.WriteError(w, http.StatusNotFound, "Kategori tidak ditemukan")
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Detail data kategori", c)
		case http.MethodPut:
			var b map[string]string
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
				return
			}
			err := h.useCase.UpdateCategory(r.Context(), id, b["code"], b["name"])
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Kategori berhasil diperbarui", nil)
		case http.MethodDelete:
			err := h.useCase.DeleteCategory(r.Context(), id)
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Kategori berhasil dihapus (soft-delete)", nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}
}

func (h *InventoryHandler) CategoryRouterDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) == 3 {
		switch r.Method {
		case http.MethodGet:
			c, _ := h.useCase.GetAllCategories(r.Context())
			response.WriteSuccess(w, http.StatusOK, "Daftar kategori", c)
		case http.MethodPost:
			var b map[string]string
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
				return
			}
			if err := h.useCase.CreateCategory(r.Context(), b["code"], b["name"]); err != nil {
				response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusCreated, "Kategori berhasil dibuat", nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 4 {
		id := parts[3]
		switch r.Method {
		case http.MethodGet:
			c, _ := h.useCase.GetCategoryByID(r.Context(), id)
			if c == nil {
				response.WriteError(w, http.StatusNotFound, "Kategori tidak ditemukan")
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Detail data kategori", c)
		case http.MethodPut:
			var b map[string]string
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
				return
			}
			if err := h.useCase.UpdateCategory(r.Context(), id, b["code"], b["name"]); err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Kategori diperbarui", nil)
		case http.MethodDelete:
			if err := h.useCase.DeleteCategory(r.Context(), id); err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Kategori dihapus", nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

// UnitRouterDispatcher mengatur rute dinamis terpadu untuk Unit Satuan
func (h *InventoryHandler) UnitRouterDispatcher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) == 3 {
		switch r.Method {
		case http.MethodGet:
			u, _ := h.useCase.GetAllUnits(r.Context())
			response.WriteSuccess(w, http.StatusOK, "Daftar satuan unit", u)
		case http.MethodPost:
			var b map[string]string
			json.NewDecoder(r.Body).Decode(&b)
			h.useCase.CreateUnit(r.Context(), b["code"], b["name"])
			response.WriteSuccess(w, http.StatusCreated, "Satuan unit berhasil dibuat", nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 4 {
		id := parts[3]
		switch r.Method {
		case http.MethodGet:
			u, _ := h.useCase.GetUnitByID(r.Context(), id)
			response.WriteSuccess(w, http.StatusOK, "Detail data unit satuan", u)
		case http.MethodPut:
			var b map[string]string
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
				return
			}
			err := h.useCase.UpdateUnit(r.Context(), id, b["code"], b["name"])
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteSuccess(w, http.StatusOK, "Satuan unit berhasil diperbarui", nil)
		case http.MethodDelete:
			h.useCase.DeleteUnit(r.Context(), id)
			response.WriteSuccess(w, http.StatusOK, "Satuan unit berhasil dihapus", nil)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

// Tambahan Baru: Menangani GET /api/inventory/products/search
func (h *InventoryHandler) ProductSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	queryParam := r.URL.Query().Get("query")
	if queryParam == "" {
		queryParam = r.URL.Query().Get("q")
	}

	products, err := h.useCase.SearchProducts(r.Context(), queryParam)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Hasil pencarian produk retail", products)
}

// Tambahan Baru: Menangani GET /api/inventory/sync/version
func (h *InventoryHandler) SyncVersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	currentVersion, err := h.useCase.GetLatestSyncVersion(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Berhasil mengambil versi sinkronisasi global server terbaru", map[string]int64{
		"latest_version_number": currentVersion,
	})
}

// KOREKSI PARSING: Sesuaikan pengambilan segmen karena rute /products/barcode/{code} membuat kode bergeser ke parts[4]
func (h *InventoryHandler) BarcodeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "Barcode tidak dicantumkan")
		return
	}
	code := parts[4] // Perbaikan dari parts[5] ke parts[4]
	p, _ := h.useCase.GetProductByBarcode(r.Context(), code)
	if p == nil {
		response.WriteError(w, http.StatusNotFound, "Produk dengan barcode tersebut tidak ditemukan")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Data barcode produk", p)
}

// KOREKSI PARSING: Sesuaikan pengambilan segmen untuk SKU
func (h *InventoryHandler) SKUHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "SKU tidak dicantumkan")
		return
	}
	sku := parts[4] // Perbaikan dari parts[5] ke parts[4]
	p, _ := h.useCase.GetProductBySKU(r.Context(), sku)
	if p == nil {
		response.WriteError(w, http.StatusNotFound, "Produk dengan SKU tersebut tidak ditemukan")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Data SKU produk", p)
}

func (h *InventoryHandler) CatalogSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	data, err := h.useCase.GetCatalogSyncData(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Seluruh katalog master berhasil ditarik untuk PWA Offline caching", data)
}

func (h *InventoryHandler) GlobalStocksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	stocks, err := h.useCase.GetAllStocksData(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Daftar seluruh stok produk retail", stocks)
}

func (h *InventoryHandler) UpdateStockThresholdsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "Product ID wajib dicantumkan")
		return
	}
	productID := parts[3]

	var req dto.StockThresholdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload tidak valid")
		return
	}

	if err := h.useCase.UpdateStockThresholds(r.Context(), productID, req); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Berhasil memperbarui batas stok minimum dan safety stock", nil)
}

func (h *InventoryHandler) StockCardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "Product ID wajib dicantumkan")
		return
	}
	productID := parts[4] // segmen {product_id}

	card, err := h.useCase.GetStockCard(r.Context(), productID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Kartu kendali stok produk berhasil dipetakan", card)
}

// GET /api/inventory/stock-movements
func (h *InventoryHandler) StockMovementsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	movements, err := h.useCase.GetAllMovements(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Daftar mutasi pergerakan barang global", movements)
}

// GET /api/inventory/stock-movements/{id}
func (h *InventoryHandler) StockMovementsByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		response.WriteError(w, http.StatusBadRequest, "Movement ID tidak valid")
		return
	}
	id := parts[3] // mengambil segmen {id}

	m, err := h.useCase.GetMovementByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Data mutasi pergerakan stok tidak ditemukan")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Detail mutasi dokumen log pergerakan stok", m)
}

// POST /internal/inventory/stock/deduct
func (h *InventoryHandler) InternalDeductStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req dto.InternalStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
		return
	}

	if err := h.useCase.InternalDeductStock(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Pengurangan stok internal berhasil dicatat", nil)
}

// POST /internal/inventory/stock/restore
func (h *InventoryHandler) InternalRestoreStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req dto.InternalStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
		return
	}

	if err := h.useCase.InternalRestoreStock(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Pengembalian stok internal berhasil dicatat", nil)
}
