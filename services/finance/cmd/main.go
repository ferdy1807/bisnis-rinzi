package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	dbMinio "bisnis-rinzi/packages/backend/database/minio"
	"bisnis-rinzi/packages/backend/logger"
	delivery "bisnis-rinzi/services/finance/delivery/http"
	"bisnis-rinzi/services/finance/repository"
	"bisnis-rinzi/services/finance/routes"
	"bisnis-rinzi/services/finance/usecase"
	"bisnis-rinzi/services/finance/worker"
)

func main() {
	logger.InitLogger()
	logger.Info("Starting Finance Microservice...")

	// 1. Inisialisasi Postgres Pool khusus untuk finance_db
	dbDSN := os.Getenv("FINANCE_DB_URL")
	if dbDSN == "" {
		dbDSN = "postgres://postgres:postgres@localhost:5432/finance_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := postgres.Connect(ctx, postgres.Config{DSN: dbDSN, MaxConns: 25})
	if err != nil {
		log.Fatalf("Failed to connect to finance_db: %v", err)
	}
	defer dbClient.Close()

	// 2a. Inisialisasi MinIO
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
		log.Fatalf("Failed to connect to minio: %v", err)
	}

	// 2. Perakitan Komponen Clean Architecture
	financeRepo := repository.NewFinanceRepository(dbClient)
	financeUC := usecase.NewFinanceUseCase(financeRepo, minioClient, "laporan-income-harian")
	financeHandler := delivery.NewFinanceHandler(financeUC)

	// 3. Menjalankan Polling Outbox Ticker Berdurasi Presisi 5 Menit
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	// Menghubungkan DB Client Pool, Repository Akuntansi, Interval 5 Menit, dan Batch Size 100
	pollingInterval := 5 * time.Second
	batchSize := 100

	outboxWorker := worker.NewFinanceOutboxWorker(
		dbClient.Pool, // Mengambil basis *pgxpool.Pool dari dbClient wrapper
		financeRepo,   // Repo lokal yang mengimplementasikan interface FinanceJournalRepository
		pollingInterval,
		batchSize,
	)

	// Menjalankan worker di dalam background goroutine terisolasi
	go outboxWorker.Start(workerCtx)

	// Menjalankan Analytics Aggregator Worker (Interval 1 Menit)
	go worker.StartAnalyticsAggregator(workerCtx, financeRepo, 1*time.Minute)

	// 4. Konfigurasi Jaringan Server HTTP Multiplexer Internal
	mux := http.NewServeMux()
	routes.RegisterFinanceRoutes(mux, financeHandler)

	port := os.Getenv("FINANCE_SERVICE_PORT")
	if port == "" {
		port = "8086" // Sesuai Reverse Proxy API Gateway Anda
	}

	logger.Info("Finance Service is running on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Finance Service stopped: %v", err)
	}
}
