<template>

    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary">Laporan Harian (Shift)</h2>
        <p class="text-body-md text-on-surface-variant">
          Ringkasan performa finansial pada shift yang sedang berjalan saat ini.
        </p>
      </div>
    </div>

    <!-- Konten -->
    <div v-if="isLoading" class="flex justify-center items-center py-20">
      <span class="animate-spin material-symbols-outlined text-4xl text-primary">sync</span>
    </div>
    <div
      v-else-if="!currentSession"
      class="bg-surface rounded-2xl border border-outline-variant/30 p-12 text-center shadow-sm"
    >
      <div
        class="w-16 h-16 bg-error/10 text-error rounded-full flex items-center justify-center mx-auto mb-4"
      >
        <span class="material-symbols-outlined text-3xl">lock</span>
      </div>
      <h3 class="text-headline-sm font-bold mb-2">Shift Belum Dibuka</h3>
      <p class="text-on-surface-variant mb-6">
        Anda harus membuka shift kasir untuk melihat ringkasan performa hari ini.
      </p>
      <router-link
        to="/shifts"
        class="bg-primary text-on-primary px-6 py-3 rounded-xl shadow-md hover:bg-primary/90 transition-colors font-bold text-label-md"
      >
        Buka Shift Sekarang
      </router-link>
    </div>
    <div v-else class="space-y-6">
      <!-- Kartu Ringkasan -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Saldo Awal -->
        <div
          class="bg-surface p-6 rounded-2xl shadow-sm border border-outline-variant/30 relative overflow-hidden"
        >
          <p class="text-label-sm text-on-surface-variant uppercase tracking-wider">
            Saldo Awal Laci
          </p>
          <h3 class="text-headline-md font-bold mt-1">
            Rp {{ (summary?.opening_cash || 0).toLocaleString('id-ID') }}
          </h3>
        </div>
        <!-- Penjualan -->
        <div
          class="bg-surface p-6 rounded-2xl shadow-sm border border-outline-variant/30 relative overflow-hidden"
        >
          <p class="text-label-sm text-on-surface-variant uppercase tracking-wider">
            Total Penjualan
          </p>
          <h3 class="text-headline-md font-bold mt-1 text-primary">
            Rp {{ (summary?.total_income || 0).toLocaleString('id-ID') }}
          </h3>
        </div>
        <!-- Kas Keluar -->
        <div
          class="bg-surface p-6 rounded-2xl shadow-sm border border-outline-variant/30 relative overflow-hidden"
        >
          <p class="text-label-sm text-on-surface-variant uppercase tracking-wider">
            Kas Keluar (Expenses)
          </p>
          <h3 class="text-headline-md font-bold mt-1 text-error">
            - Rp {{ (summary?.total_expense || 0).toLocaleString('id-ID') }}
          </h3>
        </div>
        <!-- Expected Cash -->
        <div class="bg-primary p-6 rounded-2xl shadow-lg relative overflow-hidden text-on-primary">
          <p class="text-label-sm uppercase tracking-wider opacity-90">Estimasi Saldo Akhir</p>
          <h3 class="text-headline-md font-bold mt-1">
            Rp {{ (summary?.expected_cash || 0).toLocaleString('id-ID') }}
          </h3>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { cashApi } from '@frontend/api/cash'
import type { CashierSession, ShiftSummary } from '@frontend/types/cash'

const isLoading = ref(true)
const currentSession = ref<CashierSession | null>(null)
const summary = ref<ShiftSummary | null>(null)

const fetchSummary = async () => {
  isLoading.value = true
  try {
    currentSession.value = await cashApi.getCurrentSession()
    if (currentSession.value) {
      summary.value = await cashApi.getShiftSummary(currentSession.value.id)
    }
  } catch (error) {
    console.error('Gagal mengambil laporan shift', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchSummary()
})
</script>
