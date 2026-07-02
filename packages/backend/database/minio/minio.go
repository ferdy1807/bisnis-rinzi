package minio

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient membungkus *minio.Client untuk kebutuhan perpanjangan fungsi internal.
type MinioClient struct {
	Client *minio.Client
}

// Config menampung parameter koneksi ke Object Storage MinIO sesuai .env.
type Config struct {
	Endpoint        string // Contoh: localhost:9000
	AccessKeyID     string // Kredensial MINIO_ROOT_USER
	SecretAccessKey string // Kredensial MINIO_ROOT_PASSWORD
	UseSSL          bool   // Set false untuk lingkungan development lokal
}

// Connect menginisialisasi client koneksi ke server MinIO Object Storage.
func Connect(ctx context.Context, cfg Config) (*MinioClient, error) {
	// Inisialisasi minio client core
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi struktur client MinIO: %w", err)
	}

	// Perbaikan: IsOnline() pada SDK terbaru tidak menerima parameter context.Context
	if !client.IsOnline() {
		return nil, fmt.Errorf("gagal terhubung: server MinIO tidak merespons (offline)")
	}

	log.Printf("Integritas Koneksi Sukses: Berhasil terhubung ke MinIO Object Storage")

	return &MinioClient{Client: client}, nil
}

// CreateBucketIfNotExist adalah fungsi pembantu untuk memastikan bucket media telah siap digunakan.
func (m *MinioClient) CreateBucketIfNotExist(ctx context.Context, bucketName string, location string) error {
	err := m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Jika bucket sudah ada, jangan kembalikan error (abaikan)
		exists, errBucketExists := m.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket '%s' sudah tersedia dan siap digunakan.", bucketName)
			return nil
		}
		return fmt.Errorf("gagal membuat bucket '%s': %w", bucketName, err)
	}
	log.Printf("Bucket '%s' berhasil dibuat.", bucketName)
	return nil
}

// MakeBucketPublic mengatur bucket agar isinya dapat diakses oleh publik secara read-only
func (m *MinioClient) MakeBucketPublic(ctx context.Context, bucketName string) error {
	policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"]}]}`, bucketName)
	return m.Client.SetBucketPolicy(ctx, bucketName, policy)
}

// UploadFile mengunggah file byte slice ke MinIO
func (m *MinioClient) UploadFile(ctx context.Context, bucketName, objectName string, fileBytes []byte, contentType string) (string, error) {
	reader := bytes.NewReader(fileBytes)
	_, err := m.Client.PutObject(ctx, bucketName, objectName, reader, int64(len(fileBytes)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	
	// Format URL statis jika bucket dibuat publik, atau endpoint MinIO
	// Contoh sederhana URL public (dengan endpoint API gateway /public/...)
	return fmt.Sprintf("http://localhost:9000/%s/%s", bucketName, objectName), nil
}
