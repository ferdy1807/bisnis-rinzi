package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
	"time"

	"bisnis-rinzi/packages/backend/database/minio"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/packages/backend/utils"
	"bisnis-rinzi/services/rental/dto"
	"bisnis-rinzi/services/rental/entity"
	"bisnis-rinzi/services/rental/repository"

	minioSDK "github.com/minio/minio-go/v7"
)

type rentalUseCase struct {
	rentalRepo  repository.RentalRepository
	minioClient *minio.MinioClient
	mediaBucket string
}

func NewRentalUseCase(repo repository.RentalRepository, mc *minio.MinioClient, bucket string) RentalUseCase {
	return &rentalUseCase{
		rentalRepo:  repo,
		minioClient: mc,
		mediaBucket: bucket,
	}
}

func (u *rentalUseCase) CreateCategory(ctx context.Context, input dto.CreateCategoryRequest) error {
	now := time.Now()
	c := &entity.RentalCategory{
		ID:        utils.GenerateUUIDv4(),
		Code:      input.Code,
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := c.Validate(); err != nil {
		return err
	}
	return u.rentalRepo.SaveCategory(ctx, c)
}

func (u *rentalUseCase) CreateProduct(ctx context.Context, input dto.CreateProductRequest) error {
	existing, err := u.rentalRepo.FindProductByCode(ctx, input.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("kode unit rental sudah terdaftar di sistem")
	}

	now := time.Now()
	p := &entity.RentalProduct{
		ID:                utils.GenerateUUIDv4(),
		CategoryID:        input.CategoryID,
		Code:              input.Code,
		Name:              input.Name,
		Description:       input.Description,
		RentalPrice:       input.RentalPrice,
		QuantityAvailable: input.QuantityAvailable,
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := p.Validate(); err != nil {
		return err
	}
	return u.rentalRepo.SaveProduct(ctx, p)
}

func (u *rentalUseCase) CheckStockAvailability(ctx context.Context, input dto.CheckAvailabilityRequest) (bool, error) {
	prod, err := u.rentalRepo.FindProductByID(ctx, input.RentalProductID)
	if err != nil {
		return false, err
	}
	if prod == nil {
		return false, errors.New("produk sewa tidak ditemukan")
	}

	maxStockPhysicalStore := int(prod.QuantityAvailable)
	if maxStockPhysicalStore <= 0 {
		return false, nil
	}

	reservedQty, err := u.rentalRepo.GetTotalStockReservedOnDates(ctx, input.RentalProductID, input.StartDate, input.EndDate)
	if err != nil {
		return false, err
	}

	availableQty := maxStockPhysicalStore - reservedQty
	return availableQty >= input.QtyRequested, nil
}

func (u *rentalUseCase) GetProductCalendar(ctx context.Context, productID string, start, end time.Time) ([]*entity.StockReservation, error) {
	return u.rentalRepo.FindStockReservationsByProduct(ctx, productID, start, end)
}

func (u *rentalUseCase) GetReservationsCalendar(ctx context.Context, start, end time.Time) ([]*entity.Reservation, error) {
	return u.rentalRepo.FindReservationsByDateRange(ctx, start, end)
}

func (u *rentalUseCase) CreateReservation(ctx context.Context, cashierID string, input dto.CreateReservationRequest) (string, string, error) {
	if cashierID == "" {
		cashierID = "00000000-0000-0000-0000-000000000000"
	}
	if input.CustomerIdentity == "" || input.CustomerIdentity == "-" {
		return "", "", fmt.Errorf("nomor identitas KTP wajib diisi")
	}
	if input.StartDate.IsZero() || input.EndDate.IsZero() {
		return "", "", fmt.Errorf("tanggal sewa dan kembali wajib diisi")
	}
	if input.EndDate.Before(input.StartDate) {
		return "", "", fmt.Errorf("tanggal kembali tidak boleh sebelum tanggal sewa")
	}

	for _, item := range input.Items {
		isAvailable, err := u.CheckStockAvailability(ctx, dto.CheckAvailabilityRequest{
			RentalProductID: item.RentalProductID,
			StartDate:       input.StartDate,
			EndDate:         input.EndDate,
			QtyRequested:    int(item.Qty),
		})
		if err != nil {
			return "", "", err
		}
		if !isAvailable {
			return "", "", fmt.Errorf("gagal memesan: unit dengan ID %s tidak mencukupi kuantitasnya pada rentang tanggal tersebut", item.RentalProductID)
		}
	}

	durationDays := math.Ceil(input.EndDate.Sub(input.StartDate).Hours() / 24.0)
	if durationDays <= 0 {
		durationDays = 1
	}

	var subtotal float64
	var reservationItems []*entity.ReservationItem
	var reservationContents []*entity.ReservationContent
	var stockReservations []*entity.StockReservation

	resID := utils.GenerateUUIDv4()
	snapID := utils.GenerateUUIDv4()
	now := time.Now()

	snap := &entity.CustomerSnapshot{
		ID:             snapID,
		CustomerName:   input.CustomerName,
		CustomerPhone:  input.CustomerPhone,
		CustomerIDCard: input.CustomerIdentity,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	for _, item := range input.Items {
		subtotalItem := float64(item.Qty) * item.PricePerPeriod
		subtotal += subtotalItem

		prod, _ := u.rentalRepo.FindProductByID(ctx, item.RentalProductID)

		// Snapshot nama produk agar tidak berubah jika produk diedit di kemudian hari
		productName := "Produk Tidak Dikenal"
		if prod != nil {
			productName = prod.Name
		}

		reservationItems = append(reservationItems, &entity.ReservationItem{
			ID:                  utils.GenerateUUIDv4(),
			RentalReservationID: resID,
			RentalProductID:     item.RentalProductID,
			RentalProductName:   productName, // Snapshot nama produk saat reservasi
			Qty:                 item.Qty,
			PricePerPeriod:      item.PricePerPeriod,
			Subtotal:            subtotalItem,
			CreatedAt:           now,
			UpdatedAt:           now,
		})

		for d := input.StartDate; !d.After(input.EndDate); d = d.AddDate(0, 0, 1) {
			stockReservations = append(stockReservations, &entity.StockReservation{
				ID:              utils.GenerateUUIDv4(),
				RentalProductID: item.RentalProductID,
				ReservationID:   resID,
				ReserveDate:     d,
				QtyReserved:     int(item.Qty),
				CreatedAt:       now,
			})
		}
	}

	for _, c := range input.Contents {
		reservationContents = append(reservationContents, &entity.ReservationContent{
			ID:                  utils.GenerateUUIDv4(),
			RentalReservationID: resID,
			ItemName:            c.ItemName,
			Description:         c.Description,
			Quantity:            c.Quantity,
			ConditionNotes:      c.ConditionNotes,
			CreatedAt:           now,
			UpdatedAt:           now,
		})
	}

	totalAmount := subtotal - input.Discount
	// DownPayment: uang muka dari input (bisa partial), default ke totalAmount jika tidak diisi
	downPayment := input.DownPayment
	if downPayment <= 0 {
		downPayment = totalAmount
	}
	changeAmount := input.AmountPaid - downPayment
	if changeAmount < 0 {
		changeAmount = 0
	}
	invoiceNum := fmt.Sprintf("INV/RENTAL/%d", now.UnixNano())

	resEntity := &entity.Reservation{
		ID:                 resID,
		InvoiceNumber:      invoiceNum,
		CustomerSnapshotID: snapID,
		TransactionDate:    now,
		StartDate:          input.StartDate,
		EndDate:            input.EndDate,
		EventDate:          input.EventDate,
		Subtotal:           subtotal,
		DownPayment:        downPayment,
		AmountPaid:         input.AmountPaid,
		ChangeAmount:       changeAmount,
		TotalAmount:        totalAmount,
		Status:             "BOOKED",
		CashierSessionID:   input.CashierSessionID,
		CreatedBy:          cashierID, // Gunakan ID kasir dari request header, bukan UUID acak
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := resEntity.Validate(); err != nil {
		return "", "", err
	}

	eventPayload := map[string]interface{}{
		"reservation_id": resID,
		"customer_id":    snapID,
		"box_code":       "Rental Items",
		"rental_price":   totalAmount,
		"deposit_amount": downPayment,
		"status_booking": "CONFIRMED",
		"deadline_date":  input.EndDate,
	}
	event, err := outbox.CreateEvent("RENTAL_RESERVATION", resID, "RENTAL_CONFIRMED", eventPayload)
	if err != nil {
		return "", "", err
	}

	err = u.rentalRepo.SaveReservationTx(ctx, resEntity, reservationItems, stockReservations, snap, reservationContents, event)
	if err != nil {
		return "", "", err
	}

	return resID, invoiceNum, nil
}

func (u *rentalUseCase) UploadReservationInvoice(ctx context.Context, resID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error {
	// Verifikasi entitas
	res, err := u.rentalRepo.FindReservationByID(ctx, resID)
	if err != nil || res == nil {
		return errors.New("data reservasi tidak ditemukan")
	}

	targetBucket := "invoice-sewa"
	err = u.minioClient.CreateBucketIfNotExist(ctx, targetBucket, "us-east-1")
	if err == nil {
		_ = u.minioClient.MakeBucketPublic(ctx, targetBucket)
	}

	objectName := fileName
	_, err = u.minioClient.Client.PutObject(ctx, targetBucket, objectName, fileReader, fileSize, minioSDK.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return fmt.Errorf("gagal upload invoice reservasi ke storage: %w", err)
	}

	return nil
}

func (u *rentalUseCase) ProcessRentalReturn(ctx context.Context, cashierID string, input dto.ProcessReturnRequest) (*entity.RentalReturn, error) {
	if cashierID == "" {
		cashierID = "00000000-0000-0000-0000-000000000000"
	}
	res, err := u.rentalRepo.FindReservationByID(ctx, input.ReservationID)
	if err != nil || res == nil {
		return nil, errors.New("data dokumen reservasi sewa tidak ditemukan")
	}
	if res.Status == "RETURNED" {
		return nil, errors.New("transaksi rental ini sudah berstatus selesai dikembalikan")
	}
	if res.Status == "CANCELLED" {
		return nil, errors.New("reservasi yang dibatalkan tidak dapat diproses pengembaliannya")
	}

	now := time.Now()
	returnID := utils.GenerateUUIDv4()

	// ── 1. Hitung keterlambatan dalam hari ──────────────────────────────────
	// Keterlambatan dihitung dari end_date reservasi hingga tanggal pengembalian aktual.
	lateDays := 0
	if now.After(res.EndDate) {
		// math.Ceil: terlambat 1 jam tetap dihitung 1 hari penuh
		lateDays = int(math.Ceil(now.Sub(res.EndDate).Hours() / 24.0))
	}

	// ── 2. Hitung denda keterlambatan (level reservasi, bukan per item) ──────
	// Rumus: 10% per hari × total tagihan awal × jumlah hari terlambat
	// Contoh: 10% × 120.000 × 1 hari = 12.000
	totalLateFees := float64(lateDays) * (res.TotalAmount * 0.10)

	// ── 3. Hitung denda kerusakan (level item, dari input kasir) ─────────────
	// Setiap item memiliki denda kerusakan sendiri yang diinput manual oleh kasir.
	var totalDamageFees float64 = input.ManualDamageFee
	var returnItems []*entity.RentalReturnItem

	// Ambil data item reservasi untuk snapshot nama produk
	dbItems, _ := u.rentalRepo.FindReservationItems(ctx, input.ReservationID)
	productNameMap := make(map[string]string)
	for _, dbItem := range dbItems {
		productNameMap[dbItem.RentalProductID] = dbItem.RentalProductName
	}

	if len(input.ReturnItems) == 0 {
		// Jika kasir frontend tidak mengirimkan item per item, asumsikan semua kembali
		// dan catat denda manual ke keseluruhan pengembalian.
		for _, dbItem := range dbItems {
			returnItems = append(returnItems, &entity.RentalReturnItem{
				ID:                utils.GenerateUUIDv4(),
				RentalReturnID:    returnID,
				RentalProductID:   dbItem.RentalProductID,
				RentalProductName: dbItem.RentalProductName,
				QtyReturned:       dbItem.Qty,
				ConditionStatus:   "GOOD",
				DamageFee:         0, // Agregat denda dicatat di ManualDamageFee
				ConditionNotes:    input.ManualReturnNotes,
				CreatedAt:         now,
				UpdatedAt:         now,
			})
		}
	} else {
		for _, itemInput := range input.ReturnItems {
			totalDamageFees += itemInput.DamageFee

			condStatus := itemInput.ConditionStatus
			if condStatus == "" {
				if itemInput.DamageFee > 0 {
					condStatus = "DAMAGED"
				} else {
					condStatus = "GOOD"
				}
			}

			productName := productNameMap[itemInput.RentalProductID]
			if productName == "" {
				productName = "Produk Tidak Dikenal"
			}

			returnItems = append(returnItems, &entity.RentalReturnItem{
				ID:                utils.GenerateUUIDv4(),
				RentalReturnID:    returnID,
				RentalProductID:   itemInput.RentalProductID,
				RentalProductName: productName,
				QtyReturned:       itemInput.QtyReturned,
				ConditionStatus:   condStatus,
				DamageFee:         itemInput.DamageFee,
				ConditionNotes:    itemInput.ConditionNotes,
				CreatedAt:         now,
				UpdatedAt:         now,
			})
		}
	}

	// ── 4. Hitung sisa tagihan dan total yang harus dibayar saat pengembalian ──
	// Sisa tagihan = total_amount - down_payment (uang muka yang sudah dibayar di awal)
	// Contoh: 120.000 - 100.000 = 20.000
	remainingPayment := res.TotalAmount - res.DownPayment
	if remainingPayment < 0 {
		remainingPayment = 0
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endDay := time.Date(res.EndDate.Year(), res.EndDate.Month(), res.EndDate.Day(), 0, 0, 0, 0, time.Local)

	if today.After(endDay) {
		lateDays = int(today.Sub(endDay).Hours() / 24)
	}
	totalLateFees = float64(lateDays) * 50000.0

	// TOTAL KEWAJIBAN
	totalDamageFees = input.ManualDamageFee
	sisaPokok := res.TotalAmount - res.DownPayment
	remainingPayment = sisaPokok + totalLateFees + totalDamageFees

	// GRAND TOTAL PAID = DownPayment + AmountPaid ini
	grandTotalPaid := res.DownPayment + input.AmountPaid

	retEntity := &entity.RentalReturn{
		ID:               returnID,
		ReservationID:    input.ReservationID,
		ReturnDate:       now,
		LateDays:         lateDays,
		TotalLateFees:    totalLateFees,
		TotalDamageFees:  totalDamageFees,
		RemainingPayment: remainingPayment,
		AmountPaid:       input.AmountPaid,
		ChangeAmount:     input.ChangeAmount,
		GrandTotalPaid:   grandTotalPaid,
		ReceivedBy:       cashierID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	eventPayload := map[string]interface{}{
		"reservation_id":    input.ReservationID,
		"return_id":         returnID,
		"late_days":         lateDays,
		"total_late_fees":   totalLateFees,
		"total_damage_fees": totalDamageFees,
		"remaining_payment": sisaPokok, // Hanya sisa pokok sewa, denda dipisah
		"grand_total_paid":  grandTotalPaid,
		"timestamp":         now,
	}
	event, err := outbox.CreateEvent("RENTAL_RESERVATION", returnID, "PRODUCT_RETURN_PROCESSED", eventPayload)

	if err != nil {
		return nil, err
	}

	if err := u.rentalRepo.SaveReturnTx(ctx, retEntity, returnItems, event); err != nil {
		return nil, err
	}
	return retEntity, nil
}

func (u *rentalUseCase) UploadProductMedia(ctx context.Context, productID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error {
	prod, err := u.rentalRepo.FindProductByID(ctx, productID)
	if err != nil || prod == nil {
		return errors.New("produk sewa tidak terdaftar")
	}

	err = u.minioClient.CreateBucketIfNotExist(ctx, u.mediaBucket, "us-east-1")
	if err != nil {
		return fmt.Errorf("gagal verifikasi minio cluster: %w", err)
	}

	// 1.5 Hapus media lama di MinIO jika ada (agar saat edit foto lama terhapus)
	if prod.ObjectName != nil && *prod.ObjectName != "" {
		_ = u.minioClient.Client.RemoveObject(ctx, u.mediaBucket, *prod.ObjectName, minioSDK.RemoveObjectOptions{})
	}

	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg"
	}
	safeName := strings.ReplaceAll(prod.Name, " ", "_")
	objectName := fmt.Sprintf("%s_%s%s", safeName, time.Now().Format("20060102_150405"), ext)

	_, err = u.minioClient.Client.PutObject(ctx, u.mediaBucket, objectName, fileReader, fileSize, minioSDK.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return fmt.Errorf("gagal transfer berkas biner ke storage: %w", err)
	}

	return u.rentalRepo.UpdateProductPhoto(ctx, productID, objectName, fileName, mimeType)
}

func (u *rentalUseCase) DeleteProductPhoto(ctx context.Context, id string) error {
	prod, err := u.rentalRepo.FindProductByID(ctx, id)
	if err != nil || prod == nil {
		return errors.New("produk sewa tidak ditemukan")
	}

	if prod.ObjectName != nil && *prod.ObjectName != "" {
		_ = u.minioClient.Client.RemoveObject(ctx, u.mediaBucket, *prod.ObjectName, minioSDK.RemoveObjectOptions{})
	}

	return u.rentalRepo.UpdateProductPhoto(ctx, id, "", "", "")
}

func (u *rentalUseCase) GetAllReservations(ctx context.Context) ([]*entity.Reservation, error) {
	return u.rentalRepo.FindAllReservations(ctx)
}

func (u *rentalUseCase) PickupRentalItems(ctx context.Context, id string) error {
	res, err := u.rentalRepo.FindReservationByID(ctx, id)
	if err != nil || res == nil {
		return errors.New("berkas reservasi sewa tidak ditemukan")
	}
	if res.Status != "READY_FOR_PICKUP" {
		return errors.New("unit sewa hanya dapat diambil jika status pesanan adalah READY_FOR_PICKUP (sudah melalui tahap dekorasi)")
	}

	return u.rentalRepo.UpdateReservationStatus(ctx, id, "PICKED_UP")
}

func (u *rentalUseCase) UndoPickupRentalItems(ctx context.Context, id string) error {
	res, err := u.rentalRepo.FindReservationByID(ctx, id)
	if err != nil || res == nil {
		return errors.New("berkas reservasi sewa tidak ditemukan")
	}
	if res.Status != "PICKED_UP" {
		return errors.New("batal serah terima hanya dapat dilakukan jika status pesanan adalah PICKED_UP")
	}

	return u.rentalRepo.UpdateReservationStatus(ctx, id, "READY_FOR_PICKUP")
}

func (u *rentalUseCase) MarkAsReadyForPickup(ctx context.Context, id string) error {
	res, err := u.rentalRepo.FindReservationByID(ctx, id)
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("reservasi tidak ditemukan")
	}
	if res.Status != "BOOKED" && res.Status != "CONTENTS_RECEIVED" {
		return errors.New("hanya pesanan baru yang dapat dipindah ke tahap siap diambil")
	}

	// Validasi keamanan operasional: Pastikan ada barang titipan (rentals contents)
	count, err := u.rentalRepo.CountReservationContents(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal memvalidasi barang titipan: %v", err)
	}
	if count == 0 {
		return errors.New("barang titipan belum dicatat. Silakan lengkapi data barang titipan pelanggan terlebih dahulu melalui menu Kelola")
	}

	return u.rentalRepo.UpdateReservationStatus(ctx, id, "READY_FOR_PICKUP")
}

func (u *rentalUseCase) UndoReadyForPickup(ctx context.Context, id string) error {
	res, err := u.rentalRepo.FindReservationByID(ctx, id)
	if err != nil || res == nil {
		return errors.New("reservasi tidak ditemukan")
	}
	if res.Status != "READY_FOR_PICKUP" {
		return errors.New("pembatalan tahap siap diambil hanya dapat dilakukan pada status READY_FOR_PICKUP")
	}

	return u.rentalRepo.UpdateReservationStatus(ctx, id, "BOOKED")
}

func (u *rentalUseCase) CancelReservation(ctx context.Context, input dto.CancelReservationRequest) error {
	res, err := u.rentalRepo.FindReservationByID(ctx, input.ReservationID)
	if err != nil || res == nil {
		return errors.New("data reservasi sewa tidak ditemukan")
	}
	if res.Status != "BOOKED" && res.Status != "CONTENTS_RECEIVED" && res.Status != "DECORATING" && res.Status != "READY_FOR_PICKUP" {
		return errors.New("pesanan sewa yang sudah diambil atau selesai tidak dapat dibatalkan")
	}

	cancelPayload := map[string]interface{}{
		"reservation_id":      input.ReservationID,
		"invoice_number":      res.InvoiceNumber,
		"refund_rental_fee":   0.0,
		"refund_down_payment": res.DownPayment,
		"penalty_fee":         input.PenaltyFee,
		"timestamp":           time.Now(),
	}
	event, err := outbox.CreateEvent("RENTAL_RESERVATION", input.ReservationID, "RESERVATION_CANCELLED", cancelPayload)
	if err != nil {
		return err
	}

	return u.rentalRepo.CancelReservationTx(ctx, input.ReservationID, event)
}

func (u *rentalUseCase) SaveDepositItems(ctx context.Context, reservationID string, input dto.ReservationContentPayload) error {
	res, err := u.rentalRepo.FindReservationByID(ctx, reservationID)
	if err != nil || res == nil {
		return errors.New("data reservasi sewa tidak ditemukan")
	}

	content := &entity.ReservationContent{
		ID:                  utils.GenerateUUIDv4(),
		RentalReservationID: reservationID,
		ItemName:            input.ItemName,
		Description:         input.Description,
		Quantity:            input.Quantity,
		ConditionNotes:      input.ConditionNotes,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	err = u.rentalRepo.SaveReservationContent(ctx, content)
	if err != nil {
		return err
	}

	// Otomatisasi: Update status menjadi Menunggu Diambil (READY_FOR_PICKUP)
	if res.Status == "BOOKED" || res.Status == "CONTENTS_RECEIVED" || res.Status == "DECORATING" {
		err = u.rentalRepo.UpdateReservationStatus(ctx, reservationID, "READY_FOR_PICKUP")
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *rentalUseCase) GetAllReturns(ctx context.Context) ([]*entity.RentalReturn, error) {
	return u.rentalRepo.FindAllReturns(ctx)
}

func (u *rentalUseCase) GetReturnDetail(ctx context.Context, id string) (*entity.RentalReturn, error) {
	return u.rentalRepo.FindReturnByID(ctx, id)
}

func (u *rentalUseCase) GetReturnItems(ctx context.Context, returnID string) ([]*entity.RentalReturnItem, error) {
	return u.rentalRepo.FindReturnItemsByReturnID(ctx, returnID)
}

func (u *rentalUseCase) GetReturnByReservationID(ctx context.Context, resID string) (*entity.RentalReturn, error) {
	return u.rentalRepo.FindReturnByReservationID(ctx, resID)
}

func (u *rentalUseCase) UpdateCategory(ctx context.Context, id string, input dto.CreateCategoryRequest) error {
	cat, err := u.rentalRepo.FindCategoryByID(ctx, id)
	if err != nil || cat == nil {
		return errors.New("kategori tidak ditemukan")
	}
	cat.Code = input.Code
	cat.Name = input.Name
	cat.UpdatedAt = time.Now()
	return u.rentalRepo.UpdateCategory(ctx, cat)
}

func (u *rentalUseCase) DeleteCategory(ctx context.Context, id string) error {
	cat, err := u.rentalRepo.FindCategoryByID(ctx, id)
	if err != nil || cat == nil {
		return errors.New("kategori tidak ditemukan")
	}
	return u.rentalRepo.DeleteCategory(ctx, id)
}

func (u *rentalUseCase) UpdateProduct(ctx context.Context, id string, input dto.CreateProductRequest) error {
	prod, err := u.rentalRepo.FindProductByID(ctx, id)
	if err != nil || prod == nil {
		return errors.New("produk sewa tidak ditemukan")
	}
	prod.CategoryID = input.CategoryID
	prod.Code = input.Code
	prod.Name = input.Name
	prod.Description = input.Description
	prod.RentalPrice = input.RentalPrice
	prod.QuantityAvailable = input.QuantityAvailable
	prod.IsActive = input.IsActive
	prod.UpdatedAt = time.Now()
	return u.rentalRepo.UpdateProduct(ctx, prod)
}

func (u *rentalUseCase) DeleteProduct(ctx context.Context, id string) error {
	prod, err := u.rentalRepo.FindProductByID(ctx, id)
	if err != nil || prod == nil {
		return errors.New("produk sewa tidak ditemukan")
	}
	return u.rentalRepo.DeleteProduct(ctx, id)
}

func (u *rentalUseCase) GetActiveReservations(ctx context.Context) ([]*entity.Reservation, error) {
	return u.rentalRepo.FindActiveReservations(ctx)
}

func (u *rentalUseCase) GetUpcomingReservations(ctx context.Context) ([]*entity.Reservation, error) {
	return u.rentalRepo.FindUpcomingReservations(ctx)
}

func (u *rentalUseCase) GetOverdueReservations(ctx context.Context) ([]*entity.Reservation, error) {
	return u.rentalRepo.FindOverdueReservations(ctx)
}

func (u *rentalUseCase) GetCategories(ctx context.Context) ([]*entity.RentalCategory, error) {
	return u.rentalRepo.FindAllCategories(ctx)
}
func (u *rentalUseCase) GetCategoryByID(ctx context.Context, id string) (*entity.RentalCategory, error) {
	return u.rentalRepo.FindCategoryByID(ctx, id)
}
func (u *rentalUseCase) GetProducts(ctx context.Context) ([]*entity.RentalProduct, error) {
	return u.rentalRepo.FindAllProducts(ctx)
}
func (u *rentalUseCase) GetProductByID(ctx context.Context, id string) (*entity.RentalProduct, error) {
	return u.rentalRepo.FindProductByID(ctx, id)
}
func (u *rentalUseCase) GetReservationDetail(ctx context.Context, id string) (*entity.Reservation, error) {
	return u.rentalRepo.FindReservationByID(ctx, id)
}
func (u *rentalUseCase) GetReservationItemsDetail(ctx context.Context, resID string) ([]*entity.ReservationItem, error) {
	return u.rentalRepo.FindReservationItems(ctx, resID)
}

func (u *rentalUseCase) GetReturnPenaltySummary(ctx context.Context, returnID string) (map[string]interface{}, error) {
	ret, err := u.rentalRepo.FindReturnByID(ctx, returnID)
	if err != nil || ret == nil {
		return nil, errors.New("dokumen pengembalian tidak ditemukan")
	}

	items, _ := u.rentalRepo.FindReturnItemsByReturnID(ctx, returnID)

	summary := map[string]interface{}{
		"late_days":         ret.LateDays,
		"total_late_fees":   ret.TotalLateFees,
		"total_damage_fees": ret.TotalDamageFees,
		"remaining_payment": ret.RemainingPayment,
		"grand_total_paid":  ret.GrandTotalPaid,
		"items_detail":      items,
	}
	return summary, nil
}

func (u *rentalUseCase) UploadReturnPhoto(ctx context.Context, returnID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error {
	ret, err := u.rentalRepo.FindReturnByID(ctx, returnID)
	if err != nil || ret == nil {
		return errors.New("dokumen pengembalian tidak ditemukan")
	}

	res, err := u.rentalRepo.FindReservationByID(ctx, ret.ReservationID)
	if err != nil || res == nil {
		return errors.New("dokumen reservasi sewa tidak ditemukan")
	}

	targetBucket := "foto-kerusakan-sewa"
	err = u.minioClient.CreateBucketIfNotExist(ctx, targetBucket, "us-east-1")
	if err != nil {
		return fmt.Errorf("gagal verifikasi minio cluster: %w", err)
	}

	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg"
	}

	// res.CustomerSnapshotID contains "Name (Phone)" based on DB query
	safeName := strings.ReplaceAll(res.CustomerSnapshotID, " ", "_")
	safeName = strings.ReplaceAll(safeName, "(", "")
	safeName = strings.ReplaceAll(safeName, ")", "")
	objectName := fmt.Sprintf("%s_%s%s", safeName, time.Now().Format("20060102_150405"), ext)

	_, err = u.minioClient.Client.PutObject(ctx, "foto-kerusakan-sewa", objectName, fileReader, fileSize, minioSDK.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return fmt.Errorf("gagal upload foto ke storage: %w", err)
	}

	photoEntity := &entity.ReturnPhoto{
		ID:               utils.GenerateUUIDv4(),
		RentalReturnID:   returnID,
		BucketName:       targetBucket,
		ObjectName:       objectName,
		OriginalFileName: fileName,
		MimeType:         mimeType,
		FileSizeValues:   fileSize,
		CreatedAt:        time.Now(),
	}
	return u.rentalRepo.SaveReturnPhoto(ctx, photoEntity)
}

func (u *rentalUseCase) GetReturnPhotos(ctx context.Context, returnID string) ([]*entity.ReturnPhoto, error) {
	return u.rentalRepo.FindReturnPhotosByReturnID(ctx, returnID)
}

func (u *rentalUseCase) DeleteReturnPhoto(ctx context.Context, id string) error {
	p, err := u.rentalRepo.FindReturnPhotoByID(ctx, id)
	if err != nil || p == nil {
		return errors.New("foto tidak ditemukan")
	}

	_ = u.minioClient.Client.RemoveObject(ctx, p.BucketName, p.ObjectName, minioSDK.RemoveObjectOptions{})
	return u.rentalRepo.DeleteReturnPhotoByID(ctx, id)
}

func (u *rentalUseCase) UploadReturnReceipt(ctx context.Context, returnID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error {
	ret, err := u.rentalRepo.FindReturnByID(ctx, returnID)
	if err != nil || ret == nil {
		return errors.New("dokumen pengembalian tidak ditemukan")
	}

	res, err := u.rentalRepo.FindReservationByID(ctx, ret.ReservationID)
	if err != nil || res == nil {
		return errors.New("dokumen reservasi sewa tidak ditemukan")
	}

	targetBucket := "invoice-return"
	err = u.minioClient.CreateBucketIfNotExist(ctx, targetBucket, "us-east-1")
	if err != nil {
		return fmt.Errorf("gagal verifikasi minio cluster: %w", err)
	}

	// Pastikan fileName aman
	safeFileName := filepath.Base(fileName)
	objectName := safeFileName

	_, err = u.minioClient.Client.PutObject(ctx, targetBucket, objectName, fileReader, fileSize, minioSDK.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return fmt.Errorf("gagal upload invoice ke storage: %w", err)
	}

	// Update the receipt URL in the database
	receiptURL := fmt.Sprintf("/%s/%s", targetBucket, objectName)
	return u.rentalRepo.UpdateReturnReceiptURL(ctx, returnID, receiptURL)
}

func (u *rentalUseCase) GetDamagedItems(ctx context.Context) ([]*dto.DamagedItemAudit, error) {
	return u.rentalRepo.FindAllDamagedItems(ctx)
}

func (u *rentalUseCase) SettleDamage(ctx context.Context, damageID string, input dto.SettleDamageRequest) error {
	return u.rentalRepo.SettleDamagedItem(ctx, damageID, input.PaymentAction, input.AuditNotes)
}
