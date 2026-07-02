<!-- apps/admin-dashboard/src/views/DashboardAnalyticsView.vue -->
<template>
  <div class="space-y-8">
    <!-- HEADER UTAMA & PANEL NAVIGASI KONTROL -->
    <div
      class="flex flex-col xl:flex-row xl:items-center justify-between gap-6 border-b border-outline-variant/20 pb-6"
    >
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">
          Dashboard Pemantauan dan analisa Bisnis Rinzi
        </h2>
        <p class="text-body-md text-on-surface-variant mt-1">
          Selamat datang kembali,
          <span class="font-bold text-primary">{{ authStore.user?.full_name || 'OWNER' }}</span
          >. Konsol pemantauan makro real-time multi-database.
        </p>
      </div>
      <div
        class="flex flex-wrap items-center gap-4 bg-surface-container-low p-3 rounded-2xl border border-outline-variant/30 shadow-sm"
      >
        <div class="flex items-center gap-2">
          <span class="material-symbols-outlined text-on-surface-variant text-md"
            >shield_person</span
          >
          <span class="text-label-md font-mono text-on-surface-variant"
            >Hak Akses: {{ authStore.user?.role || 'OWNER' }}</span
          >
        </div>
        <div class="h-6 w-px bg-outline-variant/40 hidden sm:block"></div>

        <!-- Tombol Sinkronisasi Store -->
        <button
          @click="syncAllStoresData(true)"
          :disabled="isManualSyncing || financeStore.loading"
          class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2.5 rounded-xl shadow-sm hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all text-label-md font-bold disabled:opacity-60"
        >
          <span
            class="material-symbols-outlined text-md"
            :class="isManualSyncing || financeStore.loading ? 'animate-spin' : ''"
            >sync</span
          >
          Sinkronisasi Data
        </button>
      </div>
    </div>

    <!-- BANNER FEEDBACK NOTIFIKASI -->
    <div
      v-if="errorMessage"
      class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium flex items-center gap-2"
    >
      <span class="material-symbols-outlined">error</span>
      <span>{{ errorMessage }}</span>
    </div>
    <div
      v-if="successMessage"
      class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl text-sm font-medium flex items-center gap-2"
    >
      <span class="material-symbols-outlined">check_circle</span>
      <span>{{ successMessage }}</span>
    </div>

    <!-- LOADING STATE UTAMA -->
    <div
      v-if="financeStore.loading && !isManualSyncing"
      class="flex flex-col items-center justify-center min-h-[400px] text-primary bg-surface/50 backdrop-blur-sm rounded-2xl border border-outline-variant/20"
    >
      <span class="animate-spin material-symbols-outlined text-6xl mb-4">sync</span>
      <p class="text-title-medium font-bold tracking-widest uppercase animate-pulse">
        Menghubungkan ke Lini Database Pusat...
      </p>
    </div>

    <div v-else class="space-y-8">
      <!-- ===== ROW 1: KPI CARDS DENGAN TOMBOL REDIRECT ===== -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Gross Retail Sales -->
        <div
          @click="navigateTo('/catalog/pos-sales')"
          class="bg-surface border border-outline-variant/30 rounded-2xl p-6 shadow-sm relative overflow-hidden group hover:border-primary/50 transition-all cursor-pointer">
          <div class="flex items-center justify-between mb-4">
            <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest">Penjualan Barang Toko</span>
            <div class="p-2.5 bg-primary/10 text-primary rounded-xl">
              <span class="material-symbols-outlined text-md">storefront</span>
            </div>
          </div>
          <h3 class="text-headline-small font-black text-on-surface">
            Rp {{ formatNumber(grossRetailSales) }}
          </h3>
          <p class="text-[11px] font-bold text-primary mt-3 flex items-center gap-1">
            <span>Buka Riwayat Penjualan POS</span>
            <span class="material-symbols-outlined text-xs">arrow_forward</span>
          </p>
          <div
            class="absolute bottom-0 inset-x-0 h-1 bg-primary scale-x-0 group-hover:scale-x-100 transition-transform origin-left"
          ></div>
        </div>

        <!-- Rental Gross Income -->
        <div
          @click="navigateTo('/rental/reservations')"
          class="bg-surface border border-outline-variant/30 rounded-2xl p-6 shadow-sm relative overflow-hidden group hover:border-secondary/50 transition-all cursor-pointer"
        >
          <div class="flex items-center justify-between mb-4">
            <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest"
              >Rental Gross Income</span
            >
            <div class="p-2.5 bg-secondary/10 text-secondary rounded-xl">
              <span class="material-symbols-outlined text-md">layers</span>
            </div>
          </div>
          <h3 class="text-headline-small font-black text-on-surface">
            Rp {{ formatNumber(rentalGrossIncome) }}
          </h3>
          <p class="text-[11px] font-bold text-secondary mt-3 flex items-center gap-1">
            <span>Lihat Manifes Kontrak Aktif</span>
            <span class="material-symbols-outlined text-xs">arrow_forward</span>
          </p>
          <div
            class="absolute bottom-0 inset-x-0 h-1 bg-secondary scale-x-0 group-hover:scale-x-100 transition-transform origin-left"
          ></div>
        </div>

        <!-- Pendapatan Harian -->
        <div
          @click="navigateTo('/finance/internal-incomes')"
          class="bg-surface border border-outline-variant/30 rounded-2xl p-6 shadow-sm relative overflow-hidden group hover:border-warning/50 transition-all cursor-pointer"
        >
          <div class="flex items-center justify-between mb-4">
            <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest">Pendapatan Hari Ini</span>
            <div class="p-2.5 bg-amber-500/10 text-amber-600 rounded-xl">
              <span class="material-symbols-outlined text-md">payments</span>
            </div>
          </div>
          <h3 class="text-headline-small font-black text-on-surface">
            Rp {{ formatNumber(dailyRevenue) }}
          </h3>
          <p class="text-[11px] font-bold text-amber-600 mt-3 flex items-center gap-1">
            <span>Audit Kliring Harian</span>
            <span class="material-symbols-outlined text-xs">arrow_forward</span>
          </p>
          <div
            class="absolute bottom-0 inset-x-0 h-1 bg-amber-500 scale-x-0 group-hover:scale-x-100 transition-transform origin-left"
          ></div>
        </div>

        <!-- Margin Persentase Laba -->
        <div
          @click="navigateTo('/finance/reports')"
          class="bg-surface border border-outline-variant/30 rounded-2xl p-6 shadow-sm relative overflow-hidden group hover:border-rose-500/50 transition-all cursor-pointer"
        >
          <div class="flex items-center justify-between mb-4">
            <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest"
              >Margin Rasio Untung</span>
            <div class="p-2.5 bg-rose-500/10 text-rose-600 rounded-xl">
              <span class="material-symbols-outlined text-md">trending_up</span>
            </div>
          </div>
          <h3 class="text-headline-small font-black text-rose-600">
            {{ formatPercent(profitMarginPercent) }}%
          </h3>
          <p class="text-[11px] font-bold text-rose-600 mt-3 flex items-center gap-1">
            <span>Buka Laporan Neraca & Rugi Laba</span>
            <span class="material-symbols-outlined text-xs">arrow_forward</span>
          </p>
          <div
            class="absolute bottom-0 inset-x-0 h-1 bg-rose-500 scale-x-0 group-hover:scale-x-100 transition-transform origin-left"
          ></div>
        </div>
      </div>

      <!-- ===== ROW 2: VISUALISASI GRAFIK BAR & DONUT CHART ===== -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Bar Chart Tren Pendapatan Bulanan -->
        <div
          class="bg-surface border border-outline-variant/30 rounded-2xl p-6 shadow-sm lg:col-span-2 flex flex-col"
        >
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-6">
            <h4 class="font-title-medium text-primary font-bold flex items-center gap-2">
              <span class="material-symbols-outlined text-md">equalizer</span> Tren Arus Kas
              Pendapatan Bulanan
            </h4>
          </div>

          <div
            class="h-44 flex items-end justify-between px-2 pt-4 border-b border-outline-variant/30 font-mono overflow-x-auto gap-1">
            <div
              v-for="(bar, idx) in monthlyRevenueTrend"
              :key="idx"
              class="flex flex-col items-center gap-2 flex-1 min-w-[50px] group relative">
              <div
                class="absolute -top-6 text-on-surface-variant text-[9px] font-bold z-10 pointer-events-none whitespace-nowrap">
                {{ formatNumber(bar.amount) }}
              </div>
              <div class="w-full flex justify-center items-end h-32">
                <div
                  :style="{ height: Math.max((bar.amount / maxTrendValue) * 100, 4) + '%' }"
                  class="w-8 bg-gradient-to-t from-primary to-primary/70 rounded-t-sm hover:brightness-110 transition-all"
                ></div>
              </div>
              <span class="text-[10px] font-bold text-on-surface-variant uppercase mt-1">{{
                bar.month
              }}</span>
            </div>
          </div>
        </div>

        <!-- Donut Pie Chart Distribusi Arus Kas -->
        <div
          class="bg-surface border border-outline-variant/30 rounded-2xl p-6 shadow-sm flex flex-col justify-between"
        >
          <h4 class="font-title-medium text-primary font-bold mb-3 flex items-center gap-2">
            <span class="material-symbols-outlined text-md">donut_large</span> Proporsi Sumber Omzet
          </h4>

          <div class="flex justify-center items-center py-2">
            <div class="relative w-28 h-28">
              <svg class="w-full h-full transform -rotate-90" viewBox="0 0 32 32">
                <circle
                  cx="16"
                  cy="16"
                  r="14"
                  fill="transparent"
                  stroke="#f1f5f9"
                  stroke-width="3"
                />
                <circle
                  cx="16"
                  cy="16"
                  r="14"
                  fill="transparent"
                  stroke="#1580d0"
                  stroke-width="3"
                  :stroke-dasharray="`${(grossRetailSales / (totalGrossRevenue || 1)) * 88} 88`"
                />
                <circle
                  cx="16"
                  cy="16"
                  r="14"
                  fill="transparent"
                  stroke="#59e81d"
                  stroke-width="3"
                  :stroke-dasharray="`${(rentalGrossIncome / (totalGrossRevenue || 1)) * 88} 88`"
                  :stroke-dashoffset="`-${(grossRetailSales / (totalGrossRevenue || 1)) * 88}`"
                />
              </svg>
            </div>
          </div>

          <div class="space-y-2 text-xs font-semibold mt-4">
            <div
              class="flex items-center justify-between p-2 bg-surface-container-low rounded-xl border border-outline-variant/10"
            >
              <span class="text-on-surface flex items-center gap-2"
                ><span class="w-2 h-2 bg-primary rounded-full"></span>POS Ritel</span
              >
              <span class="font-mono text-on-surface-variant"
                >Rp {{ formatNumber(grossRetailSales) }}</span
              >
            </div>
            <div
              class="flex items-center justify-between p-2 bg-surface-container-low rounded-xl border border-outline-variant/10"
            >
              <span class="text-on-surface flex items-center gap-2"
                ><span class="w-2 h-2 bg-secondary rounded-full"></span>Sewa Rental</span
              >
              <span class="font-mono text-on-surface-variant"
                >Rp {{ formatNumber(rentalGrossIncome) }}</span
              >
            </div>
          </div>
        </div>
      </div>

      <!-- ===== ROW 3: FAST MOVING GOODS & RISK CONTROL ===== -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Bulanan Barang Terlaris (POS) -->
        <div
          class="bg-surface border border-outline-variant/30 rounded-2xl shadow-sm overflow-hidden flex flex-col"
        >
          <div
            class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex items-center justify-between"
          >
            <h4 class="font-title-medium text-primary font-bold flex items-center gap-2">
              <span class="material-symbols-outlined text-md">star</span> Produk Ritel Terlaris
              Bulan Ini
            </h4>
          </div>
          <div
            class="p-4 divide-y divide-outline-variant/10 overflow-y-auto max-h-52 custom-scrollbar"
          >
            <div
              v-if="!topProducts || topProducts.length === 0"
              class="py-6 text-center text-on-surface-variant text-xs italic"
            >
              Belum ada riwayat transaksi penjualan dari kasir ritel.
            </div>
            <div
              v-for="(prod, idx) in topProducts"
              :key="idx"
              class="py-2.5 flex items-center justify-between text-xs font-semibold"
            >
              <div>
                <span class="font-mono text-outline mr-2">#{{ Number(idx) + 1 }}</span>
                <span class="text-on-surface font-bold">{{ prod.product_name }}</span>
                <p class="text-[10px] text-on-surface-variant font-normal mt-0.5">
                  {{ prod.qty }} Unit Terjual
                </p>
              </div>
              <span class="font-mono text-on-surface">Rp {{ formatNumber(prod.total) }}</span>
            </div>
          </div>
        </div>

        <!-- Manifes Risiko & Otoritas Rental -->
        <div
          class="bg-surface border border-outline-variant/30 rounded-2xl shadow-sm overflow-hidden flex flex-col justify-between"
        >
          <div
            class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex items-center justify-between"
          >
            <h4 class="font-title-medium text-primary font-bold flex items-center gap-2">
              <span class="material-symbols-outlined text-md">gavel</span> Pengawasan Aset & Risiko
              Rental
            </h4>
          </div>

          <div class="grid grid-cols-2 gap-4 p-6 flex-1 items-center">
            <div
              class="bg-surface-container-low border border-outline-variant/10 rounded-xl p-4 text-center cursor-pointer hover:bg-surface-container-high transition-all"
              @click="navigateTo('/rental/reservations')"
            >
              <span class="text-[10px] font-bold text-on-surface-variant uppercase tracking-wide"
                >Sewa Aktif</span
              >
              <p class="text-2xl font-black text-on-surface font-mono mt-1">
                {{ pickedUpList?.length || 0 }}
              </p>
            </div>
            <div
              class="bg-rose-50 border border-rose-100 rounded-xl p-4 text-center cursor-pointer hover:bg-rose-100 transition-all"
              @click="navigateTo('/rental-damages')"
            >
              <span class="text-[10px] font-bold text-rose-600 uppercase tracking-wide"
                >Terlambat Kembali</span
              >
              <p class="text-2xl font-black text-rose-700 font-mono mt-1">
                {{ overdueReservations?.length || 0 }}
              </p>
            </div>
          </div>

          <div class="p-4 bg-surface-container-low border-t border-outline-variant/10">
            <button
              @click="navigateTo('/rental/reversals')"
              class="w-full py-2 bg-surface text-on-surface font-bold text-xs rounded-xl border border-outline-variant/60 hover:bg-surface-container-high transition-all flex items-center justify-center gap-2"
            >
              <span class="material-symbols-outlined text-sm">restart_alt</span>
              Akses Modul Reversal Kontrak (Rollback)
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@frontend/stores/auth'
import { useFinanceStore } from '@frontend/stores/finance'
import { useRentalStore } from '@frontend/stores/rental'
import { financeApi } from '@frontend/api/finance'

const router = useRouter()

// --- INITIALIZE PINIA SYSTEM STORES ---
const authStore = useAuthStore()
const financeStore = useFinanceStore()
const rentalStore = useRentalStore()

// Membungkus properti state store menggunakan storeToRefs demi menjaga asas reaktivitas data Vue 3
const { dashboardData } = storeToRefs(financeStore)
const { pickedUpList, overdueReservations } = storeToRefs(rentalStore)

// --- LOCAL ANCILLARY STATE ---
const isManualSyncing = ref(false)
const successMessage = ref('')
const errorMessage = ref('')

// --- REACTIVE METRIC CALCULATION ---
const grossRetailSales = computed(() => Number((dashboardData.value as any)?.gross_retail_sales || 0))
const rentalGrossIncome = computed(() => Number((dashboardData.value as any)?.rental_gross_income || 0))
const dailyRevenue = computed(() => Number((dashboardData.value as any)?.daily_revenue || 0))
const monthlyRevenue = computed(() => Number((dashboardData.value as any)?.monthly_revenue || 0))
const totalHpp = computed(() => Number((dashboardData.value as any)?.total_cost_of_goods_sold || 0))
const totalExpenses = computed(() => Number((dashboardData.value as any)?.total_expenses || 0))
const totalGrossRevenue = computed(() => grossRetailSales.value + rentalGrossIncome.value)

const grossProfit = computed(() => {
  const apiGross = (dashboardData.value as any)?.gross_profit
  if (apiGross !== undefined && apiGross !== null) return Number(apiGross)
  return totalGrossRevenue.value - totalHpp.value
})

const netIncome = computed(() => {
  const apiNet = (dashboardData.value as any)?.net_income
  if (apiNet !== undefined && apiNet !== null) return Number(apiNet)
  return grossProfit.value - totalExpenses.value
})

const topProducts = computed(() => (dashboardData.value as any)?.top_products || [])

const profitMarginPercent = computed(() => {
  if (totalGrossRevenue.value === 0) return 0
  return (netIncome.value / totalGrossRevenue.value) * 100
})

const monthlyRevenueTrend = computed(() => {
  const trend = (dashboardData.value as any)?.sales_trend
  if (Array.isArray(trend) && trend.length === 12) {
    return trend
  }
  // Fallback default jika data belum termuat
  return [
    { month: 'Jan', amount: 0 },
    { month: 'Feb', amount: 0 },
    { month: 'Mar', amount: 0 },
    { month: 'Apr', amount: 0 },
    { month: 'Mei', amount: 0 },
    { month: 'Jun', amount: 0 },
    { month: 'Jul', amount: 0 },
    { month: 'Ags', amount: 0 },
    { month: 'Sep', amount: 0 },
    { month: 'Okt', amount: 0 },
    { month: 'Nov', amount: 0 },
    { month: 'Des', amount: 0 },
  ]
})

const maxTrendValue = computed(() => {
  return Math.max(...monthlyRevenueTrend.value.map((t) => t.amount), 1)
})

// --- CORE DISPATCHER DATA RECONCILIATION ---
const syncAllStoresData = async (isManual = false) => {
  if (isManual) isManualSyncing.value = true
  errorMessage.value = ''
  if (isManual) successMessage.value = ''

  try {
    // Menjalankan instruksi muat ulang state secara paralel pada store Pinia terdaftar
    const [_, __, ___, dailyClosings] = await Promise.all([
      financeStore.fetchDashboardSummary(),
      rentalStore.fetchActiveReservations(),
      rentalStore.fetchOverdueReservations(),
      financeApi.getDailyClosings()
    ])

    // Pendapatan hari ini kini langsung diambil secara real-time dari metric dashboard_data (backend)

    if (isManual) {
      successMessage.value =
        'Sinkronisasi berhasil! Arus data antar-database telah terekonsiliasi sempurna.'
    }
  } catch (err: any) {
    errorMessage.value = err?.message || 'Gagal memicu sinkronisasi data store.'
  } finally {
    isManualSyncing.value = false
  }
}

// --- NAVIGATION REDIRECTS ---
const navigateTo = (path: string) => {
  router.push(path)
}

// --- FORMATTERS UTILITIES ---
const formatNumber = (value: number | undefined | null): string => {
  return new Intl.NumberFormat('id-ID', { minimumFractionDigits: 0 }).format(value ?? 0)
}

const formatPercent = (value: number): string => {
  return value.toFixed(1)
}

onMounted(() => {
  syncAllStoresData(false)
})
</script>

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
