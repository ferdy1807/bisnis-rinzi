package entity

import "time"

type StockReservation struct {
	ID              string    `json:"id"`
	RentalProductID string    `json:"rental_product_id"`
	ReservationID   string    `json:"reservation_id"`
	ReserveDate     time.Time `json:"reserve_date"` // Dipecah per hari untuk alokasi matriks inventaris
	QtyReserved     int       `json:"qty_reserved"`
	CreatedAt       time.Time `json:"created_at"`
}
