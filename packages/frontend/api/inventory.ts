// packages/frontend/api/inventory.ts
import { BaseApi } from './base';
import type { Product, Category, Brand, Unit, ProductCreateRequest, ProductStock, ProductCostHistory, ProductMedia, StockMovement } from '../types/inventory';

export class InventoryApi extends BaseApi {
    public async getProducts(): Promise<Product[]> {
        const response: any = await this.http.get('/api/inventory/products');
        // Interceptor sudah meng-unwrap { success, data } menjadi data langsung
        if (Array.isArray(response.data)) {
            return response.data;
        }
        // Fallback jika interceptor tidak aktif (respons masih terbungkus)
        if (response.data && Array.isArray(response.data.data)) {
            return response.data.data;
        }
        return [];
    }
    // --- CATEGORY ---
    public async getCategories(): Promise<Category[]> {
        const response: any = await this.http.get('/api/inventory/categories');
        if (Array.isArray(response.data)) {
            return response.data;
        }
        if (response.data && Array.isArray(response.data.data)) {
            return response.data.data;
        }
        return [];
    }
    public async createCategory(payload: { code: string; name: string }): Promise<Category> {
        const response: any = await this.http.post('/api/inventory/categories', payload);
        return response.data?.data ?? response.data;
    }
    public async updateCategory(id: string, payload: { code: string; name: string }): Promise<Category> {
        const response: any = await this.http.put(`/api/inventory/categories/${id}`, payload);
        return response.data?.data ?? response.data;
    }
    public async deleteCategory(id: string): Promise<any> {
        const response: any = await this.http.delete(`/api/inventory/categories/${id}`);
        return response.data;
    }

    // --- BRAND ---
    public async getBrands(): Promise<Brand[]> {
        const response: any = await this.http.get('/api/inventory/brands');
        if (Array.isArray(response.data)) {
            return response.data;
        }
        if (response.data && Array.isArray(response.data.data)) {
            return response.data.data;
        }
        return [];
    }
    public async createBrand(payload: { code: string; name: string }): Promise<Brand> {
        const response: any = await this.http.post('/api/inventory/brands', payload);
        return response.data?.data ?? response.data;
    }
    public async updateBrand(id: string, payload: { code: string; name: string }): Promise<Brand> {
        const response: any = await this.http.put(`/api/inventory/brands/${id}`, payload);
        return response.data?.data ?? response.data;
    }
    public async deleteBrand(id: string): Promise<any> {
        const response: any = await this.http.delete(`/api/inventory/brands/${id}`);
        return response.data;
    }

    // --- UNIT ---
    public async getUnits(): Promise<Unit[]> {
        const response: any = await this.http.get('/api/inventory/units');
        if (Array.isArray(response.data)) {
            return response.data;
        }
        if (response.data && Array.isArray(response.data.data)) {
            return response.data.data;
        }
        return [];
    }
    public async createUnit(payload: { code: string; name: string }): Promise<Unit> {
        const response: any = await this.http.post('/api/inventory/units', payload);
        return response.data?.data ?? response.data;
    }
    public async updateUnit(id: string, payload: { code: string; name: string }): Promise<Unit> {
        const response: any = await this.http.put(`/api/inventory/units/${id}`, payload);
        return response.data?.data ?? response.data;
    }
    public async deleteUnit(id: string): Promise<any> {
        const response: any = await this.http.delete(`/api/inventory/units/${id}`);
        return response.data;
    }

    // --- PRODUCT ---
    public async createProduct(payload: ProductCreateRequest): Promise<Product> {
        const response: any = await this.http.post('/api/inventory/products', payload);
        return response.data?.data ?? response.data ?? {};
    }

    public async updateProduct(id: string, payload: any): Promise<Product> {
        const response: any = await this.http.put(`/api/inventory/products/${id}`, payload);
        return response.data?.data ?? response.data ?? {};
    }

    public async deleteProduct(id: string): Promise<any> {
        const response: any = await this.http.delete(`/api/inventory/products/${id}`);
        return response.data;
    }

    public async getProductById(id: string): Promise<Product> {
        const response: any = await this.http.get(`/api/inventory/products/${id}`);
        return response.data?.data ?? response.data;
    }

    public async getProductByBarcode(barcode: string): Promise<Product> {
        const response: any = await this.http.get(`/api/inventory/products/barcode/${barcode}`);
        return response.data?.data ?? response.data;
    }

    public async getProductStock(id: string): Promise<ProductStock> {
        const response: any = await this.http.get(`/api/inventory/products/${id}/stock`);
        return response.data?.data ?? response.data;
    }

    public async updateStockThreshold(id: string, payload: { min_stock: number, safety_stock: number }): Promise<any> {
        const response: any = await this.http.put(`/api/inventory/products/${id}/thresholds`, payload);
        return response.data;
    }

    public async adjustStock(payload: any): Promise<any> {
        const response: any = await this.http.post('/api/inventory/stocks/adjust', payload);
        return response.data;
    }

    public async getCostHistories(id: string): Promise<ProductCostHistory[]> {
        const response: any = await this.http.get(`/api/inventory/products/${id}/cost-histories`);
        if (Array.isArray(response.data)) {
            return response.data;
        }
        return response.data?.data || [];
    }

    public async addCostHistory(id: string, payload: { average_cost: number }): Promise<any> {
        const response: any = await this.http.post(`/api/inventory/products/${id}/cost-histories`, payload);
        return response.data;
    }

    public async getProductMedia(id: string): Promise<ProductMedia[]> {
        const response: any = await this.http.get(`/api/inventory/products/${id}/media`);
        if (Array.isArray(response.data)) {
            return response.data;
        }
        return response.data?.data || [];
    }

    public async uploadProductMedia(id: string, formData: FormData): Promise<any> {
        const response: any = await this.http.post(`/api/inventory/products/${id}/media`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        return response.data;
    }

    public async deleteProductMedia(id: string, mediaId: string): Promise<any> {
        const response: any = await this.http.delete(`/api/inventory/products/${id}/media/${mediaId}`);
        return response.data;
    }

    public getMediaUrl(mediaId: string): string {
        const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
        return `${baseUrl}/api/public/inventory/media/${mediaId}`;
    }

    public async getStockMovements(): Promise<StockMovement[]> {
        const response: any = await this.http.get('/api/inventory/stock-movements');
        if (Array.isArray(response.data)) {
            return response.data;
        }
        return response.data?.data || [];
    }

    public async getLowStockProducts(): Promise<Product[]> {
        const response: any = await this.http.get('/api/inventory/products/low-stock');
        if (Array.isArray(response.data)) {
            return response.data;
        }
        return response.data?.data || [];
    }
}

export const inventoryApi = new InventoryApi();
