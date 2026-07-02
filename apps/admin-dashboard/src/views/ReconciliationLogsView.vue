<template>
  <div class="p-1 animate-fade-in space-y-6">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-outline-variant/20 pb-6">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">Log Hasil Rekonsiliasi</h2>
        <p class="text-body-md text-on-surface-variant mt-1">
          Otoritas Pemeriksaan Deviasi Owner: Riwayat kecocokan saldo kas, pelacakan ketidakseimbangan sistem, dan audit trail kliring harian.
        </p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="fetchLogsFromService"
          class="flex items-center gap-2 bg-surface border border-outline-variant/60 text-primary px-4 py-2.5 rounded-xl shadow-sm hover:bg-primary/5 active:scale-95 transition-all text-label-md font-bold"
        >
          <span class="material-symbols-outlined text-md">sync</span>
          Muat Ulang Data
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
      <div class="bg-surface border border-outline-variant/30 rounded-2xl p-5 shadow-sm">
        <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest block">Total Sesi Balanced</span>
        <h3 class="text-title-xl font-black text-teal-600 mt-1">{{ totalBalancedSessions }} Sesi</h3>
      </div>
      <div class="bg-surface border border-outline-variant/30 rounded-2xl p-5 shadow-sm">
        <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest block">Sesi Deviasi Discrepancy</span>
        <h3 class="text-title-xl font-black text-error mt-1">{{ totalDiscrepancySessions }} Sesi</h3>
      </div>
      <div class="bg-surface border border-outline-variant/30 rounded-2xl p-5 shadow-sm">
        <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest block">Akumulasi Selisih Kas</span>
        <h3 class="text-title-xl font-black text-on-surface mt-1" :class="accumulatedDiscrepancy === 0 ? 'text-teal-600' : 'text-error'">
          Rp {{ formatCurrency(accumulatedDiscrepancy) }}
        </h3>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[350px] text-primary">
      <span class="animate-spin material-symbols-outlined text-4xl mb-2">sync</span>
      <p class="text-label-md font-bold tracking-wider">Menarik Riwayat Rekonsiliasi Kasir...</p>
    </div>

    <div v-else class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden">
      <div class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex items-center justify-between">
        <h3 class="font-title-medium text-primary font-bold flex items-center gap-2">
          <span class="material-symbols-outlined text-md">rule</span> Manifes Rekonsiliasi Multi-Sistem (reconciliation_logs)
        </h3>
        <span class="text-[10px] font-mono bg-surface-container-high px-3 py-1 rounded-full font-bold text-on-surface-variant">
          Endpoint: /api/finance/reconciliation/logs
        </span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold border-b border-outline-variant/20">
              <th class="px-6 py-4">Waktu Rekon</th>
              <th class="px-6 py-4">ID Penutupan Harian</th>
              <th class="px-6 py-4">Sistem Target</th>
              <th class="px-6 py-4 text-right">Jumlah Sistem</th>
              <th class="px-6 py-4 text-right">Jumlah Fisik</th>
              <th class="px-6 py-4 text-right">Selisih Deviasi</th>
              <th class="px-6 py-4">Berita Acara / Catatan Staf</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/10 text-body-sm font-mono text-xs">
            <tr v-if="logsList.length === 0">
              <td colspan="7" class="px-6 py-8 text-center font-sans text-on-surface-variant">Tidak ditemukan riwayat log hasil rekonsiliasi.</td>
            </tr>
            <tr v-for="log in logsList" :key="log.id" class="hover:bg-surface-container-low/30 transition-colors">
              <td class="px-6 py-4 font-sans text-on-surface-variant text-xs">{{ log.created_at }}</td>
              <td class="px-6 py-4 text-primary font-bold">{{ log.daily_closing_id?.substring(0, 8) }}...</td>
              <td class="px-6 py-4 font-sans"><span class="px-2 py-0.5 rounded bg-surface-container-high border text-on-surface font-bold text-[10px] tracking-wider uppercase">{{ log.target_system }}</span></td>
              <td class="px-6 py-4 text-right">Rp {{ formatCurrency(log.system_amount) }}</td>
              <td class="px-6 py-4 text-right">Rp {{ formatCurrency(log.actual_amount) }}</td>
              <td class="px-6 py-4 text-right font-bold" :class="log.discrepancy === 0 ? 'text-teal-600' : 'text-error'">
                Rp {{ formatCurrency(log.discrepancy) }}
              </td>
              <td class="px-6 py-4 font-sans text-on-surface-variant italic max-w-xs truncate" :title="log.notes">
                "{{ log.notes || 'Sesi kliring seimbang.' }}"
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { financeApi } from '@frontend/api/finance';
import type { ReconciliationLog } from '@frontend/types/finance';

const isLoading = ref(true);
const logsList = ref<ReconciliationLog[]>([]);

const fetchLogsFromService = async () => {
  isLoading.value = true;
  try {
    // Memanggil GET /api/finance/reconciliation/logs sesuai manifestasi berkas endpoint.md
    const response = await financeApi.getReconciliationLogs();
    logsList.value = Array.isArray(response) ? response : (response as any).data || [];
  } catch (err) {
    console.error('Gagal menarik manifes riwayat log rekonsiliasi:', err);
  } {
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchLogsFromService();
});

// --- LOGIKA AGREGASI KPI HISTORIS ---
const totalBalancedSessions = computed(() => logsList.value.filter(log => log.discrepancy === 0).length);
const totalDiscrepancySessions = computed(() => logsList.value.filter(log => log.discrepancy !== 0).length);
const accumulatedDiscrepancy = computed(() => logsList.value.reduce((sum, log) => sum + (log.discrepancy || 0), 0));

const formatCurrency = (val: number): string => {
  return new Intl.NumberFormat('id-ID', { minimumFractionDigits: 0 }).format(val);
};
</script>