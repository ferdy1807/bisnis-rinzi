<!-- apps/portal-sewa/src/views/OverdueView.vue -->
<template>
  <div class="space-y-8 pb-12 font-sans text-[#1a1c20]">
    <!-- HEADER -->
    <div
      class="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-[#c4c6d1]/60 pb-6"
    >
      <div>
        <p class="text-xs font-extrabold tracking-widest text-[#ba1a1a] uppercase mb-1">
          Peringatan Keamanan Aset
        </p>
        <h2 class="text-3xl font-black text-[#93000a] tracking-tight flex items-center gap-2.5">
          <i class="pi pi-exclamation-triangle"></i> Kontrol Keterlambatan (Overdue)
        </h2>
      </div>
      <Button
        label="Perbarui Data Deteksi"
        icon="pi pi-refresh"
        class="bg-[#ba1a1a] hover:bg-[#93000a] text-white border-none font-bold text-xs px-5 py-3 rounded-xl shadow-xs"
        :loading="loading"
        @click="loadOverdue"
      />
    </div>

    <Toast />

    <!-- URGENT BENTO NOTICE CARD (Replikasi "Critical Alerts" Desain Anda) -->
    <div
      class="bg-[#ffdad6]/40 border-2 border-[#ba1a1a]/30 rounded-3xl p-6 shadow-xs flex flex-col md:flex-row items-start md:items-center justify-between gap-4"
    >
      <div class="flex items-center gap-4">
        <span class="p-3.5 bg-[#ba1a1a] text-white rounded-2xl shadow-xs"
          ><i class="pi pi-bell text-2xl animate-bounce"></i
        ></span>
        <div class="space-y-0.5">
          <h3 class="font-extrabold text-[#93000a] text-base uppercase tracking-wide">
            SOP Penarikan Box Hantaran Terlambat
          </h3>
          <p class="text-xs text-[#444650] leading-relaxed">
            Penyewa di bawah ini telah melewati batas akhir pengembalian. Klik tombol
            <strong>"Eksekusi Penarikan"</strong> untuk mengalihkan dokumen ke meja kasir pelunasan
            beserta denda Rp50.000/hari.
          </p>
        </div>
      </div>
      <div
        class="bg-white px-5 py-3 rounded-2xl border border-[#ba1a1a]/20 text-center shadow-2xs w-full md:w-auto"
      >
        <span class="text-[10px] font-extrabold text-[#ba1a1a] block uppercase tracking-wider"
          >Aset Dalam Risiko:</span
        >
        <span class="text-2xl font-black text-[#93000a] font-mono"
          >{{ overdueList.length }} <span class="text-xs font-bold">Nota</span></span
        >
      </div>
    </div>

    <!-- OVERDUE DATATABLE -->
    <div
      class="bg-white border-2 border-[#ba1a1a]/20 rounded-3xl overflow-hidden shadow-xs flex flex-col"
    >
      <div
        class="p-6 border-b border-[#c4c6d1]/40 flex items-center justify-between bg-[#ffdad6]/25"
      >
        <h3 class="font-bold text-[#93000a] text-base">Daftar Penyewa Melewati Tenggat</h3>
        <span class="text-xs font-black text-[#ba1a1a] uppercase tracking-wider animate-pulse"
          >SIAGA PENAGIHAN</span
        >
      </div>

      <DataTable
        :value="overdueList"
        :loading="loading"
        paginator
        :rows="10"
        class="p-datatable-sm text-xs border-none"
        responsiveLayout="scroll"
      >
        <template #empty>
          <div class="text-center py-16 bg-green-50/50">
            <i class="pi pi-verified text-4xl text-green-600 mb-2 block"></i>
            <span class="text-sm font-bold text-green-800"
              >Kondisi Fisik Sempurna! Tidak ada keterlambatan pengembalian box hari ini.</span
            >
          </div>
        </template>

        <Column
          field="invoice_number"
          header="No. Invoice"
          class="font-mono text-[#93000a] font-black"
        ></Column>
        <Column
          field="customer_snapshot_id"
          header="Penyewa (Snapshot)"
          class="font-bold text-[#1a1c20]"
        ></Column>
        <Column header="Tgl Diambil">
          <template #body="slotProps">{{ formatDate(slotProps.data.start_date) }}</template>
        </Column>
        <Column header="Tenggat Seharusnya">
          <template #body="slotProps"
            ><span class="font-black text-[#ba1a1a] underline">{{
              formatDate(slotProps.data.end_date)
            }}</span></template
          >
        </Column>
        <Column header="Lama Menunggak" class="text-center">
          <template #body="slotProps">
            <span
              class="bg-[#ba1a1a] text-white font-black px-2.5 py-1 rounded-md font-mono text-xs shadow-2xs"
            >
              +{{ calculateDaysOverdue(slotProps.data.end_date) }} Hari
            </span>
          </template>
        </Column>
        <Column header="Potensi Denda" class="text-right">
          <template #body="slotProps">
            <span class="font-mono font-bold text-[#93000a]">
              {{ formatCurrency(calculateDaysOverdue(slotProps.data.end_date) * 50000) }}
            </span>
          </template>
        </Column>
        <Column header="Tindakan" class="w-36 text-center">
          <template #body="slotProps">
            <Button
              label="Eksekusi Retur"
              icon="pi pi-replay"
              class="p-button-xs bg-[#ba1a1a] hover:bg-[#93000a] text-white border-none font-bold text-[11px] py-1.5 px-3 rounded-lg shadow-2xs"
              @click="jumpToReturns(slotProps.data)"
            />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { rentalApi } from '@frontend/api/rental'
import type { RentalReservation } from '@frontend/types/rental'
import { useToast } from 'primevue/usetoast'
import { useRouter } from 'vue-router'

const toast = useToast()
const router = useRouter()
const loading = ref(false)
const overdueList = ref<RentalReservation[]>([])

const loadOverdue = async () => {
  loading.value = true
  try {
    const res = await rentalApi.getOverdue()
    overdueList.value = Array.isArray(res) ? res : (res as any).data || []
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: 'Gagal',
      detail: 'Sistem gagal mengidentifikasi data keterlambatan',
      life: 3000,
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOverdue()
})

const calculateDaysOverdue = (endStr: string) => {
  if (!endStr) return 0
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const end = new Date(endStr)
  end.setHours(0, 0, 0, 0)
  const diff = Math.floor((today.getTime() - end.getTime()) / (1000 * 3600 * 24))
  return diff > 0 ? diff : 0
}

// Lompat cerdas: Membawa pegawai langsung ke halaman pengembalian
const jumpToReturns = (item: RentalReservation) => {
  router.push('/returns')
}

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)
const formatDate = (d: string) =>
  new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
</script>
