package repository

import (
	"context"
	"errors"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/pos/dto"
	"bisnis-rinzi/services/pos/entity"
	"time"

	"github.com/jackc/pgx/v5"
)

type pgPOSRepository struct {
	db *postgres.DBClient
}

func NewPOSRepository(db *postgres.DBClient) POSRepository {
	// Setup FDW untuk inventory_db
	fdwQuery := `
		CREATE EXTENSION IF NOT EXISTS postgres_fdw;
		
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_foreign_server WHERE srvname = 'inventory_server') THEN
				CREATE SERVER inventory_server FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'inventory_db', port '5432');
			END IF;
		END $$;

		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_user_mappings WHERE srvname = 'inventory_server' AND usename = 'postgres') THEN
				CREATE USER MAPPING FOR postgres SERVER inventory_server OPTIONS (user 'postgres', password 'postgres');
			END IF;
		END $$;

		DROP FOREIGN TABLE IF EXISTS inv_products CASCADE;
		IMPORT FOREIGN SCHEMA public LIMIT TO (products) FROM SERVER inventory_server INTO public;
		ALTER FOREIGN TABLE products RENAME TO inv_products;
	`
	_, _ = db.Pool.Exec(context.Background(), fdwQuery)

	return &pgPOSRepository{db: db}
}

func (r *pgPOSRepository) GetProductCostPrice(ctx context.Context, productID string) (float64, error) {
	var costPrice float64
	query := `SELECT COALESCE(cost_price, 0) FROM inv_products WHERE id = $1`
	err := r.db.Pool.QueryRow(ctx, query, productID).Scan(&costPrice)
	if err != nil {
		return 0, err
	}
	return costPrice, nil
}

func (r *pgPOSRepository) SaveTransaction(ctx context.Context, sale *entity.Sale, items []*entity.SaleItem, event *outbox.Event) error {
	// Membuka database transaksi ACID lokal pos_db
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert ke tabel sales
	qSale := `INSERT INTO sales (id, idempotency_key, invoice_number, transaction_date, subtotal, discount, total, amount_paid, change_amount, payment_method, payment_status, cashier_id, cashier_session_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err = tx.Exec(ctx, qSale, sale.ID, sale.IdempotencyKey, sale.InvoiceNumber, sale.TransactionDate, sale.Subtotal, sale.Discount, sale.Total, sale.AmountPaid, sale.ChangeAmount, sale.PaymentMethod, sale.PaymentStatus, sale.CashierID, sale.CashierSessionID, sale.CreatedAt, sale.UpdatedAt)
	if err != nil {
		return err
	}

	// 2. Insert item baris belanja ke tabel sale_items
	qItem := `INSERT INTO sale_items (id, sale_id, product_id, product_name, unit_code, qty, unit_price, discount, subtotal, cost_price, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	for _, item := range items {
		_, err = tx.Exec(ctx, qItem, item.ID, item.SaleID, item.ProductID, item.ProductName, item.UnitCode, item.Qty, item.UnitPrice, item.Discount, item.Subtotal, item.CostPrice, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return err
		}
	}

	// 3. Suntikkan Outbox Event ke tabel outbox_events dalam transaksi yang sama [cite: 50]
	err = outbox.SaveEventTx(ctx, tx, event)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgPOSRepository) CheckIdempotency(ctx context.Context, key string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sales WHERE idempotency_key = $1)`
	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, key).Scan(&exists)
	return exists, err
}

func (r *pgPOSRepository) FindAll(ctx context.Context) ([]*entity.Sale, error) {
	query := `SELECT id, idempotency_key, invoice_number, transaction_date, subtotal, discount, total, amount_paid, change_amount, payment_method, payment_status, cashier_id, cashier_session_id, invoice_url FROM sales ORDER BY transaction_date DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*entity.Sale
	for rows.Next() {
		var s entity.Sale
		if err := rows.Scan(&s.ID, &s.IdempotencyKey, &s.InvoiceNumber, &s.TransactionDate, &s.Subtotal, &s.Discount, &s.Total, &s.AmountPaid, &s.ChangeAmount, &s.PaymentMethod, &s.PaymentStatus, &s.CashierID, &s.CashierSessionID, &s.InvoiceURL); err != nil {
			return nil, err
		}
		sales = append(sales, &s)
	}
	return sales, nil
}

func (r *pgPOSRepository) FindByID(ctx context.Context, id string) (*entity.Sale, error) {
	query := `SELECT id, idempotency_key, invoice_number, transaction_date, subtotal, discount, total, amount_paid, change_amount, payment_method, payment_status, cashier_id, cashier_session_id, invoice_url FROM sales WHERE id = $1`
	var s entity.Sale
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.IdempotencyKey, &s.InvoiceNumber, &s.TransactionDate, &s.Subtotal, &s.Discount, &s.Total, &s.AmountPaid, &s.ChangeAmount, &s.PaymentMethod, &s.PaymentStatus, &s.CashierID, &s.CashierSessionID, &s.InvoiceURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *pgPOSRepository) FindByInvoiceNumber(ctx context.Context, invoiceNum string) (*entity.Sale, error) {
	query := `SELECT id, idempotency_key, invoice_number, transaction_date, subtotal, discount, total, amount_paid, change_amount, payment_method, payment_status, cashier_id, cashier_session_id, invoice_url FROM sales WHERE invoice_number = $1`
	var s entity.Sale
	err := r.db.Pool.QueryRow(ctx, query, invoiceNum).Scan(&s.ID, &s.IdempotencyKey, &s.InvoiceNumber, &s.TransactionDate, &s.Subtotal, &s.Discount, &s.Total, &s.AmountPaid, &s.ChangeAmount, &s.PaymentMethod, &s.PaymentStatus, &s.CashierID, &s.CashierSessionID, &s.InvoiceURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *pgPOSRepository) FindItemsBySaleID(ctx context.Context, saleID string) ([]*entity.SaleItem, error) {
	query := `SELECT id, sale_id, product_id, product_name, unit_code, qty, unit_price, discount, subtotal, cost_price FROM sale_items WHERE sale_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.SaleItem
	for rows.Next() {
		var i entity.SaleItem
		if err := rows.Scan(&i.ID, &i.SaleID, &i.ProductID, &i.ProductName, &i.UnitCode, &i.Qty, &i.UnitPrice, &i.Discount, &i.Subtotal, &i.CostPrice); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	return items, nil
}

func (r *pgPOSRepository) SaveSyncLog(ctx context.Context, log *entity.SyncLog) error {
	query := `INSERT INTO sync_logs (id, entity_type, entity_id, sync_status, error_message, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Pool.Exec(ctx, query, log.ID, log.EntityType, log.EntityID, log.SyncStatus, log.ErrorMessage, log.CreatedAt, log.UpdatedAt)
	return err
}

func (r *pgPOSRepository) FindFailedSyncLogs(ctx context.Context) ([]*entity.SyncLog, error) {
	query := `SELECT id, entity_type, entity_id, sync_status, error_message, created_at FROM sync_logs WHERE sync_status = 'FAILED' ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.SyncLog
	for rows.Next() {
		var l entity.SyncLog
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.SyncStatus, &l.ErrorMessage, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}

func (r *pgPOSRepository) UpdateInvoiceURL(ctx context.Context, saleID string, url string) error {
	query := `UPDATE sales SET invoice_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, url, saleID)
	return err
}

func (r *pgPOSRepository) UpdateSyncLogStatus(ctx context.Context, entityID string, status string, errMsg *string) error {
	query := `
		UPDATE sync_logs 
		SET sync_status = $1, error_message = $2, updated_at = $3
		WHERE entity_id = $4
	`
	_, err := r.db.Pool.Exec(ctx, query, status, errMsg, time.Now(), entityID)
	return err
}

func (r *pgPOSRepository) GetTopProducts(ctx context.Context, limit int, sessionID *string) ([]*dto.TopProductResponse, error) {
	var query string
	var args []interface{}
	
	if sessionID != nil && *sessionID != "" {
		query = `
			SELECT si.product_id, si.product_name as product_name, SUM(si.qty) as total_qty, SUM(si.subtotal) as total_revenue
			FROM sale_items si
			JOIN sales s ON si.sale_id = s.id
			WHERE s.cashier_session_id = $1 AND s.payment_status != 'CANCELLED'
			GROUP BY si.product_id, si.product_name
			ORDER BY total_qty DESC
			LIMIT $2
		`
		args = append(args, *sessionID, limit)
	} else {
		query = `
			SELECT si.product_id, si.product_name as product_name, SUM(si.qty) as total_qty, SUM(si.subtotal) as total_revenue
			FROM sale_items si
			JOIN sales s ON si.sale_id = s.id
			WHERE si.created_at >= NOW() - INTERVAL '1 month' AND s.payment_status != 'CANCELLED'
			GROUP BY si.product_id, si.product_name
			ORDER BY total_qty DESC
			LIMIT $1
		`
		args = append(args, limit)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tops []*dto.TopProductResponse
	for rows.Next() {
		var t dto.TopProductResponse
		if err := rows.Scan(&t.ProductID, &t.ProductName, &t.TotalQty, &t.TotalRevenue); err != nil {
			return nil, err
		}
		tops = append(tops, &t)
	}
	return tops, nil
}

func (r *pgPOSRepository) FindSalesHistoryByProductID(ctx context.Context, productID string) ([]*dto.ProductSalesHistoryResponse, error) {
	query := `
		SELECT s.transaction_date, s.invoice_number, si.qty, si.unit_price, si.subtotal
		FROM sale_items si
		JOIN sales s ON si.sale_id = s.id
		WHERE si.product_id = $1 AND s.payment_status != 'CANCELLED'
		ORDER BY s.transaction_date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*dto.ProductSalesHistoryResponse
	for rows.Next() {
		var h dto.ProductSalesHistoryResponse
		if err := rows.Scan(&h.TransactionDate, &h.InvoiceNumber, &h.Qty, &h.UnitPrice, &h.Subtotal); err != nil {
			return nil, err
		}
		histories = append(histories, &h)
	}
	return histories, nil
}

