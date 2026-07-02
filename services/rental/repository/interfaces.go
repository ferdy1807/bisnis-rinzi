package repository

import (
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/rental/dto"
	"bisnis-rinzi/services/rental/entity"
	"context"
	"time"
)

type RentalRepository interface {
	// Kategori & Produk
	SaveCategory(ctx context.Context, cat *entity.RentalCategory) error
	FindAllCategories(ctx context.Context) ([]*entity.RentalCategory, error)
	FindCategoryByID(ctx context.Context, id string) (*entity.RentalCategory, error)
	UpdateCategory(ctx context.Context, cat *entity.RentalCategory) error
	DeleteCategory(ctx context.Context, id string) error
	
	SaveProduct(ctx context.Context, prod *entity.RentalProduct) error
	FindAllProducts(ctx context.Context) ([]*entity.RentalProduct, error)
	FindProductByID(ctx context.Context, id string) (*entity.RentalProduct, error)
	UpdateProduct(ctx context.Context, prod *entity.RentalProduct) error
	DeleteProduct(ctx context.Context, id string) error
	FindProductByCode(ctx context.Context, code string) (*entity.RentalProduct, error)

	// Media Gambar Produk (MinIO References)
	UpdateProductPhoto(ctx context.Context, id string, objectName, originalName, mimeType string) error

	// Matriks Ketersediaan Inventaris
	GetTotalStockReservedOnDates(ctx context.Context, productID string, start, end time.Time) (int, error)
	FindStockReservationsByProduct(ctx context.Context, productID string, start, end time.Time) ([]*entity.StockReservation, error)
	FindReservationsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Reservation, error)

	// Transaksional Reservasi Operasional
	SaveReservationTx(ctx context.Context, res *entity.Reservation, items []*entity.ReservationItem, stocks []*entity.StockReservation, snap *entity.CustomerSnapshot, contents []*entity.ReservationContent, event *outbox.Event) error
	FindAllReservations(ctx context.Context) ([]*entity.Reservation, error)
	FindReservationByID(ctx context.Context, id string) (*entity.Reservation, error)
	FindReservationItems(ctx context.Context, resID string) ([]*entity.ReservationItem, error)
	UpdateReservationStatus(ctx context.Context, id string, status string) error
	CancelReservationTx(ctx context.Context, id string, event *outbox.Event) error
	FindActiveReservations(ctx context.Context) ([]*entity.Reservation, error)
	FindUpcomingReservations(ctx context.Context) ([]*entity.Reservation, error)
	FindOverdueReservations(ctx context.Context) ([]*entity.Reservation, error)
	SaveReservationContent(ctx context.Context, content *entity.ReservationContent) error
	CountReservationContents(ctx context.Context, resID string) (int, error)

	// Transaksional Pengembalian Unit (Return)
	SaveReturnTx(ctx context.Context, ret *entity.RentalReturn, items []*entity.RentalReturnItem, event *outbox.Event) error
	FindAllReturns(ctx context.Context) ([]*entity.RentalReturn, error)
	FindReturnByID(ctx context.Context, id string) (*entity.RentalReturn, error)
	FindReturnByReservationID(ctx context.Context, resID string) (*entity.RentalReturn, error)
	FindReturnItemsByReturnID(ctx context.Context, returnID string) ([]*entity.RentalReturnItem, error)
	
	// Media Foto Pengembalian
	SaveReturnPhoto(ctx context.Context, photo *entity.ReturnPhoto) error
	FindReturnPhotosByReturnID(ctx context.Context, returnID string) ([]*entity.ReturnPhoto, error)
	FindReturnPhotoByID(ctx context.Context, id string) (*entity.ReturnPhoto, error)
	DeleteReturnPhotoByID(ctx context.Context, id string) error
	
	// Invoice / Receipt Pengembalian
	UpdateReturnReceiptURL(ctx context.Context, returnID string, receiptURL string) error

	// Audit Kerusakan (Rental Damages)
	FindAllDamagedItems(ctx context.Context) ([]*dto.DamagedItemAudit, error)
	SettleDamagedItem(ctx context.Context, damageID string, paymentAction string, auditNotes string) error
}
