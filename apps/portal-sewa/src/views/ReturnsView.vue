<template>
  <div class="space-y-4 pb-8 font-sans text-[#1a1c20] w-full max-w-[1600px] mx-auto px-2">
    <Toast />

    <!-- HEADER: Dibuat Lebih Ringkas & Hemat Ruang Vertikal -->
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-[#c4c6d1]/60 pb-3"
    >
      <div>
        <p class="text-[10px] font-black tracking-widest text-[#ba1a1a] uppercase mb-0.5">
          Kliring & Penalti
        </p>
        <h2 class="text-xl font-black text-[#1a1c20] tracking-tight flex items-center gap-2">
          <i class="pi pi-box text-[#ba1a1a] text-base"></i>
          Pengembalian Aset (Return)
        </h2>
      </div>
      <div class="flex items-center gap-3">
        <span class="p-input-icon-left w-full sm:w-72">
          <i class="pi pi-search text-gray-400 text-xs" />
          <InputText
            v-model="searchQuery"
            placeholder="Cari Nama atau Invoice..."
            class="w-full rounded-xl border-gray-300 shadow-xs p-inputtext-sm text-xs pl-8"
          />
        </span>
        <Button
          label="Segarkan"
          icon="pi pi-refresh"
          class="p-button-sm bg-white hover:bg-gray-50 text-[#58627c] border border-gray-200 rounded-xl shadow-xs font-bold px-3 py-2 transition-all text-xs"
          :loading="rentalStore.isLoading.active"
          @click="rentalStore.fetchActiveReservations()"
        />
      </div>
    </div>

    <!-- MASTER DETAIL LAYOUT: Full View Optimized -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">
      <!-- LEFT PANEL: MASTER DATATABLE (col-span-3 untuk memperluas area form kanan) -->
      <div
        class="lg:col-span-3 bg-white border border-[#c4c6d1]/60 rounded-2xl p-4 shadow-xs min-h-[450px]"
      >
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <span
              class="w-6 h-6 rounded-lg bg-orange-50 text-orange-600 flex items-center justify-center font-bold border border-orange-100"
            >
              <i class="pi pi-list text-xs"></i>
            </span>
            <h3 class="font-bold text-[#1a1c20] text-xs">Daftar Sewa Aktif</h3>
          </div>
          <span class="text-[9px] font-bold text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full"
            >{{ filteredList.length }} Kontrak</span
          >
        </div>

        <DataTable
          :value="filteredList"
          :paginator="true"
          :rows="6"
          class="p-datatable-sm custom-table-returns text-xs"
          responsiveLayout="scroll"
        >
          <template #empty>
            <div
              class="p-6 text-center text-gray-400 flex flex-col items-center border border-dashed border-gray-200 rounded-xl my-2"
            >
              <i class="pi pi-check-square text-2xl mb-1 text-gray-300"></i>
              <p class="font-bold text-gray-600 text-[10px]">Tidak ada kontrak aktif</p>
            </div>
          </template>

          <Column header="Invoice & Klien" class="py-1">
            <template #body="{ data }">
              <div class="flex flex-col gap-0.5">
                <span class="font-bold text-gray-800 text-[11px] truncate max-w-[140px]">{{
                  parseCustomer(data.customer_snapshot_id).name
                }}</span>
                <span
                  class="text-[9px] font-mono font-bold text-orange-700 bg-orange-50 w-max px-1 rounded border border-orange-200"
                  >{{ data.invoice_number }}</span
                >
                <span v-if="isOverdue(data.end_date)" class="text-[9px] font-bold text-red-600 mt-0.5 flex items-center gap-0.5"
                  ><i class="pi pi-exclamation-triangle text-[8px]"></i> Telat</span
                >
              </div>
            </template>
          </Column>

          <Column header="Aksi" class="w-12 text-right">
            <template #body="{ data }">
              <Button
                v-if="selectedActive?.id !== data.id"
                icon="pi pi-chevron-right"
                class="p-button-rounded p-button-text p-button-sm text-blue-600 hover:bg-blue-50 w-7 h-7"
                @click="selectReservation(data)"
              />
              <span v-else class="text-xs font-bold text-emerald-600"
                ><i class="pi pi-check-circle"></i
              ></span>
            </template>
          </Column>
        </DataTable>
      </div>

      <!-- RIGHT PANEL: SIDE-BY-SIDE SIDEBAR FORM (col-span-9) -->
      <div class="lg:col-span-9 space-y-4">
        <!-- Empty State Panel -->
        <div
          v-if="!selectedActive"
          class="bg-slate-50 border border-slate-200 border-dashed rounded-2xl p-8 text-center flex flex-col items-center justify-center min-h-[450px] text-slate-400"
        >
          <i class="pi pi-arrow-left text-3xl mb-3 text-slate-300"></i>
          <h3 class="font-bold text-slate-600 text-sm mb-1">Pilih Kontrak di Sisi Kiri</h3>
          <p class="text-xs max-w-sm mx-auto text-slate-500">
            Silakan pilih salah satu data pemesanan aktif untuk melakukan pemeriksaan logistik fisik, entri penalti cacat, dan kalkulasi pembayaran akhir.
          </p>
        </div>

        <!-- Form Panel: Dioptimalkan Menjadi Dua Kolom Utama Menyamping Sepenuhnya -->
        <div
          v-else
          class="bg-[#1a1c20] text-white rounded-2xl shadow-xl overflow-hidden flex flex-col min-h-[450px] border border-gray-800 animate-fade-in relative"
        >
          <!-- Top Section Profile Header -->
          <div
            class="p-4 border-b border-gray-800 bg-gray-900/50 flex items-center justify-between"
          >
            <div class="flex items-center gap-4">
              <div class="flex flex-col">
                <span class="text-[9px] font-black tracking-widest text-orange-400 uppercase mb-0.5"
                  >Panel Inspeksi Pengembalian</span
                >
                <div class="flex items-baseline gap-2">
                  <span class="text-lg font-black text-white leading-none">{{
                    parseCustomer(selectedActive.customer_snapshot_id).name
                  }}</span>
                  <span class="text-xs font-mono font-bold text-orange-400 bg-orange-950/40 px-1.5 py-0.5 rounded border border-orange-900/40">{{
                    selectedActive.invoice_number
                  }}</span>
                </div>
                <span class="text-[10px] text-gray-400 mt-1 flex items-center gap-1"
                  ><i class="pi pi-whatsapp text-green-500 text-[9px]"></i>
                  {{ parseCustomer(selectedActive.customer_snapshot_id).phone }}</span
                >
              </div>
            </div>

            <div class="flex items-center gap-3">
              <div class="bg-gray-800/90 rounded-xl px-3 py-1 border border-gray-700 text-right flex items-center gap-3">
                <div>
                  <p class="text-[8px] text-gray-400 uppercase tracking-widest font-bold">Keterlambatan</p>
                  <p class="font-black text-sm" :class="autoCalculatedLateDays > 0 ? 'text-red-400' : 'text-emerald-400'">
                    {{ autoCalculatedLateDays }} Hari
                  </p>
                </div>
              </div>
              <button
                @click="selectedActive = null"
                class="w-7 h-7 rounded-full bg-gray-800 hover:bg-gray-700 flex items-center justify-center text-gray-400 transition-colors cursor-pointer"
              >
                <i class="pi pi-times text-xs"></i>
              </button>
            </div>
          </div>

          <!-- LAYOUT UTAMA KIRI & KANAN (Mencegah Menumpuk Kebawah) -->
          <div class="p-5 grid grid-cols-1 md:grid-cols-12 gap-6 items-start flex-1 bg-gradient-to-b from-transparent to-gray-900/40">
            
            <!-- SUB-PANEL KIRI: LOGISTIK & FOTO (col-span-6) -->
            <div class="md:col-span-6 space-y-4">
              <!-- List Barang Sewa -->
              <div class="space-y-1.5">
                <label class="text-[10px] font-extrabold text-orange-400 uppercase tracking-widest block"
                  >Daftar Aset Sewa</label
                >
                <div class="bg-gray-800/40 rounded-xl border border-gray-700/40 p-1.5 max-h-[110px] overflow-y-auto">
                  <table class="w-full text-[11px] text-left">
                    <tbody>
                      <tr
                        v-for="item in selectedActive.items"
                        :key="item.id"
                        class="border-b border-gray-800 last:border-0"
                      >
                        <td class="py-1 pl-1.5 text-gray-300">{{ item.rental_product_name }}</td>
                        <td class="py-1 pr-1.5 text-right font-bold text-orange-200">{{ item.qty }} Unit</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- List Barang Titipan (Jika ada) -->
              <div v-if="selectedActive.contents && selectedActive.contents.length > 0" class="space-y-1.5">
                <label class="text-[10px] font-extrabold text-blue-400 uppercase tracking-widest block"
                  >Daftar Barang Titipan (Deposit)</label
                >
                <div class="bg-blue-950/20 rounded-xl border border-blue-900/30 p-1.5 max-h-[90px] overflow-y-auto">
                  <table class="w-full text-[11px] text-left">
                    <tbody>
                      <tr
                        v-for="c in selectedActive.contents"
                        :key="c.id"
                        class="border-b border-blue-950 last:border-0"
                      >
                        <td class="py-1 pl-1.5 text-blue-200">{{ c.item_name }}</td>
                        <td class="py-1 pr-1.5 text-right text-blue-300 font-bold">{{ c.quantity }} Pcs</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Upload Foto -->
              <div class="space-y-1.5">
                <label class="text-[10px] font-extrabold text-gray-300 uppercase tracking-widest block"
                  >Unggah Foto Bukti Fisik</label
                >
                <div class="relative group cursor-pointer">
                  <input
                    type="file"
                    accept="image/*"
                    @change="onFileSelected"
                    class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-20"
                  />
                  <div v-if="previewUrl" class="w-full relative rounded-xl overflow-hidden border border-gray-700 bg-black/40">
                    <img :src="previewUrl" alt="Preview Foto" class="w-full h-24 object-cover opacity-90" />
                    <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity">
                      <span class="text-white font-bold text-[10px] bg-black/60 px-2 py-1 rounded"><i class="pi pi-sync mr-1"></i>Ganti Foto</span>
                    </div>
                  </div>
                  <div v-else class="w-full flex items-center gap-2 bg-gray-800/60 border border-dashed border-gray-700 rounded-xl p-2.5 transition-colors group-hover:border-orange-500/60">
                    <span class="w-7 h-7 rounded-full bg-gray-700 flex items-center justify-center text-gray-400 text-xs"><i class="pi pi-camera"></i></span>
                    <div class="flex-1 overflow-hidden">
                      <p class="text-xs font-bold text-gray-300 truncate">{{ selectedFileName || 'Pilih Gambar Bukti Fisik' }}</p>
                      <p class="text-[9px] text-gray-500">Maks. 10MB (JPG/PNG)</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- SUB-PANEL KANAN: DENDA & INPUT MANUAL (col-span-6) -->
            <div class="md:col-span-6 space-y-3">
              <label class="text-[10px] font-black text-gray-400 uppercase tracking-widest block"
                >Parameter Penalti & Kompensasi</label
              >

              <div class="grid grid-cols-2 gap-3">
                <div class="flex flex-col gap-1">
                  <label class="text-[10px] font-extrabold text-red-400">Denda Terlambat (Auto)</label>
                  <div class="p-inputgroup h-8">
                    <span class="p-inputgroup-addon bg-gray-800 border-gray-700 text-gray-400 text-[10px] px-2 font-bold border-r-0 rounded-l-lg">Rp</span>
                    <InputNumber v-model="lateFee" :disabled="true" class="p-inputtext-sm font-bold text-red-400 bg-gray-800/40 border-gray-700 rounded-r-lg" inputClass="p-1 text-xs" />
                  </div>
                </div>

                <div class="flex flex-col gap-1">
                  <label class="text-[10px] font-extrabold text-orange-400">Denda Kerusakan (Manual)</label>
                  <div class="p-inputgroup h-8">
                    <span class="p-inputgroup-addon bg-gray-800 border-gray-700 text-gray-400 text-[10px] px-2 font-bold border-r-0 rounded-l-lg">Rp</span>
                    <InputNumber v-model="damageFee" class="p-inputtext-sm font-bold text-white bg-gray-800 border-gray-700 rounded-r-lg focus:border-orange-500" inputClass="p-1 text-xs" />
                  </div>
                </div>
              </div>

              <!-- Pilih Aset Rusak -->
              <div class="flex flex-col gap-1">
                <label class="text-[10px] font-bold text-gray-400">Pilih Aset Rusak (Jika Ada)</label>
                <select
                  v-model="damagedItemId"
                  class="p-inputtext-sm rounded-lg border-gray-700 bg-gray-800 text-white text-xs p-1.5 focus:border-orange-500 outline-none"
                >
                  <option value="">-- Tetapkan ke Aset Pertama (Default) --</option>
                  <option v-for="item in selectedActive.items" :key="item.id" :value="item.rental_product_id">
                    {{ item.rental_product_name }}
                  </option>
                </select>
              </div>

              <!-- Catatan Inspeksi -->
              <div class="flex flex-col gap-1">
                <label class="text-[10px] font-bold text-gray-400">Catatan Kondisi Khusus Fisik</label>
                <Textarea
                  v-model="returnNotes"
                  rows="2"
                  placeholder="Kondisi kelengkapan box fisik..."
                  class="p-inputtext-sm rounded-lg border-gray-700 bg-gray-800 text-white placeholder-gray-500 text-xs p-1.5"
                />
              </div>

              <!-- Input Kasir Tunai -->
              <div class="flex flex-col gap-1 pt-1 border-t border-gray-800/60">
                <label class="text-[10px] font-extrabold text-emerald-400 uppercase tracking-widest"
                  >Nominal Uang Diterima Kasir</label
                >
                <div class="p-inputgroup h-9">
                  <span class="p-inputgroup-addon bg-gray-800 border-gray-700 text-gray-400 text-xs px-2.5 font-bold border-r-0 rounded-l-lg">Rp</span>
                  <InputNumber
                    v-model="amountPaid"
                    class="p-inputtext-sm font-black text-emerald-400 bg-gray-800 border-gray-700 rounded-r-lg focus:border-emerald-500"
                    inputClass="p-1.5 text-sm"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- FOOTER CALCULATOR: Bagian Penghitungan Finansial Sejajar Horisontal -->
          <div
            class="p-4 bg-black border-t border-gray-800 relative z-10 flex flex-col sm:flex-row items-center justify-between gap-4"
          >
            <!-- Kolom Neraca Pembukuan (Horizontal flex) -->
            <div class="w-full sm:w-auto flex flex-wrap items-center gap-x-6 gap-y-1 text-xs text-gray-400 font-medium">
              <div class="flex items-center gap-2">
                <span>Sisa Pokok:</span>
                <span class="font-bold text-white">{{ formatCurrency(selectedActive.total_amount - selectedActive.down_payment) }}</span>
              </div>
              <div class="flex items-center gap-2 text-orange-400">
                <span>Total Penalti:</span>
                <span class="font-bold">+{{ formatCurrency(lateFee + damageFee) }}</span>
              </div>
              <div class="flex items-center gap-2 text-emerald-400" v-if="amountPaid > 0">
                <span>Kembalian:</span>
                <span class="font-bold font-mono">{{ formatCurrency(changeAmount) }}</span>
              </div>
              
              <!-- Grand Total Highlighter -->
              <div class="sm:border-l sm:border-gray-800 sm:pl-4 flex items-center gap-2">
                <span class="text-[10px] text-gray-400 font-bold uppercase tracking-wider">Wajib Dibayar:</span>
                <span class="text-lg font-black text-white font-mono bg-gray-900 px-2 py-0.5 rounded border border-gray-800">{{ formatCurrency(totalBill) }}</span>
              </div>
            </div>

            <!-- Tombol Utama Eksekusi Kunci -->
            <Button
              label="Tarik Aset & Selesaikan"
              icon="pi pi-check-circle"
              class="w-full sm:w-auto bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white border-none font-black text-xs py-2.5 px-5 rounded-xl transition-all shadow-md cursor-pointer"
              :loading="rentalStore.isLoading.mutate"
              @click="executeReturn"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- HIDDEN INVOICE TEMPLATE FOR PDF GENERATION -->
    <div class="fixed pointer-events-none" style="top: 0; left: -9999px; width: 800px; z-index: -1">
      <div id="invoice-template" class="w-[800px] p-10 bg-white text-black font-sans">
        <div class="text-center border-b-2 border-gray-800 pb-4 mb-6">
          <h1 class="text-3xl font-black uppercase">RINZI RENTAL</h1>
          <p class="text-gray-600 text-sm">INVOICE Pengembalian</p>
        </div>

        <div class="flex justify-between mb-8">
          <div>
            <p class="text-xs font-bold text-gray-500 uppercase">Pelanggan</p>
            <p class="font-bold text-lg">
              {{ selectedActive ? parseCustomer(selectedActive.customer_snapshot_id).name : '-' }}
            </p>
            <p class="text-sm">
              {{ selectedActive ? parseCustomer(selectedActive.customer_snapshot_id).phone : '-' }}
            </p>
          </div>
          <div class="text-right">
            <p class="text-xs font-bold text-gray-500 uppercase">No. Invoice</p>
            <p class="font-bold text-lg">{{ selectedActive?.invoice_number ?? '-' }}</p>
            <p class="text-sm">Tanggal Kembali: {{ new Date().toLocaleDateString('id-ID') }}</p>
          </div>
        </div>

        <table class="w-full mb-8 border-collapse">
          <thead>
            <tr class="bg-gray-200">
              <th class="p-3 text-left text-sm font-bold border border-gray-300">Deskripsi Barang Sewa</th>
              <th class="p-3 text-right text-sm font-bold border border-gray-300">Qty</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in selectedActive?.items ?? []" :key="item.id">
              <td class="p-3 text-sm border border-gray-300">{{ item.rental_product_name }}</td>
              <td class="p-3 text-sm text-right border border-gray-300">{{ item.qty }}</td>
            </tr>
          </tbody>
        </table>

        <div class="flex justify-end">
          <div class="w-1/2">
            <div class="flex justify-between py-2 border-b border-gray-200">
              <span class="text-sm font-bold">Keterlambatan ({{ autoCalculatedLateDays }} Hari):</span>
              <span class="text-sm">{{ formatCurrency(lateFee) }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-gray-200">
              <span class="text-sm font-bold">Denda Kerusakan:</span>
              <span class="text-sm">{{ formatCurrency(damageFee) }}</span>
            </div>
            <div class="flex justify-between py-3 border-b-2 border-gray-800">
              <span class="text-lg font-black uppercase">Total Denda &amp; Penalti:</span>
              <span class="text-lg font-black">{{ formatCurrency(lateFee + damageFee) }}</span>
            </div>
            <div class="flex justify-between py-3">
              <span class="text-lg font-black uppercase text-red-600">Sisa Tagihan Pokok:</span>
              <span class="text-lg font-black">{{ formatCurrency((selectedActive?.total_amount ?? 0) - (selectedActive?.down_payment ?? 0)) }}</span>
            </div>
            <div class="flex justify-between py-3 bg-gray-100 p-2 mt-2">
              <span class="text-xl font-black uppercase">Grand Total Dibayar:</span>
              <span class="text-xl font-black">{{ formatCurrency((selectedActive?.total_amount ?? 0) - (selectedActive?.down_payment ?? 0) + lateFee + damageFee) }}</span>
            </div>
          </div>
        </div>

        <div class="mt-12 pt-8 border-t border-gray-300 text-center">
          <p class="text-xs text-gray-500 italic">
            Terima kasih atas kepercayaan Anda menyewa aset di Rinzi Rental. Dokumen ini sah sebagai tanda serah terima pengembalian.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRentalStore } from '@frontend/stores/rental'
import type { RentalReservation } from '@frontend/types/rental'
import { useToast } from 'primevue/usetoast'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'
import Toast from 'primevue/toast'

import { jsPDF } from 'jspdf'

const toast = useToast()
const rentalStore = useRentalStore()
const searchQuery = ref('')
const selectedActive = ref<RentalReservation | null>(null)
const damageFee = ref(0)
const returnNotes = ref('')
const damagedItemId = ref('')

const returnPhotoFile = ref<File | null>(null)
const selectedFileName = ref('')
const previewUrl = ref<string | null>(null)
const amountPaid = ref(0)

onMounted(() => {
  rentalStore.fetchActiveReservations()
})

const parseCustomer = (str?: string) => {
  if (!str) return { name: 'Tanpa Nama', phone: '-' }
  const parts = str.split(' (')
  return { name: parts[0] || 'Tanpa Nama', phone: parts[1] ? parts[1].replace(')', '') : '-' }
}

const filteredList = computed(() => {
  let list = rentalStore.pickedUpList
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(
      (r) =>
        r.invoice_number.toLowerCase().includes(q) ||
        parseCustomer(r.customer_snapshot_id).name.toLowerCase().includes(q),
    )
  }
  return list
})

const selectReservation = (data: RentalReservation) => {
  selectedActive.value = data
  damageFee.value = 0
  returnNotes.value = ''
  damagedItemId.value = ''
  amountPaid.value = 0
  returnPhotoFile.value = null
  selectedFileName.value = ''
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = null
  }
}

const onFileSelected = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    if (!file) return
    if (file.size > 10 * 1024 * 1024) {
      toast.add({
        severity: 'warn',
        summary: 'File Terlalu Besar',
        detail: 'Maksimal ukuran foto adalah 10MB',
        life: 3000,
      })
      target.value = ''
      return
    }
    returnPhotoFile.value = file
    selectedFileName.value = file.name

    if (previewUrl.value) {
      URL.revokeObjectURL(previewUrl.value)
    }
    previewUrl.value = URL.createObjectURL(file)
  }
}

const autoCalculatedLateDays = computed(() => {
  if (!selectedActive.value) return 0
  const end = new Date(selectedActive.value.end_date)
  end.setHours(0, 0, 0, 0)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const diff = Math.floor((today.getTime() - end.getTime()) / (1000 * 3600 * 24))
  return diff > 0 ? diff : 0
})

const isOverdue = (endDateStr: string) => {
  const end = new Date(endDateStr)
  end.setHours(0, 0, 0, 0)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return today.getTime() > end.getTime()
}

const lateFee = computed(() => autoCalculatedLateDays.value * 50000)

const totalBill = computed(() => {
  if (!selectedActive.value) return 0
  return selectedActive.value.total_amount - selectedActive.value.down_payment + lateFee.value + damageFee.value
})

const changeAmount = computed(() => {
  return Math.max(0, amountPaid.value - totalBill.value)
})

const generateAndUploadPdf = async (returnId: string, invoiceNumber: string, customerNameRaw: string) => {
  const customerName = parseCustomer(customerNameRaw).name.replace(/[^a-zA-Z0-9]/g, '_')
  const safeInvoice = invoiceNumber.replace(/[\/\\]/g, '_')

  const doc = new jsPDF('p', 'mm', 'a4')
  const pageWidth = doc.internal.pageSize.getWidth()

  doc.setFillColor(15, 23, 42)
  doc.rect(0, 0, pageWidth, 42, 'F')
  
  doc.setTextColor(255, 255, 255)
  doc.setFont('helvetica', 'bold')
  doc.setFontSize(22)
  doc.text('INVOICE RETURN', 15, 18)
  
  doc.setFont('helvetica', 'normal')
  doc.setFontSize(10)
  doc.setTextColor(148, 163, 184)
  doc.text('Bukti Penyelesaian Kontrak, Kliring Penalti & Serah Terima Aset', 15, 26)

  doc.setFont('helvetica', 'bold')
  doc.setFontSize(11)
  doc.setTextColor(255, 255, 255)
  doc.text(`NO. NOTA: ${invoiceNumber}`, pageWidth - 15, 18, { align: 'right' })
  
  doc.setFont('helvetica', 'normal')
  doc.setFontSize(9)
  doc.setTextColor(148, 163, 184)
  doc.text(`Tanggal Cetak: ${new Date().toLocaleDateString('id-ID')}`, pageWidth - 15, 26, { align: 'right' })

  let currentY = 52
  doc.setDrawColor(226, 232, 240)
  doc.setFillColor(248, 250, 252)
  doc.rect(15, currentY, pageWidth - 30, 18, 'FD')

  doc.setTextColor(100, 116, 139)
  doc.setFontSize(8)
  doc.setFont('helvetica', 'bold')
  doc.text('NAMA PENYEWA', 20, currentY + 5.5)
  doc.setTextColor(30, 41, 59)
  doc.setFontSize(10)
  doc.text(parseCustomer(customerNameRaw).name, 20, currentY + 12.5)

  doc.setTextColor(100, 116, 139)
  doc.setFontSize(8)
  doc.setFont('helvetica', 'bold')
  doc.text('KONTAK WHATSAPP', 95, currentY + 5.5)
  doc.setTextColor(30, 41, 59)
  doc.setFontSize(10)
  doc.text(parseCustomer(customerNameRaw).phone, 95, currentY + 12.5)

  doc.setTextColor(100, 116, 139)
  doc.setFontSize(8)
  doc.setFont('helvetica', 'bold')
  doc.text('STATUS DOKUMEN', 155, currentY + 5.5)
  doc.setTextColor(22, 163, 74)
  doc.setFontSize(10)
  doc.setFont('helvetica', 'bold')
  doc.text('CLOSED / LUNAS', 155, currentY + 12.5)

  currentY = 82
  doc.setTextColor(15, 23, 42)
  doc.setFontSize(10)
  doc.setFont('helvetica', 'bold')
  doc.text('I. RINCIAN PAKET TERPESAN DI DALAM NOTA', 15, currentY)
  
  currentY += 4
  doc.setFillColor(15, 23, 42)
  doc.rect(15, currentY, pageWidth - 30, 8, 'F')
  
  doc.setTextColor(255, 255, 255)
  doc.setFontSize(9)
  doc.text('Nama Barang / Box', 20, currentY + 5.5)
  doc.text('Qty', 110, currentY + 5.5, { align: 'center' })
  doc.text('Harga Satuan', 145, currentY + 5.5, { align: 'right' })
  doc.text('Subtotal', pageWidth - 20, currentY + 5.5, { align: 'right' })

  currentY += 8
  doc.setFont('helvetica', 'normal')
  doc.setTextColor(51, 65, 85)
  
  const itemsList = selectedActive.value?.items ?? []
  itemsList.forEach((item, index) => {
    if (index % 2 === 0) {
      doc.setFillColor(250, 250, 250)
      doc.rect(15, currentY, pageWidth - 30, 7.5, 'F')
    }
    doc.setDrawColor(241, 245, 249)
    doc.line(15, currentY + 7.5, pageWidth - 15, currentY + 7.5)

    doc.text(item.rental_product_name, 20, currentY + 5)
    doc.text(item.qty.toString(), 110, currentY + 5, { align: 'center' })
    doc.text(`Rp ${item.price_per_period.toLocaleString('id-ID')}`, 145, currentY + 5, { align: 'right' })
    doc.text(`Rp ${item.subtotal.toLocaleString('id-ID')}`, pageWidth - 20, currentY + 5, { align: 'right' })
    currentY += 7.5
  })

  currentY += 12
  doc.setTextColor(15, 23, 42)
  doc.setFontSize(10)
  doc.setFont('helvetica', 'bold')
  doc.text('II. PERHITUNGAN AKUNTANSI FINANSIAL & PENALTI', 15, currentY)

  currentY += 4
  doc.setFillColor(51, 65, 85)
  doc.rect(15, currentY, pageWidth - 30, 8, 'F')

  doc.setTextColor(255, 255, 255)
  doc.setFontSize(9)
  doc.text('Komponen Pembukuan Kas / Kliring Gudang', 20, currentY + 5.5)
  doc.text('Sifat Arus', 125, currentY + 5.5, { align: 'center' })
  doc.text('Nilai Buku Positif / Negatif', pageWidth - 20, currentY + 5.5, { align: 'right' })

  const ledgerEntries = [
    { label: 'Total Pokok Sewa Box', op: 'Debit Awal', val: selectedActive.value?.total_amount ?? 0 },
    { label: `Denda Keterlambatan (${autoCalculatedLateDays.value} Hari)`, op: '(+) Penalti', val: lateFee.value },
    { label: 'Denda Kerusakan / Cacat Fisik', op: '(+) Penalti', val: damageFee.value || 0 },
    { label: 'Total Tagihan Kumulatif', op: '(=) Gross', val: (selectedActive.value?.total_amount ?? 0) + lateFee.value + (damageFee.value || 0) },
    { label: 'Uang Muka Terbayar (DP Di Muka)', op: '(-) Kredit', val: selectedActive.value?.down_payment ?? 0 },
    { label: 'Tagihan Akhir (Kekurangan Kliring)', op: '(=) Net Due', val: totalBill.value },
    { label: 'Dibayarkan (Pelunasan Tunai/Transfer)', op: 'Setoran Kas', val: amountPaid.value || 0 },
    { label: 'Kembalian Penyewa', op: 'Sisa Refund', val: changeAmount.value || 0 }
  ]

  currentY += 8
  ledgerEntries.forEach((entry, idx) => {
    if (idx === 3 || idx === 5 || idx === 7) {
      doc.setFillColor(241, 245, 249)
      doc.rect(15, currentY, pageWidth - 30, 8, 'F')
    }

    doc.setDrawColor(226, 232, 240)
    doc.line(15, currentY + 8, pageWidth - 15, currentY + 8)

    if (idx === 3 || idx === 5 || idx === 7) {
      doc.setFont('helvetica', 'bold')
      doc.setTextColor(15, 23, 42)
    } else {
      doc.setFont('helvetica', 'normal')
      doc.setTextColor(71, 85, 105)
    }

    doc.text(entry.label, 20, currentY + 5.5)
    doc.text(entry.op, 125, currentY + 5.5, { align: 'center' })
    
    const sign = entry.op.includes('(-)') ? '-' : ''
    doc.text(`${sign}Rp ${entry.val.toLocaleString('id-ID')}`, pageWidth - 20, currentY + 5.5, { align: 'right' })
    
    currentY += 8
  })

  if (returnNotes.value) {
    currentY += 10
    doc.setDrawColor(226, 232, 240)
    doc.setFillColor(254, 254, 255)
    doc.rect(15, currentY, pageWidth - 30, 16, 'D')
    
    doc.setFont('helvetica', 'bold')
    doc.setFontSize(8)
    doc.setTextColor(100, 116, 139)
    doc.text('CATATAN KONDISI LOGISTIK FISIK BARANG:', 19, currentY + 5)
    
    doc.setFont('helvetica', 'normal')
    doc.setFontSize(9)
    doc.setTextColor(51, 65, 85)
    doc.text(returnNotes.value, 19, currentY + 11, { maxWidth: pageWidth - 38 })
  }

  currentY += 24
  doc.setFontSize(9)
  doc.setFont('helvetica', 'normal')
  doc.setTextColor(100, 116, 139)
  
  doc.text('Penyewa / Customer,', 40, currentY, { align: 'center' })
  doc.line(20, currentY + 18, 60, currentY + 18)
  doc.setFont('helvetica', 'bold')
  doc.setTextColor(15, 23, 42)
  doc.text(parseCustomer(customerNameRaw).name, 40, currentY + 23, { align: 'center' })

  doc.setFont('helvetica', 'normal')
  doc.setTextColor(100, 116, 139)
  doc.text('Kasir Otorisasi / Gudang,', pageWidth - 40, currentY, { align: 'center' })
  doc.line(pageWidth - 60, currentY + 18, pageWidth - 20, currentY + 18)
  doc.setFont('helvetica', 'bold')
  doc.setTextColor(15, 23, 42)
  doc.text('Sistem Arsip Valid', pageWidth - 40, currentY + 23, { align: 'center' })

  const footerPageY = 282
  doc.setDrawColor(226, 232, 240)
  doc.line(15, footerPageY - 6, pageWidth - 15, footerPageY - 6)
  doc.setFont('helvetica', 'normal')
  doc.setFontSize(8)
  doc.setTextColor(148, 163, 184)
  doc.text('Seluruh transaksi pengembalian ini telah tervalidasi secara komparatif oleh sistem pergudangan Rinzi Rental.', pageWidth / 2, footerPageY, { align: 'center' })

  const pdfBlob = doc.output('blob')
  const fileName = `${customerName}_${safeInvoice}.pdf`

  doc.save(`${customerName}_${safeInvoice}.pdf`)
  await rentalStore.uploadReturnReceipt(returnId, pdfBlob, fileName)
}

const executeReturn = async () => {
  if (!selectedActive.value) return

  const snapshotInvoice = selectedActive.value.invoice_number
  const snapshotReservationId = selectedActive.value.id
  const snapshotCustomer = selectedActive.value.customer_snapshot_id

  try {
    const returnItemsPayload = (selectedActive.value.items || []).map((item: any, index: number) => {
      const isDamaged = damagedItemId.value ? item.rental_product_id === damagedItemId.value : index === 0;
      return {
        rental_product_id: item.rental_product_id,
        qty_returned: item.qty,
        condition_status: (isDamaged && damageFee.value > 0) ? 'DAMAGED' : 'GOOD',
        damage_fee: isDamaged ? (damageFee.value || 0) : 0,
        condition_notes: isDamaged ? returnNotes.value : ''
      };
    })

    const returnResponse = await rentalStore.executeProcessReturnKliring(
      {
        reservation_id: snapshotReservationId,
        amount_paid: amountPaid.value,
        change_amount: changeAmount.value,
        manual_damage_fee: damageFee.value || 0,
        manual_return_notes: returnNotes.value,
        return_items: returnItemsPayload,
      },
      (returnPhotoFile.value as unknown as File) || undefined,
    )

    toast.add({
      severity: 'success',
      summary: 'Aset Berhasil Ditarik',
      detail: `Kontrak ${snapshotInvoice} telah diselesaikan!`,
      life: 4000,
    })

    const returnData: any = returnResponse
    const returnId: string | undefined = returnData?.id ?? returnData?.data?.id

    if (returnId) {
      try {
        await generateAndUploadPdf(returnId, snapshotInvoice, snapshotCustomer)
      } catch (pdfErr: any) {
        console.error('PDF Return Error:', pdfErr)
      }
    }

    await rentalStore.fetchActiveReservations()
    selectedActive.value = null
    damageFee.value = 0
    returnNotes.value = ''
    returnPhotoFile.value = null
    selectedFileName.value = ''
    if (previewUrl.value) {
      URL.revokeObjectURL(previewUrl.value)
      previewUrl.value = null
    }
  } catch (e: any) {
    toast.add({
      severity: 'error',
      summary: 'Gagal Eksekusi',
      detail: rentalStore.errorMessage || e.response?.data?.message || 'Gagal menutup kontrak kliring',
      life: 4000,
    })
  }
}

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
</script>

<style scoped>
:deep(.custom-table-returns .p-datatable-thead > tr > th) {
  background-color: #f8fafc;
  color: #64748b;
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0.5rem;
  border-bottom: 1px solid #e2e8f0;
}
:deep(.custom-table-returns .p-datatable-tbody > tr > td) {
  padding: 0.5rem;
  font-size: 0.8rem;
  border-bottom: 1px solid #f1f5f9;
}
:deep(.custom-table-returns .p-datatable-tbody > tr:hover) {
  background-color: #f8fafc;
}
.animate-fade-in {
  animation: fadeIn 0.2s ease-out forwards;
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>