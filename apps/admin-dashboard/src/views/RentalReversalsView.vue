<!-- apps/admin-dashboard/src/views/RentalReversalsView.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { rentalApi } from '@frontend/api/rental';
import type { RentalReservation } from '@frontend/types/rental';

// --- STATE MANAGEMENT ---
const reservations = ref<RentalReservation[]>([]);
const isLoading = ref(false);
const isActionLoading = ref(false);
const searchQuery = ref('');
const errorMessage = ref('');
const successMessage = ref('');

// State Otoritas Alasan Perubahan Status
const reversalReason = ref('');
const activeAction = ref<{ id: string; type: 'CANCEL' | 'ROLLBACK_PICKUP' | 'UNDO_READY' } | null>(null);

// --- FETCH DATA ---
const fetchReservations = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    reservations.value = await rentalApi.getAllReservations();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat manifes kontrak rental.';
  } finally {
    isLoading.value = false;
  }
};

// --- ACTION EXECUTORS ---
const initReversalAction = (id: string, type: 'CANCEL' | 'ROLLBACK_PICKUP' | 'UNDO_READY') => {
  activeAction.value = { id, type };
  reversalReason.value = '';
  errorMessage.value = '';
  successMessage.value = '';
};

const executeReversal = async () => {
  if (!activeAction.value) return;
  if (!reversalReason.value.trim()) {
    errorMessage.value = 'Alasan pembatalan/pemulihan status wajib diisi untuk catatan log audit.';
    return;
  }

  isActionLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    const { id, type } = activeAction.value;

    if (type === 'CANCEL') {
      await rentalApi.cancelReservation(id, 0); // No penalty fee in this view for now
      successMessage.value = 'Kontrak reservasi sewa berhasil dibatalkan secara permanen.';
    } else if (type === 'ROLLBACK_PICKUP') {
      await rentalApi.rollbackPickupReservation(id);
      successMessage.value = 'Status pengambilan barang berhasil dikembalikan (Rollback) ke status Booking.';
    } else if (type === 'UNDO_READY') {
      await rentalApi.undoReady(id);
      successMessage.value = 'Kesiapan unit sewa berhasil dibatalkan (Undo Ready) ke status persiapan.';
    }

    activeAction.value = null;
    await fetchReservations();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal mengeksekusi perintah pemulihan status.';
  } finally {
    isActionLoading.value = false;
  }
};

// --- SEARCH & FILTER LOGIC ---
const filteredReservations = computed(() => {
  if (!searchQuery.value) return reservations.value;
  const query = searchQuery.value.toLowerCase();
  return reservations.value.filter(res => 
    res.customer_name?.toLowerCase().includes(query) ||
    res.invoice_number?.toLowerCase().includes(query) ||
    res.id?.toLowerCase().includes(query)
  );
});

onMounted(() => {
  fetchReservations();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <!-- Header Modul -->
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-published-with-changes text-rose-600"></i>
        Otoritas Pembatalan & Pemulihan (Reversals)
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Menu khusus Owner untuk melakukan intervensi, mengoreksi salah klik status oleh staf lapangan, serta memulihkan data transaksi sewa.
      </p>
    </div>

    <!-- Banner Notifikasi -->
    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl text-sm font-medium">
      {{ successMessage }}
    </div>

    <!-- TOOLBAR PENCARIAN & MANUAL SYNC -->
    <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center gap-3 justify-between">
      <div class="relative w-full sm:w-80">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <i class="pi pi-search text-sm"></i>
        </span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari berdasarkan nama customer atau no invoice..."
          class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
        />
      </div>

      <button 
        @click="fetchReservations"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border transition-colors"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Segarkan Data
      </button>
    </div>

    <!-- TABEL MANIFEST INTERVENSI STATUS -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Menyinkronkan status data rental...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">Nomor Kontrak / Invoice</th>
              <th class="py-4 px-6">Data Penyewa</th>
              <th class="py-4 px-6 text-center">Status Berjalan</th>
              <th class="py-4 px-6 text-right">Aksi Otoritas Eksekutif</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="filteredReservations.length === 0 && !isLoading">
              <td colspan="4" class="py-12 text-center text-slate-400">
                Tidak ada data persewaan aktif yang terdeteksi.
              </td>
            </tr>
            <tr v-for="res in filteredReservations" :key="res.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900">{{ res.invoice_number || 'REG-RENTAL' }}</div>
                <div class="text-[10px] text-slate-400 font-mono select-all">{{ res.id }}</div>
              </td>
              <td class="py-4 px-6">
                <div class="font-semibold text-slate-900">{{ res.customer_name }}</div>
              </td>
              <td class="py-4 px-6 text-center">
                <span class="px-2 py-0.5 text-[10px] font-extrabold bg-slate-100 border text-slate-700 rounded-md uppercase">
                  {{ (res as any).status || 'ACTIVE' }}
                </span>
              </td>
              <td class="py-4 px-6 text-right">
                <div class="flex items-center justify-end gap-2">
                  <!-- Button Undo Ready -->
                  <button 
                    @click="initReversalAction(res.id, 'UNDO_READY')"
                    class="px-2.5 py-1.5 bg-amber-50 hover:bg-amber-100 text-amber-700 border border-amber-200 rounded-lg text-xs font-bold transition-all flex items-center gap-1"
                    title="Batalkan Status Siap / Balikkan ke Persiapan"
                  >
                    <i class="pi pi-undo text-[10px]"></i> Undo Ready
                  </button>

                  <!-- Button Rollback Pickup -->
                  <button 
                    @click="initReversalAction(res.id, 'ROLLBACK_PICKUP')"
                    class="px-2.5 py-1.5 bg-orange-50 hover:bg-orange-100 text-orange-700 border border-orange-200 rounded-lg text-xs font-bold transition-all flex items-center gap-1"
                    title="Batalkan Pengambilan / Balikkan ke Booking"
                  >
                    <i class="pi pi-refresh text-[10px]"></i> Rollback Pickup
                  </button>

                  <!-- Button Cancel Permanen -->
                  <button 
                    @click="initReversalAction(res.id, 'CANCEL')"
                    class="px-2.5 py-1.5 bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200 rounded-lg text-xs font-bold transition-all flex items-center gap-1"
                    title="Batalkan Kontrak Secara Total"
                  >
                    <i class="pi pi-times-circle text-[10px]"></i> Cancel Total
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL CONFIRMATION & LOG AUDIT ENTRIES -->
    <div v-if="activeAction" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-200 overflow-hidden transform transition-all">
        <div class="p-6 border-b border-slate-100 flex items-center justify-between bg-slate-50">
          <h3 class="text-base font-bold text-slate-900 flex items-center gap-2">
            <i class="pi pi-shield text-rose-600"></i>
            Konfirmasi Otoritas Kontrol
          </h3>
          <button @click="activeAction = null" class="text-slate-400 hover:text-slate-600"><i class="pi pi-times"></i></button>
        </div>
        
        <div class="p-6 space-y-4">
          <p class="text-xs text-slate-600 leading-relaxed">
            Anda sedang mengeksekusi perintah intervensi sistem bertipe <span class="px-1.5 py-0.5 font-mono font-bold bg-slate-900 text-white rounded text-[10px]">{{ activeAction.type }}</span> terhadap dokumen ID <span class="font-mono text-slate-900 font-bold bg-slate-100 px-1 rounded">{{ activeAction.id }}</span>.
          </p>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Alasan Intervensi (Masuk Log Audit)</label>
            <textarea 
              v-model="reversalReason"
              rows="3"
              required
              placeholder="Contoh: Salah klik status oleh staf administrasi toko, Unit mendadak trouble..."
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all实时 resize-none"
            ></textarea>
          </div>

          <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-100 mt-6">
            <button type="button" @click="activeAction = null" class="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 rounded-lg border border-slate-200 transition-colors">
              Batal
            </button>
            <button 
              @click="executeReversal"
              :disabled="isActionLoading"
              class="px-4 py-2 text-sm font-bold text-white bg-rose-600 hover:bg-rose-700 rounded-lg shadow-md transition-colors disabled:opacity-50"
            >
              <span>{{ isActionLoading ? 'Memulihkan Status...' : 'Eksekusi Perubahan' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>