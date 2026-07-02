import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { rentalApi, type CreateReservationPayload, type ProcessReturnPayload } from '@frontend/api/rental';
import type { RentalReservation, RentalProduct } from '@frontend/types/rental';

export const useRentalStore = defineStore('rental', () => {
    // ================= 1. STATE TERPUSAT =================
    const masterReservations = ref<RentalReservation[]>([]);
    const activeReservations = ref<RentalReservation[]>([]);
    const upcomingReservations = ref<RentalReservation[]>([]);
    const overdueReservations = ref<RentalReservation[]>([]);
    const products = ref<RentalProduct[]>([]);
    const currentDetail = ref<RentalReservation | null>(null);

    // Pemisahan bendera loading agar antarmuka tetap responsif
    const isLoading = ref({
        master: true,
        active: true,
        upcoming: true,
        overdue: true,
        products: true,
        detail: false,
        mutate: false // Digunakan khusus saat aksi POST/PUT berjalan
    });

    const errorMessage = ref<string | null>(null);

    // ================= 2. GETTERS (Komputasi Reaktif) =================
    const bookedList = computed(() => activeReservations.value.filter((r: RentalReservation) => r.status === 'BOOKED'));
    const readyList = computed(() => activeReservations.value.filter((r: RentalReservation) => r.status === 'READY_FOR_PICKUP'));
    const pickedUpList = computed(() => activeReservations.value.filter((r: RentalReservation) => r.status === 'PICKED_UP'));

    const totalOverduePotentialPenalties = computed(() => {
        return overdueReservations.value.reduce((acc: number, curr: RentalReservation) => {
            const overdueDays = calculateOverdueDays(curr.end_date);
            return acc + (overdueDays * 50000);
        }, 0);
    });


    // ================= 3. ACTIONS (Manipulasi Asinkron) =================

    function calculateOverdueDays(endDateStr: string): number {
        if (!endDateStr) return 0;
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const expectedEnd = new Date(endDateStr);
        expectedEnd.setHours(0, 0, 0, 0);
        const diffTime = today.getTime() - expectedEnd.getTime();
        const diffDays = Math.floor(diffTime / (1000 * 3600 * 24));
        return diffDays > 0 ? diffDays : 0;
    }

    function clearError() {
        errorMessage.value = null;
    }

    // --- FETCH ACTIONS ---
    async function fetchMasterReservations(forceRefresh = false) {
        if (masterReservations.value.length > 0 && !forceRefresh) return;
        isLoading.value.master = true;
        clearError();
        try {
            const response = await rentalApi.getAllReservations();
            masterReservations.value = response || [];
        } catch (error: any) {
            errorMessage.value = 'Sistem gagal mengambil seluruh arsip master reservasi';
            throw error;
        } finally {
            isLoading.value.master = false;
        }
    }

    async function fetchActiveReservations() {
        isLoading.value.active = true;
        clearError();
        try {
            const response = await rentalApi.getActiveReservations();
            activeReservations.value = response || [];
        } catch (error) {
            errorMessage.value = 'Gagal menarik data pesanan harian aktif';
        } finally {
            isLoading.value.active = false;
        }
    }

    async function fetchUpcomingReservations(forceRefresh?: boolean) {
        if (!forceRefresh && upcomingReservations.value.length > 0 && !isLoading.value.upcoming) return;
        isLoading.value.upcoming = true;
        try {
            const response = await rentalApi.getUpcomingReservations();
            upcomingReservations.value = response || [];
        } catch (error) {
            errorMessage.value = 'Gagal menarik daftar agenda pesanan mendatang';
        } finally {
            isLoading.value.upcoming = false;
        }
    }

    async function fetchOverdueReservations() {
        isLoading.value.overdue = true;
        try {
            const response = await rentalApi.getOverdueReservations();
            overdueReservations.value = response || [];
        } catch (error) {
            errorMessage.value = 'Gagal mendeteksi riwayat penunggakan box hantaran';
        } finally {
            isLoading.value.overdue = false;
        }
    }

    async function fetchProducts() {
        isLoading.value.products = true;
        try {
            const response = await rentalApi.getProducts();
            products.value = response || [];
        } catch (error) {
            errorMessage.value = 'Gagal memuat master katalog produk sewa';
        } finally {
            isLoading.value.products = false;
        }
    }

    async function loadReservationDetail(id: string) {
        isLoading.value.detail = true;
        currentDetail.value = null;
        try {
            const response = await rentalApi.getReservationDetail(id);
            currentDetail.value = response;
            return response;
        } catch (error) {
            errorMessage.value = 'Arsip nota tidak ditemukan di dalam database';
            throw error;
        } finally {
            isLoading.value.detail = false;
        }
    }

    // --- MUTATION ACTIONS (Write to Backend) ---
    async function executeCreateContract(payload: CreateReservationPayload) {
        isLoading.value.mutate = true;
        clearError();
        try {
            const response = await rentalApi.createReservation(payload);
            // Sinkronisasi ulang data lokal secara otomatis setelah transaksi berhasil
            await Promise.all([
                fetchActiveReservations(),
                fetchUpcomingReservations(),
                fetchProducts()
            ]);
            return response.data;
        } catch (error: any) {
            errorMessage.value = error.response?.data?.message || 'Gagal menerbitkan kontrak sewa baru';
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function executePickupUnit(id: string) {
        isLoading.value.mutate = true;
        clearError();
        try {
            await rentalApi.pickupReservation(id);
            await fetchActiveReservations();
        } catch (error: any) {
            errorMessage.value = error.response?.data?.message || 'Gagal mengubah status unit sewa';
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function executeMarkReady(id: string) {
        isLoading.value.mutate = true;
        clearError();
        try {
            await rentalApi.markReady(id);
            await fetchActiveReservations();
        } catch (error: any) {
            errorMessage.value = error.response?.data?.message || 'Gagal mengubah status unit sewa';
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function executeUndoReady(id: string) {
        isLoading.value.mutate = true;
        clearError();
        try {
            await rentalApi.undoReady(id);
            await fetchActiveReservations();
        } catch (error: any) {
            errorMessage.value = error.response?.data?.message || 'Gagal mengubah status unit sewa';
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function executeRollbackPickup(id: string) {
        isLoading.value.mutate = true;
        clearError();
        try {
            await rentalApi.rollbackPickupReservation(id);
            // Pembaruan optimistis (Optimistic Update)
            const targetActive = activeReservations.value.find((r: RentalReservation) => r.id === id);
            if (targetActive) targetActive.status = 'BOOKED';

            const targetMaster = masterReservations.value.find((r: RentalReservation) => r.id === id);
            if (targetMaster) targetMaster.status = 'BOOKED';
        } catch (error: any) {
            errorMessage.value = 'Gagal membatalkan pengeluaran box';
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function executeProcessReturnKliring(payload: ProcessReturnPayload, photoFile?: File) {
        isLoading.value.mutate = true;
        clearError();
        try {
            const response = await rentalApi.processReturn(payload);
            console.log("Response from processReturn:", response);

            const returnData = response.data as any; // Due to axios unwrapper
            const returnId = returnData?.id;
            console.log("Parsed returnId:", returnId);
            console.log("Photo File attached:", !!photoFile);

            if (photoFile && returnId) {
                console.log("Triggering uploadReturnPhoto...");
                try {
                    const photoRes = await rentalApi.uploadReturnPhoto(returnId, photoFile);
                    console.log("Photo upload response:", photoRes);
                } catch (photoErr) {
                    console.error("Gagal mengunggah foto, namun retur berhasil dicatat di DB:", photoErr);
                    // Kami tidak men-throw error agar proses pembersihan state (UI) tetap berjalan.
                }
            } else {
                console.warn("Skipping upload. photoFile:", !!photoFile, "returnId:", returnId);
            }

            // Pembersihan baris data secara instan dari tabel pemantauan aktif dan penunggakan
            activeReservations.value = activeReservations.value.filter((r: RentalReservation) => r.id !== payload.reservation_id);
            overdueReservations.value = overdueReservations.value.filter((r: RentalReservation) => r.id !== payload.reservation_id);
            await fetchProducts(); // Kembalikan kuota ketersediaan stok produk
            return response.data;
        } catch (error: any) {
            console.error("Error in executeProcessReturnKliring:", error);
            const detailMsg = error.response?.data?.message || error.message || "Unknown error";
            errorMessage.value = `Sistem gagal melakukan kliring: ${detailMsg}`;
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function getReturnByReservationId(resId: string) {
        isLoading.value.detail = true;
        try {
            const res = await rentalApi.getReturnByReservationId(resId);
            return res.data;
        } catch (error: any) {
            console.error("Gagal menarik data pengembalian", error);
            throw error;
        } finally {
            isLoading.value.detail = false;
        }
    }

    async function uploadReturnReceipt(returnId: string, receiptBlob: Blob, filename: string) {
        isLoading.value.mutate = true;
        clearError();
        try {
            const response = await rentalApi.uploadReturnReceipt(returnId, receiptBlob, filename);
            return response;
        } catch (error: any) {
            errorMessage.value = "Gagal mengunggah dokumen nota pengembalian PDF";
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function uploadReservationInvoice(resId: string, invoiceBlob: Blob, filename: string) {
        isLoading.value.mutate = true;
        clearError();
        try {
            const response = await rentalApi.uploadReservationInvoice(resId, invoiceBlob, filename);
            return response;
        } catch (error: any) {
            errorMessage.value = "Gagal mengunggah dokumen invoice reservasi PDF";
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    async function executeAnulirContract(id: string, penaltyFee: number) {
        isLoading.value.mutate = true;
        clearError();
        try {
            await rentalApi.cancelReservation(id, penaltyFee);
            await fetchMasterReservations(true);
            await fetchActiveReservations();
        } catch (error: any) {
            errorMessage.value = 'Gagal memproses pembatalan pesanan';
            throw error;
        } finally {
            isLoading.value.mutate = false;
        }
    }

    return {
        // State
        masterReservations,
        activeReservations,
        upcomingReservations,
        overdueReservations,
        products,
        currentDetail,
        isLoading,
        errorMessage,
        // Getters
        bookedList,
        pickedUpList,
        readyList,
        totalOverduePotentialPenalties,
        // Actions
        fetchMasterReservations,
        fetchActiveReservations,
        fetchUpcomingReservations,
        fetchOverdueReservations,
        fetchProducts,
        loadReservationDetail,
        executeCreateContract,
        executePickupUnit,
        executeMarkReady,
        executeUndoReady,
        executeRollbackPickup,
        executeProcessReturnKliring,
        executeAnulirContract,
        getReturnByReservationId,
        uploadReturnReceipt,
        uploadReservationInvoice,
    };
});