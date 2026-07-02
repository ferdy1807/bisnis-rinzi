<template>
  <!-- Perubahan Kontainer Utama menjadi Full View / Full Width -->
  <div class="w-full max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-6 font-sans text-slate-900">
      
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-200 pb-1">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-primary/10 rounded-2xl flex items-center justify-center text-primary shrink-0">
            <span class="material-symbols-outlined text-3xl">point_of_sale</span>
          </div>
          <div>
            <h2 class="text-2xl font-black tracking-tight text-slate-900">Manajemen Sesi Operasional Kasir</h2>
            <p class="text-xs text-slate-500 mt-0.5">Lacak akurasi kas laci fisik, kelola rekonsiliasi finansial, dan pengawasan shift.</p>
          </div>
        </div>
      </div>

      <div v-if="isLoading" class="flex flex-col justify-center items-center py-32 text-primary gap-3">
        <span class="animate-spin material-symbols-outlined text-5xl">sync</span>
        <p class="text-xs font-bold text-slate-400 animate-pulse">Sinkronisasi data pembukuan kas aktif...</p>
      </div>

      <div v-else-if="!currentSession" class="bg-white border border-slate-200 shadow-xl rounded-3xl overflow-hidden max-w-xl mx-auto my-12 transition-all">
        <div class="bg-slate-50 p-10 text-center border-b border-slate-100 relative overflow-hidden">
          <div class="absolute -right-6 -top-6 w-24 h-24 bg-primary/5 rounded-full blur-xl pointer-events-none"></div>
          <span class="material-symbols-outlined text-6xl text-primary mb-4 block">lock_open</span>
          <h3 class="text-xl font-black text-slate-900 mb-2">Buka Shift Sesi Baru</h3>
          <p class="text-xs text-slate-500 max-w-sm mx-auto leading-relaxed">Masukkan saldo awal uang tunai (laci kas) sebelum memulai transaksi untuk menjaga integritas audit finansial.</p>
        </div>
        <form @submit.prevent="handleOpenSession" class="p-8 space-y-6">
          <div class="space-y-2">
            <label class="text-xs font-black text-slate-600 uppercase tracking-wider block">Modal Awal Laci (Opening Cash)</label>
            <div class="relative rounded-xl shadow-xs">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 font-bold text-slate-400 text-sm">Rp</span>
              <input 
                v-model="formattedOpeningCash" 
                type="text" 
                required
                placeholder="0"
                class="w-full bg-slate-50 border border-slate-200 rounded-xl pl-12 pr-4 py-3.5 font-mono font-bold text-slate-900 focus:ring-2 focus:ring-primary focus:border-primary focus:bg-white focus:outline-none transition-all"/>
            </div>
          </div>
          <button 
            type="submit" 
            class="w-full bg-primary text-white py-4 rounded-xl font-bold text-sm shadow-md hover:bg-primary/90 hover:-translate-y-0.5 active:translate-y-0 transition-all flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
            :disabled="isSubmitting"
          >
            <span v-if="isSubmitting" class="animate-spin material-symbols-outlined text-sm">sync</span>
            <span v-else class="material-symbols-outlined text-sm">play_arrow</span>
            {{ isSubmitting ? 'Membuka Sesi...' : 'Aktifkan Sesi Shift Sekarang' }}
          </button>
        </form>
      </div>

      <!-- Tampilan Grid Meja Kasir (Diubah menjadi Full View Grid 12 Kolom) -->
      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        
        <!-- Summary Card (Left Side - 7 Kolom) -->
        <div class="lg:col-span-7 bg-white border border-slate-200 shadow-sm rounded-3xl overflow-hidden flex flex-col min-h-[480px]">
          <div class="p-6 bg-slate-50/80 border-b border-slate-100 flex justify-between items-center">
            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <span class="relative flex h-2.5 w-2.5">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
                </span>
                <span class="text-[10px] font-black text-emerald-700 tracking-wider bg-emerald-50 px-2 py-0.5 rounded-md border border-emerald-100">SHIFT OPERASIONAL AKTIF</span>
              </div>
              <h3 class="text-lg font-black text-slate-900">Ringkasan Arus Buku Kas Sesi</h3>
            </div>
            <div class="text-right">
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Waktu Mulai Shift</p>
              <p class="text-xs font-mono font-bold text-slate-700 mt-0.5">{{ formatDate(currentSession.open_time) }}</p>
            </div>
          </div>
        
          <div class="p-6 flex-1 space-y-4 bg-white">
            <div class="flex justify-between items-center pb-3.5 border-b border-slate-100 text-sm">
              <span class="text-slate-500 font-medium">Saldo Awal Laci Pembukaan:</span>
              <span class="font-mono font-bold text-slate-900">Rp {{ summary?.opening_cash?.toLocaleString('id-ID') || 0 }}</span>
            </div>
            <div class="flex justify-between items-center pb-3.5 border-b border-slate-100 text-sm text-blue-700">
              <span class="font-medium">Total Penerimaan POS (+):</span>
              <span class="font-mono font-black">Rp {{ summary?.total_income?.toLocaleString('id-ID') || 0 }}</span>
            </div>
            <div class="flex justify-between items-center pb-3.5 border-b border-slate-100 text-sm text-emerald-600">
              <span class="font-medium">Kas Masuk Logistik Ekstra (+):</span>
              <span class="font-mono font-black">Rp {{ ((summary?.total_deposit || 0) - (summary?.total_income || 0)).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between items-center pb-1 text-sm text-rose-600">
              <span class="font-medium">Total Pengeluaran Kas Operasional (-):</span>
              <span class="font-mono font-black">Rp {{ summary?.total_expense?.toLocaleString('id-ID') || 0 }}</span>
            </div>
          </div>

          <div class="p-6 bg-slate-50 border-t border-slate-100 text-center space-y-1">
            <p class="text-[11px] font-black text-slate-400 uppercase tracking-widest">Ekspektasi Uang Fisik Sistem di Laci</p>
            <p class="text-3xl font-black font-mono text-primary tracking-tight">Rp {{ summary?.expected_cash?.toLocaleString('id-ID') || 0 }}</p>
          </div>
        </div>

        <!-- Close Shift Form (Right Side - 5 Kolom) -->
        <div class="lg:col-span-5 bg-white border border-slate-200 shadow-lg rounded-3xl overflow-hidden min-h-[480px] flex flex-col">
          <div class="p-6 border-b border-slate-100 bg-slate-50/40">
            <h3 class="text-lg font-black text-slate-900">Penutupan & Rekonsiliasi Shift</h3>
            <p class="text-xs text-slate-500 mt-0.5">Hitung ulang fisik uang tunai di laci untuk pencatatan laporan neraca penutupan.</p>
          </div>
          
          <form @submit.prevent="handleCloseSession" class="p-6 space-y-6 flex-1 flex flex-col justify-between">
            <div class="space-y-5">
              <div class="space-y-2">
                <label class="text-xs font-black text-slate-600 uppercase tracking-wider block">Uang Fisik Aktual Lapangan (Actual Cash)</label>
                <div class="relative rounded-xl shadow-xs">
                  <span class="absolute left-4 top-1/2 -translate-y-1/2 font-bold text-slate-400 text-sm">Rp</span>
                  <input 
                    v-model="formattedActualCash" 
                    type="text" 
                    required
                    class="w-full bg-slate-50 border border-slate-200 rounded-xl pl-12 pr-4 py-3.5 font-mono font-bold text-slate-900 text-lg focus:ring-2 focus:ring-primary focus:border-primary focus:bg-white focus:outline-none transition-all"
                    placeholder="0"
                  />
                </div>
              </div>

              <!-- Indicator Dinamis Selisih Neraca -->
              <div v-if="actualCash !== null" :class="['p-4 rounded-xl flex items-center gap-3 border transition-all', difference === 0 ? 'bg-emerald-50 border-emerald-100 text-emerald-800' : (difference > 0 ? 'bg-blue-50 border-blue-100 text-blue-800' : 'bg-rose-50 border-rose-100 text-rose-800')]">
                <span class="material-symbols-outlined text-3xl shrink-0">
                  {{ difference === 0 ? 'check_circle' : (difference > 0 ? 'trending_up' : 'trending_down') }}
                </span>
                <div>
                  <p class="text-[10px] font-black uppercase tracking-wider">{{ difference === 0 ? 'Status Balance (Seimbang)' : (difference > 0 ? 'Surplus Neraca (Selisih Lebih)' : 'Defisit Neraca (Selisih Kurang)') }}</p>
                  <p class="text-lg font-black font-mono mt-0.5">Rp {{ Math.abs(difference).toLocaleString('id-ID') }}</p>
                </div>
              </div>
            </div>

            <button 
              type="submit" 
              class="w-full bg-rose-600 hover:bg-rose-700 text-white py-4 rounded-xl font-bold text-sm shadow-md hover:-translate-y-0.5 active:translate-y-0 transition-all flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer mt-6"
              :disabled="isSubmitting">
              <span v-if="isSubmitting" class="animate-spin material-symbols-outlined text-sm">sync</span>
              <span v-else class="material-symbols-outlined text-sm">stop_circle</span>
              {{ isSubmitting ? 'Memproses Tutup Sesi...' : 'Tutup Shift & Unduh Laporan' }}
            </button>
          </form>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { cashApi } from '@frontend/api/cash';
import { posApi } from '@frontend/api/pos';
import { useAuthStore } from '@frontend/stores/auth';
import type { CashierSession, ShiftSummary } from '@frontend/types/cash';
import { jsPDF } from 'jspdf';
import autoTable from 'jspdf-autotable';

const authStore = useAuthStore();

const isLoading = ref(true);
const isSubmitting = ref(false);

const currentSession = ref<CashierSession | null>(null);
const summary = ref<ShiftSummary | null>(null);

const openingCash = ref<number | null>(null);
const actualCash = ref<number | null>(null);

const formatCurrencyInput = (val: number | null): string => {
  if (val === null || val === undefined) return '';
  return val.toLocaleString('id-ID');
};

const parseCurrencyInput = (val: string): number | null => {
  if (!val) return null;
  const parsed = parseInt(val.replace(/\D/g, ''), 10);
  return isNaN(parsed) ? null : parsed;
};

const formattedOpeningCash = computed({
  get: () => formatCurrencyInput(openingCash.value),
  set: (val: string) => {
    openingCash.value = parseCurrencyInput(val);
  }
});

const formattedActualCash = computed({
  get: () => formatCurrencyInput(actualCash.value),
  set: (val: string) => {
    actualCash.value = parseCurrencyInput(val);
  }
});

const difference = computed(() => {
  if (actualCash.value === null || !summary.value) return 0;
  return actualCash.value - summary.value.expected_cash;
});

const loadData = async () => {
  isLoading.value = true;
  try {
    const session = await cashApi.getCurrentSession();
    
    if (session && session.id && session.status === 'OPEN') {
      currentSession.value = session;
      const summ = await cashApi.getShiftSummary(session.id);
      summary.value = summ;
    } else {
      currentSession.value = null;
      summary.value = null;
    }
  } catch (error: any) {
    console.error("Gagal memuat sesi kasir:", error);
    alert("Gagal memuat data sesi kasir: " + (error.response?.data?.message || error.message));
    currentSession.value = null;
    summary.value = null;
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadData();
});

const handleOpenSession = async () => {
  if (openingCash.value === null || openingCash.value < 0) return;
  
  isSubmitting.value = true;
  try {
    await cashApi.openSession(openingCash.value);
    alert('Sesi shift kasir berhasil dibuka!');
    await loadData();
  } catch (error: any) {
    alert("Gagal membuka shift: " + (error.response?.data?.message || error.message));
  } finally {
    isSubmitting.value = false;
  }
};

const handleCloseSession = async () => {
  if (actualCash.value === null || actualCash.value < 0) return;
  
  if (!confirm(`Apakah Anda yakin ingin menutup shift ini dengan uang fisik sebesar Rp ${actualCash.value.toLocaleString('id-ID')}?`)) {
    return;
  }

  isSubmitting.value = true;
  try {
    const sessionId = currentSession.value?.id;
    if (!sessionId) throw new Error("Sesi tidak valid");
    
    // Generate PDF Laporan
    const products = await posApi.getTopProducts(sessionId);
    const summ = summary.value;
    const cashierName = authStore.user?.username || 'kasir';

    await cashApi.closeSession(actualCash.value);
    
    if (summ) {
      const doc = new jsPDF('p', 'mm', 'a4');
      doc.setFontSize(18);
      doc.text('Laporan Penjualan Sesi Shift', 14, 20);
      
      doc.setFontSize(12);
      doc.text(`Kasir: ${cashierName}`, 14, 30);
      doc.text(`Waktu Buka: ${formatDate(currentSession.value!.open_time)}`, 14, 36);
      doc.text(`Waktu Tutup: ${formatDate(new Date().toISOString())}`, 14, 42);
      
      const expected = summ.expected_cash || 0;
      const actual = actualCash.value || 0;
      const diff = actual - expected;

      autoTable(doc, {
        startY: 50,
        head: [['Keterangan Ringkasan', 'Jumlah (Rp)']],
        body: [
          ['Saldo Awal Laci', (summ.opening_cash || 0).toLocaleString('id-ID')],
          ['Total Penerimaan POS', (summ.total_income || 0).toLocaleString('id-ID')],
          ['Kas Masuk Ekstra (Manual)', ((summ.total_deposit || 0) - (summ.total_income || 0)).toLocaleString('id-ID')],
          ['Total Kas Keluar (Expenses)', (summ.total_expense || 0).toLocaleString('id-ID')],
          ['Saldo Akhir Ekspektasi', expected.toLocaleString('id-ID')],
          ['Saldo Akhir Aktual', actual.toLocaleString('id-ID')],
          ['Selisih (Variance)', `${diff < 0 ? '-' : diff > 0 ? '+' : ''}${Math.abs(diff).toLocaleString('id-ID')} ${diff === 0 ? '(Seimbang)' : diff > 0 ? '(Surplus)' : '(Minus)'}`],
        ],
        theme: 'grid',
        headStyles: { fillColor: [63, 81, 181], textColor: 255 },
        styles: { fontSize: 11, cellPadding: 5 },
        columnStyles: {
          0: { fontStyle: 'bold' },
          1: { halign: 'right' }
        }
      });
      
      const finalY = (doc as any).lastAutoTable.finalY + 10;
      doc.setFontSize(12);
      doc.text('Rincian Penjualan Produk:', 14, finalY);
      autoTable(doc, {
        startY: finalY + 5,
        head: [['Nama Produk', 'Terjual (Qty)', 'Total Rp']],
        body: products.map(p => [p.product_name, p.total_qty, (p.total_revenue || 0).toLocaleString('id-ID')]),
        theme: 'striped',
        headStyles: { fillColor: [41, 128, 185] },
        styles: { fontSize: 10 }
      });

      const pdfBlob = doc.output('blob');
      const filename = `${cashierName}.pdf`;

      // Upload file ke Minio di latar belakang
      cashApi.uploadShiftReport(sessionId, pdfBlob, filename).catch(e => {
        console.error("Gagal mengunggah laporan PDF:", e);
      });
    }

    alert('Sesi shift kasir berhasil ditutup!');
    currentSession.value = null;
    summary.value = null;
    openingCash.value = null;
    actualCash.value = null;
    await loadData();
  } catch (error: any) {
    alert("Gagal menutup shift: " + (error.response?.data?.message || error.message));
  } finally {
    isSubmitting.value = false;
  }
};

const formatDate = (dateString: string) => {
  if (!dateString) return '-';
  const d = new Date(dateString);
  return d.toLocaleString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
};

const formatNumber = (amount: number) => {
  return new Intl.NumberFormat('id-ID', { minimumFractionDigits: 0 }).format(amount);
};

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value);
};
</script>

