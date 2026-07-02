<template>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary">Laporan Tutup Shift</h2>
        <p class="text-body-md text-on-surface-variant">
          Daftar riwayat shift kasir yang telah selesai (ditutup).
        </p>
      </div>
      <router-link
        to="/daily-report"
        class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2 rounded-xl shadow-md hover:bg-primary/90 transition-colors"
      >
        <span class="material-symbols-outlined">today</span>
        <span class="text-label-md font-bold">Lihat Shift Berjalan</span>
      </router-link>
    </div>

    <!-- Tabel -->
    <div
      class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]"
    >
      <div
        class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col sm:flex-row justify-between items-center gap-4"
      >
        <h3 class="font-headline-md text-primary flex items-center gap-2">
          <span class="material-symbols-outlined">history</span> Riwayat Shift
        </h3>
      </div>

      <div class="overflow-x-auto flex-1">
        <table class="w-full text-left">
          <thead>
            <tr
              class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
            >
              <th class="px-6 py-4">Waktu Buka</th>
              <th class="px-6 py-4">Waktu Tutup</th>
              <th class="px-6 py-4 text-right">Saldo Awal</th>
              <th class="px-6 py-4 text-right">Saldo Akhir Aktual</th>
              <th class="px-6 py-4 text-right">Selisih (Variance)</th>
              <th class="px-6 py-4 text-center">Status</th>
              <th class="px-6 py-4 text-center">Preview Laporan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/20">
            <tr v-if="isLoading" class="text-center">
              <td colspan="6" class="py-8 text-on-surface-variant">Memuat riwayat shift...</td>
            </tr>
            <tr v-else-if="closedShifts.length === 0" class="text-center">
              <td colspan="6" class="py-8 text-on-surface-variant">
                Belum ada riwayat shift yang ditutup.
              </td>
            </tr>
            <tr
              v-for="shift in closedShifts"
              :key="shift.id"
              class="hover:bg-surface-container-lowest/50 transition-colors">
              <td class="px-6 py-4 text-body-sm">
                {{ new Date(shift.open_time).toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 text-body-sm">
                {{ shift.close_time ? new Date(shift.close_time).toLocaleString('id-ID') : '-' }}
              </td>
              <td class="px-6 py-4 text-right">
                Rp {{ (shift.opening_cash || 0).toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 text-right font-bold text-primary">
                Rp {{ (shift.actual_cash || 0).toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 text-right">
                <span
                  :class="[
                    'px-2 py-1 rounded-md text-[10px] font-bold',
                    (shift.difference || 0) < 0
                      ? 'bg-error/10 text-error'
                      : (shift.difference || 0) > 0
                        ? 'bg-primary/10 text-primary'
                        : 'bg-surface-variant/20 text-on-surface-variant',]">
                  {{ (shift.difference || 0) < 0 ? '-' : ((shift.difference || 0) > 0 ? '+' : '') }}Rp
                  {{ Math.abs(shift.difference || 0).toLocaleString('id-ID') }}
                  {{ (shift.difference || 0) === 0 ? '(Seimbang)' : ((shift.difference || 0) > 0 ? '(Surplus)' : '(Minus)') }}
                </span>
              </td>
              <td class="px-6 py-4 text-center flex justify-center items-center gap-2">
                <span
                  class="px-2 py-1 rounded-md text-[10px] font-bold bg-error/10 text-error uppercase">{{ shift.status}}</span>
              </td>
              <td class="px-6 py-4 text-center">
                <button
                  @click="handlePreviewPDF(shift.id)"
                  class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-secondary text-on-secondary hover:bg-secondary/80 transition-colors shadow-sm hover:shadow-md"
                >
                  <span class="material-symbols-outlined text-sm">visibility</span>
                  <span class="hidden sm:inline text-label-sm font-bold">Preview</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- PDF Preview Modal -->
    <div v-if="selectedPdfUrl" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-surface rounded-2xl w-full max-w-4xl h-[85vh] flex flex-col overflow-hidden shadow-2xl">
        <div class="p-4 border-b border-outline-variant/30 flex justify-between items-center bg-surface-container-lowest">
          <h3 class="font-headline-sm text-primary flex items-center gap-2">
            <span class="material-symbols-outlined">visibility</span> Preview Laporan Tutup Shift
          </h3>
          <button @click="closePreview" class="p-2 rounded-full hover:bg-error/10 text-error transition-colors">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="flex-1 w-full bg-surface-variant/20 p-2">
          <iframe :src="selectedPdfUrl" class="w-full h-full rounded-xl border border-outline-variant/30 bg-white" frameborder="0"></iframe>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { cashApi } from '@frontend/api/cash'
import type { CashierSession } from '@frontend/types/cash'
import { useAuthStore } from '@frontend/stores/auth'

const authStore = useAuthStore()
const shifts = ref<CashierSession[]>([])
const isLoading = ref(true)
const selectedPdfUrl = ref<string | null>(null)

const openPreview = (url: string) => {
  selectedPdfUrl.value = url
}

const handlePreviewPDF = (shiftId: string) => {
  const shift = shifts.value.find(s => s.id === shiftId)
  if (shift?.receipt_url) {
    openPreview(shift.receipt_url)
  } else {
    alert("PDF laporan belum tersedia untuk shift ini.")
  }
}

const closePreview = () => {
  selectedPdfUrl.value = null
}

// Hanya tampilkan shift yang sudah ditutup (close_time ada isinya) dan milik user login
const closedShifts = computed(() => {
  return shifts.value.filter((s) => {
    const isClosed = s.close_time !== null && s.close_time !== '0001-01-01T00:00:00Z'
    const isOwner = authStore.user?.role === 'OWNER'
    const isMine = s.cashier_id === authStore.user?.id
    return isClosed && (isOwner || isMine)
  })
})

const fetchShifts = async () => {
  isLoading.value = true
  try {
    const data = await cashApi.getShifts()
    shifts.value = data.sort(
      (a, b) => new Date(b.open_time).getTime() - new Date(a.open_time).getTime(),
    )
  } catch (error) {
    console.error('Gagal mengambil riwayat shift', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchShifts()
})
</script>
