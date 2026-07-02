<template>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary">Kas Masuk (Internal Income)</h2>
        <p class="text-body-md text-on-surface-variant">
          Catat pemasukan kasir di luar transaksi penjualan (contoh: suntikan modal laci awal).
        </p>
      </div>
      <button
        @click="openCreateModal"
        class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2 rounded-xl shadow-md hover:bg-primary/90 transition-colors"
      >
        <span class="material-symbols-outlined">add</span>
        <span class="text-label-md font-bold">Catat Kas Masuk</span>
      </button>
    </div>

    <!-- Tabel -->
    <div
      class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]"
    >
      <div
        class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex justify-between items-center"
      >
        <h3 class="font-headline-md text-primary flex items-center gap-2">
          <span class="material-symbols-outlined">payments</span> Daftar Kas Masuk
        </h3>
      </div>

      <div class="overflow-x-auto flex-1">
        <table class="w-full text-left">
          <thead>
            <tr
              class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
            >
              <th class="px-6 py-4">Waktu</th>
              <th class="px-6 py-4">Referensi</th>
              <th class="px-6 py-4">Keterangan</th>
              <th class="px-6 py-4 text-right">Nominal (Rp)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/20">
            <tr v-if="isLoading" class="text-center">
              <td colspan="4" class="py-8 text-on-surface-variant">Memuat data kas masuk...</td>
            </tr>
            <tr v-else-if="incomes.length === 0" class="text-center">
              <td colspan="4" class="py-8 text-on-surface-variant">Belum ada catatan kas masuk.</td>
            </tr>
            <tr
              v-for="inc in incomes"
              :key="inc.id"
              class="hover:bg-surface-container-lowest/50 transition-colors"
            >
              <td class="px-6 py-4 text-body-sm">
                {{ inc.created_at ? new Date(inc.created_at).toLocaleString('id-ID') : '-' }}
              </td>
              <td class="px-6 py-4 font-medium">{{ inc.reference_type }}</td>
              <td class="px-6 py-4 text-body-sm text-on-surface-variant">{{ inc.notes || '-' }}</td>
              <td class="px-6 py-4 text-right font-bold text-primary">
                + Rp {{ inc.amount?.toLocaleString('id-ID') || 0 }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    >
      <div class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl transform transition-all">
        <h3 class="text-title-lg font-bold text-primary mb-4">Catat Kas Masuk</h3>

        <form @submit.prevent="saveIncome" class="space-y-4">
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1">Kategori Pemasukan</label>
            <select
              v-model="formData.title"
              required
              class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-primary focus:ring-1 focus:ring-primary outline-none bg-surface-container-lowest"
            >
              <option value="" disabled>Pilih Kategori Pendapatan</option>
              <option v-for="cat in incomeCategories" :key="cat.id" :value="cat.name">
                [{{ cat.code }}] {{ cat.name }}
              </option>
            </select>
          </div>
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1">Nominal (Rp)</label>
            <input
              v-model.number="formData.amount"
              type="number"
              min="0"
              required
              class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-primary focus:ring-1 focus:ring-primary outline-none bg-surface-container-lowest"
              placeholder="0"
            />
          </div>
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1">Keterangan Tambahan</label>
            <textarea
              v-model="formData.description"
              rows="2"
              class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-primary focus:ring-1 focus:ring-primary outline-none bg-surface-container-lowest"
              placeholder="Catatan opsional..."></textarea>
          </div>

          <div class="flex justify-end gap-3 mt-8">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 text-label-md font-medium text-on-surface-variant hover:bg-surface-variant/20 rounded-lg transition-colors"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSaving"
              class="px-4 py-2 text-label-md font-medium bg-primary text-on-primary rounded-lg shadow hover:bg-primary/90 transition-colors disabled:opacity-70"
            >
              {{ isSaving ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { cashApi } from '@frontend/api/cash'
import type { CashTransaction, CashierSession } from '@frontend/types/cash'

const transactions = ref<CashTransaction[]>([])
const isLoading = ref(true)

const showModal = ref(false)
const isSaving = ref(false)
const formData = ref({ title: '', amount: 0, description: '' })

// Filter hanya transaksi DEPOSIT (Kas Masuk)
const incomes = computed(() => {
  return transactions.value.filter((t) => t.transaction_type === 'DEPOSIT')
})

const fetchTransactions = async () => {
  isLoading.value = true
  try {
    const session = await cashApi.getCurrentSession()
    if (session) {
      const allTx = await cashApi.getTransactions()
      transactions.value = allTx.filter((t) => t.session_id === session.id)
    } else {
      transactions.value = []
    }
  } catch (error) {
    console.error('Gagal mengambil histori transaksi kas', error)
  } finally {
    isLoading.value = false
  }
}

const incomeCategories = ref<any[]>([])

const fetchIncomeCategories = async () => {
  try {
    incomeCategories.value = await cashApi.getOtherIncomeCategories()
  } catch (error) {
    console.error('Gagal mengambil kategori pendapatan', error)
  }
}

const openCreateModal = () => {
  formData.value = { title: '', amount: 0, description: '' }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveIncome = async () => {
  isSaving.value = true
  try {
    const currentSession = await cashApi.getCurrentSession()
    if (!currentSession) {
      alert('Anda harus membuka Shift Kasir terlebih dahulu sebelum mencatat kas masuk.')
      isSaving.value = false
      return
    }

    await cashApi.createTransaction({
      transaction_type: 'DEPOSIT',
      reference_type: 'MANUAL',
      amount: formData.value.amount,
      notes: `${formData.value.title} - ${formData.value.description}`,
    })

    closeModal()
    await fetchTransactions()
  } catch (error: any) {
    console.error('Gagal menyimpan kas masuk', error)
    alert(error.response?.data?.message || 'Terjadi kesalahan saat menyimpan data mutasi kas.')
  } finally {
    isSaving.value = false
  }
}

onMounted(() => {
  fetchTransactions()
  fetchIncomeCategories()
})
</script>
