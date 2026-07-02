package entity

import "time"

// ProductStock merepresentasikan tabel 'product_stocks'
type ProductStock struct {
	ProductID      string    `json:"product_id"`
	Qty            float64   `json:"qty"`
	QtyMinStock    float64   `json:"qty_min_stock"`
	QtySafetyStock float64   `json:"qty_safety_stock"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
