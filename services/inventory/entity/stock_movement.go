package entity

import "time"

// StockMovement merepresentasikan riwayat mutasi dari tabel 'stock_movements'
type StockMovement struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name,omitempty"`
	SKU          string    `json:"sku,omitempty"`
	MovementType string    `json:"movement_type"` // e.g., "IN", "OUT", "ADJUSTMENT"
	Qty          float64   `json:"qty"`
	Reference    string    `json:"reference"` // e.g., "INVOICE-001", "STOCK_ADJUST"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
