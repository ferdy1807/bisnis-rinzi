package dto

type OpenSessionRequest struct {
	OpeningCash float64 `json:"opening_cash"`
}

type CloseSessionRequest struct {
	ActualCash float64 `json:"actual_cash"`
}

type CashTransactionRequest struct {
	TransactionType string  `json:"transaction_type"` // "DEPOSIT", "WITHDRAWAL"
	Amount          float64 `json:"amount"`
	Notes           string  `json:"notes"`
}

type InternalIncomeRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Source      string  `json:"source"`
	Reference   string  `json:"reference"`
}

type CreateExpenseCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CreateExpenseRequest struct {
	CategoryID  string  `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type CashierSessionResponse struct {
	ID           string   `json:"id"`
	CashierID    string   `json:"cashier_id"`
	OpeningCash  float64  `json:"opening_cash"`
	ExpectedCash *float64 `json:"expected_cash,omitempty"`
	ActualCash   *float64 `json:"actual_cash,omitempty"`
	Difference   *float64 `json:"difference,omitempty"`
	Status       string   `json:"status"`
	ReceiptURL   *string  `json:"receipt_url,omitempty"`
	OpenTime     string   `json:"open_time"`
	CloseTime    *string  `json:"close_time,omitempty"`
}

type ShiftSummaryResponse struct {
	SessionID       string  `json:"session_id"`
	CashierID       string  `json:"cashier_id"`
	OpeningCash     float64 `json:"opening_cash"`
	TotalIncome     float64 `json:"total_income"`
	TotalExpense    float64 `json:"total_expense"`
	TotalDeposit    float64 `json:"total_deposit"`
	TotalWithdrawal float64 `json:"total_withdrawal"`
	ExpectedCash    float64 `json:"expected_cash"`
	ActualCash      float64 `json:"actual_cash,omitempty"`
	Difference      float64 `json:"difference,omitempty"`
	Status          string  `json:"status"`
}
