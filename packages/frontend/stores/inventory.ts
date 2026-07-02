// apps/admin-dashboard/src/stores/inventory.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { inventoryApi } from '@frontend/api/inventory';
import type { Product, ProductStock, StockMovement } from '@frontend/types/inventory';

export const useInventoryStore = defineStore('inventory', () => {
    // ==========================================
    // 1. STATE MANAGEMENT TERPUSAT
    // ==========================================
    const products = ref<Product[]>([]);
    const lowStockProducts = ref<Product[]>([]);
    const stockMovements = ref<StockMovement[]>([]);
    const currentProductStock = ref<ProductStock | null>(null);

    // Bendera indikator pemuatan data untuk efisiensi antarmuka UI/UX
    const isLoading = ref({
        products: false,
        lowStock: false,
        movements: false,
        stockDetail: false,
        submitAction: false
    });

    const errorMessage = ref<string | null>(null);
    const successMessage = ref<string | null>(null);

    // ==========================================
    // 2. GETTERS (KOMPUTASI REAKTIF)
    // ==========================================
    // Menghitung total variasi produk aktif yang terdaftar di sistem katalog
    const totalActiveProducts = computed(() => {
        return products.value.filter((p: Product) => p.is_active).length;
    });

    // Menghitung jumlah produk yang saat ini berada di zona kritis di bawah ambang batas minimum
    const totalCriticalStockCount = computed(() => lowStockProducts.value.length);

    // Pemetaan produk berdasarkan SKU untuk pencarian instan pada modul POS/Kasir
    const productSkuMap = computed(() => {
        return products.value.reduce((map: Record<string, Product>, product: Product) => {
            if (product.sku) {
                map[product.sku] = product;
            }
            return map;
        }, {});
    });

    // ==========================================
    // 3. ACTIONS (OPERASI ASINKRON & MUTASI DATA)
    // ==========================================

    function clearMessages() {
        errorMessage.value = null;
        successMessage.value = null;
    }

    /**
     * Mengambil seluruh katalog produk master dari database inventory_db
     */
    async function fetchProducts(forceRefresh = false) {
        if (products.value.length > 0 && !forceRefresh) return;
        isLoading.value.products = true;
        clearMessages();
        try {
            // Memanggil API master inventory untuk menarik daftar produk
            const response = await inventoryApi.getProducts();
            products.value = response || [];
        } catch (error: any) {
            errorMessage.value = error?.message || 'Sistem gagal mengambil repositori katalog produk master.';
            throw error;
        } finally {
            isLoading.value.products = false;
        }
    }

    /**
     * Memindai dan mengambil daftar produk yang menyentuh batas minimum keselamatan stok (safety stock)
     */
    async function fetchLowStockProducts() {
        isLoading.value.lowStock = true;
        clearMessages();
        try {
            const response = await inventoryApi.getLowStockProducts();
            lowStockProducts.value = response || [];
        } catch (error: any) {
            errorMessage.value = error?.message || 'Gagal menarik data indikasi peringatan dini stok kritis.';
        } finally {
            isLoading.value.lowStock = false;
        }
    }

    /**
     * Memuat berkas jejak audit kronologis mutasi keluar-masuk barang
     */
    async function fetchStockMovements() {
        isLoading.value.movements = true;
        clearMessages();
        try {
            const response = await inventoryApi.getStockMovements();
            stockMovements.value = response || [];
        } catch (error: any) {
            errorMessage.value = error?.message || 'Gagal mengekstrak berkas log riwayat mutasi barang.';
        } finally {
            isLoading.value.movements = false;
        }
    }

    /**
     * Memeriksa rincian konfigurasi kuantitas stok minimum dan safety stock suatu produk
     */
    async function fetchProductStockDetail(productId: string) {
        isLoading.value.stockDetail = true;
        currentProductStock.value = null;
        clearMessages();
        try {
            const response = await inventoryApi.getProductStock(productId);
            currentProductStock.value = response || null;
            return response;
        } catch (error: any) {
            errorMessage.value = error?.message || 'Gagal memuat parameter ambang batas stok produk.';
            throw error;
        } finally {
            isLoading.value.stockDetail = false;
        }
    }

    /**
     * Mengeksekusi modul Stock Opname (Koreksi Selisih Fisik) dengan validasi ledger sistem
     */
    async function executeStockAdjustment(payload: { product_id: string; quantity: number; reason: string }) {
        isLoading.value.submitAction = true;
        clearMessages();
        try {
            // Mengirimkan permintaan koreksi fisik menuju endpoint Inventory API
            await inventoryApi.adjustStock(payload);

            successMessage.value = 'Data penyesuaian stok fisik berhasil direkam dan divalidasi ke sistem ledger.';

            // Sinkronisasi data lokal secara otomatis agar visualisasi dashboard tetap valid dan sinkron
            await Promise.all([
                fetchProducts(true),
                fetchLowStockProducts(),
                fetchStockMovements()
            ]);
        } catch (error: any) {
            errorMessage.value = error?.message || 'Gagal mengeksekusi instruksi penyesuaian stok opname.';
            throw error;
        } finally {
            isLoading.value.submitAction = false;
        }
    }

    return {
        // State
        products,
        lowStockProducts,
        stockMovements,
        currentProductStock,
        isLoading,
        errorMessage,
        successMessage,

        // Getters
        totalActiveProducts,
        totalCriticalStockCount,
        productSkuMap,

        // Actions
        clearMessages,
        fetchProducts,
        fetchLowStockProducts,
        fetchStockMovements,
        fetchProductStockDetail,
        executeStockAdjustment
    };
});