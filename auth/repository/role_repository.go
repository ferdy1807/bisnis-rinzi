package repository

import (
	"bisnis-rinzi/auth/entity"
	"bisnis-rinzi/packages/backend/database/postgres"
	"context"
	"time"
)

type pgRoleRepository struct {
	db *postgres.DBClient
}

func NewRoleRepository(db *postgres.DBClient) RoleRepository {
	return &pgRoleRepository{db: db}
}

func (r *pgRoleRepository) FindAll(ctx context.Context) ([]*entity.Role, error) {
	query := `SELECT code, name, dashboard_url, created_at, updated_at FROM roles`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		var rl entity.Role
		if err := rows.Scan(&rl.Code, &rl.Name, &rl.DashboardURL, &rl.CreatedAt, &rl.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, &rl)
	}
	return roles, nil
}

func (r *pgRoleRepository) FindByCode(ctx context.Context, code string) (*entity.Role, error) {
	query := `SELECT code, name, dashboard_url, created_at, updated_at FROM roles WHERE code = $1`
	var rl entity.Role
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(&rl.Code, &rl.Name, &rl.DashboardURL, &rl.CreatedAt, &rl.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rl, nil
}

func (r *pgRoleRepository) UpdateDashboardURL(ctx context.Context, code string, url string) error {
	query := `UPDATE roles SET dashboard_url = $1, updated_at = $2 WHERE code = $3`
	_, err := r.db.Pool.Exec(ctx, query, url, time.Now(), code)
	return err
}
