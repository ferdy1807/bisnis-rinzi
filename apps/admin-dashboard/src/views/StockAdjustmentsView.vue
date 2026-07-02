<!-- apps/admin-dashboard/src/views/StockAdjustmentsView.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { inventoryApi } from '@frontend/api/inventory';
import type { Product } from '@frontend/types/inventory';

// --- STATE MANAGEMENT ---
const products = ref<Product[]>([]);
const lowStockProducts = ref<Product[]>([]);
const isLoading = ref(false);
const isSubmitLoading = ref(false);
const isModalOpen = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

// Form State penyesuaian stok
const selectedProductId = ref('');
const form = ref({
  adjustmentQty: 0, // bisa positif (tambah) atau negatif (kurang)
  reason: ''
});

// --- FETCH DATA ---
const fetchData = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    // Ambil data produk secara paralel
    const [allProds, lowProds] = await Promise.all([
      inventoryApi.getProducts(),
      inventoryApi.getLowStockProducts()
    ]);
    products.value = allProds;
    lowStockProducts.value = lowProds;
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal mengambil data dari server inventaris.';
  } finally {
    isLoading.value = false;
  }
};

// --- SUBMIT ADJUSTMENT ---
const handleAdjustStock = async () => {
  if (!selectedProductId.value) {
    errorMessage.value = 'Silakan pilih produk terlebih dahulu.';
    return;
  }
  if (form.value.adjustmentQty === 0) {
    errorMessage.value = 'Kuantitas penyesuaian tidak boleh nol.';
    return;
  }
  if (!form.value.reason.trim()) {
    errorMessage.value = 'Alasan penyesuaian (keterangan audit) wajib diisi.';
    return;
  }

  isSubmitLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    const currentStock = (selectedProductDetail.value as any)?.qty || 0;
    const calculatedNewQty = currentStock + form.value.adjustmentQty;

    // Memanggil API adjustStock dari InventoryApi
    await inventoryApi.adjustStock({
      product_id: selectedProductId.value,
      new_qty: calculatedNewQty,
      reference: form.value.reason
    });

    successMessage.value = 'Penyesuaian stok fisik berhasil direkam dan divalidasi ke sistem ledger.';
    isModalOpen.value = false;
    
    // Reset Form & Refresh Data
    selectedProductId.value = '';
    form.value = { adjustmentQty: 0, reason: '' };
    await fetchData();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal mengeksekusi penyesuaian stok.';
  } finally {
    isSubmitLoading.value = false;
  }
};

const openAdjustmentModal = (productId: string = '') => {
  selectedProductId.value = productId;
  form.value = { adjustmentQty: 0, reason: '' };
  errorMessage.value = '';
  successMessage.value = '';
  isModalOpen.value = true;
};

// Mencari info produk ter-update berdasarkan ID yang dipilih dalam form modal
const selectedProductDetail = computed(() => {
  return products.value.find(p => p.id === selectedProductId.value);
});

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <!-- Header Modul -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
          <i class="pi pi-sliders-h text-indigo-600"></i>
          Koreksi & Aturan Stok (Stock Opname)
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Otoritas mutlak Owner dalam menyesuaikan kuantitas sistem dengan kondisi riil fisik barang di gudang.
        </p>
      </div>
      <button 
        @click="openAdjustmentModal('')"
        class="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-semibold px-4 py-2.5 rounded-xl shadow-md transition-all duration-200"
      >
        <i class="pi pi-plus text-xs"></i>
        Buat Penyesuaian Baru
      </button>
    </div>

    <!-- Banner Notifikasi -->
    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl text-sm font-medium">
      {{ successMessage }}
    </div>

    <!-- SEKSI ALERT: PRODUK STOK MENIPIS / KRITIS -->
    <div class="bg-amber-50 border border-amber-200 rounded-xl p-5 space-y-3 shadow-sm">
      <div class="flex items-center gap-2 text-amber-800 font-bold text-sm">
        <i class="pi pi-exclamation-triangle text-base"></i>
        <span>Peringatan Dini: Produk Kritis di Bawah Ambang Batas Minimum ({{ lowStockProducts.length }})</span>
      </div>
      
      <div v-if="lowStockProducts.length === 0" class="text-xs text-amber-700 font-medium">
        Kondisi gudang aman. Seluruh produk saat ini berada di atas batas minimum keselamatan stok (*safety stock*).
      </div>
      
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
        <div 
          v-for="prod in lowStockProducts" 
          :key="prod.id" 
          class="bg-white border border-amber-100 rounded-lg p-3 flex items-center justify-between shadow-xs"
        >
          <div class="min-w-0 pr-2">
            <p class="text-xs font-bold text-slate-900 truncate">{{ prod.name }}</p>
            <p class="text-[10px] text-slate-400 font-mono mt-0.5">SKU: {{ prod.sku || '-' }}</p>
          </div>
          <div class="flex items-center gap-3 shrink-0">
            <span class="text-xs px-2 py-1 bg-rose-50 text-rose-700 rounded-md font-extrabold font-mono">
              {{ (prod as any).stock || 0 }} Pcs
            </span>
            <button 
              @click="openAdjustmentModal(prod.id)"
              class="p-1.5 text-indigo-600 hover:bg-indigo-50 rounded transition-colors"
              title="Koreksi Langsung"
            >
              <i class="pi pi-pencil text-xs"></i>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TABEL LOG DATA UTAMA PRODUK UNTUK ADJ -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Memuat status inventaris master...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">ID Produk</th>
              <th class="py-4 px-6">Nama Komoditas</th>
              <th class="py-4 px-6">SKU / Kode Bar</th>
              <th class="py-4 px-6 text-center">Stok Berjalan</th>
              <th class="py-4 px-6 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="products.length === 0 && !isLoading">
              <td colspan="5" class="py-12 text-center text-slate-400">Belum ada komoditas terdaftar dalam katalog master.</td>
            </tr>
            <tr v-for="product in products" :key="product.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6 font-mono text-xs text-slate-400 select-all">{{ product.id }}</td>
              <td class="py-4 px-6 text-slate-900 font-semibold">{{ product.name }}</td>
              <td class="py-4 px-6 font-mono text-xs text-slate-500">{{ product.sku || '-' }}</td>
              <td class="py-4 px-6 text-center font-mono font-bold">
                {{ (product as any).qty || 0 }}
              </td>
              <td class="py-4 px-6 text-right">
                <button 
                  @click="openAdjustmentModal(product.id)"
                  class="inline-flex items-center gap-1.5 text-xs text-indigo-600 hover:text-indigo-700 bg-indigo-50 hover:bg-indigo-100 font-bold px-2.5 py-1.5 rounded-lg transition-all">
                  <i class="pi pi-cog text-[10px]"></i> Opname
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL POPUP FORM PENYESUAIAN STOK (OPNAME) -->
    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-200 overflow-hidden transform transition-all">
        <div class="p-6 border-b border-slate-100 flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-900">Form Formulir Penyesuaian Stok</h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600"><i class="pi pi-times"></i></button>
        </div>
        
        <form @submit.prevent="handleAdjustStock" class="p-6 space-y-4">
          <!-- Dropdown Pemilihan Produk (Hanya jika belum dipilih dari tabel) -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Komoditas Produk</label>
            <select 
              v-model="selectedProductId"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 transition-all font-semibold"
            >
              <option value="" disabled>-- Pilih Item --</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }} (Stok: {{ (p as any).qty || 0 }})
              </option>
            </select>
          </div>

          <!-- Kuantitas Penyesuaian -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Kuantitas Koreksi Selisih</label>
            <input 
              v-model.number="form.adjustmentQty"
              type="number"
              required
              placeholder="Masukkan angka contoh: -5 untuk susut, 10 untuk tambah..."
              class="w-full px-3 py-2 font-mono border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
            />
            <span class="text-[10px] text-slate-400 block leading-normal">
              Gunakan tanda **minus (-)** jika barang fisik di gudang hilang/rusak/berkurang, dan angka **positif** jika terdapat kelebihan fisik.
            </span>
          </div>

          <!-- Alasan Opname -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Alasan & Dokumen Keterangan Audit</label>
            <textarea 
              v-model="form.reason"
              rows="3"
              required
              placeholder="Contoh: Barang kedaluwarsa, Rusak di rak display, Selisih hitung timbangan supplier..."
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all resize-none"
            ></textarea>
          </div>

          <!-- Aksi Eksekusi Form -->
          <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-100 mt-6">
            <button type="button" @click="isModalOpen = false" class="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 rounded-lg border border-slate-200 transition-colors">
              Batal
            </button>
            <button 
              type="submit" 
              :disabled="isSubmitLoading"
              class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-md transition-colors disabled:opacity-50"
            >
              <span>{{ isSubmitLoading ? 'Memproses Audit...' : 'Terapkan Opname Stok' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>