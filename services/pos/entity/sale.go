package entity

import (
	"errors"
	"strings"
	"time"
)

type Sale struct {
	ID               string    `json:"id"`
	IdempotencyKey   string    `json:"idempotency_key"`
	InvoiceNumber    string    `json:"invoice_number"`
	TransactionDate  time.Time `json:"transaction_date"`
	Subtotal         float64   `json:"subtotal"`
	Discount         float64   `json:"discount"`
	Total            float64   `json:"total"`
	AmountPaid       float64   `json:"amount_paid"`
	ChangeAmount     float64   `json:"change_amount"`
	PaymentMethod    string    `json:"payment_method"` // CASH, QRIS, TRANSFER
	PaymentStatus    string    `json:"payment_status"` // COMPLETED, PENDING, CANCELLED
	CashierID        string    `json:"cashier_id"`
	CashierSessionID string    `json:"cashier_session_id"`
	InvoiceURL       *string   `json:"invoice_url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (s *Sale) Validate() error {
	if strings.TrimSpace(s.IdempotencyKey) == "" {
		return errors.New("idempotency key wajib diisi untuk mencegah duplikasi data transaksi")
	}
	if strings.TrimSpace(s.InvoiceNumber) == "" {
		return errors.New("nomor invoice tidak boleh kosong")
	}
	if s.Total < 0 {
		return errors.New("total nilai transaksi penjualan tidak boleh negatif")
	}
	if s.CashierSessionID == "" {
		return errors.New("transaksi wajib terikat dengan sesi harian kasir yang aktif")
	}
	if s.PaymentMethod == "CASH" && s.AmountPaid < s.Total {
		return errors.New("nominal pembayaran tunai tidak boleh kurang dari total transaksi")
	}
	return nil
}
