package repository

import (
	"context"
	"errors"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/rental/dto"
	"bisnis-rinzi/services/rental/entity"

	"github.com/jackc/pgx/v5"
)

type pgRentalRepository struct {
	db *postgres.DBClient
}

func NewRentalRepository(db *postgres.DBClient) RentalRepository {
	return &pgRentalRepository{db: db}
}

func (r *pgRentalRepository) SaveCategory(ctx context.Context, c *entity.RentalCategory) error {
	query := `INSERT INTO rental_categories (id, code, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Pool.Exec(ctx, query, c.ID, c.Code, c.Name, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *pgRentalRepository) FindAllCategories(ctx context.Context) ([]*entity.RentalCategory, error) {
	query := `SELECT id, code, name FROM rental_categories WHERE deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.RentalCategory
	for rows.Next() {
		var c entity.RentalCategory
		if err := rows.Scan(&c.ID, &c.Code, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, nil
}

func (r *pgRentalRepository) FindCategoryByID(ctx context.Context, id string) (*entity.RentalCategory, error) {
	query := `SELECT id, code, name FROM rental_categories WHERE id = $1 AND deleted_at IS NULL`
	var c entity.RentalCategory
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.Code, &c.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *pgRentalRepository) SaveProduct(ctx context.Context, p *entity.RentalProduct) error {
	objName, origName, mime := "", "", ""
	if p.ObjectName != nil {
		objName = *p.ObjectName
	}
	if p.OriginalFileName != nil {
		origName = *p.OriginalFileName
	}
	if p.MimeType != nil {
		mime = *p.MimeType
	}

	query := `INSERT INTO rental_products (id, category_id, code, name, description, rental_price, quantity_available, is_active, object_name, original_file_name, mime_type, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.Pool.Exec(ctx, query, p.ID, p.CategoryID, p.Code, p.Name, p.Description, p.RentalPrice, p.QuantityAvailable, p.IsActive, objName, origName, mime, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *pgRentalRepository) FindAllProducts(ctx context.Context) ([]*entity.RentalProduct, error) {
	query := `SELECT p.id, p.category_id, p.code, p.name, COALESCE(p.description, ''), p.rental_price, 
	                 COALESCE(p.quantity_available, 0) - COALESCE((
	                     SELECT SUM(ri.qty) 
	                     FROM rental_items ri 
	                     JOIN rental_reservations rr ON ri.rental_reservation_id = rr.id 
	                     WHERE ri.rental_product_id = p.id AND rr.status IN ('BOOKED', 'READY_FOR_PICKUP', 'PICKED_UP')
	                 ), 0) AS quantity_available, 
	                 COALESCE(p.is_active, true), p.object_name, p.original_file_name, p.mime_type 
	          FROM rental_products p WHERE p.deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.RentalProduct
	for rows.Next() {
		var p entity.RentalProduct
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Code, &p.Name, &p.Description, &p.RentalPrice, &p.QuantityAvailable, &p.IsActive, &p.ObjectName, &p.OriginalFileName, &p.MimeType); err != nil {
			return nil, err
		}
		list = append(list, &p)
	}
	return list, nil
}

func (r *pgRentalRepository) FindProductByID(ctx context.Context, id string) (*entity.RentalProduct, error) {
	query := `SELECT p.id, p.category_id, p.code, p.name, COALESCE(p.description, ''), p.rental_price, 
	                 COALESCE(p.quantity_available, 0) - COALESCE((
	                     SELECT SUM(ri.qty) 
	                     FROM rental_items ri 
	                     JOIN rental_reservations rr ON ri.rental_reservation_id = rr.id 
	                     WHERE ri.rental_product_id = p.id AND rr.status IN ('BOOKED', 'READY_FOR_PICKUP', 'PICKED_UP')
	                 ), 0) AS quantity_available, 
	                 COALESCE(p.is_active, true), p.object_name, p.original_file_name, p.mime_type 
	          FROM rental_products p WHERE p.id = $1 AND p.deleted_at IS NULL`
	var p entity.RentalProduct
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.CategoryID, &p.Code, &p.Name, &p.Description, &p.RentalPrice, &p.QuantityAvailable, &p.IsActive, &p.ObjectName, &p.OriginalFileName, &p.MimeType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgRentalRepository) FindProductByCode(ctx context.Context, code string) (*entity.RentalProduct, error) {
	query := `SELECT id, category_id, code, name, COALESCE(description, ''), rental_price, COALESCE(quantity_available, 0), COALESCE(is_active, true), object_name, original_file_name, mime_type FROM rental_products WHERE code = $1 AND deleted_at IS NULL`
	var p entity.RentalProduct
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(&p.ID, &p.CategoryID, &p.Code, &p.Name, &p.Description, &p.RentalPrice, &p.QuantityAvailable, &p.IsActive, &p.ObjectName, &p.OriginalFileName, &p.MimeType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgRentalRepository) GetTotalStockReservedOnDates(ctx context.Context, productID string, start, end time.Time) (int, error) {
	query := `SELECT COALESCE(MAX(daily_total), 0) FROM (
	            SELECT reserved_date, SUM(qty_reserved) as daily_total 
	            FROM stock_reservations 
	            WHERE rental_product_id = $1 AND reserved_date BETWEEN $2 AND $3
	            GROUP BY reserved_date
	          ) as daily_sums`
	var maxReserved int
	err := r.db.Pool.QueryRow(ctx, query, productID, start, end).Scan(&maxReserved)
	return maxReserved, err
}

func (r *pgRentalRepository) FindStockReservationsByProduct(ctx context.Context, productID string, start, end time.Time) ([]*entity.StockReservation, error) {
	query := `SELECT id, rental_product_id, rental_reservation_id, reserved_date, qty_reserved, created_at 
	          FROM stock_reservations 
	          WHERE rental_product_id = $1 AND reserved_date BETWEEN $2 AND $3 
	          ORDER BY reserved_date ASC`
	rows, err := r.db.Pool.Query(ctx, query, productID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.StockReservation
	for rows.Next() {
		var stk entity.StockReservation
		if err := rows.Scan(&stk.ID, &stk.RentalProductID, &stk.ReservationID, &stk.ReserveDate, &stk.QtyReserved, &stk.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &stk)
	}
	return list, nil
}

func (r *pgRentalRepository) FindReservationsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Reservation, error) {
	query := `SELECT r.id, r.invoice_number, COALESCE(c.customer_name, 'Tanpa Nama') || ' (' || COALESCE(c.customer_phone, '-') || ')', r.transaction_date, r.start_date, r.end_date, r.event_date, r.subtotal, r.down_payment, r.amount_paid, r.change_amount, r.total_amount, r.status, r.picked_up_by, r.picked_up_at, r.cashier_session_id, r.created_by 
	          FROM rental_reservations r
			  LEFT JOIN customer_snapshots c ON r.customer_snapshot_id = c.id
	          WHERE r.start_date <= $2 AND r.end_date >= $1 
	          AND r.status IN ('BOOKED', 'PICKED_UP', 'RETURNED') 
	          ORDER BY r.start_date ASC`
	return r.fetchReservationsArgs(ctx, query, start, end)
}

func (r *pgRentalRepository) SaveReservationTx(ctx context.Context, res *entity.Reservation, items []*entity.ReservationItem, stocks []*entity.StockReservation, snap *entity.CustomerSnapshot, contents []*entity.ReservationContent, event *outbox.Event) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qSnap := `INSERT INTO customer_snapshots (id, customer_name, customer_phone, customer_id_card, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(ctx, qSnap, snap.ID, snap.CustomerName, snap.CustomerPhone, snap.CustomerIDCard, snap.CreatedAt, snap.UpdatedAt)
	if err != nil {
		return err
	}

	qRes := `INSERT INTO rental_reservations (id, invoice_number, customer_snapshot_id, transaction_date, start_date, end_date, event_date, subtotal, down_payment, amount_paid, change_amount, total_amount, status, picked_up_by, picked_up_at, cashier_session_id, created_by, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	_, err = tx.Exec(ctx, qRes, res.ID, res.InvoiceNumber, res.CustomerSnapshotID, res.TransactionDate, res.StartDate, res.EndDate, res.EventDate, res.Subtotal, res.DownPayment, res.AmountPaid, res.ChangeAmount, res.TotalAmount, res.Status, res.PickedUpBy, res.PickedUpAt, res.CashierSessionID, res.CreatedBy, res.CreatedAt, res.UpdatedAt)
	if err != nil {
		return err
	}

	qItem := `INSERT INTO rental_items (id, rental_reservation_id, rental_product_id, rental_product_name, qty, price_per_period, subtotal, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	for _, item := range items {
		_, err = tx.Exec(ctx, qItem, item.ID, item.RentalReservationID, item.RentalProductID, item.RentalProductName, item.Qty, item.PricePerPeriod, item.Subtotal, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return err
		}
	}

	qStock := `INSERT INTO stock_reservations (id, rental_product_id, rental_reservation_id, reserved_date, qty_reserved, created_at) 
	           VALUES ($1, $2, $3, $4, $5, $6)`
	for _, stk := range stocks {
		_, err = tx.Exec(ctx, qStock, stk.ID, stk.RentalProductID, stk.ReservationID, stk.ReserveDate, stk.QtyReserved, stk.CreatedAt)
		if err != nil {
			return err
		}
	}

	qContent := `INSERT INTO rental_reservation_contents (id, rental_reservation_id, item_name, description, quantity, condition_notes, created_at, updated_at) 
	             VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	for _, c := range contents {
		_, err = tx.Exec(ctx, qContent, c.ID, c.RentalReservationID, c.ItemName, c.Description, c.Quantity, c.ConditionNotes, c.CreatedAt, c.UpdatedAt)
		if err != nil {
			return err
		}
	}

	err = outbox.SaveEventTx(ctx, tx, event)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgRentalRepository) FindReservationByID(ctx context.Context, id string) (*entity.Reservation, error) {
	query := `SELECT r.id, r.invoice_number, COALESCE(c.customer_name, 'Tanpa Nama') || ' (' || COALESCE(c.customer_phone, '-') || ')', r.transaction_date, r.start_date, r.end_date, r.event_date, r.subtotal, r.down_payment, r.amount_paid, r.change_amount, r.total_amount, r.status, r.picked_up_by, r.picked_up_at, r.cashier_session_id, r.created_by 
			  FROM rental_reservations r
			  LEFT JOIN customer_snapshots c ON r.customer_snapshot_id = c.id 
			  WHERE r.id = $1`
	var res entity.Reservation
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&res.ID, &res.InvoiceNumber, &res.CustomerSnapshotID, &res.TransactionDate, &res.StartDate, &res.EndDate, &res.EventDate, &res.Subtotal, &res.DownPayment, &res.AmountPaid, &res.ChangeAmount, &res.TotalAmount, &res.Status, &res.PickedUpBy, &res.PickedUpAt, &res.CashierSessionID, &res.CreatedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *pgRentalRepository) FindReservationItems(ctx context.Context, resID string) ([]*entity.ReservationItem, error) {
	query := `SELECT id, rental_reservation_id, rental_product_id, rental_product_name, qty, price_per_period, subtotal FROM rental_items WHERE rental_reservation_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, resID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.ReservationItem
	for rows.Next() {
		var i entity.ReservationItem
		if err := rows.Scan(&i.ID, &i.RentalReservationID, &i.RentalProductID, &i.RentalProductName, &i.Qty, &i.PricePerPeriod, &i.Subtotal); err != nil {
			return nil, err
		}
		list = append(list, &i)
	}
	return list, nil
}

func (r *pgRentalRepository) UpdateReservationStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE rental_reservations SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, status, time.Now(), id)
	return err
}

func (r *pgRentalRepository) SaveReservationContent(ctx context.Context, content *entity.ReservationContent) error {
	query := `INSERT INTO rental_reservation_contents (id, rental_reservation_id, item_name, description, quantity, condition_notes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Pool.Exec(ctx, query, content.ID, content.RentalReservationID, content.ItemName, content.Description, content.Quantity, content.ConditionNotes, content.CreatedAt, content.UpdatedAt)
	return err
}

func (r *pgRentalRepository) CountReservationContents(ctx context.Context, resID string) (int, error) {
	query := `SELECT COUNT(*) FROM rental_reservation_contents WHERE rental_reservation_id = $1`
	var count int
	err := r.db.Pool.QueryRow(ctx, query, resID).Scan(&count)
	return count, err
}

func (r *pgRentalRepository) SaveReturnTx(ctx context.Context, ret *entity.RentalReturn, items []*entity.RentalReturnItem, event *outbox.Event) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qRet := `INSERT INTO rental_returns 
	         (id, rental_reservation_id, return_date, late_days, total_late_fees, total_damage_fees, remaining_payment, amount_paid, change_amount, grand_total_paid, notes, received_by, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err = tx.Exec(ctx, qRet,
		ret.ID, ret.ReservationID, ret.ReturnDate, ret.LateDays,
		ret.TotalLateFees, ret.TotalDamageFees, ret.RemainingPayment,
		ret.AmountPaid, ret.ChangeAmount, ret.GrandTotalPaid,
		ret.Notes, ret.ReceivedBy, ret.CreatedAt, ret.UpdatedAt,
	)
	if err != nil {
		return err
	}

	qItem := `INSERT INTO rental_return_items 
	          (id, rental_return_id, rental_product_id, rental_product_name, qty_returned, condition_status, damage_fee, condition_notes, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	for _, item := range items {
		_, err = tx.Exec(ctx, qItem,
			item.ID, item.RentalReturnID, item.RentalProductID, item.RentalProductName,
			item.QtyReturned, item.ConditionStatus, item.DamageFee, item.ConditionNotes,
			item.CreatedAt, item.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	qReleaseStock := `DELETE FROM stock_reservations WHERE rental_reservation_id = $1`
	_, err = tx.Exec(ctx, qReleaseStock, ret.ReservationID)
	if err != nil {
		return err
	}

	qUpdateRes := `UPDATE rental_reservations SET status = 'RETURNED', amount_paid = amount_paid + $1, change_amount = change_amount + $2, updated_at = $3 WHERE id = $4`
	_, err = tx.Exec(ctx, qUpdateRes, ret.AmountPaid, ret.ChangeAmount, time.Now(), ret.ReservationID)
	if err != nil {
		return err
	}

	err = outbox.SaveEventTx(ctx, tx, event)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgRentalRepository) UpdateProductPhoto(ctx context.Context, id string, objectName, originalName, mimeType string) error {
	query := `UPDATE rental_products SET object_name = $1, original_file_name = $2, mime_type = $3, updated_at = NOW() WHERE id = $4 AND deleted_at IS NULL`
	// Karena kolom diset NOT NULL di database, kita kirimkan string kosong alih-alih NULL
	_, err := r.db.Pool.Exec(ctx, query, objectName, originalName, mimeType, id)
	return err
}

func (r *pgRentalRepository) FindAllReservations(ctx context.Context) ([]*entity.Reservation, error) {
	query := `SELECT r.id, r.invoice_number, COALESCE(c.customer_name, 'Tanpa Nama') || ' (' || COALESCE(c.customer_phone, '-') || ')', r.transaction_date, r.start_date, r.end_date, r.event_date, r.subtotal, r.down_payment, r.amount_paid, r.change_amount, r.total_amount, r.status, r.picked_up_by, r.picked_up_at, r.cashier_session_id, r.created_by 
	          FROM rental_reservations r
			  LEFT JOIN customer_snapshots c ON r.customer_snapshot_id = c.id 
			  ORDER BY r.created_at DESC`
	return r.fetchReservations(ctx, query)
}

func (r *pgRentalRepository) CancelReservationTx(ctx context.Context, id string, event *outbox.Event) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qUpdate := `UPDATE rental_reservations SET status = 'CANCELLED', updated_at = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, qUpdate, time.Now(), id)
	if err != nil {
		return err
	}

	qRelease := `DELETE FROM stock_reservations WHERE rental_reservation_id = $1`
	_, err = tx.Exec(ctx, qRelease, id)
	if err != nil {
		return err
	}

	err = outbox.SaveEventTx(ctx, tx, event)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgRentalRepository) FindAllReturns(ctx context.Context) ([]*entity.RentalReturn, error) {
	query := `SELECT id, rental_reservation_id, return_date, late_days, total_late_fees, total_damage_fees, remaining_payment, grand_total_paid, received_by, COALESCE(receipt_url, '') 
	          FROM rental_returns ORDER BY return_date DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.RentalReturn
	for rows.Next() {
		var rt entity.RentalReturn
		err := rows.Scan(&rt.ID, &rt.ReservationID, &rt.ReturnDate, &rt.LateDays, &rt.TotalLateFees, &rt.TotalDamageFees, &rt.RemainingPayment, &rt.GrandTotalPaid, &rt.ReceivedBy, &rt.ReceiptURL)
		if err != nil {
			return nil, err
		}
		list = append(list, &rt)
	}
	return list, nil
}

func (r *pgRentalRepository) FindReturnByID(ctx context.Context, id string) (*entity.RentalReturn, error) {
	query := `SELECT id, rental_reservation_id, return_date, late_days, total_late_fees, total_damage_fees, remaining_payment, grand_total_paid, received_by, COALESCE(receipt_url, '') 
	          FROM rental_returns WHERE id = $1`
	var rt entity.RentalReturn
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&rt.ID, &rt.ReservationID, &rt.ReturnDate, &rt.LateDays, &rt.TotalLateFees, &rt.TotalDamageFees, &rt.RemainingPayment, &rt.GrandTotalPaid, &rt.ReceivedBy, &rt.ReceiptURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *pgRentalRepository) FindReturnByReservationID(ctx context.Context, resID string) (*entity.RentalReturn, error) {
	query := `SELECT id, rental_reservation_id, return_date, late_days, total_late_fees, total_damage_fees, remaining_payment, grand_total_paid, received_by, COALESCE(receipt_url, '') 
	          FROM rental_returns WHERE rental_reservation_id = $1`
	var rt entity.RentalReturn
	err := r.db.Pool.QueryRow(ctx, query, resID).Scan(&rt.ID, &rt.ReservationID, &rt.ReturnDate, &rt.LateDays, &rt.TotalLateFees, &rt.TotalDamageFees, &rt.RemainingPayment, &rt.GrandTotalPaid, &rt.ReceivedBy, &rt.ReceiptURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *pgRentalRepository) FindReturnItemsByReturnID(ctx context.Context, returnID string) ([]*entity.RentalReturnItem, error) {
	query := `SELECT id, rental_return_id, rental_product_id, rental_product_name, qty_returned, condition_status, damage_fee, condition_notes 
	          FROM rental_return_items WHERE rental_return_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, returnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.RentalReturnItem
	for rows.Next() {
		var i entity.RentalReturnItem
		err := rows.Scan(&i.ID, &i.RentalReturnID, &i.RentalProductID, &i.RentalProductName, &i.QtyReturned, &i.ConditionStatus, &i.DamageFee, &i.ConditionNotes)
		if err != nil {
			return nil, err
		}
		list = append(list, &i)
	}
	return list, nil
}

func (r *pgRentalRepository) UpdateCategory(ctx context.Context, cat *entity.RentalCategory) error {
	query := `UPDATE rental_categories SET code = $1, name = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Pool.Exec(ctx, query, cat.Code, cat.Name, cat.UpdatedAt, cat.ID)
	return err
}

func (r *pgRentalRepository) DeleteCategory(ctx context.Context, id string) error {
	query := `UPDATE rental_categories SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *pgRentalRepository) UpdateProduct(ctx context.Context, prod *entity.RentalProduct) error {
	objName, origName, mime := "", "", ""
	if prod.ObjectName != nil {
		objName = *prod.ObjectName
	}
	if prod.OriginalFileName != nil {
		origName = *prod.OriginalFileName
	}
	if prod.MimeType != nil {
		mime = *prod.MimeType
	}

	query := `UPDATE rental_products 
	          SET category_id = $1, code = $2, name = $3, description = $4, rental_price = $5, quantity_available = $6, is_active = $7, object_name = $8, original_file_name = $9, mime_type = $10, updated_at = $11 
	          WHERE id = $12 AND deleted_at IS NULL`
	_, err := r.db.Pool.Exec(ctx, query, prod.CategoryID, prod.Code, prod.Name, prod.Description, prod.RentalPrice, prod.QuantityAvailable, prod.IsActive, objName, origName, mime, prod.UpdatedAt, prod.ID)
	return err
}

func (r *pgRentalRepository) DeleteProduct(ctx context.Context, id string) error {
	query := `UPDATE rental_products SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *pgRentalRepository) FindActiveReservations(ctx context.Context) ([]*entity.Reservation, error) {
	query := `SELECT r.id, r.invoice_number, COALESCE(c.customer_name, 'Tanpa Nama') || ' (' || COALESCE(c.customer_phone, '-') || ')', r.transaction_date, r.start_date, r.end_date, r.event_date, r.subtotal, r.down_payment, r.amount_paid, r.change_amount, r.total_amount, r.status, r.picked_up_by, r.picked_up_at, r.cashier_session_id, r.created_by 
	          FROM rental_reservations r
			  LEFT JOIN customer_snapshots c ON r.customer_snapshot_id = c.id 
	          WHERE r.status IN ('BOOKED', 'READY_FOR_PICKUP', 'PICKED_UP') 
	          ORDER BY r.start_date ASC`
	return r.fetchReservations(ctx, query)
}

func (r *pgRentalRepository) FindUpcomingReservations(ctx context.Context) ([]*entity.Reservation, error) {
	query := `SELECT r.id, r.invoice_number, COALESCE(c.customer_name, 'Tanpa Nama') || ' (' || COALESCE(c.customer_phone, '-') || ')', r.transaction_date, r.start_date, r.end_date, r.event_date, r.subtotal, r.down_payment, r.amount_paid, r.change_amount, r.total_amount, r.status, r.picked_up_by, r.picked_up_at, r.cashier_session_id, r.created_by 
	          FROM rental_reservations r
			  LEFT JOIN customer_snapshots c ON r.customer_snapshot_id = c.id 
	          WHERE r.status IN ('BOOKED', 'READY_FOR_PICKUP') AND DATE(r.start_date) >= CURRENT_DATE
	          ORDER BY r.start_date ASC`
	return r.fetchReservations(ctx, query)
}

func (r *pgRentalRepository) FindOverdueReservations(ctx context.Context) ([]*entity.Reservation, error) {
	query := `SELECT r.id, r.invoice_number, COALESCE(c.customer_name, 'Tanpa Nama') || ' (' || COALESCE(c.customer_phone, '-') || ')', r.transaction_date, r.start_date, r.end_date, r.event_date, r.subtotal, r.down_payment, r.amount_paid, r.change_amount, r.total_amount, r.status, r.picked_up_by, r.picked_up_at, r.cashier_session_id, r.created_by 
	          FROM rental_reservations r
			  LEFT JOIN customer_snapshots c ON r.customer_snapshot_id = c.id 
	          WHERE r.status = 'PICKED_UP' AND r.end_date < $1 
	          ORDER BY r.end_date ASC`
	return r.fetchReservationsArgs(ctx, query, time.Now())
}

func (r *pgRentalRepository) fetchReservations(ctx context.Context, query string) ([]*entity.Reservation, error) {
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, err := scanReservationsFull(rows)
	if err != nil {
		return nil, err
	}
	if err := r.attachContentsToReservations(ctx, list); err != nil {
		return nil, err
	}
	if err := r.attachItemsToReservations(ctx, list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *pgRentalRepository) fetchReservationsArgs(ctx context.Context, query string, args ...interface{}) ([]*entity.Reservation, error) {
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, err := scanReservationsFull(rows)
	if err != nil {
		return nil, err
	}
	if err := r.attachContentsToReservations(ctx, list); err != nil {
		return nil, err
	}
	if err := r.attachItemsToReservations(ctx, list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *pgRentalRepository) attachContentsToReservations(ctx context.Context, list []*entity.Reservation) error {
	if len(list) == 0 {
		return nil
	}

	ids := make([]string, len(list))
	resMap := make(map[string]*entity.Reservation)
	for i, res := range list {
		ids[i] = res.ID
		resMap[res.ID] = res
		res.Contents = make([]*entity.ReservationContent, 0)
	}

	query := `SELECT id, rental_reservation_id, item_name, COALESCE(description, ''), quantity, COALESCE(condition_notes, ''), created_at, updated_at 
	          FROM rental_reservation_contents 
	          WHERE rental_reservation_id = ANY($1) 
	          ORDER BY created_at ASC`
	rows, err := r.db.Pool.Query(ctx, query, ids)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var c entity.ReservationContent
		if err := rows.Scan(&c.ID, &c.RentalReservationID, &c.ItemName, &c.Description, &c.Quantity, &c.ConditionNotes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return err
		}
		if res, ok := resMap[c.RentalReservationID]; ok {
			res.Contents = append(res.Contents, &c)
		}
	}
	return nil
}

func (r *pgRentalRepository) attachItemsToReservations(ctx context.Context, list []*entity.Reservation) error {
	if len(list) == 0 {
		return nil
	}

	ids := make([]string, len(list))
	resMap := make(map[string]*entity.Reservation)
	for i, res := range list {
		ids[i] = res.ID
		resMap[res.ID] = res
		res.Items = make([]*entity.ReservationItem, 0)
	}

	query := `SELECT id, rental_reservation_id, rental_product_id, rental_product_name, qty, price_per_period, subtotal 
	          FROM rental_items 
	          WHERE rental_reservation_id = ANY($1)`
	rows, err := r.db.Pool.Query(ctx, query, ids)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.ReservationItem
		if err := rows.Scan(&item.ID, &item.RentalReservationID, &item.RentalProductID, &item.RentalProductName, &item.Qty, &item.PricePerPeriod, &item.Subtotal); err != nil {
			return err
		}
		if res, ok := resMap[item.RentalReservationID]; ok {
			res.Items = append(res.Items, &item)
		}
	}
	return nil
}

func scanReservationsFull(rows pgx.Rows) ([]*entity.Reservation, error) {
	var list []*entity.Reservation
	for rows.Next() {
		var res entity.Reservation
		err := rows.Scan(&res.ID, &res.InvoiceNumber, &res.CustomerSnapshotID, &res.TransactionDate, &res.StartDate, &res.EndDate, &res.EventDate, &res.Subtotal, &res.DownPayment, &res.AmountPaid, &res.ChangeAmount, &res.TotalAmount, &res.Status, &res.PickedUpBy, &res.PickedUpAt, &res.CashierSessionID, &res.CreatedBy)
		if err != nil {
			return nil, err
		}
		list = append(list, &res)
	}
	return list, nil
}

func (r *pgRentalRepository) SaveReturnPhoto(ctx context.Context, m *entity.ReturnPhoto) error {
	query := `INSERT INTO rental_return_photos (id, rental_return_id, bucket_name, object_name, original_file_name, mime_type, file_size_bytes, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Pool.Exec(ctx, query, m.ID, m.RentalReturnID, m.BucketName, m.ObjectName, m.OriginalFileName, m.MimeType, m.FileSizeValues, m.CreatedAt)
	return err
}

func (r *pgRentalRepository) FindReturnPhotosByReturnID(ctx context.Context, returnID string) ([]*entity.ReturnPhoto, error) {
	query := `SELECT id, rental_return_id, bucket_name, object_name, original_file_name, mime_type, file_size_bytes, created_at 
	          FROM rental_return_photos WHERE rental_return_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, returnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.ReturnPhoto
	for rows.Next() {
		var m entity.ReturnPhoto
		err := rows.Scan(&m.ID, &m.RentalReturnID, &m.BucketName, &m.ObjectName, &m.OriginalFileName, &m.MimeType, &m.FileSizeValues, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, &m)
	}
	return list, nil
}

func (r *pgRentalRepository) FindReturnPhotoByID(ctx context.Context, id string) (*entity.ReturnPhoto, error) {
	query := `SELECT id, rental_return_id, bucket_name, object_name, original_file_name, mime_type, file_size_bytes, created_at 
	          FROM rental_return_photos WHERE id = $1`
	var m entity.ReturnPhoto
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.RentalReturnID, &m.BucketName, &m.ObjectName, &m.OriginalFileName, &m.MimeType, &m.FileSizeValues, &m.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *pgRentalRepository) DeleteReturnPhotoByID(ctx context.Context, id string) error {
	query := `DELETE FROM rental_return_photos WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *pgRentalRepository) UpdateReturnReceiptURL(ctx context.Context, returnID string, receiptURL string) error {
	query := `UPDATE rental_returns SET receipt_url = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, returnID, receiptURL)
	return err
}

func (r *pgRentalRepository) FindAllDamagedItems(ctx context.Context) ([]*dto.DamagedItemAudit, error) {
	query := `SELECT 
		i.id, 
		i.rental_product_name as item_name, 
		COALESCE(c.customer_name, 'Tanpa Nama') as customer_name,
		i.condition_status, 
		i.condition_notes, 
		i.damage_fee,
		CASE WHEN i.condition_status = 'SETTLED' THEN 'SETTLED' ELSE 'PENDING_AUDIT' END as status
	FROM rental_return_items i
	JOIN rental_returns ret ON i.rental_return_id = ret.id
	JOIN rental_reservations res ON ret.rental_reservation_id = res.id
	LEFT JOIN customer_snapshots c ON res.customer_snapshot_id = c.id
	WHERE i.condition_status != 'GOOD' OR i.damage_fee > 0
	ORDER BY ret.return_date DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*dto.DamagedItemAudit
	for rows.Next() {
		var d dto.DamagedItemAudit
		if err := rows.Scan(&d.ID, &d.ItemName, &d.CustomerName, &d.ConditionStatus, &d.ConditionNotes, &d.DamageFee, &d.Status); err != nil {
			return nil, err
		}
		list = append(list, &d)
	}
	return list, nil
}

func (r *pgRentalRepository) SettleDamagedItem(ctx context.Context, damageID string, paymentAction string, auditNotes string) error {
	query := `UPDATE rental_return_items 
	          SET condition_status = 'SETTLED', condition_notes = condition_notes || ' | Audit: ' || $1 || ' (' || $2 || ')', updated_at = NOW() 
	          WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, auditNotes, paymentAction, damageID)
	return err
}
