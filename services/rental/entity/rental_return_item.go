package entity

import "time"

type RentalReturnItem struct {
	ID              string    `json:"id"`
	RentalReturnID  string    `json:"rental_return_id"`
	RentalProductID string    `json:"rental_product_id"`
	RentalProductName string  `json:"rental_product_name"` // Snapshot nama produk saat pengembalian
	QtyReturned     float64   `json:"qty_returned"`
	ConditionStatus string    `json:"condition_status"`  // 'GOOD', 'DAMAGED', 'LOST'
	DamageFee       float64   `json:"damage_fee"`        // Denda kerusakan untuk item ini (input kasir)
	ConditionNotes  string    `json:"condition_notes"`   // Catatan kondisi barang
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
