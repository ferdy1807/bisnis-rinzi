package dto

import "time"

type SaleItemPayload struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitCode    string  `json:"unit_code"`
	Qty       float64 `json:"qty"`
	UnitPrice float64 `json:"unit_price"`
	Discount  float64 `json:"discount"`
}

type CreateSaleRequest struct {
	IdempotencyKey   string            `json:"idempotency_key"`
	PaymentMethod    string            `json:"payment_method"`
	Discount         float64           `json:"discount"`
	AmountPaid       float64           `json:"amount_paid"`
	CashierSessionID string            `json:"cashier_session_id"`
	Items            []SaleItemPayload `json:"items"`
}

type OfflineSyncPayload struct {
	Transactions []CreateSaleRequest `json:"transactions"`
}

type ReceiptResponse struct {
	InvoiceNumber   string            `json:"invoice_number"`
	TransactionDate time.Time         `json:"transaction_date"`
	PaymentMethod   string            `json:"payment_method"`
	Subtotal        float64           `json:"subtotal"`
	Discount        float64           `json:"discount"`
	Total           float64           `json:"total"`
	AmountPaid      float64           `json:"amount_paid"`
	ChangeAmount    float64           `json:"change_amount"`
	CashierName     string            `json:"cashier_name"`
	Items           []SaleItemPayload `json:"items"`
}

type TopProductResponse struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	TotalQty     float64 `json:"total_qty"`
	TotalRevenue float64 `json:"total_revenue"`
}

type ProductSalesHistoryResponse struct {
	TransactionDate time.Time `json:"transaction_date"`
	InvoiceNumber   string    `json:"invoice_number"`
	Qty             float64   `json:"qty"`
	UnitPrice       float64   `json:"unit_price"`
	Subtotal        float64   `json:"subtotal"`
}
