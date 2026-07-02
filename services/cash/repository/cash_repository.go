package repository

import (
	"context"
	"errors"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/cash/entity"

	"github.com/jackc/pgx/v5"
)

type pgCashRepository struct {
	db *postgres.DBClient
}

func NewCashRepository(db *postgres.DBClient) CashRepository {
	return &pgCashRepository{db: db}
}

func (r *pgCashRepository) SaveSession(ctx context.Context, s *entity.CashierSession) error {
	query := `INSERT INTO cashier_sessions (id, cashier_id, opening_cash, expected_cash, actual_cash, difference, status, receipt_url, open_time, close_time, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Pool.Exec(ctx, query, s.ID, s.CashierID, s.OpeningCash, s.ExpectedCash, s.ActualCash, s.Difference, s.Status, s.ReceiptURL, s.OpenTime, s.CloseTime, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *pgCashRepository) UpdateSession(ctx context.Context, s *entity.CashierSession) error {
	query := `UPDATE cashier_sessions SET expected_cash=$1, actual_cash=$2, difference=$3, status=$4, close_time=$5, updated_at=$6, receipt_url=$7 WHERE id=$8`
	_, err := r.db.Pool.Exec(ctx, query, s.ExpectedCash, s.ActualCash, s.Difference, s.Status, s.CloseTime, s.UpdatedAt, s.ReceiptURL, s.ID)
	return err
}

func (r *pgCashRepository) FindActiveSessionByCashierID(ctx context.Context, cashierID string) (*entity.CashierSession, error) {
	query := `SELECT id, cashier_id, opening_cash, expected_cash, actual_cash, difference, status, receipt_url, open_time, close_time, created_at, updated_at 
	          FROM cashier_sessions WHERE cashier_id = $1 AND status = 'OPEN' LIMIT 1`
	var s entity.CashierSession
	err := r.db.Pool.QueryRow(ctx, query, cashierID).Scan(&s.ID, &s.CashierID, &s.OpeningCash, &s.ExpectedCash, &s.ActualCash, &s.Difference, &s.Status, &s.ReceiptURL, &s.OpenTime, &s.CloseTime, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}
func (r *pgCashRepository) FindLatestOpenSession(ctx context.Context) (*entity.CashierSession, error) {
	query := `SELECT id, cashier_id, opening_cash, expected_cash, actual_cash, difference, status, receipt_url, open_time, close_time, created_at, updated_at 
	          FROM cashier_sessions WHERE status = 'OPEN' ORDER BY created_at DESC LIMIT 1`
	var s entity.CashierSession
	err := r.db.Pool.QueryRow(ctx, query).Scan(&s.ID, &s.CashierID, &s.OpeningCash, &s.ExpectedCash, &s.ActualCash, &s.Difference, &s.Status, &s.ReceiptURL, &s.OpenTime, &s.CloseTime, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *pgCashRepository) FindSessionByID(ctx context.Context, id string) (*entity.CashierSession, error) {
	query := `SELECT id, cashier_id, opening_cash, expected_cash, actual_cash, difference, status, receipt_url, open_time, close_time, created_at, updated_at 
	          FROM cashier_sessions WHERE id = $1`
	var s entity.CashierSession
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.CashierID, &s.OpeningCash, &s.ExpectedCash, &s.ActualCash, &s.Difference, &s.Status, &s.ReceiptURL, &s.OpenTime, &s.CloseTime, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *pgCashRepository) FindAllSessions(ctx context.Context) ([]*entity.CashierSession, error) {
	query := `SELECT id, cashier_id, opening_cash, expected_cash, actual_cash, difference, status, receipt_url, open_time, close_time, created_at, updated_at,
	          COALESCE((SELECT SUM(amount) FROM cash_transactions WHERE session_id = cashier_sessions.id AND transaction_type = 'DEPOSIT' AND reference_type = 'MANUAL'), 0) as total_manual_income
	          FROM cashier_sessions ORDER BY open_time DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*entity.CashierSession
	for rows.Next() {
		var s entity.CashierSession
		if err := rows.Scan(&s.ID, &s.CashierID, &s.OpeningCash, &s.ExpectedCash, &s.ActualCash, &s.Difference, &s.Status, &s.ReceiptURL, &s.OpenTime, &s.CloseTime, &s.CreatedAt, &s.UpdatedAt, &s.TotalManualIncome); err == nil {
			sessions = append(sessions, &s)
		}
	}
	return sessions, nil
}

func (r *pgCashRepository) CloseSessionTx(ctx context.Context, s *entity.CashierSession, event *outbox.Event) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qUpdate := `UPDATE cashier_sessions SET expected_cash=$1, actual_cash=$2, difference=$3, status=$4, close_time=$5, updated_at=$6 WHERE id=$7`
	_, err = tx.Exec(ctx, qUpdate, s.ExpectedCash, s.ActualCash, s.Difference, s.Status, s.CloseTime, s.UpdatedAt, s.ID)
	if err != nil {
		return err
	}

	err = outbox.SaveEventTx(ctx, tx, event)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgCashRepository) SaveCashTransaction(ctx context.Context, t *entity.CashTransaction) error {
	query := `INSERT INTO cash_transactions (id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.SessionID, t.TransactionType, t.ReferenceType, t.ReferenceID, t.Amount, t.Notes, t.CreatedBy, t.CreatedAt, t.UpdatedAt)
	return err
}

// SaveCashTransactionAndUpdateSessionTx menyimpan mutasi kas dan memperbarui saldo session dalam 1 transaksi DB atomik
func (r *pgCashRepository) SaveCashTransactionAndUpdateSessionTx(ctx context.Context, t *entity.CashTransaction, s *entity.CashierSession) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qTrx := `INSERT INTO cash_transactions (id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(ctx, qTrx, t.ID, t.SessionID, t.TransactionType, t.ReferenceType, t.ReferenceID, t.Amount, t.Notes, t.CreatedBy, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}

	qUpdate := `UPDATE cashier_sessions SET expected_cash=$1, updated_at=$2 WHERE id=$3`
	_, err = tx.Exec(ctx, qUpdate, s.ExpectedCash, s.UpdatedAt, s.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgCashRepository) FindAllTransactions(ctx context.Context) ([]*entity.CashTransaction, error) {
	query := `SELECT id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at FROM cash_transactions ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.CashTransaction
	for rows.Next() {
		var t entity.CashTransaction
		if err := rows.Scan(&t.ID, &t.SessionID, &t.TransactionType, &t.ReferenceType, &t.ReferenceID, &t.Amount, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, nil
}

func (r *pgCashRepository) FindInternalIncomes(ctx context.Context) ([]*entity.CashTransaction, error) {
	query := `SELECT id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at
	          FROM cash_transactions 
			  WHERE transaction_type = 'DEPOSIT' 
			  ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.CashTransaction
	for rows.Next() {
		var t entity.CashTransaction
		if err := rows.Scan(&t.ID, &t.SessionID, &t.TransactionType, &t.ReferenceType, &t.ReferenceID, &t.Amount, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, nil
}

func (r *pgCashRepository) FindTransactionByID(ctx context.Context, id string) (*entity.CashTransaction, error) {
	query := `SELECT id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at
	          FROM cash_transactions WHERE id = $1`
	var t entity.CashTransaction
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.SessionID, &t.TransactionType, &t.ReferenceType, &t.ReferenceID, &t.Amount, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *pgCashRepository) FindTransactionsBySession(ctx context.Context, sessionID string) ([]*entity.CashTransaction, error) {
	query := `SELECT id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at 
	          FROM cash_transactions WHERE session_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.CashTransaction
	for rows.Next() {
		var t entity.CashTransaction
		if err := rows.Scan(&t.ID, &t.SessionID, &t.TransactionType, &t.ReferenceType, &t.ReferenceID, &t.Amount, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, nil
}

func (r *pgCashRepository) SaveExpenseCategory(ctx context.Context, ec *entity.ExpenseCategory) error {
	query := `INSERT INTO expense_categories (id, code, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Pool.Exec(ctx, query, ec.ID, ec.Code, ec.Name, ec.CreatedAt, ec.UpdatedAt)
	return err
}

func (r *pgCashRepository) UpdateExpenseCategory(ctx context.Context, ec *entity.ExpenseCategory) error {
	query := `UPDATE expense_categories SET code = $1, name = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Pool.Exec(ctx, query, ec.Code, ec.Name, ec.UpdatedAt, ec.ID)
	return err
}

func (r *pgCashRepository) DeleteExpenseCategory(ctx context.Context, id string) error {
	query := `DELETE FROM expense_categories WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *pgCashRepository) FindAllExpenseCategories(ctx context.Context) ([]*entity.ExpenseCategory, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM expense_categories ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.ExpenseCategory
	for rows.Next() {
		var ec entity.ExpenseCategory
		if err := rows.Scan(&ec.ID, &ec.Code, &ec.Name, &ec.CreatedAt, &ec.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &ec)
	}
	return categories, nil
}

func (r *pgCashRepository) FindExpenseCategoryByID(ctx context.Context, id string) (*entity.ExpenseCategory, error) {
	query := `SELECT id, code, name, created_at, updated_at FROM expense_categories WHERE id = $1`
	var ec entity.ExpenseCategory
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&ec.ID, &ec.Code, &ec.Name, &ec.CreatedAt, &ec.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ec, nil
}

func (r *pgCashRepository) SaveExpenseTx(ctx context.Context, exp *entity.Expense, txRecord *entity.CashTransaction) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Save expense
	qExp := `INSERT INTO expenses (id, expense_date, category_id, description, amount, created_by, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(ctx, qExp, exp.ID, exp.ExpenseDate, exp.CategoryID, exp.Description, exp.Amount, exp.CreatedBy, exp.CreatedAt, exp.UpdatedAt)
	if err != nil {
		return err
	}

	// 2. Save cash transaction
	qTx := `INSERT INTO cash_transactions (id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at) 
	        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(ctx, qTx, txRecord.ID, txRecord.SessionID, txRecord.TransactionType, txRecord.ReferenceType, txRecord.ReferenceID, txRecord.Amount, txRecord.Notes, txRecord.CreatedBy, txRecord.CreatedAt, txRecord.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// SaveExpenseAndUpdateSessionTx menyimpan pengeluaran, mutasi kas, dan saldo session dalam 1 transaksi DB atomik
func (r *pgCashRepository) SaveExpenseAndUpdateSessionTx(ctx context.Context, exp *entity.Expense, txRecord *entity.CashTransaction, session *entity.CashierSession) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Save expense
	qExp := `INSERT INTO expenses (id, expense_date, category_id, description, amount, created_by, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(ctx, qExp, exp.ID, exp.ExpenseDate, exp.CategoryID, exp.Description, exp.Amount, exp.CreatedBy, exp.CreatedAt, exp.UpdatedAt)
	if err != nil {
		return err
	}

	// 2. Save cash transaction
	qTx := `INSERT INTO cash_transactions (id, session_id, transaction_type, reference_type, reference_id, amount, notes, created_by, created_at, updated_at) 
	        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(ctx, qTx, txRecord.ID, txRecord.SessionID, txRecord.TransactionType, txRecord.ReferenceType, txRecord.ReferenceID, txRecord.Amount, txRecord.Notes, txRecord.CreatedBy, txRecord.CreatedAt, txRecord.UpdatedAt)
	if err != nil {
		return err
	}

	// 3. Update session expected_cash
	qUpdate := `UPDATE cashier_sessions SET expected_cash=$1, updated_at=$2 WHERE id=$3`
	_, err = tx.Exec(ctx, qUpdate, session.ExpectedCash, session.UpdatedAt, session.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgCashRepository) UpdateExpense(ctx context.Context, exp *entity.Expense) error {
	query := `UPDATE expenses 
	          SET category_id = $1, amount = $2, description = $3, updated_at = $4
	          WHERE id = $5`
	_, err := r.db.Pool.Exec(ctx, query, exp.CategoryID, exp.Amount, exp.Description, exp.UpdatedAt, exp.ID)
	return err
}

func (r *pgCashRepository) DeleteExpense(ctx context.Context, id string) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *pgCashRepository) FindAllExpenses(ctx context.Context) ([]*entity.Expense, error) {
	query := `SELECT id, expense_date, category_id, description, amount, created_by, created_at, updated_at FROM expenses ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.Expense
	for rows.Next() {
		var e entity.Expense
		if err := rows.Scan(&e.ID, &e.ExpenseDate, &e.CategoryID, &e.Description, &e.Amount, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	return list, nil
}

func (r *pgCashRepository) FindExpenseByID(ctx context.Context, id string) (*entity.Expense, error) {
	query := `SELECT id, expense_date, category_id, description, amount, created_by, created_at, updated_at FROM expenses WHERE id = $1`
	var e entity.Expense
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&e.ID, &e.ExpenseDate, &e.CategoryID, &e.Description, &e.Amount, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}
