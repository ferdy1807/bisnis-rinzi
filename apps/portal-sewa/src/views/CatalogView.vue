<template>
  <div class="space-y-8 pb-12 font-sans text-[#1a1c20]">
    <div
      class="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-[#c4c6d1]/60 pb-6"
    >
      <div>
        <p class="text-xs font-extrabold tracking-widest text-[#254582] uppercase mb-1">
          Informasi Katalog
        </p>
        <h2 class="text-3xl font-black text-[#1a1c20] tracking-tight">
          Katalog Master Box Hantaran
        </h2>
      </div>
      <Button
        label="Perbarui Live Stock"
        icon="pi pi-refresh"
        class="p-button-outlined p-button-sm text-[#58627c] border-[#c4c6d1] rounded-xl font-bold"
        :loading="loading"
        @click="loadProducts"
      />
    </div>

    <div
      class="bg-[#d9e2ff]/50 border border-[#d9e2ff] p-5 rounded-3xl flex items-center gap-4 text-[#001944]"
    >
      <span class="p-3.5 bg-[#254582] text-white rounded-2xl shadow-xs"
        ><i class="pi pi-lock text-xl"></i
      ></span>
      <div class="text-xs space-y-0.5 leading-relaxed">
        <span class="font-black uppercase tracking-wider block text-[10px] text-[#254582]"
          >Mode Akses Pegawai (Read-Only)</span
        >
        <p class="text-[#444650]">
          Anda sedang melihat tarif resmi dan sisa kuota fisik secara <em>real-time</em>. Penambahan
          tipe box baru atau modifikasi harga pokok hanya dapat dilakukan oleh Owner melalui panel
          <strong>admin-dashboard</strong>.
        </p>
      </div>
    </div>

    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div
        v-for="i in 4"
        :key="i"
        class="bg-white border border-[#c4c6d1]/60 rounded-3xl p-4 space-y-3 shadow-2xs"
      >
        <Skeleton height="12rem" class="rounded-2xl" />
        <Skeleton width="75%" height="1.2rem" />
        <Skeleton width="45%" />
      </div>
    </div>

    <div
      v-else-if="products.length === 0"
      class="text-center py-16 bg-white border border-[#c4c6d1]/60 rounded-3xl text-[#58627c] text-xs italic shadow-2xs"
    >
      <i class="pi pi-box text-4xl mb-2 block text-[#c4c6d1]"></i> Belum ada produk hantaran yang
      terdaftar di dalam database
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div
        v-for="prod in products"
        :key="prod.id"
        class="bg-white border border-[#c4c6d1]/60 rounded-3xl overflow-hidden hover:-translate-y-1.5 hover:shadow-xl hover:border-[#254582] transition-all duration-300 flex flex-col group shadow-2xs"
      >
        <div class="h-48 bg-[#f3f3f9] relative flex items-center justify-center overflow-hidden">
          <i
            class="pi pi-image text-4xl text-[#c4c6d1] group-hover:scale-110 transition-transform duration-500"
          ></i>
          <div class="absolute top-3 right-3">
            <span
              class="px-3 py-1 rounded-xl text-[10px] font-black uppercase tracking-wider backdrop-blur-md shadow-xs inline-block"
              :class="
                prod.quantity_available > 0
                  ? 'bg-green-700 text-white'
                  : 'bg-[#ba1a1a] text-white animate-pulse'
              "
            >
              {{ prod.quantity_available > 0 ? `${prod.quantity_available} Tersedia` : 'Habis' }}
            </span>
          </div>

          <span
            class="absolute bottom-2.5 left-3 bg-white/90 backdrop-blur-md px-2.5 py-0.5 rounded-md text-[10px] font-extrabold text-[#254582]"
          >
            {{ prod.category_id || 'Box Premium' }}
          </span>
        </div>

        <div class="p-5 flex-1 flex flex-col justify-between space-y-4">
          <div>
            <h3 class="font-black text-[#1a1c20] text-base line-clamp-1 mb-1">{{ prod.name }}</h3>
            <p class="text-xs text-[#58627c] line-clamp-2 leading-relaxed">
              Paket nampan seserahan eksklusif Rinzi siap hias kelengkapan pengantin.
            </p>
          </div>

          <div class="pt-3.5 border-t border-[#e2e2e7] flex items-center justify-between">
            <div>
              <span class="text-[10px] text-[#58627c] block uppercase font-bold tracking-wider"
                >Tarif Sewa</span
              >
              <span class="text-[#254582] font-mono font-black text-lg">{{
                formatCurrency(prod.rental_price)
              }}</span>
            </div>
            <span
              class="text-[10px] font-black bg-[#d9e2ff] text-[#254582] px-2.5 py-1.5 rounded-xl"
            >
              / Event
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { rentalApi } from '@frontend/api/rental'
import type { RentalProduct } from '@frontend/types/rental'

const loading = ref(false)
const products = ref<RentalProduct[]>([])

const loadProducts = async () => {
  loading.value = true
  try {
    const res = await rentalApi.getProducts()
    products.value = Array.isArray(res) ? res : (res as any).data || []
  } catch (e) {
  } finally {
    loading.value = false
  }
}

const formatCurrency = (val: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val)

onMounted(() => {
  loadProducts()
})
</script>
