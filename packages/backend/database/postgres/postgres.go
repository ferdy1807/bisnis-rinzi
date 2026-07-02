package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBClient membungkus pgxpool.Pool untuk kustomisasi lebih lanjut jika diperlukan di masa depan.
type DBClient struct {
	Pool *pgxpool.Pool
}

// Config menampung parameter konfigurasi untuk koneksi database per service.
type Config struct {
	DSN             string        // Format: postgres://username:password@host:port/dbname?sslmode=disable
	MaxConns        int32         // Jumlah maksimum koneksi dalam pool
	MinConns        int32         // Jumlah minimum koneksi standby dalam pool
	MaxConnIdleTime time.Duration // Durasi maks koneksi idle sebelum ditutup
}

// Connect menginisialisasi connection pool yang aman untuk database yang ditentukan.
func Connect(ctx context.Context, cfg Config) (*DBClient, error) {
	// Pasang konfigurasi default jika tidak ditentukan
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 10
	}
	if cfg.MinConns == 0 {
		cfg.MinConns = 2
	}
	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = 15 * time.Minute
	}

	// Parsing DSN string ke objek konfigurasi pgxpool
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("gagal parsing database DSN: %w", err)
	}

	// Atur performa optimal connection pool
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	// Membuat connection pool ke PostgreSQL (Versi 17-alpine)
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat database connection pool: %w", err)
	}

	// Ping database untuk memastikan koneksi fisik benar-benar valid dan terhubung
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("gagal melakukan ping ke database: %w", err)
	}

	log.Printf("Integritas Koneksi Sukses: Berhasil terhubung ke database")

	return &DBClient{Pool: pool}, nil
}

// Close menutup seluruh jaringan koneksi pool dengan aman (graceful shutdown).
func (c *DBClient) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}
