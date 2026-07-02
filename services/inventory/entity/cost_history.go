package entity

import "time"

// CostHistory merepresentasikan log pergerakan Harga Pokok Penjualan (HPP) / Modal dari suatu produk
type CostHistory struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	AverageCost   float64   `json:"average_cost"`
	EffectiveDate time.Time `json:"effective_date"`
	CreatedAt     time.Time `json:"created_at"`
}
