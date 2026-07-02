package entity

import "time"

type Income struct {
	ID               string    `json:"id"`
	CashierSessionID string    `json:"cashier_session_id"`
	Amount           float64   `json:"amount"`
	Source           string    `json:"source"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
