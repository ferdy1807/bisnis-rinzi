package entity

import "time"

type Expense struct {
	ID          string    `json:"id"`
	ExpenseDate time.Time `json:"expense_date"`
	CategoryID  string    `json:"category_id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
