package entity

import "time"

type SyncLog struct {
	ID           string    `json:"id"`
	EntityType   string    `json:"entity_type"` // "SALE"
	EntityID     string    `json:"entity_id"`
	SyncStatus   string    `json:"sync_status"` // SUCCESS, FAILED
	ErrorMessage *string   `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
