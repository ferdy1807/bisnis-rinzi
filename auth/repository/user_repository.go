package repository

import (
	"context"
	"errors"
	"time"

	"bisnis-rinzi/auth/entity"
	"bisnis-rinzi/packages/backend/database/postgres"

	"github.com/jackc/pgx/v5"
)

type pgUserRepository struct {
	db *postgres.DBClient
}

func NewUserRepository(db *postgres.DBClient) UserRepository {
	return &pgUserRepository{db: db}
}

func (r *pgUserRepository) Save(ctx context.Context, u *entity.User) error {
	query := `INSERT INTO users (id, username, password_hash, full_name, role, is_active, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Pool.Exec(ctx, query, u.ID, u.Username, u.PasswordHash, u.FullName, u.Role, u.IsActive, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *pgUserRepository) Update(ctx context.Context, u *entity.User) error {
	query := `UPDATE users SET username = $1, full_name = $2, role = $3, is_active = $4, updated_at = $5 WHERE id = $6 AND deleted_at IS NULL`
	_, err := r.db.Pool.Exec(ctx, query, u.Username, u.FullName, u.Role, u.IsActive, u.UpdatedAt, u.ID)
	return err
}

func (r *pgUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, username, password_hash, full_name, role, is_active, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL`
	var u entity.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, username, password_hash, full_name, role, is_active, created_at, updated_at FROM users WHERE username = $1 AND deleted_at IS NULL`
	var u entity.User
	err := r.db.Pool.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) UpdatePassword(ctx context.Context, userID string, newHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`
	_, err := r.db.Pool.Exec(ctx, query, newHash, time.Now(), userID)
	return err
}

func (r *pgUserRepository) Delete(ctx context.Context, id string) error {
	// Menerapkan Soft Delete demi menjaga validitas integritas data historis
	query := `UPDATE users SET deleted_at = $1, is_active = false WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *pgUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	query := `SELECT id, username, full_name, role, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Username, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
