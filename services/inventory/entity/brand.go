package entity

import "time"

// Brand merepresentasikan model data dari tabel 'brands'
type Brand struct {
	ID        string     `json:"id"` // UUID
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete support
}
