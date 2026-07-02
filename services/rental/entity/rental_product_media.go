package entity

import "time"

type RentalProductMedia struct {
	ID               string    `json:"id"`
	RentalProductID  string    `json:"rental_product_id"`
	BucketName       string    `json:"bucket_name"`
	ObjectName       string    `json:"object_name"`
	OriginalFileName string    `json:"original_file_name"`
	MimeType         string    `json:"mime_type"`
	FileSizeValues   int64     `json:"file_size_bytes"`
	CreatedAt        time.Time `json:"created_at"`
}
