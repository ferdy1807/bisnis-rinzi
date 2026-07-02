package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	// Berikan alias eksplisit "dbMinio" di depan path import paket internal Anda
	dbMinio "bisnis-rinzi/packages/backend/database/minio"
	"bisnis-rinzi/packages/backend/logger"
	delivery "bisnis-rinzi/services/rental/delivery/http"
	"bisnis-rinzi/services/rental/repository"
	"bisnis-rinzi/services/rental/routes"
	"bisnis-rinzi/services/rental/usecase"
	"bisnis-rinzi/services/rental/worker"
)

func main() {
	logger.InitLogger()
	logger.Info("Starting Rental Microservice...")

	// 1. Inisialisasi Postgres Pool untuk rental_db
	dbDSN := os.Getenv("RENTAL_DB_URL")
	if dbDSN == "" {
		dbDSN = "postgres://postgres:postgres@localhost:5432/rental_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := postgres.Connect(ctx, postgres.Config{DSN: dbDSN, MaxConns: 15})
	if err != nil {
		log.Fatalf("Failed to connect to rental_db: %v", err)
	}
	defer dbClient.Close()

	// 2. Setup Object Storage Client (MinIO Cluster) menggunakan alias dbMinio yang valid
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioClient, err := dbMinio.Connect(ctx, dbMinio.Config{
		Endpoint:        minioEndpoint,
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		UseSSL:          false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO Client: %v", err)
	}

	// 3. Perakitan Komponen Clean Architecture Dengan Injeksi Ketergantungan MinIO
	rentalRepo := repository.NewRentalRepository(dbClient)
	rentalUC := usecase.NewRentalUseCase(rentalRepo, minioClient, "foto-produk-sewa")
	rentalHandler := delivery.NewRentalHandler(rentalUC)

	// 4. Menjalankan Polling Outbox Ticker Berdurasi Presisi 5 Menit
	// Inisialisasi Ticker Outbox Worker Rental Service
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	rentalOutboxWorker := worker.NewRentalOutboxWorker(dbClient)
	go rentalOutboxWorker.Start(workerCtx, 2*time.Second) // Polling interval 2 detik

	// 5. Konfigurasi Endpoint Multiplexer Jaringan Internal
	mux := http.NewServeMux()
	routes.RegisterRentalRoutes(mux, rentalHandler)

	port := os.Getenv("RENTAL_SERVICE_PORT")
	if port == "" {
		port = "8085"
	}

	logger.Info("Rental Service is running on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Rental Service stopped: %v", err)
	}
}
