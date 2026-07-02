package repository

import (
	"context"
	"errors"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/services/inventory/entity"

	"github.com/jackc/pgx/v5"
)

type pgBrandRepository struct {
	db *postgres.DBClient
}

func NewBrandRepository(db *postgres.DBClient) BrandRepository {
	return &pgBrandRepository{db: db}
}

func (r *pgBrandRepository) Save(ctx context.Context, b *entity.Brand) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO brands (id, code, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, query, b.ID, b.Code, b.Name, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('brand', $1, 'INSERT', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, b.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgBrandRepository) Update(ctx context.Context, b *entity.Brand) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE brands SET code = $1, name = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL`
	_, err = tx.Exec(ctx, query, b.Code, b.Name, b.UpdatedAt, b.ID)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('brand', $1, 'UPDATE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, b.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgBrandRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE brands SET deleted_at = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('brand', $1, 'DELETE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgBrandRepository) FindAll(ctx context.Context) ([]*entity.Brand, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM brands WHERE deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []*entity.Brand
	for rows.Next() {
		var b entity.Brand
		if err := rows.Scan(&b.ID, &b.Code, &b.Name, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		brands = append(brands, &b)
	}
	return brands, nil
}

func (r *pgBrandRepository) FindByID(ctx context.Context, id string) (*entity.Brand, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM brands WHERE id = $1 AND deleted_at IS NULL`
	var b entity.Brand
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.Code, &b.Name, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}
