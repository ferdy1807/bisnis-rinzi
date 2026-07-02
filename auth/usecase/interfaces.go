package usecase

import (
	"bisnis-rinzi/auth/dto"
	"bisnis-rinzi/auth/entity"
	"context"
)

type AuthUseCase interface {
	Register(ctx context.Context, input dto.RegisterRequest) error
	Login(ctx context.Context, input dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, userID string, refreshToken string, accessToken string) error

	// Manajemen User
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	UpdateUser(ctx context.Context, id string, username, fullName, role string, isActive bool) error
	UpdatePassword(ctx context.Context, id string, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id string) error
	ForceLogout(ctx context.Context, userID string) error

	// Roles & Audit & Tokens Lintas User
	GetAllRoles(ctx context.Context) ([]*entity.Role, error)
	UpdateRoleDashboard(ctx context.Context, code string, url string) error
	GetAuditLogs(ctx context.Context) ([]*entity.AuditLog, error)
	GetAuditLogByID(ctx context.Context, id string) (*entity.AuditLog, error)
	GetAllActiveTokens(ctx context.Context) ([]*entity.RefreshToken, error)
	RevokeTokenByID(ctx context.Context, id string) error
}
