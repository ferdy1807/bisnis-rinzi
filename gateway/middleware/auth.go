package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bisnis-rinzi/packages/backend/database/redis"
	"bisnis-rinzi/packages/backend/jwt"
	"bisnis-rinzi/packages/backend/response"
)

type contextKey string

const UserContextKey contextKey = "user_claims"

func AuthMiddleware(next http.Handler, redisClient *redis.RedisClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			cookie, err := r.Cookie("access_token")
			if err == nil && cookie.Value != "" {
				authHeader = "Bearer " + cookie.Value
			} else {
				response.WriteError(w, http.StatusUnauthorized, "Token otorisasi tidak ditemukan")
				return
			}
		}

		// Format token harus berupa: Bearer <token_string>
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.WriteError(w, http.StatusUnauthorized, "Format token otorisasi salah")
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "rinzi_secret_key_2026" // Fallback fallback untuk dev mode
		}

		// Validasi menggunakan shared package jwt
		tokenStr := tokenParts[1]
		claims, err := jwt.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, "Token tidak valid atau telah kadaluwarsa")
			return
		}

		if redisClient != nil && redisClient.Client != nil {
			ctxRedis := r.Context()

			// 1. Periksa apakah Access Token spesifik ini di-blacklist (Logout biasa)
			isBlacklisted, _ := redisClient.Client.Exists(ctxRedis, "blacklist:token:"+tokenStr).Result()
			if isBlacklisted > 0 {
				response.WriteError(w, http.StatusUnauthorized, "Sesi Anda telah berakhir. Silakan login kembali.")
				return
			}

			// 2. Periksa apakah ada Force Logout untuk user ini setelah token diterbitkan
			forceLogoutAtStr, _ := redisClient.Client.Get(ctxRedis, "user:"+claims.UserID+":force_logout_at").Result()
			if forceLogoutAtStr != "" && claims.IssuedAt != nil {
				forceLogoutAt, err := strconv.ParseInt(forceLogoutAtStr, 10, 64)
				if err == nil {
					// Jika token diterbitkan SEBELUM timestamp force logout, tolak aksesnya
					if claims.IssuedAt.Time.Unix() < forceLogoutAt {
						response.WriteError(w, http.StatusUnauthorized, "Sesi Anda telah dihentikan secara paksa. Silakan login kembali.")
						return
					}
				}
			}
		}

		// Menyisipkan data klaim user (ID, Username, Role) ke dalam context request
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
