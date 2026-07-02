// packages/frontend/api/rental.ts
import { BaseApi } from './base';
import { getToken } from '../services/token';
import type {
    RentalReservation,
    RentalProduct,
    RentalReturn
} from '../types/rental';

export interface CreateReservationPayload {
    customer_identity: string;
    customer_name: string;
    customer_phone: string;
    start_date: string;
    end_date: string;
    down_payment: number;
    amount_paid: number;
    cashier_session_id?: string;
    items: Array<{
        rental_product_id: string;
        qty: number;
        price_per_period: number;
    }>;
}

export interface ProcessReturnPayload {
    reservation_id: string;
    amount_paid: number;
    change_amount: number;
    manual_damage_fee: number;
    manual_return_notes?: string;
    return_items: Array<{
        rental_product_id: string;
        condition_status: string;
        damage_fee: number;
        condition_notes?: string;
    }>;
}export class RentalApi extends BaseApi {
    // =========================================================
    // 1. DASHBOARD & OPERASIONAL
    // =========================================================
    getActiveReservations() {
        return this.http.get<RentalReservation[]>('/api/rental/reservations/active').then(res => res.data || []);
    }
    getActive() { return this.getActiveReservations(); } // <--- Penambal error 'getActive'

    getUpcomingReservations() {
        return this.http.get<RentalReservation[]>('/api/rental/reservations/upcoming').then(res => res.data || []);
    }
    getUpcoming() { return this.getUpcomingReservations(); } // <--- Penambal error 'getUpcoming'

    getOverdueReservations() {
        return this.http.get<RentalReservation[]>('/api/rental/reservations/overdue').then(res => res.data || []);
    }
    getOverdue() { return this.getOverdueReservations(); } // <--- Penambal error 'getOverdue'


    // =========================================================
    // 2. RESERVASI & KONTRAK
    // =========================================================
    getAllReservations() {
        return this.http.get<RentalReservation[]>('/api/rental/reservations').then(res => res.data || []);
    }

    getReservationDetail(id: string) {
        return this.http.get<RentalReservation>(`/api/rental/reservations/${id}`).then(res => res.data);
    }

    createReservation(payload: CreateReservationPayload) {
        return this.http.post<{ id: string, invoice_number: string }>('/api/rental/reservations', payload);
    }

    cancelReservation(id: string, penaltyFee: number = 0) {
        return this.http.post(`/api/rental/reservations/${id}/cancel`, {
            reservation_id: id,
            penalty_fee: penaltyFee
        });
    }


    // =========================================================
    // 3. LOGISTIK & PENERIMAAN/PENYERAHAN
    // =========================================================
    pickupReservation(id: string) {
        return this.http.post(`/api/rental/reservations/${id}/pickup`);
    }
    pickupUnit(id: string) { return this.pickupReservation(id); }

    // Disesuaikan dengan skema fisik tabel rental_reservation_contents
    saveDepositItems(id: string, payload: {
        item_name: string;
        description: string;
        quantity: number;
        condition_notes?: string;
    }) {
        return this.http.post(`/api/rental/reservations/${id}/contents_received`, payload);
    }

    rollbackPickupReservation(id: string) {
        return this.http.post(`/api/rental/reservations/${id}/undo_pickup`);
    }

    markReady(id: string) {
        return this.http.post(`/api/rental/reservations/${id}/ready`);
    }

    undoReady(id: string) {
        return this.http.post(`/api/rental/reservations/${id}/undo_ready`);
    }


    // =========================================================
    // 4. PENGEMBALIAN (RETURNS) & KATALOG
    // =========================================================
    processReturn(payload: ProcessReturnPayload) {
        return this.http.post<{ data: RentalReturn }>('/api/rental/returns', payload);
    }

    getReturnByReservationId(resId: string) {
        return this.http.get<RentalReturn>(`/api/rental/reservations/${resId}/return`);
    }

    async uploadReturnReceipt(returnId: string, receipt: Blob, filename: string) {
        const formData = new FormData();
        formData.append('receipt', receipt, filename);

        const response: any = await this.http.post(`/api/rental/returns/${returnId}/receipt`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        return response.data;
    }

    async uploadReservationInvoice(resId: string, invoice: Blob, filename: string) {
        const formData = new FormData();
        formData.append('invoice', invoice, filename);

        const response: any = await this.http.post(`/api/rental/reservations/${resId}/invoice`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        return response.data;
    }

    async uploadReturnPhoto(returnId: string, photo: File) {
        const formData = new FormData();
        formData.append('photo', photo);

        const response: any = await this.http.post(`/api/rental/returns/${returnId}/photos`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        return response.data;
    }

    getProducts() {
        return this.http.get<RentalProduct[]>('/api/rental/products').then(res => res.data || []);
    }

    getCategories() {
        return this.http.get<any[]>('/api/rental/categories').then(res => res.data);
    }

    createProduct(payload: any) {
        return this.http.post('/api/rental/products', payload);
    }

    updateProduct(id: string, payload: any) {
        return this.http.put(`/api/rental/products/${id}`, payload);
    }

    deleteProduct(id: string) {
        return this.http.delete(`/api/rental/products/${id}`);
    }

    async uploadProductMedia(id: string, photo: File, filename: string) {
        const formData = new FormData();
        formData.append('media', photo, filename);

        const response: any = await this.http.post(`/api/rental/products/${id}/media`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        return response.data;
    }

    checkAvailability(params: { product_id: string; start_date: string; end_date: string; qty: number }) {
        return this.http.get<{ is_available: boolean }>('/api/rental/availability', { params }).then(res => res.data);
    }
}

export const rentalApi = new RentalApi();