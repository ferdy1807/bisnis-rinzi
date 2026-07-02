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

// Pemicu perpindahan tab laporan
const changeTab = (tab: 'profit-loss' | 'balance-sheet' | 'cash-flow') => {
  activeTab.value = tab;
  loadReportData();
};

onMounted(() => {
  loadReportData();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="border-b border-slate-200 pb-5 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
          <i class="pi pi-analytics text-indigo-600"></i>
          Laporan Finansial Standar
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Konsolidasi data akuntansi riil untuk pemantauan kesehatan ekosistem bisnis retail dan rental Anda.
        </p>
      </div>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl flex items-start gap-3">
      <i class="pi pi-exclamation-circle text-lg mt-0.5"></i>
      <span class="text-sm font-medium">{{ errorMessage }}</span>
    </div>

    <div class="flex border-b border-slate-200 space-x-2">
      <button
        @click="changeTab('profit-loss')"
        class="px-5 py-3 text-sm font-semibold border-b-2 transition-all"
        :class="activeTab === 'profit-loss' ? 'border-indigo-600 text-indigo-600 font-bold' : 'border-transparent text-slate-500 hover:text-slate-800'"
      >
        Laporan Laba Rugi
      </button>
      <button
        @click="changeTab('balance-sheet')"
        class="px-5 py-3 text-sm font-semibold border-b-2 transition-all"
        :class="activeTab === 'balance-sheet' ? 'border-indigo-600 text-indigo-600 font-bold' : 'border-transparent text-slate-500 hover:text-slate-800'"
      >
        Neraca Keuangan
      </button>
      <button
        @click="changeTab('cash-flow')"
        class="px-5 py-3 text-sm font-semibold border-b-2 transition-all"
        :class="activeTab === 'cash-flow' ? 'border-indigo-600 text-indigo-600 font-bold' : 'border-transparent text-slate-500 hover:text-slate-800'"
      >
        Arus Kas (Cash Flow)
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm p-6 relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-3xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Mengkonsolidasi Buku Besar Finansial...</span>
      </div>

      <div v-if="activeTab === 'profit-loss' && profitLossData" class="space-y-6">
        <div class="border-b border-slate-100 pb-4">
          <h2 class="text-lg font-bold text-slate-800">Income Statement (Profit & Loss)</h2>
          <p class="text-xs text-slate-400 font-mono mt-0.5">Periode: Akumulasi Berjalan</p>
        </div>

        <div class="space-y-3 max-w-2xl">
          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2">
            <span class="text-slate-600 font-medium pl-2">Pendapatan Bersih (Revenue)</span>
            <span class="font-mono font-bold text-slate-900">{{ formatCurrency((profitLossData as any).total_revenue || 0) }}</span>
          </div>

          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2">
            <span class="text-slate-600 font-medium pl-2">Harga Pokok Penjualan (COGS)</span>
            <span class="font-mono font-bold text-rose-600">({{ formatCurrency((profitLossData as any).total_cogs || 0) }})</span>
          </div>

          <div class="flex justify-between items-center text-sm border-b-2 border-slate-200 py-2 bg-slate-50/50">
            <span class="text-slate-800 font-bold pl-2">Laba Kotor (Gross Profit)</span>
            <span class="font-mono font-bold" :class="(profitLossData as any).gross_profit >= 0 ? 'text-emerald-600' : 'text-rose-600'">
              {{ formatCurrency((profitLossData as any).gross_profit || 0) }}
            </span>
          </div>

          <div class="flex justify-between items-center text-sm border-b border-slate-100 py-2">
            <span class="text-slate-600 font-medium pl-2">Beban Operasional (Operating Expenses)</span>
            <span class="font-mono font-bold text-rose-600">({{ formatCurrency((profitLossData as any).total_expense || 0) }})</span>
          </div>

          <div class="flex justify-between items-center text-base bg-slate-50 p-3 rounded-lg border border-slate-200 mt-4">
            <span class="font-bold text-slate-800">Laba Bersih (Net Income)</span>
            <span class="font-mono font-extrabold text-xl" :class="(profitLossData as any).net_income >= 0 ? 'text-emerald-600' : 'text-rose-600'">
              {{ formatCurrency((profitLossData as any).net_income || 0) }}
            </span>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'balance-sheet' && balanceSheetData" class="space-y-8">
        <div class="border-b border-slate-100 pb-4">
          <h2 class="text-lg font-bold text-slate-800">Balance Sheet</h2>
          <p class="text-xs text-slate-400 font-mono mt-0.5">Persamaan Dasar Akuntansi: Aset = Liabilitas + Ekuitas</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div class="space-y-4">
            <h3 class="text-sm font-bold uppercase tracking-wider text-indigo-600 border-b pb-1">Aset (Assets)</h3>
            <div class="flex justify-between text-sm py-1 font-medium text-slate-700">
              <span>Total Aset</span>
              <span class="font-mono font-bold">{{ formatCurrency((balanceSheetData as any).total_assets || 0) }}</span>
            </div>
          </div>

          <div class="space-y-6">
            <div class="space-y-2">
              <h3 class="text-sm font-bold uppercase tracking-wider text-slate-500 border-b pb-1">Kewajiban & Modal</h3>
              <div class="flex justify-between text-sm py-1 font-medium text-slate-700">
                <span>Total Liabilitas (Kewajiban)</span>
                <span class="font-mono font-bold">{{ formatCurrency((balanceSheetData as any).total_liab || 0) }}</span>
              </div>
              <div class="flex justify-between text-sm py-1 font-medium text-slate-700">
                <span>Total Ekuitas (Modal Saham)</span>
                <span class="font-mono font-bold">{{ formatCurrency((balanceSheetData as any).total_equity || 0) }}</span>
              </div>
            </div>

            <div class="flex justify-between text-sm p-3 bg-slate-50 rounded-lg font-bold text-slate-800 border">
              <span>Total Pasiva (Liabilitas + Ekuitas)</span>
              <span class="font-mono text-indigo-600">{{ formatCurrency(((balanceSheetData as any).total_liab || 0) + ((balanceSheetData as any).total_equity || 0)) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'cash-flow' && cashFlowData" class="space-y-6">
        <div class="border-b border-slate-100 pb-4">
          <h2 class="text-lg font-bold text-slate-800">Statement of Cash Flows</h2>
          <p class="text-xs text-slate-400 font-mono mt-0.5">Laporan Rekonsiliasi Aliran Likuiditas Riil Kas</p>
        </div>

        <div class="space-y-3 max-w-2xl">
          <div class="flex justify-between text-sm border-b pb-2">
            <span class="text-slate-600 font-medium">Aktivitas Operasional</span>
            <span class="font-mono font-bold" :class="(cashFlowData.operating_activities || 0) >= 0 ? 'text-slate-900' : 'text-rose-600'">
              {{ formatCurrency(cashFlowData.operating_activities || 0) }}
            </span>
          </div>
          <div class="flex justify-between text-sm border-b pb-2">
            <span class="text-slate-600 font-medium">Aktivitas Investasi</span>
            <span class="font-mono font-bold" :class="(cashFlowData.investing_activities || 0) >= 0 ? 'text-slate-900' : 'text-rose-600'">
              {{ formatCurrency(cashFlowData.investing_activities || 0) }}
            </span>
          </div>
          <div class="flex justify-between text-sm border-b pb-2">
            <span class="text-slate-600 font-medium">Aktivitas Pendanaan</span>
            <span class="font-mono font-bold" :class="(cashFlowData.financing_activities || 0) >= 0 ? 'text-slate-900' : 'text-rose-600'">
              {{ formatCurrency(cashFlowData.financing_activities || 0) }}
            </span>
          </div>

          <div class="flex justify-between items-center text-base bg-slate-900 text-white p-3 rounded-lg font-bold mt-6">
            <span>Kenaikan/Penurunan Bersih Kas</span>
            <span class="font-mono text-xl text-emerald-400">{{ formatCurrency(cashFlowData.net_cash_flow || 0) }}</span>
          </div>
        </div>
      </div>
      
      <div v-if="!isLoading && !profitLossData && !balanceSheetData && !cashFlowData" class="py-12 text-center text-slate-400">
        Tidak ada data laporan yang tersedia atau belum dikonsolidasi oleh sistem akuntansi pusat.
      </div>
    </div>
  </div>
</template>