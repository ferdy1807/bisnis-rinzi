package usecase

import (
	"bisnis-rinzi/services/rental/dto"
	"bisnis-rinzi/services/rental/entity"
	"context"
	"io"
	"time"
)

type RentalUseCase interface {
	// Kategori & Produk
	CreateCategory(ctx context.Context, input dto.CreateCategoryRequest) error
	GetCategories(ctx context.Context) ([]*entity.RentalCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*entity.RentalCategory, error)
	UpdateCategory(ctx context.Context, id string, input dto.CreateCategoryRequest) error
	DeleteCategory(ctx context.Context, id string) error
	CreateProduct(ctx context.Context, input dto.CreateProductRequest) error
	GetProducts(ctx context.Context) ([]*entity.RentalProduct, error)
	GetProductByID(ctx context.Context, id string) (*entity.RentalProduct, error)
	UpdateProduct(ctx context.Context, id string, input dto.CreateProductRequest) error
	DeleteProduct(ctx context.Context, id string) error

	// Gambar Media (MinIO Integration)
	UploadProductMedia(ctx context.Context, productID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error
	DeleteProductPhoto(ctx context.Context, id string) error

	// Operasional Sirkulasi Rental
	CheckStockAvailability(ctx context.Context, input dto.CheckAvailabilityRequest) (bool, error)
	GetProductCalendar(ctx context.Context, productID string, start, end time.Time) ([]*entity.StockReservation, error)
	GetReservationsCalendar(ctx context.Context, start, end time.Time) ([]*entity.Reservation, error)

	// Reservasi
	CreateReservation(ctx context.Context, cashierID string, input dto.CreateReservationRequest) (string, string, error)
	GetReservationDetail(ctx context.Context, id string) (*entity.Reservation, error)
	GetReservationItemsDetail(ctx context.Context, resID string) ([]*entity.ReservationItem, error)
	GetAllReservations(ctx context.Context) ([]*entity.Reservation, error)
	PickupRentalItems(ctx context.Context, id string) error
	UndoPickupRentalItems(ctx context.Context, id string) error
	MarkAsReadyForPickup(ctx context.Context, id string) error
	UndoReadyForPickup(ctx context.Context, id string) error
	CancelReservation(ctx context.Context, input dto.CancelReservationRequest) error
	SaveDepositItems(ctx context.Context, reservationID string, input dto.ReservationContentPayload) error
	GetAllReturns(ctx context.Context) ([]*entity.RentalReturn, error)
	GetActiveReservations(ctx context.Context) ([]*entity.Reservation, error)
	GetUpcomingReservations(ctx context.Context) ([]*entity.Reservation, error)
	GetOverdueReservations(ctx context.Context) ([]*entity.Reservation, error)

	// Returns Management
	ProcessRentalReturn(ctx context.Context, cashierID string, req dto.ProcessReturnRequest) (*entity.RentalReturn, error)
	GetReturnDetail(ctx context.Context, id string) (*entity.RentalReturn, error)
	GetReturnByReservationID(ctx context.Context, resID string) (*entity.RentalReturn, error)
	GetReturnItems(ctx context.Context, returnID string) ([]*entity.RentalReturnItem, error)
	GetReturnPenaltySummary(ctx context.Context, returnID string) (map[string]interface{}, error)

	UploadReturnPhoto(ctx context.Context, returnID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error
	GetReturnPhotos(ctx context.Context, returnID string) ([]*entity.ReturnPhoto, error)
	DeleteReturnPhoto(ctx context.Context, id string) error
	UploadReturnReceipt(ctx context.Context, returnID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error
	UploadReservationInvoice(ctx context.Context, resID string, fileReader io.Reader, fileName, mimeType string, fileSize int64) error

	// Audit Kerusakan
	GetDamagedItems(ctx context.Context) ([]*dto.DamagedItemAudit, error)
	SettleDamage(ctx context.Context, damageID string, input dto.SettleDamageRequest) error
}
