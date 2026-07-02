<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { inventoryApi } from '@frontend/api/inventory';
import type { StockMovement } from '@frontend/types/inventory';

// --- STATE MANAGEMENT ---
const movements = ref<StockMovement[]>([]);
const isLoading = ref(false);
const searchQuery = ref('');
const filterType = ref<string>('ALL'); // ALL, IN, OUT, ADJUSTMENT
const errorMessage = ref('');

// --- UTILITIES / FORMATTERS ---
const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// --- FETCH DATA ---
const fetchMovements = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    movements.value = await inventoryApi.getStockMovements();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat riwayat mutasi stok barang.';
  } {
    isLoading.value = false;
  }
};

// --- FILTER & SEARCH LOGIC ---
const filteredMovements = computed(() => {
  let result = movements.value;

  // Filter Berdasarkan Tipe Log Penjualan / Mutasi
  if (filterType.value !== 'ALL') {
    result = result.filter((m: any) => (m as any).movement_type === filterType.value);
  }

  // Filter Berdasarkan Isian Kata Kunci Pencarian
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter((m: any) => 
      (m as any).reference?.toLowerCase().includes(query) ||
      m.id?.toLowerCase().includes(query) ||
      m.product_id?.toLowerCase().includes(query) ||
      m.product_name?.toLowerCase().includes(query) ||
      m.sku?.toLowerCase().includes(query)
    );
  }

  return result;
});

onMounted(() => {
  fetchMovements();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-swap-horiz text-indigo-600"></i>
        Riwayat Mutasi Barang (Stock Movements)
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Jejak audit kronologis alur keluar masuk komoditas barang dagangan untuk mendeteksi penyusutan stok fisik dan manipulasi data.
      </p>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>

    <div class="bg-white p-4 rounded-xl border border-slate-200 flex flex-col md:flex-row items-center gap-4 shadow-sm justify-between">
      <div class="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
        <div class="relative w-full sm:w-72">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
            <i class="pi pi-search text-sm"></i>
          </span>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari produk atau nomor referensi..."
            class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
          />
        </div>

        <div class="relative w-full sm:w-48">
          <select 
            v-model="filterType"
            class="w-full pl-3 pr-10 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 appearance-none font-semibold text-slate-700"
          >
            <option value="ALL">Semua Jenis Mutasi</option>
            <option value="IN">Stok Masuk (IN)</option>
            <option value="OUT">Stok Keluar (OUT)</option>
            <option value="ADJUSTMENT">Penyesuaian (ADJUSTMENT)</option>
          </select>
          <span class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-slate-400">
            <i class="pi pi-chevron-down text-xs"></i>
          </span>
        </div>
      </div>

      <button 
        @click="fetchMovements"
        class="w-full md:w-auto inline-flex items-center justify-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border transition-colors"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Segarkan Log
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Mengekstrak berkas jejak audit barang...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">Waktu Eksekusi</th>
              <th class="py-4 px-6">ID Log</th>
              <th class="py-4 px-6">Deskripsi Komoditas Produk</th>
              <th class="py-4 px-6">Nomor Referensi Dokumen</th>
              <th class="py-4 px-6 text-center">Tipe</th>
              <th class="py-4 px-6 text-right">Kuantitas Perubahan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="filteredMovements.length === 0 && !isLoading">
              <td colspan="6" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada riwayat pergerakan stok yang cocok dengan filter pencarian.
              </td>
            </tr>
            <tr v-for="movement in filteredMovements" :key="movement.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6 text-slate-500 font-normal text-xs whitespace-nowrap">
                {{ formatDate(movement.created_at) }}
              </td>
              <td class="py-4 px-6 font-mono text-xs text-slate-400 select-all">
                {{ movement.id }}
              </td>
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900">{{ movement.product_name || 'Tidak Diketahui' }}</div>
                <div class="text-[10px] text-slate-400 font-mono mt-0.5">SKU: {{ movement.sku || '-' }}</div>
              </td>
              <td class="py-4 px-6 font-mono text-xs text-slate-600">
                {{ (movement as any).reference || 'SISTEM_AUTOMATION' }}
              </td>
              <td class="py-4 px-6 text-center">
                <span 
                  class="px-2.5 py-1 text-[10px] font-extrabold rounded-md uppercase border"
                  :class="{
                    'bg-emerald-50 text-emerald-700 border-emerald-200': (movement as any).movement_type === 'IN',
                    'bg-rose-50 text-rose-700 border-rose-200': (movement as any).movement_type === 'OUT',
                    'bg-amber-50 text-amber-700 border-amber-200': (movement as any).movement_type === 'ADJUSTMENT'
                  }"
                >
                  {{ (movement as any).movement_type || 'UNKNOWN' }}
                </span>
              </td>
              <td 
                class="py-4 px-6 text-right font-mono font-bold text-base"
                :class="(movement.qty || 0) >= 0 ? 'text-emerald-600' : 'text-rose-600'">
                {{ (movement.qty || 0) >= 0 ? `+${movement.qty}` : movement.qty }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>