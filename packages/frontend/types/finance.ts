// packages/frontend/types/finance.ts

/**
 * ENTITAS GUDANG DATA SQL (Berdasarkan skema finance_db di init.sql)
 */

export interface ChartOfAccount {
    id: string;
    account_code: string;
    account_name: string;
    account_group: 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'EXPENSE' | string;
    normal_balance: 'DEBIT' | 'CREDIT';
    current_balance: number;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
}

export interface AccountingPeriod {
    id: string;
    name: string;
    start_date: string;
    end_date: string;
    is_closed: boolean;
    created_at?: string;
    updated_at?: string;
}

export interface DailyClosing {
    id: string;
    closing_date: string;
    total_sales_retail: number;
    total_rental_income: number;
    total_other_income: number;
    total_expenses: number;
    net_cash_flow: number;
    actual_cash?: number;
    opening_cash?: number;
    is_reconciled: boolean;
    created_at?: string;
    updated_at?: string;
}

export interface ReconciliationLog {
    id: string;
    daily_closing_id: string;
    target_system: 'POS' | 'RENTAL' | 'CASH' | string;
    system_amount: number;
    actual_amount: number;
    discrepancy: number;
    notes?: string;
    reconciled_by?: string;
    created_at?: string;
}

export interface Journal {
    id: string;
    journal_code: string;
    name: string;
    description?: string;
    created_at?: string;
    updated_at?: string;
}

export interface JournalEntryDetail {
    id?: string;
    journal_entry_id?: string;
    account_id: string;
    debit_amount: number;
    credit_amount: number;
    created_at?: string;
}

export interface JournalEntry {
    id: string;
    journal_id: string;
    accounting_period_id: string;
    reference_number: string;
    entry_date: string;
    narration?: string;
    is_posted: boolean;
    details?: JournalEntryDetail[];
    created_at?: string;
    updated_at?: string;
}

export interface PeriodLock {
    id: string;
    accounting_period_id: string;
    locked_by: string;
    lock_reason?: string;
    created_at?: string;
    updated_at?: string;
}

/**
 * DATA TRANSFER OBJECTS (DTO) RESPONS (Berdasarkan finance_usecase.go / finance_handler.go)
 */

export interface LedgerLine {
    entry_date: string;
    reference_number: string;
    narration: string;
    debit: number;
    credit: number;
    running_balance: number;
}

export interface LedgerReportResponse {
    account_code: string;
    account_name: string;
    initial_balance: number;
    final_balance: number;
    lines: LedgerLine[];
}

export interface TrialBalanceLine {
    account_code: string;
    account_name: string;
    debit: number;
    credit: number;
}

export interface IncomeStatementResponse {
    period_name: string;
    total_revenue: number;
    total_cogs: number;
    gross_profit: number;
    total_expense: number;
    net_income: number;
    items: Record<string, number>;
}

export interface BalanceSheetResponse {
    period_name: string;
    total_assets: number;
    total_liab: number;
    total_equity: number;
    items: Record<string, number>;
}

export interface CorporateDashboardAnalytics {
    sales_trend?: any;
    rental_trend?: any;
    top_selling_products?: any[];
    category_contribution?: Record<string, number>;
    inventory_valuation?: any;
    cashier_performance?: any[];
}