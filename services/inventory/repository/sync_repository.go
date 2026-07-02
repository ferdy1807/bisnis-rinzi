package repository

import (
	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/services/inventory/entity"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type pgSyncRepository struct {
	db *postgres.DBClient
}

func NewSyncRepository(db *postgres.DBClient) SyncRepository {
	return &pgSyncRepository{db: db}
}

func (r *pgSyncRepository) GetLatestVersion(ctx context.Context) (int64, error) {
	// Mengambil nomor urut versi monotonik global tertinggi saat ini
	query := `SELECT COALESCE(MAX(version_number), 0) FROM sync_versions`
	var currentVer int64
	err := r.db.Pool.QueryRow(ctx, query).Scan(&currentVer)
	return currentVer, err
}

func (r *pgSyncRepository) GetChangesFromVersion(ctx context.Context, fromVersion int64) ([]*entity.SyncVersion, error) {
	// Menarik data delta (hanya data yang berubah dari versi terakhir milik PWA klien)
	query := `SELECT id, entity_type, entity_id, operation, version_number, changed_at 
	          FROM sync_versions WHERE version_number > $1 ORDER BY version_number ASC`
	rows, err := r.db.Pool.Query(ctx, query, fromVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []*entity.SyncVersion
	for rows.Next() {
		var sv entity.SyncVersion
		err := rows.Scan(&sv.ID, &sv.EntityType, &sv.EntityID, &sv.Operation, &sv.VersionNumber, &sv.ChangedAt)
		if err != nil {
			return nil, err
		}
		changes = append(changes, &sv)
	}
	return changes, nil
}

func (r *pgSyncRepository) LogSyncVersion(ctx context.Context, tx pgx.Tx, entityType, entityID, operation string) error {
	query := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ($1, $2, $3, nextval('sync_global_version_seq'))`
	_, err := tx.Exec(ctx, query, entityType, entityID, operation)
	return err
}

func (r *pgSyncRepository) GetFullCatalogSync(ctx context.Context) (map[string]interface{}, error) {
	// 1. Ambil data categories
	cRows, err := r.db.Pool.Query(ctx, "SELECT id, code, name FROM categories WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer cRows.Close()

	var categories []map[string]string
	for cRows.Next() {
		var id, code, name string
		_ = cRows.Scan(&id, &code, &name)
		categories = append(categories, map[string]string{"id": id, "code": code, "name": name})
	}

	// 2. Ambil data products teraktif
	pRows, err := r.db.Pool.Query(ctx, "SELECT id, sku, category_id, name, base_unit_code, selling_price, barcode FROM products WHERE deleted_at IS NULL AND is_active = true")
	if err != nil {
		return nil, err
	}
	defer pRows.Close()

	var products []map[string]interface{}
	for pRows.Next() {
		var id, sku, catID, name, unit string
		var price float64
		var bar *string
		_ = pRows.Scan(&id, &sku, &catID, &name, &unit, &price, &bar)

		barcodeVal := ""
		if bar != nil {
			barcodeVal = *bar
		}

		products = append(products, map[string]interface{}{
			"id":             id,
			"sku":            sku,
			"category_id":    catID,
			"name":           name,
			"base_unit_code": unit,
			"selling_price":  price,
			"barcode":        barcodeVal,
		})
	}

	return map[string]interface{}{
		"categories": categories,
		"products":   products,
		"synced_at":  time.Now(),
	}, nil
}
