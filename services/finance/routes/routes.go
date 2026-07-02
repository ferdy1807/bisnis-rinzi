package routes

import (
	delivery "bisnis-rinzi/services/finance/delivery/http"
	"net/http"
)

func RegisterFinanceRoutes(mux *http.ServeMux, handler *delivery.FinanceHandler) {
	// 1. Daily Closing
	mux.HandleFunc("/api/finance/daily-closings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ClosingsQueryHandler(w, r)
		case http.MethodPost:
			handler.DailyClosingHandler(w, r)
		}
	})
	mux.HandleFunc("/api/finance/daily-closings/", handler.ClosingDetailDispatcher)

	// 2. Period Lock & Periods
	mux.HandleFunc("/api/finance/accounting-periods", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.PeriodsQueryHandler(w, r)
		case http.MethodPost:
			handler.PeriodOpenHandler(w, r)
		}
	})
	mux.HandleFunc("/api/finance/period-locks", handler.PeriodLockHandler)

	// 3. Accounting / General Ledger
	mux.HandleFunc("/api/finance/accounts", handler.COAHandler)
	mux.HandleFunc("/api/finance/accounts/", handler.AccountDetailDispatcher)

	mux.HandleFunc("/api/finance/journal-entries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.JournalsQueryHandler(w, r)
		case http.MethodPost:
			handler.JournalEntryHandler(w, r)
		}
	})
	mux.HandleFunc("/api/finance/journal-entries/", handler.JournalDetailDispatcher)

	// 4. Dashboard
	mux.HandleFunc("/api/finance/analytics/dashboard", handler.AnalyticsDashboardHandler)

	// 5. Reconciliation
	mux.HandleFunc("/api/finance/reconciliation", handler.ReconciliationHandler)

	// 6. Reports
	mux.HandleFunc("/api/finance/reports/profit-loss", handler.IncomeStatementReportHandler)
	mux.HandleFunc("/api/finance/reports/cash-flow", handler.CashFlowReportHandler)
	mux.HandleFunc("/api/finance/reports/balance-sheet", handler.BalanceSheetReportHandler)
	mux.HandleFunc("/api/finance/reports/export/", handler.ExportReportHandler)

	// 7. Analytics
	mux.HandleFunc("/api/finance/analytics/", handler.AnalyticsDispatcher)
	mux.HandleFunc("/api/finance/inbox-events", handler.ProcessInboxEventHandler)

	// 8. PDF Export & Upload
	mux.HandleFunc("/api/finance/daily-incomes/upload-pdf", handler.UploadDailyIncomeHandler)
}
