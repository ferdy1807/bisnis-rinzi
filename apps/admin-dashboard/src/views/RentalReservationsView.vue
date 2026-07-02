<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { rentalApi } from '@frontend/api/rental';
import type { RentalReservation } from '@frontend/types/rental';

// --- STATE MANAGEMENT ---
type FilterStatus = 'ALL' | 'PICKED_UP' | 'BOOKED' | 'OVERDUE';
const activeTab = ref<FilterStatus>('ALL');
const isLoading = ref(false);
const searchQuery = ref('');
const errorMessage = ref('');

// Data Buckets
const allReservations = ref<RentalReservation[]>([]);
const activeReservations = ref<RentalReservation[]>([]);
const upcomingReservations = ref<RentalReservation[]>([]);
const overdueReservations = ref<RentalReservation[]>([]);

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

// --- FETCH DATA OPERATIONS ---
const fetchReservationsData = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    // Memuat seluruh kategori status reservasi logistik secara paralel dari rentalApi
    const [allRes, activeRes, upcomingRes, overdueRes] = await Promise.all([
      rentalApi.getAllReservations(),
      rentalApi.getActiveReservations(),
      rentalApi.getUpcomingReservations(),
      rentalApi.getOverdueReservations()
    ]);

    allReservations.value = allRes;
    activeReservations.value = activeRes.filter(r => r.status === 'PICKED_UP');
    upcomingReservations.value = upcomingRes;
    overdueReservations.value = overdueRes;
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat repositori data reservasi kontrak rental.';
  } finally {
    isLoading.value = false;
  }
};

// --- SWITCH TAB CONTROLLER ---
const currentReservationsList = computed(() => {
  let targetList = allReservations.value;
  if (activeTab.value === 'PICKED_UP') targetList = activeReservations.value;
  if (activeTab.value === 'BOOKED') targetList = upcomingReservations.value;
  if (activeTab.value === 'OVERDUE') targetList = overdueReservations.value;

  // Filter out RETURNED data globally based on user request
  targetList = targetList.filter(res => res.status !== 'RETURNED');

  // Filter pencarian berdasarkan nama pelanggan atau nomor invoice
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    return targetList.filter(res => 
      res.customer_name?.toLowerCase().includes(query) ||
      res.invoice_number?.toLowerCase().includes(query) ||
      res.id?.toLowerCase().includes(query)
    );
  }
  return targetList;
});

onMounted(() => {
  fetchReservationsData();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-calendar text-indigo-600"></i>
        Kontrak & Reservasi Rental Aktif
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Pusat kendali dan pengawasan dokumen kontrak sewa unit untuk memastikan validitas sirkulasi aset operasional.
      </p>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>

    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div @click="activeTab = 'ALL'" class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm cursor-pointer hover:border-slate-300 transition-all" :class="{ 'ring-2 ring-indigo-500/20 bg-indigo-50/20': activeTab === 'ALL' }">
        <span class="text-xs text-slate-400 font-bold uppercase">Semua Berkas</span>
        <p class="text-2xl font-black text-slate-800 mt-1 font-mono">{{ allReservations.length }}</p>
      </div>
      <div @click="activeTab = 'PICKED_UP'" class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm cursor-pointer hover:border-slate-300 transition-all" :class="{ 'ring-2 ring-emerald-500/20 bg-emerald-50/20': activeTab === 'PICKED_UP' }">
        <span class="text-xs text-slate-400 font-bold uppercase text-emerald-600">Sedang Disewa (PICKED_UP)</span>
        <p class="text-2xl font-black text-emerald-600 mt-1 font-mono">{{ activeReservations.length }}</p>
      </div>
      <div @click="activeTab = 'BOOKED'" class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm cursor-pointer hover:border-slate-300 transition-all" :class="{ 'ring-2 ring-blue-500/20 bg-blue-50/20': activeTab === 'BOOKED' }">
        <span class="text-xs text-slate-400 font-bold uppercase text-blue-600">Belum Dihias (BOOKED)</span>
        <p class="text-2xl font-black text-blue-600 mt-1 font-mono">{{ upcomingReservations.length }}</p>
      </div>
      <div @click="activeTab = 'OVERDUE'" class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm cursor-pointer hover:border-slate-300 transition-all" :class="{ 'ring-2 ring-rose-500/20 bg-rose-50/20': activeTab === 'OVERDUE' }">
        <span class="text-xs text-slate-400 font-bold uppercase text-rose-600">Terlambat Kembali</span>
        <p class="text-2xl font-black text-rose-600 mt-1 font-mono">{{ overdueReservations.length }}</p>
      </div>
    </div>

    <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center gap-3 justify-between">
      <div class="relative w-full sm:w-80">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <i class="pi pi-search text-sm"></i>
        </span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari nama pelanggan atau nomor invoice..."
          class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
        />
      </div>

      <button 
        @click="fetchReservationsData"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border transition-colors"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Sinkronisasi Manifest
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Mengekstrak logistik data rental...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">ID / No. Invoice</th>
              <th class="py-4 px-6">Nama Pelanggan</th>
              <th class="py-4 px-6">Rentang Waktu Sewa</th>
              <th class="py-4 px-6 text-center">Status</th>
              <th class="py-4 px-6 text-right">Uang Muka (DP)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="currentReservationsList.length === 0 && !isLoading">
              <td colspan="5" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada data reservasi rental yang sesuai dengan kriteria filter saat ini.
              </td>
            </tr>
            <tr v-for="res in currentReservationsList" :key="res.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900">{{ res.invoice_number || 'REG-RENTAL' }}</div>
                <div class="text-[10px] text-slate-400 font-mono select-all">{{ res.id }}</div>
              </td>
              <td class="py-4 px-6">
                <div class="font-semibold text-slate-900">{{ res.customer_name }}</div>
                <div class="text-[11px] text-slate-500 font-normal">Tlp: {{ res.customer_phone || '-' }}</div>
              </td>
              <td class="py-4 px-6 text-xs text-slate-600 font-normal">
                <div>Mulai: {{ formatDate(res.start_date) }}</div>
                <div class="mt-0.5">Selesai: {{ formatDate(res.end_date) }}</div>
              </td>
              <td class="py-4 px-6 text-center">
                <span 
                  class="px-2.5 py-1 text-[10px] font-extrabold rounded-md uppercase border"
                  :class="{
                    'bg-emerald-50 text-emerald-700 border-emerald-200': (res as any).status === 'ACTIVE' || activeTab === 'BOOKED',
                    'bg-blue-50 text-blue-700 border-blue-200': (res as any).status === 'UPCOMING' || activeTab === 'PICKED_UP',
                    'bg-rose-50 text-rose-700 border-rose-200': (res as any).status === 'OVERDUE' || activeTab === 'OVERDUE',
                    'bg-slate-50 text-slate-700 border-slate-200': activeTab === 'ALL' && !(res as any).status
                  }"
                >
                  {{ (res as any).status || activeTab }}
                </span>
              </td>
              <td class="py-4 px-6 text-right font-mono font-bold text-slate-900">
                {{ formatCurrency(res.down_payment || 0) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>