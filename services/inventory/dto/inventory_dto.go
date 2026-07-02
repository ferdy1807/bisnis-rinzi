package dto

// ProductCreateRequest menampung payload untuk membuat produk dan stok awal
type ProductCreateRequest struct {
	SKU          string  `json:"sku"`
	CategoryID   string  `json:"category_id"`
	BrandID      *string `json:"brand_id,omitempty"`
	Name         string  `json:"name"`
	BaseUnitCode string  `json:"base_unit_code" validate:"required"`
	CostPrice    float64 `json:"cost_price"`
	SellingPrice float64 `json:"selling_price"`
	Barcode      *string `json:"barcode,omitempty"`
	InitialQty   float64 `json:"initial_qty"`
}

// StockAdjustRequest menampung payload untuk penyesuaian stok manual (Opname)
type StockAdjustRequest struct {
	ProductID string  `json:"product_id"`
	NewQty    float64 `json:"new_qty"`
	Reference string  `json:"reference"`
}

// StockThresholdRequest menampung payload untuk mengatur batas minimum & safety stock
type StockThresholdRequest struct {
	MinStock    float64 `json:"min_stock"`
	SafetyStock float64 `json:"safety_stock"`
}

// SyncChangesResponse mengembalikan daftar perubahan data delta untuk PWA
type SyncChangesResponse struct {
	LatestVersion int64         `json:"latest_version"`
	Changes       []interface{} `json:"changes"`
}

// InternalStockRequest menampung payload untuk deduksi atau restorasi stok dari service internal (misal: POS/Rental)
type InternalStockRequest struct {
	ProductID string  `json:"product_id"`
	Qty       float64 `json:"qty"`
	Reference string  `json:"reference"`
}
