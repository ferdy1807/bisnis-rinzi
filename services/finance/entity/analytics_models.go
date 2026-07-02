package entity

import "time"

type SalesAnalytics struct {
	Date              time.Time `json:"date"`
	TotalSalesAmount  float64   `json:"total_sales_amount"`
	TotalTransactions int       `json:"total_transactions"`
	AverageBasketSize float64   `json:"average_basket_size"`
}

type RentalAnalytics struct {
	Date              time.Time `json:"date"`
	TotalRentalAmount float64   `json:"total_rental_amount"`
	TotalReservations int       `json:"total_reservations"`
	ActiveRentedUnits int       `json:"active_rented_units"`
}

type ProductAnalytics struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	QtySold      float64 `json:"qty"`
	TotalRevenue float64 `json:"total"`
}

type CategoryAnalytics struct {
	CategoryName        string  `json:"category_name"`
	TotalRevenue        float64 `json:"total_revenue"`
	ContributionPercent float64 `json:"contribution_percent"`
}

type StockAnalytics struct {
	ProductID      string  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	CurrentStock   float64 `json:"current_stock"`
	StockValueCost float64 `json:"stock_value_cost"`
	TurnoverRate   float64 `json:"turnover_rate"`
}

type CashierAnalytics struct {
	CashierID          string  `json:"cashier_id"`
	CashierName        string  `json:"cashier_name"`
	TotalHandlingSales float64 `json:"total_handling_sales"`
	TotalDiscrepancies float64 `json:"total_discrepancies"`
}
