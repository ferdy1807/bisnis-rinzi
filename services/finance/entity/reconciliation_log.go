package entity

import "time"

type ReconciliationLog struct {
	ID             string    `json:"id"`
	DailyClosingID string    `json:"daily_closing_id"`
	TargetSystem   string    `json:"target_system"` // e.g., "POS", "CASH", "RENTAL"
	SystemAmount   float64   `json:"system_amount"`
	ActualAmount   float64   `json:"actual_amount"`
	Discrepancy    float64   `json:"discrepancy"`
	Notes          string    `json:"notes"`
	ReconciledBy   string    `json:"reconciled_by"`
	CreatedAt      time.Time `json:"created_at"`
}
