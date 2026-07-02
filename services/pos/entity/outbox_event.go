package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

// POSOutboxEvent merepresentasikan struktur baris tabel outbox di pos_db
type POSOutboxEvent struct {
	ID            string    `json:"id"`
	AggregateType string    `json:"aggregate_type"` // "RETAIL_ORDER"
	AggregateID   string    `json:"aggregate_id"`   // Invoice ID / Order ID
	EventType     string    `json:"event_type"`     // "SALE_COMPLETED"
	Payload       []byte    `json:"payload"`        // Data mentah transaksi kasir
	Status        string    `json:"status"`         // "PENDING", "PROCESSED", "FAILED"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// POSSaleItemPayload mendefinisikan detail item barang yang dibeli konsumen
type POSSaleItemPayload struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"qty"`
	UnitPrice float64 `json:"unit_price"`
	SubTotal  float64 `json:"subtotal"`
}

// POSSaleCompletedPayload membungkus data utuh penjualan lunas kasir POS ritel
type POSSaleCompletedPayload struct {
	InvoiceNumber    string               `json:"invoice_number"`
	CashierID        string               `json:"cashier_id"`
	TotalTransaction float64              `json:"total_transaction"` // Nilai bruto transaksi
	TaxAmount        float64              `json:"tax_amount"`        // Nominal PPN
	TotalAmount      float64              `json:"total_amount"`      // Net setelah pajak
	PaymentMethod    string               `json:"payment_method"`    // "TUNAI", "QRIS", "DEBIT"
	Items            []POSSaleItemPayload `json:"items"`
	CreatedAt        time.Time            `json:"created_at"`
}

// UnmarshalSalePayload mengurai []byte mentah menjadi struct penjualan kasir
func UnmarshalSalePayload(raw []byte) (*POSSaleCompletedPayload, error) {
	var p POSSaleCompletedPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("pos_entity: gagal unmarshal payload penjualan: %w", err)
	}
	return &p, nil
}
