package entity

import "time"

type ReservationItem struct {
	ID                  string    `json:"id"`
	RentalReservationID string    `json:"rental_reservation_id"`
	RentalProductID     string    `json:"rental_product_id"`
	RentalProductName   string    `json:"rental_product_name"` // Snapshot nama produk saat reservasi (immutable)
	Qty                 float64   `json:"qty"`
	PricePerPeriod      float64   `json:"price_per_period"`
	Subtotal            float64   `json:"subtotal"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
