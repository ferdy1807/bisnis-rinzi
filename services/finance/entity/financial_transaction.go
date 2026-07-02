package entity

import "time"

type FinancialTransaction struct {
	ID              string    `json:"id"`
	TransactionType string    `json:"transaction_type"` // e.g., "INCOME", "EXPENSE", "SALES", "RENTAL"
	ReferenceID     string    `json:"reference_id"`     // UUID dari service asal (SaleID, ReservationID, dll)
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
}
