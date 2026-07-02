package entity

import "time"

type SaleItem struct {
	ID          string    `json:"id"`
	SaleID      string    `json:"sale_id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	UnitCode    string    `json:"unit_code"`
	Qty       float64   `json:"qty"`
	UnitPrice float64   `json:"unit_price"`
	Discount  float64   `json:"discount"`
	Subtotal  float64   `json:"subtotal"`
	CostPrice float64   `json:"cost_price"` // Harga modal saat barang terjual untuk HPP
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
