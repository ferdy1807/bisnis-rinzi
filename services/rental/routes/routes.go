package routes

import (
	delivery "bisnis-rinzi/services/rental/delivery/http"
	"net/http"
	"strings"
)

func RegisterRentalRoutes(mux *http.ServeMux, handler *delivery.RentalHandler) {
	// 1. Pengelolaan Master Katalog Rental
	mux.HandleFunc("/api/rental/categories", handler.CategoriesHandler)
	mux.HandleFunc("/api/rental/categories/", handler.CategoryItemHandler)
	mux.HandleFunc("/api/rental/availability", handler.AvailabilityHandler)

	mux.HandleFunc("/api/rental/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ProductsQueryHandler(w, r)
		case http.MethodPost:
			handler.ProductsHandler(w, r)
		}
	})

	// Sub-Resource Media Gambar Terikat Parameter ID Produk & Detail Produk
	mux.HandleFunc("/api/rental/products/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) == 5 && parts[4] == "media" {
			switch r.Method {
			case http.MethodPost:
				handler.UploadProductMediaHandler(w, r)
			case http.MethodDelete:
				handler.DeleteProductMediaHandler(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		
		if len(parts) == 4 { // /api/rental/products/{id}
			handler.ProductItemHandler(w, r)
			return
		} else if len(parts) == 5 && parts[4] == "calendar" { // /api/rental/products/{id}/calendar
			handler.ProductCalendarHandler(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// 2. Alur Transaksional Reservasi Operasional Toko Rental
	mux.HandleFunc("/api/rental/reservations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.HistoryReservationsHandler(w, r)
		case http.MethodPost:
			handler.ReservationsHandler(w, r)
		}
	})

	// Jalur dinamis penanganan berkas spesifik ID Reservasi, Pembatalan, Penjemputan unit, Laporan Status
	mux.HandleFunc("/api/rental/reservations/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(strings.Trim(path, "/"), "/")

		if len(parts) >= 4 {
			idOrStatus := parts[3]
			if r.Method == http.MethodGet && (idOrStatus == "active" || idOrStatus == "upcoming" || idOrStatus == "overdue") {
				handler.ReservationStatusQueryHandler(w, r)
				return
			}
			if r.Method == http.MethodGet && idOrStatus == "calendar" {
				handler.ReservationsCalendarHandler(w, r)
				return
			}
		}

		if len(parts) >= 5 {
			subAction := parts[4]
			if subAction == "pickup" {
				switch r.Method {
				case http.MethodPost:
					handler.PickupHandler(w, r)
				case http.MethodGet:
					handler.GetPickupReportHandler(w, r)
				}
				return
			}
			if r.Method == http.MethodPost && subAction == "undo_pickup" {
				handler.UndoPickupHandler(w, r)
				return
			}
			if r.Method == http.MethodPost && subAction == "ready" {
				handler.ReadyForPickupHandler(w, r)
				return
			}
			if r.Method == http.MethodPost && subAction == "undo_ready" {
				handler.UndoReadyForPickupHandler(w, r)
				return
			}
			if r.Method == http.MethodPost && subAction == "cancel" {
				handler.CancelHandler(w, r)
				return
			}
		}
		handler.DynamicRouteDispatcher(w, r)
	})

	// 3. Alur Transaksional Dokumen Pengembalian & Audit Log Rekap Denda (Returns)
	mux.HandleFunc("/api/rental/returns", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ReturnsQueryHandler(w, r)
		case http.MethodPost:
			handler.ReturnsHandler(w, r)
		}
	})
	mux.HandleFunc("/api/rental/returns/", handler.ReturnDetailDispatcher)
}
