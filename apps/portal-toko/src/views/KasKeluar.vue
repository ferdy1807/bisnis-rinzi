<template>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-error">Kas Keluar (Expenses)</h2>
        <p class="text-body-md text-on-surface-variant">
          Catat pengeluaran dari laci kasir (contoh: bayar listrik, kebersihan, restock).
        </p>
      </div>
      <button
        @click="openCreateModal"
        class="flex items-center gap-2 bg-error text-on-error px-4 py-2 rounded-xl shadow-md hover:bg-error/90 transition-colors">
        <span class="material-symbols-outlined">remove</span>
        <span class="text-label-md font-bold">Catat Pengeluaran</span>
      </button>
    </div>

    <!-- Tabel -->
    <div
      class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]"
    >
      <div
        class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex justify-between items-center"
      >
        <h3 class="font-headline-md text-error flex items-center gap-2">
          <span class="material-symbols-outlined">receipt_long</span> Daftar Pengeluaran
        </h3>
      </div>

      <div class="overflow-x-auto flex-1">
        <table class="w-full text-left">
          <thead>
            <tr
              class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
            >
              <th class="px-6 py-4">Waktu</th>
              <th class="px-6 py-4">Kategori Pengeluaran</th>
              <th class="px-6 py-4">Keterangan Tambahan</th>
              <th class="px-6 py-4 text-right">Nominal (Rp)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/20">
            <tr v-if="isLoading" class="text-center">
              <td colspan="4" class="py-8 text-on-surface-variant">Memuat data kas keluar...</td>
            </tr>
            <tr v-else-if="expenses.length === 0" class="text-center">
              <td colspan="4" class="py-8 text-on-surface-variant">
                Belum ada catatan pengeluaran.
              </td>
            </tr>
            <tr
              v-for="exp in expenses"
              :key="exp.id"
              class="hover:bg-surface-container-lowest/50 transition-colors">
              <td class="px-6 py-4 text-body-sm">
                {{
                  exp.created_at
                    ? new Date(exp.created_at).toLocaleString('id-ID')
                    : exp.expense_date
                }}
              </td>
              <td class="px-6 py-4 font-medium">{{ getCategoryName(exp.category_id) }}</td>
              <td class="px-6 py-4 text-body-sm text-on-surface-variant">
                {{ exp.description || '-' }}
              </td>
              <td class="px-6 py-4 text-right font-bold text-error">
                - Rp {{ exp.amount?.toLocaleString('id-ID') || 0 }}
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
      <div
        class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl transform transition-all border-t-4 border-error"
      >
        <h3 class="text-title-lg font-bold text-error mb-4">Catat Kas Keluar</h3>

        <form @submit.prevent="saveExpense" class="space-y-4">
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1"
              >Kategori Pengeluaran</label
            >
            <select
              v-model="formData.category_id"
              required
              class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-error focus:ring-1 focus:ring-error outline-none bg-surface-container-lowest"
            >
              <option value="" disabled>Pilih Kategori</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }} ({{ cat.code }})
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
              class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-error focus:ring-1 focus:ring-error outline-none bg-surface-container-lowest"
              placeholder="0"
            />
          </div>
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1">Keterangan / Rincian</label>
            <textarea
              v-model="formData.description"
              rows="2"
              class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-error focus:ring-1 focus:ring-error outline-none bg-surface-container-lowest"
              placeholder="Contoh: Beli sapu lidi 2 buah"
            ></textarea>
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
              class="px-4 py-2 text-label-md font-medium bg-error text-on-error rounded-lg shadow hover:bg-error/90 transition-colors disabled:opacity-70"
            >
              {{ isSaving ? 'Menyimpan...' : 'Simpan Pengeluaran' }}
            </button>
          </div>
        </form>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { cashApi } from '@frontend/api/cash'
import type { Expense, ExpenseCategory, CashierSession } from '@frontend/types/cash'

const expenses = ref<Expense[]>([])
const categories = ref<ExpenseCategory[]>([])
const isLoading = ref(true)

const showModal = ref(false)
const isSaving = ref(false)
const formData = ref({ category_id: '', amount: 0, description: '' })

const fetchInitialData = async () => {
  isLoading.value = true
  try {
    const session = await cashApi.getCurrentSession()
    if (session) {
      const [expRes, catRes, txRes] = await Promise.all([
        cashApi.getExpenses(),
        cashApi.getExpenseCategories(),
        cashApi.getTransactions(),
      ])

      // Filter expenses dengan mencari transaksi terkait di sesi saat ini
      const sessionTxIds = new Set(
        txRes
          .filter((t) => t.session_id === session.id && t.reference_type === 'EXPENSE')
          .map((t) => t.reference_id),
      )

      expenses.value = expRes.filter((e) => sessionTxIds.has(e.id))
      categories.value = catRes
    } else {
      expenses.value = []
      categories.value = await cashApi.getExpenseCategories()
    }
  } catch (error) {
    console.error('Gagal mengambil data kas keluar', error)
  } finally {
    isLoading.value = false
  }
}

const getCategoryName = (categoryId: string) => {
  const cat = categories.value.find((c) => c.id === categoryId)
  return cat ? cat.name : 'Kategori Tidak Dikenal'
}

const openCreateModal = () => {
  formData.value = { category_id: '', amount: 0, description: '' }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveExpense = async () => {
  isSaving.value = true
  try {
    const currentSession = await cashApi.getCurrentSession()
    if (!currentSession) {
      alert('Anda harus membuka Shift Kasir terlebih dahulu sebelum mencatat pengeluaran.')
      isSaving.value = false
      return
    }

    await cashApi.createExpense({
      category_id: formData.value.category_id,
      amount: formData.value.amount,
      description: formData.value.description,
    })

    closeModal()
    // Refresh expenses
    await fetchInitialData()
  } catch (error: any) {
    console.error('Gagal menyimpan kas keluar', error)
    alert(error.response?.data?.message || 'Terjadi kesalahan saat menyimpan data pengeluaran.')
  } finally {
    isSaving.value = false
  }
}

onMounted(() => {
  fetchInitialData()
})
</script>
