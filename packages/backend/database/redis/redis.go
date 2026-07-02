package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient membungkus *redis.Client bawaan untuk mempermudah kustomisasi fungsi internal.
type RedisClient struct {
	Client *redis.Client
}

// Config menampung konfigurasi koneksi Redis yang disesuaikan dengan .env Anda.
type Config struct {
	Host         string        // Contoh: localhost
	Port         int           // Contoh: 6379
	Password     string        // Password dari .env (jika ada)
	DB           int           // Nomor Logical DB Redis (default: 0)
	PoolSize     int           // Jumlah maksimum koneksi di dalam pool
	MinIdleConns int           // Jumlah minimum koneksi standby
	DialTimeout  time.Duration // Batas waktu tunggu saat membuka koneksi baru
}

// Connect menginisialisasi connection pool ke Redis server secara aman.
func Connect(ctx context.Context, cfg Config) (*RedisClient, error) {
	// Pasang konfigurasi default jika tidak didefinisikan
	if cfg.Port == 0 {
		cfg.Port = 6379
	}
	if cfg.PoolSize == 0 {
		cfg.PoolSize = 10
	}
	if cfg.MinIdleConns == 0 {
		cfg.MinIdleConns = 3
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5 * time.Second
	}

	// Inisialisasi opsi koneksi go-redis
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
	}

	client := redis.NewClient(options)

	// Lakukan Ping untuk memastikan koneksi ke Redis Container benar-benar terjalin
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("gagal terhubung ke Redis: %w", err)
	}

	log.Printf("Integritas Koneksi Sukses: Berhasil terhubung ke Redis Server (DB %d)", cfg.DB)

	return &RedisClient{Client: client}, nil
}

// Close menutup seluruh connection pool Redis secara aman saat aplikasi dimatikan.
func (r *RedisClient) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}
