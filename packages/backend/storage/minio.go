package storage

import (
	"context"
	"io"
)

// ObjectStorageClient defines the interface for interacting with Object Storage (MinIO)
type ObjectStorageClient interface {
	UploadFile(ctx context.Context, bucketName, objectName string, file io.Reader, objectSize int64, contentType string) error
	GetFileURL(ctx context.Context, bucketName, objectName string) (string, error)
	DeleteFile(ctx context.Context, bucketName, objectName string) error
}

type minioStorageImpl struct {
	// minioClient *minio.Client
}

func NewMinioStorage() ObjectStorageClient {
	return &minioStorageImpl{}
}

func (m *minioStorageImpl) UploadFile(ctx context.Context, bucketName, objectName string, file io.Reader, objectSize int64, contentType string) error {
	// TODO: Implementation for uploading file to MinIO
	return nil
}

func (m *minioStorageImpl) GetFileURL(ctx context.Context, bucketName, objectName string) (string, error) {
	// TODO: Implementation for generating pre-signed URL from MinIO
	return "", nil
}

func (m *minioStorageImpl) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	// TODO: Implementation for deleting file from MinIO
	return nil
}
