// packages/frontend/api/pos.ts
import { BaseApi } from './base';
import type { Sale } from '@frontend/types/pos';
export interface SaleItemPayload {
    product_id: string;
    product_name: string;
    unit_code: string;
    qty: number;
    unit_price: number;
    discount: number;
}

export interface CreateSaleRequest {
    idempotency_key: string;
    payment_method: string;
    discount: number;
    amount_paid?: number;
    cashier_session_id: string;
    items: SaleItemPayload[];
}

export class PosApi extends BaseApi {
    public async checkout(data: CreateSaleRequest): Promise<Sale> {
        const response: any = await this.http.post('/api/pos/sales', data);
        if (response.data && response.data.data) {
            return response.data.data;
        }
        return response.data as Sale;
    }

    public async uploadInvoice(saleId: string, pdfBlob: Blob): Promise<any> {
        const formData = new FormData();
        formData.append('invoice', pdfBlob, 'invoice.pdf');

        const response = await this.http.post(`/api/pos/sales/${saleId}/invoice`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        return response.data;
    }

    public async getTopProducts(sessionId?: string): Promise<{ product_name: string, total_qty: number, total_revenue: number }[]> {
        const url = sessionId ? `/api/pos/reports/top-products?session_id=${sessionId}` : '/api/pos/reports/top-products';
        const response: any = await this.http.get(url);
        return response.data?.data ?? response.data ?? [];
    }

    public async getSalesHistory(): Promise<Sale[]> {
        const response: any = await this.http.get('/api/pos/sales');
        return response.data?.data ?? response.data ?? [];
    }
}

export const posApi = new PosApi();
