package routes

import (
	"net/http"
	"strings"

	delivery "bisnis-rinzi/services/inventory/delivery/http"
)

func RegisterInventoryRoutes(mux *http.ServeMux, handler *delivery.InventoryHandler) {
	// 0. Endpoint Publik (Akses Media Tanpa Token)
	mux.HandleFunc("/api/public/inventory/media/", handler.PublicMediaHandler)

	// 1. Endpoint Core Produk (Strict Path)
	mux.HandleFunc("/api/inventory/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ProductsQueryHandler(w, r)
		case http.MethodPost:
			handler.CreateProductHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Pindahkan rute statis spesifik ke atas mux utama agar tidak bertabrakan dengan ID dinamis
	mux.HandleFunc("/api/inventory/products/low-stock", handler.LowStockHandler)
	mux.HandleFunc("/api/inventory/products/search", handler.ProductSearchHandler) // Tambahan Baru
	mux.HandleFunc("/api/inventory/products/import", handler.ImportCSVHandler)
	mux.HandleFunc("/api/inventory/products/export", handler.ExportCSVHandler)

	// Tambahan Baru: Kesiapan Infrastruktur Sync PWA
	mux.HandleFunc("/api/inventory/sync/version", handler.SyncVersionHandler)
	mux.HandleFunc("/api/inventory/sync/changes", handler.SyncChangesHandler)
	mux.HandleFunc("/api/inventory/catalog/sync", handler.CatalogSyncHandler)

	// 2. JALUR DINAMIS PRODUCTS (Hanya untuk ID unik dan Sub-Resource miliknya)
	mux.HandleFunc("/api/inventory/products/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(strings.Trim(path, "/"), "/")

		// Amankan jalur barcode dan sku terlebih dahulu
		if len(parts) >= 5 && parts[3] == "barcode" {
			handler.BarcodeHandler(w, r)
			return
		}
		if len(parts) >= 5 && parts[3] == "sku" {
			handler.SKUHandler(w, r)
			return
		}

		if len(parts) < 4 {
			http.NotFound(w, r)
			return
		}

		// Jika murni /api/inventory/products/{id}
		if len(parts) == 4 {
			handler.ProductItemHandler(w, r)
			return
		}

		// Jika polanya /api/inventory/products/{id}/{sub_resource}
		subResource := parts[4]
		switch subResource {
		case "stock":
			handler.ProductStockHandler(w, r)
		case "stock-history":
			handler.ProductStockHistoryHandler(w, r)
		case "cost-histories":
			handler.ProductCostHistoriesHandler(w, r)
		case "media":
			// Jika polanya /api/inventory/products/{id}/media/{media_id} (Panjang parts == 6)
			if len(parts) == 6 {
				handler.ProductMediaItemHandler(w, r)
				return
			}
			// Jika polanya /api/inventory/products/{id}/media (Panjang parts == 5)
			switch r.Method {
			case http.MethodGet:
				handler.ProductMediaListHandler(w, r) // Panggil Handler List Baru
			case http.MethodPost:
				handler.UploadMediaHandler(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		default:
			http.NotFound(w, r)
		}
	})

	// 3. JALUR DINAMIS TERPADU: Kategori, Brand, Unit
	mux.HandleFunc("/api/inventory/categories/", handler.CategoryRouterDispatcher)
	mux.HandleFunc("/api/inventory/categories", handler.CategoryRouterDispatcher)
	mux.HandleFunc("/api/inventory/brands/", handler.BrandRouterDispatcher)
	mux.HandleFunc("/api/inventory/brands", handler.BrandRouterDispatcher)
	mux.HandleFunc("/api/inventory/units/", handler.UnitRouterDispatcher)
	mux.HandleFunc("/api/inventory/units", handler.UnitRouterDispatcher)

	// 4. Stocks & Movements Murni
	mux.HandleFunc("/api/inventory/stocks", handler.GlobalStocksHandler)
	mux.HandleFunc("/api/inventory/stocks/adjust", handler.AdjustStockHandler)
	mux.HandleFunc("/api/inventory/stocks/card/", handler.StockCardHandler)

	mux.HandleFunc("/api/inventory/stocks/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) == 5 && parts[4] == "thresholds" {
			handler.UpdateStockThresholdsHandler(w, r)
			return
		}

		// If it reaches here, we might need to handle other dynamic routes
		// but since adjust and card are already matched exactly above,
		// we just return 404 for unknown dynamic paths.
		http.NotFound(w, r)
	})
	mux.HandleFunc("/api/inventory/stock-movements", handler.StockMovementsHandler)
	mux.HandleFunc("/api/inventory/stock-movements/", handler.StockMovementsByIDHandler)

	// 5. Internal Endpoints (Service to Service)
	mux.HandleFunc("/internal/inventory/stock/deduct", handler.InternalDeductStockHandler)
	mux.HandleFunc("/internal/inventory/stock/restore", handler.InternalRestoreStockHandler)
}
