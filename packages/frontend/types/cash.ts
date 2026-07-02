export interface CashierSession {
    id: string;
    cashier_id: string;
    open_time: string;
    close_time?: string;
    opening_cash: number;
    expected_cash?: number;
    actual_cash?: number;
    difference?: number;
    status: 'OPEN' | 'CLOSED';
    receipt_url?: string;
}

export interface ExpenseCategory {
    id?: string;       // Opsional karena di-generate otomatis oleh uuid_generate_v4() di DB
    code: string;      // Wajib (NOT NULL & UNIQUE)
    name: string;      // Wajib (NOT NULL)
    created_at?: string;
    updated_at?: string;
}

export interface Expense {
    id: string;
    expense_date: string;
    category_id: string;
    description: string;
    amount: number;
    created_by: string;
    created_at?: string;
}

export interface CashTransaction {
    id: string;
    session_id: string;
    transaction_type: string;
    reference_type: string;
    reference_id?: string;
    amount: number;
    notes?: string;
    created_by: string;
    created_at?: string;
    updated_at?: string;
}

export interface ShiftSummary {
    session_id: string;
    cashier_id: string;
    opening_cash: number;
    total_income: number;
    total_expense: number;
    total_deposit: number;
    total_withdrawal: number;
    expected_cash: number;
    actual_cash?: number;
    difference?: number;
    status: 'OPEN' | 'CLOSED';
}