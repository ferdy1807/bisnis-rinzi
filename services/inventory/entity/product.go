package entity

import "time"

// Product merepresentasikan model data dari tabel 'products'
type Product struct {
	ID           string     `json:"id"` // UUID
	SKU          string     `json:"sku"`
	CategoryID   string     `json:"category_id"`
	BrandID      *string    `json:"brand_id,omitempty"` // Nullable
	Name         string     `json:"name"`
	BaseUnitCode string     `json:"base_unit_code"`
	CostPrice    float64    `json:"cost_price"`
	SellingPrice float64    `json:"selling_price"`
	IsActive     bool       `json:"is_active"`
	Barcode      *string    `json:"barcode,omitempty"` // Nullable
	Qty          *float64        `json:"qty,omitempty"`
	Media        []*ProductMedia `json:"media,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
