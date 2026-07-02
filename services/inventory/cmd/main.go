package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	dbMinio "bisnis-rinzi/packages/backend/database/minio" // Alias pengujian sukses Anda
	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/logger"
	delivery "bisnis-rinzi/services/inventory/delivery/http"
	"bisnis-rinzi/services/inventory/repository"
	"bisnis-rinzi/services/inventory/routes"
	"bisnis-rinzi/services/inventory/usecase"
	"bisnis-rinzi/services/inventory/worker"
)

func main() {
	logger.InitLogger()
	logger.Info("Starting Inventory Microservice...")

	// 1. Setup Database Postgres (inventory_db)
	dbDSN := os.Getenv("INVENTORY_DB_URL")
	if dbDSN == "" {
		dbDSN = "postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := postgres.Connect(ctx, postgres.Config{DSN: dbDSN, MaxConns: 15})
	if err != nil {
		log.Fatalf("Failed to connect to inventory_db: %v", err)
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

	// 3. Perakitan Komponen Clean Architecture
	productRepo := repository.NewProductRepository(dbClient)
	categoryRepo := repository.NewCategoryRepository(dbClient)
	brandRepo := repository.NewBrandRepository(dbClient)
	unitRepo := repository.NewUnitRepository(dbClient)
	syncRepo := repository.NewSyncRepository(dbClient)

	inventoryUC := usecase.NewInventoryUseCase(productRepo, categoryRepo, brandRepo, unitRepo, syncRepo, minioClient, "foto-produk-toko")
	inventoryHandler := delivery.NewInventoryHandler(inventoryUC)

	// Inisialisasi Ticker Outbox Worker Inventory Service
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	inventoryOutboxWorker := worker.NewInventoryOutboxWorker(dbClient)
	go inventoryOutboxWorker.Start(workerCtx, 3*time.Second) // Polling interval 3 detik

	// 4. Inisialisasi HTTP Server
	mux := http.NewServeMux()
	routes.RegisterInventoryRoutes(mux, inventoryHandler)

	port := os.Getenv("INVENTORY_SERVICE_PORT")
	if port == "" {
		port = "8082" // Port target routing internal API Gateway
	}

	logger.Info("Inventory Service is running on port :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Inventory Service stopped: %v", err)
	}
}
