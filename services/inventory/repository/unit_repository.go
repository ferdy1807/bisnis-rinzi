package repository

import (
	"context"
	"errors"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/services/inventory/entity"

	"github.com/jackc/pgx/v5"
)

type pgUnitRepository struct {
	db *postgres.DBClient
}

func NewUnitRepository(db *postgres.DBClient) UnitRepository {
	return &pgUnitRepository{db: db}
}

func (r *pgUnitRepository) Save(ctx context.Context, u *entity.Unit) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO units (id, code, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, query, u.ID, u.Code, u.Name, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('unit', $1, 'INSERT', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, u.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgUnitRepository) Update(ctx context.Context, u *entity.Unit) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE units SET code = $2, name = $3, updated_at = $4 WHERE id = $1`
	_, err = tx.Exec(ctx, query, u.ID, u.Code, u.Name, u.UpdatedAt)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('unit', $1, 'UPDATE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, u.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgUnitRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Units di tabel rancangan tidak ada deleted_at, langsung hard-delete
	query := `DELETE FROM units WHERE id = $1`
	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	qSync := `INSERT INTO sync_versions (entity_type, entity_id, operation, version_number) 
	          VALUES ('unit', $1, 'DELETE', nextval('sync_global_version_seq'))`
	_, err = tx.Exec(ctx, qSync, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgUnitRepository) FindAll(ctx context.Context) ([]*entity.Unit, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM units`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*entity.Unit
	for rows.Next() {
		var u entity.Unit
		if err := rows.Scan(&u.ID, &u.Code, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		units = append(units, &u)
	}
	return units, nil
}

func (r *pgUnitRepository) FindByID(ctx context.Context, id string) (*entity.Unit, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM units WHERE id = $1`
	var u entity.Unit
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Code, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
