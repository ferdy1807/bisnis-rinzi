package entity

import (
	"errors"
	"time"
)

type CustomerSnapshot struct {
	ID             string    `json:"id"`
	CustomerName   string    `json:"customer_name"`
	CustomerPhone  string    `json:"customer_phone"`
	CustomerIDCard string    `json:"customer_id_card"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Reservation struct {
	ID                 string             `json:"id"`
	InvoiceNumber      string             `json:"invoice_number"`
	CustomerSnapshotID string             `json:"customer_snapshot_id"`
	CustomerName       string             `json:"customer_name,omitempty"`
	CustomerPhone      string             `json:"customer_phone,omitempty"`
	TransactionDate    time.Time          `json:"transaction_date"`
	StartDate          time.Time          `json:"start_date"`
	EndDate            time.Time          `json:"end_date"`
	EventDate          *time.Time         `json:"event_date,omitempty"`
	Subtotal           float64            `json:"subtotal"`
	DownPayment        float64            `json:"down_payment"`
	AmountPaid         float64            `json:"amount_paid"`
	ChangeAmount       float64            `json:"change_amount"`
	TotalAmount        float64            `json:"total_amount"`
	Status             string             `json:"status"` // "BOOKED", "CONTENTS_RECEIVED", "DECORATING", "READY_FOR_PICKUP", "PICKED_UP", "RETURNED", "CANCELLED"
	PickedUpBy         *string            `json:"picked_up_by,omitempty"`
	PickedUpAt         *time.Time         `json:"picked_up_at,omitempty"`
	CashierSessionID   string             `json:"cashier_session_id"`
	CreatedBy          string             `json:"created_by"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	GrandTotalIncome   float64            `json:"grand_total_income"`
	Items              []*ReservationItem  `json:"items,omitempty"`
	Contents           []*ReservationContent `json:"contents,omitempty"`
}

func (r *Reservation) Validate() error {
	if !r.EndDate.After(r.StartDate) && !r.EndDate.Equal(r.StartDate) {
		return errors.New("waktu penyelesaian rental harus valid")
	}
	if r.TotalAmount < 0 {
		return errors.New("total nilai reservasi sewa tidak valid")
	}
	return nil
}

type ReservationContent struct {
	ID                  string    `json:"id"`
	RentalReservationID string    `json:"rental_reservation_id"`
	ItemName            string    `json:"item_name"`
	Description         string    `json:"description"`
	Quantity            int       `json:"quantity"`
	ConditionNotes      string    `json:"condition_notes"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
