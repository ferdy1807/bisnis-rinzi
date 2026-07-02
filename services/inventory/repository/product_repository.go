package repository

import (
	"context"
	"errors"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/utils"
	"bisnis-rinzi/services/inventory/entity"

	"github.com/jackc/pgx/v5"
)

type pgProductRepository struct {
	db *postgres.DBClient
}

func NewProductRepository(db *postgres.DBClient) ProductRepository {
	return &pgProductRepository{db: db}
}

func (r *pgProductRepository) Save(ctx context.Context, p *entity.Product, initialQty float64) error {
	// Menjalankan ACID Transaction untuk tabel products dan product_stocks sekaligus
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert ke tabel products
	qProduct := `INSERT INTO products (id, sku, category_id, brand_id, name, base_unit_code, cost_price, selling_price, is_active, barcode, created_at, updated_at) 
	             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = tx.Exec(ctx, qProduct, p.ID, p.SKU, p.CategoryID, p.BrandID, p.Name, p.BaseUnitCode, p.CostPrice, p.SellingPrice, p.IsActive, p.Barcode, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	// 2. Insert ke tabel product_stocks
	qStock := `INSERT INTO product_stocks (product_id, qty, qty_min_stock, qty_safety_stock, created_at, updated_at) 
	           VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(ctx, qStock, p.ID, initialQty, 0, 0, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	// 2.5 (Baru) Insert ke tabel stock_movements jika initialQty > 0
	if initialQty > 0 {
		movementID := utils.GenerateUUIDv4()
		qMovement := `INSERT INTO stock_movements (id, product_id, movement_type, qty, reference, created_at, updated_at) 
		              VALUES ($1, $2, 'IN', $3, 'Stok Awal Produk', $4, $5)`
		_, err = tx.Exec(ctx, qMovement, movementID, p.ID, initialQty, p.CreatedAt, p.UpdatedAt)
		if err != nil {
			return err
		}
	}

	// 3. Catat ke tabel sync_versions untuk PWA Sync
	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('product', $1, 'INSERT', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, p.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgProductRepository) Update(ctx context.Context, p *entity.Product) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE products SET category_id=$1, brand_id=$2, name=$3, base_unit_code=$4, cost_price=$5, selling_price=$6, is_active=$7, barcode=$8, updated_at=$9 
	          WHERE id=$10 AND deleted_at IS NULL`
	_, err = tx.Exec(ctx, query, p.CategoryID, p.BrandID, p.Name, p.BaseUnitCode, p.CostPrice, p.SellingPrice, p.IsActive, p.Barcode, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('product', $1, 'UPDATE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, p.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgProductRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Soft delete
	query := `UPDATE products SET deleted_at = $1, is_active = false WHERE id = $2`
	_, err = tx.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('product', $1, 'DELETE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgProductRepository) FindAll(ctx context.Context) ([]*entity.Product, error) {
	query := `SELECT p.id, p.sku, p.category_id, p.brand_id, p.name, p.base_unit_code, p.cost_price, p.selling_price, p.is_active, p.barcode, p.created_at, p.updated_at, s.qty 
	          FROM products p LEFT JOIN product_stocks s ON p.id = s.product_id WHERE p.deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var p entity.Product
		err := rows.Scan(&p.ID, &p.SKU, &p.CategoryID, &p.BrandID, &p.Name, &p.BaseUnitCode, &p.CostPrice, &p.SellingPrice, &p.IsActive, &p.Barcode, &p.CreatedAt, &p.UpdatedAt, &p.Qty)
		if err != nil {
			return nil, err
		}
		p.Media = make([]*entity.ProductMedia, 0)
		products = append(products, &p)
	}

	// Ambil semua media untuk efisiensi n+1
	mediaQuery := `SELECT id, product_id, object_name FROM product_media WHERE is_active = true`
	mediaRows, err := r.db.Pool.Query(ctx, mediaQuery)
	if err == nil {
		defer mediaRows.Close()
		mediaMap := make(map[string][]*entity.ProductMedia)
		for mediaRows.Next() {
			var m entity.ProductMedia
			if err := mediaRows.Scan(&m.ID, &m.ProductID, &m.ObjectName); err == nil {
				mediaMap[m.ProductID] = append(mediaMap[m.ProductID], &m)
			}
		}
		for _, p := range products {
			if mediaList, exists := mediaMap[p.ID]; exists {
				p.Media = mediaList
			}
		}
	}

	return products, nil
}

func (r *pgProductRepository) FindBySearch(ctx context.Context, queryStr string) ([]*entity.Product, error) {
	query := `SELECT p.id, p.sku, p.category_id, p.brand_id, p.name, p.base_unit_code, p.cost_price, p.selling_price, p.is_active, p.barcode, p.created_at, p.updated_at, s.qty 
	          FROM products p LEFT JOIN product_stocks s ON p.id = s.product_id 
	          WHERE (p.name ILIKE '%' || $1 || '%' OR p.sku ILIKE '%' || $1 || '%' OR p.barcode ILIKE '%' || $1 || '%') 
	          AND p.deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var p entity.Product
		err := rows.Scan(&p.ID, &p.SKU, &p.CategoryID, &p.BrandID, &p.Name, &p.BaseUnitCode, &p.CostPrice, &p.SellingPrice, &p.IsActive, &p.Barcode, &p.CreatedAt, &p.UpdatedAt, &p.Qty)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *pgProductRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	query := `SELECT p.id, p.sku, p.category_id, p.brand_id, p.name, p.base_unit_code, p.cost_price, p.selling_price, p.is_active, p.barcode, p.created_at, p.updated_at, s.qty 
	          FROM products p LEFT JOIN product_stocks s ON p.id = s.product_id WHERE p.id = $1 AND p.deleted_at IS NULL`
	var p entity.Product
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.SKU, &p.CategoryID, &p.BrandID, &p.Name, &p.BaseUnitCode, &p.CostPrice, &p.SellingPrice, &p.IsActive, &p.Barcode, &p.CreatedAt, &p.UpdatedAt, &p.Qty)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgProductRepository) FindBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	query := `SELECT p.id, p.sku, p.category_id, p.brand_id, p.name, p.base_unit_code, p.cost_price, p.selling_price, p.is_active, p.barcode, p.created_at, p.updated_at, s.qty 
	          FROM products p LEFT JOIN product_stocks s ON p.id = s.product_id WHERE p.sku = $1 AND p.deleted_at IS NULL`
	var p entity.Product
	err := r.db.Pool.QueryRow(ctx, query, sku).Scan(&p.ID, &p.SKU, &p.CategoryID, &p.BrandID, &p.Name, &p.BaseUnitCode, &p.CostPrice, &p.SellingPrice, &p.IsActive, &p.Barcode, &p.CreatedAt, &p.UpdatedAt, &p.Qty)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgProductRepository) FindByBarcode(ctx context.Context, barcode string) (*entity.Product, error) {
	query := `SELECT p.id, p.sku, p.category_id, p.brand_id, p.name, p.base_unit_code, p.cost_price, p.selling_price, p.is_active, p.barcode, p.created_at, p.updated_at, s.qty 
	          FROM products p LEFT JOIN product_stocks s ON p.id = s.product_id WHERE p.barcode = $1 AND p.deleted_at IS NULL`
	var p entity.Product
	err := r.db.Pool.QueryRow(ctx, query, barcode).Scan(&p.ID, &p.SKU, &p.CategoryID, &p.BrandID, &p.Name, &p.BaseUnitCode, &p.CostPrice, &p.SellingPrice, &p.IsActive, &p.Barcode, &p.CreatedAt, &p.UpdatedAt, &p.Qty)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgProductRepository) FindLowStock(ctx context.Context) ([]*entity.Product, error) {
	query := `SELECT p.id, p.sku, p.category_id, p.brand_id, p.name, p.base_unit_code, p.cost_price, p.selling_price, p.is_active, p.barcode, p.created_at, p.updated_at, s.qty 
	          FROM products p
	          JOIN product_stocks s ON p.id = s.product_id
	          WHERE (s.qty <= s.qty_min_stock OR (s.qty_min_stock = 0 AND s.qty <= 5)) AND p.deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var p entity.Product
		err := rows.Scan(&p.ID, &p.SKU, &p.CategoryID, &p.BrandID, &p.Name, &p.BaseUnitCode, &p.CostPrice, &p.SellingPrice, &p.IsActive, &p.Barcode, &p.CreatedAt, &p.UpdatedAt, &p.Qty)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *pgProductRepository) GetStockByID(ctx context.Context, productID string) (*entity.ProductStock, error) {
	query := `SELECT product_id, qty, qty_min_stock, qty_safety_stock, created_at, updated_at FROM product_stocks WHERE product_id = $1`
	var s entity.ProductStock
	err := r.db.Pool.QueryRow(ctx, query, productID).Scan(&s.ProductID, &s.Qty, &s.QtyMinStock, &s.QtySafetyStock, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *pgProductRepository) AddStockMovement(ctx context.Context, m *entity.StockMovement) error {
	query := `INSERT INTO stock_movements (id, product_id, movement_type, qty, reference, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Pool.Exec(ctx, query, m.ID, m.ProductID, m.MovementType, m.Qty, m.Reference, m.CreatedAt, m.UpdatedAt)
	return err
}

// AdjustStock Mendukung endpoint POST /api/inventory/stocks/adjust
func (r *pgProductRepository) AdjustStock(ctx context.Context, productID string, newQty float64, reference string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Ambil stock lama untuk menghitung delta pencatatan movement
	var oldQty float64
	err = tx.QueryRow(ctx, "SELECT qty FROM product_stocks WHERE product_id = $1 FOR UPDATE", productID).Scan(&oldQty)
	if err != nil {
		return err
	}

	// 2. Update stock baru
	_, err = tx.Exec(ctx, "UPDATE product_stocks SET qty = $1, updated_at = $2 WHERE product_id = $3", newQty, time.Now(), productID)
	if err != nil {
		return err
	}

	// 3. Masukkan mutasi movement ke history
	deltaQty := newQty - oldQty
	mType := "ADJUSTMENT_IN"
	if deltaQty < 0 {
		mType = "ADJUSTMENT_OUT"
	}

	qMove := `INSERT INTO stock_movements (id, product_id, movement_type, qty, reference) VALUES (gen_random_uuid(), $1, $2, $3, $4)`
	_, err = tx.Exec(ctx, qMove, productID, mType, deltaQty, reference)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgProductRepository) FindAllStocks(ctx context.Context) ([]map[string]interface{}, error) {
	query := `SELECT p.id, p.sku, p.name, s.qty, s.qty_min_stock 
	          FROM products p 
	          JOIN product_stocks s ON p.id = s.product_id 
	          WHERE p.deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []map[string]interface{}
	for rows.Next() {
		var id, sku, name string
		var qty, minStock float64
		if err := rows.Scan(&id, &sku, &name, &qty, &minStock); err != nil {
			return nil, err
		}
		stocks = append(stocks, map[string]interface{}{
			"product_id":    id,
			"sku":           sku,
			"product_name":  name,
			"qty":           qty,
			"qty_min_stock": minStock,
		})
	}
	return stocks, nil
}

func (r *pgProductRepository) FindStockCardByProductID(ctx context.Context, productID string) ([]*entity.StockMovement, error) {
	query := `SELECT sm.id, sm.product_id, p.name, p.sku, sm.movement_type, sm.qty, sm.reference, sm.created_at, sm.updated_at 
	          FROM stock_movements sm
	          LEFT JOIN products p ON sm.product_id = p.id
	          WHERE sm.product_id = $1 
	          ORDER BY sm.created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*entity.StockMovement
	for rows.Next() {
		var m entity.StockMovement
		var ref *string
		var pName, pSku *string
		if err := rows.Scan(&m.ID, &m.ProductID, &pName, &pSku, &m.MovementType, &m.Qty, &ref, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if ref != nil {
			m.Reference = *ref
		}
		if pName != nil {
			m.ProductName = *pName
		}
		if pSku != nil {
			m.SKU = *pSku
		}
		movements = append(movements, &m)
	}
	return movements, nil
}

func (r *pgProductRepository) FindAllStockMovements(ctx context.Context) ([]*entity.StockMovement, error) {
	query := `SELECT sm.id, sm.product_id, p.name, p.sku, sm.movement_type, sm.qty, sm.reference, sm.created_at, sm.updated_at 
	          FROM stock_movements sm
	          LEFT JOIN products p ON sm.product_id = p.id
	          ORDER BY sm.created_at DESC LIMIT 100`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*entity.StockMovement
	for rows.Next() {
		var m entity.StockMovement
		var ref *string
		var pName, pSku *string
		if err := rows.Scan(&m.ID, &m.ProductID, &pName, &pSku, &m.MovementType, &m.Qty, &ref, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if ref != nil {
			m.Reference = *ref
		}
		if pName != nil {
			m.ProductName = *pName
		}
		if pSku != nil {
			m.SKU = *pSku
		}
		movements = append(movements, &m)
	}
	return movements, nil
}

func (r *pgProductRepository) FindStockMovementByID(ctx context.Context, id string) (*entity.StockMovement, error) {
	query := `SELECT sm.id, sm.product_id, p.name, p.sku, sm.movement_type, sm.qty, sm.reference, sm.created_at, sm.updated_at 
	          FROM stock_movements sm
	          LEFT JOIN products p ON sm.product_id = p.id
	          WHERE sm.id = $1`
	var m entity.StockMovement
	var ref *string
	var pName, pSku *string
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.ProductID, &pName, &pSku, &m.MovementType, &m.Qty, &ref, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if ref != nil {
		m.Reference = *ref
	}
	if pName != nil {
		m.ProductName = *pName
	}
	if pSku != nil {
		m.SKU = *pSku
	}
	return &m, nil
}

func (r *pgProductRepository) UpdateStockThresholds(ctx context.Context, productID string, minStock, safetyStock float64) error {
	query := `UPDATE product_stocks SET qty_min_stock = $1, qty_safety_stock = $2, updated_at = CURRENT_TIMESTAMP WHERE product_id = $3`
	_, err := r.db.Pool.Exec(ctx, query, minStock, safetyStock, productID)
	return err
}

// ==========================================
// COST HISTORY & PRODUCT MEDIA IMPLEMENTATION
// ==========================================

func (r *pgProductRepository) AddCostHistory(ctx context.Context, ch *entity.CostHistory) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO product_cost_histories (id, product_id, average_cost, effective_date, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_DATE, $4, $4)`
	_, err = tx.Exec(ctx, query, ch.ID, ch.ProductID, ch.AverageCost, ch.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgProductRepository) FindCostHistoriesByProductID(ctx context.Context, productID string) ([]*entity.CostHistory, error) {
	query := `SELECT id, product_id, average_cost, effective_date, created_at 
	          FROM product_cost_histories 
	          WHERE product_id = $1 
	          ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.CostHistory
	for rows.Next() {
		var ch entity.CostHistory
		if err := rows.Scan(&ch.ID, &ch.ProductID, &ch.AverageCost, &ch.EffectiveDate, &ch.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &ch)
	}
	return result, nil
}

func (r *pgProductRepository) SaveProductMedia(ctx context.Context, m *entity.ProductMedia) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO product_media (id, product_id, media_category, bucket_name, object_name, original_file_name, mime_type, file_size_bytes, is_active, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(ctx, query, m.ID, m.ProductID, m.MediaCategory, m.BucketName, m.ObjectName, m.OriginalFileName, m.MimeType, m.FileSizeValues, m.IsActive, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) VALUES ('product_media', $1, 'INSERT', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, m.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgProductRepository) FindMediaByProductID(ctx context.Context, productID string) ([]*entity.ProductMedia, error) {
	query := `SELECT id, product_id, media_category, bucket_name, object_name, original_file_name, mime_type, file_size_bytes, is_active, created_at, updated_at 
	          FROM product_media WHERE product_id = $1 AND is_active = true`
	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.ProductMedia
	for rows.Next() {
		var m entity.ProductMedia
		if err := rows.Scan(&m.ID, &m.ProductID, &m.MediaCategory, &m.BucketName, &m.ObjectName, &m.OriginalFileName, &m.MimeType, &m.FileSizeValues, &m.IsActive, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &m)
	}
	return result, nil
}

func (r *pgProductRepository) FindMediaByID(ctx context.Context, mediaID string) (*entity.ProductMedia, error) {
	query := `SELECT id, product_id, media_category, bucket_name, object_name, original_file_name, mime_type, file_size_bytes, is_active, created_at, updated_at 
	          FROM product_media WHERE id = $1`
	var m entity.ProductMedia
	err := r.db.Pool.QueryRow(ctx, query, mediaID).Scan(&m.ID, &m.ProductID, &m.MediaCategory, &m.BucketName, &m.ObjectName, &m.OriginalFileName, &m.MimeType, &m.FileSizeValues, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *pgProductRepository) DeleteMediaByID(ctx context.Context, mediaID string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM product_media WHERE id = $1`
	_, err = tx.Exec(ctx, query, mediaID)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) VALUES ('product_media', $1, 'DELETE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, mediaID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
