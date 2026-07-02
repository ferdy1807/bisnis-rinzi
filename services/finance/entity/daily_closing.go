package entity

import "time"

type DailyClosing struct {
	ID                string    `json:"id"`
	ClosingDate       time.Time `json:"closing_date"`
	TotalSalesRetail  float64   `json:"total_sales_retail"`
	TotalRentalIncome float64   `json:"total_rental_income"`
	TotalOtherIncome  float64   `json:"total_other_income"`
	TotalExpenses     float64   `json:"total_expenses"`
	NetCashFlow       float64   `json:"net_cash_flow"`
	ActualCash        float64   `json:"actual_cash"`
	OpeningCash       float64   `json:"opening_cash"`
	IsReconciled      bool      `json:"is_reconciled"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
