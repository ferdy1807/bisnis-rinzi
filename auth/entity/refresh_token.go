package entity

import "time"

// RefreshToken merepresentasikan tabel 'refresh_tokens' pada auth_db
type RefreshToken struct {
	ID         string     `json:"id"` // UUID Primary Key
	UserID     string     `json:"user_id"`
	Token      string     `json:"token"`
	ExpiresAt  time.Time  `json:"expires_at"`
	DeviceInfo string     `json:"device_info"`
	IPAddress  string     `json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

// IsExpired Memeriksa apakah token sesi ini sudah kedaluwarsa secara waktu
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsActive Memeriksa apakah token masih aktif dan belum dicabut
func (rt *RefreshToken) IsActive() bool {
	return rt.RevokedAt == nil && !rt.IsExpired()
}
