<!-- apps/portal-sewa/src/views/UpcomingView.vue -->
<template>
  <div class="max-w-[1440px] mx-auto space-y-6 pb-16 font-sans text-slate-800 antialiased">
    
    <!-- 1. HEADER SECTION & SYNC -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <p class="text-[11px] font-bold tracking-widest text-indigo-600 uppercase mb-1">
          Radar Operasional Hantaran
        </p>
        <h2 class="text-2xl font-extrabold text-slate-900 tracking-tight flex items-center gap-2">
          Jadwal Pengambilan Terdekat
        </h2>
      </div>
      <div class="flex items-center gap-3 self-end sm:self-auto">
        <span class="text-xs text-slate-500 bg-slate-100 px-3 py-1.5 rounded-xl border border-slate-200/60">
          Sinkronisasi terakhir: <strong class="font-mono text-slate-800">{{ lastSynced }}</strong>
        </span>
        <Button
          label="Perbarui Radar"
          icon="pi pi-refresh"
          class="p-button-sm bg-indigo-600 hover:bg-indigo-700 text-white border-none rounded-xl font-bold px-4 py-2.5 shadow-xs cursor-pointer transition-all"
          :loading="rentalStore.isLoading.upcoming"
          @click="syncRadar"
        />
      </div>
    </div>

    <Toast />

    <!-- MAIN TWO-COLUMN LAYOUT -->
    <div class="flex flex-col xl:flex-row-reverse gap-8 items-start">
      
      <!-- RIGHT PANEL: BENTO URGENCY CARD COUNTERS -->
      <div class="w-full xl:w-[300px] shrink-0 xl:sticky xl:top-6">
        <div class="grid grid-cols-2 xl:grid-cols-1 gap-4">
          
          <!-- Pill 1: Total Antrean -->
          <div
            @click="setUrgencyFilter('all')"
            class="p-5 rounded-2xl border transition-all duration-300 cursor-pointer flex flex-col justify-between min-h-[120px] shadow-xs relative overflow-hidden"
            :class="
              urgencyFilter === 'all'
                ? 'bg-slate-900 text-white border-slate-900 shadow-md ring-2 ring-indigo-600/20'
                : 'bg-white text-slate-800 border-slate-200 hover:border-indigo-400 hover:bg-slate-50/50'
            "
          >
            <div class="flex items-center justify-between mb-2">
              <span
                class="w-8 h-8 rounded-xl flex items-center justify-center transition-all"
                :class="urgencyFilter === 'all' ? 'bg-white/10 text-white' : 'bg-slate-100 text-indigo-600'"
              >
                <i class="pi pi-calendar-clock text-base"></i>
              </span>
              <span
                class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md"
                :class="urgencyFilter === 'all' ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-500'"
              >
                Semua
              </span>
            </div>
            <div>
              <p class="text-[11px] opacity-75 m-0 font-medium">Antrean Terjadwal</p>
              <h3 class="text-2xl font-black font-mono mt-1 mb-0">
                {{ processedUpcoming.length }} <span class="text-xs font-normal opacity-70">Nota</span>
              </h3>
            </div>
          </div>

          <!-- Pill 2: Kritis (H-1 / Hari Ini) -->
          <div
            @click="setUrgencyFilter('besok')"
            class="p-5 rounded-2xl border transition-all duration-300 cursor-pointer flex flex-col justify-between min-h-[120px] shadow-xs relative overflow-hidden"
            :class="
              urgencyFilter === 'besok'
                ? 'bg-red-50 text-red-900 border-red-300 shadow-md ring-2 ring-red-500/10'
                : 'bg-white text-slate-800 border-slate-200 hover:border-red-300 hover:bg-red-50/20'
            "
          >
            <div class="flex items-center justify-between mb-2">
              <span
                class="w-8 h-8 rounded-xl flex items-center justify-center transition-all"
                :class="urgencyFilter === 'besok' ? 'bg-red-600 text-white' : 'bg-red-50 text-red-600'"
              >
                <i class="pi pi-bell text-base" :class="{ 'animate-wiggle': countBesok > 0 }"></i>
              </span>
              <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md bg-red-100/70 text-red-800">
                Kritis
              </span>
            </div>
            <div>
              <p class="text-[11px] opacity-80 m-0 font-bold text-red-700">Ambil H-1 / Hari Ini</p>
              <h3 class="text-2xl font-black font-mono mt-1 mb-0 text-red-700">
                {{ countBesok }} <span class="text-xs font-normal opacity-70">Nota</span>
              </h3>
            </div>
          </div>

          <!-- Pill 3: Siaga (H-2 s/d H-4) -->
          <div
            @click="setUrgencyFilter('minggu_ini')"
            class="p-5 rounded-2xl border transition-all duration-300 cursor-pointer flex flex-col justify-between min-h-[120px] shadow-xs relative overflow-hidden"
            :class="
              urgencyFilter === 'minggu_ini'
                ? 'bg-amber-50 text-amber-900 border-amber-300 shadow-md ring-2 ring-amber-500/10'
                : 'bg-white text-slate-800 border-slate-200 hover:border-amber-300 hover:bg-amber-50/20'
            "
          >
            <div class="flex items-center justify-between mb-2">
              <span class="w-8 h-8 rounded-xl bg-amber-100/70 text-amber-700 flex items-center justify-center">
                <i class="pi pi-clock text-base"></i>
              </span>
              <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md bg-amber-100 text-amber-800">
                Persiapan
              </span>
            </div>
            <div>
              <p class="text-[11px] opacity-80 m-0 font-medium text-amber-800">Ambil 2 - 4 Hari</p>
              <h3 class="text-2xl font-black font-mono mt-1 mb-0 text-amber-700">
                {{ countMingguIni }} <span class="text-xs font-normal opacity-70">Nota</span>
              </h3>
            </div>
          </div>

          <!-- Pill 4: Aman (H-5+) -->
          <div
            @click="setUrgencyFilter('santai')"
            class="p-5 rounded-2xl border transition-all duration-300 cursor-pointer flex flex-col justify-between min-h-[120px] shadow-xs relative overflow-hidden"
            :class="
              urgencyFilter === 'santai'
                ? 'bg-emerald-50 text-emerald-900 border-emerald-300 shadow-md ring-2 ring-emerald-500/10'
                : 'bg-white text-slate-800 border-slate-200 hover:border-emerald-300 hover:bg-emerald-50/20'
            "
          >
            <div class="flex items-center justify-between mb-2">
              <span class="w-8 h-8 rounded-xl bg-emerald-100/70 text-emerald-700 flex items-center justify-center">
                <i class="pi pi-shield text-base"></i>
              </span>
              <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-md bg-emerald-100 text-emerald-800">
                Aman
              </span>
            </div>
            <div>
              <p class="text-[11px] opacity-80 m-0 font-medium text-emerald-800">Ambil &gt; 5 Hari</p>
              <h3 class="text-2xl font-black font-mono mt-1 mb-0 text-emerald-700">
                {{ countSantai }} <span class="text-xs font-normal opacity-70">Nota</span>
              </h3>
            </div>
          </div>

        </div>
      </div>

      <!-- LEFT PANEL: DATA TABLES & SEARCH FILTERS -->
      <div class="flex-1 w-full space-y-4 min-w-0">
        
        <!-- SEARCH SEARCH BAR -->
        <div class="bg-white border border-slate-200 rounded-2xl p-4 shadow-xs flex flex-col sm:flex-row items-center justify-between gap-4">
          <div class="flex items-center gap-3 w-full sm:w-auto">
            <span class="p-input-icon-left w-full sm:w-80">
              <i class="pi pi-search text-slate-400"></i>
              <InputText
                v-model="searchQuery"
                placeholder="Cari Calon Pengantin / Nota..."
                class="p-inputtext-sm w-full rounded-xl border-slate-300 focus:border-indigo-500 font-medium pl-10 py-2.5"
              />
            </span>
            <Button
              v-if="searchQuery || urgencyFilter !== 'all'"
              icon="pi pi-filter-slash"
              class="p-button-outlined p-button-sm text-red-600 border-slate-200 hover:bg-red-50 rounded-xl h-[42px] px-3 transition-all"
              @click="resetFilters"
              title="Hapus Filter"
            />
          </div>

          <div class="text-xs font-bold text-slate-500 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-emerald-500 inline-block animate-pulse"></span>
            <span>Menampilkan <strong>{{ filteredList.length }}</strong> Jadwal Terurut</span>
          </div>
        </div>

        <!-- MAIN DATA TABLE -->
        <div class="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-xs flex flex-col">
          <DataTable
            :value="filteredList"
            :loading="rentalStore.isLoading.upcoming"
            paginator
            :rows="10"
            class="p-datatable-sm text-sm border-none custom-table-rows"
            responsiveLayout="scroll"
            rowHover
          >
            <template #empty>
              <div class="text-center py-16 text-slate-400 space-y-3 bg-slate-50/50">
                <i class="pi pi-radar text-4xl text-slate-300 block animate-pulse"></i>
                <p class="text-sm font-bold text-slate-700">
                  Tidak ada jadwal hantaran yang terdeteksi
                </p>
                <Button
                  label="Lihat Seluruh Antrean"
                  class="p-button-text p-button-sm text-indigo-600 font-bold rounded-lg hover:bg-indigo-50"
                  @click="resetFilters"
                />
              </div>
            </template>

            <!-- Col 1: Hitung Mundur -->
            <Column header="Hitung Mundur" class="w-36 p-3">
              <template #body="{ data }">
                <span
                  :class="getUrgencyPillClass(data.days_left)"
                  class="px-3 py-1.5 rounded-xl font-mono font-bold text-[11px] block text-center shadow-2xs transition-all"
                >
                  {{ data.countdown_text }}
                </span>
              </template>
            </Column>

            <!-- Col 2: Nota & Tanggal Pakai -->
            <Column header="Nota & Tgl Ambil" class="w-40 p-3">
              <template #body="{ data }">
                <span class="font-mono font-bold text-indigo-700 block text-xs">{{ data.invoice_number }}</span>
                <span class="text-[11px] font-semibold text-slate-500 block mt-0.5">{{ formatDate(data.start_date) }}</span>
              </template>
            </Column>

            <!-- Col 3: Calon Pengantin + Contact Integration -->
            <Column header="Calon Pengantin / Penyewa" class="min-w-[240px] p-3">
              <template #body="{ data }">
                <div class="space-y-1.5">
                  <span class="font-bold text-slate-900 text-sm block leading-tight">
                    {{ parseCustomer(data.customer_snapshot_id).name }}
                  </span>
                  <a
                    v-if="parseCustomer(data.customer_snapshot_id).phone !== '-'"
                    :href="generateWaLink(data)"
                    target="_blank"
                    class="text-[11px] font-mono text-emerald-700 hover:text-white bg-emerald-50 hover:bg-emerald-600 px-2 py-1 rounded-lg border border-emerald-200 transition-all inline-flex items-center gap-1.5 font-bold no-underline shadow-2xs"
                    @click.stop
                  >
                    <i class="pi pi-whatsapp text-xs"></i> Follow-Up WA
                  </a>
                </div>
              </template>
            </Column>

            <!-- Col 4: Status Barang -->
            <Column header="Status Barang Titipan" class="min-w-[180px] p-3">
              <template #body="{ data }">
                <Tag
                  :severity="hasDepositItems(data) ? 'success' : 'danger'"
                  :value="hasDepositItems(data) ? 'Sudah Drop Barang' : 'Belum Titip Barang'"
                  class="font-bold text-[11px] px-2.5 py-1 rounded-lg"
                />
              </template>
            </Column>

            <!-- Col 5: Sisa Tagihan -->
            <Column header="Kekurangan Bayar" class="text-right w-36 p-3">
              <template #body="{ data }">
                <span
                  class="font-mono font-bold text-xs"
                  :class="data.total_amount - data.down_payment > 0 ? 'text-red-600' : 'text-emerald-600'"
                >
                  {{ formatCurrency(data.total_amount - data.down_payment) }}
                </span>
              </template>
            </Column>

            <!-- Col 6: Aksi Dekorator -->
            <Column header="Lembar Kerja" class="text-center w-32 p-3">
              <template #body="{ data }">
                <Button
                  label="Dekorasi"
                  icon="pi pi-wrench"
                  class="p-button-xs bg-slate-900 hover:bg-slate-800 text-white border-none rounded-xl text-[11px] font-bold py-2 px-3 shadow-xs cursor-pointer transition-all"
                  @click="openPreparationSheet(data)"
                />
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

    </div>

    <!-- 5. DIALOG SHEET DEKORATOR -->
    <Dialog
      v-model:visible="showPrepModal"
      modal
      :style="{ width: '48rem' }"
      class="font-sans"
      :pt="{
        root: { class: 'rounded-2xl shadow-2xl border-none overflow-hidden bg-white' },
        mask: { class: 'backdrop-blur-xs bg-slate-900/40 shadow-xs' },
        header: { class: 'bg-slate-50 border-b border-slate-100 px-6 py-4' },
        content: { class: 'px-0 py-0' },
      }"
    >
      <template #header>
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-indigo-50 flex items-center justify-center">
            <i class="pi pi-clipboard text-lg text-indigo-600"></i>
          </div>
          <div>
            <h2 class="text-base font-extrabold text-slate-900 leading-tight m-0">Lembar Kerja Kerja Dekorator</h2>
            <p class="text-xs text-slate-500 font-medium mt-0.5 m-0">
              Ceklis persiapan & inventarisasi barang hantaran masuk
            </p>
          </div>
        </div>
      </template>

      <div v-if="selectedOrder" class="flex flex-col bg-slate-50/40">
        <!-- Client Detail Banner -->
        <div class="bg-slate-900 p-6 text-white flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-800">
          <div>
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest block mb-1">Klien Penyewa</span>
            <strong class="text-xl font-extrabold block tracking-tight">
              {{ parseCustomer(selectedOrder.customer_snapshot_id).name }}
            </strong>
            <span class="text-[11px] font-mono text-slate-400 mt-2 block">
              <i class="pi pi-hashtag text-[10px] mr-0.5"></i> {{ selectedOrder.invoice_number }}
            </span>
          </div>
          <div class="bg-slate-800 p-3 rounded-xl border border-slate-700/60 text-left sm:text-right">
            <span class="text-[10px] text-slate-400 uppercase tracking-widest font-bold block mb-1">Target Pengambilan</span>
            <strong class="text-amber-400 font-extrabold text-base font-mono">
              {{ formatDate(selectedOrder.start_date) }}
            </strong>
          </div>
        </div>

        <div class="p-6 space-y-6">
          <!-- Item Box Section -->
          <div class="space-y-3">
            <h3 class="text-xs font-bold text-slate-900 uppercase tracking-wider flex items-center gap-2 m-0">
              <span class="w-6 h-6 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center"><i class="pi pi-box text-xs"></i></span>
              Kebutuhan Box Utama
            </h3>
            <div class="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-2xs">
              <DataTable :value="selectedOrder.items || []" class="p-datatable-sm text-sm border-none">
                <template #empty>
                  <div class="p-6 text-center text-slate-400 font-medium italic">
                    Rincian tipe wadah mengikuti nota paket induk
                  </div>
                </template>
                <Column field="rental_product_name" header="Tipe Wadah / Nampan" class="font-semibold text-slate-700 p-3"></Column>
                <Column field="qty" header="Jumlah Kuota" class="text-right w-40 p-3">
                  <template #body="{ data }">
                    <span class="bg-indigo-50 text-indigo-700 font-bold font-mono px-2.5 py-1 rounded-md border border-indigo-100">
                      {{ data.qty }} Unit
                    </span>
                  </template>
                </Column>
              </DataTable>
            </div>
          </div>

          <!-- Deposit Content Item Section -->
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <h3 class="text-xs font-bold text-slate-900 uppercase tracking-wider flex items-center gap-2 m-0">
                <span class="w-6 h-6 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center"><i class="pi pi-gift text-xs"></i></span>
                Isian Barang Pengantin (Titipan)
              </h3>
              <Button
                label="Kelola Barang"
                icon="pi pi-pencil"
                class="p-button-text p-button-xs text-emerald-600 hover:bg-emerald-50 font-bold px-3 py-1.5 rounded-xl border border-emerald-100 bg-white transition-all"
                @click="openContentsModal"
              />
            </div>

            <div class="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-2xs">
              <div
                v-if="!selectedOrder.contents || selectedOrder.contents.length === 0"
                class="p-8 flex flex-col items-center justify-center bg-red-50/40 text-center border border-red-100/60 rounded-xl m-3"
              >
                <div class="w-10 h-10 bg-red-100 text-red-500 rounded-full flex items-center justify-center mb-2.5">
                  <i class="pi pi-exclamation-triangle text-base"></i>
                </div>
                <h4 class="text-xs font-extrabold text-red-900 mb-0.5">Belum Menitipkan Barang</h4>
                <p class="text-[11px] text-red-700/80 max-w-xs leading-relaxed m-0">
                  Calon pengantin belum menyerahkan item seserahan fisik ke toko. Harap hubungi penyewa untuk follow-up barang.
                </p>
              </div>
              <DataTable v-else :value="selectedOrder.contents" class="p-datatable-sm text-sm border-none">
                <Column field="item_name" header="Nama Barang" class="font-bold text-slate-800 p-3"></Column>
                <Column field="description" header="Deskripsi / Isi Paket" class="p-3 text-slate-600"></Column>
                <Column field="condition_notes" header="Kondisi Tiba" class="italic text-slate-400 text-xs p-3"></Column>
              </DataTable>
            </div>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="p-4 bg-white border-t border-slate-100 flex items-center justify-between">
          <Button
            label="Tutup Lembar"
            icon="pi pi-times"
            class="p-button-text text-slate-500 hover:bg-slate-100 font-bold p-button-sm rounded-xl px-4 py-2"
            @click="closePrepModal"
          />
          <router-link to="/" class="no-underline">
            <Button
              label="Buka Papan Kanban"
              icon="pi pi-arrow-right"
              iconPos="right"
              class="bg-indigo-600 hover:bg-indigo-700 text-white border-none font-bold p-button-sm rounded-xl px-5 py-2.5 shadow-xs cursor-pointer transition-all"
            />
          </router-link>
        </div>
      </div>
    </Dialog>

    <!-- FLOATING MANAGE CONTENTS DIALOG -->
    <ManageDepositContentsModal
      v-model:visible="showContentsModal"
      :reservation="selectedOrder"
      @saved="refreshRadar"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRentalStore } from '@frontend/stores/rental'
import type { RentalReservation } from '@frontend/types/rental'
import ManageDepositContentsModal from '../components/ManageDepositContentsModal.vue'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Tag from 'primevue/tag'
import Toast from 'primevue/toast'

const rentalStore = useRentalStore()

const searchQuery = ref('')
const urgencyFilter = ref<'all' | 'besok' | 'minggu_ini' | 'santai'>('all')
const showPrepModal = ref(false)
const showContentsModal = ref(false)
const selectedOrder = ref<RentalReservation | null>(null)

const lastSynced = ref(
  new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }),
)

const syncRadar = async () => {
  await rentalStore.fetchUpcomingReservations(true)
  lastSynced.value = new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  rentalStore.fetchUpcomingReservations(true)
})

const processedUpcoming = computed(() => {
  const rawList = rentalStore.upcomingReservations || []
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  return rawList
    .map((item: RentalReservation) => {
      const targetDate = new Date(item.start_date)
      targetDate.setHours(0, 0, 0, 0)

      const diffTime = targetDate.getTime() - today.getTime()
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

      let label = ''
      if (diffDays === 0) label = 'HARI INI'
      else if (diffDays === 1) label = 'BESOK PAGI'
      else if (diffDays > 1) label = `H - ${diffDays}`
      else label = 'PASCA JADWAL'

      return {
        ...item,
        days_left: diffDays,
        countdown_text: label,
      }
    })
    .sort((a, b) => a.days_left - b.days_left)
})

const countBesok = computed(
  () => processedUpcoming.value.filter((i) => i.days_left <= 1 && i.days_left >= 0).length,
)
const countMingguIni = computed(
  () => processedUpcoming.value.filter((i) => i.days_left >= 2 && i.days_left <= 4).length,
)
const countSantai = computed(() => processedUpcoming.value.filter((i) => i.days_left >= 5).length)

const setUrgencyFilter = (f: 'all' | 'besok' | 'minggu_ini' | 'santai') => {
  urgencyFilter.value = f
}
const resetFilters = () => {
  searchQuery.value = ''
  urgencyFilter.value = 'all'
}

const filteredList = computed(() => {
  let list = processedUpcoming.value

  if (urgencyFilter.value === 'besok')
    list = list.filter((i) => i.days_left <= 1 && i.days_left >= 0)
  else if (urgencyFilter.value === 'minggu_ini')
    list = list.filter((i) => i.days_left >= 2 && i.days_left <= 4)
  else if (urgencyFilter.value === 'santai') 
    list = list.filter((i) => i.days_left >= 5)

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter(
      (i) =>
        i.invoice_number.toLowerCase().includes(q) ||
        i.customer_snapshot_id.toLowerCase().includes(q),
    )
  }

  return list
})

const parseCustomer = (str?: string) => {
  if (!str) return { name: 'Tanpa Nama', phone: '-' }
  const m = str.match(/^(.*?)\s*\((.*?)\)$/)
  return m ? { name: m[1]?.trim() || str, phone: m[2]?.trim() || '-' } : { name: str, phone: '-' }
}

const hasDepositItems = (r: RentalReservation) => {
  return Array.isArray(r.contents) && r.contents.length > 0
}

const openPreparationSheet = (order: RentalReservation) => {
  selectedOrder.value = order
  showPrepModal.value = true
}

const closePrepModal = () => {
  showPrepModal.value = false
}

const openContentsModal = () => {
  showContentsModal.value = true
}

const refreshRadar = async () => {
  await rentalStore.fetchUpcomingReservations(true)
}

const generateWaLink = (order: any) => {
  const cust = parseCustomer(order.customer_snapshot_id)
  if (cust.phone === '-') return '#'

  const cleanPhone = '62' + cust.phone.replace(/\D/g, '').replace(/^0/, '')
  const tgl = formatDate(order.start_date)

  const text = `Halo Kak ${cust.name}, salam dari Rinzi Kotak Seserahan. Mengingatkan bahwa jadwal pengambilan pesanan box hantaran Kakak (Nota: ${order.invoice_number}) adalah pada tanggal *${tgl}*. Apakah barang-barang isian seserahannya sudah bisa dikirimkan ke toko kami hari ini? Terima kasih Kak! 🙏`

  return `https://wa.me/${cleanPhone}?text=${encodeURIComponent(text)}`
}

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)

const formatDate = (dStr: string) => {
  return new Date(dStr).toLocaleDateString('id-ID', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

const getUrgencyPillClass = (days: number) => {
  if (days === 0) return 'bg-red-600 text-white font-extrabold shadow-sm'
  if (days === 1) return 'bg-red-50 text-red-700 border border-red-200'
  if (days >= 2 && days <= 4) return 'bg-amber-50 text-amber-700 border border-amber-200'
  if (days >= 5) return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
  return 'bg-slate-100 text-slate-600'
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

@keyframes wiggle {
  0%, 100% { transform: rotate(-6deg); }
  50% { transform: rotate(6deg); }
}
.animate-wiggle {
  display: inline-block;
  animation: wiggle 0.4s ease-in-out infinite;
}
</style>