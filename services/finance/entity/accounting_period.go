package entity

import (
	"errors"
	"time"
)

type AccountingPeriod struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"` // e.g., "Januari 2026", "Q1 2026"
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsClosed  bool      `json:"is_closed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ap *AccountingPeriod) Validate() error {
	if ap.Name == "" {
		return errors.New("nama periode akuntansi wajib diisi")
	}
	if !ap.EndDate.After(ap.StartDate) {
		return errors.New("tanggal akhir periode harus setelah tanggal awal periode")
	}
	return nil
}
