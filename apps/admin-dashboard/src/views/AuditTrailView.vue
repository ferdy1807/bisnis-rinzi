<template>
  <div class="p-1 animate-fade-in space-y-6">
    <!-- TOP BAR KONTROL JEJAK AUDIT -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-outline-variant/20 pb-6">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">Log Jejak Audit Aktivitas (Audit Trail)</h2>
        <p class="text-body-md text-on-surface-variant mt-1">
          Otoritas Kemananan Owner: Rekaman kronologis mutasi data sensitif, pelacakan aksi manipulasi sistem, dan forensik digital multi-skema.
        </p>
      </div>
      <div class="flex items-center gap-3">
        <button
          @click="loadAuditTrailManifest"
          class="flex items-center gap-2 bg-surface border border-outline-variant/60 text-primary px-4 py-2.5 rounded-xl shadow-sm hover:bg-primary/5 active:scale-95 transition-all text-label-md font-bold"
        >
          <span class="material-symbols-outlined text-md">refresh</span>
          Segarkan Log
        </button>
      </div>
    </div>

    <!-- FILTER BOARD KONTROL FORENSIK -->
    <div class="bg-surface-container-low p-4 rounded-2xl border border-outline-variant/30 grid grid-cols-1 sm:grid-cols-3 gap-4 shadow-sm">
      <div>
        <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block mb-1">Cari ID / Nama Entitas</label>
        <div class="relative flex items-center">
          <span class="material-symbols-outlined absolute left-3 text-outline text-md">search</span>
          <input v-model="searchQuery" type="text" placeholder="Masukkan ID entitas target..." class="w-full bg-surface border border-outline-variant/50 rounded-xl pl-9 pr-3 py-2 text-xs focus:outline-none focus:ring-1 focus:ring-primary" />
        </div>
      </div>
      <div>
        <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block mb-1">Filter Tipe Aksi Mutasi</label>
        <select v-model="filterAction" class="w-full bg-surface border border-outline-variant/50 rounded-xl px-3 py-2 text-xs text-on-surface focus:outline-none">
          <option value="ALL">Semua Aksi (INSERT / UPDATE / DELETE)</option>
          <option value="INSERT">INSERT (Pencatatan Baru)</option>
          <option value="UPDATE">UPDATE (Pengubahan Parameter)</option>
          <option value="DELETE">DELETE (Penghapusan Rekor)</option>
        </select>
      </div>
      <div>
        <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block mb-1">Urutan Kronologis</label>
        <select v-model="sortOrder" class="w-full bg-surface border border-outline-variant/50 rounded-xl px-3 py-2 text-xs text-on-surface focus:outline-none">
          <option value="DESC">Terbaru Menuju Terlama</option>
          <option value="ASC">Terlama Menuju Terbaru</option>
        </select>
      </div>
    </div>

    <!-- MAIN LOADING INDICATOR -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[350px] text-primary">
      <span class="animate-spin material-symbols-outlined text-4xl mb-2">sync</span>
      <p class="text-label-md font-bold tracking-wider">Menyusun Manifes Kronologis Forensik Karyawan...</p>
    </div>

    <div v-else class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden">
      <div class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex items-center justify-between">
        <h3 class="font-title-medium text-primary font-bold flex items-center gap-2">
          <span class="material-symbols-outlined text-md">assignment_ind</span> Catatan Log Transaksional Sistem (audit_logs)
        </h3>
        <span class="text-[10px] font-mono bg-surface-container-high px-3 py-1 rounded-full font-bold text-on-surface-variant">
          Otorisasi Jalur: /api/auth/audit-logs
        </span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold border-b border-outline-variant/20">
              <th class="px-6 py-4">Stempel Waktu UTC</th>
              <th class="px-6 py-4">Aksi Operasional</th>
              <th class="px-6 py-4">Nama Entitas Skema</th>
              <th class="px-6 py-4 font-mono">Target Entity ID</th>
              <th class="px-6 py-4">ID Pegawai Pelaksana</th>
              <th class="px-6 py-4 text-center">Komparasi Data</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/10 text-body-sm font-mono text-xs">
            <tr v-if="filteredLogs.length === 0">
              <td colspan="6" class="px-6 py-8 text-center font-sans text-on-surface-variant">Tidak ada jejak log audit yang cocok dengan parameter kueri pencarian.</td>
            </tr>
            <tr v-for="log in filteredLogs" :key="log.id" class="hover:bg-surface-container-low/30 transition-colors">
              <td class="px-6 py-4 text-on-surface-variant">{{ log.created_at }}</td>
              <td class="px-6 py-4 font-sans">
                <span :class="getActionBadgeClass(log.action)" class="px-2 py-0.5 rounded text-[9px] font-black uppercase tracking-wider border">
                  {{ log.action }}
                </span>
              </td>
              <td class="px-6 py-4 font-sans text-on-surface font-bold text-xs uppercase tracking-wide">{{ log.user_name || '-' }}</td>
              <td class="px-6 py-4 text-primary font-bold text-xs">{{ log.entity_id || '-' }}</td>
              <td class="px-6 py-4 font-sans text-on-surface-variant">{{ log.user_id }}</td>
              <td class="px-6 py-4 text-center font-sans">
                <button
                  @click="openInspectionModal(log)"
                  class="inline-flex items-center gap-1 bg-surface border border-outline-variant/60 text-secondary px-3 py-1.5 rounded-lg hover:bg-secondary/5 active:scale-95 transition-all text-label-sm font-bold"
                >
                  <span class="material-symbols-outlined text-sm">visibility</span>
                  Inspeksi Payload
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL WINDOW PAYLOAD INSPECTOR (COMPARE OLD DATA VS NEW DATA) -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="bg-surface rounded-2xl w-full max-w-4xl p-6 shadow-xl relative border border-outline-variant/20 max-h-[85vh] overflow-y-auto">
        <button @click="showModal = false" class="absolute top-6 right-6 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1 transition-colors">
          <span class="material-symbols-outlined">close</span>
        </button>

        <h3 class="text-title-lg font-bold text-primary mb-1 flex items-center gap-1.5">
          <span class="material-symbols-outlined">analytics</span> Struktur Forensik Payload Objek JSON
        </h3>
        <p class="text-xs text-on-surface-variant font-mono mb-6">Target Entity: {{ selectedLog?.entity_name }} ({{ selectedLog?.entity_id }})</p>

        <!-- GRID DATA DUA KOLOM KOMPARASI STRUKTUR DATABAS -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Kolom A: Keadaan Objek Lama (old_data) -->
          <div class="space-y-1.5">
            <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block">Kondisi Sebelum Mutasi (`old_data`)</span>
            <pre class="bg-surface-container-lowest text-error p-4 rounded-xl border border-outline-variant/40 text-xs font-mono overflow-x-auto max-h-64 leading-relaxed"><code>{{ selectedLog?.old_data ? JSON.stringify(selectedLog.old_data, null, 2) : '// Tidak ada data awal (Aksi INSERT).' }}</code></pre>
          </div>
          <!-- Kolom B: Keadaan Objek Baru (new_data) -->
          <div class="space-y-1.5">
            <span class="text-[10px] font-black text-teal-600 uppercase tracking-wider block">Kondisi Sesudah Mutasi (`new_data`)</span>
            <pre class="bg-surface-container-lowest text-teal-600 p-4 rounded-xl border border-outline-variant/40 text-xs font-mono overflow-x-auto max-h-64 leading-relaxed"><code>{{ selectedLog?.new_data ? JSON.stringify(selectedLog.new_data, null, 2) : '// Tidak ada data sisa (Aksi DELETE).' }}</code></pre>
          </div>
        </div>

        <div class="flex justify-end pt-6 mt-4 border-t border-outline-variant/20">
          <button type="button" @click="showModal = false" class="px-5 py-2 bg-primary text-on-primary rounded-lg text-label-md font-bold shadow-md hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all">
            Selesai Tinjauan
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { authApi } from '@frontend/api/auth';
import type { AuditLog } from '@frontend/types/auth';

const isLoading = ref(true);
const showModal = ref(false);

const searchQuery = ref('');
const filterAction = ref('ALL');
const sortOrder = ref<'ASC' | 'DESC'>('DESC');

const auditLogsRawList = ref<AuditLog[]>([]);
const selectedLog = ref<AuditLog | null>(null);

const loadAuditTrailManifest = async () => {
  isLoading.value = true;
  try {
    // Memanggil authApi.getAuditLogs() untuk penarikan data bersih rute /api/auth/audit-logs
    const response = await authApi.getAuditLogs();
    auditLogsRawList.value = Array.isArray(response) ? response : (response as any).data || [];
  } catch (err) {
    console.error('Gagal memuat manifes audit trail forensik:', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadAuditTrailManifest();
});

// --- LOGIKA FILTERING & PENGURUTAN LIVE REAKTIF ---
const filteredLogs = computed(() => {
  let output = [...auditLogsRawList.value];

  // 1. Pemfilteran Kriteria Aksi
  if (filterAction.value !== 'ALL') {
    output = output.filter(log => log.action === filterAction.value);
  }

  // 2. Pemfilteran Berdasarkan String Pencarian ID/Nama Entitas
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    output = output.filter(log => 
      log.entity_id?.toLowerCase().includes(query) || 
      log.entity_name?.toLowerCase().includes(query)
    );
  }

  // 3. Pengurutan Kronologis Waktu
  output.sort((a, b) => {
    const timeA = new Date(a.created_at || '').getTime();
    const timeB = new Date(b.created_at || '').getTime();
    return sortOrder.value === 'DESC' ? timeB - timeA : timeA - timeB;
  });

  return output;
});

const openInspectionModal = (log: AuditLog) => {
  selectedLog.value = log;
  showModal.value = true;
};

const getActionBadgeClass = (action: string | undefined): string => {
  switch (action) {
    case 'INSERT': return 'bg-success/10 text-success border-success/20';
    case 'UPDATE': return 'bg-warning/10 text-warning border-warning/20';
    case 'DELETE': return 'bg-error/10 text-error border-error/20';
    default: return 'bg-outline-variant/10 text-on-surface-variant border-outline-variant/20';
  }
};
</script>