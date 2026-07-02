package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

// RentalOutboxEvent merepresentasikan struktur baris data outbox pada tabel outbox_events di rental_db.
type RentalOutboxEvent struct {
	ID            string    `json:"id"`
	AggregateType string    `json:"aggregate_type"` // "RENTAL_RESERVATION" atau "RESERVATION"
	AggregateID   string    `json:"aggregate_id"`   // ReservationID
	EventType     string    `json:"event_type"`     // "RESERVATION_CREATED", "PRODUCT_RETURN_PROCESSED", dll.
	Payload       []byte    `json:"payload"`        // Teks biner JSONB asli dari database lokal
	Status        string    `json:"status"`         // "PENDING", "SUCCESS", "FAILED"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RentalReservationPayload menangkap data operasional saat pemesanan awal hantaran dibuat.
type RentalReservationPayload struct {
	ReservationID string    `json:"reservation_id"`
	CustomerID    string    `json:"customer_id"`
	BoxCode       string    `json:"box_code"`
	RentalPrice   float64   `json:"rental_price"`
	DepositAmount float64   `json:"deposit_amount"`
	PenaltyAmount float64   `json:"penalty_amount,omitempty"`
	StatusBooking string    `json:"status_booking"` // "CONFIRMED", "RETURNED", "DAMAGED"
	DeadlineDate  time.Time `json:"deadline_date"`
}

// LocalProductReturnPayload menangkap data riil saat mika box dikembalikan oleh pelanggan.
type LocalProductReturnPayload struct {
	ReturnID         string    `json:"return_id"`
	ReservationID    string    `json:"reservation_id"`
	LateDays         int       `json:"late_days"`
	TotalLateFees    float64   `json:"total_late_fees"`
	TotalDamageFees  float64   `json:"total_damage_fees"`
	RemainingPayment float64   `json:"remaining_payment"`
	GrandTotalPaid   float64   `json:"grand_total_paid"`
	Timestamp        time.Time `json:"timestamp"`
}

// UnmarshalReservationPayload mengurai slice byte mentah database menjadi bentuk objek struktur reservasi awal.
func UnmarshalReservationPayload(raw []byte) (*RentalReservationPayload, error) {
	var p RentalReservationPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("rental_entity: gagal mengekstrak data reservasi: %w", err)
	}
	return &p, nil
}

// UnmarshalProductReturnPayload mengurai data biner transaksi pengembalian barang menjadi bentuk struct return murni.
func UnmarshalProductReturnPayload(raw []byte) (*LocalProductReturnPayload, error) {
	var p LocalProductReturnPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("rental_entity: gagal mengekstrak data pengembalian produk: %w", err)
	}
	return &p, nil
}
