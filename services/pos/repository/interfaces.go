package repository

import (
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/pos/dto"
	"bisnis-rinzi/services/pos/entity"
	"context"
)

type POSRepository interface {
	SaveTransaction(ctx context.Context, sale *entity.Sale, items []*entity.SaleItem, event *outbox.Event) error
	FindAll(ctx context.Context) ([]*entity.Sale, error)
	FindByID(ctx context.Context, id string) (*entity.Sale, error)
	FindByInvoiceNumber(ctx context.Context, invoiceNum string) (*entity.Sale, error)
	FindItemsBySaleID(ctx context.Context, saleID string) ([]*entity.SaleItem, error)
	CheckIdempotency(ctx context.Context, key string) (bool, error)
	GetTopProducts(ctx context.Context, limit int, sessionID *string) ([]*dto.TopProductResponse, error)
	FindSalesHistoryByProductID(ctx context.Context, productID string) ([]*dto.ProductSalesHistoryResponse, error)

	// Offline Sync Support
	SaveSyncLog(ctx context.Context, log *entity.SyncLog) error
	UpdateSyncLogStatus(ctx context.Context, entityID string, status string, errMsg *string) error
	FindFailedSyncLogs(ctx context.Context) ([]*entity.SyncLog, error)
	UpdateInvoiceURL(ctx context.Context, saleID string, url string) error

	// Integrasi Inventory
	GetProductCostPrice(ctx context.Context, productID string) (float64, error)
}
