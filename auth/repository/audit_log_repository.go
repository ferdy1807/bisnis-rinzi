package repository

import (
	"context"

	"bisnis-rinzi/auth/entity"
	"bisnis-rinzi/packages/backend/database/postgres"
)

type pgAuditLogRepository struct {
	db *postgres.DBClient
}

func NewAuditLogRepository(db *postgres.DBClient) AuditLogRepository {
	return &pgAuditLogRepository{db: db}
}

func (r *pgAuditLogRepository) SaveLog(ctx context.Context, l *entity.AuditLog) error {
	query := `INSERT INTO audit_logs (id, user_id, action, entity_name, entity_id, old_data, new_data, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Pool.Exec(ctx, query, l.ID, l.UserID, l.Action, l.EntityName, l.EntityID, l.OldData, l.NewData, l.CreatedAt, l.UpdatedAt)
	return err
}

func (r *pgAuditLogRepository) FindAll(ctx context.Context) ([]*entity.AuditLog, error) {
	query := `SELECT id, user_id, action, entity_name, entity_id, old_data, new_data, created_at FROM audit_logs ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		var l entity.AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.EntityName, &l.EntityID, &l.OldData, &l.NewData, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}

func (r *pgAuditLogRepository) FindByID(ctx context.Context, id string) (*entity.AuditLog, error) {
	query := `SELECT id, user_id, action, entity_name, entity_id, old_data, new_data, created_at FROM audit_logs WHERE id = $1`
	var l entity.AuditLog
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&l.ID, &l.UserID, &l.Action, &l.EntityName, &l.EntityID, &l.OldData, &l.NewData, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}
