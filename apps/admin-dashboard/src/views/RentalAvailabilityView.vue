<!-- apps/admin-dashboard/src/views/RentalAvailabilityView.vue -->
<script setup lang="ts">
import { ref } from 'vue';
import { rentalApi } from '@frontend/api/rental';

// --- STATE MANAGEMENT ---
const isLoading = ref(false);
const errorMessage = ref('');
const searchExecuted = ref(false);

// Form Filter Parameter Sesuai API
const form = ref({
  startDate: new Date().toISOString().split('T')[0],
  endDate: new Date(Date.now() + 86400000).toISOString().split('T')[0] // Besok
});

// Hasil Pengecekan Ketersediaan Kategori Unit
const availabilityResults = ref<{
  product_id: string;
  product_name: string;
  total_units: number;
  rented_units: number;
  available_units: number;
  status: 'AVAILABLE' | 'LIMITED' | 'UNAVAILABLE';
}[]>([]);

// --- UTILITIES / FORMATTERS ---
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

// --- EXECUTE CHECK AVAILABILITY ---
const handleCheckAvailability = async () => {
  if (!form.value.startDate || !form.value.endDate) {
    errorMessage.value = 'Rentang tanggal mulai dan tanggal selesai wajib ditentukan.';
    return;
  }

  if (new Date(form.value.startDate) > new Date(form.value.endDate)) {
    errorMessage.value = 'Tanggal mulai tidak boleh melebihi tanggal selesai.';
    return;
  }

  isLoading.value = true;
  errorMessage.value = '';
  searchExecuted.value = true;

  try {
    const products = await rentalApi.getProducts();
    const reservations = await rentalApi.getAllReservations();
    
    const reqStart = new Date(form.value.startDate).getTime();
    const reqEnd = new Date(form.value.endDate).getTime();

    const rentedMap: Record<string, number> = {};

    for (const res of reservations) {
      if (res.status === 'RETURNED' || res.status === 'CANCELLED') continue;
      
      const resStart = new Date(res.start_date).getTime();
      const resEnd = new Date(res.end_date).getTime();
      
      if (resStart <= reqEnd && resEnd >= reqStart) {
        let items = res.items;
        if (!items || items.length === 0) {
          try {
            const detail = await rentalApi.getReservationDetail(res.id);
            items = detail.items || [];
          } catch (e) {}
        }
        
        items?.forEach(item => {
          rentedMap[item.rental_product_id] = (rentedMap[item.rental_product_id] || 0) + item.qty;
        });
      }
    }

    availabilityResults.value = products.map(p => {
      const total = p.quantity_available || 0;
      const rented = rentedMap[p.id] || 0;
      const available = total - rented;
      
      let status: 'AVAILABLE' | 'LIMITED' | 'UNAVAILABLE' = 'AVAILABLE';
      if (available <= 0) status = 'UNAVAILABLE';
      else if (available <= 2) status = 'LIMITED';
      
      return {
        product_id: p.id,
        product_name: p.name,
        total_units: total,
        rented_units: rented,
        available_units: available < 0 ? 0 : available,
        status
      };
    });
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memeriksa ketersediaan armada unit rental.';
    availabilityResults.value = [];
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <!-- Header Modul -->
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-event-available text-indigo-600"></i>
        Cek Ketersediaan Unit Rental
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Simulasi matriks alokasi aset. Periksa status unit yang bebas tanggungan atau sedang disewa pada rentang waktu tertentu.
      </p>
    </div>

    <!-- Alert Notifikasi Eror -->
    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>

    <!-- PANEL FILTER FILTER PARAMETER JADWAL -->
    <div class="bg-white p-5 rounded-xl border border-slate-200 shadow-sm">
      <form @submit.prevent="handleCheckAvailability" class="grid grid-cols-1 sm:grid-cols-3 gap-4 items-end">
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Tanggal Mulai Sewa</label>
          <input 
            v-model="form.startDate"
            type="date"
            required
            class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all font-semibold text-slate-700"
          />
        </div>

        <div class="space-y-1.5">
          <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Tanggal Rencana Kembali</label>
          <input 
            v-model="form.endDate"
            type="date"
            required
            class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all font-semibold text-slate-700"
          />
        </div>

        <button 
          type="submit"
          :disabled="isLoading"
          class="w-full inline-flex items-center justify-center gap-2 px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-400 text-white text-sm font-bold rounded-lg shadow-md shadow-indigo-600/10 transition-colors cursor-pointer"
        >
          <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-search'" class="text-xs"></i>
          <span>{{ isLoading ? 'Memindai Kalender...' : 'Periksa Ketersediaan' }}</span>
        </button>
      </form>
    </div>

    <!-- TABEL HASIL ANALISIS KETERSEDIAAN UNIT -->
    <div v-if="searchExecuted" class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div class="p-4 bg-slate-50 border-b border-slate-200 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <span class="text-xs font-bold uppercase tracking-wider text-slate-500">
          Hasil Pemindaian Jadwal Kontrak
        </span>
        <span class="text-xs font-medium text-indigo-950 bg-indigo-50 px-2.5 py-1 rounded-md border border-indigo-100">
          Rentang: <span class="font-bold">{{ formatDate(form.startDate) }}</span> s/d <span class="font-bold">{{ formatDate(form.endDate) }}</span>
        </span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">Nama Aset / Produk Rental</th>
              <th class="py-4 px-6 text-center">Total Kapasitas Fisik</th>
              <th class="py-4 px-6 text-center">Sedang Terpakai (Booked)</th>
              <th class="py-4 px-6 text-center">Sisa Unit Bebas (Ready)</th>
              <th class="py-4 px-6 text-right">Status Alokasi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="availabilityResults.length === 0 && !isLoading">
              <td colspan="5" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada data produk rental yang ter-indeks pada sistem logistik gudang.
              </td>
            </tr>
            <tr v-for="item in availabilityResults" :key="item.product_id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900">{{ item.product_name }}</div>
                <div class="text-[10px] text-slate-400 font-mono select-all">Product ID: {{ item.product_id }}</div>
              </td>
              <td class="py-4 px-6 text-center font-mono text-slate-600 font-semibold">
                {{ item.total_units }} Unit
              </td>
              <td class="py-4 px-6 text-center font-mono font-bold text-orange-600">
                {{ item.rented_units }} Unit
              </td>
              <td class="py-4 px-6 text-center font-mono font-extrabold text-emerald-600 bg-emerald-50/20">
                {{ item.available_units }} Unit
              </td>
              <td class="py-4 px-6 text-right">
                <!-- Status Badge secara dinamis berdasarkan kalkulasi sisa unit sewa -->
                <span 
                  class="px-2.5 py-1 text-[10px] font-extrabold rounded-md uppercase border"
                  :class="{
                    'bg-emerald-50 text-emerald-700 border-emerald-200': item.available_units > 2,
                    'bg-amber-50 text-amber-700 border-amber-200': item.available_units <= 2 && item.available_units > 0,
                    'bg-rose-50 text-rose-700 border-rose-200': item.available_units === 0
                  }"
                >
                  {{ item.available_units === 0 ? 'FULL BOOKED' : item.available_units <= 2 ? 'LIMITED STOCK' : 'SAFE / AVAILABLE' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>