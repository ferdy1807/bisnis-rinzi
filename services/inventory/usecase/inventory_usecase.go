package usecase

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"bisnis-rinzi/packages/backend/database/minio" // Pustaka resmi MinIO v7
	"bisnis-rinzi/packages/backend/utils"
	"bisnis-rinzi/services/inventory/dto"
	"bisnis-rinzi/services/inventory/entity"
	"bisnis-rinzi/services/inventory/repository"

	minioSDK "github.com/minio/minio-go/v7"
)

type inventoryUseCase struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	brandRepo    repository.BrandRepository
	unitRepo     repository.UnitRepository
	syncRepo     repository.SyncRepository
	minioClient  *minio.MinioClient
	mediaBucket  string
}

func NewInventoryUseCase(
	pr repository.ProductRepository,
	cr repository.CategoryRepository,
	br repository.BrandRepository,
	ur repository.UnitRepository,
	sr repository.SyncRepository,
	mc *minio.MinioClient,
	bucket string,
) InventoryUseCase {
	return &inventoryUseCase{
		productRepo:  pr,
		categoryRepo: cr,
		brandRepo:    br,
		unitRepo:     ur,
		syncRepo:     sr,
		minioClient:  mc,
		mediaBucket:  bucket,
	}
}

// Tambahkan kode UUID murni dengan memanfaatkan string acak standar RFC 4122

func (u *inventoryUseCase) CreateProduct(ctx context.Context, input dto.ProductCreateRequest) (string, error) {
	existing, err := u.productRepo.FindBySKU(ctx, input.SKU)
	if err != nil || existing != nil {
		return "", errors.New("kode SKU produk sudah terdaftar")
	}

	now := time.Now()
	pEntity := &entity.Product{
		ID:           utils.GenerateUUIDv4(), // Perbaikan: Menggunakan UUID v4 murni standar
		SKU:          input.SKU,
		CategoryID:   input.CategoryID,
		BrandID:      input.BrandID,
		Name:         input.Name,
		BaseUnitCode: input.BaseUnitCode,
		CostPrice:    input.CostPrice,
		SellingPrice: input.SellingPrice,
		IsActive:     true,
		Barcode:      input.Barcode,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return pEntity.ID, u.productRepo.Save(ctx, pEntity, input.InitialQty)
}

func (u *inventoryUseCase) UpdateProduct(ctx context.Context, id string, p *entity.Product) error {
	p.ID = id
	p.UpdatedAt = time.Now()
	return u.productRepo.Update(ctx, p)
}

func (u *inventoryUseCase) DeleteProduct(ctx context.Context, id string) error {
	return u.productRepo.Delete(ctx, id)
}

func (u *inventoryUseCase) GetAllProducts(ctx context.Context) ([]*entity.Product, error) {
	return u.productRepo.FindAll(ctx)
}

// --- CATEGORIES LOGIC ---
func (u *inventoryUseCase) CreateCategory(ctx context.Context, code, name string) error {
	now := time.Now()
	return u.categoryRepo.Save(ctx, &entity.Category{ID: utils.GenerateUUIDv4(), Code: code, Name: name, CreatedAt: now, UpdatedAt: now})
}
func (u *inventoryUseCase) UpdateCategory(ctx context.Context, id, code, name string) error {
	return u.categoryRepo.Update(ctx, &entity.Category{ID: id, Code: code, Name: name, UpdatedAt: time.Now()})
}
func (u *inventoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	return u.categoryRepo.Delete(ctx, id)
}
func (u *inventoryUseCase) GetAllCategories(ctx context.Context) ([]*entity.Category, error) {
	return u.categoryRepo.FindAll(ctx)
}
func (u *inventoryUseCase) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	return u.categoryRepo.FindByID(ctx, id)
}

// --- BRANDS LOGIC ---
func (u *inventoryUseCase) CreateBrand(ctx context.Context, code, name string) error {
	now := time.Now()
	return u.brandRepo.Save(ctx, &entity.Brand{ID: utils.GenerateUUIDv4(), Code: code, Name: name, CreatedAt: now, UpdatedAt: now})
}
func (u *inventoryUseCase) UpdateBrand(ctx context.Context, id, code, name string) error {
	return u.brandRepo.Update(ctx, &entity.Brand{ID: id, Code: code, Name: name, UpdatedAt: time.Now()})
}
func (u *inventoryUseCase) DeleteBrand(ctx context.Context, id string) error {
	return u.brandRepo.Delete(ctx, id)
}
func (u *inventoryUseCase) GetAllBrands(ctx context.Context) ([]*entity.Brand, error) {
	return u.brandRepo.FindAll(ctx)
}
func (u *inventoryUseCase) GetBrandByID(ctx context.Context, id string) (*entity.Brand, error) {
	return u.brandRepo.FindByID(ctx, id)
}

// --- UNITS LOGIC ---
func (u *inventoryUseCase) CreateUnit(ctx context.Context, code, name string) error {
	now := time.Now()
	return u.unitRepo.Save(ctx, &entity.Unit{ID: utils.GenerateUUIDv4(), Code: code, Name: name, CreatedAt: now, UpdatedAt: now})
}
func (u *inventoryUseCase) UpdateUnit(ctx context.Context, id, code, name string) error {
	return u.unitRepo.Update(ctx, &entity.Unit{ID: id, Code: code, Name: name, UpdatedAt: time.Now()})
}
func (u *inventoryUseCase) DeleteUnit(ctx context.Context, id string) error {
	return u.unitRepo.Delete(ctx, id)
}
func (u *inventoryUseCase) GetAllUnits(ctx context.Context) ([]*entity.Unit, error) {
	return u.unitRepo.FindAll(ctx)
}
func (u *inventoryUseCase) GetUnitByID(ctx context.Context, id string) (*entity.Unit, error) {
	return u.unitRepo.FindByID(ctx, id)
}

func (u *inventoryUseCase) AdjustStock(ctx context.Context, input dto.StockAdjustRequest) error {
	if input.NewQty < 0 {
		return errors.New("kuantitas fisik hasil penyesuaian tidak boleh minus")
	}
	return u.productRepo.AdjustStock(ctx, input.ProductID, input.NewQty, input.Reference)
}

func (u *inventoryUseCase) AddCostHistory(ctx context.Context, productID string, costPrice float64, reference string) error {
	ch := &entity.CostHistory{
		ID:          utils.GenerateUUIDv4(),
		ProductID:   productID,
		AverageCost: costPrice,
		CreatedAt:   time.Now(),
	}
	return u.productRepo.AddCostHistory(ctx, ch)
}

func (u *inventoryUseCase) GetCostHistories(ctx context.Context, productID string) ([]*entity.CostHistory, error) {
	return u.productRepo.FindCostHistoriesByProductID(ctx, productID)
}

func (u *inventoryUseCase) UploadProductMedia(ctx context.Context, productID, category string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error {
	// 1. Memastikan bucket media penampung gambar di MinIO telah terbuat
	err := u.minioClient.CreateBucketIfNotExist(ctx, u.mediaBucket, "us-east-1")
	if err != nil {
		return fmt.Errorf("gagal memvalidasi kapasitas bucket storage: %w", err)
	}

	// 1.5 Hapus media lama jika ada (agar saat edit foto lama terhapus)
	existingMedia, err := u.productRepo.FindMediaByProductID(ctx, productID)
	if err == nil {
		for _, m := range existingMedia {
			_ = u.DeleteProductMedia(ctx, m.ID)
		}
	}

	// 2. Susun nama unik object di MinIO untuk mencegah tumpang tindih berkas
	objectName := fmt.Sprintf("%s_%d_%s", productID, time.Now().Unix(), fileName)

	// 3. Alirkan stream biner file ke kontainer MinIO menggunakan client aslinya
	_, err = u.minioClient.Client.PutObject(ctx, u.mediaBucket, objectName, fileReader, fileSize, minioSDK.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return fmt.Errorf("gagal mentransfer berkas media ke object storage: %w", err)
	}

	// Insert ke database retail
	now := time.Now()
	media := &entity.ProductMedia{
		ID:               utils.GenerateUUIDv4(),
		ProductID:        productID,
		MediaCategory:    category,
		BucketName:       u.mediaBucket,
		ObjectName:       objectName,
		OriginalFileName: fileName,
		MimeType:         mimeType,
		FileSizeValues:   fileSize,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return u.productRepo.SaveProductMedia(ctx, media)
}

func (u *inventoryUseCase) GetProductMedia(ctx context.Context, productID string) ([]*entity.ProductMedia, error) {
	return u.productRepo.FindMediaByProductID(ctx, productID)
}

func (u *inventoryUseCase) GetMediaItem(ctx context.Context, mediaID string) (*entity.ProductMedia, error) {
	return u.productRepo.FindMediaByID(ctx, mediaID)
}

func (u *inventoryUseCase) GetMediaStream(ctx context.Context, mediaID string) (io.ReadCloser, *entity.ProductMedia, error) {
	media, err := u.productRepo.FindMediaByID(ctx, mediaID)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mencari media: %w", err)
	}
	if media == nil {
		return nil, nil, errors.New("data media tidak ditemukan")
	}

	object, err := u.minioClient.Client.GetObject(ctx, media.BucketName, media.ObjectName, minioSDK.GetObjectOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mengambil object dari MinIO: %w", err)
	}

	return object, media, nil
}

func (u *inventoryUseCase) DeleteProductMedia(ctx context.Context, mediaID string) error {
	// 1. Ambil data media dari db
	media, err := u.productRepo.FindMediaByID(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("gagal mencari media: %w", err)
	}
	if media == nil {
		return errors.New("data media tidak ditemukan")
	}

	// 2. Hapus object di MinIO
	err = u.minioClient.Client.RemoveObject(ctx, media.BucketName, media.ObjectName, minioSDK.RemoveObjectOptions{})
	if err != nil {
		// Log error tapi lanjutkan hapus db agar tidak orphaned
	}

	// 3. Hapus dari DB
	return u.productRepo.DeleteMediaByID(ctx, mediaID)
}

func (u *inventoryUseCase) GetSyncData(ctx context.Context, clientVersion int64) (*dto.SyncChangesResponse, error) {
	latestVer, err := u.syncRepo.GetLatestVersion(ctx)
	if err != nil {
		return nil, err
	}

	dbChanges, err := u.syncRepo.GetChangesFromVersion(ctx, clientVersion)
	if err != nil {
		return nil, err
	}

	var genericChanges []interface{}
	for _, change := range dbChanges {
		genericChanges = append(genericChanges, change)
	}

	return &dto.SyncChangesResponse{
		LatestVersion: latestVer,
		Changes:       genericChanges,
	}, nil
}

func (u *inventoryUseCase) ImportProductsFromCSV(ctx context.Context, reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	// 1. Baca baris pertama (Header Kolom)
	_, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("gagal membaca header file CSV: %w", err)
	}

	var rowNumber int = 1
	for {
		rowNumber++
		record, err := csvReader.Read()
		if err == io.EOF {
			break // Akhir file CSV tercapai
		}
		if err != nil {
			return fmt.Errorf("error membaca baris ke-%d: %w", rowNumber, err)
		}

		// Validasi jumlah kolom minimal harus 6
		if len(record) < 6 {
			return fmt.Errorf("baris ke-%d gagal: jumlah kolom kurang dari 6", rowNumber)
		}

		price, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return fmt.Errorf("baris ke-%d gagal: format harga jual '%s' tidak valid", rowNumber, record[4])
		}

		initialQty, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return fmt.Errorf("baris ke-%d gagal: format kuantitas awal '%s' tidak valid", rowNumber, record[5])
		}

		req := dto.ProductCreateRequest{
			SKU:          record[0],
			Name:         record[1],
			CategoryID:   record[2],
			BaseUnitCode: record[3],
			SellingPrice: price,
			InitialQty:   initialQty,
		}

		// KOREKSI UTAMA: Tangkap error insert database, jangan diabaikan!
		_, err = u.CreateProduct(ctx, req)
		if err != nil {
			return fmt.Errorf("gagal menyimpan data baris ke-%d (SKU: %s): %w", rowNumber, req.SKU, err)
		}
	}
	return nil
}

func (u *inventoryUseCase) ExportProductsToCSV(ctx context.Context, writer io.Writer) error {
	products, err := u.productRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Tulis Header Dokumen CSV
	_ = csvWriter.Write([]string{"SKU", "Nama Produk", "ID Kategori", "Satuan", "Harga Jual"})

	for _, p := range products {
		_ = csvWriter.Write([]string{
			p.SKU,
			p.Name,
			p.CategoryID,
			p.BaseUnitCode,
			fmt.Sprintf("%.2f", p.SellingPrice),
		})
	}
	return nil
}

// Tambahan Baru: Menghubungkan fungsi pencarian teks dari delivery handler ke repository
func (u *inventoryUseCase) SearchProducts(ctx context.Context, query string) ([]*entity.Product, error) {
	if query == "" {
		return u.productRepo.FindAll(ctx)
	}
	return u.productRepo.FindBySearch(ctx, query)
}

// Tambahan Baru: Melayani penarikan nilai sequence versi server untuk endpoint /api/inventory/sync/version
func (u *inventoryUseCase) GetLatestSyncVersion(ctx context.Context) (int64, error) {
	return u.syncRepo.GetLatestVersion(ctx)
}

// Sisa fungsi pencarian bypass mapping langsung ke layer repository...
func (u *inventoryUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return u.productRepo.FindByID(ctx, id)
}
func (u *inventoryUseCase) GetProductBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	return u.productRepo.FindBySKU(ctx, sku)
}
func (u *inventoryUseCase) GetProductByBarcode(ctx context.Context, barcode string) (*entity.Product, error) {
	return u.productRepo.FindByBarcode(ctx, barcode)
}
func (u *inventoryUseCase) GetLowStockProducts(ctx context.Context) ([]*entity.Product, error) {
	return u.productRepo.FindLowStock(ctx)
}
func (u *inventoryUseCase) GetStock(ctx context.Context, productID string) (*entity.ProductStock, error) {
	return u.productRepo.GetStockByID(ctx, productID)
}

func (u *inventoryUseCase) UpdateStockThresholds(ctx context.Context, productID string, input dto.StockThresholdRequest) error {
	if input.MinStock < 0 || input.SafetyStock < 0 {
		return errors.New("min_stock dan safety_stock tidak boleh bernilai negatif")
	}

	_, err := u.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	return u.productRepo.UpdateStockThresholds(ctx, productID, input.MinStock, input.SafetyStock)
}

func (u *inventoryUseCase) GetCatalogSyncData(ctx context.Context) (map[string]interface{}, error) {
	return u.syncRepo.GetFullCatalogSync(ctx)
}

func (u *inventoryUseCase) GetAllStocksData(ctx context.Context) ([]map[string]interface{}, error) {
	return u.productRepo.FindAllStocks(ctx)
}

func (u *inventoryUseCase) GetStockCard(ctx context.Context, productID string) ([]*entity.StockMovement, error) {
	return u.productRepo.FindStockCardByProductID(ctx, productID)
}

func (u *inventoryUseCase) GetAllMovements(ctx context.Context) ([]*entity.StockMovement, error) {
	return u.productRepo.FindAllStockMovements(ctx)
}

func (u *inventoryUseCase) GetMovementByID(ctx context.Context, id string) (*entity.StockMovement, error) {
	return u.productRepo.FindStockMovementByID(ctx, id)
}

func (u *inventoryUseCase) InternalDeductStock(ctx context.Context, input dto.InternalStockRequest) error {
	if input.Qty <= 0 {
		return errors.New("kuantitas pengurangan stok harus lebih dari nol")
	}

	// Cek stok saat ini
	currentStock, err := u.productRepo.GetStockByID(ctx, input.ProductID)
	if err != nil {
		return fmt.Errorf("gagal memeriksa stok: %v", err)
	}
	if currentStock.Qty < input.Qty {
		return errors.New("stok tidak mencukupi untuk dikurangi")
	}

	newQty := currentStock.Qty - input.Qty
	return u.productRepo.AdjustStock(ctx, input.ProductID, newQty, "INTERNAL_DEDUCT: "+input.Reference)
}

func (u *inventoryUseCase) InternalRestoreStock(ctx context.Context, input dto.InternalStockRequest) error {
	if input.Qty <= 0 {
		return errors.New("kuantitas pengembalian stok harus lebih dari nol")
	}

	// Cek stok saat ini
	currentStock, err := u.productRepo.GetStockByID(ctx, input.ProductID)
	if err != nil {
		return fmt.Errorf("gagal memeriksa stok: %v", err)
	}

	newQty := currentStock.Qty + input.Qty
	return u.productRepo.AdjustStock(ctx, input.ProductID, newQty, "INTERNAL_RESTORE: "+input.Reference)
}
