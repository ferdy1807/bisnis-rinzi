package entity

import "time"

type CashTransaction struct {
	ID              string    `json:"id"`
	SessionID       string    `json:"session_id"`
	TransactionType string    `json:"transaction_type"` // "DEPOSIT", "WITHDRAWAL"
	ReferenceType   string    `json:"reference_type"`   // "EXPENSE", "MANUAL"
	ReferenceID     *string   `json:"reference_id,omitempty"`
	Amount          float64   `json:"amount"`
	Notes           *string   `json:"notes,omitempty"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
