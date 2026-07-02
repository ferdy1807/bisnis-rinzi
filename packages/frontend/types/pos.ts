export interface SaleItem {
    id?: string;
    sale_id?: string;
    product_id: string;
    product_name: string;
    unit_code: string;
    qty: number;
    unit_price: number;
    discount: number;
    subtotal: number;
    cost_price: number;
}

export interface Sale {
    id: string;
    idempotency_key: string;
    invoice_number: string;
    transaction_date: string;
    subtotal: number;
    discount: number;
    total: number;
    amount_paid: number;
    change_amount: number;
    payment_method: string;
    payment_status: string;
    cashier_id: string;
    cashier_session_id: string;
    invoice_url?: string;
    items?: SaleItem[]; // Berisi relasi ke sale_items
}

export interface SyncLog {
    id: string;
    entity_type?: string;
    entity_id?: string;
    sync_status?: string;
    error_message?: string;
    created_at?: string;
    updated_at?: string;
}