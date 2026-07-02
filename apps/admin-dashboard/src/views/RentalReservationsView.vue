<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { rentalApi } from '@frontend/api/rental';
import type { RentalReservation } from '@frontend/types/rental';

// --- STATE MANAGEMENT ---
type FilterStatus = 'ALL' | 'PICKED_UP' | 'BOOKED' | 'OVERDUE' | 'RETURNED';
const activeTab = ref<FilterStatus>('ALL');
const isLoading = ref(false);
const searchQuery = ref('');
const errorMessage = ref('');

// Data Buckets
const allReservations = ref<RentalReservation[]>([]);
const activeReservations = ref<RentalReservation[]>([]);
const upcomingReservations = ref<RentalReservation[]>([]);
const overdueReservations = ref<RentalReservation[]>([]);
const returnedReservations = ref<RentalReservation[]>([]); // State baru untuk data RETURNED

const activeTabAllCount = computed(() => {
  return allReservations.value.filter(res => res.status !== 'RETURNED' && res.status !== 'CANCELLED').length;
});

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
    // Memuat seluruh reservasi sebagai single source of truth
    const allRes = await rentalApi.getAllReservations();

    allReservations.value = allRes;
    
    // Sedang disewa (termasuk yang terlambat)
    activeReservations.value = allRes.filter(r => r.status === 'PICKED_UP');
    
    // Belum diambil (Booked / Siap diambil)
    upcomingReservations.value = allRes.filter(r => r.status === 'BOOKED' || r.status === 'READY_FOR_PICKUP');
    
    // Keterlambatan (Overdue)
    const now = new Date();
    now.setHours(0,0,0,0);
    overdueReservations.value = allRes.filter(r => {
      if (r.status !== 'PICKED_UP' || !r.end_date) return false;
      const endDate = new Date(r.end_date);
      endDate.setHours(0,0,0,0);
      return now > endDate;
    });
    
    // Selesai dikembalikan
    returnedReservations.value = allRes.filter(r => r.status === 'RETURNED');
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
  if (activeTab.value === 'RETURNED') {
    targetList = returnedReservations.value;
  } else {
    // Sembunyikan berkas selesai & batal di tab lain agar antrean aktif tetap bersih
    targetList = targetList.filter(res => res.status !== 'RETURNED' && res.status !== 'CANCELLED');
  }

  // Filter pencarian berdasarkan nama pelanggan, nomor invoice, atau ID
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase().trim();
    return targetList.filter(res => 
      res.customer_name?.toLowerCase().includes(query) ||
      res.invoice_number?.toLowerCase().includes(query) ||
      res.id?.toLowerCase().includes(query)
    );
  }
  return targetList;
});

// Helper untuk penentuan kosmetik label status pada baris tabel
const getStatusBadgeClass = (status?: string) => {
  const s = status || activeTab.value;
  if (s === 'PICKED_UP' || s === 'ACTIVE') return 'bg-emerald-50 text-emerald-700 border-emerald-200';
  if (s === 'BOOKED' || s === 'UPCOMING') return 'bg-indigo-50 text-indigo-700 border-indigo-200';
  if (s === 'OVERDUE') return 'bg-rose-50 text-rose-700 border-rose-200';
  if (s === 'RETURNED') return 'bg-slate-100 text-slate-700 border-slate-300';
  return 'bg-slate-50 text-slate-600 border-slate-200';
};

onMounted(() => {
  fetchReservationsData();
});
</script>

<template>
  <div class="p-6 max-w-[1440px] mx-auto space-y-6 font-sans text-slate-800 antialiased">
    
    <!-- HEADER SECTION -->
    <div class="border-b border-slate-200 pb-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <p class="text-[11px] font-bold tracking-widest text-indigo-600 uppercase mb-1">Manajemen Kontrak</p>
        <h1 class="text-2xl font-extrabold text-slate-900 tracking-tight flex items-center gap-2.5">
          <i class="pi pi-calendar text-indigo-600 bg-indigo-50 p-2 rounded-xl text-xl"></i>
          Kontrak & Reservasi Rental
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Pusat kendali dan pengawasan dokumen kontrak sewa unit untuk memastikan validitas sirkulasi aset operasional.
        </p>
      </div>
    </div>

    <!-- ERROR ALERT -->
    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium shadow-3xs">
      {{ errorMessage }}
    </div>

    <!-- METRIC KPI SUMMARY BENTO CARDS -->
    <div class="grid grid-cols-2 lg:grid-cols-5 gap-4">
      <!-- Card 1: ALL -->
      <div 
        @click="activeTab = 'ALL'" 
        class="bg-white p-5 rounded-2xl border transition-all duration-200 cursor-pointer shadow-3xs flex flex-col justify-between min-h-[110px]" 
        :class="activeTab === 'ALL' ? 'border-slate-900 bg-slate-50/50 ring-2 ring-slate-900/10' : 'border-slate-200 hover:border-slate-400'"
      >
        <div class="flex items-center justify-between">
          <span class="text-[11px] text-slate-400 font-bold uppercase tracking-wider">Semua Berkas</span>
          <i class="pi pi-folder text-slate-400 text-xs"></i>
        </div>
        <p class="text-2xl font-black text-slate-900 mt-2 m-0 font-mono">{{ activeTabAllCount }}</p>
      </div>

      <!-- Card 2: BOOKED -->
      <div 
        @click="activeTab = 'BOOKED'" 
        class="bg-white p-5 rounded-2xl border transition-all duration-200 cursor-pointer shadow-3xs flex flex-col justify-between min-h-[110px]" 
        :class="activeTab === 'BOOKED' ? 'border-indigo-600 bg-indigo-50/30 ring-2 ring-indigo-500/10' : 'border-slate-200 hover:border-indigo-300'"
      >
        <div class="flex items-center justify-between">
          <span class="text-[11px] text-indigo-600 font-bold uppercase tracking-wider">Belum Dihias</span>
          <i class="pi pi-bookmark text-indigo-400 text-xs"></i>
        </div>
        <p class="text-2xl font-black text-indigo-600 mt-2 m-0 font-mono">{{ upcomingReservations.length }}</p>
      </div>

      <!-- Card 3: PICKED_UP -->
      <div 
        @click="activeTab = 'PICKED_UP'" 
        class="bg-white p-5 rounded-2xl border transition-all duration-200 cursor-pointer shadow-3xs flex flex-col justify-between min-h-[110px]" 
        :class="activeTab === 'PICKED_UP' ? 'border-emerald-600 bg-emerald-50/30 ring-2 ring-emerald-500/10' : 'border-slate-200 hover:border-emerald-300'"
      >
        <div class="flex items-center justify-between">
          <span class="text-[11px] text-emerald-600 font-bold uppercase tracking-wider">Sedang Disewa</span>
          <i class="pi pi-external-link text-emerald-400 text-xs"></i>
        </div>
        <p class="text-2xl font-black text-emerald-600 mt-2 m-0 font-mono">{{ activeReservations.length }}</p>
      </div>

      <!-- Card 4: OVERDUE -->
      <div 
        @click="activeTab = 'OVERDUE'" 
        class="bg-white p-5 rounded-2xl border transition-all duration-200 cursor-pointer shadow-3xs flex flex-col justify-between min-h-[110px]" 
        :class="activeTab === 'OVERDUE' ? 'border-rose-600 bg-rose-50/40 ring-2 ring-rose-500/10' : 'border-slate-200 hover:border-rose-300'"
      >
        <div class="flex items-center justify-between">
          <span class="text-[11px] text-rose-600 font-bold uppercase tracking-wider">Keterlambatan</span>
          <i class="pi pi-clock text-rose-400 text-xs"></i>
        </div>
        <p class="text-2xl font-black text-rose-600 mt-2 m-0 font-mono">{{ overdueReservations.length }}</p>
      </div>

      <!-- Card 5: RETURNED (TAMBAHAN UTAMA) -->
      <div 
        @click="activeTab = 'RETURNED'" 
        class="bg-white p-5 rounded-2xl border transition-all duration-200 cursor-pointer shadow-3xs flex flex-col justify-between min-h-[110px] grid-cols-2" 
        :class="activeTab === 'RETURNED' ? 'border-slate-500 bg-slate-100 ring-2 ring-slate-500/10' : 'border-slate-200 hover:border-slate-400'"
      >
        <div class="flex items-center justify-between">
          <span class="text-[11px] text-slate-500 font-bold uppercase tracking-wider">Selesai Kembali</span>
          <i class="pi pi-check-circle text-slate-400 text-xs"></i>
        </div>
        <p class="text-2xl font-black text-slate-700 mt-2 m-0 font-mono">{{ returnedReservations.length }}</p>
      </div>
    </div>

    <!-- FILTER BAR & SEARCH -->
    <div class="bg-white p-4 rounded-2xl border border-slate-200 shadow-3xs flex flex-col sm:flex-row items-center gap-4 justify-between">
      <div class="relative w-full sm:w-96">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3.5 pointer-events-none text-slate-400">
          <i class="pi pi-search text-sm"></i>
        </span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari nama pelanggan, nota invoice..."
          class="w-full pl-10 pr-4 py-2.5 border border-slate-300 rounded-xl text-sm bg-slate-50/60 focus:bg-white focus:border-indigo-500 outline-none transition-all shadow-3xs"
        />
      </div>

      <button 
        @click="fetchReservationsData"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-4 py-2.5 bg-white hover:bg-slate-50 text-slate-700 text-xs font-bold rounded-xl border border-slate-300 shadow-3xs transition-all cursor-pointer"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Sinkronisasi Manifest
      </button>
    </div>

    <!-- MASTER MANIFEST DATA TABLE -->
    <div class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden relative min-h-[250px]">
      
      <!-- RE-INDEX OVERLAY LOADING -->
      <div v-if="isLoading" class="absolute inset-0 bg-white/80 backdrop-blur-xs z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-3xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-bold text-slate-500">Mengekstrak logistik data rental...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/70 text-slate-500 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">ID / No. Invoice</th>
              <th class="py-4 px-6">Nama Pelanggan</th>
              <th class="py-4 px-6">Rentang Waktu Sewa</th>
              <th class="py-4 px-6 text-center">Status Berkas</th>
              <th class="py-4 px-6 text-right">Pendapatan (DP+Pelunasan+Denda)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="currentReservationsList.length === 0 && !isLoading">
              <td colspan="5" class="py-16 text-center text-slate-400 font-medium italic bg-slate-50/10">
                Tidak ada data reservasi rental yang sesuai dengan kriteria filter saat ini.
              </td>
            </tr>
            <tr v-for="res in currentReservationsList" :key="res.id" class="hover:bg-slate-50/40 transition-colors">
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900 text-xs font-mono">{{ res.invoice_number || 'REG-RENTAL' }}</div>
                <div class="text-[10px] text-slate-400 font-mono select-all mt-0.5">{{ res.id }}</div>
              </td>
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900 text-sm">{{ res.customer_name }}</div>
                <div class="text-xs text-slate-500 font-medium mt-0.5 flex items-center gap-1">
                  <i class="pi pi-phone text-[10px]"></i> {{ res.customer_phone || '-' }}
                </div>
              </td>
              <td class="py-4 px-6 text-xs text-slate-600 font-medium">
                <div class="flex items-center gap-1"><span class="w-1.5 h-1.5 rounded-full bg-indigo-400"></span> Mulai: {{ formatDate(res.start_date) }}</div>
                <div class="flex items-center gap-1 mt-1"><span class="w-1.5 h-1.5 rounded-full bg-rose-400"></span> Akhir: {{ formatDate(res.end_date) }}</div>
              </td>
              <td class="py-4 px-6 text-center">
                <span 
                  class="px-2.5 py-1 text-[10px] font-extrabold rounded-lg uppercase border shadow-3xs"
                  :class="getStatusBadgeClass(res.status)"
                >
                  {{ res.status || activeTab }}
                </span>
              </td>
              <td class="py-4 px-6 text-right font-mono font-bold text-slate-900">
                {{ formatCurrency(res.grand_total_income || 0) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>
</template>