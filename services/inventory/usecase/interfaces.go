package usecase

import (
	"bisnis-rinzi/services/inventory/dto"
	"bisnis-rinzi/services/inventory/entity"
	"context"
	"io"
)

type InventoryUseCase interface {
	// Product Core & Custom Query
	CreateProduct(ctx context.Context, input dto.ProductCreateRequest) (string, error)
	UpdateProduct(ctx context.Context, id string, p *entity.Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetAllProducts(ctx context.Context) ([]*entity.Product, error)
	GetProductByID(ctx context.Context, id string) (*entity.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*entity.Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*entity.Product, error)
	GetLowStockProducts(ctx context.Context) ([]*entity.Product, error)

	// Tambahkan di dalam InventoryUseCase interface
	GetCatalogSyncData(ctx context.Context) (map[string]interface{}, error)
	GetAllStocksData(ctx context.Context) ([]map[string]interface{}, error)
	GetStockCard(ctx context.Context, productID string) ([]*entity.StockMovement, error)
	GetCostHistories(ctx context.Context, productID string) ([]*entity.CostHistory, error)
	GetAllMovements(ctx context.Context) ([]*entity.StockMovement, error)
	GetMovementByID(ctx context.Context, id string) (*entity.StockMovement, error)

	// Categories CRUD
	CreateCategory(ctx context.Context, code, name string) error
	UpdateCategory(ctx context.Context, id, code, name string) error
	DeleteCategory(ctx context.Context, id string) error
	GetAllCategories(ctx context.Context) ([]*entity.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*entity.Category, error)

	// Brands CRUD
	CreateBrand(ctx context.Context, code, name string) error
	UpdateBrand(ctx context.Context, id, code, name string) error
	DeleteBrand(ctx context.Context, id string) error
	GetAllBrands(ctx context.Context) ([]*entity.Brand, error)
	GetBrandByID(ctx context.Context, id string) (*entity.Brand, error)

	// Units CRUD
	CreateUnit(ctx context.Context, code, name string) error
	UpdateUnit(ctx context.Context, id, code, name string) error
	DeleteUnit(ctx context.Context, id string) error
	GetAllUnits(ctx context.Context) ([]*entity.Unit, error)
	GetUnitByID(ctx context.Context, id string) (*entity.Unit, error)

	// Cost History
	AddCostHistory(ctx context.Context, productID string, costPrice float64, reference string) error

	// Stock & Media & PWA Sync
	GetStock(ctx context.Context, productID string) (*entity.ProductStock, error)
	AdjustStock(ctx context.Context, input dto.StockAdjustRequest) error
	UpdateStockThresholds(ctx context.Context, productID string, input dto.StockThresholdRequest) error
	UploadProductMedia(ctx context.Context, productID, category string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error
	GetProductMedia(ctx context.Context, productID string) ([]*entity.ProductMedia, error)
	GetMediaItem(ctx context.Context, mediaID string) (*entity.ProductMedia, error)
	GetMediaStream(ctx context.Context, mediaID string) (io.ReadCloser, *entity.ProductMedia, error)
	DeleteProductMedia(ctx context.Context, mediaID string) error
	GetSyncData(ctx context.Context, clientVersion int64) (*dto.SyncChangesResponse, error)
	GetLatestSyncVersion(ctx context.Context) (int64, error)
	ImportProductsFromCSV(ctx context.Context, reader io.Reader) error
	ExportProductsToCSV(ctx context.Context, writer io.Writer) error
	SearchProducts(ctx context.Context, query string) ([]*entity.Product, error)
	InternalDeductStock(ctx context.Context, input dto.InternalStockRequest) error
	InternalRestoreStock(ctx context.Context, input dto.InternalStockRequest) error
}

type ProductRepository interface {
	Save(ctx context.Context, p *entity.Product, initialQty float64) error
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entity.Product, error)
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindBySKU(ctx context.Context, sku string) (*entity.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*entity.Product, error)
	FindLowStock(ctx context.Context) ([]*entity.Product, error)
	FindBySearch(ctx context.Context, query string) ([]*entity.Product, error) // Tambahan Baru

	// Stock & Movements
	GetStockByID(ctx context.Context, productID string) (*entity.ProductStock, error)
	AddStockMovement(ctx context.Context, m *entity.StockMovement) error
	AdjustStock(ctx context.Context, productID string, newQty float64, reference string) error
}
