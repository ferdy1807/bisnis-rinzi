package repository

import (
	"bisnis-rinzi/auth/entity"
	"bisnis-rinzi/packages/backend/outbox"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, userID string, newHash string) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
}

type TokenRepository interface {
	SaveToken(ctx context.Context, token *entity.RefreshToken) error
	FindToken(ctx context.Context, tokenString string) (*entity.RefreshToken, error)
	RevokeToken(ctx context.Context, tokenString string) error
	DeleteTokenByID(ctx context.Context, id string) error
	FindAllActiveTokens(ctx context.Context) ([]*entity.RefreshToken, error)
	DeleteAllSessionsByUserID(ctx context.Context, userID string) error
}

type RoleRepository interface {
	FindAll(ctx context.Context) ([]*entity.Role, error)
	FindByCode(ctx context.Context, code string) (*entity.Role, error)
	UpdateDashboardURL(ctx context.Context, code string, url string) error
}

type AuditLogRepository interface {
	SaveLog(ctx context.Context, log *entity.AuditLog) error
	FindAll(ctx context.Context) ([]*entity.AuditLog, error)
	FindByID(ctx context.Context, id string) (*entity.AuditLog, error)
}

// OutboxRepository mengelola antrean event asinkron lintas service
type OutboxRepository interface {
	SaveEvent(ctx context.Context, tx pgx.Tx, event *outbox.Event) error
	GetPendingEvents(ctx context.Context, limit int) ([]*outbox.Event, error)
	UpdateEventStatus(ctx context.Context, id string, status string, errorMsg string) error
}
