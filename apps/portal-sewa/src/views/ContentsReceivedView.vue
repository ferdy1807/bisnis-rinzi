<template>
  <div class="max-w-[1440px] mx-auto space-y-6 pb-16 font-sans text-slate-800 antialiased">
    
    <div class="flex flex-col lg:flex-row lg:items-end justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <p class="text-[11px] font-bold tracking-widest text-indigo-600 uppercase mb-1">
          Alur Kerja (Workflow)
        </p>
        <h2 class="text-2xl font-extrabold text-slate-900 tracking-tight">
          Penerimaan Barang Titipan
        </h2>
      </div>
      
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 w-full lg:w-auto">
        <Dropdown
          v-model="filterStatus"
          :options="filterStatusOptions"
          optionLabel="label"
          optionValue="value"
          class="w-full sm:w-48 rounded-xl border-slate-300 shadow-xs text-sm bg-white h-[42px]"
          :pt="{ input: { class: 'py-2.5 px-3' } }"
        />
        <span class="p-input-icon-left w-full sm:w-80">
          <i class="pi pi-search text-slate-400" />
          <InputText
            v-model="searchQuery"
            placeholder="Cari Pengantin, nota, atau item..."
            class="w-full rounded-xl border-slate-300 focus:border-indigo-500 font-medium pl-10 py-2.5 h-[42px]"
          />
        </span>
        <Button
          label="Segarkan"
          icon="pi pi-refresh"
          class="p-button-outlined bg-white hover:bg-slate-50 text-slate-600 border border-slate-300 rounded-xl shadow-xs font-bold px-4 py-2.5 h-[42px] transition-all cursor-pointer text-xs"
          :loading="rentalStore.isLoading.active"
          @click="refreshData"
        />
      </div>
    </div>

    <Toast />

    <div class="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-xs">
      <DataTable
        :value="filteredReservations"
        :paginator="true"
        :rows="10"
        :loading="rentalStore.isLoading.active"
        class="p-datatable-sm text-sm border-none custom-table-rows"
        responsiveLayout="scroll"
        rowHover
      >
        <template #empty>
          <div class="p-16 text-center text-slate-400 flex flex-col items-center justify-center bg-slate-50/50">
            <i class="pi pi-inbox text-5xl mb-3 text-slate-300"></i>
            <h4 class="font-extrabold text-slate-700 m-0">Tidak Ada Pesanan Aktif</h4>
            <p class="text-xs text-slate-500 mt-1 max-w-xs leading-relaxed m-0">
              Belum ada antrean pesanan aktif berstatus BOOKED atau READY yang sesuai dengan kriteria filter saat ini.
            </p>
          </div>
        </template>

        <Column header="Identitas Pemesan" class="w-[30%] p-4">
          <template #body="{ data }">
            <div class="flex items-center gap-3.5">
              <div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center shrink-0 border border-indigo-100 shadow-2xs">
                <i class="pi pi-user text-base"></i>
              </div>
              <div class="flex flex-col min-w-0">
                <span class="font-bold text-slate-900 text-sm truncate">
                  {{ parseCustomer(data.customer_snapshot_id).name }}
                </span>
                <div class="flex items-center gap-2 mt-1 flex-wrap">
                  <span class="text-[10px] font-mono font-bold text-indigo-700 bg-indigo-50 px-1.5 py-0.5 rounded border border-indigo-100 shadow-3xs">
                    {{ data.invoice_number }}
                  </span>
                  <span class="text-xs text-slate-500 inline-flex items-center gap-1 font-medium">
                    <i class="pi pi-phone text-[9px] text-slate-400"></i>
                    {{ parseCustomer(data.customer_snapshot_id).phone }}
                  </span>
                </div>
              </div>
            </div>
          </template>
        </Column>

        <Column header="Target Selesai" class="w-[18%] p-4">
          <template #body="{ data }">
            <div class="flex flex-col">
              <span class="font-bold text-slate-800 text-xs font-mono bg-slate-50 border border-slate-200 px-2 py-1 rounded-lg w-max shadow-3xs">
                {{ formatDate(data.start_date) }}
              </span>
              <span class="text-[10px] text-slate-400 font-semibold tracking-wide uppercase mt-1">Jadwal Acara</span>
            </div>
          </template>
        </Column>

        <Column header="Status Titipan" class="w-[18%] p-4">
          <template #body="{ data }">
            <div v-if="data.contents && data.contents.length > 0" class="inline-flex items-center gap-2 bg-emerald-50 border border-emerald-200 px-2.5 py-1 rounded-xl shadow-3xs">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
              <span class="text-[11px] font-bold text-emerald-700">
                {{ data.contents.length }} Barang Tercatat
              </span>
            </div>
            <div v-else class="inline-flex items-center gap-2 bg-rose-50 border border-rose-100 px-2.5 py-1 rounded-xl shadow-3xs">
              <span class="w-1.5 h-1.5 rounded-full bg-rose-400"></span>
              <span class="text-[11px] font-bold text-rose-700">
                Belum Ada Barang
              </span>
            </div>
          </template>
        </Column>

        <Column header="List Barang Tercatat" class="w-[24%] p-4">
          <template #body="{ data }">
            <div v-if="data.contents && data.contents.length > 0" class="flex flex-wrap gap-1.5 max-w-[280px]">
              <span
                v-for="(item, index) in data.contents"
                :key="index"
                class="text-[10px] bg-slate-50 border border-slate-200/80 text-slate-600 font-medium px-2 py-0.5 rounded-md truncate max-w-[130px] shadow-3xs"
                :title="item.item_name"
              >
                {{ item.item_name }} <span class="font-bold text-slate-400">({{ item.quantity }})</span>
              </span>
            </div>
            <span v-else class="text-xs text-slate-400 italic font-medium pl-1">-</span>
          </template>
        </Column>

        <Column header="Aksi" class="w-[10%] text-right p-4">
          <template #body="{ data }">
            <Button
              label="Kelola"
              icon="pi pi-box"
              class="p-button-sm bg-indigo-600 hover:bg-indigo-700 text-white border-none rounded-xl font-bold px-3.5 py-2 shadow-xs transition-all cursor-pointer text-xs"
              @click="openModal(data)"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <ManageDepositContentsModal
      v-model:visible="showModal"
      :reservation="selectedReservation"
      @saved="refreshData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRentalStore } from '@frontend/stores/rental'
import type { RentalReservation } from '@frontend/types/rental'
import ManageDepositContentsModal from '../components/ManageDepositContentsModal.vue'
import { useRoute, useRouter } from 'vue-router'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'

const toast = useToast()
const route = useRoute()
const router = useRouter()
const rentalStore = useRentalStore()

const searchQuery = ref('')
const showModal = ref(false)
const selectedReservation = ref<RentalReservation | null>(null)

const filterStatusOptions = [
  { label: 'Semua Status', value: 'ALL' },
  { label: 'Tercatat (Ada Barang)', value: 'TERCATAT' },
  { label: 'Belum Ada Barang', value: 'BELUM' },
]
const filterStatus = ref('ALL')

const filteredReservations = computed(() => {
  let list = [...rentalStore.bookedList, ...rentalStore.readyList]

  if (filterStatus.value === 'TERCATAT') {
    list = list.filter((r) => r.contents && r.contents.length > 0)
  } else if (filterStatus.value === 'BELUM') {
    list = list.filter((r) => !r.contents || r.contents.length === 0)
  }

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter((r) => {
      const matchInvoice = r.invoice_number.toLowerCase().includes(q)
      const matchCustomer = parseCustomer(r.customer_snapshot_id).name.toLowerCase().includes(q)
      const matchItems = r.contents?.some((c) => c.item_name.toLowerCase().includes(q)) ?? false
      return matchInvoice || matchCustomer || matchItems
    })
  }
  return list
})

const refreshData = async () => {
  // 1. Tarik data terbaru dari peladen masuk ke dalam store
  await rentalStore.fetchActiveReservations()

  // 2. AMANKAN REAKTIVITAS MODAL: Jika modal sedang terbuka, perbarui isinya secara *live*
  if (showModal.value && selectedReservation.value) {
    // Gabungkan seluruh baris antrean untuk mencari ulang data yang paling mutakhir
    const allReservations = [...rentalStore.bookedList, ...rentalStore.readyList]
    const updatedData = allReservations.find(r => r.id === selectedReservation.value?.id)
    
    if (updatedData) {
      // Perbarui referensi objek agar perubahan array 'contents' langsung meremajakan UI Modal
      selectedReservation.value = { ...updatedData }
    }
  }

  // 3. Logika auto-open dari query parameter URL (jika ada)
  if (route.query.id && !showModal.value) {
    const targetId = route.query.id as string
    const found = rentalStore.bookedList.find((r) => r.id === targetId)
    if (found) {
      openModal(found)
      const newQuery = { ...route.query }
      delete newQuery.id
      router.replace({ query: newQuery })
    }
  }
}

const openModal = (reservation: RentalReservation) => {
  selectedReservation.value = reservation
  showModal.value = true
}

const formatDate = (d: string) =>
  new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })

const parseCustomer = (str?: string) => {
  if (!str) return { name: 'Tanpa Nama', phone: '-' }
  const m = str.match(/^(.*?)\s*\((.*?)\)$/)
  return m ? { name: m[1]?.trim() || str, phone: m[2]?.trim() || '-' } : { name: str, phone: '-' }
}
</script>

<style scoped>
:deep(.custom-table-rows .p-datatable-tbody > tr > td) {
  border-width: 0 0 1px 0 !important;
  border-color: #f1f5f9 !important;
  padding: 0.85rem 1rem !important;
}
:deep(.p-datatable .p-datatable-thead > tr > th) {
  background: #f8fafc !important;
  color: #475569 !important;
  font-weight: 700 !important;
  padding: 0.75rem 1rem !important;
  border-bottom: 2px solid #e2e8f0 !important;
}
/* Memaksa tinggi komponen kontrol h-[42px] agar seragam di semua browser */
:deep(.p-dropdown) {
  display: flex;
  align-items: center;
}
</style>