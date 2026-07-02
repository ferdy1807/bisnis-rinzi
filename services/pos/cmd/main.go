package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	dbMinio "bisnis-rinzi/packages/backend/database/minio"
	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/logger"
	delivery "bisnis-rinzi/services/pos/delivery/http"
	"bisnis-rinzi/services/pos/repository"
	"bisnis-rinzi/services/pos/routes"
	"bisnis-rinzi/services/pos/usecase"
	"bisnis-rinzi/services/pos/worker"
)

func main() {
	logger.InitLogger()
	logger.Info("Starting POS Microservice...")

	// 1. Inisialisasi PostgreSQL pool khusus untuk pos_db [cite: 10]
	dbDSN := os.Getenv("POS_DB_URL")
	if dbDSN == "" {
		dbDSN = "postgres://postgres:postgres@localhost:5432/pos_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := postgres.Connect(ctx, postgres.Config{DSN: dbDSN, MaxConns: 20})
	if err != nil {
		log.Fatalf("Failed to connect to pos_db: %v", err)
	}
	defer dbClient.Close()

	// 2. Setup Object Storage Client (MinIO)
	minioClient, err := dbMinio.Connect(ctx, dbMinio.Config{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		UseSSL:          false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO Client: %v", err)
	}

	// 2.1 Ensure media bucket exists and is public
	if err := minioClient.CreateBucketIfNotExist(ctx, "invoice-toko", "us-east-1"); err != nil {
		logger.Error("Warning: failed to create bucket invoice-toko: %v", err)
	}
	if err := minioClient.MakeBucketPublic(ctx, "invoice-toko"); err != nil {
		logger.Error("Warning: failed to make bucket invoice-toko public: %v", err)
	}

	// 3. Perakitan Komponen Clean Architecture
	posRepo := repository.NewPOSRepository(dbClient)
	posUC := usecase.NewPOSUseCase(posRepo, minioClient, "invoice-toko")
	posHandler := delivery.NewPOSHandler(posUC)

	// 3. Hidupkan Background Worker Outbox dengan Durasi Ticker 5 Menit
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	posOutboxWorker := worker.NewOutboxWorker(dbClient)
	go posOutboxWorker.Start(workerCtx, 1*time.Second)

	// 4. Nyalakan HTTP Server Multiplexer Internal
	mux := http.NewServeMux()
	routes.RegisterPOSRoutes(mux, posHandler)

	port := os.Getenv("POS_SERVICE_PORT")
	if port == "" {
		port = "8084" // Port target Reverse Proxy API Gateway
	}

	logger.Info("POS Service is running on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("POS Service stopped: %v", err)
	}
}
