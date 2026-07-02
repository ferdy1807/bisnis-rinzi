<!-- apps/admin-dashboard/src/views/CoaManagementView.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { financeApi } from '@frontend/api/finance';

// --- STATE MANAGEMENT ---
const accounts = ref<any[]>([]);
const isLoading = ref(false);
const searchQuery = ref('');
const filterClassification = ref<string>('ALL');
const errorMessage = ref('');

// Klasifikasi Akun Akuntansi Standar GL
const classifications = [
  { id: 'ALL', label: 'Semua Klasifikasi' },
  { id: 'ASSET', label: 'Aset / Aktiva' },
  { id: 'LIABILITY', label: 'Liabilitas / Kewajiban' },
  { id: 'EQUITY', label: 'Ekuitas / Modal' },
  { id: 'REVENUE', label: 'Pendapatan Operasional' },
  { id: 'EXPENSE', label: 'Beban / Biaya Operasional' }
];

// --- FETCH DATA COA MASTER ---
const fetchCoaAccounts = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    // Memanggil API getAccounts dari FinanceApi untuk memuat data COA ledger resmi
    accounts.value = await financeApi.getAccounts();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat arsitektur Bagan Akun (COA) pusat.';
    // Fallback data struktural standar akuntansi jika response server kosong
    accounts.value = [
      { id: 'acc-101', account_code: '11010', account_name: 'Kas Utama Toko (Brankas)', account_group: 'ASSET', current_balance: 0 },
      { id: 'acc-102', account_code: '11020', account_name: 'Bank BCA Arus Operasional', account_group: 'ASSET', current_balance: 0 },
      { id: 'acc-201', account_code: '21010', account_name: 'Utang Dagang Konsinyasi Supplier', account_group: 'LIABILITY', current_balance: 0 },
      { id: 'acc-301', account_code: '31010', account_name: 'Modal Saham Disetor Owner', account_group: 'EQUITY', current_balance: 0 },
      { id: 'acc-401', account_code: '41010', account_name: 'Pendapatan Penjualan POS Ritel', account_group: 'REVENUE', current_balance: 0 },
      { id: 'acc-402', account_code: '41020', account_name: 'Pendapatan Sewa Unit Kontrak Rental', account_group: 'REVENUE', current_balance: 0 },
      { id: 'acc-501', account_code: '51010', account_name: 'Beban Pokok Penjualan (HPP Ritel)', account_group: 'EXPENSE', current_balance: 0 },
      { id: 'acc-502', account_code: '51020', account_name: 'Beban Kerusakan Aset Operasional', account_group: 'EXPENSE', current_balance: 0 }
    ];
  } finally {
    isLoading.value = false;
  }
};

// --- FILTER & SEARCH LOGIC ---
const filteredAccounts = computed(() => {
  let result = accounts.value;

  // Filter Berdasarkan Tipe Klasifikasi Akuntansi
  if (filterClassification.value !== 'ALL') {
    result = result.filter(acc => acc.account_group === filterClassification.value);
  }

  // Filter Berdasarkan Kata Kunci Pencarian Kode atau Nama Akun
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(acc => 
      acc.account_name?.toLowerCase().includes(query) ||
      acc.account_code?.toLowerCase().includes(query) ||
      acc.id?.toLowerCase().includes(query)
    );
  }

  // Urutkan berdasarkan Kode Akun (Ascending) untuk kerapian Ledger
  return result.sort((a, b) => (a.account_code || '').localeCompare(b.account_code || ''));
});

// Helper Badge Visual Klasifikasi
const getClassificationClass = (type: string) => {
  switch (type) {
    case 'ASSET': return 'bg-emerald-50 text-emerald-700 border-emerald-200';
    case 'LIABILITY': return 'bg-amber-50 text-amber-700 border-amber-200';
    case 'EQUITY': return 'bg-indigo-50 text-indigo-700 border-indigo-200';
    case 'REVENUE': return 'bg-blue-50 text-blue-700 border-blue-200';
    case 'EXPENSE': return 'bg-rose-50 text-rose-700 border-rose-200';
    default: return 'bg-slate-50 text-slate-700 border-slate-200';
  }
};

onMounted(() => {
  fetchCoaAccounts();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <!-- Header Modul COA -->
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-account-tree text-indigo-600"></i>
        Bagan Akun Master (Chart of Accounts - COA)
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Arsitektur pos pembukuan induk. Seluruh transaksi kasir toko (POS, pengeluaran, retail) dan kontrak pegawai (rental) dijurnalkan secara otomatis ke dalam akun-akun resmi ini.
      </p>
    </div>

    <!-- Alert Panel Eror -->
    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>

    <!-- TOOLBAR MANAJEMEN FILTER ADVANCED -->
    <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col md:flex-row items-center gap-4 justify-between">
      <div class="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
        <!-- Input Kata Kunci -->
        <div class="relative w-full sm:w-72">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
            <i class="pi pi-search text-sm"></i>
          </span>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari kode atau nama pos akun..."
            class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
          />
        </div>

        <!-- Dropdown Filter Klasifikasi Standar General Ledger -->
        <div class="relative w-full sm:w-56">
          <select 
            v-model="filterClassification"
            class="w-full pl-3 pr-10 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 appearance-none font-semibold text-slate-700"
          >
            <option v-for="cls in classifications" :key="cls.id" :value="cls.id">
              {{ cls.label }}
            </option>
          </select>
          <span class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-slate-400">
            <i class="pi pi-chevron-down text-xs"></i>
          </span>
        </div>
      </div>

      <button 
        @click="fetchCoaAccounts"
        class="w-full md:w-auto inline-flex items-center justify-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border transition-colors"
        :disabled="isLoading"
      >
        <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
        Segarkan Struktur COA
      </button>
    </div>

    <!-- TABEL LOG DATA UTAMA BAGAN AKUN -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Menyinkronkan integritas neraca COA...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">ID Sistem</th>
              <th class="py-4 px-6">Nomor Kode Akun</th>
              <th class="py-4 px-6">Nama Rekening Buku Besar</th>
              <th class="py-4 px-6 text-center">Klasifikasi Neraca</th>
              <th class="py-4 px-6 text-right">Saldo Pengarah Berjalan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="filteredAccounts.length === 0 && !isLoading">
              <td colspan="5" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada data bagan akun (COA) yang memenuhi filter klasifikasi terpilih.
              </td>
            </tr>
            <tr v-for="account in filteredAccounts" :key="account.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6 font-mono text-xs text-slate-400 select-all">
                {{ account.id }}
              </td>
              <td class="py-4 px-6">
                <span class="px-2.5 py-1 bg-slate-900 text-slate-100 font-mono text-xs font-bold rounded-md">
                  {{ account.account_code }}
                </span>
              </td>
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900">{{ account.account_name }}</div>
              </td>
              <td class="py-4 px-6 text-center">
                <span 
                  class="px-2.5 py-1 text-[10px] font-extrabold rounded-md uppercase border tracking-wider"
                  :class="getClassificationClass(account.account_group)"
                >
                  {{ account.account_group || 'UNKNOWN' }}
                </span>
              </td>
              <td class="py-4 px-6 text-right font-mono font-bold text-slate-400">
                <!-- Nominal dinonaktifkan / disembunyikan pada manajemen COA reguler, balance ditarik via Buku Besar -->
                <span class="text-xs italic font-normal text-slate-400">Automated Ledger Link</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>