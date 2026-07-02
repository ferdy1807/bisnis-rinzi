package dto

import "time"

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CreateProductRequest struct {
	CategoryID        string  `json:"category_id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	RentalPrice       float64 `json:"rental_price"`
	QuantityAvailable float64 `json:"quantity_available"`
	IsActive          bool    `json:"is_active"`
}

type ReservationItemPayload struct {
	RentalProductID string  `json:"rental_product_id"`
	Qty             float64 `json:"qty"`
	PricePerPeriod  float64 `json:"price_per_period"`
}

type CreateReservationRequest struct {
	CustomerName     string                      `json:"customer_name"`
	CustomerIdentity string                      `json:"customer_identity"`
	CustomerPhone    string                      `json:"customer_phone"`
	StartDate        time.Time                   `json:"start_date"`
	EndDate          time.Time                   `json:"end_date"`
	Discount         float64                     `json:"discount"`
	DownPayment      float64                     `json:"down_payment"` // Uang muka yang dibayar di awal (bisa partial)
	AmountPaid       float64                     `json:"amount_paid"`  // Uang yang dibayarkan konsumen saat DP
	EventDate        *time.Time                  `json:"event_date,omitempty"`
	CashierSessionID string                      `json:"cashier_session_id"`
	Items            []ReservationItemPayload    `json:"items"`
	Contents         []ReservationContentPayload `json:"contents"`
}

type ReservationContentPayload struct {
	ItemName       string `json:"item_name"`
	Description    string `json:"description"`
	Quantity       int    `json:"quantity"`
	ConditionNotes string `json:"condition_notes"`
}

type ReturnItemPayload struct {
	RentalProductID string  `json:"rental_product_id"`
	QtyReturned     float64 `json:"qty_returned"`
	ConditionStatus string  `json:"condition_status"` // 'GOOD', 'DAMAGED', 'LOST'
	DamageFee       float64 `json:"damage_fee"`       // Denda kerusakan (Rp), diinput kasir
	ConditionNotes  string  `json:"condition_notes"`  // Catatan kondisi barang
}

type ProcessReturnRequest struct {
	ReservationID        string              `json:"reservation_id"`
	DiscountCompensation float64             `json:"discount_compensation"` // Potongan manual dari pegawai jika merias lambat
	AmountPaid           float64             `json:"amount_paid"`
	ChangeAmount         float64             `json:"change_amount"`
	ManualDamageFee      float64             `json:"manual_damage_fee"`   // Tagihan denda agregat manual dari frontend
	ManualReturnNotes    string              `json:"manual_return_notes"` // Catatan inspeksi manual dari frontend
	ReturnItems          []ReturnItemPayload `json:"return_items"`
}

type CancelReservationRequest struct {
	ReservationID string  `json:"reservation_id"`
	PenaltyFee    float64 `json:"penalty_fee"` // Denda pembatalan manual yang diisi pegawai
}

type CheckAvailabilityRequest struct {
	RentalProductID string    `json:"rental_product_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	QtyRequested    int       `json:"qty_requested"`
}

type DamagedItemAudit struct {
	ID              string  `json:"id"`
	ItemName        string  `json:"item_name"`
	CustomerName    string  `json:"customer_name"`
	ConditionStatus string  `json:"condition_status"`
	ConditionNotes  string  `json:"condition_notes"`
	DamageFee       float64 `json:"damage_fee"`
	Status          string  `json:"status"`
}

type SettleDamageRequest struct {
	PaymentAction string `json:"payment_action"`
	AuditNotes    string `json:"audit_notes"`
}
