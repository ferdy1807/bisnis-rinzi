<template>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary">Riwayat Transaksi</h2>
        <p class="text-body-md text-on-surface-variant">
          Daftar transaksi penjualan yang terjadi dari mesin kasir (POS).
        </p>
      </div>
    </div>

    <!-- Tabel -->
    <div
      class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]"
    >
      <div
        class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col sm:flex-row justify-between items-center gap-4"
      >
        <h3 class="font-headline-md text-primary flex items-center gap-2">
          <span class="material-symbols-outlined">receipt_long</span> Daftar Transaksi
        </h3>
        <!-- Search bar opsional -->
        <div class="relative w-full sm:w-64">
          <span
            class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-lg"
            >search</span
          >
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari ID Transaksi..."
            class="w-full pl-10 pr-4 py-2 bg-surface border border-outline-variant/50 rounded-lg text-body-sm focus:ring-2 focus:ring-primary focus:outline-none"
          />
        </div>
      </div>

      <div class="overflow-x-auto flex-1">
        <table class="w-full text-left">
          <thead>
            <tr
              class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
            >
              <th class="px-6 py-4">Waktu Transaksi</th>
              <th class="px-6 py-4">No. Invoice</th>
              <th class="px-6 py-4">Metode Bayar</th>
              <th class="px-6 py-4">Total Belanja</th>
              <th class="px-6 py-4">Nominal Bayar</th>
              <th class="px-6 py-4">Kembali</th>
              <th class="px-6 py-4 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/20">
            <tr v-if="isLoading" class="text-center">
              <td colspan="7" class="py-8 text-on-surface-variant">Memuat data transaksi...</td>
            </tr>
            <tr v-else-if="filteredTransactions.length === 0" class="text-center">
              <td colspan="7" class="py-8 text-on-surface-variant">
                Belum ada transaksi ditemukan untuk akun Anda.
              </td>
            </tr>
            <tr
              v-for="t in filteredTransactions"
              :key="t.id"
              class="hover:bg-surface-container-lowest/50 transition-colors"
            >
              <td class="px-6 py-4 text-body-sm">
                {{ new Date(t.transaction_date).toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 font-mono text-sm text-primary">{{ t.invoice_number }}</td>
              <td class="px-6 py-4 text-body-sm capitalize">{{ t.payment_method || 'Tunai' }}</td>
              <td class="px-6 py-4 font-bold text-primary">
                Rp {{ t.total?.toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 text-body-sm">
                Rp {{ (t.amount_paid || t.total)?.toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 text-success font-bold">
                Rp {{ (t.change_amount || 0).toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 text-center">
                <button
                  v-if="t.invoice_url"
                  @click="previewInvoice(t.invoice_url)"
                  class="p-2 text-primary hover:bg-primary/10 rounded-full transition-colors"
                  title="Preview Invoice">
                  <span class="material-symbols-outlined text-sm">visibility</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Preview Invoice -->
    <div
      v-if="previewUrl"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div
        class="bg-surface rounded-2xl shadow-xl w-full max-w-4xl h-[80vh] overflow-hidden flex flex-col"
      >
        <div
          class="p-4 border-b border-outline-variant/30 flex justify-between items-center bg-primary text-on-primary"
        >
          <h3 class="font-headline-sm">Preview Invoice</h3>
          <button
            @click="previewUrl = null"
            class="hover:bg-white/20 rounded-full p-1 transition-colors"
          >
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="flex-1 bg-surface-container-lowest p-2">
          <iframe :src="previewUrl" class="w-full h-full border-none rounded-lg"></iframe>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { posApi } from '@frontend/api/pos'
import type { Sale } from '@frontend/types/pos'
import { useAuthStore } from '@frontend/stores/auth'

const authStore = useAuthStore()

const transactions = ref<Sale[]>([])
const isLoading = ref(true)
const searchQuery = ref('')
const previewUrl = ref<string | null>(null)

const previewInvoice = (url: string) => {
  previewUrl.value = url
}

const fetchTransactions = async () => {
  isLoading.value = true
  try {
    transactions.value = await posApi.getSalesHistory()
  } catch (error) {
    console.error('Gagal mengambil transaksi', error)
  } finally {
    isLoading.value = false
  }
}

const filteredTransactions = computed(() => {
  // Tampilkan semua jika OWNER, jika tidak filter berdasarkan user yang login
  const userSales =
    authStore.user?.role === 'OWNER'
      ? transactions.value
      : transactions.value.filter((t) => t.cashier_id === authStore.user?.id)

  if (!searchQuery.value) return userSales

  const q = searchQuery.value.toLowerCase()
  // Gunakan optional chaining pada invoice_number dan ID sebagai fallback
  return userSales.filter(
    (t) =>
      (t.invoice_number && t.invoice_number.toLowerCase().includes(q)) ||
      (t.id && t.id.toLowerCase().includes(q)),
  )
})

onMounted(() => {
  fetchTransactions()
})
</script>
