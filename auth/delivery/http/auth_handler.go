package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"bisnis-rinzi/auth/dto"
	"bisnis-rinzi/auth/usecase"
	"bisnis-rinzi/packages/backend/response"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: uc}
}

// RegisterHandler menangani POST /api/auth/register
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
		return
	}

	if err := h.authUseCase.Register(r.Context(), req); err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "Registrasi pengguna berhasil", nil)
}

// LoginHandler menangani POST /api/auth/login
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
		return
	}

	// Menangkap metadata perangkat untuk jejak audit token
	req.IPAddress = r.RemoteAddr
	req.DeviceInfo = r.Header.Get("User-Agent")

	res, err := h.authUseCase.Login(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set Cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(8 * 3600), // 8 jam
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(7 * 24 * 3600), // 7 hari
	})

	response.WriteSuccess(w, http.StatusOK, "Login berhasil", res)
}

// ForceLogoutHandler menangani DELETE /api/auth/users/{id}/sessions
func (h *AuthHandler) ForceLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	// Ekstrak ID dari URL Path: /api/auth/users/{id}/sessions
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "ID Pengguna tidak valid")
		return
	}
	userID := pathParts[4] // Mengambil segmen {id}

	if err := h.authUseCase.ForceLogout(r.Context(), userID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Gagal melakukan force logout massal")
		return
	}

	// Hapus cookie sesi saat ini jika kebetulan user mem-force-logout dirinya sendiri
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response.WriteSuccess(w, http.StatusOK, "Seluruh sesi perangkat berhasil dicabut", nil)
}

func (h *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	// Mengambil ID pengguna yang disuntikkan oleh Gateway ke dalam Header
	userID := r.Header.Get("X-User-Id")
	userName := r.Header.Get("X-User-Name")
	userRole := r.Header.Get("X-User-Role")

	if userID == "" {
		response.WriteError(w, http.StatusUnauthorized, "Informasi pengguna tidak ditemukan dalam sesi")
		return
	}

	// Ambil data user dari database untuk mendapatkan full_name
	user, err := h.authUseCase.GetUserByID(r.Context(), userID)
	fullName := ""
	if err == nil && user != nil {
		fullName = user.FullName
	}

	// Menyusun data profil ringkas untuk dikembalikan ke klien/PWA
	profile := map[string]string{
		"id":        userID,
		"username":  userName,
		"role":      userRole,
		"full_name": fullName,
	}

	response.WriteSuccess(w, http.StatusOK, "Berhasil mengambil data profil pengguna", profile)
}

// RefreshTokenHandler menangani POST /api/auth/refresh-token
func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Coba dari body JSON
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)

	refreshToken := body["refresh_token"]
	if refreshToken == "" {
		// Fallback ke cookie
		if c, err := r.Cookie("refresh_token"); err == nil {
			refreshToken = c.Value
		}
	}

	newToken, err := h.authUseCase.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(8 * 3600), // 8 jam
	})

	response.WriteSuccess(w, http.StatusOK, "Token berhasil diperbarui", map[string]string{"access_token": newToken})
}

// LogoutHandler menangani POST /api/auth/logout
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		response.WriteError(w, http.StatusUnauthorized, "Sesi tidak valid atau header X-User-Id tidak ditemukan")
		return
	}

	authHeader := r.Header.Get("Authorization")
	var accessToken string
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			accessToken = parts[1]
		}
	} else {
		if c, err := r.Cookie("access_token"); err == nil {
			accessToken = c.Value
		}
	}

	var body map[string]string
	// Decode JSON namun abaikan error karena bisa saja klien cuma mengirim POST kosong
	_ = json.NewDecoder(r.Body).Decode(&body)

	refreshToken := body["refresh_token"]
	if refreshToken == "" {
		if c, err := r.Cookie("refresh_token"); err == nil {
			refreshToken = c.Value
		}
	}

	if err := h.authUseCase.Logout(r.Context(), userID, refreshToken, accessToken); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hapus cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response.WriteSuccess(w, http.StatusOK, "Berhasil logout dan mengakhiri sesi aktif", nil)
}

// UsersHandler menangani GET /api/auth/users dan POST /api/auth/users
func (h *AuthHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, _ := h.authUseCase.GetAllUsers(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar pengguna", users)
	case http.MethodPost:
		var req dto.RegisterRequest
		json.NewDecoder(r.Body).Decode(&req)
		h.authUseCase.Register(r.Context(), req)
		response.WriteSuccess(w, http.StatusCreated, "User berhasil dibuat", nil)
	default:
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
	}
}

// UserIDRouteHandler mengelola sub-jalur dinamis dinamis /api/auth/users/{id}/*
func (h *AuthHandler) UserIDRouteHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "ID salah")
		return
	}
	id := pathParts[4]

	if len(pathParts) == 5 { // /api/auth/users/{id}
		switch r.Method {
		case http.MethodGet:
			u, _ := h.authUseCase.GetUserByID(r.Context(), id)
			response.WriteSuccess(w, http.StatusOK, "Detail user", u)
		case http.MethodPut:
			var req dto.UpdateUserRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				response.WriteError(w, http.StatusBadRequest, "Payload JSON tidak valid")
				return
			}

			// Ambil data user yang lama dari database terlebih dahulu untuk mempertahankan nilai jika tidak dikirim
			existingUser, err := h.authUseCase.GetUserByID(r.Context(), id)
			if err != nil || existingUser == nil {
				response.WriteError(w, http.StatusNotFound, "User tidak ditemukan")
				return
			}

			// Jika field tidak dikirim di JSON, gunakan nilai lama dari database
			isActiveValue := existingUser.IsActive
			if req.IsActive != nil {
				isActiveValue = *req.IsActive
			}
			if req.Username == "" {
				req.Username = existingUser.Username
			}
			if req.FullName == "" {
				req.FullName = existingUser.FullName
			}
			if req.Role == "" {
				req.Role = existingUser.Role
			}

			// Kirim data yang sudah divalidasi ke Use Case
			err = h.authUseCase.UpdateUser(r.Context(), id, req.Username, req.FullName, req.Role, isActiveValue)
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}

			response.WriteSuccess(w, http.StatusOK, "User berhasil diperbarui", nil)
		case http.MethodDelete:
			h.authUseCase.DeleteUser(r.Context(), id)
			response.WriteSuccess(w, http.StatusOK, "User dihapus", nil)
		}
		return
	}

	// Sub-resource path: /api/auth/users/{id}/password atau /sessions
	subResource := pathParts[5]
	if r.Method == http.MethodPut && subResource == "password" {
		var b map[string]string
		json.NewDecoder(r.Body).Decode(&b)
		err := h.authUseCase.UpdatePassword(r.Context(), id, b["old_password"], b["new_password"])
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteSuccess(w, http.StatusOK, "Password berhasil diganti", nil)
	} else if r.Method == http.MethodDelete && subResource == "sessions" {
		h.ForceLogoutHandler(w, r)
	}
}

func (h *AuthHandler) RolesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		roles, _ := h.authUseCase.GetAllRoles(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Daftar Roles", roles)
	}
}

func (h *AuthHandler) RoleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	code := pathParts[4]
	var b map[string]string
	json.NewDecoder(r.Body).Decode(&b)
	h.authUseCase.UpdateRoleDashboard(r.Context(), code, b["dashboard_url"])
	response.WriteSuccess(w, http.StatusOK, "Dashboard URL berhasil diubah", nil)
}

func (h *AuthHandler) AuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) == 5 && pathParts[4] != "" { // /api/auth/audit-logs/{id}
		l, _ := h.authUseCase.GetAuditLogByID(r.Context(), pathParts[4])
		response.WriteSuccess(w, http.StatusOK, "Detail Audit Log", l)
		return
	}
	logs, _ := h.authUseCase.GetAuditLogs(r.Context())
	response.WriteSuccess(w, http.StatusOK, "Daftar Audit Logs", logs)
}

func (h *AuthHandler) TokensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := h.authUseCase.GetAllActiveTokens(r.Context())
		response.WriteSuccess(w, http.StatusOK, "Active Tokens", t)
	}
}

func (h *AuthHandler) TokenDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.WriteError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		response.WriteError(w, http.StatusBadRequest, "ID token tidak valid")
		return
	}
	err := h.authUseCase.RevokeTokenByID(r.Context(), pathParts[4])
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Gagal menghapus token secara permanen dari database")
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Token berhasil dihapus secara permanen oleh admin", nil)
}
