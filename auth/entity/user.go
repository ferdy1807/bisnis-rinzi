package entity

import (
	"errors"
	"strings"
	"time"
)

// User merepresentasikan tabel 'users' pada auth_db
type User struct {
	ID           string     `json:"id"` // UUID Primary Key
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"` // Disembunyikan dari JSON response demi keamanan
	FullName     string     `json:"full_name"`
	Role         string     `json:"role"` // Foreign Key ke roles.code
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"` // Nullable untuk soft delete
}

// Validate Memastikan integritas data user sebelum diproses ke database
func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("username tidak boleh kosong")
	}
	if len(u.Username) < 4 {
		return errors.New("username minimal harus berukuran 4 karakter")
	}
	if strings.TrimSpace(u.FullName) == "" {
		return errors.New("nama lengkap tidak boleh kosong")
	}
	if strings.TrimSpace(u.Role) == "" {
		return errors.New("role pengguna wajib ditentukan")
	}
	return nil
}
