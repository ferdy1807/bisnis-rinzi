<template>
  <div class="max-w-[1400px] mx-auto space-y-6 pb-16 font-sans text-slate-800 antialiased">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <p class="text-[11px] font-bold tracking-widest text-indigo-600 uppercase mb-1">
          Pemesanan & Invoice
        </p>
        <h2 class="text-2xl font-extrabold text-slate-900 tracking-tight flex items-center gap-2.5">
          <i class="pi pi-file-edit text-indigo-600 text-xl bg-indigo-50 p-2 rounded-xl"></i>
          Formulir Kontrak Baru
        </h2>
      </div>
      <router-link to="/reservations-list" class="no-underline">
        <Button
          label="Kembali ke Daftar"
          icon="pi pi-arrow-left"
          class="p-button-outlined p-button-sm text-slate-600 border-slate-300 hover:bg-slate-50 font-semibold rounded-xl px-4 py-2"
        />
      </router-link>
    </div>

    <Toast />

    <div class="bg-amber-50 border border-amber-200 rounded-2xl p-4 flex items-start gap-3.5 text-amber-900 text-sm shadow-xs transition-all">
      <div class="bg-amber-100 p-2 rounded-xl shrink-0 text-amber-700">
        <i class="pi pi-exclamation-triangle text-base block"></i>
      </div>
      <div class="space-y-0.5">
        <span class="font-bold uppercase tracking-wider block text-xs text-amber-800">Peringatan Prosedur Logistik:</span>
        <p class="m-0 text-amber-700/90 leading-relaxed text-xs">
          Pastikan seluruh identitas fisik penyewa terekam utuh dan valid. Status pemesanan barang sewa ini akan terkunci sebagai 
          <span class="font-bold text-amber-900 bg-amber-200/60 px-1.5 py-0.5 rounded">BOOKED</span> 
          dan membutuhkan otorisasi verifikasi ketersediaan stok fisik gudang sebelum nota resmi dapat diterbitkan.
        </p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      
      <div class="lg:col-span-8 space-y-6">
        
        <div class="grid grid-cols-2 gap-2 bg-slate-100/80 p-1.5 rounded-2xl border border-slate-200/60">
          <button
            @click="activeStep = 1"
            :class="activeStep === 1 ? 'bg-white text-indigo-700 shadow-sm border-slate-200/50' : 'text-slate-500 hover:text-slate-800 hover:bg-white/50'"
            class="flex py-3 px-4 rounded-xl font-bold text-xs transition-all items-center justify-center gap-2.5 cursor-pointer border border-transparent bg-transparent"
          >
            <span :class="activeStep === 1 ? 'bg-indigo-600 text-white' : 'bg-slate-200 text-slate-600'" class="w-6 h-6 rounded-lg flex items-center justify-center text-[11px] font-bold transition-all">1</span>
            Data & Periode Sewa
          </button>
          
          <button
            @click="activeStep = 2"
            :class="activeStep === 2 ? 'bg-white text-indigo-700 shadow-sm border-slate-200/50' : 'text-slate-500 hover:text-slate-800 hover:bg-white/50'"
            class="flex py-3 px-4 rounded-xl font-bold text-xs transition-all items-center justify-center gap-2.5 cursor-pointer border border-transparent bg-transparent"
          >
            <span :class="activeStep === 2 ? 'bg-indigo-600 text-white' : 'bg-slate-200 text-slate-600'" class="w-6 h-6 rounded-lg flex items-center justify-center text-[11px] font-bold transition-all">2</span>
            Katalog & Keranjang
          </button>
        </div>

        <div v-show="activeStep === 1" class="space-y-6">
          <div class="bg-white border border-slate-200 rounded-2xl p-6 shadow-xs space-y-5">
            <div class="flex items-center justify-between pb-3 border-b border-slate-100">
              <div class="flex items-center gap-2.5">
                <span class="w-7 h-7 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold shadow-xs"><i class="pi pi-id-card text-sm"></i></span>
                <h3 class="font-bold text-slate-900 text-base m-0">1. Data Pemesan (Customer Snapshot)</h3>
              </div>
              <span class="text-[10px] bg-red-50 text-red-600 font-bold px-2.5 py-1 rounded-lg uppercase tracking-wider border border-red-100">Wajib</span>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5 md:col-span-2">
                <label class="text-xs font-bold text-slate-600">Nama Lengkap Penyewa <span class="text-red-500">*</span></label>
                <InputText v-model="formCustomer.customer_name" placeholder="Masukkan nama lengkap sesuai KTP" class="p-inputtext-sm rounded-xl border-slate-300 focus:border-indigo-500 font-medium text-slate-900 w-full p-2.5" />
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-bold text-slate-600">Nomor WhatsApp Aktif <span class="text-red-500">*</span></label>
                <InputText v-model="formCustomer.customer_phone" placeholder="081234567xxx" class="p-inputtext-sm rounded-xl border-slate-300 font-mono text-slate-900 w-full p-2.5" />
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-bold text-slate-600">Nomor KTP / NIK (Opsional)</label>
                <InputText v-model="formCustomer.customer_id_card" placeholder="327101234567xxxx" class="p-inputtext-sm rounded-xl border-slate-300 font-mono text-slate-900 w-full p-2.5" />
              </div>
            </div>
          </div>

          <div class="bg-white border border-slate-200 rounded-2xl p-6 shadow-xs space-y-5">
            <div class="flex items-center justify-between pb-3 border-b border-slate-100">
              <div class="flex items-center gap-2.5">
                <span class="w-7 h-7 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold shadow-xs"><i class="pi pi-calendar text-sm"></i></span>
                <h3 class="font-bold text-slate-900 text-base m-0">2. Periode Peminjaman Box</h3>
              </div>
              <Tag :value="`Durasi: ${calculatedDays} Hari`" severity="info" class="font-mono font-bold text-[11px] uppercase px-3 py-1 rounded-lg" />
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-bold text-slate-600">Tanggal Diambil (Mulai Sewa) <span class="text-red-500">*</span></label>
                <Calendar v-model="formDates.start_date" :minDate="today" showIcon class="w-full custom-calendar" dateFormat="yy-mm-dd" placeholder="Pilih tanggal ambil..." @update:modelValue="resetLock" />
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-bold text-slate-600">Tanggal Kembali (Selesai Sewa) <span class="text-red-500">*</span></label>
                <Calendar v-model="formDates.end_date" :minDate="formDates.start_date || today" showIcon class="w-full custom-calendar" dateFormat="yy-mm-dd" placeholder="Pilih tanggal kembali..." @update:modelValue="resetLock" />
              </div>
            </div>
          </div>

          <div class="flex justify-end pt-2">
            <Button label="Lanjut ke Katalog Produk" icon="pi pi-arrow-right" iconPos="right" class="bg-indigo-600 hover:bg-indigo-700 text-white border-none rounded-xl font-bold px-5 py-3 text-xs shadow-sm transition-all cursor-pointer" @click="activeStep = 2" />
          </div>
        </div>

        <div v-show="activeStep === 2" class="space-y-6">
          <div class="bg-white border border-slate-200 rounded-2xl p-6 shadow-xs space-y-5">
            <div class="flex items-center justify-between pb-3 border-b border-slate-100">
              <div class="flex items-center gap-2.5">
                <span class="w-7 h-7 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold shadow-xs"><i class="pi pi-box text-sm"></i></span>
                <h3 class="font-bold text-slate-900 text-base m-0">3. Pilihan Kelengkapan Hantaran</h3>
              </div>
              <span class="text-xs font-bold text-indigo-700 font-mono bg-indigo-50 px-3 py-1 rounded-lg border border-indigo-100">Total: {{ cart.length }} Macam Box</span>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 items-stretch sm:items-end bg-slate-50 p-4 rounded-xl border border-slate-200/60">
              <div class="flex flex-col gap-1.5 flex-1 min-w-[240px]">
                <label class="text-xs font-bold text-slate-600">Cari Produk Master</label>
                <Dropdown v-model="selectedProduct" :options="rentalStore.products" optionLabel="name" placeholder="Pilih paket box akrilik/rotan..." class="w-full rounded-xl border-slate-300 text-sm shadow-xs bg-white h-[42px]" :filter="true">
                  <template #value="slotProps">
                    <div v-if="slotProps.value" class="flex items-center gap-2 text-xs py-0.5">
                      <img v-if="slotProps.value.object_name" :src="`http://localhost:9000/foto-produk-sewa/${slotProps.value.object_name}`" class="w-6 h-6 object-cover rounded-md border shadow-xs" />
                      <i v-else class="pi pi-box text-slate-400"></i>
                      <span class="font-semibold text-slate-800">{{ slotProps.value.name }}</span>
                    </div>
                    <span v-else class="text-xs text-slate-400 self-center">{{ slotProps.placeholder }}</span>
                  </template>
                  <template #option="slotProps">
                    <div class="flex items-center gap-3 py-1 border-b border-slate-50 last:border-none">
                      <img v-if="slotProps.option.object_name" :src="`http://localhost:9000/foto-produk-sewa/${slotProps.option.object_name}`" class="w-10 h-10 object-cover rounded-lg border shadow-xs shrink-0" />
                      <div v-else class="w-10 h-10 bg-slate-100 rounded-lg flex items-center justify-center text-slate-400 border shrink-0"><i class="pi pi-box text-sm"></i></div>
                      <div class="min-w-0 flex-1">
                        <div class="font-bold text-xs text-slate-800 truncate">{{ slotProps.option.name }}</div>
                        <div class="text-[11px] text-indigo-600 font-mono font-bold mt-0.5">
                          {{ formatCurrency(slotProps.option.rental_price) }} / periode
                        </div>
                      </div>
                    </div>
                  </template>
                </Dropdown>
              </div>
              <div class="flex flex-col gap-1.5 w-full sm:w-28">
                <label class="text-xs font-bold text-slate-600">Jumlah Qty</label>
                <InputNumber v-model="selectedQty" :min="1" class="w-full rounded-xl custom-input-number h-[42px] bg-white" inputClass="w-full text-center font-bold text-slate-800 border-slate-300 rounded-xl" showButtons buttonLayout="horizontal" incrementButtonIcon="pi pi-plus" decrementButtonIcon="pi pi-minus" />
              </div>
              <Button label="Tambah" icon="pi pi-plus" class="bg-indigo-600 hover:bg-indigo-700 text-white border-none px-5 rounded-xl text-xs font-bold shadow-xs h-[42px] cursor-pointer shrink-0 transition-all" @click="addItemToCart" />
            </div>

            <div class="border border-slate-200 rounded-xl overflow-hidden shadow-xs bg-white">
              <DataTable :value="cart" class="p-datatable-sm text-sm" scrollable scrollHeight="280px">
                <template #empty>
                  <div class="p-10 text-center text-slate-400 italic font-medium bg-white flex flex-col items-center justify-center gap-2">
                    <i class="pi pi-shopping-bag text-2xl text-slate-300"></i>
                    Belum ada paket hantaran di dalam nota belanja.
                  </div>
                </template>
                <Column field="name" header="Nama Paket / Box" class="p-3">
                  <template #body="slotProps">
                    <div class="flex items-center gap-3">
                      <img v-if="slotProps.data.object_name" :src="`http://localhost:9000/foto-produk-sewa/${slotProps.data.object_name}`" class="w-10 h-10 object-cover rounded-lg border shadow-xs bg-white shrink-0" />
                      <div v-else class="w-10 h-10 bg-slate-50 rounded-lg border flex items-center justify-center text-slate-300 shrink-0"><i class="pi pi-box"></i></div>
                      <div class="min-w-0">
                        <span class="font-bold text-slate-800 block truncate">{{ slotProps.data.name }}</span>
                        <span class="text-[10px] font-mono text-slate-400 uppercase tracking-wider block mt-0.5">ID: {{ slotProps.data.rental_product_id.substring(0,8) }}</span>
                      </div>
                    </div>
                  </template>
                </Column>
                <Column header="Tarif Pokok" class="font-mono text-right pr-4 text-slate-600">
                  <template #body="slotProps">{{ formatCurrency(slotProps.data.price) }}</template>
                </Column>
                <Column field="qty" header="Qty" class="text-center font-mono font-bold w-20 text-slate-800"></Column>
                <Column header="Subtotal" class="font-bold text-indigo-600 text-right font-mono pr-4 w-36">
                  <template #body="slotProps">{{ formatCurrency(slotProps.data.price * slotProps.data.qty) }}</template>
                </Column>
                <Column header="Aksi" class="w-16 text-center">
                  <template #body="slotProps">
                    <Button icon="pi pi-trash" class="p-button-text p-button-sm text-red-600 rounded-lg hover:bg-red-50 cursor-pointer" @click="removeItem(slotProps.index)" />
                  </template>
                </Column>
              </DataTable>
            </div>

            <div class="flex justify-start pt-2">
              <Button label="Kembali ke Data Pemesan" icon="pi pi-arrow-left" class="p-button-text p-button-sm text-slate-500 hover:bg-slate-100 rounded-xl font-bold px-4 py-2" @click="activeStep = 1" />
            </div>
          </div>
        </div>
      </div>
      <div class="lg:col-span-4 sticky top-6 space-y-6">
        <div class="bg-slate-900 text-white p-6 rounded-2xl shadow-xl space-y-6 flex flex-col justify-between min-h-[420px] border border-slate-800">
          <div>
            <div class="flex items-center justify-between pb-3 border-b border-slate-800">
              <span class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">Total Rincian Kontrak</span>
              <button @click="toggleDevBypass" class="px-2 py-0.5 rounded text-[9px] font-mono font-bold transition-all cursor-pointer border" :class="devBypass ? 'bg-amber-400 text-slate-950 border-amber-300' : 'bg-slate-800 text-slate-400 border-slate-700 hover:text-white'">
                ⚡ Bypass: {{ devBypass ? 'ON' : 'OFF' }}
              </button>
            </div>

            <div class="space-y-4 pt-4 text-xs text-slate-300">
              <div class="flex justify-between items-center">
                <span class="text-slate-400">Akumulasi Subtotal:</span>
                <span class="font-bold text-white font-mono text-lg">{{ formatCurrency(calculatedSubtotal) }}</span>
              </div>

              <div class="flex flex-col gap-1.5 pt-2">
                <label class="text-[10px] font-bold uppercase tracking-wider text-amber-400">Setoran Uang Muka (DP)</label>
                <InputNumber v-model="downPayment" mode="currency" currency="IDR" locale="id-ID" class="w-full font-bold text-slate-950 custom-input-dp" placeholder="Rp 0" inputClass="rounded-xl font-bold p-2.5 bg-white" />
              </div>
            </div>
          </div>

          <div class="pt-4 border-t border-slate-800 space-y-4">
            <div class="flex justify-between items-baseline">
              <span class="text-xs text-slate-400">Sisa Pelunasan:</span>
              <span class="text-2xl font-extrabold text-amber-400 font-mono tracking-tight">{{ formatCurrency(calculatedSubtotal - downPayment) }}</span>
            </div>

            <div class="space-y-2.5 pt-2">
              <Button 
                label="1. Verifikasi Ketersediaan Stok" 
                icon="pi pi-shield" 
                class="w-full bg-indigo-600 hover:bg-indigo-700 text-white border-none font-bold text-xs py-3 rounded-xl transition-all shadow-xs cursor-pointer" 
                :loading="checkingAvailability" 
                @click="checkAllAvailability" 
              />

              <div v-if="!isSubmitUnlocked" class="text-[11px] bg-slate-800/80 border border-amber-500/20 text-amber-300/90 p-2.5 rounded-xl text-center leading-relaxed">
                🔒 Klik <strong class="font-bold text-amber-400">Verifikasi Stok</strong> terlebih dahulu untuk membuka gembok validasi aset.
              </div>

              <Button 
                label="2. Simpan & Terbitkan Nota" 
                icon="pi pi-check-circle" 
                class="w-full bg-emerald-600 hover:bg-emerald-500 text-white border-none font-bold text-xs py-3.5 rounded-xl transition-all shadow-md cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed" 
                :loading="submitting" 
                :disabled="!isSubmitUnlocked" 
                @click="submitReservation" 
              />
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
/* Mengharmonisasikan Input PrimeVue agar fit dengan gaya Tailwind moderen */
:deep(.p-calendar .p-inputtext) {
  border-radius: 12px !important;
  border-color: #cbd5e1 !important;
  padding: 0.625rem !important;
  font-size: 0.875rem !important;
}
:deep(.p-calendar .p-datepicker-trigger) {
  background: transparent !important;
  color: #64748b !important;
  border: 1px solid #cbd5e1 !important;
  border-left: none !important;
  border-top-right-radius: 12px !important;
  border-bottom-right-radius: 12px !important;
}
:deep(.custom-input-number .p-inputtext) {
  padding: 0.625rem !important;
  font-size: 0.875rem !important;
}
:deep(.custom-input-dp .p-inputtext) {
  width: 100% !important;
}
</style>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useRentalStore } from '@frontend/stores/rental'
import { useAuthStore } from '@frontend/stores/auth'
import { rentalApi } from '@frontend/api/rental'
import type { CreateReservationPayload } from '@frontend/api/rental'
import { cashApi } from '@frontend/api/cash'
import type { RentalProduct } from '@frontend/types/rental'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Calendar from 'primevue/calendar'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Toast from 'primevue/toast'
import Tag from 'primevue/tag'

import { jsPDF } from 'jspdf'
import html2canvas from 'html2canvas'

const router = useRouter()
const toast = useToast()
const rentalStore = useRentalStore()
const authStore = useAuthStore()

const checkingAvailability = ref(false)
const isAvailabilityVerified = ref(false)
const submitting = ref(false)
const devBypass = ref(false)
const currentSessionId = ref('')

const formCustomer = ref({ customer_name: '', customer_phone: '', customer_id_card: '' })
const formDates = ref({ start_date: null as Date | null, end_date: null as Date | null })
const downPayment = ref(0)
const today = ref(new Date())
const activeStep = ref(1)

interface CartItem {
  rental_product_id: string
  name: string
  price: number
  qty: number
  object_name?: string
}
const cart = ref<CartItem[]>([])

const selectedProduct = ref<RentalProduct | null>(null)
const selectedQty = ref(1)

onMounted(async () => {
  rentalStore.fetchProducts()
  
  if (authStore.role === 'CASHIER') {
    try {
      const session = await cashApi.getCurrentSession()
      if (session && session.id) {
        currentSessionId.value = session.id
      } else {
        toast.add({ severity: 'warn', summary: 'Sesi Kasir Tutup', detail: 'Anda belum membuka shift kasir. Transaksi mungkin akan ditolak.', life: 5000 })
      }
    } catch (error) {
      console.error('Gagal menarik data sesi kasir:', error)
      toast.add({ severity: 'error', summary: 'Koneksi Kasir Gagal', detail: 'Sistem tidak dapat memvalidasi sesi kasir aktif.', life: 4000 })
    }
  } else {
    // Gunakan UUID user Pegawai sebagai pelacak sesi untuk non-kasir
    currentSessionId.value = authStore.user?.id || '00000000-0000-0000-0000-000000000000'
  }
})

const calculatedDays = computed(() => {
  if (!formDates.value.start_date || !formDates.value.end_date) return 0
  const diffTime = Math.abs(formDates.value.end_date.getTime() - formDates.value.start_date.getTime())
  const diff = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  return diff > 0 ? diff : 1
})

const calculatedSubtotal = computed(() =>
  cart.value.reduce((sum, item) => sum + item.price * item.qty, 0)
)

const isSubmitUnlocked = computed(() => devBypass.value || isAvailabilityVerified.value)

const resetLock = () => {
  isAvailabilityVerified.value = false
}
const toggleDevBypass = () => {
  devBypass.value = !devBypass.value
}

const addItemToCart = () => {
  if (!selectedProduct.value) return
  const existing = cart.value.find((i: CartItem) => i.rental_product_id === selectedProduct.value!.id)
  if (existing) existing.qty += selectedQty.value
  else {
    cart.value.push({
      rental_product_id: selectedProduct.value.id,
      name: selectedProduct.value.name,
      price: selectedProduct.value.rental_price,
      qty: selectedQty.value,
      object_name: selectedProduct.value.object_name,
    })
  }
  selectedProduct.value = null
  selectedQty.value = 1
  resetLock()
}

const removeItem = (idx: number) => {
  cart.value.splice(idx, 1)
  resetLock()
}

const checkAllAvailability = async () => {
  if (!formDates.value.start_date || !formDates.value.end_date) {
    toast.add({ severity: 'warn', summary: 'Pilih Tanggal', detail: 'Tentukan tanggal mulai dan selesai pemakaian box', life: 3000 })
    return
  }
  if (cart.value.length === 0) {
    toast.add({ severity: 'warn', summary: 'Keranjang Kosong', detail: 'Masukkan minimal 1 produk hantaran ke dalam nota', life: 3000 })
    return
  }

  checkingAvailability.value = true
  isAvailabilityVerified.value = false

  try {
    const startIso = formDates.value.start_date.toISOString()
    const endIso = formDates.value.end_date.toISOString()

    for (const item of cart.value) {
      const params = { product_id: item.rental_product_id, start_date: startIso, end_date: endIso, qty: item.qty }
      const res = await rentalApi.checkAvailability(params)

      let isAvailable = false
      if ((res as any)?.is_available !== undefined) {
        isAvailable = Boolean((res as any).is_available)
      } else {
        isAvailable = true
      }

      if (!isAvailable) {
        toast.add({ severity: 'error', summary: 'Stok Habis', detail: `Paket "${item.name}" tidak mencukupi pada tanggal tersebut`, life: 4500 })
        checkingAvailability.value = false
        return
      }
    }

    isAvailabilityVerified.value = true
    toast.add({ severity: 'success', summary: 'Stok Terverifikasi!', detail: 'Seluruh paket tersedia. Silakan klik tombol Hijau di bawah', life: 3000 })
  } catch (e: any) {
    toast.add({ severity: 'error', summary: 'Gagal Menghubungi Server', detail: e.response?.data?.message || 'Sistem gagal memvalidasi ketersediaan stok', life: 3500 })
  } finally {
    checkingAvailability.value = false
  }
}

const submitReservation = async () => {
  if (!formCustomer.value.customer_name || !formCustomer.value.customer_phone) {
    toast.add({ severity: 'warn', summary: 'Data Penyewa Tidak Lengkap', detail: 'Nama lengkap dan nomor WhatsApp wajib diisi', life: 3500 })
    return
  }
  if (cart.value.length === 0) {
    toast.add({ severity: 'warn', summary: 'Keranjang Kosong', detail: 'Silakan tambahkan setidaknya satu paket ke dalam keranjang', life: 3500 })
    return
  }

  submitting.value = true
  try {
    const payload: CreateReservationPayload = {
      customer_name: formCustomer.value.customer_name,
      customer_phone: formCustomer.value.customer_phone,
      customer_identity: formCustomer.value.customer_id_card || '-',
      start_date: formDates.value.start_date!.toISOString(),
      end_date: formDates.value.end_date!.toISOString(),
      down_payment: downPayment.value,
      amount_paid: downPayment.value,
      cashier_session_id: currentSessionId.value,
      items: cart.value.map((c: CartItem) => ({
        rental_product_id: c.rental_product_id,
        qty: c.qty,
        price_per_period: c.price,
      })),
    }

    const res = await rentalStore.executeCreateContract(payload)
    const invoiceNumber = res.invoice_number
    toast.add({ severity: 'success', summary: 'Nota Sah Terbit!', detail: `Dokumen berhasil dicetak. No. Invoice: ${invoiceNumber || 'TEREKAM'}`, life: 3000 })

    const reservationId = res.id
    if (reservationId) {
      toast.add({ severity: 'info', summary: 'Mencetak Invoice', detail: 'Sedang mengunggah PDF reservasi...', life: 2000 })

      try {
        const rawCustomerName = formCustomer.value.customer_name || 'Guest'
        const safeCustomerName = rawCustomerName.replace(/[^a-zA-Z0-9]/g, '_')
        const dateStr = formDates.value.start_date ? formDates.value.start_date.toISOString().split('T')[0] : 'nodate'
        const invoiceStr = invoiceNumber || reservationId
        const safeInvoiceNo = invoiceStr.replace(/[\/\\]/g, '_')
        const filename = `${safeCustomerName}_${dateStr}_${safeInvoiceNo}.pdf`

        // NATIVE PREMIUM jsPDF Generator
        const doc = new jsPDF('p', 'mm', 'a4')
        const pageWidth = doc.internal.pageSize.getWidth()

        // 1. HEADER BANNER (Modern Minimalist Navy Accent)
        doc.setFillColor(37, 69, 130) // Brand Color #254582
        doc.rect(0, 0, pageWidth, 42, 'F')
        
        doc.setTextColor(255, 255, 255)
        doc.setFont('helvetica', 'bold')
        doc.setFontSize(22)
        doc.text('INVOICE SEWA BOX', 15, 18)
        
        doc.setFont('helvetica', 'normal')
        doc.setFontSize(10)
        doc.setTextColor(202, 216, 255)
        doc.text('Nota Tanda Terima Resmi & Validasi Kontrak Digital', 15, 25)
        
        doc.setFont('font-mono', 'bold')
        doc.setFontSize(11)
        doc.setTextColor(255, 255, 255)
        doc.text(`NO. NOTA: ${invoiceStr}`, pageWidth - 15, 18, { align: 'right' })
        
        doc.setFont('helvetica', 'normal')
        doc.setFontSize(9)
        doc.setTextColor(202, 216, 255)
        doc.text(`Dibuat: ${new Date().toLocaleDateString('id-ID')}`, pageWidth - 15, 25, { align: 'right' })

        // 2. METADATA GRID (Customer & Date Information Block)
        let infoY = 54
        doc.setDrawColor(226, 232, 240)
        doc.setFillColor(248, 250, 252)
        doc.rect(15, infoY, pageWidth - 30, 22, 'FD') // Light Background box

        // Column 1: Customer Details
        doc.setTextColor(100, 116, 139)
        doc.setFontSize(8)
        doc.setFont('helvetica', 'bold')
        doc.text('NAMA PENYEWA', 20, infoY + 6)
        doc.setTextColor(30, 41, 59)
        doc.setFontSize(10)
        doc.setFont('helvetica', 'bold')
        doc.text(rawCustomerName, 20, infoY + 14)

        doc.setTextColor(100, 116, 139)
        doc.setFontSize(8)
        doc.setFont('helvetica', 'bold')
        doc.text('KONTAK WHATSAPP', 85, infoY + 6)
        doc.setTextColor(22, 163, 74) // Green WA
        doc.setFontSize(10)
        doc.setFont('helvetica', 'bold')
        doc.text(formCustomer.value.customer_phone || '-', 85, infoY + 14)

        // Column 2: Rental Dates
        const sDate = formDates.value.start_date ? formDates.value.start_date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '-'
        const eDate = formDates.value.end_date ? formDates.value.end_date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '-'

        doc.setTextColor(100, 116, 139)
        doc.setFontSize(8)
        doc.setFont('helvetica', 'bold')
        doc.text('TANGGAL ACARA', 145, infoY + 6)
        doc.setTextColor(30, 41, 59)
        doc.setFontSize(9)
        doc.setFont('helvetica', 'normal')
        doc.text(sDate, 145, infoY + 14)

        doc.setTextColor(100, 116, 139)
        doc.setFontSize(8)
        doc.setFont('helvetica', 'bold')
        doc.text('TENGGAT KEMBALI', 145, infoY + 11)
        doc.setTextColor(220, 38, 38) // Red Warning Accent
        doc.setFontSize(9)
        doc.setFont('helvetica', 'bold')
        doc.text(eDate, 145, infoY + 18)

        // 3. TABLE PRODUCTS (Rincian Paket Terpesan)
        let currentY = 90
        doc.setTextColor(37, 69, 130)
        doc.setFontSize(10)
        doc.setFont('helvetica', 'bold')
        doc.text('RINCIAN PAKET TERPESAN DI DALAM NOTA', 15, currentY)
        
        currentY += 4
        // Table Header Fill
        doc.setFillColor(37, 69, 130)
        doc.rect(15, currentY, pageWidth - 30, 8, 'F')
        
        doc.setTextColor(255, 255, 255)
        doc.setFontSize(9)
        doc.setFont('helvetica', 'bold')
        doc.text('Nama Barang / Box', 19, currentY + 5.5)
        doc.text('Qty', 110, currentY + 5.5, { align: 'center' })
        doc.text('Harga Satuan', 145, currentY + 5.5, { align: 'right' })
        doc.text('Subtotal', pageWidth - 19, currentY + 5.5, { align: 'right' })

        // Table Rows Loop
        currentY += 8
        doc.setFont('helvetica', 'normal')
        doc.setFontSize(9)
        
        cart.value.forEach((item, index) => {
          // Zebra striping effect
          if (index % 2 === 0) {
            doc.setFillColor(250, 250, 250)
            doc.rect(15, currentY, pageWidth - 30, 7.5, 'F')
          }
          
          // Row bottom line divider
          doc.setDrawColor(241, 245, 249)
          doc.line(15, currentY + 7.5, pageWidth - 15, currentY + 7.5)

          doc.setTextColor(51, 65, 85)
          doc.text(item.name, 19, currentY + 5)
          
          doc.setTextColor(30, 41, 59)
          doc.text(item.qty.toString(), 110, currentY + 5, { align: 'center' })
          doc.text(`Rp ${item.price.toLocaleString('id-ID')}`, 145, currentY + 5, { align: 'right' })
          
          const subtot = item.qty * item.price
          doc.setFont('helvetica', 'bold')
          doc.text(`Rp ${subtot.toLocaleString('id-ID')}`, pageWidth - 19, currentY + 5, { align: 'right' })
          doc.setFont('helvetica', 'normal')
          
          currentY += 7.5
        })

        // 4. FINANCIAL SUMMARY BLOCK
        currentY += 12
        const summaryWidth = 100
        const startX = pageWidth - 15 - summaryWidth
        
        // Background card for totals
        doc.setFillColor(248, 250, 252)
        doc.setDrawColor(226, 232, 240)
        doc.rect(startX, currentY, summaryWidth, 38, 'FD')

        doc.setFontSize(9)
        let rowY = currentY + 6

        // Subtotal Pokok
        doc.setTextColor(100, 116, 139)
        doc.text('Total Pokok Sewa Box:', startX + 5, rowY)
        doc.setTextColor(30, 41, 59)
        doc.setFont('helvetica', 'bold')
        doc.text(`Rp ${calculatedSubtotal.value.toLocaleString('id-ID')}`, pageWidth - 20, rowY, { align: 'right' })

        // Total Tagihan
        rowY += 7
        doc.setFont('helvetica', 'normal')
        doc.setTextColor(100, 116, 139)
        doc.text('Total Tagihan:', startX + 5, rowY)
        doc.setTextColor(30, 41, 59)
        doc.setFont('helvetica', 'bold')
        doc.text(`Rp ${calculatedSubtotal.value.toLocaleString('id-ID')}`, pageWidth - 20, rowY, { align: 'right' })

        // Down Payment (DP)
        rowY += 7
        doc.setFont('helvetica', 'normal')
        doc.setTextColor(22, 163, 74) // Clean green for paid values
        doc.text('Uang Muka Terbayar (DP):', startX + 5, rowY)
        doc.setFont('helvetica', 'bold')
        doc.text(`-Rp ${downPayment.value.toLocaleString('id-ID')}`, pageWidth - 20, rowY, { align: 'right' })

        // Divider Inside Box
        rowY += 3
        doc.setDrawColor(203, 213, 225)
        doc.line(startX + 4, rowY, pageWidth - 19, rowY)

        // Final Due Bill
        rowY += 6
        doc.setFont('helvetica', 'bold')
        doc.setTextColor(37, 69, 130) // Deep Navy Highlight
        doc.text('Tagihan Akhir:', startX + 5, rowY)
        const finalBill = calculatedSubtotal.value - downPayment.value
        doc.setFontSize(10)
        doc.text(`Rp ${finalBill.toLocaleString('id-ID')}`, pageWidth - 20, rowY, { align: 'right' })

        // 5. LEGAL SIGNATURE FOOTER
        doc.setDrawColor(226, 232, 240)
        doc.line(15, 255, pageWidth - 15, 255)
        
        doc.setFontSize(8)
        doc.setFont('helvetica', 'normal')
        doc.setTextColor(148, 163, 184)
        doc.text('Syarat & Ketentuan berlaku. Harap kembalikan box tepat waktu sesuai dengan tenggat tanggal di atas untuk menghindari denda.', 15, 261, { maxWidth: pageWidth - 30 })

        doc.setFontSize(9)
        doc.setTextColor(100, 116, 139)
        doc.text('Terima kasih atas kepercayaan Anda pada Bisnis Rinzi.', pageWidth / 2, 276, { align: 'center' })

        const pdfBlob = doc.output('blob')
        doc.save(filename)

        await rentalStore.uploadReservationInvoice(reservationId, pdfBlob, filename)

      } catch (pdfErr: any) {
        console.error('PDF Generation/Upload failed: ', pdfErr)
        toast.add({ severity: 'warn', summary: 'Nota Gagal Dibuat', detail: 'Reservasi berhasil disimpan, namun gagal mengunggah PDF nota.', life: 4000 })
      }
    }

    setTimeout(() => {
      router.push('/reservations-list')
    }, 1200)
  } catch (e: any) {
    console.error('Contract/Upload Error: ', e)
    toast.add({ severity: 'error', summary: 'Gagal Menyimpan Kontrak', detail: e.response?.data?.message || 'Sistem menolak pencatatan dokumen baru', life: 4000 })
  } finally {
    submitting.value = false
  }
}

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
</script>