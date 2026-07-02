package entity

import "time"

// Role merepresentasikan tabel 'roles' pada auth_db
type Role struct {
	Code         string    `json:"code"`          // Primary Key (e.g., OWNER, CASHIER)
	Name         string    `json:"name"`          // Nama peran terjemahan
	DashboardURL string    `json:"dashboard_url"` // URL pengalihan halaman dashboard setelah login
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
