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
	delivery "bisnis-rinzi/services/cash/delivery/http"
	"bisnis-rinzi/services/cash/repository"
	"bisnis-rinzi/services/cash/routes"
	"bisnis-rinzi/services/cash/usecase"
	"bisnis-rinzi/services/cash/worker"
)

func main() {
	logger.InitLogger()
	logger.Info("Starting Cash Microservice...")

	// 1. Inisialisasi Postgres Pool khusus untuk cash_db
	dbDSN := os.Getenv("CASH_DB_URL")
	if dbDSN == "" {
		dbDSN = "postgres://postgres:postgres@localhost:5432/cash_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := postgres.Connect(ctx, postgres.Config{DSN: dbDSN, MaxConns: 25})
	if err != nil {
		log.Fatalf("Failed to connect to cash_db: %v", err)
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
	cashRepo := repository.NewCashRepository(dbClient)
	cashUC := usecase.NewCashUseCase(cashRepo, minioClient, "laporan-tutup-shift")
	cashHandler := delivery.NewCashHandler(cashUC) // <-- Variabel dideklarasikan di sini

	// 3. Menjalankan Polling Outbox Ticker untuk Layanan Kas
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	cashOutboxWorker := worker.NewCashOutboxWorker(dbClient)
	go cashOutboxWorker.Start(workerCtx, 2*time.Second)

	// 4. Konfigurasi Jaringan Server HTTP Multiplexer Internal
	// SOLUSI UNDEFINED MUX: Melakukan inisialisasi instance multiplexer standar Go
	mux := http.NewServeMux()

	// SOLUSI DECLARED AND NOT USED: Mengonsumsi variabel 'cashHandler' ke dalam router registrasi
	routes.RegisterCashRoutes(mux, cashHandler)

	port := os.Getenv("CASH_SERVICE_PORT")
	if port == "" {
		port = "8083" // Port internal cash_service
	}

	logger.Info("Cash Service is running on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Cash Service stopped: %v", err)
	}
}
