<!-- apps/admin-dashboard/src/views/InternalIncomesView.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { cashApi } from '@frontend/api/cash'
import { rentalApi } from '@frontend/api/rental'
import { financeApi } from '@frontend/api/finance'
import { authApi } from '@frontend/api/auth'
import type { CashTransaction } from '@frontend/types/cash'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'

// --- STATE MANAGEMENT ---
const incomes = ref<CashTransaction[]>([])
const isLoading = ref(false)
const searchQuery = ref('')
const selectedDate = ref('')
const errorMessage = ref('')
const isExporting = ref(false)

// --- UTILITIES / FORMATTERS ---
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(value)
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// --- FETCH DATA ---
const fetchInternalIncomes = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const rawIncomes = await financeApi.getDailyIncomes()
    
    // Urutkan dari yang terbaru ke terlama
    rawIncomes.sort((a, b) => new Date(b.created_at || 0).getTime() - new Date(a.created_at || 0).getTime())
    incomes.value = rawIncomes
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat data pendapatan internal.'
  } finally {
    isLoading.value = false
  }
}

// --- FILTER & SEARCH ---
const filteredIncomes = computed(() => {
  let list = incomes.value
  
  if (selectedDate.value) {
    list = list.filter((income) => {
      if (!income.created_at) return false
      const d = new Date(income.created_at)
      const dateStr = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
      return dateStr === selectedDate.value
    })
  }

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    list = list.filter(
      (income) =>
        (income.notes && income.notes.toLowerCase().includes(query)) ||
        (income.created_by && income.created_by.toLowerCase().includes(query)) ||
        (income.id && income.id.toLowerCase().includes(query)),
    )
  }
  
  return list
})

// --- TOTAL ACCUMULATION ---
const totalAmount = computed(() => {
  return filteredIncomes.value.reduce((sum, item) => sum + (item.amount || 0), 0)
})

// --- GROUPING BY DATE ---
const groupedIncomes = computed(() => {
  const groups: { [key: string]: { dateStr: string; dateLabel: string; total: number; items: CashTransaction[] } } = {}
  
  filteredIncomes.value.forEach(income => {
    if (!income.created_at) return
    const d = new Date(income.created_at)
    // Avoid timezone offset issues by using local year/month/date
    const dateStr = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    
    if (!groups[dateStr]) {
      groups[dateStr] = {
        dateStr,
        dateLabel: d.toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }),
        total: 0,
        items: []
      }
    }
    
    groups[dateStr].items.push(income)
    groups[dateStr].total += (income.amount || 0)
  })
  
  return Object.values(groups).sort((a, b) => new Date(b.dateStr).getTime() - new Date(a.dateStr).getTime())
})

const generatePDF = (dateStr: string) => {
  const targetGroup = groupedIncomes.value.find(g => g.dateStr === dateStr)
  if (!targetGroup) return null

  const doc = new jsPDF()
  doc.setFontSize(18)
  doc.text(`Laporan Pendapatan Harian - ${targetGroup.dateLabel}`, 14, 22)
  doc.setFontSize(11)
  
  let yPos = 30
  
  doc.text(`Total Pendapatan: ${formatCurrency(targetGroup.total)}`, 14, yPos)
  yPos += 12

  const tableBody: any[][] = []
  
  tableBody.push([{ content: `Tanggal: ${targetGroup.dateLabel} (Total: ${formatCurrency(targetGroup.total)})`, colSpan: 4, styles: { fillColor: [240, 240, 240], fontStyle: 'bold' } }])
  
  targetGroup.items.forEach(item => {
    tableBody.push([
      formatDate(item.created_at),
      item.notes || '-',
      item.created_by || '-',
      formatCurrency(item.amount || 0)
    ])
  })

  autoTable(doc, {
    startY: yPos,
    head: [['Waktu', 'Catatan', 'Petugas', 'Nominal']],
    body: tableBody,
    theme: 'grid',
    headStyles: { fillColor: [79, 70, 229] }
  })

  return doc
}

// --- EXPORT PDF ---
const exportToPDF = async (dateStr: string) => {
  const doc = generatePDF(dateStr)
  if (!doc) return

  isExporting.value = true
  errorMessage.value = ''
  try {
    const filename = `Laporan_Pendapatan_Harian_${dateStr}.pdf`
    const pdfBlob = doc.output('blob')
    
    // Upload ke bucket MinIO via API
    await financeApi.uploadDailyIncomeReport(dateStr, pdfBlob, filename)
    // alert('PDF berhasil diekspor dan diunggah ke laporan-income-harian!')
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal mengekspor PDF.'
  } finally {
    isExporting.value = false
  }
}

// --- VIEW PDF ---
const viewPDF = (dateStr: string) => {
  const doc = generatePDF(dateStr)
  if (!doc) return
  
  const pdfBlob = doc.output('blob')
  const blobUrl = URL.createObjectURL(pdfBlob)
  window.open(blobUrl, '_blank')
}

onMounted(() => {
  fetchInternalIncomes()
})
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <!-- Header Modul -->
    <div class="border-b border-slate-200 pb-4">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-1">
        <i class="pi pi-download text-emerald-600"></i>
        Transaksi Pendapatan Internal
      </h1>
      <p class="text-sm text-slate-500 mt-1">
        Daftar aliran dana masuk non-operasional retail (contoh: suntikan modal owner, pendapatan
        sewa space iklan vendor, dll) yang tercatat di kasir.
      </p>
    </div>

    <!-- Alert Panel -->
    <div
      v-if="errorMessage"
      class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl flex items-start gap-3">
      <i class="pi pi-exclamation-circle text-lg mt-0.5"></i>
      <span class="text-sm font-medium">{{ errorMessage }}</span>
    </div>

    <!-- Ringkasan Ringkas & Kontrol Pencarian -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Widget Total Akumulasi -->
      <div
        class="bg-gradient-to-br from-emerald-600 to-teal-700 text-white p-5 rounded-xl shadow-md shadow-emerald-600/10 flex flex-col justify-between"
      >
        <span class="text-xs font-bold tracking-wider uppercase opacity-80"
          >Total Pendapatan Terfilter</span
        >
        <span class="text-2xl font-extrabold mt-2 font-mono">{{
          formatCurrency(totalAmount)
        }}</span>
      </div>

      <!-- Bilah Pencarian & Kontrol (Span 2 Kolom di Desktop) -->
      <div
        class="md:col-span-2 bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center gap-4 justify-between"
      >
        <div class="flex-1 flex flex-col sm:flex-row gap-3 w-full">
          <div class="relative w-full sm:w-2/3">
            <span
              class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400 pointer-events-none"
            >
              <i class="pi pi-search text-sm"></i>
            </span>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari berdasarkan catatan, penerima..."
              class="w-full pl-9 pr-4 py-2.5 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none transition-all"
            />
          </div>
          <div class="relative w-full sm:w-1/3">
            <input
              v-model="selectedDate"
              type="date"
              class="w-full px-4 py-2.5 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none transition-all text-slate-600"
            />
          </div>
        </div>
        <div class="flex items-center gap-2 w-full sm:w-auto">
          <button
            @click="fetchInternalIncomes"
            class="flex-1 sm:flex-none inline-flex items-center justify-center gap-2 px-4 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-lg border border-slate-200 transition-colors"
            :disabled="isLoading">
            <i :class="isLoading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs"></i>
            Segarkan
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Utama Log Pendapatan Internal -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr
              class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200"
            >
              <th class="py-4 px-6">ID Transaksi</th>
              <th class="py-4 px-6">Tanggal Masuk</th>
              <th class="py-4 px-6">Petugas Penerima</th>
              <th class="py-4 px-6">Catatan Sumber Aliran Dana</th>
              <th class="py-4 px-6 text-right">Jumlah (Nominal)</th>
              <th class="py-4 px-6 text-right">File Laporan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <!-- Loading State -->
            <tr v-if="isLoading && filteredIncomes.length === 0">
              <td colspan="5" class="py-12 text-center text-slate-400 font-normal">
                <i class="pi pi-spin pi-spinner text-2xl text-indigo-500 mb-2 block"></i>
                Mengamankan jaringan dan mengunduh log kas operasional...
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-else-if="filteredIncomes.length === 0">
              <td colspan="5" class="py-12 text-center text-slate-400 font-normal">
                Tidak ada riwayat transaksi internal masuk yang terdata di sistem.
              </td>
            </tr>

            <!-- Iterasi Data Grouping -->
            <template v-for="group in groupedIncomes" :key="group.dateStr">
              <!-- Header Grup -->
              <tr class="bg-indigo-50/50 border-y border-indigo-100/50">
                <td colspan="4" class="py-3 px-6 font-bold text-indigo-900 text-xs">
                  <i class="pi pi-calendar mr-2 text-indigo-500"></i>{{ group.dateLabel }}
                </td>
                <td class="py-3 px-6 text-right font-black text-indigo-700">
                  {{ formatCurrency(group.total) }}
                </td>
                <td class="py-3 px-6 text-right space-x-2 flex justify-end items-center">
                  <button
                    @click="viewPDF(group.dateStr)"
                    class="inline-flex items-center gap-2 px-3 py-1.5 bg-indigo-100 hover:bg-indigo-200 text-indigo-700 text-xs font-semibold rounded-lg border border-transparent transition-colors shadow-sm">
                    <i class="pi pi-eye text-xs"></i>
                    View
                  </button>
                  <button
                    @click="exportToPDF(group.dateStr)"
                    class="inline-flex items-center gap-2 px-3 py-1.5 bg-rose-600 hover:bg-rose-700 text-white text-xs font-semibold rounded-lg border border-transparent transition-colors shadow-sm"
                    :disabled="isExporting">
                    <i :class="isExporting ? 'pi pi-spin pi-spinner' : 'pi pi-file-pdf'" class="text-xs"></i>
                    Export PDF
                  </button>
                </td>
              </tr>
              
              <!-- Baris Transaksi -->
              <tr
                v-for="income in group.items"
                :key="income.id"
                class="hover:bg-slate-50/60 transition-colors">
                <td class="py-4 px-6 font-mono text-xs text-slate-400 select-all">
                  {{ income.id }}
                </td>
                <td class="py-4 px-6 text-slate-500 font-normal text-xs">
                  {{ formatDate(income.created_at) }}
                </td>
                <td class="py-4 px-6">
                  <div class="flex items-center gap-2">
                    <div
                      class="w-6 h-6 rounded-full bg-slate-100 border border-slate-200 flex items-center justify-center text-[10px] text-slate-600 font-bold uppercase">
                      {{ (income.created_by || 'ST').substring(0, 2) }}
                    </div>
                    <span class="text-slate-900 font-semibold">{{
                      income.created_by || 'Sistem Otomatis'}}</span>
                  </div>
                </td>
                <td
                  class="py-4 px-6 text-slate-600 font-normal max-w-xs truncate"
                  :title="income.notes">
                  {{ income.notes || 'Tanpa keterangan tambahan.' }}
                </td>
                <td class="py-4 px-6 text-right font-mono font-bold text-emerald-600">
                  {{ formatCurrency(income.amount || 0) }}
                </td>
                <td class="py-4 px-6 text-right">
                  <!-- Aksi dieksekusi di level grup tanggal -->
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
