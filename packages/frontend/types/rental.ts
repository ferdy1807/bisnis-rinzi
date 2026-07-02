export interface RentalProduct {
    id: string;
    category_id: string;
    code: string;
    name: string;
    description?: string;
    rental_price: number;
    deposit_amount: number;
    quantity_available: number;
    is_active: boolean;
    object_name?: string;
    original_file_name?: string;
    mime_type?: string;
}

export interface RentalReservation {
    id: string;
    invoice_number: string;
    customer_snapshot_id: string;
    customer_name?: string;
    customer_phone?: string;
    transaction_date: string;
    start_date: string;
    end_date: string;
    subtotal: number;
    down_payment: number;
    amount_paid: number;
    change_amount: number;
    total_amount: number;
    status: 'BOOKED' | 'CONTENTS_RECEIVED' | 'DECORATING' | 'READY_FOR_PICKUP' | 'PICKED_UP' | 'RETURNED' | 'CANCELLED';
    picked_up_by?: string;
    picked_up_at?: string;
    cashier_session_id: string;
    created_by: string;
    items?: RentalItem[];
    contents?: ReservationContent[];
}

export interface RentalReturn {
    id: string;
    rental_reservation_id: string;
    return_date: string;
    late_days: number;
    total_late_fees: number;
    total_damage_fees: number;
    amount_paid: number;
    change_amount: number;
    grand_total_paid: number;
    notes: string;
    received_by: string;
    receipt_url?: string;
    created_at: string;
    updated_at: string;
}

export interface CustomerSnapshot {
    id: string;
    customer_name: string;
    customer_phone: string;
    customer_id_card: string;
    created_at?: string;
    updated_at?: string;
}

export interface RentalCategory {
    id: string;
    code: string;
    name: string;
    created_at?: string;
    updated_at?: string;
    deleted_at?: string;
}

export interface StockReservation {
    id: string;
    rental_product_id: string;
    rental_reservation_id: string;
    reserved_date: string;
    qty_reserved: number;
    created_at?: string;
    updated_at?: string;
}

export interface RentalItem {
    id: string;
    rental_reservation_id: string;
    rental_product_id: string;
    rental_product_name: string;
    qty: number;
    price_per_period: number;
    subtotal: number;
    created_at?: string;
    updated_at?: string;
}


export interface RentalReturnItem {
    id: string;
    rental_return_id: string;
    rental_product_id: string;
    rental_product_name: string;
    qty_returned: number;
    condition_status: string;
    damage_fee: number;
    condition_notes?: string;
    created_at?: string;
    updated_at?: string;
}

export interface RentalReturnPhoto {
    id: string;
    rental_return_id: string;
    rental_return_item_id?: string;
    bucket_name: string;
    object_name: string;
    original_file_name?: string;
    mime_type?: string;
    file_size_bytes?: number;
    created_at?: string;
    updated_at?: string;
}

export interface ReservationContent {
    id: string;
    rental_reservation_id: string;
    item_name: string;
    description: string;
    quantity: number;
    condition_notes?: string;
    created_at?: string;
    updated_at?: string;
}