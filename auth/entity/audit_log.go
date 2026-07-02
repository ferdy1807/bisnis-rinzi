package entity

import "time"

// AuditLog merepresentasikan tabel 'audit_logs' pada auth_db untuk tracking aktivitas data
type AuditLog struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Action     string    `json:"action"`      // Jenis aksi (e.g., CREATE_USER, CHANGE_PASSWORD)
	EntityName string    `json:"entity_name"` // Nama tabel terkait (e.g., users)
	EntityID   string    `json:"entity_id"`
	OldData    []byte    `json:"old_data,omitempty"` // Menggunakan JSONB di PostgreSQL
	NewData    []byte    `json:"new_data,omitempty"` // Menggunakan JSONB di PostgreSQL
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
