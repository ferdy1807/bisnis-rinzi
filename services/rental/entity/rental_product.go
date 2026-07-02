package entity

import (
	"errors"
	"strings"
	"time"
)

type RentalProduct struct {
	ID                string     `json:"id"`
	CategoryID        string     `json:"category_id"`
	Code              string     `json:"code"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	RentalPrice       float64    `json:"rental_price"`
	QuantityAvailable float64    `json:"quantity_available"`
	IsActive          bool       `json:"is_active"`
	ObjectName        *string    `json:"object_name,omitempty"`
	OriginalFileName  *string    `json:"original_file_name,omitempty"`
	MimeType          *string    `json:"mime_type,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func (rp *RentalProduct) Validate() error {
	if strings.TrimSpace(rp.Code) == "" {
		return errors.New("kode produk rental tidak boleh kosong")
	}
	if strings.TrimSpace(rp.Name) == "" {
		return errors.New("nama produk rental tidak boleh kosong")
	}
	if rp.RentalPrice <= 0 {
		return errors.New("tarif sewa harus lebih besar dari nol")
	}
	return nil
}
