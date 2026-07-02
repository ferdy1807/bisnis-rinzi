<!-- apps/admin-dashboard/src/views/ExpenseIncomeManagerView.vue -->
<template>
  <div class="space-y-8">
    <!-- HEADER BAR MODUL -->
    <div
      class="flex flex-col xl:flex-row xl:items-center justify-between gap-6 border-b border-outline-variant/20 pb-6"
    >
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">
          Aturan Finansial: Pengeluaran & Pendapatan
        </h2>
        <p class="text-body-md text-on-surface-variant mt-1">
          Otoritas Owner dalam menetapkan klasifikasi biaya operasional dan pendapatan tambahan
          untuk menjaga integritas laporan laba rugi.
        </p>
      </div>
      <button
        @click="openCreateModal"
        class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2.5 rounded-xl shadow-sm hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all text-label-md font-bold cursor-pointer"
      >
        <span class="material-symbols-outlined text-md">add</span>
        <span>Tambah Aturan Atribut</span>
      </button>
    </div>

    <!-- FEEDBACK ALERT RESPONS -->
    <div
      v-if="errorMessage"
      class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium flex items-center gap-2 animate-fade-in"
    >
      <span class="material-symbols-outlined">error</span>
      <span>{{ errorMessage }}</span>
    </div>
    <div
      v-if="successMessage"
      class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl text-sm font-medium flex items-center gap-2 animate-fade-in"
    >
      <span class="material-symbols-outlined">check_circle</span>
      <span>{{ successMessage }}</span>
    </div>

    <!-- TABS NAVIGATION CONTROL -->
    <div class="flex border-b border-outline-variant/30 gap-6 text-title-small font-bold">
      <button
        @click="activeTab = 'EXPENSE'"
        class="pb-3 border-b-2 transition-all px-2 cursor-pointer border-none bg-transparent"
        :class="
          activeTab === 'EXPENSE'
            ? 'border-primary text-primary font-black'
            : 'border-transparent text-on-surface-variant hover:text-on-surface'
        "
      >
        Kategori Pengeluaran ({{ expenseCategories.length }})
      </button>
      <button
        @click="activeTab = 'INCOME'"
        class="pb-3 border-b-2 transition-all px-2 cursor-pointer border-none bg-transparent"
        :class="
          activeTab === 'INCOME'
            ? 'border-primary text-primary font-black'
            : 'border-transparent text-on-surface-variant hover:text-on-surface'
        "
      >
        Pendapatan Tambahan ({{ incomeCategories.length }})
      </button>
    </div>

    <!-- UTALITAS TOOLBAR PENCARIAN -->
    <div
      class="bg-surface-container-low p-4 rounded-2xl border border-outline-variant/30 flex flex-col sm:flex-row items-center gap-4 justify-between"
    >
      <div class="relative w-full sm:w-80">
        <span
          class="absolute inset-y-0 left-0 flex items-center pl-3 text-on-surface-variant pointer-events-none"
        >
          <span class="material-symbols-outlined text-md">search</span>
        </span>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari nama atau kode klasifikasi..."
          class="w-full pl-9 pr-4 py-2 border border-outline-variant/60 rounded-xl text-xs bg-surface focus:outline-none focus:border-primary transition-all font-medium text-on-surface"
        />
      </div>
      <div class="text-[11px] font-mono text-on-surface-variant">
        Total Terfilter:
        <span class="font-bold text-primary">{{ filteredManifest.length }} Klasifikasi</span>
      </div>
    </div>

    <!-- TABEL UTAMA MANIFES DATA -->
    <div
      class="bg-surface border border-outline-variant/30 rounded-2xl shadow-sm overflow-hidden relative"
    >
      <div
        v-if="isLoading"
        class="absolute inset-0 bg-surface/60 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center"
      >
        <span class="animate-spin material-symbols-outlined text-4xl text-primary">sync</span>
        <span class="text-xs font-mono text-on-surface-variant mt-2"
          >Sinkronisasi aturan master kasir...</span
        >
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-body-sm border-collapse">
          <thead>
            <tr
              class="bg-surface-container-low text-on-surface-variant uppercase text-[9px] tracking-widest font-black border-b border-outline-variant/20"
            >
              <th class="px-6 py-4">Kode Aturan</th>
              <th class="px-6 py-4">Nama Klasifikasi Finansial</th>
              <th class="px-6 py-4 text-center">Tipe Aliran Kas</th>
              <th class="px-6 py-4">ID Entri Sistem</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/10 text-xs font-semibold text-on-surface">
            <tr v-if="filteredManifest.length === 0 && !isLoading">
              <td
                colspan="4"
                class="px-6 py-8 text-center text-on-surface-variant font-normal font-sans italic"
              >
                Tidak ada data aturan klasifikasi transaksi finansial yang terdaftar.
              </td>
            </tr>
            <tr
              v-for="item in filteredManifest"
              :key="item.id"
              class="hover:bg-surface-container-low/40 transition-colors"
            >
              <td class="px-6 py-4 font-mono text-on-surface-variant text-[11px] select-all">
                {{ item.code || '-' }}
              </td>
              <td class="px-6 py-4">
                <span
                  class="font-bold text-sm"
                  :class="item.flow === 'EXPENSE' ? 'text-primary' : 'text-emerald-600'"
                >
                  {{ item?.name || '-' }}
                </span>
              </td>
              <td class="px-6 py-4 text-center">
                <span
                  class="px-2.5 py-0.5 rounded text-[9px] font-black tracking-widest uppercase border"
                  :class="
                    item.flow === 'EXPENSE'
                      ? 'bg-rose-50 text-rose-700 border-rose-200'
                      : 'bg-emerald-50 text-emerald-700 border-emerald-200'
                  "
                >
                  {{ item.flow }}
                </span>
              </td>
              <td class="px-6 py-4 font-mono text-on-surface-variant text-[10px]">
                {{ item.id || 'Tanpa ID Sistem' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- POPUP FORM MODAL (UNIVERSAL CREATE ENTRY) -->
    <div
      v-if="isModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm"
    >
      <div
        class="bg-surface rounded-2xl shadow-xl w-full max-w-md border border-outline-variant/30 overflow-hidden transform transition-all"
      >
        <div
          class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex items-center justify-between"
        >
          <h4 class="font-title-medium text-primary font-bold">
            Daftarkan Aturan
            {{ form.flow === 'EXPENSE' ? 'Pengeluaran' : 'Pendapatan Tambahan' }} Baru
          </h4>
          <button
            @click="isModalOpen = false"
            class="text-on-surface-variant hover:text-on-surface cursor-pointer border-none bg-transparent"
          >
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
          <!-- Switch Pilihan Arus Finansial Secara Fleksibel di Modal -->
          <div class="flex flex-col gap-1.5">
            <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider"
              >Tipe Aliran Transaksi</label
            >
            <div class="grid grid-cols-2 gap-2 bg-slate-100 p-1 rounded-xl border">
              <button
                type="button"
                @click="form.flow = 'EXPENSE'"
                :class="
                  form.flow === 'EXPENSE'
                    ? 'bg-white text-primary shadow-xs font-black'
                    : 'text-slate-500 hover:text-slate-800'
                "
                class="py-1.5 rounded-lg text-xs font-bold transition-all border-none cursor-pointer"
              >
                EXPENSE (Beban)
              </button>
              <button
                type="button"
                @click="form.flow = 'INCOME'"
                :class="
                  form.flow === 'INCOME'
                    ? 'bg-white text-emerald-600 shadow-xs font-black'
                    : 'text-slate-500 hover:text-slate-800'
                "
                class="py-1.5 rounded-lg text-xs font-bold transition-all border-none cursor-pointer"
              >
                INCOME (Luar Usaha)
              </button>
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider"
              >Kode Klasifikasi Kategori</label
            >
            <input
              v-model="form.code"
              type="text"
              required
              :placeholder="
                form.flow === 'EXPENSE'
                  ? 'Contoh: EXP-SAMPAH, EXP-LISTRIK'
                  : 'Contoh: INC-KARTON, INC-BUNGA'
              "
              class="w-full px-3 py-2 border border-outline-variant/60 rounded-xl text-xs bg-surface outline-none focus:border-primary font-medium uppercase font-mono"
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider"
              >Nama / Deskripsi Aturan Jenis Biaya</label
            >
            <textarea
              v-model="form.name"
              rows="3"
              placeholder="Berikan panduan penggunaan nama biaya/pendapatan untuk memudahkan kasir toko..."
              class="w-full px-3 py-2 border border-outline-variant/60 rounded-xl text-xs bg-surface outline-none focus:border-primary font-medium resize-none"
            ></textarea>
          </div>

          <div
            class="pt-4 flex items-center justify-end gap-3 border-t border-outline-variant/20 mt-6"
          >
            <button
              type="button"
              @click="isModalOpen = false"
              class="px-4 py-2 text-xs font-bold border border-outline-variant/60 text-on-surface hover:bg-surface-container-low rounded-xl bg-transparent cursor-pointer"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isLoading"
              class="px-4 py-2 text-xs font-bold bg-primary text-on-primary hover:bg-primary-container hover:text-on-primary-container disabled:bg-slate-400 rounded-xl shadow-sm cursor-pointer"
            >
              <span>{{ isLoading ? 'Memproses...' : 'Simpan Kategori' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { cashApi } from '@frontend/api/cash'
import { useFinanceStore } from '@frontend/stores/finance'
import type { ExpenseCategory } from '@frontend/types/cash'

// --- INITIALIZE SYSTEM STORES ---
const financeStore = useFinanceStore()

// --- STATE MANAGEMENT ---
const activeTab = ref<'EXPENSE' | 'INCOME'>('EXPENSE')
const expenseCategories = ref<ExpenseCategory[]>([])
const incomeCategories = ref<any[]>([])
const isLoading = ref(false)
const searchQuery = ref('')
const isModalOpen = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// Form State komprehensif penampung payload
const form = ref({
  flow: 'EXPENSE' as 'EXPENSE' | 'INCOME',
  code: '',
  name: '',
})

// --- FETCH DATA MASTER CORE ---
const fetchAllMasterCategories = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    // Memanggil API fisik getExpenseCategories serta penambahan data income tambahan
    const [expensesRes, incomeRes] = await Promise.allSettled([
      cashApi.getExpenseCategories(),
      cashApi.getOtherIncomeCategories?.() ?? Promise.resolve([]),
    ])

    if (expensesRes.status === 'fulfilled') {
      expenseCategories.value = expensesRes.value
    }
    if (incomeRes.status === 'fulfilled') {
      // Jika endpoint lain belum diimplementasi sempurna, di-fallback array kosong agar aman
      incomeCategories.value = incomeRes.value || []
    }
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal merekonsiliasi aturan kategori finansial.'
  } finally {
    isLoading.value = false
  }
}

// --- FILTER SEARCH UNIFIED MANIFEST ---
const filteredManifest = computed(() => {
  // 1. Pilih dataset berbasis Tab Aktif
  const activeDataset =
    activeTab.value === 'EXPENSE'
      ? expenseCategories.value.map((e) => ({ ...e, flow: 'EXPENSE' }))
      : incomeCategories.value.map((i) => ({ ...i, flow: 'INCOME' }))

  if (!searchQuery.value) return activeDataset
  const query = searchQuery.value.toLowerCase().trim()

  return activeDataset.filter(
    (item) =>
      item.name.toLowerCase().includes(query) ||
      (item.code && item.code.toLowerCase().includes(query)) ||
      (item.id && item.id.toLowerCase().includes(query)),
  )
})

// --- FORM MODAL ACTIONS ---
const openCreateModal = () => {
  form.value = { flow: activeTab.value, code: '', name: '' }
  errorMessage.value = ''
  successMessage.value = ''
  isModalOpen.value = true
}

// --- SUBMIT WORKFLOW EXECUTION ---
const handleSubmit = async () => {
  if (!form.value.name.trim() || !form.value.code.trim()) {
    errorMessage.value = 'Seluruh data parameter wajib terisi utuh.'
    return
  }

  isLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    if (form.value.flow === 'EXPENSE') {
      // Pemanggilan post ke endpoint /api/cash/expense-categories
      await cashApi.createExpenseCategory({
        code: form.value.code.toUpperCase().trim(),
        name: form.value.name,
      })
      successMessage.value = 'Kategori biaya pengeluaran toko baru berhasil didaftarkan.'
    } else {
      // Mengirim POST menuju sub-handler pendapatan tambahan luar usaha (/api/cash/income-categories)
      await cashApi.createOtherIncomeCategory({
        code: form.value.code.toUpperCase().trim(),
        name: form.value.name,
      })
      successMessage.value = 'Kategori pendapatan tambahan baru berhasil didaftarkan.'
    }

    isModalOpen.value = false

    // Refresh visualisasi tabel lokal & global dashboard
    await fetchAllMasterCategories()
    await financeStore.fetchDashboardSummary()
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal menyinkronkan aturan kategori baru ke server.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchAllMasterCategories()
})
</script>
