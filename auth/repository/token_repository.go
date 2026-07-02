package repository

import (
	"context"
	"errors"
	"time"

	"bisnis-rinzi/auth/entity"
	"bisnis-rinzi/packages/backend/database/postgres"

	"github.com/jackc/pgx/v5"
)

type pgTokenRepository struct {
	db *postgres.DBClient
}

func NewTokenRepository(db *postgres.DBClient) TokenRepository {
	return &pgTokenRepository{db: db}
}

func (r *pgTokenRepository) SaveToken(ctx context.Context, t *entity.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at, device_info, ip_address, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.UserID, t.Token, t.ExpiresAt, t.DeviceInfo, t.IPAddress, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *pgTokenRepository) FindToken(ctx context.Context, tokenString string) (*entity.RefreshToken, error) {
	query := `SELECT id, user_id, token, expires_at, device_info, ip_address, created_at, updated_at, revoked_at 
	          FROM refresh_tokens WHERE token = $1`
	var t entity.RefreshToken
	err := r.db.Pool.QueryRow(ctx, query, tokenString).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.DeviceInfo, &t.IPAddress, &t.CreatedAt, &t.UpdatedAt, &t.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// DeleteAllSessionsByUserID Melakukan force logout massal pada seluruh perangkat user dengan menghapus dari database
func (r *pgTokenRepository) DeleteAllSessionsByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

func (r *pgTokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1, updated_at = $2 WHERE token = $3 AND revoked_at IS NULL`
	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, query, now, now, tokenString)
	return err
}

func (r *pgTokenRepository) DeleteTokenByID(ctx context.Context, id string) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *pgTokenRepository) FindAllActiveTokens(ctx context.Context) ([]*entity.RefreshToken, error) {
	query := `SELECT id, user_id, token, expires_at, device_info, ip_address, created_at, updated_at, revoked_at 
	          FROM refresh_tokens WHERE revoked_at IS NULL AND expires_at > $1`
	rows, err := r.db.Pool.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.RefreshToken
	for rows.Next() {
		var t entity.RefreshToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.DeviceInfo, &t.IPAddress, &t.CreatedAt, &t.UpdatedAt, &t.RevokedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, &t)
	}
	return tokens, nil
}
