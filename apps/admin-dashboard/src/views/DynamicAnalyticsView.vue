<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { financeApi } from '@frontend/api/finance'

// --- STATE MANAGEMENT ---
const isLoading = ref(false)
const errorMessage = ref('')
const rawAnalyticsData = ref<any>(null)

// Pilihan aksi dinamis yang didukung oleh backend dispatcher
const selectedAction = ref<string>('performance') // Default action

const availableActions = [
  {
    id: 'performance',
    label: 'Metrik Performa Kerja Staf',
    icon: 'pi-users',
    desc: 'Analisis produktivitas, pencapaian shift, dan efisiensi operasional kru toko.',
  },
  {
    id: 'loss-analytics',
    label: 'Laporan Kerugian & Depresiasi Fisik',
    icon: 'pi-percentage',
    desc: 'Evaluasi nilai penyusutan barang retail dan akumulasi kerugian dari unit rental rusak.',
  },
  {
    id: 'retail-insights',
    label: 'Metrik Khusus Performa Ritel',
    icon: 'pi-shopping-bag',
    desc: 'Statistik perputaran produk, profit margin per kategori, dan tren keranjang POS.',
  },
]

// --- UTILITIES / FORMATTERS ---
const formatJson = (data: any): string => {
  if (!data) return ''
  return JSON.stringify(data, null, 2)
}

// --- FETCH ANALYTICS DATA ---
const fetchDynamicAnalytics = async () => {
  isLoading.value = true
  errorMessage.value = ''
  rawAnalyticsData.value = null

  try {
    // Memanggil API dinamis dispatcher berdasarkan string action yang dipilih Owner
    const response = await financeApi.getAnalyticsData(selectedAction.value)
    rawAnalyticsData.value = response
  } catch (error: any) {
    errorMessage.value =
      error?.message ||
      `Gagal memuat log dispatcher analitik untuk tindakan: ${selectedAction.value}`
  } finally {
    isLoading.value = false
  }
}

// Pemicu pergantian aksi analitik
const handleActionChange = (actionId: string) => {
  selectedAction.value = actionId
  fetchDynamicAnalytics()
}

onMounted(() => {
  fetchDynamicAnalytics()
})
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="border-b border-slate-200 pb-5">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
        <i class="pi pi-sliders-h text-indigo-600"></i>
        Dispatcher Analitik Dinamis
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Pusat panggilan metrik lanjutan. Eksekusi tindakan analitik mikro secara dinamis langsung
        menuju repositori pengolahan data pusat.
      </p>
    </div>

    <div
      v-if="errorMessage"
      class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium animate-fade-in"
    >
      <i class="pi pi-exclamation-circle mr-2"></i>
      <span>{{ errorMessage }}</span>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
      <div class="space-y-3 lg:col-span-1">
        <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1"
          >1. Pilih Parameter Tindakan</span
        >

        <div
          v-for="action in availableActions"
          :key="action.id"
          @click="handleActionChange(action.id)"
          class="bg-white p-4 rounded-xl border cursor-pointer transition-all hover:border-slate-300 shadow-xs flex items-start gap-4"
          :class="
            selectedAction === action.id
              ? 'border-indigo-600 bg-indigo-50/30 ring-2 ring-indigo-500/10'
              : 'border-slate-200'
          "
        >
          <div
            class="p-2.5 rounded-lg shrink-0 border"
            :class="
              selectedAction === action.id
                ? 'bg-indigo-600 text-white border-indigo-600'
                : 'bg-slate-50 text-slate-500 border-slate-200'
            "
          >
            <i :class="['pi', action.icon, 'text-sm']"></i>
          </div>
          <div class="space-y-0.5">
            <h4
              class="text-xs font-bold text-slate-900"
              :class="{ 'text-indigo-900': selectedAction === action.id }"
            >
              {{ action.label }}
            </h4>
            <p class="text-[11px] text-slate-500 font-normal leading-normal">
              {{ action.desc }}
            </p>
          </div>
        </div>

        <button
          @click="fetchDynamicAnalytics"
          :disabled="isLoading"
          class="w-full inline-flex items-center justify-center gap-2 px-4 py-2.5 bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold rounded-xl transition-colors disabled:bg-slate-400 mt-2 shadow-sm"
        >
          <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-[10px]"></i>
          <span>Paksa Segarkan Metrik</span>
        </button>
      </div>

      <div
        class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden lg:col-span-2 flex flex-col h-[520px] relative"
      >
        <div
          v-if="isLoading"
          class="absolute inset-0 bg-white/80 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center"
        >
          <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
          <span class="text-xs font-semibold text-slate-400"
            >Menyusun kompilasi data dispatcher...</span
          >
        </div>

        <div
          class="p-4 bg-slate-50 border-b border-slate-200 flex items-center justify-between shrink-0"
        >
          <span
            class="text-xs font-bold uppercase tracking-wider text-slate-500 flex items-center gap-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-indigo-500"></span>
            Data Stream Object Inspector
          </span>
          <span
            class="font-mono text-[10px] bg-slate-200 text-slate-700 px-2 py-0.5 rounded-md border font-bold"
          >
            action: /api/finance/analytics/{{ selectedAction }}
          </span>
        </div>

        <div
          class="flex-1 p-4 overflow-auto bg-slate-950 text-slate-200 font-mono text-xs selection:bg-indigo-500/30 custom-scrollbar"
        >
          <pre v-if="rawAnalyticsData" class="whitespace-pre-wrap break-all leading-relaxed">{{
            formatJson(rawAnalyticsData)
          }}</pre>

          <div
            v-else-if="!isLoading"
            class="h-full flex flex-col items-center justify-center text-slate-500 font-normal"
          >
            <i class="pi pi-code text-2xl mb-2 text-slate-600"></i>
            <span>Tidak ada data objek stream yang ditangkap atau respons kosong.</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: #020617;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #475569;
}
</style>
