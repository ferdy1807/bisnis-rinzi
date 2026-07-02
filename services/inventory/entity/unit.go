package entity

import "time"

// Unit merepresentasikan model data dari tabel 'units'
type Unit struct {
	ID        string    `json:"id"` // UUID
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
