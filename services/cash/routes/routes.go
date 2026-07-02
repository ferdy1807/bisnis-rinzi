package routes

import (
	delivery "bisnis-rinzi/services/cash/delivery/http"
	"net/http"
	"strings"
)

func RegisterCashRoutes(mux *http.ServeMux, handler *delivery.CashHandler) {
	// 1. Tata Kelola Shift Kasir (Shift Management)
	mux.HandleFunc("/api/cash/shifts/open", handler.OpenSessionHandler)
	mux.HandleFunc("/api/cash/shifts/close", handler.CloseSessionHandler)
	mux.HandleFunc("/api/cash/shifts/current", handler.CurrentSessionHandler)
	mux.HandleFunc("/api/cash/shifts", handler.ShiftsHandler)

	// Penanganan Dinamis Detail Shift & Ringkasan Berita Acara
	mux.HandleFunc("/api/cash/shifts/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Pola Cocok: /api/cash/shifts/{id} -> len(parts) == 4 [id berada di parts[3]]
		if len(parts) == 4 {
			handler.ShiftItemHandler(w, r)
			return
		}

		// Pola Cocok: /api/cash/shifts/{id}/summary -> len(parts) == 5 dan parts[4] == "summary"
		if len(parts) == 5 && parts[4] == "summary" {
			handler.ShiftSummaryHandler(w, r)
			return
		}

		// Pola Cocok: /api/cash/shifts/{id}/report -> len(parts) == 5 dan parts[4] == "report"
		if len(parts) == 5 && parts[4] == "report" {
			handler.UploadShiftReportHandler(w, r)
			return
		}

		http.NotFound(w, r)
	})

	// 2. Transaksi Utama Laci Kasir (Transactions Core Collections)
	mux.HandleFunc("/api/cash/transactions", handler.TransactionsHandler)
	mux.HandleFunc("/api/cash/transactions/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Pola Cocok: /api/cash/transactions/{id} -> len(parts) == 4
		if len(parts) == 4 {
			handler.TransactionItemHandler(w, r)
			return
		}

		http.NotFound(w, r)
	})

	// 3. Kategori Pengeluaran Toko (Expense Categories)
	mux.HandleFunc("/api/cash/expense-categories", handler.ExpenseCategoriesHandler)
	mux.HandleFunc("/api/cash/expense-categories/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Pola Cocok: /api/cash/expense-categories/{id} -> len(parts) == 4
		if len(parts) == 4 {
			handler.ExpenseCategoryIDHandler(w, r)
			return
		}

		http.NotFound(w, r)
	})

	// 4. Pengeluaran Operasional Toko (Expenses Routing)
	mux.HandleFunc("/api/cash/expenses", handler.ExpensesHandler)
	mux.HandleFunc("/api/cash/expenses/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Pola Cocok: /api/cash/expenses/{id} -> len(parts) == 4
		if len(parts) == 4 {
			handler.ExpenseItemHandler(w, r)
			return
		}

		http.NotFound(w, r)
	})

	// 5. Jaringan Endpoint Internal Sistem Lintas-Skema
	// Menangani hook masuk otomatis dari POS outbox untuk pencatatan nominal cash_db.transactions
	mux.HandleFunc("/internal/cash/income", handler.InternalIncomeHandler)

	// Pendapatan Internal (GET query)
	mux.HandleFunc("/api/cash/internal-incomes", handler.InternalIncomesQueryHandler)
}
