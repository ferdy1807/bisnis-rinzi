package routes

import (
	delivery "bisnis-rinzi/services/pos/delivery/http"
	"net/http"
	"strings"
)

func RegisterPOSRoutes(mux *http.ServeMux, handler *delivery.POSHandler) {
	mux.HandleFunc("/api/pos/sales", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.HistoryHandler(w, r)
		} else if r.Method == http.MethodPost {
			handler.CheckoutHandler(w, r)
		}
	})

	mux.HandleFunc("/api/pos/sync", handler.SyncHandler)
	mux.HandleFunc("/api/pos/sync/retry", handler.RetrySyncHandler)
	mux.HandleFunc("/api/pos/sync/status", handler.SyncStatusHandler)
	mux.HandleFunc("/api/pos/sync/failed", handler.FailedSyncLogsHandler)
	mux.HandleFunc("/api/pos/reports/top-products", handler.TopProductsHandler)
	mux.HandleFunc("/api/pos/products/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/sales") {
			handler.ProductSalesHistoryHandler(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// Route dinamis dengan parameter prefix
	mux.HandleFunc("/api/pos/sales/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/invoice/") || strings.HasSuffix(r.URL.Path, "/invoice") {
			if r.Method == http.MethodPost {
				handler.UploadInvoiceHandler(w, r)
			} else {
				handler.InvoiceHandler(w, r)
			}
			return
		}
		handler.DynamicRouteDispatcher(w, r)
	})
}
