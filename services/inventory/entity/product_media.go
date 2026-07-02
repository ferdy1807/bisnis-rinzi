package entity

import "time"

// ProductMedia merepresentasikan model data dari tabel 'product_media'
type ProductMedia struct {
	ID               string    `json:"id"`
	ProductID        string    `json:"product_id"`
	MediaCategory    string    `json:"media_category"` // e.g., "IMAGE", "THUMBNAIL"
	BucketName       string    `json:"bucket_name"`
	ObjectName       string    `json:"object_name"`
	OriginalFileName string    `json:"original_file_name"`
	MimeType         string    `json:"mime_type"`
	FileSizeValues   int64     `json:"file_size_bytes"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
