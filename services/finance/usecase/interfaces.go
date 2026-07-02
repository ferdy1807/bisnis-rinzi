package usecase

import (
	"bisnis-rinzi/services/finance/dto"
	"bisnis-rinzi/services/finance/entity"
	"context"
	"time"
)

type FinanceUseCase interface {
	EnqueueInboxEvent(ctx context.Context, req dto.JournalIncomingRequest) error
	// 1. Setup Master Chart of Accounts (COA)
	CreateAccount(ctx context.Context, req dto.CreateCOARequest) error
	GetAccounts(ctx context.Context) ([]*entity.ChartOfAccount, error)
	GetAccountByID(ctx context.Context, id string) (*entity.ChartOfAccount, error)
	UpdateAccount(ctx context.Context, id string, req dto.CreateCOARequest) error
	DeleteAccount(ctx context.Context, id string) error

	// 2. Pemrosesan Entri Jurnal Double-Entry
	RecordManualJournal(ctx context.Context, req dto.CreateJournalEntryRequest) error
	RecordAutoJournalSystem(ctx context.Context, sourceModule string, refID string, lines []entity.JournalEntryDetail) error
	GetAllJournals(ctx context.Context) ([]*entity.JournalEntry, error)
	GetJournalByID(ctx context.Context, id string) (*entity.JournalEntry, error)

	// 3. Manajemen Siklus Periode Akuntansi
	OpenAccountingPeriod(ctx context.Context, req dto.OpenPeriodRequest) error
	UpdateAccountingPeriod(ctx context.Context, id string, req dto.OpenPeriodRequest) error
	DeleteAccountingPeriod(ctx context.Context, id string) error
	LockPeriod(ctx context.Context, periodID string, lockedBy, reason string) error
	GetAllPeriods(ctx context.Context) ([]*entity.AccountingPeriod, error)

	// 4. Proses Rekonsiliasi & Tutup Buku Kasir
	ProcessClosingAndReconciliation(ctx context.Context, cashierID string, req dto.ProcessDailyClosingRequest) error
	GetAllDailyClosings(ctx context.Context) ([]*entity.DailyClosing, error)
	GetDailyClosingDetail(ctx context.Context, id string) (*entity.DailyClosing, error)
	GetReconciliationLogs(ctx context.Context, closingID string) ([]*entity.ReconciliationLog, error)

	// 5. Engine Pelaporan SAK
	GenerateLedgerReport(ctx context.Context, accountCode string, start, end time.Time) (*dto.LedgerReportResponse, error)
	GenerateTrialBalance(ctx context.Context) ([]dto.TrialBalanceLine, error)
	GenerateIncomeStatement(ctx context.Context) (*dto.IncomeStatementResponse, error)
	GenerateBalanceSheet(ctx context.Context) (*dto.BalanceSheetResponse, error)
	GenerateCashFlowStatement(ctx context.Context) (map[string]interface{}, error)

	// 6. Central Analytics
	GetCorporateDashboardAnalytics(ctx context.Context, start, end time.Time) (map[string]interface{}, error)
	GetAnalyticsData(ctx context.Context, action string) (map[string]interface{}, error)

	// 7. Ekspor & Laporan Fisik
	UploadDailyIncomeReport(ctx context.Context, filename string, fileBytes []byte) (string, error)
}
