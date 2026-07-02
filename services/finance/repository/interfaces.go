package repository

import (
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/finance/dto"
	"bisnis-rinzi/services/finance/entity"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type FinanceRepository interface {
	// Chart of Accounts (COA)
	SaveCOA(ctx context.Context, coa *entity.ChartOfAccount) error
	FindAllCOA(ctx context.Context) ([]*entity.ChartOfAccount, error)
	FindCOAByID(ctx context.Context, id string) (*entity.ChartOfAccount, error)
	FindCOAByCode(ctx context.Context, code string) (*entity.ChartOfAccount, error)
	UpdateCOABalance(ctx context.Context, id string, amount float64) error
	UpdateCOA(ctx context.Context, coa *entity.ChartOfAccount) error
	DeleteCOA(ctx context.Context, id string) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
	EnqueueInboxEvent(ctx context.Context, req dto.JournalIncomingRequest) error

	// REQUIRED BY FINANACE WORKER (Mendukung transaksi lintas-tabel)
	InsertJournalEntry(ctx context.Context, tx pgx.Tx, aggregateID, narasi string, debetAcc, kreditAcc string, total float64) error
	IsEventProcessed(ctx context.Context, tx pgx.Tx, eventID string) (bool, error)
	MarkEventAsProcessed(ctx context.Context, tx pgx.Tx, eventID string) error

	// Periods & Locking Control
	SavePeriod(ctx context.Context, period *entity.AccountingPeriod) error
	FindAllPeriods(ctx context.Context) ([]*entity.AccountingPeriod, error)
	FindActivePeriod(ctx context.Context) (*entity.AccountingPeriod, error)
	FindPeriodByID(ctx context.Context, id string) (*entity.AccountingPeriod, error)
	UpdatePeriod(ctx context.Context, period *entity.AccountingPeriod) error
	DeletePeriod(ctx context.Context, id string) error
	ClosePeriodTx(ctx context.Context, periodID string, lock *entity.PeriodLock) error

	// Journals & Jurnal Entri
	SaveJournal(ctx context.Context, journal *entity.Journal) error
	FindJournalByCode(ctx context.Context, code string) (*entity.Journal, error)
	SaveJournalEntryTx(ctx context.Context, entry *entity.JournalEntry, details []*entity.JournalEntryDetail) error
	FindJournalEntriesByPeriod(ctx context.Context, periodID string) ([]*entity.JournalEntry, error)
	FindJournalEntryDetails(ctx context.Context, entryID string) ([]*entity.JournalEntryDetail, error)
	FindAllJournalEntries(ctx context.Context) ([]*entity.JournalEntry, error)
	FindJournalEntryByID(ctx context.Context, id string) (*entity.JournalEntry, error)

	// Financial Transactions
	SaveFinancialTransaction(ctx context.Context, tx *entity.FinancialTransaction) error
	SumTransactionAmountByType(ctx context.Context, txType string, date time.Time) (float64, error)

	// Closings
	SaveDailyClosingTx(ctx context.Context, dc *entity.DailyClosing, logs []*entity.ReconciliationLog, outboxEvent *outbox.Event) error
	FindAllDailyClosings(ctx context.Context) ([]*entity.DailyClosing, error)
	FindDailyClosingByID(ctx context.Context, id string) (*entity.DailyClosing, error)
	FindDailyClosingByDate(ctx context.Context, d time.Time) (*entity.DailyClosing, error)
	FindReconciliationLogsByClosingID(ctx context.Context, closingID string) ([]*entity.ReconciliationLog, error)

	// Engine Report Generator
	GetLedgerLines(ctx context.Context, accountID string, start, end time.Time) ([]dto.LedgerReportLine, error)
	
	// Analytics
	GetDashboardMetrics(ctx context.Context) (map[string]interface{}, error)
	GetTotalMonthlyAnalytics(ctx context.Context) (float64, float64, float64, error)
	GetMonthlyRevenueTrend(ctx context.Context, year int) ([]map[string]interface{}, error)
	GetSalesAnalytics(ctx context.Context, start, end time.Time) ([]entity.SalesAnalytics, error)
	GetRentalAnalytics(ctx context.Context, start, end time.Time) ([]entity.RentalAnalytics, error)
	GetTopProductsAnalytics(ctx context.Context, limit int) ([]entity.ProductAnalytics, error)
	GetCategoryContributionAnalytics(ctx context.Context) ([]entity.CategoryAnalytics, error)
	GetStockValueAnalytics(ctx context.Context) ([]entity.StockAnalytics, error)
	GetCashierPerformanceAnalytics(ctx context.Context) ([]entity.CashierAnalytics, error)
	GetAnalyticsData(ctx context.Context, action string) (map[string]interface{}, error)

	// Background Analytics Aggregation
	RefreshMonthlyAnalytics(ctx context.Context) error
	RefreshProductAnalytics(ctx context.Context) error
	RefreshDailyClosings(ctx context.Context) error
}
