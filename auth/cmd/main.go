package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	delivery "bisnis-rinzi/auth/delivery/http"
	"bisnis-rinzi/auth/repository"
	"bisnis-rinzi/auth/routes"
	"bisnis-rinzi/auth/usecase"
	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/database/redis"
	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/packages/backend/utils"
)

func main() {
	logger.InitLogger()
	logger.Info("Starting Auth Microservice...")

	// Muat .env dari root direktori proyek jika ada
	_ = utils.LoadEnv(".env")

	// 1. Inisialisasi Koneksi Database Terisolasi (auth_db)
	dsn := os.Getenv("AUTH_DB_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := postgres.Connect(ctx, postgres.Config{
		DSN:             dsn,
		MaxConns:        10,
		MinConns:        2,
		MaxConnIdleTime: 15 * time.Minute,
	})
	if err != nil {
		log.Fatalf("❌ Gagal menghubungkan ke database auth_db: %v", err)
	}
	defer dbClient.Close()

	// 2. Kunci JWT Secret Key
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "rinzi_secret_key_2026"
	}

	// 2.5 Koneksi Redis untuk token blacklisting
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisClient, err := redis.Connect(context.Background(), redis.Config{
		Host:     redisHost,
		Port:     6379,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err != nil {
		logger.Error("Gagal terhubung ke Redis: %v", err)
	} else {
		defer redisClient.Close()
	}

	// 3. Dependency Injection Semuanya (Clean Architecture)
	userRepo := repository.NewUserRepository(dbClient)
	tokenRepo := repository.NewTokenRepository(dbClient)
	roleRepo := repository.NewRoleRepository(dbClient)
	auditRepo := repository.NewAuditLogRepository(dbClient)
	authUC := usecase.NewAuthUseCase(userRepo, tokenRepo, roleRepo, auditRepo, jwtSecret, redisClient)
	authHandler := delivery.NewAuthHandler(authUC)

	// 4. Setup Jaringan Routing
	mux := http.NewServeMux()
	routes.RegisterAuthRoutes(mux, authHandler)

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081" // Port internal sesuai target routing API Gateway
	}

	logger.Info("Auth Service is running on port :%s", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, mux); err != nil {
		logger.Error("Auth Service stopped: %v", err)
	}
}
