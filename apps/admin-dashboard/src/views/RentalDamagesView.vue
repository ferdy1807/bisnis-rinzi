<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { financeApi } from '@frontend/api/finance';

// --- STATE MANAGEMENT ---
const damages = ref<any[]>([]);
const isLoading = ref(false);
const isSubmitLoading = ref(false);
const isModalOpen = ref(false);
const errorMessage = ref('');
const successMessage = ref('');
const searchQuery = ref('');

// Form State Penyelesaian Ganti Rugi (/api/finance/analytics/rental-damages/{id}/settle)
const selectedDamageId = ref<string>('');
const form = ref({
  paymentAction: 'CASH', // CASH, WAIVED
  auditNotes: ''
});

// --- UTILITIES / FORMATTERS ---
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value);
};

// --- FETCH DATA DAMAGES ---
const fetchRentalDamages = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    const res = await financeApi.getRentalDamages();
    // Hanya tampilkan yang kondisinya tidak GOOD (misal: DAMAGED atau LOST)
    damages.value = res.filter((d: any) => d.condition_status !== 'GOOD');
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat log audit klaim kerusakan rental.';
  } finally {
    isLoading.value = false;
  }
};

// --- SETTLE CLAlM ACTION ---
const handleSettleDamage = async () => {
  if (!selectedDamageId.value) return;
  if (!form.value.auditNotes.trim()) {
    errorMessage.value = 'Catatan pemeriksaan audit wajib diisi untuk rilis penyelesaian.';
    return;
  }

  isSubmitLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    // Mengeksekusi settleRentalDamage via FinanceApi menuju endpoint analitik finansial
    await financeApi.settleRentalDamage(selectedDamageId.value, {
      payment_action: form.value.paymentAction,
      audit_notes: form.value.auditNotes
    });

    successMessage.value = 'Klaim ganti rugi kerusakan unit berhasil diselesaikan dan dijurnalkan ke arus kas masuk.';
    isModalOpen.value = false;
    
    // Reset form state & reload data
    selectedDamageId.value = '';
    form.value = { paymentAction: 'CASH', auditNotes: '' };
    await fetchRentalDamages();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal mengeksekusi penyelesaian klaim denda kerusakan.';
  } finally {
    isSubmitLoading.value = false;
  }
};

const openSettleModal = (id: string) => {
  selectedDamageId.value = id;
  form.value = { paymentAction: 'CASH', auditNotes: '' };
  errorMessage.value = '';
  successMessage.value = '';
  isModalOpen.value = true;
};

// --- FILTER SEARCH ---
const filteredDamages = computed(() => {
  if (!searchQuery.value) return damages.value;
  const query = searchQuery.value.toLowerCase();
  return damages.value.filter(d => 
    d.id?.toLowerCase().includes(query) ||
    d.item_name?.toLowerCase().includes(query) ||
    d.customer_name?.toLowerCase().includes(query)
  );
});

onMounted(() => {
  fetchRentalDamages();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-gavel text-indigo-600"></i>
        Audit Kerusakan Rental & Ganti Rugi
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Konsol penegakan denda logistik fisik. Kelola rekaman unit yang kembali dalam kondisi rusak atau cacat untuk pemulihan nilai aset perusahaan.
      </p>
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
          placeholder="Cari ID berkas, nama aset, atau penyewa..."
          class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
        />
      </div>

      <button 
        @click="fetchRentalDamages"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border transition-colors"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Segarkan Manifes Klaim
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Menyinkronkan denda operasional rental...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">ID Berkas</th>
              <th class="py-4 px-6">Informasi Aset & Penyewa</th>
              <th class="py-4 px-6">Kondisi & Catatan Pegawai</th>
              <th class="py-4 px-6 text-center">Status Penyelesaian</th>
              <th class="py-4 px-6 text-right">Nilai Denda Klaim</th>
              <th class="py-4 px-6 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="filteredDamages.length === 0 && !isLoading">
              <td colspan="6" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada data rekam jejak kerusakan fisik unit rental yang perlu diaudit.
              </td>
            </tr>
            <tr v-for="damage in filteredDamages" :key="damage.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6 font-mono text-xs text-slate-400 select-all">
                {{ damage.id }}
              </td>
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900">{{ damage.item_name || 'Unit Armada Rental' }}</div>
                <div class="text-[11px] text-slate-500 font-normal mt-0.5">Penyewa: {{ damage.customer_name || 'Kontrak Regular' }}</div>
              </td>
              <td class="py-4 px-6 max-w-xs font-normal text-xs text-slate-600 leading-normal">
                <div class="font-semibold text-rose-700">{{ damage.condition_status || 'RUSAK' }}</div>
                <div class="mt-0.5 truncate" :title="damage.condition_notes">{{ damage.condition_notes || 'Tidak ada catatan kondisi khusus.' }}</div>
              </td>
              <td class="py-4 px-6 text-center">
                <span 
                  class="px-2.5 py-1 text-[10px] font-extrabold rounded-md border tracking-wider uppercase"
                  :class="damage.status === 'SETTLED' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-rose-50 text-rose-700 border-rose-200'"
                >
                  {{ damage.status || 'PENDING_AUDIT' }}
                </span>
              </td>
              <td class="py-4 px-6 text-right font-mono font-bold text-rose-600">
                {{ formatCurrency(damage.damage_fee || 0) }}
              </td>
              <td class="py-4 px-6 text-right">
                <button 
                  v-if="damage.status !== 'SETTLED'"
                  @click="openSettleModal(damage.id)"
                  class="inline-flex items-center gap-1 text-xs bg-indigo-600 hover:bg-indigo-700 text-white font-bold px-3 py-1.5 rounded-lg shadow-sm transition-all"
                >
                  <i class="pi pi-check text-[9px]"></i> Settle / Eksekusi
                </button>
                <span v-else class="text-xs text-slate-400 font-normal italic flex items-center justify-end gap-1">
                  <i class="pi pi-check-circle text-emerald-500"></i> Terpembukukan
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-200 overflow-hidden transform transition-all animate-fade-in">
        <div class="p-6 border-b border-slate-100 flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-900">Penyelesaian Klaim Ganti Rugi Aset</h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600"><i class="pi pi-times"></i></button>
        </div>
        
        <form @submit.prevent="handleSettleDamage" class="p-6 space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Aksi Tindakan Finansial</label>
            <select 
              v-model="form.paymentAction"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 transition-all font-semibold text-slate-700"
            >
              <option value="CASH">CASH (Bayar tunai langsung di kasir)</option>
              <option value="WAIVED">WAIVED (Dianulir / Dibebaskan oleh Owner)</option>
            </select>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Catatan Pemeriksaan & Rekonsiliasi Owner</label>
            <textarea 
              v-model="form.auditNotes"
              rows="3"
              required
              placeholder="Contoh: Pembayaran tunai lunas diterima untuk penggantian sparepart unit, status aset dikembalikan ke antrean maintenance..."
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
              <span>{{ isSubmitLoading ? 'Memproses Jurnal...' : 'Selesaikan Klaim Denda' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>