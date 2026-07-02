package entity

import "time"

// Category merepresentasikan model data dari tabel 'categories'
type Category struct {
	ID        string     `json:"id"` // UUID
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete support
}
