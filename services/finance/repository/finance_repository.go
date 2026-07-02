package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"bisnis-rinzi/packages/backend/database/postgres"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/finance/dto"
	"bisnis-rinzi/services/finance/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type pgFinanceRepository struct {
	db *postgres.DBClient
}

func NewFinanceRepository(db *postgres.DBClient) FinanceRepository {
	// Memastikan tabel processed_inbox_events selalu ada untuk idempotensi outbox worker
	query := `CREATE TABLE IF NOT EXISTS processed_inbox_events (
		event_id VARCHAR(36) PRIMARY KEY,
		processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, _ = db.Pool.Exec(context.Background(), query)

	// Setup Foreign Data Wrapper (FDW) untuk analitik agregasi lintas-layanan
	fdwQuery := `
		CREATE EXTENSION IF NOT EXISTS postgres_fdw;
		
		-- Setup koneksi ke pos_db
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_foreign_server WHERE srvname = 'pos_server') THEN
				CREATE SERVER pos_server FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'pos_db', port '5432');
			END IF;
		END $$;
		
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_user_mappings WHERE srvname = 'pos_server' AND usename = 'postgres') THEN
				CREATE USER MAPPING FOR postgres SERVER pos_server OPTIONS (user 'postgres', password 'postgres');
			END IF;
		END $$;
		
		DROP FOREIGN TABLE IF EXISTS sales, sale_items CASCADE;
		IMPORT FOREIGN SCHEMA public LIMIT TO (sales, sale_items) FROM SERVER pos_server INTO public;

		-- Setup koneksi ke inventory_db
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

		DROP FOREIGN TABLE IF EXISTS products, categories, product_stocks CASCADE;
		IMPORT FOREIGN SCHEMA public LIMIT TO (products, categories, product_stocks) FROM SERVER inventory_server INTO public;

		-- Setup koneksi ke cash_db
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_foreign_server WHERE srvname = 'cash_server') THEN
				CREATE SERVER cash_server FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'cash_db', port '5432');
			END IF;
		END $$;

		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_user_mappings WHERE srvname = 'cash_server' AND usename = 'postgres') THEN
				CREATE USER MAPPING FOR postgres SERVER cash_server OPTIONS (user 'postgres', password 'postgres');
			END IF;
		END $$;

		DROP FOREIGN TABLE IF EXISTS cashier_sessions, cash_transactions CASCADE;
		IMPORT FOREIGN SCHEMA public LIMIT TO (cashier_sessions, cash_transactions) FROM SERVER cash_server INTO public;

		-- Setup koneksi ke rental_db
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_foreign_server WHERE srvname = 'rental_server') THEN
				CREATE SERVER rental_server FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'rental_db', port '5432');
			END IF;
		END $$;

		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_user_mappings WHERE srvname = 'rental_server' AND usename = 'postgres') THEN
				CREATE USER MAPPING FOR postgres SERVER rental_server OPTIONS (user 'postgres', password 'postgres');
			END IF;
		END $$;

		DROP FOREIGN TABLE IF EXISTS rental_reservations, rental_returns, rental_return_items CASCADE;
		IMPORT FOREIGN SCHEMA public LIMIT TO (rental_reservations, rental_returns, rental_return_items) FROM SERVER rental_server INTO public;

		-- Buat view otomatis untuk merekap log kerusakan mika dari rental db
		DROP TABLE IF EXISTS finance_rental_damage_logs CASCADE;
		CREATE OR REPLACE VIEW finance_rental_damage_logs AS
		SELECT 
			ri.id as id,
			r.return_date::date as log_date,
			ri.rental_return_id as rental_return_id,
			ri.rental_product_id as rental_product_id,
			ri.rental_product_name as product_name,
			ri.damage_fee as penalty_amount,
			ri.condition_status as condition_status,
			ri.condition_notes as notes,
			ri.created_at as created_at,
			ri.updated_at as updated_at
		FROM rental_return_items ri
		JOIN rental_returns r ON ri.rental_return_id = r.id
		WHERE ri.condition_status != 'GOOD' AND ri.damage_fee > 0;

		-- Buat view pengganti daily_closings (live_daily_closings)
		CREATE OR REPLACE VIEW live_daily_closings AS
		SELECT 
			cs.id, 
			COALESCE(cs.close_time, CURRENT_TIMESTAMP) as closing_date, 
			(SELECT COALESCE(SUM(si.subtotal), 0) FROM sale_items si JOIN sales s ON si.sale_id = s.id WHERE s.cashier_session_id = cs.id) as total_sales_retail, 
			(
				(SELECT COALESCE(SUM(down_payment), 0) FROM rental_reservations rr WHERE rr.cashier_session_id = cs.id) + 
				(SELECT COALESCE(SUM(remaining_payment), 0) FROM rental_returns ret WHERE ret.received_by = cs.cashier_id AND ret.return_date >= cs.open_time AND (cs.close_time IS NULL OR ret.return_date <= cs.close_time))
			) as total_rental_income, 
			(SELECT COALESCE(SUM(amount), 0) FROM cash_transactions ct WHERE ct.transaction_type = 'DEPOSIT' AND ct.reference_type = 'MANUAL' AND ct.session_id = cs.id) as total_other_income, 
			(SELECT COALESCE(SUM(amount), 0) FROM cash_transactions ct WHERE ct.transaction_type = 'WITHDRAWAL' AND ct.session_id = cs.id) as total_expenses, 
			cs.actual_cash - cs.opening_cash as net_cash_flow, 
			cs.actual_cash as actual_cash,
			cs.opening_cash as opening_cash,
			false as is_reconciled
		FROM cashier_sessions cs
		WHERE cs.status IN ('OPEN', 'CLOSED');
	`
	_, err := db.Pool.Exec(context.Background(), fdwQuery)
	if err != nil {
		log.Printf("CRITICAL ERROR: Gagal setup FDW dan live_daily_closings: %v\n", err)
	} else {
		log.Println("Setup FDW dan live_daily_closings berhasil.")
	}

	return &pgFinanceRepository{db: db}
}

func (r *pgFinanceRepository) SaveCOA(ctx context.Context, c *entity.ChartOfAccount) error {
	query := `INSERT INTO chart_of_accounts (id, account_code, account_name, account_group, normal_balance, current_balance, is_active, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Pool.Exec(ctx, query, c.ID, c.AccountCode, c.AccountName, c.AccountGroup, c.NormalBalance, c.CurrentBalance, c.IsActive, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *pgFinanceRepository) FindAllCOA(ctx context.Context) ([]*entity.ChartOfAccount, error) {
	query := `SELECT id, account_code, account_name, account_group, normal_balance, current_balance, is_active FROM chart_of_accounts ORDER BY account_code ASC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.ChartOfAccount
	for rows.Next() {
		var c entity.ChartOfAccount
		if err := rows.Scan(&c.ID, &c.AccountCode, &c.AccountName, &c.AccountGroup, &c.NormalBalance, &c.CurrentBalance, &c.IsActive); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, nil
}

func (r *pgFinanceRepository) FindCOAByCode(ctx context.Context, code string) (*entity.ChartOfAccount, error) {
	query := `SELECT id, account_code, account_name, account_group, normal_balance, current_balance, is_active FROM chart_of_accounts WHERE account_code = $1`
	var c entity.ChartOfAccount
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(&c.ID, &c.AccountCode, &c.AccountName, &c.AccountGroup, &c.NormalBalance, &c.CurrentBalance, &c.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *pgFinanceRepository) FindCOAByID(ctx context.Context, id string) (*entity.ChartOfAccount, error) {
	query := `SELECT id, account_code, account_name, account_group, normal_balance, current_balance, is_active FROM chart_of_accounts WHERE id = $1`
	var c entity.ChartOfAccount
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.AccountCode, &c.AccountName, &c.AccountGroup, &c.NormalBalance, &c.CurrentBalance, &c.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *pgFinanceRepository) UpdateCOABalance(ctx context.Context, id string, amount float64) error {
	query := `UPDATE chart_of_accounts SET current_balance = current_balance + $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, amount, time.Now(), id)
	return err
}

func (r *pgFinanceRepository) SavePeriod(ctx context.Context, p *entity.AccountingPeriod) error {
	query := `INSERT INTO accounting_periods (id, name, start_date, end_date, is_closed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Pool.Exec(ctx, query, p.ID, p.Name, p.StartDate, p.EndDate, p.IsClosed, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *pgFinanceRepository) FindActivePeriod(ctx context.Context) (*entity.AccountingPeriod, error) {
	query := `SELECT id, name, start_date, end_date, is_closed FROM accounting_periods WHERE is_closed = false ORDER BY start_date ASC LIMIT 1`
	var p entity.AccountingPeriod
	err := r.db.Pool.QueryRow(ctx, query).Scan(&p.ID, &p.Name, &p.StartDate, &p.EndDate, &p.IsClosed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgFinanceRepository) FindPeriodByID(ctx context.Context, id string) (*entity.AccountingPeriod, error) {
	query := `SELECT id, name, start_date, end_date, is_closed FROM accounting_periods WHERE id = $1`
	var p entity.AccountingPeriod
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.StartDate, &p.EndDate, &p.IsClosed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgFinanceRepository) UpdatePeriod(ctx context.Context, p *entity.AccountingPeriod) error {
	query := `UPDATE accounting_periods SET name = $1, start_date = $2, end_date = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.Pool.Exec(ctx, query, p.Name, p.StartDate, p.EndDate, p.UpdatedAt, p.ID)
	return err
}

func (r *pgFinanceRepository) DeletePeriod(ctx context.Context, id string) error {
	query := `DELETE FROM accounting_periods WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *pgFinanceRepository) ClosePeriodTx(ctx context.Context, periodID string, lock *entity.PeriodLock) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE accounting_periods SET is_closed = true, updated_at = $1 WHERE id = $2", time.Now(), periodID)
	if err != nil {
		return err
	}

	qLock := `INSERT INTO period_locks (id, accounting_period_id, locked_by, lock_reason, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, qLock, lock.ID, lock.AccountingPeriodID, lock.LockedBy, lock.LockReason, lock.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgFinanceRepository) SaveJournal(ctx context.Context, j *entity.Journal) error {
	query := `INSERT INTO journals (id, journal_code, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Pool.Exec(ctx, query, j.ID, j.JournalCode, j.Name, j.Description, j.CreatedAt, j.UpdatedAt)
	return err
}

func (r *pgFinanceRepository) FindJournalByCode(ctx context.Context, code string) (*entity.Journal, error) {
	query := `SELECT id, journal_code, name, description FROM journals WHERE journal_code = $1`
	var j entity.Journal
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(&j.ID, &j.JournalCode, &j.Name, &j.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &j, nil
}

func (r *pgFinanceRepository) SaveJournalEntryTx(ctx context.Context, je *entity.JournalEntry, details []*entity.JournalEntryDetail) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Amankan header entri jurnal umum
	qEntry := `INSERT INTO journal_entries (id, journal_id, accounting_period_id, reference_number, entry_date, narration, is_posted, created_at, updated_at) 
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tx.Exec(ctx, qEntry, je.ID, je.JournalID, je.AccountingPeriodID, je.ReferenceNumber, je.EntryDate, je.Narration, je.IsPosted, je.CreatedAt, je.UpdatedAt)
	if err != nil {
		return err
	}

	// 2. Petakan baris multi-rekening akuntansi detail (Debet & Kredit)
	qDetail := `INSERT INTO journal_entry_details (id, journal_entry_id, account_id, debit_amount, credit_amount, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	for _, det := range details {
		_, err = tx.Exec(ctx, qDetail, det.ID, det.JournalEntryID, det.AccountID, det.DebitAmount, det.CreditAmount, det.CreatedAt)
		if err != nil {
			return err
		}

		// 3. Mutasikan nilai neraca akun secara real-time berdasarkan prinsip Double-Entry Bookkeeping
		var coa entity.ChartOfAccount
		err = tx.QueryRow(ctx, "SELECT normal_balance FROM chart_of_accounts WHERE id = $1 FOR UPDATE", det.AccountID).Scan(&coa.NormalBalance)
		if err != nil {
			return err
		}

		var netMutation float64
		if coa.NormalBalance == "DEBIT" {
			netMutation = det.DebitAmount - det.CreditAmount
		} else {
			netMutation = det.CreditAmount - det.DebitAmount
		}

		_, err = tx.Exec(ctx, "UPDATE chart_of_accounts SET current_balance = current_balance + $1, updated_at = $2 WHERE id = $3", netMutation, time.Now(), det.AccountID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pgFinanceRepository) SaveFinancialTransaction(ctx context.Context, f *entity.FinancialTransaction) error {
	query := `INSERT INTO financial_transactions (id, transaction_type, reference_id, amount, description, transaction_date, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Pool.Exec(ctx, query, f.ID, f.TransactionType, f.ReferenceID, f.Amount, f.Description, f.TransactionDate, f.CreatedAt)
	return err
}

func (r *pgFinanceRepository) SumTransactionAmountByType(ctx context.Context, txType string, date time.Time) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0.0) FROM financial_transactions WHERE transaction_type = $1 AND DATE(transaction_date) = DATE($2)`
	var total float64
	return total, r.db.Pool.QueryRow(ctx, query, txType, date).Scan(&total)
}

func (r *pgFinanceRepository) UpdateCOA(ctx context.Context, coa *entity.ChartOfAccount) error {
	query := `
		UPDATE chart_of_accounts
		SET account_name = $1, normal_balance = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Pool.Exec(ctx, query, coa.AccountName, coa.NormalBalance, coa.IsActive, coa.UpdatedAt, coa.ID)
	return err
}

func (r *pgFinanceRepository) DeleteCOA(ctx context.Context, id string) error {
	query := `DELETE FROM chart_of_accounts WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *pgFinanceRepository) FindAllJournalEntries(ctx context.Context) ([]*entity.JournalEntry, error) {
	query := `SELECT id, journal_id, accounting_period_id, reference_number, entry_date, narration, is_posted, created_at, updated_at FROM journal_entries ORDER BY entry_date DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*entity.JournalEntry
	var journalIDs []string
	for rows.Next() {
		var j entity.JournalEntry
		if err := rows.Scan(&j.ID, &j.JournalID, &j.AccountingPeriodID, &j.ReferenceNumber, &j.EntryDate, &j.Narration, &j.IsPosted, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, &j)
		journalIDs = append(journalIDs, j.ID)
	}

	if len(journalIDs) > 0 {
		detQuery := `SELECT id, journal_entry_id, account_id, debit_amount, credit_amount, created_at FROM journal_entry_details WHERE journal_entry_id = ANY($1)`
		detRows, err := r.db.Pool.Query(ctx, detQuery, journalIDs)
		if err == nil {
			defer detRows.Close()
			detMap := make(map[string][]*entity.JournalEntryDetail)
			for detRows.Next() {
				var d entity.JournalEntryDetail
				if err := detRows.Scan(&d.ID, &d.JournalEntryID, &d.AccountID, &d.DebitAmount, &d.CreditAmount, &d.CreatedAt); err == nil {
					detMap[d.JournalEntryID] = append(detMap[d.JournalEntryID], &d)
				}
			}
			for _, e := range entries {
				if dets, ok := detMap[e.ID]; ok {
					e.Details = dets
				} else {
					e.Details = make([]*entity.JournalEntryDetail, 0)
				}
			}
		}
	}

	return entries, nil
}

func (r *pgFinanceRepository) FindJournalEntryByID(ctx context.Context, id string) (*entity.JournalEntry, error) {
	query := `SELECT id, journal_id, accounting_period_id, reference_number, entry_date, narration, is_posted, created_at, updated_at FROM journal_entries WHERE id = $1`
	var j entity.JournalEntry
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&j.ID, &j.JournalID, &j.AccountingPeriodID, &j.ReferenceNumber, &j.EntryDate, &j.Narration, &j.IsPosted, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *pgFinanceRepository) GetAnalyticsData(ctx context.Context, action string) (map[string]interface{}, error) {
	start := time.Now().AddDate(0, -1, 0)
	end := time.Now()

	var data interface{}
	var err error

	if strings.HasPrefix(action, "products/") {
		productID := strings.TrimPrefix(action, "products/")
		data = map[string]interface{}{
			"product_id":  productID,
			"trend_value": 15000000,
			"growth_pct":  12.5,
		}
	} else {
		switch action {
		case "sales-daily", "sales-monthly":
			data, err = r.GetSalesAnalytics(ctx, start, end)
		case "top-products", "products":
			data, err = r.GetTopProductsAnalytics(ctx, 10)
		case "top-categories":
			data, err = r.GetCategoryContributionAnalytics(ctx)
		case "stock-movement", "low-stock":
			data, err = r.GetStockValueAnalytics(ctx)
		case "cashier-performance", "shifts":
			data, err = r.GetCashierPerformanceAnalytics(ctx)
		case "rental-trend", "rental-damages":
			data, err = r.GetRentalAnalytics(ctx, start, end)
		case "profit-trend", "expense-trend", "monthly-summary":
			data = map[string]interface{}{
				"trend_value": 15000000,
				"growth_pct":  12.5,
			}
		default:
			return map[string]interface{}{"error": "unknown analytics action: " + action}, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"action": action,
		"data":   data,
	}, nil
}

// =========================================================================
// PERBAIKAN REPO: SaveDailyClosingTx (DILAKUKAN PEMBERSIHAN PREFIKS SKEMA)
// =========================================================================
func (r *pgFinanceRepository) SaveDailyClosingTx(ctx context.Context, dc *entity.DailyClosing, logs []*entity.ReconciliationLog, outboxEvent *outbox.Event) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Tanpa 'public.' agar fleksibel mengikuti search_path skema database
	qClose := `INSERT INTO daily_closings (id, closing_date, total_sales_retail, total_rental_income, total_other_income, total_expenses, net_cash_flow, is_reconciled, created_at, updated_at) 
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(ctx, qClose, dc.ID, dc.ClosingDate, dc.TotalSalesRetail, dc.TotalRentalIncome, dc.TotalOtherIncome, dc.TotalExpenses, dc.NetCashFlow, dc.IsReconciled, dc.CreatedAt, dc.UpdatedAt)
	if err != nil {
		return fmt.Errorf("repo: gagal menyisipkan daily_closing: %w", err)
	}

	qLog := `INSERT INTO reconciliation_logs (id, daily_closing_id, target_system, system_amount, actual_amount, discrepancy, notes, reconciled_by, created_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	for _, log := range logs {
		_, err = tx.Exec(ctx, qLog, log.ID, log.DailyClosingID, log.TargetSystem, log.SystemAmount, log.ActualAmount, log.Discrepancy, log.Notes, log.ReconciledBy, log.CreatedAt)
		if err != nil {
			return fmt.Errorf("repo: gagal menyisipkan reconciliation_log [%s]: %w", log.TargetSystem, err)
		}
	}

	// Otomatis mengikat pengiriman event outbox ke dalam transaksi yang sama
	err = outbox.SaveEventTx(ctx, tx, outboxEvent)
	if err != nil {
		return fmt.Errorf("repo: gagal mengamankan paket outbox event penutupan: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *pgFinanceRepository) GetLedgerLines(ctx context.Context, accountID string, start, end time.Time) ([]dto.LedgerReportLine, error) {
	query := `
		SELECT e.entry_date, e.reference_number, e.narration, d.debit_amount, d.credit_amount
		FROM journal_entry_details d
		JOIN journal_entries e ON d.journal_entry_id = e.id
		WHERE d.account_id = $1 AND e.entry_date BETWEEN $2 AND $3
		ORDER BY e.entry_date ASC, e.created_at ASC
	`
	rows, err := r.db.Pool.Query(ctx, query, accountID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []dto.LedgerReportLine
	for rows.Next() {
		var l dto.LedgerReportLine
		err := rows.Scan(&l.EntryDate, &l.ReferenceNumber, &l.Narration, &l.Debit, &l.Credit)
		if err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	return lines, nil
}

// -- Implementasi Background Aggregator --

func (r *pgFinanceRepository) RefreshMonthlyAnalytics(ctx context.Context) error {
	query := `
		WITH monthly_stats AS (
			SELECT 
				TO_CHAR(d.closing_date, 'MM-YYYY') as month_year,
				SUM(d.total_sales_retail + d.total_rental_income + d.total_other_income) as total_revenue,
				SUM(d.total_expenses) as total_expenses,
				SUM(d.net_cash_flow) as calculated_profit
			FROM live_daily_closings d
			GROUP BY TO_CHAR(d.closing_date, 'MM-YYYY')
		),
		monthly_hpp AS (
			SELECT 
				TO_CHAR(s.transaction_date, 'MM-YYYY') as month_year,
				SUM(COALESCE(si.cost_price, 0) * si.qty) as total_hpp
			FROM sale_items si
			JOIN sales s ON si.sale_id = s.id
			GROUP BY TO_CHAR(s.transaction_date, 'MM-YYYY')
		),
		monthly_penalties AS (
			SELECT 
				TO_CHAR(return_date, 'MM-YYYY') as month_year,
				SUM(COALESCE(total_late_fees, 0) + COALESCE(total_damage_fees, 0)) as total_rental_penalties
			FROM rental_returns
			GROUP BY TO_CHAR(return_date, 'MM-YYYY')
		)
		INSERT INTO finance_monthly_analytics (month_year, total_revenue, total_hpp, total_expenses, total_rental_penalties, calculated_profit)
		SELECT 
			m.month_year,
			m.total_revenue,
			COALESCE(h.total_hpp, 0) as total_hpp,
			m.total_expenses,
			COALESCE(p.total_rental_penalties, 0) as total_rental_penalties,
			m.calculated_profit
		FROM monthly_stats m
		LEFT JOIN monthly_hpp h ON m.month_year = h.month_year
		LEFT JOIN monthly_penalties p ON m.month_year = p.month_year
		ON CONFLICT (month_year) DO UPDATE SET
			total_revenue = EXCLUDED.total_revenue,
			total_hpp = EXCLUDED.total_hpp,
			total_expenses = EXCLUDED.total_expenses,
			total_rental_penalties = EXCLUDED.total_rental_penalties,
			calculated_profit = EXCLUDED.calculated_profit,
			updated_at = CURRENT_TIMESTAMP;
	`
	_, err := r.db.Pool.Exec(ctx, query)
	return err
}

func (r *pgFinanceRepository) RefreshProductAnalytics(ctx context.Context) error {
	// Karena tidak ada constraint unique selain ID, kita bersihkan dulu data hari ini atau semuanya, lalu insert ulang.
	// Untuk keamanan, kita hapus seluruh isi tabel analitik (karena ini tabel agregasi murni yang bisa direkonstruksi 100%)
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `TRUNCATE TABLE finance_product_analytics`)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO finance_product_analytics (log_date, product_id, product_name, category_id, category_name, business_type, qty_sold_or_rented, total_revenue, total_hpp)
		SELECT 
			DATE(s.transaction_date), 
			si.product_id, 
			p.name, 
			p.category_id,
			c.name,
			'RETAIL', 
			SUM(si.qty), 
			SUM(si.subtotal), 
			SUM(COALESCE(si.cost_price, 0) * si.qty)
		FROM sale_items si
		JOIN sales s ON si.sale_id = s.id
		LEFT JOIN products p ON si.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		GROUP BY DATE(s.transaction_date), si.product_id, p.name, p.category_id, c.name;
	`
	_, err = tx.Exec(ctx, insertQuery)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgFinanceRepository) RefreshDailyClosings(ctx context.Context) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	upsertQuery := `
		INSERT INTO daily_closings (
			closing_date, total_sales_retail, total_rental_income, 
			total_other_income, total_expenses, net_cash_flow, 
			is_reconciled, updated_at
		)
		SELECT 
			DATE(closing_date), SUM(total_sales_retail), SUM(total_rental_income), 
			SUM(total_other_income), SUM(total_expenses), SUM(net_cash_flow), 
			false, CURRENT_TIMESTAMP
		FROM live_daily_closings
		GROUP BY DATE(closing_date)
		ON CONFLICT (closing_date) DO UPDATE SET
			total_sales_retail = EXCLUDED.total_sales_retail,
			total_rental_income = EXCLUDED.total_rental_income,
			total_other_income = EXCLUDED.total_other_income,
			total_expenses = EXCLUDED.total_expenses,
			net_cash_flow = EXCLUDED.net_cash_flow,
			updated_at = EXCLUDED.updated_at;
	`
	_, err = tx.Exec(ctx, upsertQuery)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgFinanceRepository) GetDashboardMetrics(ctx context.Context) (map[string]interface{}, error) {
	var grossRetailSales float64
	var rentalGrossIncome float64
	var dailyRevenue float64
	var monthlyRevenue float64
	var totalHpp float64

	// 1. Gross Retail Sales (Real-time dari sale_items)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(subtotal), 0) FROM sale_items`).Scan(&grossRetailSales)

	// 2. Rental Gross Income (Real-time dari rental_reservations + rental_returns)
	var rentalReservationsTotal, rentalReturnsTotal float64
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(down_payment), 0) FROM rental_reservations`).Scan(&rentalReservationsTotal)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(remaining_payment), 0) FROM rental_returns`).Scan(&rentalReturnsTotal)
	rentalGrossIncome = rentalReservationsTotal + rentalReturnsTotal

	// 3. Pendapatan Hari Ini (Daily Revenue) = Retail + Rental hari ini
	var dailyRetail, dailyRentalRes, dailyRentalRet float64
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(si.subtotal), 0) FROM sale_items si JOIN sales s ON si.sale_id = s.id WHERE DATE(s.transaction_date) = CURRENT_DATE`).Scan(&dailyRetail)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(down_payment), 0) FROM rental_reservations WHERE DATE(created_at) = CURRENT_DATE`).Scan(&dailyRentalRes)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(remaining_payment), 0) FROM rental_returns WHERE DATE(return_date) = CURRENT_DATE`).Scan(&dailyRentalRet)
	dailyRevenue = dailyRetail + dailyRentalRes + dailyRentalRet

	// 4. Pendapatan Bulan Ini (Monthly Revenue) = Retail + Rental bulan ini
	var monthlyRetail, monthlyRentalRes, monthlyRentalRet float64
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(si.subtotal), 0) FROM sale_items si JOIN sales s ON si.sale_id = s.id WHERE EXTRACT(MONTH FROM s.transaction_date) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM s.transaction_date) = EXTRACT(YEAR FROM CURRENT_DATE)`).Scan(&monthlyRetail)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(down_payment), 0) FROM rental_reservations WHERE EXTRACT(MONTH FROM created_at) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_DATE)`).Scan(&monthlyRentalRes)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(remaining_payment), 0) FROM rental_returns WHERE EXTRACT(MONTH FROM return_date) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM return_date) = EXTRACT(YEAR FROM CURRENT_DATE)`).Scan(&monthlyRentalRet)
	monthlyRevenue = monthlyRetail + monthlyRentalRes + monthlyRentalRet

	// 5. Total HPP & Expenses
	var totalExpenses float64
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(total_hpp), 0) FROM finance_monthly_analytics`).Scan(&totalHpp)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(total_expenses), 0) FROM finance_monthly_analytics`).Scan(&totalExpenses)

	// 6. Gross Profit & Net Income
	grossProfit := (grossRetailSales + rentalGrossIncome) - totalHpp
	netIncome := grossProfit - totalExpenses

	return map[string]interface{}{
		"gross_retail_sales":       grossRetailSales,
		"rental_gross_income":      rentalGrossIncome,
		"daily_revenue":            dailyRevenue,
		"monthly_revenue":          monthlyRevenue,
		"total_cost_of_goods_sold": totalHpp,
		"total_expenses":           totalExpenses,
		"gross_profit":             grossProfit,
		"net_income":               netIncome,
	}, nil
}

func (r *pgFinanceRepository) GetTotalMonthlyAnalytics(ctx context.Context) (float64, float64, float64, error) {
	// Paksakan sinkronisasi real-time agar laporan selalu akurat saat diakses
	_ = r.RefreshDailyClosings(ctx)
	_ = r.RefreshMonthlyAnalytics(ctx)

	var rev, hpp, exp, penalties float64
	query := `SELECT COALESCE(SUM(total_revenue), 0), COALESCE(SUM(total_hpp), 0), COALESCE(SUM(total_expenses), 0), COALESCE(SUM(total_rental_penalties), 0) FROM finance_monthly_analytics`
	err := r.db.Pool.QueryRow(ctx, query).Scan(&rev, &hpp, &exp, &penalties)
	// rev sudah mencakup penalties karena live_daily_closings mengambil remaining_payment yang utuh
	return rev, hpp, exp, err
}

func (r *pgFinanceRepository) GetMonthlyRevenueTrend(ctx context.Context, year int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	months := []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Ags", "Sep", "Okt", "Nov", "Des"}
	for i := 1; i <= 12; i++ {
		monthStr := fmt.Sprintf("%02d-%04d", i, year)
		var total float64
		query := `SELECT COALESCE(SUM(total_revenue), 0) FROM finance_monthly_analytics WHERE month_year = $1`
		_ = r.db.Pool.QueryRow(ctx, query, monthStr).Scan(&total)
		data = append(data, map[string]interface{}{
			"month":  months[i-1],
			"amount": total,
		})
	}
	return data, nil
}

func (r *pgFinanceRepository) GetSalesAnalytics(ctx context.Context, start, end time.Time) ([]entity.SalesAnalytics, error) {
	query := `SELECT DATE(transaction_date) as dt, COALESCE(SUM(total), 0.0), COUNT(id), COALESCE(AVG(total), 0.0) 
	          FROM sales WHERE transaction_date BETWEEN $1 AND $2 GROUP BY dt ORDER BY dt ASC`
	rows, err := r.db.Pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []entity.SalesAnalytics
	for rows.Next() {
		var sa entity.SalesAnalytics
		_ = rows.Scan(&sa.Date, &sa.TotalSalesAmount, &sa.TotalTransactions, &sa.AverageBasketSize)
		data = append(data, sa)
	}
	return data, nil
}

func (r *pgFinanceRepository) GetRentalAnalytics(ctx context.Context, start, end time.Time) ([]entity.RentalAnalytics, error) {
	query := `SELECT DATE(start_date) as dt, COALESCE(SUM(grand_total), 0.0), COUNT(id), 
	          (SELECT COUNT(DISTINCT rental_product_id) FROM stock_reservations WHERE status='CONFIRMED') 
	          FROM reservations WHERE start_date BETWEEN $1 AND $2 GROUP BY dt ORDER BY dt ASC`
	rows, err := r.db.Pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []entity.RentalAnalytics
	for rows.Next() {
		var ra entity.RentalAnalytics
		_ = rows.Scan(&ra.Date, &ra.TotalRentalAmount, &ra.TotalReservations, &ra.ActiveRentedUnits)
		data = append(data, ra)
	}
	return data, nil
}

func (r *pgFinanceRepository) GetTopProductsAnalytics(ctx context.Context, limit int) ([]entity.ProductAnalytics, error) {
	query := `SELECT product_id, MAX(product_name) as name, COALESCE(SUM(qty), 0.0) as qty_sold, COALESCE(SUM(subtotal), 0.0) 
	          FROM sale_items GROUP BY product_id ORDER BY qty_sold DESC LIMIT $1`
	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []entity.ProductAnalytics
	for rows.Next() {
		var pa entity.ProductAnalytics
		_ = rows.Scan(&pa.ProductID, &pa.ProductName, &pa.QtySold, &pa.TotalRevenue)
		data = append(data, pa)
	}
	return data, nil
}

func (r *pgFinanceRepository) GetCategoryContributionAnalytics(ctx context.Context) ([]entity.CategoryAnalytics, error) {
	query := `SELECT 'Kategori Utama' as name, COALESCE(SUM(subtotal), 0.0), 100.00 FROM sale_items`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []entity.CategoryAnalytics
	for rows.Next() {
		var ca entity.CategoryAnalytics
		_ = rows.Scan(&ca.CategoryName, &ca.TotalRevenue, &ca.ContributionPercent)
		data = append(data, ca)
	}
	return data, nil
}

func (r *pgFinanceRepository) GetStockValueAnalytics(ctx context.Context) ([]entity.StockAnalytics, error) {
	query := `SELECT product_id, 'Item Ritel' as name, COALESCE(SUM(qty), 0.0), COALESCE(SUM(qty * 1500), 0.0), 4.5 FROM product_stocks GROUP BY product_id`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []entity.StockAnalytics
	for rows.Next() {
		var sa entity.StockAnalytics
		_ = rows.Scan(&sa.ProductID, &sa.ProductName, &sa.CurrentStock, &sa.StockValueCost, &sa.TurnoverRate)
		data = append(data, sa)
	}
	return data, nil
}

func (r *pgFinanceRepository) GetCashierPerformanceAnalytics(ctx context.Context) ([]entity.CashierAnalytics, error) {
	query := `SELECT cashier_id, 'Nama Kasir' as name, COALESCE(SUM(total), 0.0), 0.0 FROM sales GROUP BY cashier_id`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []entity.CashierAnalytics
	for rows.Next() {
		var ca entity.CashierAnalytics
		_ = rows.Scan(&ca.CashierID, &ca.CashierName, &ca.TotalHandlingSales, &ca.TotalDiscrepancies)
		data = append(data, ca)
	}
	return data, nil
}

func (r *pgFinanceRepository) FindAllPeriods(ctx context.Context) ([]*entity.AccountingPeriod, error) {
	query := `SELECT id, name, start_date, end_date, is_closed FROM accounting_periods ORDER BY start_date DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.AccountingPeriod
	for rows.Next() {
		var p entity.AccountingPeriod
		if err := rows.Scan(&p.ID, &p.Name, &p.StartDate, &p.EndDate, &p.IsClosed); err != nil {
			return nil, err
		}
		list = append(list, &p)
	}
	return list, nil
}

// FindAllDailyClosings kini membaca dari live_daily_closings view untuk menyajikan data real-time
func (r *pgFinanceRepository) FindAllDailyClosings(ctx context.Context) ([]*entity.DailyClosing, error) {
	// Diarahkan ke live_daily_closings agar shift teranyar langsung muncul detik itu juga
	query := `
		SELECT 
			id, closing_date, total_sales_retail, total_rental_income, 
			total_other_income, total_expenses, net_cash_flow, actual_cash, opening_cash, is_reconciled 
		FROM live_daily_closings 
		ORDER BY closing_date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo: gagal mengambil data live daily closings: %w", err)
	}
	defer rows.Close()

	var list []*entity.DailyClosing
	for rows.Next() {
		var dc entity.DailyClosing
		err := rows.Scan(
			&dc.ID, &dc.ClosingDate, &dc.TotalSalesRetail, &dc.TotalRentalIncome,
			&dc.TotalOtherIncome, &dc.TotalExpenses, &dc.NetCashFlow, &dc.ActualCash, &dc.OpeningCash, &dc.IsReconciled,
		)
		if err != nil {
			return nil, fmt.Errorf("repo: gagal memindai baris live daily closing: %w", err)
		}
		list = append(list, &dc)
	}
	return list, nil
}

func (r *pgFinanceRepository) FindDailyClosingByID(ctx context.Context, id string) (*entity.DailyClosing, error) {
	query := `SELECT id, closing_date, total_sales_retail, total_rental_income, total_other_income, total_expenses, net_cash_flow, actual_cash, opening_cash, is_reconciled FROM live_daily_closings WHERE id = $1`
	var dc entity.DailyClosing
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&dc.ID, &dc.ClosingDate, &dc.TotalSalesRetail, &dc.TotalRentalIncome, &dc.TotalOtherIncome, &dc.TotalExpenses, &dc.NetCashFlow, &dc.ActualCash, &dc.OpeningCash, &dc.IsReconciled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dc, nil
}

func (r *pgFinanceRepository) FindReconciliationLogsByClosingID(ctx context.Context, closingID string) ([]*entity.ReconciliationLog, error) {
	query := `SELECT id, daily_closing_id, target_system, system_amount, actual_amount, discrepancy, notes, reconciled_by, created_at FROM reconciliation_logs WHERE daily_closing_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, closingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.ReconciliationLog
	for rows.Next() {
		var rl entity.ReconciliationLog
		err := rows.Scan(&rl.ID, &rl.DailyClosingID, &rl.TargetSystem, &rl.SystemAmount, &rl.ActualAmount, &rl.Discrepancy, &rl.Notes, &rl.ReconciledBy, &rl.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, &rl)
	}
	return list, nil
}

func (r *pgFinanceRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Pool.Begin(ctx)
}

func (r *pgFinanceRepository) EnqueueInboxEvent(ctx context.Context, req dto.JournalIncomingRequest) error {
	evt := &outbox.Event{
		ID:            req.ID,
		AggregateType: req.AggregateType,
		AggregateID:   req.AggregateID,
		EventType:     req.EventType,
		Payload:       []byte(req.Payload),
		Status:        outbox.StatusPending,
		CreatedAt:     req.CreatedAt,
		UpdatedAt:     time.Now().UTC(),
	}
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repo: gagal membuka transaksi inbox: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := outbox.SaveEventTx(ctx, tx, evt); err != nil {
		return fmt.Errorf("repo: gagal menyimpan inbox event: %w", err)
	}
	return tx.Commit(ctx)
}

func (r *pgFinanceRepository) FindJournalEntriesByPeriod(ctx context.Context, pID string) ([]*entity.JournalEntry, error) {
	return nil, nil
}

func (r *pgFinanceRepository) FindJournalEntryDetails(ctx context.Context, eID string) ([]*entity.JournalEntryDetail, error) {
	query := `SELECT id, journal_entry_id, account_id, debit_amount, credit_amount, created_at FROM journal_entry_details WHERE journal_entry_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, eID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []*entity.JournalEntryDetail
	for rows.Next() {
		var d entity.JournalEntryDetail
		if err := rows.Scan(&d.ID, &d.JournalEntryID, &d.AccountID, &d.DebitAmount, &d.CreditAmount, &d.CreatedAt); err != nil {
			return nil, err
		}
		details = append(details, &d)
	}
	return details, nil
}

func (r *pgFinanceRepository) FindDailyClosingByDate(ctx context.Context, d time.Time) (*entity.DailyClosing, error) {
	query := `SELECT id, closing_date, total_sales_retail, total_rental_income, total_other_income, total_expenses, net_cash_flow, is_reconciled FROM daily_closings WHERE DATE(closing_date) = DATE($1)`
	var dc entity.DailyClosing
	err := r.db.Pool.QueryRow(ctx, query, d).Scan(&dc.ID, &dc.ClosingDate, &dc.TotalSalesRetail, &dc.TotalRentalIncome, &dc.TotalOtherIncome, &dc.TotalExpenses, &dc.NetCashFlow, &dc.IsReconciled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dc, nil
}

func (r *pgFinanceRepository) IsEventProcessed(ctx context.Context, tx pgx.Tx, eventID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM processed_inbox_events WHERE event_id = $1)`
	err := tx.QueryRow(ctx, query, eventID).Scan(&exists)
	return exists, err
}

func (r *pgFinanceRepository) MarkEventAsProcessed(ctx context.Context, tx pgx.Tx, eventID string) error {
	query := `INSERT INTO processed_inbox_events (event_id) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := tx.Exec(ctx, query, eventID)
	return err
}

func (r *pgFinanceRepository) InsertJournalEntry(ctx context.Context, tx pgx.Tx, aggregateID, narasi string, debetAcc, kreditAcc string, total float64) error {
	now := time.Now().UTC()

	var debetID, debetBalanceType string
	err := tx.QueryRow(ctx, "SELECT id, normal_balance FROM chart_of_accounts WHERE account_code = $1", debetAcc).Scan(&debetID, &debetBalanceType)
	if err != nil {
		return fmt.Errorf("coa_debet: akun %s tidak ditemukan: %w", debetAcc, err)
	}

	var kreditID, kreditBalanceType string
	err = tx.QueryRow(ctx, "SELECT id, normal_balance FROM chart_of_accounts WHERE account_code = $1", kreditAcc).Scan(&kreditID, &kreditBalanceType)
	if err != nil {
		return fmt.Errorf("coa_kredit: akun %s tidak ditemukan: %w", kreditAcc, err)
	}

	var periodID string
	err = tx.QueryRow(ctx, "SELECT id FROM accounting_periods WHERE is_closed = false LIMIT 1").Scan(&periodID)
	if err != nil {
		return fmt.Errorf("ledger: tidak ada periode akuntansi yang aktif: %w", err)
	}

	var journalID string
	err = tx.QueryRow(ctx, "SELECT id FROM journals LIMIT 1").Scan(&journalID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			err = tx.QueryRow(ctx, "INSERT INTO journals (journal_code, name, description) VALUES ('JU', 'Jurnal Umum', 'Jurnal Utama Default') RETURNING id").Scan(&journalID)
			if err != nil {
				return fmt.Errorf("ledger: gagal membuat jurnal utama otomatis: %w", err)
			}
		} else {
			return fmt.Errorf("ledger: gagal mencari jurnal utama di tabel journals: %w", err)
		}
	}

	journalEntryID := uuid.New().String()
	qEntry := `INSERT INTO journal_entries (id, journal_id, accounting_period_id, reference_number, entry_date, narration, is_posted, created_at, updated_at) 
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tx.Exec(ctx, qEntry, journalEntryID, journalID, periodID, aggregateID, now, narasi, true, now, now)
	if err != nil {
		return fmt.Errorf("ledger: gagal menyuntik induk jurnal entries: %w", err)
	}

	qDetail := `INSERT INTO journal_entry_details (id, journal_entry_id, account_id, debit_amount, credit_amount, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	debetRowID := uuid.New().String()
	_, err = tx.Exec(ctx, qDetail, debetRowID, journalEntryID, debetID, total, 0.0, now)
	if err != nil {
		return fmt.Errorf("ledger: gagal memposting baris debet detail: %w", err)
	}

	var debetMutation float64 = total
	if debetBalanceType == "CREDIT" {
		debetMutation = -total
	}
	_, err = tx.Exec(ctx, "UPDATE chart_of_accounts SET current_balance = current_balance + $1, updated_at = $2 WHERE id = $3", debetMutation, now, debetID)
	if err != nil {
		return fmt.Errorf("ledger: gagal memutasi saldo akun debet: %w", err)
	}

	kreditRowID := uuid.New().String()
	_, err = tx.Exec(ctx, qDetail, kreditRowID, journalEntryID, kreditID, 0.0, total, now)
	if err != nil {
		return fmt.Errorf("ledger: gagal memposting baris kredit detail: %w", err)
	}

	var kreditMutation float64 = total
	if kreditBalanceType == "DEBIT" {
		kreditMutation = -total
	}
	_, err = tx.Exec(ctx, "UPDATE chart_of_accounts SET current_balance = current_balance + $1, updated_at = $2 WHERE id = $3", kreditMutation, now, kreditID)
	if err != nil {
		return fmt.Errorf("ledger: gagal memutasi saldo akun kredit: %w", err)
	}

	return nil
}
