package dto

import "time"

type CreateCOARequest struct {
	AccountCode   string `json:"account_code"`
	AccountName   string `json:"account_name"`
	AccountGroup  string `json:"account_group"`
	NormalBalance string `json:"normal_balance"`
	IsActive      bool   `json:"is_active"`
}

type OpenPeriodRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type CreateJournalEntryDetailPayload struct {
	AccountID    string  `json:"account_id"`
	DebitAmount  float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
}

type CreateJournalEntryRequest struct {
	JournalID       string                            `json:"journal_id"`
	ReferenceNumber string                            `json:"reference_number"`
	Narration       string                            `json:"narration"`
	Details         []CreateJournalEntryDetailPayload `json:"details"`
}

type ReconcilePayload struct {
	TargetSystem string  `json:"target_system"`
	ActualAmount float64 `json:"actual_amount"`
	Notes        string  `json:"notes"`
}

type ProcessDailyClosingRequest struct {
	ClosingDate     time.Time          `json:"closing_date"`
	Reconciliations []ReconcilePayload `json:"reconciliations"`
}

type LedgerReportLine struct {
	EntryDate       time.Time `json:"entry_date"`
	ReferenceNumber string    `json:"reference_number"`
	Narration       string    `json:"narration"`
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	RunningBalance  float64   `json:"running_balance"`
}

type LedgerReportResponse struct {
	AccountCode    string             `json:"account_code"`
	AccountName    string             `json:"account_name"`
	InitialBalance float64            `json:"initial_balance"`
	FinalBalance   float64            `json:"final_balance"`
	Lines          []LedgerReportLine `json:"lines"`
}

type TrialBalanceLine struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}

type BalanceSheetResponse struct {
	PeriodName  string             `json:"period_name"`
	TotalAssets float64            `json:"total_assets"`
	TotalLiab   float64            `json:"total_liabilities"`
	TotalEquity float64            `json:"total_equity"`
	Items       map[string]float64 `json:"items"`
}

type IncomeStatementResponse struct {
	PeriodName   string             `json:"period_name"`
	TotalRevenue float64            `json:"total_revenue"`
	TotalCOGS    float64            `json:"total_cogs"`
	GrossProfit  float64            `json:"gross_profit"`
	TotalExpense float64            `json:"total_expense"`
	NetIncome    float64            `json:"net_income"`
	Items        map[string]float64 `json:"items"`
}

type JournalIncomingRequest struct {
	ID            string    `json:"id"`
	AggregateType string    `json:"aggregate_type"`
	AggregateID   string    `json:"aggregate_id"`
	EventType     string    `json:"event_type"`
	Payload       string    `json:"payload"` // Dibaca sebagai string mentah teks JSON
	CreatedAt     time.Time `json:"created_at"`
}
