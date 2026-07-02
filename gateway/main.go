package main

import (
	"context"
	"net/http"
	"os"

	"bisnis-rinzi/gateway/routes"
	"bisnis-rinzi/packages/backend/database/redis"
	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/packages/backend/utils"
)

func main() {
	logger.InitLogger()
	logger.Info("Menginisialisasi API Gateway Bisnis-Rinzi...")

	// Muat .env dari root direktori proyek jika ada
	_ = utils.LoadEnv(".env")

	// Koneksi Redis untuk token blacklisting
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
		// Tetap berjalan walaupun Redis gagal, atau bisa juga dihentikan
	} else {
		defer redisClient.Close()
	}

	// Mendaftarkan seluruh rute proxy dan middleware
	handler := routes.RegisterRoutes(redisClient)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("API Gateway berhasil berjalan di port :%s 🚀", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, handler); err != nil {
		logger.Error("API Gateway berhenti secara tidak normal: %v", err)
	}
}
