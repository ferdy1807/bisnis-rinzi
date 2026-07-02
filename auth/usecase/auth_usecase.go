package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"bisnis-rinzi/auth/dto"
	"bisnis-rinzi/auth/entity"
	"bisnis-rinzi/auth/repository"
	"bisnis-rinzi/packages/backend/database/redis"
	"bisnis-rinzi/packages/backend/jwt"
	"bisnis-rinzi/packages/backend/utils"
)

type authUseCase struct {
	userRepo    repository.UserRepository
	tokenRepo   repository.TokenRepository
	roleRepo    repository.RoleRepository
	auditRepo   repository.AuditLogRepository
	jwtSecret   string
	redisClient *redis.RedisClient
}

func NewAuthUseCase(ur repository.UserRepository, tr repository.TokenRepository, rr repository.RoleRepository, ar repository.AuditLogRepository, secret string, rClient *redis.RedisClient) AuthUseCase {
	return &authUseCase{
		userRepo:    ur,
		tokenRepo:   tr,
		roleRepo:    rr,
		auditRepo:   ar,
		jwtSecret:   secret,
		redisClient: rClient,
	}
}

// Register menangani pembuatan akun user baru dengan enkripsi password Bcrypt
func (u *authUseCase) Register(ctx context.Context, input dto.RegisterRequest) error {
	// 1. Validasi ketersediaan Username unik
	existingUser, err := u.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username sudah digunakan oleh pegawai lain")
	}

	// 2. Hash Password menggunakan shared package utils
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("gagal memproses enkripsi password: %w", err)
	}

	// 3. Bangun objek entitas domain baru
	now := time.Now()
	newUser := &entity.User{
		ID:           utils.GenerateUUIDv4(), // Di lapangan, backend bisa auto-generate lewat db atau uuid library
		Username:     input.Username,
		PasswordHash: hashedPassword,
		FullName:     input.FullName,
		Role:         input.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 4. Validasi Integritas Data Aturan Bisnis
	if err := newUser.Validate(); err != nil {
		return err
	}

	return u.userRepo.Save(ctx, newUser)
}

// Login mematangkan pencocokan password dan melahirkan token akses JWT terproteksi
func (u *authUseCase) Login(ctx context.Context, input dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Cari user berdasarkan username
	user, err := u.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsActive {
		return nil, errors.New("kredensial login salah atau akun dinonaktifkan")
	}

	// 2. Verifikasi Password Bcrypt
	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, errors.New("kredensial login salah")
	}

	// 3. Generate Access Token JWT (Berlaku 15 Menit) via Shared Package
	accessToken, err := jwt.GenerateToken(user.ID, user.Username, user.Role, u.jwtSecret, 1*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("gagal merilis access token: %w", err)
	}

	// 4. Generate Opaque Refresh Token (Sesi Panjang)
	refreshTokenStr, err := generateSecureRandomString(32)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	refreshTokenEntity := &entity.RefreshToken{
		ID:         utils.GenerateUUIDv4(),
		UserID:     user.ID,
		Token:      refreshTokenStr,
		ExpiresAt:  now.Add(1 * time.Hour), // Berlaku 7 Hari
		DeviceInfo: input.DeviceInfo,
		IPAddress:  input.IPAddress,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 5. Simpan Refresh Token ke database auth_db
	if err := u.tokenRepo.SaveToken(ctx, refreshTokenEntity); err != nil {
		return nil, fmt.Errorf("gagal mengamankan token sesi: %w", err)
	}

	// 6. Mapping Dashboard URL Fallback
	dashboardURL := "/portal-toko"
	roleEntity, errRole := u.roleRepo.FindByCode(ctx, user.Role)
	if errRole == nil && roleEntity != nil && roleEntity.DashboardURL != "" {
		dashboardURL = roleEntity.DashboardURL
	} else {
		// Fallback manual jika db tidak sinkron
		switch user.Role {
		case "OWNER":
			dashboardURL = "/admin-dashboard/"
		case "PEGAWAI":
			dashboardURL = "/portal-sewa/"
		default:
			dashboardURL = "/portal-toko/"
		}
	}

	// 7. Catat Audit Log
	_ = u.auditRepo.SaveLog(ctx, &entity.AuditLog{
		ID:         utils.GenerateUUIDv4(),
		UserID:     user.ID,
		Action:     "USER_LOGIN",
		EntityName: "users",
		EntityID:   user.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		DashboardURL: dashboardURL,
	}, nil
}

// ForceLogout mematikan hak akses seluruh sesi perangkat dari user terkait
func (u *authUseCase) ForceLogout(ctx context.Context, userID string) error {
	// 1. Catat force logout di Redis agar semua access token yg lama tertolak
	if u.redisClient != nil && u.redisClient.Client != nil {
		nowStr := fmt.Sprintf("%d", time.Now().Unix())
		// TTL untuk force logout bisa diatur sepanjang max umur token (misal 15 menit) atau dibiarkan tanpa TTL (atau TTL panjang)
		u.redisClient.Client.Set(ctx, "user:"+userID+":force_logout_at", nowStr, 24*time.Hour)
	}

	// 2. Physical Delete semua refresh tokens di DB
	// Menggunakan DeleteAllSessionsByUserID yang akan diperbarui implementasinya (atau jika repo masih soft delete, kita ganti perintahnya)
	// Namun agar lebih clean, mari panggil DeleteAllSessionsByUserID, kita akan modifikasi repo-nya menjadi physical delete.
	err := u.tokenRepo.DeleteAllSessionsByUserID(ctx, userID)

	now := time.Now()
	_ = u.auditRepo.SaveLog(ctx, &entity.AuditLog{
		ID:         utils.GenerateUUIDv4(),
		UserID:     userID,
		Action:     "USER_FORCE_LOGOUT",
		EntityName: "users",
		EntityID:   userID,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	return err
}

// Helper internal untuk kebutuhan keunikan biner id token/user dummy generator
func generateSecureRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (u *authUseCase) RefreshToken(ctx context.Context, rToken string) (string, error) {
	tEntity, err := u.tokenRepo.FindToken(ctx, rToken)
	if err != nil || tEntity == nil || !tEntity.IsActive() {
		return "", errors.New("refresh token tidak valid atau sudah dicabut")
	}
	user, err := u.userRepo.FindByID(ctx, tEntity.UserID)
	if err != nil || user == nil || !user.IsActive {
		return "", errors.New("pengenang sesi tidak valid")
	}
	return jwt.GenerateToken(user.ID, user.Username, user.Role, u.jwtSecret, 1*time.Hour)
}

func (u *authUseCase) Logout(ctx context.Context, userID string, rToken string, accessToken string) error {
	tEntity, err := u.tokenRepo.FindToken(ctx, rToken)
	if err != nil || tEntity == nil {
		return errors.New("refresh token tidak valid")
	}

	// Verifikasi apakah token benar-benar milik user yang me-request logout
	if tEntity.UserID != userID {
		return errors.New("anda tidak memiliki izin untuk mengakhiri sesi ini")
	}

	// 1. Blacklist Access Token di Redis
	if u.redisClient != nil && u.redisClient.Client != nil && accessToken != "" {
		// TTL diset 15 menit (maksimal umur access token)
		u.redisClient.Client.Set(ctx, "blacklist:token:"+accessToken, "revoked", 15*time.Minute)
	}

	// 2. Tandai token sebagai revoked (soft delete) di database
	err = u.tokenRepo.RevokeToken(ctx, tEntity.Token)

	now := time.Now()
	_ = u.auditRepo.SaveLog(ctx, &entity.AuditLog{
		ID:         utils.GenerateUUIDv4(),
		UserID:     userID,
		Action:     "USER_LOGOUT",
		EntityName: "users",
		EntityID:   userID,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	return err
}

func (u *authUseCase) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.FindAll(ctx)
}

func (u *authUseCase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *authUseCase) UpdateUser(ctx context.Context, id string, username, fullName, role string, isActive bool) error {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan")
	}
	user.Username = username
	user.FullName = fullName
	user.Role = role
	user.IsActive = isActive
	user.UpdatedAt = time.Now()
	return u.userRepo.Update(ctx, user)
}

func (u *authUseCase) UpdatePassword(ctx context.Context, id string, oldPassword, newPassword string) error {
	// 1. Cari user berdasarkan ID terlebih dahulu untuk mendapatkan username-nya
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan")
	}

	// 2. Ambil data user lengkap (termasuk password_hash) menggunakan FindByUsername
	dbUser, err := u.userRepo.FindByUsername(ctx, user.Username)
	if err != nil || dbUser == nil {
		return errors.New("gagal memvalidasi data internal user")
	}

	// 3. Lakukan komparasi Bcrypt hash yang valid
	if !utils.CheckPasswordHash(oldPassword, dbUser.PasswordHash) {
		return errors.New("password lama tidak sesuai")
	}

	// 4. Enkripsi password baru
	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 5. Simpan ke database
	return u.userRepo.UpdatePassword(ctx, id, newHash)
}

func (u *authUseCase) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *authUseCase) GetAllRoles(ctx context.Context) ([]*entity.Role, error) {
	return u.roleRepo.FindAll(ctx)
}

func (u *authUseCase) UpdateRoleDashboard(ctx context.Context, code string, url string) error {
	return u.roleRepo.UpdateDashboardURL(ctx, code, url)
}

func (u *authUseCase) GetAuditLogs(ctx context.Context) ([]*entity.AuditLog, error) {
	return u.auditRepo.FindAll(ctx)
}

func (u *authUseCase) GetAuditLogByID(ctx context.Context, id string) (*entity.AuditLog, error) {
	return u.auditRepo.FindByID(ctx, id)
}

func (u *authUseCase) GetAllActiveTokens(ctx context.Context) ([]*entity.RefreshToken, error) {
	return u.tokenRepo.FindAllActiveTokens(ctx)
}

func (u *authUseCase) RevokeTokenByID(ctx context.Context, id string) error {
	return u.tokenRepo.DeleteTokenByID(ctx, id)
}
