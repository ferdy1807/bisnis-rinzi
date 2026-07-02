<template>
  <div class="space-y-8 pb-12 font-sans text-[#1a1c20]">
    <div
      class="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-[#c4c6d1]/60 pb-6"
    >
      <div>
        <p class="text-xs font-extrabold tracking-widest text-[#254582] uppercase mb-1">
          Informasi Katalog
        </p>
        <h2 class="text-3xl font-black text-[#1a1c20] tracking-tight">Peta Okupansi & Kalender</h2>
      </div>
      <div class="flex items-center gap-2">
        <Button
          label="Sorot Bulan Ini"
          class="bg-[#254582] hover:bg-[#3f5d9b] text-white border-none font-bold text-xs px-4 py-2.5 rounded-xl shadow-xs"
        />
        <Button
          icon="pi pi-refresh"
          class="p-button-rounded p-button-text text-[#58627c]"
          @click="fetchSchedule"
        />
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <div class="lg:col-span-4 space-y-6 sticky top-20">
        <div class="bg-white border border-[#c4c6d1]/60 rounded-3xl p-6 shadow-xs space-y-5">
          <div class="flex items-center gap-2.5 pb-3 border-b border-[#e2e2e7]">
            <span
              class="w-7 h-7 rounded-lg bg-[#d9e2ff] text-[#254582] flex items-center justify-center font-bold"
              ><i class="pi pi-filter text-sm"></i
            ></span>
            <h3 class="font-bold text-[#1a1c20] text-base">Filter & Petunjuk Status</h3>
          </div>

          <div class="space-y-3 text-xs">
            <div
              class="flex items-center justify-between p-3 rounded-2xl bg-[#f3f3f9] border border-gray-100"
            >
              <span class="flex items-center gap-2 font-bold text-[#254582]">
                <span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> Siap Diambil (BOOKED)
              </span>
              <span
                class="font-mono font-black text-[#254582] bg-[#d9e2ff] px-2 py-0.5 rounded-md"
                >{{ bookedCount }}</span
              >
            </div>
            <div
              class="flex items-center justify-between p-3 rounded-2xl bg-[#ffdeac]/25 border border-amber-100"
            >
              <span class="flex items-center gap-2 font-bold text-[#7e5700]">
                <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span> Sedang Keluar
                (PICKED_UP)
              </span>
              <span
                class="font-mono font-black text-[#7e5700] bg-[#ffdeac] px-2 py-0.5 rounded-md"
                >{{ pickedUpCount }}</span
              >
            </div>
          </div>

          <div
            class="p-4 bg-[#d9e2ff]/40 rounded-2xl border border-[#d9e2ff] text-xs text-[#254582] space-y-1"
          >
            <p class="font-black uppercase tracking-wider text-[10px]">Tips Operasional Pegawai</p>
            <p class="leading-relaxed text-[#58627c]">
              Periksa peta garis waktu di samping sebelum menyetujui tanggal pengambilan box baru
              agar tidak terjadi bentrok fisik di toko.
            </p>
          </div>
        </div>
      </div>

      <div
        class="lg:col-span-8 bg-white border border-[#c4c6d1]/60 rounded-3xl p-6 shadow-xs min-h-[500px] flex flex-col"
      >
        <div class="flex items-center justify-between pb-4 mb-6 border-b border-[#e2e2e7]">
          <h3 class="font-bold text-[#1a1c20] text-base">Garis Waktu Okupansi Terdekat</h3>
          <span class="text-xs text-[#58627c] font-mono font-semibold"
            >Menampilkan {{ scheduleList.length }} Agenda</span
          >
        </div>

        <div v-if="loading" class="space-y-4">
          <Skeleton height="5rem" rounded class="rounded-2xl" v-for="i in 4" :key="i" />
        </div>
        <div
          v-else-if="scheduleList.length === 0"
          class="m-auto text-center py-16 text-[#58627c] text-xs italic"
        >
          Tidak ada jadwal penyewaan aktif
        </div>
        <div v-else class="relative border-l-2 border-[#c4c6d1]/80 ml-4 pl-6 space-y-6 py-2 flex-1">
          <div v-for="item in scheduleList" :key="item.id" class="relative group">
            <span
              class="absolute -left-[31px] top-2 w-3.5 h-3.5 rounded-full border-2 border-white ring-4 ring-white transition-transform group-hover:scale-125 shadow-xs"
              :class="item.status === 'PICKED_UP' ? 'bg-amber-500' : 'bg-[#254582]'"
            ></span>

            <div
              class="bg-[#f3f3f9]/70 hover:bg-[#f3f3f9] p-4 rounded-2xl border border-[#e2e2e7] transition-all space-y-2.5"
            >
              <div class="flex justify-between items-center">
                <div class="flex items-center gap-2.5">
                  <span
                    class="font-mono text-xs font-black text-[#254582] bg-[#d9e2ff] px-2 py-0.5 rounded"
                    >{{ item.invoice_number }}</span
                  >
                  <span class="text-xs font-black text-[#1a1c20]">{{
                    item.customer_snapshot_id
                  }}</span>
                </div>
                <Tag
                  :value="item.status"
                  :severity="item.status === 'PICKED_UP' ? 'warning' : 'info'"
                  class="text-[10px]"
                />
              </div>

              <div
                class="flex items-center gap-3 text-xs text-[#58627c] pt-1 border-t border-gray-200/50"
              >
                <span class="flex items-center gap-1.5"
                  ><i class="pi pi-calendar text-[11px] text-[#254582]"></i>
                  <strong>{{ formatDate(item.start_date) }}</strong> s/d
                  <strong>{{ formatDate(item.end_date) }}</strong></span
                >
                <span>•</span>
                <span class="font-mono text-[#254582] font-bold"
                  >Durasi: {{ calcDays(item.start_date, item.end_date) }} Hari</span
                >
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { rentalApi } from '@frontend/api/rental'
import type { RentalReservation } from '@frontend/types/rental'

const loading = ref(false)
const scheduleList = ref<RentalReservation[]>([])

const fetchSchedule = async () => {
  loading.value = true
  try {
    const res = await rentalApi.getActive()
    scheduleList.value = Array.isArray(res) ? res : (res as any).data || []
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const bookedCount = computed(() => scheduleList.value.filter((s) => s.status === 'BOOKED').length)
const pickedUpCount = computed(
  () => scheduleList.value.filter((s) => s.status === 'PICKED_UP').length,
)

const formatDate = (d: string) =>
  new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
const calcDays = (s: string, e: string) =>
  Math.max(1, Math.round((new Date(e).getTime() - new Date(s).getTime()) / (1000 * 3600 * 24)))

onMounted(() => {
  fetchSchedule()
})
</script>
