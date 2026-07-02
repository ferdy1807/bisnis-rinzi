package entity

import "time"

type ExpenseCategory struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"` // e.g., "UTIL", "LOGISTIC"
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
