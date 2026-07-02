package routes

import (
	delivery "bisnis-rinzi/auth/delivery/http"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux, handler *delivery.AuthHandler) {
	mux.HandleFunc("/api/auth/register", handler.RegisterHandler)
	mux.HandleFunc("/api/auth/login", handler.LoginHandler)
	mux.HandleFunc("/api/auth/refresh-token", handler.RefreshTokenHandler)
	mux.HandleFunc("/api/auth/logout", handler.LogoutHandler)
	mux.HandleFunc("/api/auth/me", handler.MeHandler)

	mux.HandleFunc("/api/auth/users", handler.UsersHandler)
	mux.HandleFunc("/api/auth/users/", handler.UserIDRouteHandler)

	mux.HandleFunc("/api/auth/roles", handler.RolesHandler)
	mux.HandleFunc("/api/auth/roles/", handler.RoleUpdateHandler)

	mux.HandleFunc("/api/auth/audit-logs", handler.AuditLogsHandler)
	mux.HandleFunc("/api/auth/audit-logs/", handler.AuditLogsHandler)

	mux.HandleFunc("/api/auth/tokens", handler.TokensHandler)
	mux.HandleFunc("/api/auth/tokens/", handler.TokenDeleteHandler)
}
