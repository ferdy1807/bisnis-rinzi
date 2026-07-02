<template>
  <div class="space-y-8 pb-12 font-sans text-[#1a1c20]">
    <div
      class="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-[#c4c6d1]/60 pb-6"
    >
      <div>
        <p class="text-xs font-extrabold tracking-widest text-[#254582] uppercase mb-1">
          Pemesanan & Invoice
        </p>
        <h2 class="text-3xl font-black text-[#1a1c20] tracking-tight">Master Arsip Kontrak Sewa</h2>
      </div>
      <div class="flex items-center gap-3">
        <span class="text-xs text-[#58627c] hidden sm:inline"
          >Pembaruan Terakhir: <strong class="font-mono">{{ lastSynced }}</strong></span
        >
        <Button
          label="Segarkan Arsip"
          icon="pi pi-refresh"
          class="p-button-sm bg-[#254582] hover:bg-[#3f5d9b] text-white border-none rounded-xl font-bold px-4 py-2.5 shadow-xs"
          :loading="rentalStore.isLoading.master"
          @click="syncData"
        />
      </div>
    </div>

    <Toast />

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div
        @click="setStatusFilter(null)"
        class="p-5 rounded-3xl border transition-all duration-300 cursor-pointer flex flex-col justify-between relative overflow-hidden group"
        :class="
          selectedStatusFilter === null
            ? 'bg-[#254582] text-white shadow-xl scale-[1.02] border-[#254582]'
            : 'bg-white text-gray-800 border-gray-200 hover:border-[#254582]/50 shadow-2xs'
        "
      >
        <div class="flex items-center justify-between mb-2">
          <span
            class="p-2.5 rounded-2xl"
            :class="
              selectedStatusFilter === null
                ? 'bg-white/10 text-white'
                : 'bg-[#f3f3f9] text-[#254582]'
            "
            ><i class="pi pi-folder-open text-xl"></i
          ></span>
          <span
            class="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md"
            :class="
              selectedStatusFilter === null ? 'bg-white/20 text-white' : 'bg-gray-100 text-gray-500'
            "
            >Filter Aktif</span
          >
        </div>
        <div>
          <p class="text-xs opacity-80 mb-0.5 font-medium">Total Seluruh Arsip</p>
          <h3 class="text-2xl font-black font-mono">
            {{ countAll }} <span class="text-xs font-normal opacity-75">Nota</span>
          </h3>
        </div>
      </div>

      <div
        @click="setStatusFilter('BOOKED')"
        class="p-5 rounded-3xl border transition-all duration-300 cursor-pointer flex flex-col justify-between relative overflow-hidden group"
        :class="
          selectedStatusFilter === 'BOOKED'
            ? 'bg-[#d9e2ff] text-[#001944] shadow-xl scale-[1.02] border-[#b0c2ff]'
            : 'bg-white text-gray-800 border-gray-200 hover:border-blue-300 shadow-2xs'
        "
      >
        <div class="flex items-center justify-between mb-2">
          <span class="p-2.5 rounded-2xl bg-blue-50 text-blue-600"
            ><i class="pi pi-bookmark text-xl"></i
          ></span>
          <span
            class="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-blue-100 text-blue-800"
            >Menunggu</span
          >
        </div>
        <div>
          <p class="text-xs opacity-80 mb-0.5 font-medium">Siap Diambil / Dihias</p>
          <h3 class="text-2xl font-black font-mono text-blue-700">
            {{ countBooked }} <span class="text-xs font-normal opacity-75">Nota</span>
          </h3>
        </div>
      </div>

      <div
        @click="setStatusFilter('PICKED_UP')"
        class="p-5 rounded-3xl border transition-all duration-300 cursor-pointer flex flex-col justify-between relative overflow-hidden group"
        :class="
          selectedStatusFilter === 'PICKED_UP'
            ? 'bg-[#ffdeac] text-[#5f4100] shadow-xl scale-[1.02] border-[#ffd18c]'
            : 'bg-white text-gray-800 border-gray-200 hover:border-amber-300 shadow-2xs'
        "
      >
        <div class="flex items-center justify-between mb-2">
          <span class="p-2.5 rounded-2xl bg-amber-50 text-amber-600"
            ><i class="pi pi-sign-out text-xl"></i
          ></span>
          <span
            class="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-amber-200 text-amber-900"
            >Keluar</span
          >
        </div>
        <div>
          <p class="text-xs opacity-80 mb-0.5 font-medium">Sedang Disewa Customer</p>
          <h3 class="text-2xl font-black font-mono text-amber-700">
            {{ countPickedUp }} <span class="text-xs font-normal opacity-75">Nota</span>
          </h3>
        </div>
      </div>

      <div
        @click="setStatusFilter('RETURNED')"
        class="p-5 rounded-3xl border transition-all duration-300 cursor-pointer flex flex-col justify-between relative overflow-hidden group"
        :class="
          selectedStatusFilter === 'RETURNED'
            ? 'bg-emerald-100 text-emerald-900 shadow-xl scale-[1.02] border-emerald-300'
            : 'bg-white text-gray-800 border-gray-200 hover:border-emerald-300 shadow-2xs'
        "
      >
        <div class="flex items-center justify-between mb-2">
          <span class="p-2.5 rounded-2xl bg-emerald-50 text-emerald-600"
            ><i class="pi pi-check-circle text-xl"></i
          ></span>
          <span
            class="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-emerald-200 text-emerald-900"
            >Selesai</span
          >
        </div>
        <div>
          <p class="text-xs opacity-80 mb-0.5 font-medium">Box Kembali & Lunas</p>
          <h3 class="text-2xl font-black font-mono text-emerald-700">
            {{ countReturned }} <span class="text-xs font-normal opacity-75">Nota</span>
          </h3>
        </div>
      </div>
    </div>

    <div
      class="bg-white border border-[#c4c6d1]/60 rounded-3xl p-6 shadow-xs flex flex-col md:flex-row items-center justify-between gap-4"
    >
      <div class="flex items-center gap-3 w-full md:w-auto">
        <span class="p-input-icon-left w-full md:w-80">
          <i class="pi pi-search text-[#58627c]"></i>
          <InputText
            v-model="searchQuery"
            placeholder="Ketik No. Invoice / Nama Penyewa..."
            class="p-inputtext-sm w-full rounded-xl border-[#c4c6d1] font-medium pl-10"
          />
        </span>
        <Button
          v-if="searchQuery || selectedStatusFilter"
          icon="pi pi-filter-slash"
          class="p-button-outlined p-button-sm text-red-500 border-red-200 hover:bg-red-50 rounded-xl"
          @click="resetFilters"
          title="Hapus Filter"
        />
      </div>

      <div class="text-xs font-bold text-[#58627c] flex items-center gap-2 self-end md:self-center">
        <span
          >Menampilkan:
          <strong class="text-[#254582] font-mono">{{ filteredReservations.length }}</strong>
          Rekor</span
        >
      </div>
    </div>

    <div
      class="bg-white border border-[#c4c6d1]/60 rounded-3xl overflow-hidden shadow-xs flex flex-col"
    >
      <DataTable
        :value="filteredReservations"
        :loading="rentalStore.isLoading.master"
        paginator
        :rows="10"
        :rowsPerPageOptions="[5, 10, 25]"
        class="p-datatable-sm text-xs border-none custom-table-rows"
        responsiveLayout="scroll"
        rowHover
      >
        <template #empty>
          <div class="text-center py-16 text-[#58627c] space-y-3 bg-[#f3f3f9]/30">
            <i class="pi pi-inbox text-4xl text-[#c4c6d1] block"></i>
            <p class="text-sm font-bold">Tidak ada dokumen kontrak yang memenuhi kriteria filter</p>
            <Button
              label="Tampilkan Semua Dokumen"
              class="p-button-text p-button-sm text-[#254582] font-bold"
              @click="resetFilters"
            />
          </div>
        </template>

        <Column field="invoice_number" header="No. Nota" :sortable="true" class="w-32">
          <template #body="{ data }">
            <span
              class="font-mono text-xs font-black text-[#254582] bg-[#d9e2ff]/50 px-2.5 py-1 rounded-lg border border-[#d9e2ff]"
              >{{ data.invoice_number }}</span
            >
          </template>
        </Column>

        <Column header="Penyewa / Calon Pengantin" class="min-w-[200px]">
          <template #body="{ data }">
            <div>
              <span class="font-bold text-[#1a1c20] text-sm block leading-tight mb-0.5">{{
                parseCustomer(data.customer_snapshot_id).name
              }}</span>
              <a
                v-if="parseCustomer(data.customer_snapshot_id).phone !== '-'"
                :href="`https://wa.me/62${parseCustomer(data.customer_snapshot_id).phone.substring(1)}`"
                target="_blank"
                class="text-[11px] font-mono text-green-700 hover:underline flex items-center gap-1 w-fit bg-green-50 px-2 py-0.5 rounded border border-green-200"
                @click.stop
              >
                <i class="pi pi-whatsapp text-[11px]"></i>
                {{ parseCustomer(data.customer_snapshot_id).phone }}
              </a>
            </div>
          </template>
        </Column>

        <Column header="Masa Pemakaian" class="min-w-[160px]">
          <template #body="{ data }">
            <div
              class="flex items-center gap-1.5 bg-[#f3f3f9] px-2.5 py-1 rounded-xl w-fit border border-gray-200/60 text-[11px] font-medium text-[#58627c]"
            >
              <strong class="text-[#254582]">{{ formatDateShort(data.start_date) }}</strong>
              <i class="pi pi-arrow-right text-[9px] text-gray-400"></i>
              <strong class="text-[#1a1c20]">{{ formatDateShort(data.end_date) }}</strong>
            </div>
          </template>
        </Column>

        <Column field="total_amount" header="Total Nilai" :sortable="true" class="text-right">
          <template #body="{ data }">
            <span class="font-mono font-black text-sm text-[#1a1c20]">{{
              formatCurrency(data.total_amount)
            }}</span>
          </template>
        </Column>

        <Column header="Uang Muka (DP)" class="text-right">
          <template #body="{ data }">
            <span
              class="font-mono font-bold text-xs text-green-700 bg-green-50 px-2 py-0.5 rounded"
              >{{ formatCurrency(data.down_payment) }}</span
            >
          </template>
        </Column>

        <Column field="status" header="Status Nota" :sortable="true" class="text-center w-32">
          <template #body="{ data }">
            <span :class="getCustomBadge(data.status)">
              <i :class="getStatusIcon(data.status)" class="text-[10px]"></i> {{ data.status }}
            </span>
          </template>
        </Column>

        <Column header="Rincian" class="text-center w-20">
          <template #body="{ data }">
            <Button
              icon="pi pi-eye"
              class="p-button-rounded p-button-text bg-[#f3f3f9] hover:bg-[#254582] text-[#254582] hover:text-white transition-colors shadow-2xs w-8 h-8"
              @click="openDigitalReceipt(data.id)"
              title="Buka Lembar Nota"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog
      v-model:visible="displayDetailDialog"
      modal
      :showHeader="false"
      :style="{ width: '48rem' }"
      class="rounded-3xl overflow-hidden shadow-2xl p-0 font-sans"
    >
      <div v-if="rentalStore.isLoading.detail" class="p-16 text-center space-y-3 bg-white">
        <ProgressSpinner strokeWidth="4" class="w-12 h-12" />
        <p class="text-xs font-bold text-gray-500 animate-pulse">
          Menarik arsip rincian nota dari pangkalan data...
        </p>
      </div>

      <div v-else-if="rentalStore.currentDetail" class="bg-white text-[#1a1c20]">
        <div
          class="bg-[#254582] text-white p-6 flex items-center justify-between relative overflow-hidden"
        >
          <div
            class="absolute -right-10 -bottom-10 w-32 h-32 bg-white/5 rounded-full blur-xl pointer-events-none"
          ></div>
          <div class="relative z-10">
            <span class="text-[10px] font-black tracking-widest text-[#cad8ff] uppercase block mb-1"
              >Nota Tanda Terima Resmi</span
            >
            <h3 class="text-2xl font-black font-mono tracking-tight">
              {{ rentalStore.currentDetail.invoice_number }}
            </h3>
          </div>
          <div class="flex items-center gap-3 relative z-10">
            <span
              :class="getCustomBadge(rentalStore.currentDetail.status)"
              class="bg-white/20 text-white border-white/30 text-xs px-3 py-1 shadow-xs"
            >
              {{ rentalStore.currentDetail.status }}
            </span>
            <button
              @click="openReservationInvoicePdf"
              class="w-8 h-8 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center transition-colors cursor-pointer"
              title="Lihat Nota Reservasi"
            >
              <i class="pi pi-file-pdf text-xs"></i>
            </button>
            <button
              @click="displayDetailDialog = false"
              class="w-8 h-8 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center transition-colors cursor-pointer"
            >
              <i class="pi pi-times text-xs"></i>
            </button>
          </div>
        </div>

        <div class="p-6 space-y-6">
          <div
            class="bg-[#f3f3f9] p-4 rounded-2xl border border-gray-200/80 grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs"
          >
            <div>
              <span class="text-[#58627c] block text-[10px] uppercase font-bold mb-0.5"
                >Nama Penyewa:</span
              >
              <strong class="text-[#1a1c20] text-sm block truncate">{{
                parseCustomer(rentalStore.currentDetail.customer_snapshot_id).name
              }}</strong>
            </div>
            <div>
              <span class="text-[#58627c] block text-[10px] uppercase font-bold mb-0.5"
                >Kontak WhatsApp:</span
              >
              <strong class="font-mono text-green-700 block">{{
                parseCustomer(rentalStore.currentDetail.customer_snapshot_id).phone
              }}</strong>
            </div>
            <div>
              <span class="text-[#58627c] block text-[10px] uppercase font-bold mb-0.5"
                >Tanggal Acara:</span
              >
              <strong class="text-[#254582] block">{{
                formatDateShort(rentalStore.currentDetail.start_date)
              }}</strong>
            </div>
            <div>
              <span class="text-[#58627c] block text-[10px] uppercase font-bold mb-0.5"
                >Tenggat Kembali:</span
              >
              <strong class="text-red-600 block">{{
                formatDateShort(rentalStore.currentDetail.end_date)
              }}</strong>
            </div>
          </div>

          <div class="space-y-2">
            <label
              class="text-xs font-black text-[#254582] uppercase tracking-wider flex items-center gap-1.5"
            >
              <i class="pi pi-box text-xs"></i> Rincian Paket Terpesan di Dalam Nota
            </label>

            <div class="border border-gray-200 rounded-2xl overflow-hidden shadow-2xs">
              <DataTable
                :value="rentalStore.currentDetail.items || []"
                class="p-datatable-sm text-xs border-none"
              >
                <template #empty
                  ><div class="p-4 text-center text-gray-400 italic">
                    Rincian produk teragregasi pada nota induk
                  </div></template
                >
                <Column field="rental_product_name" header="Nama Barang / Box"></Column>
                <Column field="qty" header="Qty" class="text-center w-16"></Column>
                <Column header="Harga Satuan" class="text-right">
                  <template #body="{ data }">{{ formatCurrency(data.price_per_period) }}</template>
                </Column>
                <Column header="Subtotal" class="text-right font-bold text-[#254582]">
                  <template #body="{ data }">{{ formatCurrency(data.subtotal) }}</template>
                </Column>
              </DataTable>
            </div>
          </div>

          <div class="bg-white border-2 border-gray-200/80 rounded-2xl p-4 space-y-2 text-xs">
            <!-- Rincian Awal & Denda -->
            <div class="flex justify-between text-gray-600">
              <span>Total Pokok Sewa Box:</span>
              <strong class="font-mono text-gray-900">{{
                formatCurrency(rentalStore.currentDetail.total_amount)
              }}</strong>
            </div>
            <template v-if="currentReturnDetail">
              <div class="flex justify-between text-gray-600">
                <span>Denda Keterlambatan:</span>
                <strong class="font-mono text-red-600"
                  >+{{ formatCurrency(currentReturnDetail.total_late_fees) }}</strong
                >
              </div>
              <div class="flex justify-between text-gray-600">
                <span>Denda Kerusakan:</span>
                <strong class="font-mono text-orange-600"
                  >+{{ formatCurrency(currentReturnDetail.total_damage_fees) }}</strong
                >
              </div>
            </template>

            <hr class="border-gray-200 my-1" />

            <!-- Kalkulasi Tagihan -->
            <div class="flex justify-between text-gray-900 font-bold">
              <span>Total Tagihan:</span>
              <span class="font-mono">{{
                formatCurrency(
                  rentalStore.currentDetail.total_amount +
                    (currentReturnDetail?.total_late_fees || 0) +
                    (currentReturnDetail?.total_damage_fees || 0),
                )
              }}</span>
            </div>
            <div class="flex justify-between text-green-700 font-semibold">
              <span>Uang Muka Terbayar (DP):</span>
              <span class="font-mono"
                >-{{ formatCurrency(rentalStore.currentDetail.down_payment) }}</span
              >
            </div>

            <hr class="border-gray-200 my-1" />

            <!-- Tagihan Akhir yang wajib dibayar -->
            <div class="flex justify-between text-gray-900 font-black text-[13px]">
              <span>Tagihan Akhir:</span>
              <span class="font-mono">{{
                formatCurrency(
                  rentalStore.currentDetail.total_amount +
                    (currentReturnDetail?.total_late_fees || 0) +
                    (currentReturnDetail?.total_damage_fees || 0) -
                    rentalStore.currentDetail.down_payment,
                )
              }}</span>
            </div>

            <!-- Pelunasan -->
            <template
              v-if="rentalStore.currentDetail.amount_paid > rentalStore.currentDetail.down_payment"
            >
              <div
                class="flex justify-between text-blue-700 font-bold mt-1.5 pt-1.5 border-t border-gray-100"
              >
                <span>Dibayarkan (Pelunasan Kasir):</span>
                <span class="font-mono">{{formatCurrency(rentalStore.currentDetail.amount_paid - rentalStore.currentDetail.down_payment)}}</span>
              </div>
              <div class="flex justify-between text-emerald-600 font-bold">
                <span>Kembalian:</span>
                <span class="font-mono">{{
                  formatCurrency(rentalStore.currentDetail.change_amount || 0)
                }}</span>
              </div>
            </template>

            <!-- Status Lunas -->
            <div
              class="mt-4 pt-2 border-t border-gray-200 flex justify-between items-center font-black text-base"
              v-if="
                rentalStore.currentDetail.status === 'RETURNED' ||
                rentalStore.currentDetail.amount_paid >= rentalStore.currentDetail.total_amount
              "
            >
              <span class="text-emerald-700">Status Pembayaran:</span>
              <span
                class="font-mono text-xl text-emerald-700 bg-emerald-50 px-3 py-1 rounded-xl border border-emerald-200"
              >
                LUNAS
              </span>
            </div>
          </div>

          <div class="pt-2">
            <div
              v-if="rentalStore.currentDetail.status === 'BOOKED'"
              class="p-4 bg-[#ffdad6]/30 border border-[#ba1a1a]/30 rounded-2xl flex flex-col sm:flex-row gap-3 items-end"
            >
              <div class="flex-1 w-full flex flex-col gap-1">
                <label class="text-[10px] font-black text-[#93000a] uppercase tracking-wider"
                  >Potongan Denda Batal (Penalty Fee)</label
                >
                <InputNumber
                  v-model="cancellationFee"
                  mode="currency"
                  currency="IDR"
                  locale="id-ID"
                  placeholder="Rp 0"
                  class="w-full font-bold text-red-700"
                />
              </div>
              <Button
                label="Batalkan Kontrak Ini"
                icon="pi pi-user-minus"
                class="bg-[#ba1a1a] hover:bg-[#93000a] text-white border-none font-bold text-xs px-5 py-2.5 rounded-xl shadow-md h-[40px] w-full sm:w-auto"
                :loading="rentalStore.isLoading.mutate"
                @click="executeAnulir"
              />
            </div>

            <div
              v-else-if="rentalStore.currentDetail.status === 'PICKED_UP'"
              class="p-4 bg-[#ffdeac]/30 border border-amber-300/50 rounded-2xl flex items-center justify-between"
            >
              <div class="flex items-center gap-3">
                <span class="p-2 bg-amber-500 text-white rounded-xl"
                  ><i class="pi pi-compass text-lg"></i
                ></span>
                <div>
                  <h4 class="font-bold text-xs text-[#7e5700]">
                    Box Sedang Berada di Tangan Penyewa
                  </h4>
                  <p class="text-[11px] text-gray-600">
                    Klik tombol di samping untuk memproses pengembalian & denda telat.
                  </p>
                </div>
              </div>
              <Button
                label="Proses Retur & Pelunasan"
                icon="pi pi-arrow-right"
                iconPos="right"
                class="bg-[#7e5700] hover:bg-[#5f4100] text-white border-none text-xs font-bold py-2.5 px-4 rounded-xl shadow-md"
                @click="jumpToReturnForm"
              />
            </div>

            <div
              v-else-if="rentalStore.currentDetail.status === 'RETURNED'"
              class="p-4 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-2xl flex flex-col sm:flex-row items-center justify-between gap-4 font-black text-xs tracking-wide uppercase"
            >
              <div class="flex items-center gap-2">
                <i class="pi pi-verified text-emerald-600 text-base"></i> Transaksi Selesai •
                Seluruh Aset Telah Kembali
              </div>
              <Button
                v-if="currentReturnDetail?.receipt_url"
                label="Preview PDF Invoice"
                icon="pi pi-file-pdf"
                class="bg-emerald-600 hover:bg-emerald-700 text-white border-none font-bold text-xs py-2 px-4 rounded-xl shadow-md"
                @click="openInvoicePdf"
              />
            </div>
          </div>
        </div>

        <div
          class="bg-gray-50 p-4 border-t border-gray-200/80 text-center text-[11px] font-mono text-gray-400"
        >
          ID Dokumen: {{ rentalStore.currentDetail.id }}
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useRentalStore } from '@frontend/stores/rental'
import type { RentalReservation, RentalReturn } from '@frontend/types/rental'

// 1. IMPOR KOMPONEN PRIMEVUE
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import ProgressSpinner from 'primevue/progressspinner'
import Toast from 'primevue/toast'

const router = useRouter()
const toast = useToast()
const rentalStore = useRentalStore()

// State Filter & Tampilan
const searchQuery = ref('')
const selectedStatusFilter = ref<string | null>(null)
const displayDetailDialog = ref(false)
const cancellationFee = ref(0)
const currentReturnDetail = ref<RentalReturn | null>(null)

// Catatan Waktu Sinkronisasi
const lastSynced = ref(
  new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' }),
)

const syncData = async () => {
  await rentalStore.fetchMasterReservations(true)
  lastSynced.value = new Date().toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

onMounted(() => {
  rentalStore.fetchMasterReservations()
})

// ============================================================================
// KOMPUTASI METRIK BENTO ATAS
// ============================================================================
const countAll = computed(() => rentalStore.masterReservations.length)
const countBooked = computed(
  () =>
    rentalStore.masterReservations.filter((r: RentalReservation) => r.status === 'BOOKED' || r.status === 'READY_FOR_PICKUP').length,
)
const countPickedUp = computed(
  () =>
    rentalStore.masterReservations.filter((r: RentalReservation) => r.status === 'PICKED_UP')
      .length,
)
const countReturned = computed(
  () =>
    rentalStore.masterReservations.filter((r: RentalReservation) => r.status === 'RETURNED').length,
)

// ============================================================================
// FILTER GANDA (Pencarian Teks + Sakelar Kartu Bento)
// ============================================================================
const setStatusFilter = (status: string | null) => {
  selectedStatusFilter.value = selectedStatusFilter.value === status ? null : status
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedStatusFilter.value = null
}

const filteredReservations = computed(() => {
  let list = rentalStore.masterReservations

  if (selectedStatusFilter.value) {
    if (selectedStatusFilter.value === 'BOOKED') {
      list = list.filter((r: RentalReservation) => r.status === 'BOOKED' || r.status === 'READY_FOR_PICKUP')
    } else {
      list = list.filter((r: RentalReservation) => r.status === selectedStatusFilter.value)
    }
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter(
      (r: RentalReservation) =>
        r.invoice_number.toLowerCase().includes(q) ||
        r.customer_snapshot_id.toLowerCase().includes(q),
    )
  }

  return list
})

// ============================================================================
// HELPER CERDAS PEMECAH TEKS PENYEWA (RegEx Parser)
// ============================================================================
const parseCustomer = (snapshotStr?: string) => {
  if (!snapshotStr) return { name: 'Penyewa Tanpa Nama', phone: '-' }
  // Memecah teks "Natasha Wilona (08123456xxx)"
  const match = snapshotStr.match(/^(.*?)\s*\((.*?)\)$/)
  if (match) {
    return { name: match[1]?.trim() || 'Penyewa Tanpa Nama', phone: match[2]?.trim() || '-' }
  }
  return { name: snapshotStr, phone: '-' }
}

// ============================================================================
// AKSI MODAL & NAVIGASI
// ============================================================================
const openDigitalReceipt = async (id: string) => {
  displayDetailDialog.value = true
  cancellationFee.value = 0
  currentReturnDetail.value = null
  try {
    await rentalStore.loadReservationDetail(id)
    if (rentalStore.currentDetail?.status === 'RETURNED') {
      currentReturnDetail.value = await rentalStore.getReturnByReservationId(id)
    }
  } catch (e) {
    displayDetailDialog.value = false
    toast.add({ severity: 'error', summary: 'Ralat', detail: 'Nota gagal diekstraksi dari server' })
  }
}

const openInvoicePdf = () => {
  if (currentReturnDetail.value?.receipt_url) {
    let url = currentReturnDetail.value.receipt_url
    if (!url.startsWith('http')) {
      if (url.startsWith('/')) {
        url = 'http://localhost:9000' + url
      } else {
        url = 'http://localhost:9000/invoice-return/' + url.split('/').pop()
      }
    }
    window.open(url, '_blank')
  } else {
    toast.add({
      severity: 'warn',
      summary: 'File Tidak Ditemukan',
      detail: 'Invoice PDF pengembalian belum diunggah untuk transaksi ini.',
      life: 3000,
    })
  }
}

const openReservationInvoicePdf = () => {
  if (rentalStore.currentDetail?.invoice_number) {
    const rawCustomerName = parseCustomer(rentalStore.currentDetail.customer_snapshot_id).name || 'Guest'
    const safeCustomerName = rawCustomerName.replace(/[^a-zA-Z0-9]/g, '_')
    const safeInvoiceNo = rentalStore.currentDetail.invoice_number.replace(/[\/\\]/g, '_')

    let url = ''

    if (rentalStore.currentDetail.status === 'RETURNED') {
      const filename = `${safeCustomerName}_${safeInvoiceNo}.pdf`
      url = `http://localhost:9000/invoice-return/${filename}`
    } else {
      let dateStr = 'nodate'
      if (rentalStore.currentDetail.start_date) {
        dateStr = rentalStore.currentDetail.start_date.split('T')[0]
      }
      const filename = `${safeCustomerName}_${dateStr}_${safeInvoiceNo}.pdf`
      url = `http://localhost:9000/invoice-sewa/${filename}`
    }

    window.open(url, '_blank')
  }
}

const executeAnulir = async () => {
  if (!rentalStore.currentDetail) return
  try {
    await rentalStore.executeAnulirContract(rentalStore.currentDetail.id, cancellationFee.value)
    toast.add({
      severity: 'success',
      summary: 'Kontrak Dianulir',
      detail: 'Dokumen berhasil dibatalkan secara legal',
    })
    displayDetailDialog.value = false
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: 'Gagal Membatalkan',
      detail: rentalStore.errorMessage || 'Server menolak aksi',
    })
  }
}

const jumpToReturnForm = () => {
  displayDetailDialog.value = false
  router.push('/returns')
}

// ============================================================================
// FORMATTER & PALET DESAIN
// ============================================================================
const formatCurrency = (val?: number) => {
  if (val === undefined) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
}

const formatDateShort = (dStr?: string) => {
  if (!dStr) return '-'
  return new Date(dStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

const getCustomBadge = (st?: string) => {
  const b =
    'px-3 py-1 rounded-full text-[10px] font-black tracking-wider uppercase inline-flex items-center gap-1.5 shadow-2xs '
  switch (st) {
    case 'BOOKED':
      return b + 'bg-[#d9e2ff] text-[#001944] border border-[#b0c2ff]'
    case 'PICKED_UP':
      return b + 'bg-[#ffdeac] text-[#5f4100] border border-[#ffd18c]'
    case 'RETURNED':
      return b + 'bg-emerald-100 text-emerald-800 border border-emerald-300'
    case 'CANCELLED':
      return b + 'bg-[#ffdad6] text-[#93000a] border border-[#ffb4ab]'
    default:
      return b + 'bg-gray-100 text-gray-600 border border-gray-200'
  }
}

const getStatusIcon = (st?: string) => {
  switch (st) {
    case 'BOOKED':
      return 'pi pi-bookmark'
    case 'PICKED_UP':
      return 'pi pi-sign-out'
    case 'RETURNED':
      return 'pi pi-check-circle'
    case 'CANCELLED':
      return 'pi pi-times-circle'
    default:
      return 'pi pi-info-circle'
  }
}
</script>

<style scoped>
/* Menghilangkan garis vertikal kaku pada tabel PrimeVue */
:deep(.custom-table-rows .p-datatable-tbody > tr > td) {
  border-width: 0 0 1px 0 !important;
  border-color: rgba(196, 198, 209, 0.4) !important;
  padding: 0.85rem 1rem !important;
}
</style>
