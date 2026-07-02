<!-- apps/admin-dashboard/src/views/FinancialClosingAuditView.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { financeApi } from '@frontend/api/finance';
import { cashApi } from '@frontend/api/cash';

// --- STATE MANAGEMENT ---
const dailyClosings = ref<any[]>([]);
const shifts = ref<any[]>([]);
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
    day: 'numeric'
  });
};

// --- FETCH CLOSING & SHIFT AUDIT DATA ---
const fetchAuditData = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    // Menarik data penutupan harian makro keuangan dan shift mikro kasir secara paralel
    const [closingsRes, shiftsRes] = await Promise.all([
      financeApi.getDailyClosings(),
      cashApi.getShifts()
    ]);
    dailyClosings.value = closingsRes;
    shifts.value = shiftsRes;
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat berkas kliring audit penutupan harian.';
  } finally {
    isLoading.value = false;
  }
};

// --- FILTER & SEARCH ---
const filteredClosings = computed(() => {
  if (!searchQuery.value) return dailyClosings.value;
  const query = searchQuery.value.toLowerCase();
  return dailyClosings.value.filter(closing => 
    closing.id?.toLowerCase().includes(query) ||
    closing.closed_by?.toLowerCase().includes(query) ||
    closing.date?.toLowerCase().includes(query)
  );
});

onMounted(() => {
  fetchAuditData();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <!-- Header Modul -->
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-account-balance text-indigo-600"></i>
        Kliring Tutup Shift & Penutupan Harian
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Konsol verifikasi silang akuntansi Owner. Rekonsiliasi pendapatan tercatat sistem terhadap hitungan fisik uang kasir lapor lunas harian.
      </p>
    </div>

    <!-- Alert Notifikasi Eror -->
    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>

    <!-- STRUKTUR ANALISIS GRID DUA KOLOM -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
      
      <!-- TABEL KIRI: LOG PENUTUPAN MACRO HARIAN (2/3 KOLOM) -->
      <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden lg:col-span-2 flex flex-col relative">
        <div class="p-4 bg-slate-50 border-b border-slate-200 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-500">Jurnal Buku Penutupan Harian (Finance)</span>
          <!-- Bilah Pencarian Mikro -->
          <div class="relative w-full sm:w-64">
            <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400 pointer-events-none">
              <i class="pi pi-search text-xs"></i>
            </span>
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Cari tanggal atau petugas..."
              class="w-full pl-8 pr-3 py-1.5 border border-slate-200 rounded-lg text-xs bg-white outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all"
            />
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50/50 text-slate-400 text-[10px] font-bold tracking-wider uppercase border-b border-slate-200">
                <th class="py-3 px-4">Tanggal Buku</th>
                <th class="py-3 px-4 text-right">Sistem Tercatat</th>
                <th class="py-3 px-4 text-right">Fisik Dilaporkan</th>
                <th class="py-3 px-4 text-right">Selisih Kas (Variance)</th>
                <th class="py-3 px-4 text-center">Status Audit</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-xs font-medium text-slate-700">
              <tr v-if="isLoading && filteredClosings.length === 0">
                <td colspan="5" class="py-8 text-center text-slate-400 font-normal">
                  <i class="pi pi-spin pi-spinner text-indigo-500 mr-1"></i> Membuka brankas data kliring...
                </td>
              </tr>
              <tr v-else-if="filteredClosings.length === 0">
                <td colspan="5" class="py-8 text-center text-slate-400 font-normal">Tidak ada catatan penutupan buku harian.</td>
              </tr>
              <tr v-for="closing in filteredClosings" :key="closing.id" class="hover:bg-slate-50/60 transition-colors">
                <td class="py-4 px-4">
                  <div class="font-bold text-slate-900">{{ formatDate(closing.closing_date) }}</div>
                  <div class="text-[10px] text-slate-400 font-mono">ID: {{ closing.id }}</div>
                </td>
                <td class="py-4 px-4 text-right font-mono font-semibold">{{ formatCurrency((closing.total_sales_retail || 0) + (closing.opening_cash || 0) + (closing.total_other_income || 0) - (closing.total_expenses || 0)) }}</td>
                <td class="py-4 px-4 text-right font-mono font-bold text-slate-900">{{ formatCurrency(closing.actual_cash || 0) }}</td>
                <td 
                  class="py-4 px-4 text-right font-mono font-black"
                  :class="((closing.actual_cash || 0) - ((closing.total_sales_retail || 0) + (closing.opening_cash || 0) + (closing.total_other_income || 0) - (closing.total_expenses || 0))) === 0 ? 'text-slate-500' : ((closing.actual_cash || 0) - ((closing.total_sales_retail || 0) + (closing.opening_cash || 0) + (closing.total_other_income || 0) - (closing.total_expenses || 0))) > 0 ? 'text-emerald-600' : 'text-rose-600'"
                >
                  {{ ((closing.actual_cash || 0) - ((closing.total_sales_retail || 0) + (closing.opening_cash || 0) + (closing.total_other_income || 0) - (closing.total_expenses || 0))) > 0 ? '+' : '' }}{{ formatCurrency((closing.actual_cash || 0) - ((closing.total_sales_retail || 0) + (closing.opening_cash || 0) + (closing.total_other_income || 0) - (closing.total_expenses || 0))) }}
                </td>
                <td class="py-4 px-4 text-center">
                  <span 
                    class="px-2 py-0.5 text-[9px] font-extrabold rounded-md border uppercase"
                    :class="((closing.actual_cash || 0) - ((closing.total_sales_retail || 0) + (closing.opening_cash || 0) + (closing.total_other_income || 0) - (closing.total_expenses || 0))) === 0 ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-amber-50 text-amber-700 border-amber-200'"
                  >
                    {{ ((closing.actual_cash || 0) - ((closing.total_sales_retail || 0) + (closing.opening_cash || 0) + (closing.total_other_income || 0) - (closing.total_expenses || 0))) === 0 ? 'MATCHED' : 'DISCREPANCY' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- PANEL KANAN: MONITOR MONITOR SHIFT KASIR AKTIF (1/3 KOLOM) -->
      <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden flex flex-col relative">
        <div class="p-4 bg-slate-50 border-b border-slate-200 flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-slate-500">Log Aktivitas Shift Lapangan (Cashier)</span>
          <button @click="fetchAuditData" class="p-1 text-slate-400 hover:text-slate-600 transition-colors"><i class="pi pi-refresh text-xs"></i></button>
        </div>

        <div class="divide-y divide-slate-100 overflow-y-auto max-h-[420px] custom-scrollbar">
          <div v-if="isLoading && shifts.length === 0" class="p-6 text-center text-slate-400 text-xs">Memindai mesin laci kasir...</div>
          <div v-else-if="shifts.length === 0" class="p-6 text-center text-slate-400 text-xs">Tidak ada riwayat shift aktif hari ini.</div>
          
          <div v-for="shift in shifts" :key="shift.id" class="p-4 space-y-2 hover:bg-slate-50/50 transition-all">
            <div class="flex items-center justify-between text-xs">
              <span class="font-bold text-slate-900 flex items-center gap-1.5">
                <i class="pi pi-user text-slate-400"></i> {{ shift.cashier_name }}
              </span>
              <span 
                class="px-2 py-0.5 text-[9px] font-bold rounded-sm uppercase tracking-wide"
                :class="shift.status === 'OPEN' ? 'bg-emerald-100 text-emerald-800' : 'bg-slate-100 text-slate-600'">
                {{ shift.status || 'CLOSED' }}
              </span>
            </div>
            
            <div class="grid grid-cols-2 gap-2 text-[11px] font-medium text-slate-500 pt-1 font-mono">
              <div>Buka: <span class="text-slate-800 font-bold">{{ formatCurrency(shift.opening_cash || 0)}}</span></div>
              <div class="text-right">Setor Fisik: <span class="text-indigo-600 font-bold">{{ shift.actual_cash ? formatCurrency(shift.actual_cash) : '-' }}</span></div>
            </div>
            
            <div class="text-[10px] text-slate-400 font-normal border-t border-dashed pt-1.5 flex justify-between items-center">
              <span>Tambahan Masukan: <span class="text-slate-800 font-bold">{{ shift.total_manual_income ? formatCurrency(shift.total_manual_income) : '-' }}</span></span>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}
</style>