<template>
  <ConfirmDialog
    :pt="{
      root: {
        class:
          'border-none shadow-2xl rounded-2xl overflow-hidden bg-white/95 backdrop-blur-xl max-w-sm w-full mx-4',
      },
      header: { class: 'pb-0 pt-6 px-6 border-none text-[#1a1c20]' },
      title: { class: 'text-xl font-black' },
      content: { class: 'px-6 py-3 text-[#58627c] text-[13px] leading-relaxed' },
      footer: { class: 'px-6 pb-6 pt-3 border-none flex justify-end gap-2' },
      icon: { class: 'text-3xl mr-4' },
    }"
  ></ConfirmDialog>
  <div class="space-y-8 pb-12 font-sans text-[#1a1c20]">
    <div
      class="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-[#c4c6d1]/60 pb-6"
    >
      <div>
        <p class="text-xs font-extrabold tracking-widest text-[#254582] uppercase mb-1">
          Executive Overview
        </p>
        <h2 class="text-3xl font-black text-[#1a1c20] tracking-tight">Pusat Kendali Sewa</h2>
      </div>
      <div class="flex items-center gap-3">
        <Button
          icon="pi pi-refresh"
          class="p-button-rounded p-button-text p-button-sm text-[#58627c] hover:bg-[#e2e2e7]"
          :loading="rentalStore.isLoading.active || rentalStore.isLoading.products"
          @click="refreshDashboard"
        />
        <router-link to="/reservations-create">
          <button
            class="bg-[#254582] text-white px-6 py-3 rounded-xl font-semibold text-sm flex items-center gap-2 shadow-lg shadow-[#254582]/20 hover:scale-[1.02] active:scale-95 transition-all"
          >
            <i class="pi pi-plus text-xs"></i> Buat Reservasi Baru
          </button>
        </router-link>
      </div>
    </div>
    <Toast />
    <!-- TOP SECTION: KPI (2x2) & Agenda -->
    <div class="grid grid-cols-1 xl:grid-cols-12 gap-6 items-stretch mb-6">
      <!-- KIRI: KPI Cards (2x2 Grid) -->
      <div class="xl:col-span-5 grid grid-cols-1 sm:grid-cols-2 gap-2">
        <router-link to="/catalog" class="block h-full cursor-pointer">
          <div
            class="h-full bg-white/80 backdrop-blur-md border border-white/50 p-5 rounded-2xl shadow-xs hover:shadow-md hover:-translate-y-1 transition-all duration-300"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="text-[#254582] bg-[#254582]/10 p-2 rounded-xl"
                ><i class="pi pi-box text-lg"></i
              ></span>
              <span
                class="text-[10px] font-bold text-green-700 bg-green-50 px-2 py-0.5 rounded-full"
                >Katalog</span
              >
            </div>
            <p class="text-[11px] font-semibold text-[#58627c] mb-0.5">Stok Tersedia</p>
            <h3 class="text-xl font-black text-[#1a1c20]">
              {{ rentalStore.products.reduce((acc, p) => acc + p.quantity_available, 0) }}
            </h3>
          </div>
        </router-link>

        <router-link to="/workflow-contents-received" class="block h-full cursor-pointer">
          <div
            class="h-full bg-white/80 backdrop-blur-md border border-white/50 p-5 rounded-2xl shadow-xs hover:shadow-md hover:-translate-y-1 transition-all duration-300"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="text-blue-600 bg-blue-50 p-2 rounded-xl"
                ><i class="pi pi-bookmark text-lg"></i
              ></span>
              <span
                class="text-[10px] font-bold text-blue-700 bg-blue-50 px-2 py-0.5 rounded-full animate-pulse"
                >Action Req.</span
              >
            </div>
            <p class="text-[11px] font-semibold text-[#58627c] mb-0.5">Menunggu Pengisian</p>
            <h3 class="text-xl font-black text-blue-600">
              {{ rentalStore.bookedList.length }}
              <span class="text-[10px] font-normal text-gray-400">Box</span>
            </h3>
          </div>
        </router-link>

        <router-link to="/returns" class="block h-full cursor-pointer">
          <div
            class="h-full bg-white/80 backdrop-blur-md border border-white/50 p-5 rounded-2xl shadow-xs hover:shadow-md hover:-translate-y-1 transition-all duration-300"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="text-amber-600 bg-amber-50 p-2 rounded-xl"
                ><i class="pi pi-sync text-lg"></i
              ></span>
              <span
                class="text-[10px] font-bold text-amber-800 bg-amber-100 px-2 py-0.5 rounded-full"
                >Active</span
              >
            </div>
            <p class="text-[11px] font-semibold text-[#58627c] mb-0.5">List Pengembalian Aset</p>
            <h3 class="text-xl font-black text-amber-600">
              {{ rentalStore.pickedUpList.length }}
              <span class="text-[10px] font-normal text-gray-400">Box</span>
            </h3>
          </div>
        </router-link>

        <router-link to="/returns" class="block h-full cursor-pointer">
          <div
            class="h-full bg-white/80 backdrop-blur-md border-l-4 border-[#ba1a1a] p-5 rounded-2xl shadow-xs hover:shadow-md hover:-translate-y-1 transition-all duration-300"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="text-[#ba1a1a] bg-[#ba1a1a]/10 p-2 rounded-xl"
                ><i class="pi pi-exclamation-triangle text-lg"></i
              ></span>
              <span
                class="text-[10px] font-bold text-[#ba1a1a] bg-[#ffdad6] px-2 py-0.5 rounded-full"
                >Kritis</span
              >
            </div>
            <p class="text-[11px] font-semibold text-[#58627c] mb-0.5">Lewat Batas</p>
            <h3 class="text-xl font-black text-[#ba1a1a]">
              {{ rentalStore.overdueReservations.length }}
            </h3>
          </div>
        </router-link>
      </div>

      <!-- KANAN: Agenda Kalender (Lebar 4 Kolom) -->
      <div class="xl:col-span-7 flex flex-col h-full">
        <div
          class="bg-white border border-[#c4c6d1]/60 rounded-3xl p-6 shadow-xs flex flex-col flex-1 min-h-0"
          style="max-height: 250px"
        >
          <div class="flex justify-between items-center mb-4">
            <h3 class="font-bold text-[#1a1c20] text-lg">Agenda Kalender Terdekat</h3>
            <router-link to="/upcoming" class="text-xs font-bold text-[#254582] hover:underline"
              >Semua &rarr;</router-link
            >
          </div>

          <div class="flex-1 overflow-y-auto pr-3 custom-scrollbar">
            <div v-if="rentalStore.isLoading.upcoming" class="space-y-4">
              <Skeleton height="3rem" v-for="i in 3" :key="i" />
            </div>
            <div
              v-else-if="rentalStore.upcomingReservations.length === 0"
              class="text-xs text-[#58627c] text-center py-8"
            >
              Belum ada pesanan di masa depan
            </div>

            <div
              v-else
              class="space-y-5 relative before:absolute before:left-[11px] before:top-2 before:bottom-2 before:w-[2px] before:bg-[#c4c6d1]"
            >
              <div
                v-for="(agenda, idx) in rentalStore.upcomingReservations"
                :key="agenda.id"
                class="relative pl-8 group"
              >
                <div
                  class="absolute left-0 top-0.5 w-6 h-6 rounded-full flex items-center justify-center border-4 border-white shadow-xs transition-transform group-hover:scale-125 bg-[#254582]"
                ></div>
                <p class="text-[11px] font-extrabold font-mono text-[#254582] mb-0.5">
                  {{ formatDateWithDay(agenda.start_date) }}
                </p>
                <p class="text-xs font-bold text-[#1a1c20]">{{ agenda.customer_snapshot_id }}</p>
                <div class="flex justify-between items-center mt-1">
                  <span class="text-[11px] font-mono text-[#58627c]">{{
                    agenda.invoice_number
                  }}</span>
                  <span
                    class="text-[10px] bg-[#f3f3f9] text-[#58627c] px-1.5 py-0.5 rounded font-semibold"
                    >Sewa: {{ calculateDays(agenda.start_date, agenda.end_date) }} Hari</span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- BOTTOM SECTION: Kanban Board Full Width -->
    <div
      class="bg-white border border-[#c4c6d1]/60 rounded-3xl p-5 shadow-xs flex flex-col min-h-[520px]"
    >
      <div class="flex justify-between items-center mb-6 pb-4 border-b border-[#e2e2e7]">
        <div class="flex items-center gap-3">
          <span class="p-2 bg-[#d9e2ff] text-[#254582] rounded-xl font-black"
            ><i class="pi pi-sitemap text-lg"></i
          ></span>
          <div>
            <h3 class="text-lg font-bold text-[#1a1c20]">Papan Alur Operasional (Kanban)</h3>
            <p class="text-xs text-[#58627c]">Geser atau klik aksi untuk memindahkan tahapan box</p>
          </div>
        </div>
        <Badge
          :value="rentalStore.activeReservations.length + ' Total'"
          class="bg-[#254582] text-white"
        />
      </div>

      <div class="grid grid-cols-1 xl:grid-cols-3 gap-5 flex-1">
        <!-- KOLOM 1: Antrean Baru / Booking -->
        <div class="bg-[#f8fafc] border border-[#e2e8f0] p-4 rounded-2xl flex flex-col">
          <div class="flex justify-between items-center mb-3 pb-2 border-b border-[#cbd5e1]/40">
            <span
              class="text-xs font-extrabold text-[#475569] flex items-center gap-1.5 uppercase tracking-wider"
            >
              <span class="w-2.5 h-2.5 rounded-full bg-slate-400"></span> 1. Antrean Baru
            </span>
            <span class="text-xs font-bold bg-[#e2e8f0] text-[#475569] px-2 py-0.5 rounded-md">{{
              rentalStore.bookedList.length
            }}</span>
          </div>

          <div v-if="rentalStore.isLoading.active" class="space-y-3">
            <Skeleton height="7rem" rounded v-for="i in 2" :key="i" />
          </div>
          <div
            v-else-if="rentalStore.bookedList.length === 0"
            class="m-auto text-center py-12 text-xs text-[#58627c] italic"
          >
            Tidak ada pesanan baru
          </div>
          <div v-else class="space-y-3 flex-1 overflow-y-auto max-h-[400px] pr-1">
            <div
              v-for="item in rentalStore.bookedList"
              :key="item.id"
              class="bg-white p-3.5 rounded-xl border border-[#e2e8f0] shadow-2xs hover:border-slate-400 transition-all group"
            >
              <div class="flex justify-between items-start mb-1.5">
                <span
                  class="font-mono text-xs font-bold text-slate-700 bg-slate-50 px-1.5 py-0.5 rounded"
                  >{{ item.invoice_number }}</span
                >
                <span class="text-[11px] text-[#58627c]"
                  >Pakai: {{ formatDateShort(item.start_date) }}</span
                >
              </div>
              <p class="font-bold text-sm text-[#1a1c20] truncate">
                {{ item.customer_snapshot_id }}
              </p>
              <div class="mt-3 pt-2 border-t border-gray-100 flex justify-end items-center">
                <Button
                  label="Siap Diambil"
                  icon="pi pi-check"
                  iconPos="right"
                  class="p-button-xs bg-slate-600 hover:bg-slate-700 text-white border-none text-[11px] py-1 px-2.5 rounded-lg shadow-xs"
                  :loading="rentalStore.isLoading.mutate"
                  @click="safeConfirmReady(item)"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- KOLOM 2: Menunggu Diambil -->
        <div class="bg-[#f3f3f9] border border-[#d9e2ff] p-4 rounded-2xl flex flex-col">
          <div class="flex justify-between items-center mb-3 pb-2 border-b border-[#c4c6d1]/40">
            <span
              class="text-xs font-extrabold text-[#254582] flex items-center gap-1.5 uppercase tracking-wider"
            >
              <span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> 2. Menunggu Diambil
            </span>
            <span class="text-xs font-bold bg-[#d9e2ff] text-[#254582] px-2 py-0.5 rounded-md">{{
              rentalStore.readyList.length
            }}</span>
          </div>

          <div v-if="rentalStore.isLoading.active" class="space-y-3">
            <Skeleton height="7rem" rounded v-for="i in 2" :key="i" />
          </div>
          <div
            v-else-if="rentalStore.readyList.length === 0"
            class="m-auto text-center py-12 text-xs text-[#58627c] italic"
          >
            Tidak ada barang siap diambil
          </div>
          <div v-else class="space-y-3 flex-1 overflow-y-auto max-h-[400px] pr-1">
            <div
              v-for="item in rentalStore.readyList"
              :key="item.id"
              class="bg-white p-3.5 rounded-xl border border-[#e2e2e7] shadow-2xs hover:border-blue-400 transition-all group"
            >
              <div class="flex justify-between items-start mb-1.5">
                <span
                  class="font-mono text-xs font-bold text-blue-700 bg-blue-50 px-1.5 py-0.5 rounded"
                  >{{ item.invoice_number }}</span
                >
                <span class="text-[11px] text-[#58627c]"
                  >Pakai: {{ formatDateShort(item.start_date) }}</span
                >
              </div>
              <p class="font-bold text-sm text-[#1a1c20] truncate">
                {{ item.customer_snapshot_id }}
              </p>
              <div class="mt-3 pt-2 border-t border-gray-100 flex justify-between items-center">
                <Button
                  icon="pi pi-undo"
                  class="rounded-full w-6 h-6 bg-slate-300 hover:bg-slate-500 border-none text-white cursor-pointer transition-all duration-200"
                  title="Batal Siap (Kembalikan ke Antrean)"
                  @click="safeConfirmUndoReady(item)"
                />
                <Button
                  label="Konfirmasi Pickup"
                  icon="pi pi-arrow-right"
                  iconPos="right"
                  class="p-button-xs bg-blue-600 hover:bg-blue-700 text-white border-none text-[11px] py-1 px-2.5 rounded-lg shadow-xs"
                  :loading="rentalStore.isLoading.mutate"
                  @click="safeConfirmPickup(item)"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- KOLOM 3: Sedang Disewa -->

        <div class="bg-[#ffdeac]/20 border border-[#ffdeac] p-4 rounded-2xl flex flex-col">
          <div class="flex justify-between items-center mb-3 pb-2 border-b border-[#ffdeac]/60">
            <span
              class="text-xs font-extrabold text-[#7e5700] flex items-center gap-1.5 uppercase tracking-wider"
            >
              <span class="w-2.5 h-2.5 rounded-full bg-amber-500 animate-ping"></span> 3. Sedang
              Disewa
            </span>
            <span class="text-xs font-bold bg-[#ffdeac] text-[#7e5700] px-2 py-0.5 rounded-md">{{
              rentalStore.pickedUpList.length
            }}</span>
          </div>

          <div v-if="rentalStore.isLoading.active" class="space-y-3">
            <Skeleton height="7rem" rounded v-for="i in 2" :key="i" />
          </div>
          <div
            v-else-if="rentalStore.pickedUpList.length === 0"
            class="m-auto text-center py-12 text-xs text-[#58627c] italic"
          >
            Tidak ada box yang sedang keluar
          </div>
          <div v-else class="space-y-3 flex-1 overflow-y-auto max-h-[400px] pr-1">
            <div
              v-for="item in rentalStore.pickedUpList"
              :key="item.id"
              class="bg-white p-3.5 rounded-xl border border-amber-200/80 shadow-2xs hover:border-amber-500 transition-all"
            >
              <div class="flex justify-between items-start mb-1.5">
                <span
                  class="font-mono text-xs font-bold text-amber-800 bg-amber-50 px-1.5 py-0.5 rounded"
                  >{{ item.invoice_number }}</span
                >
                <span class="text-[11px] font-bold text-red-600"
                  >Kembali: {{ formatDateShort(item.end_date) }}</span
                >
              </div>
              <p class="font-bold text-sm text-[#1a1c20] truncate">
                {{ item.customer_snapshot_id }}
              </p>
              <div class="mt-3 pt-2 border-t border-amber-50 flex justify-between items-center">
                <span class="text-[11px] text-gray-400"
                  >Tagihan:
                  <strong class="text-gray-800">{{
                    formatCurrency(item.total_amount)
                  }}</strong></span
                >
                <div class="flex items-center gap-3">
                  <!-- TOMBOL BARU: UNDO / KEMBALIKAN KE ANTREAN (Lapisan Kuratif) -->
                  <Button
                    icon="pi pi-undo"
                    class="rounded-full w-5 h-5 bg-red-300 hover:bg-red-700 border-none text-white cursor-pointer transition-all duration-200 hover:scale-105"
                    title="Batalkan status Keluar (Kembalikan ke antrean)"
                    @click="safeConfirmRollback(item)"
                  />
                  <router-link to="/returns" title="Prose ke pengembalian">
                    <Button
                      label="Proses Retur"
                      icon="pi pi-check-circle"
                      class="rounded-full w-5 h-5 bg-green-300 hover:bg-green-700 border-none text-white cursor-pointer transition-all duration-200 hover:scale-105"
                    />
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRentalStore } from '@frontend/stores/rental'
import type { RentalReservation } from '@frontend/types/rental'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import ConfirmDialog from 'primevue/confirmdialog'
import Toast from 'primevue/toast'

const confirm = useConfirm()
const safeConfirmReady = (item: RentalReservation) => {
  confirm.require({
    header: 'Selesai Dekorasi?',
    message: `Pindahkan ${item.invoice_number} ke tahap "Menunggu Diambil"? Pastikan barang titipan pelanggan sudah lengkap.`,
    icon: 'pi pi-check text-slate-500',
    acceptLabel: 'Ya, Pindahkan',
    rejectLabel: 'Batal',
    acceptClass:
      'p-button-sm bg-slate-600 hover:bg-slate-700 text-white border-none font-bold px-5 py-2.5 rounded-xl shadow-md transition-all',
    rejectClass:
      'p-button-sm p-button-text text-gray-500 font-bold px-4 hover:bg-gray-100 rounded-xl transition-all',
    accept: async () => {
      try {
        await rentalStore.executeMarkReady(item.id)
        toast.add({ severity: 'success', summary: 'Sukses', detail: 'Siap diambil!', life: 2500 })
      } catch (e: any) {
        toast.add({
          severity: 'error',
          summary: 'Titipan Kosong',
          detail: rentalStore.errorMessage || 'Gagal memindah tahap',
          life: 5000,
        })
      }
    },
  })
}

const safeConfirmUndoReady = (item: RentalReservation) => {
  confirm.require({
    header: 'Batal Siap Diambil?',
    message: `Kembalikan ${item.invoice_number} ke Antrean Baru?`,
    icon: 'pi pi-undo text-slate-500',
    acceptLabel: 'Ya, Kembalikan',
    rejectLabel: 'Batal',
    acceptClass:
      'p-button-sm bg-slate-500 hover:bg-slate-600 text-white border-none font-bold px-5 py-2.5 rounded-xl shadow-md transition-all',
    rejectClass:
      'p-button-sm p-button-text text-gray-500 font-bold px-4 hover:bg-gray-100 rounded-xl transition-all',
    accept: async () => {
      try {
        await rentalStore.executeUndoReady(item.id)
        toast.add({
          severity: 'info',
          summary: 'Dibatalkan',
          detail: 'Kembali ke antrean awal',
          life: 2500,
        })
      } catch (e: any) {
        toast.add({
          severity: 'error',
          summary: 'Gagal',
          detail: rentalStore.errorMessage || 'Gagal rollback',
          life: 3000,
        })
      }
    },
  })
}

const safeConfirmPickup = (item: RentalReservation) => {
  confirm.require({
    header: 'Serahkan Box?',
    message: `Pastikan perlengkapan nota ${item.invoice_number} telah dibawa oleh ${item.customer_snapshot_id}. Lanjutkan serah terima?`,
    icon: 'pi pi-box text-blue-500',
    acceptLabel: 'Ya, Serahkan',
    rejectLabel: 'Batal',
    acceptClass:
      'p-button-sm bg-blue-600 hover:bg-blue-700 text-white border-none font-bold px-5 py-2.5 rounded-xl shadow-md transition-all',
    rejectClass:
      'p-button-sm p-button-text text-gray-500 font-bold px-4 hover:bg-gray-100 rounded-xl transition-all',
    accept: async () => {
      try {
        await rentalStore.executePickupUnit(item.id)
        toast.add({ severity: 'success', summary: 'Sukses', detail: 'Box diserahkan', life: 2500 })
      } catch (e: any) {
        toast.add({
          severity: 'error',
          summary: 'Gagal',
          detail: rentalStore.errorMessage || 'Sistem menolak mutasi',
          life: 3000,
        })
      }
    },
  })
}

// --- PENGAMAN LAPISAN 2 (Mengobati Salah Klik lewat Undo) ---
const safeConfirmRollback = (item: RentalReservation) => {
  confirm.require({
    header: 'Batalkan Keluar?',
    message: `Kembalikan pesanan ${item.invoice_number} ke antrean awal "Menunggu Diambil"?`,
    icon: 'pi pi-undo text-amber-500',
    acceptLabel: 'Ya, Tarik Balik',
    rejectLabel: 'Batal',
    acceptClass:
      'p-button-sm bg-amber-500 hover:bg-amber-600 text-white border-none font-bold px-5 py-2.5 rounded-xl shadow-md transition-all',
    rejectClass:
      'p-button-sm p-button-text text-gray-500 font-bold px-4 hover:bg-gray-100 rounded-xl transition-all',
    accept: async () => {
      try {
        await rentalStore.executeRollbackPickup(item.id)
        toast.add({
          severity: 'info',
          summary: 'Dibatalkan',
          detail: 'Nota ditarik kembali',
          life: 2500,
        })
      } catch (e: any) {
        toast.add({
          severity: 'error',
          summary: 'Gagal',
          detail: rentalStore.errorMessage || 'Gagal rollback',
          life: 3000,
        })
      }
    },
  })
}

const toast = useToast()
const rentalStore = useRentalStore()
const refreshDashboard = async () => {
  await Promise.all([
    rentalStore.fetchActiveReservations(),
    rentalStore.fetchUpcomingReservations(),
    rentalStore.fetchOverdueReservations(),
    rentalStore.fetchProducts(),
  ])
}

onMounted(() => {
  refreshDashboard()
})
const formatCurrency = (val: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
const formatDateShort = (d: string) =>
  new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
const formatDateWithDay = (d: string) =>
  new Date(d)
    .toLocaleDateString('id-ID', { weekday: 'short', day: 'numeric', month: 'short' })
    .toUpperCase()
const calculateDays = (s: string, e: string) =>
  Math.max(1, Math.round((new Date(e).getTime() - new Date(s).getTime()) / (1000 * 3600 * 24)))
</script>
