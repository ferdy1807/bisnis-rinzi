<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { posApi } from '@frontend/api/pos';
import type { Sale } from '@frontend/types/pos';

// --- STATE MANAGEMENT ---
const salesHistory = ref<Sale[]>([]);
const topProducts = ref<{ product_name: string, total_qty: number, total_revenue: number }[]>([]);
const isLoading = ref(false);
const searchQuery = ref('');
const errorMessage = ref('');

// --- UTILITIES / FORMATTERS ---
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value);
};

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

// --- FETCH DATA OPERATIONS ---
const fetchData = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    const [historyRes, topRes] = await Promise.all([
      posApi.getSalesHistory(),
      posApi.getTopProducts()
    ]);
    salesHistory.value = historyRes;
    topProducts.value = topRes;
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat rekam jejak transaksi penjualan POS.';
  } finally {
    isLoading.value = false;
  }
};

// --- FILTER & SEARCH ---
const filteredSales = computed(() => {
  if (!searchQuery.value) return salesHistory.value;
  const query = searchQuery.value.toLowerCase();
  return salesHistory.value.filter(sale => 
    sale.id?.toLowerCase().includes(query) ||
    sale.invoice_number?.toLowerCase().includes(query) ||
    sale.payment_method?.toLowerCase().includes(query)
  );
});

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-receipt text-indigo-600"></i>
        Riwayat Transaksi POS Ritel
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Monitor manifestasi nota penjualan kasir secara real-time untuk validasi arus kas operasional harian toko.
      </p>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm p-5 space-y-4">
      <div class="flex items-center gap-2 text-slate-800 font-bold text-sm border-b border-slate-100 pb-2">
        <i class="pi pi-star-fill text-amber-500"></i>
        <span>Metrik Performa 3 Produk Terlaris (Top Moving Goods)</span>
      </div>
      
      <div v-if="topProducts.length === 0" class="text-xs text-slate-400 font-normal">
        Belum ada visualisasi data penjualan produk terlaris yang terakumulasi.
      </div>
      
      <div v-else class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div 
          v-for="(prod, idx) in topProducts.slice(0, 3)" 
          :key="idx" 
          class="bg-slate-50 border border-slate-100 rounded-xl p-4 flex flex-col justify-between"
        >
          <div>
            <span class="text-[10px] font-extrabold uppercase tracking-wider px-2 py-0.5 bg-indigo-50 text-indigo-700 rounded-md">
              Peringkat #{{ idx + 1 }}
            </span>
            <p class="text-xs font-bold text-slate-900 mt-2 truncate" :title="prod.product_name">
              {{ prod.product_name }}
            </p>
          </div>
          <div class="mt-3 flex items-baseline justify-between">
            <span class="text-[11px] text-slate-400 font-medium">{{ prod.total_qty }} Unit Terjual</span>
            <span class="text-xs font-mono font-bold text-slate-800">{{ formatCurrency(prod.total_revenue) }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white p-4 rounded-xl border border-slate-200 flex flex-col sm:flex-row items-center gap-3 shadow-sm justify-between">
      <div class="relative w-full sm:w-80">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <i class="pi pi-search text-sm"></i>
        </span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari ID transaksi / metode bayar..."
          class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
        />
      </div>

      <button 
        @click="fetchData"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border transition-colors"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Segarkan Data Pen penjualan
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Sinkronisasi jurnal penjualan toko...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">Waktu Input Nota</th>
              <th class="py-4 px-6">ID Transaksi Pen penjualan</th>
              <th class="py-4 px-6">Metode Pembayaran</th>
              <th class="py-4 px-6 text-center">Diskon Transaksi</th>
              <th class="py-4 px-6 text-right">Total Net Penerimaan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="filteredSales.length === 0 && !isLoading">
              <td colspan="5" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada log transaksi nota kasir yang terdata.
              </td>
            </tr>
            <tr v-for="sale in filteredSales" :key="sale.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6 text-slate-500 font-normal text-xs whitespace-nowrap">
                {{ formatDate(sale.transaction_date) }}
              </td>
              <td class="py-4 px-6 font-mono text-xs text-slate-900 font-bold">
                {{ sale.invoice_number || 'INV-UNKNOWN' }}
              </td>
              <td class="py-4 px-6">
                <span class="px-2 py-0.5 border border-slate-200 text-slate-700 bg-slate-50 rounded text-xs font-bold tracking-wide uppercase">
                  {{ sale.payment_method || 'CASH' }}
                </span>
              </td>
              <td class="py-4 px-6 text-center font-mono text-xs text-slate-500 font-bold">
                {{ sale.discount ? `${sale.discount}%` : '-' }}
              </td>
              <td class="py-4 px-6 text-right font-mono font-extrabold text-slate-900">
                {{ formatCurrency(sale.total || 0) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>