<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { rentalApi } from '@frontend/api/rental';
import type { RentalReservation } from '@frontend/types/rental';

// --- STATE MANAGEMENT ---
const reservations = ref<RentalReservation[]>([]);
const isLoading = ref(false);
const isSubmitLoading = ref(false);
const isModalOpen = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

// Filter & Selection State
const selectedReservationId = ref('');
const searchQuery = ref('');

// Form State untuk titipan barang yang akan dirias di dalam box hantaran
const form = ref({
  itemName: '',
  description: '',
  quantity: 1,
  conditionNotes: ''
});

// --- FETCH DATA ---
const fetchReservations = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    // Menarik data reservasi dari rental_db
    reservations.value = await rentalApi.getAllReservations();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat manifes data barang titipan hantaran.';
  } finally {
    isLoading.value = false;
  }
};

// --- SUBMIT HANTARAN ITEMS ---
const handleSaveDeposit = async () => {
  if (!selectedReservationId.value) {
    errorMessage.value = 'Silakan pilih nomor kontrak/invoice terlebih dahulu.';
    return;
  }
  if (!form.value.itemName.trim()) {
    errorMessage.value = 'Nama barang titipan wajib diisi.';
    return;
  }
  if (form.value.quantity <= 0) {
    errorMessage.value = 'Jumlah kuantitas minimal adalah 1.';
    return;
  }

  isSubmitLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    // Memanggil API yang menyimpan item_name, quantity, dan condition_notes ke skema rental_reservation_contents
    await rentalApi.saveDepositItems(selectedReservationId.value, {
      item_name: form.value.itemName,
      description: form.value.description,
      quantity: form.value.quantity,
      condition_notes: form.value.conditionNotes
    });

    successMessage.value = 'Data barang titipan hantaran konsumen berhasil dicatat ke dalam sistem.';
    isModalOpen.value = false;
    
    // Reset form state
    form.value = { itemName: '', description: '', quantity: 1, conditionNotes: '' };
    await fetchReservations();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal menyimpan data barang titipan baru.';
  } finally {
    isSubmitLoading.value = false;
  }
};

const openDepositModal = (resId: string = '') => {
  selectedReservationId.value = resId;
  form.value = { itemName: '', description: '', quantity: 1, conditionNotes: '' };
  errorMessage.value = '';
  successMessage.value = '';
  isModalOpen.value = true;
};

// --- SEARCH & FILTER LOGIC ---
const filteredReservations = computed(() => {
  if (!searchQuery.value) return reservations.value;
  const query = searchQuery.value.toLowerCase();
  return reservations.value.filter(res => 
    res.customer_name?.toLowerCase().includes(query) ||
    res.invoice_number?.toLowerCase().includes(query) ||
    res.id?.toLowerCase().includes(query)
  );
});

// Helper style untuk status kontrak riil (BOOKED, PICKED_UP, RETURNED)
const getStatusBadgeClass = (status: string) => {
  switch (status?.toUpperCase()) {
    case 'BOOKED':
      return 'bg-amber-50 border-amber-200 text-amber-700';
    case 'PICKED_UP':
      return 'bg-indigo-50 border-indigo-200 text-indigo-700';
    case 'RETURNED':
      return 'bg-emerald-50 border-emerald-200 text-emerald-700';
    default:
      return 'bg-slate-50 border-slate-200 text-slate-700';
  }
};

onMounted(() => {
  fetchReservations();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
          <i class="pi pi-box text-indigo-600"></i>
          Pencatatan Barang Titipan Hantaran Konsumen
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Otoritas pendataan, monitoring, dan verifikasi barang-barang milik konsumen yang dititipkan untuk dihias/dirias ke dalam box hantaran sesuai manifestasi item persewaan.
        </p>
      </div>
      <button 
        @click="openDepositModal('')"
        class="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-semibold px-4 py-2.5 rounded-xl shadow-md transition-all duration-200"
      >
        <i class="pi pi-plus text-xs"></i>
        Terima Titipan Barang Baru
      </button>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl text-sm font-medium">
      {{ successMessage }}
    </div>

    <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center gap-3 justify-between">
      <div class="relative w-full sm:w-80">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <i class="pi pi-search text-sm"></i>
        </span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari nama customer atau no invoice..."
          class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
        />
      </div>

      <button 
        @click="fetchReservations"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border transition-colors"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Segarkan Data
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Menyinkronkan data manifestasi hantaran...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">Nomor Kontrak / Invoice</th>
              <th class="py-4 px-6">Pelanggan (Customer)</th>
              <th class="py-4 px-6">Status Kontrak Sewa</th>
              <th class="py-4 px-6">Manifes Isi Barang Titipan Hantaran</th>
              <th class="py-4 px-6 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="filteredReservations.length === 0 && !isLoading">
              <td colspan="5" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada data kontrak persewaan aktif yang ditemukan.
              </td>
            </tr>
            <tr v-for="res in filteredReservations" :key="res.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900">{{ res.invoice_number || 'REG-RENTAL' }}</div>
                <div class="text-[10px] text-slate-400 font-mono select-all">{{ res.id }}</div>
              </td>
              <td class="py-4 px-6">
                <div class="font-semibold text-slate-900">{{ res.customer_name }}</div>
                <div class="text-xs text-slate-400 font-normal">Telp: {{ res.customer_phone || '-' }}</div>
              </td>
              <td class="py-4 px-6">
                <span 
                  :class="getStatusBadgeClass(res.status)"
                  class="px-2 py-0.5 text-[10px] font-extrabold border rounded-md uppercase"
                >
                  {{ res.status || 'UNKNOWN' }}
                </span>
              </td>
              <td class="py-4 px-6 max-w-xs">
                <div v-if="!res.contents || res.contents.length === 0" class="text-xs text-amber-600 font-normal flex items-center gap-1">
                  <i class="pi pi-exclamation-circle"></i> Belum ada isi barang titipan dirias
                </div>
                <div v-else class="space-y-1">
                  <div v-for="(item, idx) in res.contents" :key="idx" class="text-xs bg-slate-50 p-1.5 border rounded border-slate-200">
                    <p class="font-bold text-slate-800">{{ (item as any).item_name || (item as any).ItemName }} ({{ (item as any).quantity || (item as any).Quantity }}x)</p>
                    <p class="text-[10px] text-slate-500 font-normal mt-0.5" v-if="(item as any).condition_notes || (item as any).ConditionNotes">Kondisi: {{ (item as any).condition_notes || (item as any).ConditionNotes }}</p>
                  </div>
                </div>
              </td>
              <td class="py-4 px-6 text-right">
                <button 
                  @click="openDepositModal(res.id)"
                  class="inline-flex items-center gap-1.5 text-xs text-indigo-600 hover:text-indigo-700 bg-indigo-50 hover:bg-indigo-100 font-bold px-3 py-1.5 rounded-lg transition-all"
                >
                  <i class="pi pi-plus text-[9px]"></i> Tambah Barang
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-200 overflow-hidden transform transition-all">
        <div class="p-6 border-b border-slate-100 flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-900">Registrasi Isi Barang Titipan</h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600"><i class="pi pi-times"></i></button>
        </div>
        
        <form @submit.prevent="handleSaveDeposit" class="p-6 space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Nomor Invoice / Pelanggan</label>
            <select 
              v-model="selectedReservationId"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 transition-all font-semibold"
            >
              <option value="" disabled>-- Pilih Invoice Kontrak --</option>
              <option v-for="r in reservations" :key="r.id" :value="r.id">
                {{ r.invoice_number || 'REG-RENTAL' }} - {{ r.customer_name }}
              </option>
            </select>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Barang Titipan (Untuk Dirias)</label>
            <input 
              v-model="form.itemName"
              type="text"
              required
              placeholder="Contoh: Kain Kebaya, Sepasang Sepatu, Set Kosmetik"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
            />
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div class="space-y-1.5 col-span-1">
              <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Jumlah (Qty)</label>
              <input 
                v-model.number="form.quantity"
                type="number"
                required
                min="1"
                class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 font-mono text-center"
              />
            </div>
            <div class="space-y-1.5 col-span-2">
              <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Keterangan Tambahan</label>
              <input 
                v-model="form.description"
                type="text"
                placeholder="Contoh: Warna merah marun / Merk Chanel"
                class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all"
              />
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Catatan Kondisi Fisik Barang Titipan</label>
            <textarea 
              v-model="form.conditionNotes"
              rows="2"
              placeholder="Contoh: Baru / Segel terbuka, Kotak luar ada penyok sedikit..."
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all resize-none"
            ></textarea>
          </div>

          <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-100 mt-6">
            <button type="button" @click="isModalOpen = false" class="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 rounded-lg border border-slate-200 transition-colors">
              Batal
            </button>
            <button 
              type="submit" 
              :disabled="isSubmitLoading"
              class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-md transition-colors disabled:opacity-50"
            >
              <span>{{ isSubmitLoading ? 'Menyimpan...' : 'Simpan Barang' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>