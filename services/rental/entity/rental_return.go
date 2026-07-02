package entity

import "time"

type RentalReturn struct {
	ID               string    `json:"id"`
	ReservationID    string    `json:"reservation_id"`
	ReturnDate       time.Time `json:"return_date"`
	LateDays         int       `json:"late_days"`         // Jumlah hari keterlambatan
	TotalLateFees    float64   `json:"total_late_fees"`   // Denda keterlambatan: 10% × total_amount × late_days
	TotalDamageFees  float64   `json:"total_damage_fees"` // Total denda kerusakan dari semua item
	RemainingPayment float64   `json:"remaining_payment"` // Sisa tagihan belum terbayar (total_amount - down_payment)
	AmountPaid       float64   `json:"amount_paid"`
	ChangeAmount     float64   `json:"change_amount"`
	GrandTotalPaid   float64   `json:"grand_total_paid"`  // Total yang wajib dibayar saat pengembalian
	Notes            string    `json:"notes"`
	ReceivedBy       string    `json:"received_by"`
	ReceiptURL       string    `json:"receipt_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
