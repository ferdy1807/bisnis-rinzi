package entity

import "time"

type PeriodLock struct {
	ID                 string    `json:"id"`
	AccountingPeriodID string    `json:"accounting_period_id"`
	LockedBy           string    `json:"locked_by"` // UserID pengunci
	LockReason         string    `json:"lock_reason"`
	CreatedAt          time.Time `json:"created_at"`
}
