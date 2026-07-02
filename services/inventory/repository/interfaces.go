package repository

import (
	"bisnis-rinzi/services/inventory/entity"
	"context"

	"github.com/jackc/pgx/v5"
)

// ProductRepository mengelola query tabel products, product_stocks, dan cost_histories
type ProductRepository interface {
	Save(ctx context.Context, p *entity.Product, initialQty float64) error
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entity.Product, error)
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindBySKU(ctx context.Context, sku string) (*entity.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*entity.Product, error)
	FindBySearch(ctx context.Context, query string) ([]*entity.Product, error)
	FindLowStock(ctx context.Context) ([]*entity.Product, error)

	FindAllStocks(ctx context.Context) ([]map[string]interface{}, error)
	FindStockCardByProductID(ctx context.Context, productID string) ([]*entity.StockMovement, error)
	FindAllStockMovements(ctx context.Context) ([]*entity.StockMovement, error)
	FindStockMovementByID(ctx context.Context, id string) (*entity.StockMovement, error)
	FindCostHistoriesByProductID(ctx context.Context, productID string) ([]*entity.CostHistory, error)

	// Stock & Movements
	GetStockByID(ctx context.Context, productID string) (*entity.ProductStock, error)
	AddStockMovement(ctx context.Context, m *entity.StockMovement) error
	AdjustStock(ctx context.Context, productID string, newQty float64, reference string) error
	UpdateStockThresholds(ctx context.Context, productID string, minStock, safetyStock float64) error

	// Cost History
	AddCostHistory(ctx context.Context, ch *entity.CostHistory) error

	// Product Media
	SaveProductMedia(ctx context.Context, media *entity.ProductMedia) error
	FindMediaByProductID(ctx context.Context, productID string) ([]*entity.ProductMedia, error)
	FindMediaByID(ctx context.Context, mediaID string) (*entity.ProductMedia, error)
	DeleteMediaByID(ctx context.Context, mediaID string) error
}

// CategoryRepository mengelola query tabel categories
type CategoryRepository interface {
	Save(ctx context.Context, c *entity.Category) error
	Update(ctx context.Context, c *entity.Category) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entity.Category, error)
	FindByID(ctx context.Context, id string) (*entity.Category, error)
}

// BrandRepository mengelola query tabel brands
type BrandRepository interface {
	Save(ctx context.Context, b *entity.Brand) error
	Update(ctx context.Context, b *entity.Brand) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entity.Brand, error)
	FindByID(ctx context.Context, id string) (*entity.Brand, error)
}

// UnitRepository mengelola query tabel units
type UnitRepository interface {
	Save(ctx context.Context, u *entity.Unit) error
	Update(ctx context.Context, u *entity.Unit) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entity.Unit, error)
	FindByID(ctx context.Context, id string) (*entity.Unit, error)
}

// SyncRepository mengelola sinkronisasi data delta untuk PWA Offline
type SyncRepository interface {
	GetLatestVersion(ctx context.Context) (int64, error)
	GetChangesFromVersion(ctx context.Context, fromVersion int64) ([]*entity.SyncVersion, error)
	LogSyncVersion(ctx context.Context, tx pgx.Tx, entityType, entityID, operation string) error
	GetFullCatalogSync(ctx context.Context) (map[string]interface{}, error)
}
