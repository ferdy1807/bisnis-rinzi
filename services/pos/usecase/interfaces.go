package usecase

import (
	"bisnis-rinzi/services/pos/dto"
	"bisnis-rinzi/services/pos/entity"
	"context"
)

type POSUseCase interface {
	Checkout(ctx context.Context, cashierID string, input dto.CreateSaleRequest) (*entity.Sale, error)
	GetSalesHistory(ctx context.Context) ([]*entity.Sale, error)
	GetSaleDetail(ctx context.Context, id string) (*entity.Sale, error)
	GetSaleItems(ctx context.Context, saleID string) ([]*entity.SaleItem, error)
	GetInvoiceData(ctx context.Context, invoiceNum string) (*entity.Sale, error)
	GetReceipt(ctx context.Context, saleID string, cashierName string) (*dto.ReceiptResponse, error)
	GetTopProducts(ctx context.Context, limit int, sessionID *string) ([]*dto.TopProductResponse, error)
	GetProductSalesHistory(ctx context.Context, productID string) ([]*dto.ProductSalesHistoryResponse, error)

	// PWA Offline Sync Managers
	SyncOfflineTransactions(ctx context.Context, cashierID string, payload dto.OfflineSyncPayload) error
	RetryFailedSync(ctx context.Context, cashierID string, payload dto.OfflineSyncPayload) error
	GetFailedSyncLogs(ctx context.Context) ([]*entity.SyncLog, error)

	UploadInvoice(ctx context.Context, saleID string, cashierName string, fileBytes []byte, contentType string) (string, error)
}
