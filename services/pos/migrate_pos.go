//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgres://postgres:postgres@localhost:5432/pos_db?sslmode=disable"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	query := `ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS product_name VARCHAR(255) NOT NULL DEFAULT 'Unknown';`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		fmt.Printf("Failed to execute migration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migration successful: added product_name to sale_items.")
}
