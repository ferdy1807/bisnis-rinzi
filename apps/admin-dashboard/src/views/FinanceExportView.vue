<script setup lang="ts">
import { ref } from 'vue';
import { financeApi } from '@frontend/api/finance';

// --- STATE MANAGEMENT ---
const isExporting = ref(false);
const selectedFormat = ref<'pdf' | 'excel' | 'csv'>('pdf');
const exportSuccessMessage = ref('');
const errorMessage = ref('');

// Form Filter Tambahan untuk Parameter Cetak (Opsional)
const filterOptions = ref({
  includeNotes: true,
  signatureAuthority: 'OWNER - Ferdy Perdana'
});

// --- EXECUTE EXPORT DOWNLOAD ---
const handleExportDownload = () => {
  isExporting.value = true;
  errorMessage.value = '';
  exportSuccessMessage.value = '';

  try {
    // Ambil endpoint URL unduhan fisik dari helper generator URL FinanceApi
    const targetUrl = financeApi.getExportUrl(selectedFormat.value);
    
    // Validasi dasar kesiapan URL dokumen resmi
    if (!targetUrl) {
      throw new Error('Alamat server unduhan laporan tidak terkonfigurasi dengan benar.');
    }

    // Eksekusi trigger download via native window anchor bypass
    const downloadAnchor = document.createElement('a');
    downloadAnchor.href = targetUrl;
    downloadAnchor.target = '_blank';
    // Berikan nama file tentatif yang rapi
    downloadAnchor.download = `Laporan_Konsolidasi_Finansial_${new Date().getFullYear()}.${selectedFormat.value === 'excel' ? 'xlsx' : selectedFormat.value}`;
    
    document.body.appendChild(downloadAnchor);
    downloadAnchor.click();
    document.body.removeChild(downloadAnchor);

    exportSuccessMessage.value = `Dokumen resmi berformat ${selectedFormat.value.toUpperCase()} berhasil dialihkan ke antrean unduhan browser Anda.`;
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal mengekstrak berkas laporan resmi dari server akuntansi.';
  } {
    isExporting.value = false;
  }
};
</script>

<template>
  <div class="p-6 max-w-4xl mx-auto space-y-6">
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-download text-indigo-600"></i>
        Ekspor Dokumen Finansial Resmi
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Gerbang validasi fisik pembukuan perusahaan. Unduh data finansial konsolidasi yang bersih dan sah langsung dari repositori ledger utama.
      </p>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl flex items-start gap-3 animate-fade-in">
      <i class="pi pi-exclamation-circle text-lg mt-0.5"></i>
      <span class="text-sm font-medium">{{ errorMessage }}</span>
    </div>

    <div v-if="exportSuccessMessage" class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl flex items-start gap-3 animate-fade-in">
      <i class="pi pi-check-circle text-lg mt-0.5"></i>
      <span class="text-sm font-medium">{{ exportSuccessMessage }}</span>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden grid grid-cols-1 md:grid-cols-3 divide-y md:divide-y-0 md:divide-x divide-slate-200">
      
      <div class="p-6 space-y-4 md:col-span-1">
        <h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider">1. Pilih Format Berkas</h3>
        <div class="space-y-2">
          <label 
            class="flex items-center gap-3 p-3 border rounded-xl cursor-pointer transition-all"
            :class="selectedFormat === 'pdf' ? 'border-indigo-600 bg-indigo-50/40 text-indigo-900 font-semibold shadow-sm' : 'border-slate-200 text-slate-600 hover:bg-slate-50'"
          >
            <input type="radio" v-model="selectedFormat" value="pdf" class="text-indigo-600 focus:ring-indigo-500" />
            <i class="pi pi-file-pdf text-xl" :class="selectedFormat === 'pdf' ? 'text-indigo-600' : 'text-slate-400'"></i>
            <div class="text-xs">
              <p class="font-bold">Portable Document (PDF)</p>
              <p class="text-[10px] text-slate-400 font-normal">Arsip Siap Cetak & Legalitas</p>
            </div>
          </label>

          <label 
            class="flex items-center gap-3 p-3 border rounded-xl cursor-pointer transition-all"
            :class="selectedFormat === 'excel' ? 'border-emerald-600 bg-emerald-50/40 text-emerald-900 font-semibold shadow-sm' : 'border-slate-200 text-slate-600 hover:bg-slate-50'"
          >
            <input type="radio" v-model="selectedFormat" value="excel" class="text-emerald-600 focus:ring-emerald-500" />
            <i class="pi pi-file-excel text-xl" :class="selectedFormat === 'excel' ? 'text-emerald-600' : 'text-slate-400'"></i>
            <div class="text-xs">
              <p class="font-bold">Spreadsheet (Excel)</p>
              <p class="text-[10px] text-slate-400 font-normal">Analisis Formula & Pivot Lanjutan</p>
            </div>
          </label>

          <label 
            class="flex items-center gap-3 p-3 border rounded-xl cursor-pointer transition-all"
            :class="selectedFormat === 'csv' ? 'border-amber-600 bg-amber-50/40 text-amber-900 font-semibold shadow-sm' : 'border-slate-200 text-slate-600 hover:bg-slate-50'"
          >
            <input type="radio" v-model="selectedFormat" value="csv" class="text-amber-600 focus:ring-amber-500" />
            <i class="pi pi-file text-xl" :class="selectedFormat === 'csv' ? 'text-amber-600' : 'text-slate-400'"></i>
            <div class="text-xs">
              <p class="font-bold">Comma Separated (CSV)</p>
              <p class="text-[10px] text-slate-400 font-normal">Integrasi Data Mentah antar Aplikasi</p>
            </div>
          </label>
        </div>
      </div>

      <div class="p-6 flex flex-col justify-between md:col-span-2 space-y-6 bg-slate-50/50">
        <div class="space-y-4">
          <h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider">2. Aturan Metadata Penandatanganan</h3>
          
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-[10px] font-bold text-slate-500 uppercase">Otoritas Penandatangan Dokumen</label>
              <input 
                v-model="filterOptions.signatureAuthority"
                type="text"
                disabled
                class="w-full px-3 py-2 border border-slate-200 rounded-lg text-xs bg-slate-100 text-slate-500 font-medium cursor-not-allowed"
              />
            </div>
            <div class="flex items-center h-full pt-4 sm:pt-6 pl-1">
              <label class="inline-flex items-center gap-2 cursor-pointer select-none">
                <input 
                  type="checkbox" 
                  v-model="filterOptions.includeNotes" 
                  class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
                />
                <span class="text-xs font-semibold text-slate-600">Lampirkan Log Catatan Audit Sistem</span>
              </label>
            </div>
          </div>

          <div class="bg-indigo-50/80 border border-indigo-100 rounded-xl p-4 text-xs text-indigo-800 leading-relaxed space-y-1">
            <p class="font-bold flex items-center gap-1.5">
              <i class="pi pi-info-circle"></i> Catatan Kepatuhan Akuntansi (Compliance)
            </p>
            <p class="text-indigo-950/80">
              Dokumen hasil ekstraksi ini dilindungi oleh modul log jejak audit internal sistem. Setiap berkas fisik yang diterbitkan dari halaman ini akan mencerminkan angka penutupan harian (*Daily Closing*) yang sah dan tidak dapat direkayasa sepihak.
            </p>
          </div>
        </div>

        <div class="pt-4 border-t border-slate-200/60 flex items-center justify-end">
          <button 
            @click="handleExportDownload"
            :disabled="isExporting"
            class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-6 py-3 text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-400 rounded-xl shadow-md shadow-indigo-600/10 transition-all duration-200 cursor-pointer disabled:cursor-not-allowed"
          >
            <i :class="isExporting ? 'pi pi-spin pi-spinner' : 'pi pi-cloud-download'"></i>
            <span>{{ isExporting ? 'Memproses Berkas Resmi...' : `Unduh Dokumen ${selectedFormat.toUpperCase()}` }}</span>
          </button>
        </div>

      </div>
    </div>
  </div>
</template>