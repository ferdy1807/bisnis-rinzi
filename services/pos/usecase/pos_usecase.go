package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	dbMinio "bisnis-rinzi/packages/backend/database/minio"
	"bisnis-rinzi/packages/backend/logger"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/packages/backend/utils"
	"bisnis-rinzi/services/pos/dto"
	"bisnis-rinzi/services/pos/entity"
	"bisnis-rinzi/services/pos/repository"
)

type posUseCase struct {
	posRepo     repository.POSRepository
	minioClient *dbMinio.MinioClient
	bucketName  string
}

func NewPOSUseCase(repo repository.POSRepository, minioClient *dbMinio.MinioClient, bucketName string) POSUseCase {
	return &posUseCase{
		posRepo:     repo,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (u *posUseCase) Checkout(ctx context.Context, cashierID string, input dto.CreateSaleRequest) (*entity.Sale, error) {
	// 1. Cek Idempotency Key untuk menjamin integritas data terhindar dari double submit
	isDuplicate, err := u.posRepo.CheckIdempotency(ctx, input.IdempotencyKey)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("transaksi sudah pernah diproses (idempotent)")
	}

	// 2. Hitung Ulang Total Harga di Backend untuk menghindari tampering
	var subtotal float64
	var items []*entity.SaleItem
	now := time.Now()
	saleID := utils.GenerateUUIDv4()

	for _, itemInput := range input.Items {
		itemSubtotal := (itemInput.Qty * itemInput.UnitPrice) - itemInput.Discount
		subtotal += itemSubtotal

		// Mengambil harga modal sesungguhnya dari inventory lewat pos_repository
		costPrice, err := u.posRepo.GetProductCostPrice(ctx, itemInput.ProductID)
		if err != nil {
			logger.Error("Gagal mengambil harga modal untuk produk %s: %v", itemInput.ProductID, err)
			costPrice = 0 // Fallback ke 0 jika gagal (misal produk dihapus sementara)
		}

		items = append(items, &entity.SaleItem{
			ID:          utils.GenerateUUIDv4(),
			SaleID:      saleID,
			ProductID:   itemInput.ProductID,
			ProductName: itemInput.ProductName,
			UnitCode:    itemInput.UnitCode,
			Qty:         itemInput.Qty,
			UnitPrice:   itemInput.UnitPrice,
			Discount:    itemInput.Discount,
			Subtotal:    itemSubtotal,
			CostPrice:   costPrice, // Menggunakan harga modal riil dari inventory_db
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	totalAmount := subtotal - input.Discount
	invoiceNumber := fmt.Sprintf("INV/RETAIL/%d/%s", now.Unix(), input.IdempotencyKey[:6])

	amountPaid := input.AmountPaid
	if input.PaymentMethod != "CASH" || amountPaid < totalAmount {
		amountPaid = totalAmount // Jika bukan tunai, atau kurang (misal tidak diisi), samakan dengan total
	}
	changeAmount := amountPaid - totalAmount

	saleEntity := &entity.Sale{
		ID:               saleID,
		IdempotencyKey:   input.IdempotencyKey,
		InvoiceNumber:    invoiceNumber,
		TransactionDate:  now,
		Subtotal:         subtotal,
		Discount:         input.Discount,
		Total:            totalAmount,
		AmountPaid:       amountPaid,
		ChangeAmount:     changeAmount,
		PaymentMethod:    input.PaymentMethod,
		PaymentStatus:    "COMPLETED", // Langsung tuntas
		CashierID:        cashierID,
		CashierSessionID: input.CashierSessionID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := saleEntity.Validate(); err != nil {
		return nil, err
	}

	// 4. BUNGKUS PAYLOAD KEDALAM OUTBOX EVENT (Vital untuk sinkronisasi asinkron)
	// Payload ini akan dikonsumsi oleh Inventory Service & Finance Service
	outboxEventPayload := map[string]interface{}{
		"sale_id":            saleID,
		"invoice_number":     invoiceNumber,
		"cashier_id":         cashierID,
		"total_transaction":  subtotal,
		"tax_amount":         0.0,
		"total_amount":       totalAmount,
		"payment_method":     input.PaymentMethod,
		"cashier_session_id": input.CashierSessionID,
		"items":              input.Items,
		"created_at":         now,
	}

	event, err := outbox.CreateEvent("SALE", saleID, "SALE_COMPLETED", outboxEventPayload)
	if err != nil {
		return nil, fmt.Errorf("gagal merakit outbox event data: %w", err)
	}

	// 5. Simpan ke database melalui repositori tunggal terisolasi [cite: 50]
	err = u.posRepo.SaveTransaction(ctx, saleEntity, items, event)
	if err != nil {
		return nil, fmt.Errorf("gagal mengeksekusi transaksi checkout kasir: %w", err)
	}

	return saleEntity, nil
}

func (u *posUseCase) SyncOfflineTransactions(ctx context.Context, cashierID string, payload dto.OfflineSyncPayload) error {
	for _, trx := range payload.Transactions {
		_, err := u.Checkout(ctx, cashierID, trx)

		logEntity := &entity.SyncLog{
			ID:         utils.GenerateUUIDv4(),
			EntityType: "SALE",
			EntityID:   trx.IdempotencyKey,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err != nil {
			errMsg := err.Error()
			logEntity.SyncStatus = "FAILED"
			logEntity.ErrorMessage = &errMsg
		} else {
			logEntity.SyncStatus = "SUCCESS"
		}

		_ = u.posRepo.SaveSyncLog(ctx, logEntity)
	}
	return nil
}

func (u *posUseCase) GetReceipt(ctx context.Context, saleID string, cashierName string) (*dto.ReceiptResponse, error) {
	sale, err := u.posRepo.FindByID(ctx, saleID)
	if err != nil || sale == nil {
		return nil, errors.New("transaksi penjualan tidak ditemukan")
	}
	dbItems, _ := u.posRepo.FindItemsBySaleID(ctx, saleID)

	var itemsPayload []dto.SaleItemPayload
	for _, item := range dbItems {
		itemsPayload = append(itemsPayload, dto.SaleItemPayload{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			UnitCode:    item.UnitCode,
			Qty:         item.Qty,
			UnitPrice:   item.UnitPrice,
			Discount:    item.Discount,
		})
	}

	return &dto.ReceiptResponse{
		InvoiceNumber:   sale.InvoiceNumber,
		TransactionDate: sale.TransactionDate,
		PaymentMethod:   sale.PaymentMethod,
		Subtotal:        sale.Subtotal,
		Discount:        sale.Discount,
		Total:           sale.Total,
		AmountPaid:      sale.AmountPaid,
		ChangeAmount:    sale.ChangeAmount,
		CashierName:     cashierName,
		Items:           itemsPayload,
	}, nil
}

func (u *posUseCase) GetSalesHistory(ctx context.Context) ([]*entity.Sale, error) {
	return u.posRepo.FindAll(ctx)
}
func (u *posUseCase) GetSaleDetail(ctx context.Context, id string) (*entity.Sale, error) {
	return u.posRepo.FindByID(ctx, id)
}

func (u *posUseCase) GetProductSalesHistory(ctx context.Context, productID string) ([]*dto.ProductSalesHistoryResponse, error) {
	return u.posRepo.FindSalesHistoryByProductID(ctx, productID)
}

func (u *posUseCase) GetSaleItems(ctx context.Context, saleID string) ([]*entity.SaleItem, error) {
	return u.posRepo.FindItemsBySaleID(ctx, saleID)
}
func (u *posUseCase) GetInvoiceData(ctx context.Context, invoiceNum string) (*entity.Sale, error) {
	return u.posRepo.FindByInvoiceNumber(ctx, invoiceNum)
}
func (u *posUseCase) GetFailedSyncLogs(ctx context.Context) ([]*entity.SyncLog, error) {
	return u.posRepo.FindFailedSyncLogs(ctx)
}

func (u *posUseCase) RetryFailedSync(ctx context.Context, cashierID string, payload dto.OfflineSyncPayload) error {
	for _, trx := range payload.Transactions {
		// Panggil ulang checkout untuk setiap transaksi
		_, err := u.Checkout(ctx, cashierID, trx)
		if err != nil {
			// Cek apakah errornya karena transaksi sudah ada
			if strings.Contains(err.Error(), "IdempotencyKey ini telah diproses sebelumnya") {
				// Anggap sukses (sebelumnya sudah masuk tapi mungkin terlanjur masuk log)
				u.posRepo.UpdateSyncLogStatus(ctx, trx.IdempotencyKey, "SUCCESS", nil)
				continue
			}

			// Jika error lain, biarkan tetap FAILED namun perbarui pesannya
			errMsg := err.Error()
			u.posRepo.UpdateSyncLogStatus(ctx, trx.IdempotencyKey, "FAILED", &errMsg)
			continue
		}

		// Jika sukses, ubah status log sebelumnya menjadi SUCCESS
		u.posRepo.UpdateSyncLogStatus(ctx, trx.IdempotencyKey, "SUCCESS", nil)
	}

	return nil
}

func (u *posUseCase) GetTopProducts(ctx context.Context, limit int, sessionID *string) ([]*dto.TopProductResponse, error) {
	return u.posRepo.GetTopProducts(ctx, limit, sessionID)
}

func (u *posUseCase) UploadInvoice(ctx context.Context, saleID string, cashierName string, fileBytes []byte, contentType string) (string, error) {
	// Dapatkan detail transaksi
	sale, err := u.posRepo.FindByID(ctx, saleID)
	if err != nil {
		return "", fmt.Errorf("failed to get sale details: %v", err)
	}

	// Bersihkan nama kasir dan nomor invoice dari karakter khusus (spasi, garis miring)
	safeCashierName := strings.ReplaceAll(cashierName, " ", "_")
	safeInvoiceNumber := strings.ReplaceAll(sale.InvoiceNumber, "/", "_")

	// Upload ke MinIO dengan path: {nama_kasir}/{invoice_number}.pdf
	objectName := fmt.Sprintf("%s/%s.pdf", safeCashierName, safeInvoiceNumber)
	url, err := u.minioClient.UploadFile(ctx, u.bucketName, objectName, fileBytes, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload invoice to minio: %v", err)
	}

	// Update kolom invoice_url di database
	if err := u.posRepo.UpdateInvoiceURL(ctx, saleID, url); err != nil {
		return "", fmt.Errorf("failed to update invoice_url in db: %v", err)
	}

	return url, nil
}
