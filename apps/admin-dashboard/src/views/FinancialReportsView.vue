<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { financeApi } from '@frontend/api/finance';
import type { IncomeStatementResponse, BalanceSheetResponse } from '@frontend/types/finance';

// --- STATE MANAGEMENT ---
const activeTab = ref<'profit-loss' | 'balance-sheet' | 'cash-flow'>('profit-loss');
const isLoading = ref(false);
const errorMessage = ref('');

// Data Reports State
const profitLossData = ref<IncomeStatementResponse | null>(null);
const balanceSheetData = ref<BalanceSheetResponse | null>(null);
const cashFlowData = ref<any | null>(null);

// --- UTILITIES / FORMATTERS ---
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value);
};

// --- FETCH DATA OPERATIONS ---
const fetchProfitLoss = async () => {
  try {
    profitLossData.value = await financeApi.getProfitLossReport();
  } catch (error: any) {
    throw new Error(`Laba Rugi: ${error?.message || 'Gagal memuat'}`);
  }
};

const fetchBalanceSheet = async () => {
  try {
    balanceSheetData.value = await financeApi.getBalanceSheetReport();
  } catch (error: any) {
    throw new Error(`Neraca: ${error?.message || 'Gagal memuat'}`);
  }
};

const fetchCashFlow = async () => {
  try {
    cashFlowData.value = await financeApi.getCashFlowReport();
  } catch (error: any) {
    throw new Error(`Arus Kas: ${error?.message || 'Gagal memuat'}`);
  }
};

const loadReportData = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    if (activeTab.value === 'profit-loss') await fetchProfitLoss();
    if (activeTab.value === 'balance-sheet') await fetchBalanceSheet();
    if (activeTab.value === 'cash-flow') await fetchCashFlow();
  } catch (error: any) {
    errorMessage.value = error.message;
  } finally {
    isLoading.value = false;
  }
};

const changeTab = (tab: 'profit-loss' | 'balance-sheet' | 'cash-flow') => {
  activeTab.value = tab;
  loadReportData();
};

onMounted(() => {
  loadReportData();
});
</script>

<template>
  <div class="p-6 max-w-[1400px] mx-auto space-y-6 font-sans text-slate-800 antialiased">
    
    <div class="border-b border-slate-200 pb-5 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <p class="text-[11px] font-bold tracking-widest text-indigo-600 uppercase mb-1">Akuntansi & Finansial</p>
        <h1 class="text-2xl font-extrabold text-slate-900 tracking-tight flex items-center gap-2.5">
          <i class="pi pi-analytics text-indigo-600 bg-indigo-50 p-2 rounded-xl text-xl"></i>
          Laporan Finansial Standar
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Konsolidasi data akuntansi riil untuk pemantauan kesehatan ekosistem bisnis retail dan rental Anda.
        </p>
      </div>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl flex items-start gap-3 shadow-3xs">
      <i class="pi pi-exclamation-circle text-lg mt-0.5"></i>
      <span class="text-sm font-medium">{{ errorMessage }}</span>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <div class="bg-white border border-slate-200 p-5 rounded-2xl flex items-center justify-between shadow-3xs">
        <div class="space-y-1">
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Total Pendapatan</span>
          <h3 class="text-xl font-extrabold font-mono text-slate-900 m-0">
            {{ profitLossData ? formatCurrency((profitLossData as any).total_revenue || 0) : 'Rp 0' }}
          </h3>
        </div>
        <span class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center"><i class="pi pi-money-bill text-base"></i></span>
      </div>

      <div class="bg-white border border-slate-200 p-5 rounded-2xl flex items-center justify-between shadow-3xs">
        <div class="space-y-1">
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Laba Bersih Berjalan</span>
          <h3 class="text-xl font-extrabold font-mono m-0" :class="profitLossData && (profitLossData as any).net_income >= 0 ? 'text-emerald-600' : 'text-rose-600'">
            {{ profitLossData ? formatCurrency((profitLossData as any).net_income || 0) : 'Rp 0' }}
          </h3>
        </div>
        <span class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center"><i class="pi pi-percentage text-base"></i></span>
      </div>

      <div class="bg-white border border-slate-200 p-5 rounded-2xl flex items-center justify-between shadow-3xs">
        <div class="space-y-1">
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Total Aset (Neraca)</span>
          <h3 class="text-xl font-extrabold font-mono text-slate-900 m-0">
            {{ balanceSheetData ? formatCurrency((balanceSheetData as any).total_assets || 0) : 'Rp 0' }}
          </h3>
        </div>
        <span class="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center"><i class="pi pi-briefcase text-base"></i></span>
      </div>
    </div>

    <div class="flex bg-slate-100 p-1.5 rounded-2xl border border-slate-200/60 max-w-max">
      <button
        @click="changeTab('profit-loss')"
        :class="activeTab === 'profit-loss' ? 'bg-white text-indigo-700 shadow-xs border-slate-200/50 font-bold' : 'text-slate-500 hover:text-slate-800'"
        class="px-5 py-2.5 text-xs font-bold rounded-xl border border-transparent bg-transparent transition-all cursor-pointer"
      >
        Laporan Laba Rugi
      </button>
      <button
        @click="changeTab('balance-sheet')"
        :class="activeTab === 'balance-sheet' ? 'bg-white text-indigo-700 shadow-xs border-slate-200/50 font-bold' : 'text-slate-500 hover:text-slate-800'"
        class="px-5 py-2.5 text-xs font-bold rounded-xl border border-transparent bg-transparent transition-all cursor-pointer"
      >
        Neraca Keuangan
      </button>
      <button
        @click="changeTab('cash-flow')"
        :class="activeTab === 'cash-flow' ? 'bg-white text-indigo-700 shadow-xs border-slate-200/50 font-bold' : 'text-slate-500 hover:text-slate-800'"
        class="px-5 py-2.5 text-xs font-bold rounded-xl border border-transparent bg-transparent transition-all cursor-pointer"
      >
        Arus Kas (Cash Flow)
      </button>
    </div>

    <div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 relative min-h-[300px]">
      
      <div v-if="isLoading" class="absolute inset-0 bg-white/80 backdrop-blur-xs z-10 flex flex-col items-center justify-center rounded-2xl">
        <i class="pi pi-spin pi-spinner text-3xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-bold text-slate-500">Mengkonsolidasi Buku Besar Finansial...</span>
      </div>

      <div v-if="activeTab === 'profit-loss' && profitLossData" class="space-y-6">
        <div class="border-b border-slate-100 pb-3">
          <h2 class="text-base font-extrabold text-slate-900 m-0">Income Statement (Profit & Loss)</h2>
          <p class="text-xs text-slate-400 font-mono mt-0.5">Periode: Akumulasi Berjalan</p>
        </div>

        <div class="space-y-2 max-w-3xl">
          
          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2.5">
            <span class="text-slate-600 font-medium pl-2 flex items-center gap-2">
              <i class="pi pi-shopping-bag text-slate-400 text-xs"></i> Pendapatan Penjualan Toko
            </span>
            <span class="font-mono text-slate-700 font-semibold">
              {{ formatCurrency((profitLossData as any).retail_revenue || 0) }}
            </span>
          </div>

          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2.5 bg-indigo-50/20 rounded-lg">
            <span class="text-indigo-900 font-bold pl-2 flex items-center gap-2">
              <i class="pi pi-box text-indigo-500 text-xs"></i> Pendapatan Rental (Sewa Box)
            </span>
            <span class="font-mono text-indigo-700 font-bold">
              {{ formatCurrency((profitLossData as any).rental_revenue || 0) }}
            </span>
          </div>

          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2.5 bg-slate-50/50">
            <span class="text-slate-800 font-bold pl-2">Total Pendapatan Bersih (Gross Revenue)</span>
            <span class="font-mono font-bold text-slate-900 text-base">
              {{ formatCurrency((profitLossData as any).total_revenue || 0) }}
            </span>
          </div>

          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2.5">
            <span class="text-slate-600 font-medium pl-2">Harga Pokok Penjualan (COGS)</span>
            <span class="font-mono font-bold text-rose-600">({{ formatCurrency((profitLossData as any).total_cogs || 0) }})</span>
          </div>

          <div class="flex justify-between items-center text-sm border-b-2 border-slate-200 py-3 bg-slate-50/80 rounded-lg px-1">
            <span class="text-slate-900 font-extrabold pl-2">Laba Kotor (Gross Profit)</span>
            <span class="font-mono font-extrabold text-base" :class="(profitLossData as any).gross_profit >= 0 ? 'text-emerald-600' : 'text-rose-600'">
              {{ formatCurrency((profitLossData as any).gross_profit || 0) }}
            </span>
          </div>

          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2.5">
            <span class="text-slate-600 font-medium pl-2">Beban Operasional (Operating Expenses)</span>
            <span class="font-mono font-bold text-rose-600">({{ formatCurrency((profitLossData as any).total_expense || 0) }})</span>
          </div>

          <div class="flex justify-between items-center text-base bg-slate-900 text-white p-4 rounded-xl shadow-xs mt-6">
            <span class="font-bold tracking-wide">Laba Bersih Setelah Pajak (Net Income)</span>
            <span class="font-mono font-black text-xl" :class="(profitLossData as any).net_income >= 0 ? 'text-emerald-400' : 'text-rose-400'">
              {{ formatCurrency((profitLossData as any).net_income || 0) }}
            </span>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'balance-sheet' && balanceSheetData" class="space-y-6">
        <div class="border-b border-slate-100 pb-3">
          <h2 class="text-base font-extrabold text-slate-900 m-0">Balance Sheet</h2>
          <p class="text-xs text-slate-400 font-mono mt-0.5">Persamaan Dasar Akuntansi: Aset = Liabilitas + Ekuitas</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div class="space-y-4">
            <h3 class="text-xs font-bold uppercase tracking-wider text-indigo-600 border-b border-indigo-100 pb-1.5 m-0">Aset (Assets)</h3>
            <div class="flex justify-between text-sm py-2 px-2 font-semibold text-slate-700 bg-slate-50/50 rounded-lg">
              <span>Total Nilai Aset</span>
              <span class="font-mono font-bold text-slate-900">{{ formatCurrency((balanceSheetData as any).total_assets || 0) }}</span>
            </div>
          </div>

          <div class="space-y-5">
            <div class="space-y-2">
              <h3 class="text-xs font-bold uppercase tracking-wider text-slate-500 border-b border-slate-200 pb-1.5 m-0">Kewajiban & Modal (Liabilities & Equity)</h3>
              <div class="flex justify-between text-sm py-2 px-1 font-medium text-slate-700 border-b border-slate-100">
                <span>Total Liabilitas (Kewajiban)</span>
                <span class="font-mono font-bold text-slate-800">{{ formatCurrency((balanceSheetData as any).total_liab || 0) }}</span>
              </div>
              <div class="flex justify-between text-sm py-2 px-1 font-medium text-slate-700 border-b border-slate-100">
                <span>Total Ekuitas (Modal Saham)</span>
                <span class="font-mono font-bold text-slate-800">{{ formatCurrency((balanceSheetData as any).total_equity || 0) }}</span>
              </div>
            </div>

            <div class="flex justify-between text-sm p-3.5 bg-slate-900 text-white rounded-xl font-bold border shadow-xs">
              <span>Total Pasiva (Liabilitas + Ekuitas)</span>
              <span class="font-mono text-amber-400 font-extrabold">{{ formatCurrency(((balanceSheetData as any).total_liab || 0) + ((balanceSheetData as any).total_equity || 0)) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'cash-flow' && cashFlowData" class="space-y-6">
        <div class="border-b border-slate-100 pb-3">
          <h2 class="text-base font-extrabold text-slate-900 m-0">Statement of Cash Flows</h2>
          <p class="text-xs text-slate-400 font-mono mt-0.5">Laporan Rekonsiliasi Aliran Likuiditas Riil Kas</p>
        </div>

        <div class="space-y-2 max-w-3xl">
          <div class="flex justify-between text-sm border-b border-slate-100 py-2.5 pl-1">
            <span class="text-slate-600 font-medium">Aktivitas Operasional (Operational)</span>
            <span class="font-mono font-bold" :class="(cashFlowData.operating_activities || 0) >= 0 ? 'text-slate-900' : 'text-rose-600'">
              {{ formatCurrency(cashFlowData.operating_activities || 0) }}
            </span>
          </div>
          <div class="flex justify-between text-sm border-b border-slate-100 py-2.5 pl-1">
            <span class="text-slate-600 font-medium">Aktivitas Investasi (Investing)</span>
            <span class="font-mono font-bold" :class="(cashFlowData.investing_activities || 0) >= 0 ? 'text-slate-900' : 'text-rose-600'">
              {{ formatCurrency(cashFlowData.investing_activities || 0) }}
            </span>
          </div>
          <div class="flex justify-between text-sm border-b border-slate-100 py-2.5 pl-1">
            <span class="text-slate-600 font-medium">Aktivitas Pendanaan (Financing)</span>
            <span class="font-mono font-bold" :class="(cashFlowData.financing_activities || 0) >= 0 ? 'text-slate-900' : 'text-rose-600'">
              {{ formatCurrency(cashFlowData.financing_activities || 0) }}
            </span>
          </div>

          <div class="flex justify-between items-center text-base bg-slate-900 text-white p-4 rounded-xl font-bold shadow-xs mt-6">
            <span>Kenaikan/Penurunan Bersih Kas</span>
            <span class="font-mono text-xl text-emerald-400 font-extrabold">{{ formatCurrency(cashFlowData.net_cash_flow || 0) }}</span>
          </div>
        </div>
      </div>
      
      <div v-if="!isLoading && !profitLossData && !balanceSheetData && !cashFlowData" class="py-16 text-center text-slate-400 flex flex-col items-center justify-center gap-2">
        <i class="pi pi-folder-open text-3xl text-slate-300"></i>
        <span class="text-sm font-medium">Tidak ada data laporan finansial yang tersedia saat ini.</span>
      </div>
    </div>

  </div>
</template>