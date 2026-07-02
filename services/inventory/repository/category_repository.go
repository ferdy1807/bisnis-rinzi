package repository

import (
	"context"
	"errors"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/services/inventory/entity"

	"github.com/jackc/pgx/v5"
)

type pgCategoryRepository struct {
	db *postgres.DBClient
}

func NewCategoryRepository(db *postgres.DBClient) CategoryRepository {
	return &pgCategoryRepository{db: db}
}

func (r *pgCategoryRepository) Save(ctx context.Context, c *entity.Category) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO categories (id, code, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, query, c.ID, c.Code, c.Name, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('category', $1, 'INSERT', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, c.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgCategoryRepository) Update(ctx context.Context, c *entity.Category) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE categories SET code = $1, name = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL`
	_, err = tx.Exec(ctx, query, c.Code, c.Name, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('category', $1, 'UPDATE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, c.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgCategoryRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE categories SET deleted_at = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('category', $1, 'DELETE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgCategoryRepository) FindAll(ctx context.Context) ([]*entity.Category, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM categories WHERE deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *pgCategoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM categories WHERE id = $1 AND deleted_at IS NULL`
	var c entity.Category
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.Code, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
