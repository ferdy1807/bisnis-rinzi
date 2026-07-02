package entity

import (
	"errors"
	"strings"
	"time"
)

type RentalCategory struct {
	ID          string     `json:"id"`
	Code        string     `json:"code"` // e.g., "CAM-DSLR", "LENS-FIX"
	Name        string     `json:"name"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func (rc *RentalCategory) Validate() error {
	if strings.TrimSpace(rc.Code) == "" {
		return errors.New("kode kategori rental tidak boleh kosong")
	}
	if strings.TrimSpace(rc.Name) == "" {
		return errors.New("nama kategori rental tidak boleh kosong")
	}
	return nil
}
