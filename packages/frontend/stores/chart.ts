import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

// Mengimpor interface dari shared package
import type { SaleItem } from '@frontend/types/pos';
import type { Product } from '@frontend/types/inventory';
import { posApi } from '@frontend/api/pos';
import { cashApi } from '@frontend/api/cash';

export const useCartStore = defineStore('cart', () => {
    const items = ref<SaleItem[]>([]);
    const discountAmount = ref<number>(0);
    const taxRate = ref<number>(0); // PPN 11%

    // Perhitungan Subtotal dengan penulisan tipe eksplisit
    const subtotal = computed(() => {
        return items.value.reduce((total: number, item: SaleItem) => total + item.subtotal, 0);
    });

    // Perhitungan Pajak
    const tax = computed(() => {
        return (subtotal.value - discountAmount.value) * taxRate.value;
    });

    // Perhitungan Total Akhir
    const grandTotal = computed(() => {
        return subtotal.value - discountAmount.value + tax.value;
    });

    // Fungsi menambah barang ke keranjang
    function addItem(product: Product, qty: number = 1) {
        const existingItem = items.value.find((item: SaleItem) => item.product_id === product.id);

        if (existingItem) {
            existingItem.qty += qty;
            const rawSubtotal = (existingItem.qty * existingItem.unit_price) - existingItem.discount;
            existingItem.subtotal = rawSubtotal > 0 ? rawSubtotal : 0;
        } else {
            items.value.push({
                product_id: product.id,
                product_name: product.name,
                unit_code: product.base_unit_code,
                qty: qty,
                unit_price: product.selling_price,
                discount: 0,
                subtotal: qty * product.selling_price,
                cost_price: 0
            });
        }
    }

    // Fungsi mengubah jumlah barang
    function updateItemQty(productId: string, qty: number) {
        const item = items.value.find((i: SaleItem) => i.product_id === productId);
        if (item && qty > 0) {
            item.qty = qty;
            const rawSubtotal = (item.qty * item.unit_price) - item.discount;
            item.subtotal = rawSubtotal > 0 ? rawSubtotal : 0;
        } else if (item && qty === 0) {
            removeItem(productId);
        }
    }

    // Fungsi mengubah diskon per barang
    function updateItemDiscount(productId: string, discount: number) {
        const item = items.value.find((i: SaleItem) => i.product_id === productId);
        if (item) {
            item.discount = discount >= 0 ? discount : 0;
            const rawSubtotal = (item.qty * item.unit_price) - item.discount;
            item.subtotal = rawSubtotal > 0 ? rawSubtotal : 0;
        }
    }

    // Fungsi menghapus barang dari keranjang
    function removeItem(productId: string) {
        items.value = items.value.filter((item: SaleItem) => item.product_id !== productId);
    }

    // Mengosongkan keranjang
    function clearCart() {
        items.value = [];
        discountAmount.value = 0;
    }

    // Fungsi pembantu untuk membuat UUID jika crypto.randomUUID tidak tersedia
    function generateUUID() {
        if (typeof crypto !== 'undefined' && crypto.randomUUID) {
            return crypto.randomUUID();
        }
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            const r = Math.random() * 16 | 0;
            const v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    const isCheckoutLoading = ref(false);

    // Proses Checkout
    async function checkout(paymentMethod: string = 'CASH', amountPaid?: number) {
        if (items.value.length === 0) return null;

        isCheckoutLoading.value = true;
        try {
            // Ambil sesi kasir aktif
            const session = await cashApi.getCurrentSession();
            if (!session || !session.id) {
                throw new Error("Sesi kasir (shift) tidak aktif. Silakan buka shift Anda terlebih dahulu di menu Kas > Shift.");
            }

            const payload = {
                idempotency_key: generateUUID(),
                payment_method: paymentMethod,
                discount: discountAmount.value,
                amount_paid: amountPaid,
                cashier_session_id: session.id,
                items: items.value.map(i => ({
                    product_id: i.product_id,
                    product_name: i.product_name,
                    unit_code: i.unit_code,
                    qty: i.qty,
                    unit_price: i.unit_price,
                    discount: i.discount
                }))
            };

            const response = await posApi.checkout(payload);
            clearCart();
            return response; // Mengembalikan data misal invoice_number
        } catch (error: any) {
            console.error("Gagal melakukan checkout:", error);
            const msg = error.response?.data?.message || error.message || 'Unknown error';
            alert(`Terjadi kesalahan saat memproses pembayaran:\n${msg}`);
            throw error;
        } finally {
            isCheckoutLoading.value = false;
        }
    }

    return {
        items,
        discountAmount,
        taxRate,
        subtotal,
        tax,
        grandTotal,
        isCheckoutLoading,
        addItem,
        updateItemQty,
        updateItemDiscount,
        removeItem,
        clearCart,
        checkout
    };
});