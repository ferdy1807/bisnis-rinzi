export interface Category {
    id: string;
    code: string;
    name: string;
    created_at?: string;
    updated_at?: string;
}

export interface Brand {
    id: string;
    code: string;
    name: string;
    created_at?: string;
    updated_at?: string;
}

export interface Unit {
    id: string;
    code: string;
    name: string;
    created_at?: string;
    updated_at?: string;
}

export interface ProductCreateRequest {
    sku: string;
    category_id: string;
    brand_id?: string;
    name: string;
    base_unit_code: string;
    cost_price: number;
    selling_price: number;
    barcode?: string;
    initial_qty: number;
}

export interface Product {
    id: string;
    sku: string;
    category_id: string;
    brand_id?: string;
    name: string;
    base_unit_code: string;
    cost_price: number;
    selling_price: number;
    is_active: boolean;
    barcode?: string;
    image?: string;
    qty?: number;
    stock?: number;
    created_at?: string;
    updated_at?: string;
}

export interface ProductStock {
    product_id: string;
    qty: number;
    qty_min_stock: number;
    qty_safety_stock: number;
    updated_at?: string;
}

export interface ProductMedia {
    id: string;
    product_id: string;
    media_category: string;
    bucket_name: string;
    object_name: string;
    original_file_name?: string;
    mime_type?: string;
    file_size_bytes?: number;
    is_active: boolean;
}

export interface StockMovement {
    id: string;
    product_id: string;
    product_name?: string;
    sku?: string;
    movement_type: string;
    qty: number;
    reference?: string;
    created_at?: string;
    updated_at?: string;
}

export interface ProductCostHistory {
    id: string;
    product_id: string;
    effective_date: string;
    average_cost: number;
    created_at?: string;
    updated_at?: string;
}