<template>
  <div class="space-y-8 animate-fade-in pb-12">
    <div class="bg-surface p-6 sm:p-8 rounded-3xl border border-outline-variant/30 shadow-sm flex flex-col lg:flex-row lg:items-center justify-between gap-6 relative overflow-hidden">
      <div class="absolute -right-10 -bottom-10 w-60 h-60 bg-primary/5 rounded-full blur-2xl pointer-events-none"></div>
      
      <div class="flex items-center gap-4 z-10">
        <div class="w-14 h-14 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-2xl shadow-inner">
          👋
        </div>
        <div>
          <h2 class="text-headline-sm font-bold text-on-surface tracking-tight">
            {{ timeBasedGreeting }}, <span class="text-primary">{{ authStore.user?.full_name || 'Rekan Kasir' }}</span>
          </h2>
          <p class="text-body-sm text-on-surface-variant mt-0.5 font-sans">
            Berikut adalah ringkasan lalu lintas laci kasir ritel Anda hari ini.
          </p>
        </div>
      </div>

      <div class="flex flex-wrap items-center gap-3 bg-surface-container-lowest p-2 px-4 rounded-2xl border border-outline-variant/40 shadow-inner z-10">
        <div class="flex items-center gap-2 border-r border-outline-variant/40 pr-4 py-1">
          <span class="material-symbols-outlined text-primary text-lg">calendar_clock</span>
          <span class="text-label-sm font-mono font-bold text-on-surface">{{ currentDateTime }}</span>
        </div>
        <div class="flex items-center gap-2 pl-1 py-1">
          <span class="relative flex h-2.5 w-2.5">
            <span :class="isShiftOpen ? 'bg-success animate-ping' : 'bg-error'" class="absolute inline-flex h-full w-full rounded-full opacity-75"></span>
            <span :class="isShiftOpen ? 'bg-success' : 'bg-error'" class="relative inline-flex rounded-full h-2.5 w-2.5"></span>
          </span>
          <span class="text-label-sm font-bold tracking-wide uppercase" :class="isShiftOpen ? 'text-success' : 'text-error'">
            {{ isShiftOpen ? `Shift Aktif (${shiftDuration})` : 'Shift Ditutup' }}
          </span>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      
      <div class="lg:col-span-8 space-y-6">
        
        <div
          v-if="!isLoading && !isShiftOpen"
          class="bg-gradient-to-r from-error/20 via-error/10 to-transparent border-2 border-error/40 rounded-3xl p-6 sm:p-8 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-6 shadow-md relative overflow-hidden"
        >
          <div class="absolute -right-4 -top-4 text-error/10 font-black text-8xl select-none pointer-events-none">!</div>
          <div class="flex items-start gap-4 z-10">
            <div class="p-3 bg-error text-on-error rounded-2xl shadow-lg animate-bounce">
              <span class="material-symbols-outlined text-3xl block">lock_open</span>
            </div>
            <div>
              <h3 class="text-title-lg font-black text-error">Laci Kasir Masih Terkunci!</h3>
              <p class="text-body-sm text-on-surface-variant mt-1 max-w-md leading-relaxed font-sans">
                Sistem menolak otorisasi transaksi POS. Silakan masukkan uang modal fisik awal Anda di laci untuk memulai shift.
              </p>
            </div>
          </div>
          <router-link
            to="/shifts"
            class="w-full sm:w-auto px-6 py-3.5 bg-error text-on-error rounded-xl font-label-md font-bold shadow-lg hover:bg-error-container hover:text-on-error-container active:scale-95 transition-all text-center whitespace-nowrap z-10"
          >
            Buka Shift Kasir
          </router-link>
        </div>

        <div v-if="isShiftOpen" class="grid grid-cols-2 sm:grid-cols-3 gap-4">
          <div class="bg-surface p-4 rounded-2xl border border-outline-variant/30 shadow-sm flex items-center gap-3">
            <div class="p-2.5 bg-primary/10 text-primary rounded-xl"><span class="material-symbols-outlined text-md">receipt_long</span></div>
            <div>
              <span class="text-[10px] font-bold text-on-surface-variant uppercase block">Struk Terbit</span>
              <span class="text-title-md font-black text-primary font-mono">{{ todayTxCount }} Transaksi</span>
            </div>
          </div>
          <div class="bg-surface p-4 rounded-2xl border border-outline-variant/30 shadow-sm flex items-center gap-3">
            <div class="p-2.5 bg-success/10 text-success rounded-xl"><span class="material-symbols-outlined text-md">trending_up</span></div>
            <div>
              <span class="text-[10px] font-bold text-on-surface-variant uppercase block">Rata-rata Transaksi</span>
              <span class="text-title-md font-black text-success font-mono">Rp {{ formatNumber(averageBasketSize) }}</span>
            </div>
          </div>
          <div class="bg-surface p-4 rounded-2xl border border-outline-variant/30 shadow-sm col-span-2 sm:col-span-1 flex items-center gap-3">
            <div class="p-2.5 bg-tertiary/10 text-tertiary rounded-xl"><span class="material-symbols-outlined text-md">point_of_sale</span></div>
            <div>
              <span class="text-[10px] font-bold text-on-surface-variant uppercase block">Terminal Sesi</span>
              <span class="text-xs font-bold text-on-surface font-mono">{{ currentSession?.cashier_id}}</span>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <router-link
            to="/pos"
            :class="isShiftOpen ? 'bg-gradient-to-br from-primary to-primary-container text-on-primary shadow-primary/20 hover:shadow-lg' : 'bg-surface-variant/40 text-on-surface-variant opacity-60 cursor-not-allowed'"
            class="p-6 rounded-3xl shadow-md transition-all duration-300 hover:-translate-y-1 group relative overflow-hidden flex flex-col justify-between min-h-[160px] border border-white/10"
            @click="checkAccess($event, '/pos')"
          >
            <div class="absolute -right-8 -bottom-8 w-32 h-32 bg-white/10 rounded-full group-hover:scale-150 transition-transform duration-700"></div>
            <div class="flex justify-between items-start z-10">
              <span class="material-symbols-outlined text-4xl p-2 bg-white/10 rounded-2xl backdrop-blur-md">point_of_sale</span>
              <span class="material-symbols-outlined text-xl opacity-60 group-hover:opacity-100 group-hover:translate-x-1 transition-all">arrow_forward</span>
            </div>
            <div class="z-10">
              <h3 class="font-black text-title-lg tracking-wide">Layar Kasir (POS)</h3>
              <p class="text-xs opacity-80 mt-0.5 font-sans">Scan barcode & proses pembayaran klien</p>
            </div>
          </router-link>

          <router-link
            to="/cash-in"
            :class="isShiftOpen ? 'hover:border-success hover:shadow-success/10' : 'opacity-60 cursor-not-allowed'"
            class="bg-surface p-6 rounded-3xl border border-outline-variant/30 shadow-sm transition-all duration-300 hover:-translate-y-1 group flex flex-col justify-between min-h-[160px]"
            @click="checkAccess($event, '/cash-in')"
          >
            <div class="flex justify-between items-start text-success">
              <span class="material-symbols-outlined text-4xl p-2 bg-success/10 rounded-2xl">payments</span>
              <span class="material-symbols-outlined text-xl opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all">arrow_forward</span>
            </div>
            <div>
              <h3 class="font-black text-title-lg text-on-surface">Input Kas Masuk (+)</h3>
              <p class="text-xs text-on-surface-variant mt-0.5 font-sans">Catat uang kembalian / tambahan laci</p>
            </div>
          </router-link>

          <router-link
            to="/cash-out"
            :class="isShiftOpen ? 'hover:border-error hover:shadow-error/10' : 'opacity-60 cursor-not-allowed'"
            class="bg-surface p-6 rounded-3xl border border-outline-variant/30 shadow-sm transition-all duration-300 hover:-translate-y-1 group flex flex-col justify-between min-h-[160px]"
            @click="checkAccess($event, '/cash-out')"
          >
            <div class="flex justify-between items-start text-error">
              <span class="material-symbols-outlined text-4xl p-2 bg-error/10 rounded-2xl">receipt_long</span>
              <span class="material-symbols-outlined text-xl opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all">arrow_forward</span>
            </div>
            <div>
              <h3 class="font-black text-title-lg text-on-surface">Input Kas Keluar (-)</h3>
              <p class="text-xs text-on-surface-variant mt-0.5 font-sans">Catat bon beli galon, parkir, atau lakban</p>
            </div>
          </router-link>

          <router-link
            to="/daily-report"
            class="bg-surface p-6 rounded-3xl border border-outline-variant/30 shadow-sm hover:border-secondary hover:shadow-secondary/10 transition-all duration-300 hover:-translate-y-1 group flex flex-col justify-between min-h-[160px]"
          >
            <div class="flex justify-between items-start text-secondary">
              <span class="material-symbols-outlined text-4xl p-2 bg-secondary/10 rounded-2xl">assessment</span>
              <span class="material-symbols-outlined text-xl opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all">arrow_forward</span>
            </div>
            <div>
              <h3 class="font-black text-title-lg text-on-surface">Rekapitulasi Harian</h3>
              <p class="text-xs text-on-surface-variant mt-0.5 font-sans">Cetak lembar rekap omzet fisik & POS</p>
            </div>
          </router-link>
        </div>
      </div>

      <div class="lg:col-span-4 space-y-6">
        <div class="bg-surface p-6 sm:p-7 rounded-3xl border border-outline-variant/30 shadow-sm flex flex-col justify-between h-full relative">
          
          <div>
            <div class="flex items-center justify-between pb-4 mb-6 border-b border-outline-variant/20">
              <h3 class="font-bold text-title-md text-on-surface flex items-center gap-2">
                <span class="material-symbols-outlined text-primary text-xl">account_balance_wallet</span>
                Neraca Laci Shift
              </h3>
              <button
                v-if="isShiftOpen"
                @click="refreshDashboardData"
                :disabled="isRefreshing"
                class="p-2 bg-primary/10 hover:bg-primary/20 text-primary rounded-xl transition-colors disabled:opacity-50"
                title="Sinkronisasi ulang laci"
              >
                <span :class="isRefreshing ? 'animate-spin' : ''" class="material-symbols-outlined text-sm block">sync</span>
              </button>
            </div>

            <div v-if="isLoading" class="py-24 flex flex-col items-center justify-center text-primary">
              <span class="animate-spin material-symbols-outlined text-4xl mb-2">sync</span>
              <span class="text-xs font-mono font-bold">Membaca Mutasi...</span>
            </div>

            <div v-else-if="!isShiftOpen" class="py-20 flex flex-col items-center justify-center text-center text-on-surface-variant">
              <div class="w-16 h-16 rounded-full bg-surface-variant/30 flex items-center justify-center mb-3">
                <span class="material-symbols-outlined text-3xl opacity-40">money_off</span>
              </div>
              <p class="text-label-md font-bold">Laci Sedang Dikosongkan</p>
              <p class="text-xs text-outline mt-1 max-w-[200px] font-sans">Buka shift kasir di sebelah kiri untuk melacak uang laci.</p>
            </div>

            <div v-else class="space-y-3.5 font-mono text-xs">
              <div class="flex justify-between items-center p-2.5 rounded-xl bg-surface-container-lowest">
                <span class="text-on-surface-variant font-sans">Uang Modal Awal</span>
                <span class="font-bold text-on-surface">Rp {{ formatNumber(shiftSummary?.opening_cash || 0) }}</span>
              </div>
              <div class="flex justify-between items-center p-2.5 rounded-xl bg-primary/5 text-primary font-bold">
                <span class="font-sans">Penjualan POS (+)</span>
                <span>Rp {{ formatNumber(shiftSummary?.total_income || 0) }}</span>
              </div>
              <div class="flex justify-between items-center p-2.5 rounded-xl bg-success/5 text-success font-bold">
                <span class="font-sans">Kas Masuk Ekstra (+)</span>
                <span>Rp {{ formatNumber(shiftSummary?.total_deposit || 0) }}</span>
              </div>
              <div class="flex justify-between items-center p-2.5 rounded-xl bg-error/5 text-error font-bold">
                <span class="font-sans">Kas Keluar (-)</span>
                <span>- Rp {{ formatNumber(shiftSummary?.total_expense || 0) }}</span>
              </div>

              <div class="border-b-2 border-dashed border-outline-variant/30 my-4"></div>

              <div class="bg-surface-container-low p-4 rounded-2xl border border-outline-variant/40 text-center shadow-inner">
                <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest block font-sans">
                  Target Uang Riil di Laci
                </span>
                <span class="text-2xl font-black text-primary block mt-1">
                  Rp {{ formatNumber(shiftSummary?.expected_cash || 0) }}
                </span>
              </div>
            </div>
          </div>

          <div class="mt-8 pt-4 border-t border-outline-variant/20 space-y-2">
            <router-link
              v-if="isShiftOpen"
              to="/shifts"
              class="w-full py-2.5 bg-error/10 hover:bg-error/20 text-error font-bold rounded-xl text-center block text-xs tracking-wider uppercase transition-colors"
            >
              Tutup Shift Kasir
            </router-link>
            <router-link
              to="/report/shift"
              class="w-full py-2.5 bg-surface-container hover:bg-surface-container-high text-on-surface font-bold rounded-xl text-center block text-xs tracking-wider uppercase transition-colors"
            >
              Riwayat Shift Teraudit
            </router-link>
          </div>

        </div>
      </div>

    </div>

    <router-link
      v-if="isShiftOpen"
      to="/pos"
      class="md:hidden fixed bottom-6 right-6 w-16 h-16 bg-primary text-on-primary rounded-full flex items-center justify-center shadow-2xl active:scale-95 transition-transform z-50 border-2 border-white"
    >
      <span class="material-symbols-outlined text-3xl">point_of_sale</span>
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useAuthStore } from '@frontend/stores/auth';
import { cashApi } from '@frontend/api/cash';
import { posApi } from '@frontend/api/pos';
import type { CashierSession, ShiftSummary } from '@frontend/types/cash';

const authStore = useAuthStore();
const isLoading = ref(true);
const isRefreshing = ref(false);

const currentSession = ref<CashierSession | null>(null);
const shiftSummary = ref<ShiftSummary | null>(null);
const todayTxCount = ref(0);

// Sistem Jam & Pengukur Durasi
const currentDateTime = ref('');
const shiftDuration = ref('');
let timer: number;
let shiftOpenTime: Date | null = null;

// --- SOLUSI BUG LOGIKA: Evaluasi Status Terbuka ---
const isShiftOpen = computed(() => currentSession.value?.status === 'OPEN');

const timeBasedGreeting = computed(() => {
  const hour = new Date().getHours();
  if (hour < 11) return 'Selamat Pagi';
  if (hour < 15) return 'Selamat Siang';
  if (hour < 18) return 'Selamat Sore';
  return 'Selamat Malam';
});

const averageBasketSize = ref(0);

const refreshDashboardData = async () => {
  isRefreshing.value = true;
  try {
    const session = await cashApi.getCurrentSession();
    
    // Tarik rincian laci + hitung struk dari endpoint POS
    const [summary, txsArray] = await Promise.all([
      session && session.status === 'OPEN' ? cashApi.getShiftSummary(session.id) : Promise.resolve(null),
      posApi.getSalesHistory()
    ]);

    const txs = (txsArray as any)?.data || txsArray || [];

    // Rata-rata transaksi per sesi user tersebut
    const userSales = txs.filter((tx: any) => tx.cashier_id === authStore.user?.id);
    const uniqueSessions = new Set(userSales.map((tx: any) => tx.cashier_session_id));
    const totalUserSales = userSales.reduce((sum: number, tx: any) => sum + (tx.total || 0), 0);
    averageBasketSize.value = uniqueSessions.size > 0 ? Math.round(totalUserSales / uniqueSessions.size) : 0;

    // Kunci Validasi: Hanya simpan jika status == 'OPEN'
    if (session && session.status === 'OPEN') {
      currentSession.value = session;
      shiftOpenTime = session.open_time ? new Date(session.open_time) : null;
      shiftSummary.value = summary;

      // Struk terbit = jumlah transaksi di sesi aktif ini
      const currentSessionSales = txs.filter((tx: any) => tx.cashier_session_id === session.id);
      todayTxCount.value = currentSessionSales.length;
    } else {
      currentSession.value = null;
      shiftSummary.value = null;
      shiftOpenTime = null;
      todayTxCount.value = 0;
    }
  } catch (err) {
    console.error('Gagal menyinkronkan data dasbor:', err);
    currentSession.value = null;
  } finally {
    isLoading.value = false;
    isRefreshing.value = false;
  }
};

const updateClocks = () => {
  const now = new Date();
  const options: Intl.DateTimeFormatOptions = {
    day: 'numeric', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  };
  currentDateTime.value = now.toLocaleDateString('id-ID', options).replace(/\./g, ':');

  if (shiftOpenTime && isShiftOpen.value) {
    const diff = Math.max(0, Math.floor((now.getTime() - shiftOpenTime.getTime()) / 1000));
    const h = Math.floor(diff / 3600);
    const m = Math.floor((diff % 3600) / 60);
    const s = diff % 60;
    shiftDuration.value = `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  } else {
    shiftDuration.value = '00:00:00';
  }
};

onMounted(async () => {
  updateClocks();
  timer = window.setInterval(updateClocks, 1000);
  await refreshDashboardData();
});

onUnmounted(() => {
  if (timer) clearInterval(timer);
});

const checkAccess = (event: Event, path: string) => {
  if (!isShiftOpen.value) {
    event.preventDefault();
    alert('Akses ditolak: Anda harus membuka shift terlebih dahulu sebelum mengakses menu ini.');
  }
};

const formatNumber = (amount: number) => {
  return new Intl.NumberFormat('id-ID', { minimumFractionDigits: 0 }).format(amount);
};
</script>