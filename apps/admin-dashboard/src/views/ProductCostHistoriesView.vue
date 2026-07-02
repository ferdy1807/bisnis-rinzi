<!-- apps/admin-dashboard/src/views/ProductCostHistoriesView.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { inventoryApi } from '@frontend/api/inventory';
import type { Product } from '@frontend/types/inventory';

// --- STATE MANAGEMENT ---
const products = ref<Product[]>([]);
const selectedProductId = ref<string>('');
const costHistories = ref<any[]>([]);

const isLoading = ref(false);
const isHistoryLoading = ref(false);
const isModalOpen = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

// Form state untuk menyuntikkan HPP baru
const form = ref({
  costPrice: 0,
  effectiveDate: new Date().toISOString().split('T')[0],
  notes: ''
});

// --- UTILITIES / FORMATTERS ---
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value);
};

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

// --- FETCH PRODUCTS (FOR SELECT BOX) ---
const fetchProducts = async () => {
  isLoading.value = true;
  try {
    products.value = await inventoryApi.getProducts();
    // Pilih produk pertama secara otomatis jika tersedia
    if (products.value.length > 0) {
      selectedProductId.value = products.value[0].id || '';
      await fetchCostHistories();
    }
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat katalog produk.';
  } finally {
    isLoading.value = false;
  }
};

// --- FETCH COST HISTORIES ---
const fetchCostHistories = async () => {
  if (!selectedProductId.value) return;
  isHistoryLoading.value = true;
  errorMessage.value = '';
  try {
    costHistories.value = await inventoryApi.getCostHistories(selectedProductId.value);
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat riwayat harga modal produk ini.';
  } finally {
    isHistoryLoading.value = false;
  }
};

// --- SUBMIT NEW COST HISTORY ---
const handleSubmitCost = async () => {
  if (form.value.costPrice <= 0) {
    errorMessage.value = 'Harga modal baru harus lebih besar dari Rp 0.';
    return;
  }

  isHistoryLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';
  
  try {
    await inventoryApi.addCostHistory(selectedProductId.value, {
      average_cost: form.value.costPrice
    });
    
    successMessage.value = 'Harga pokok produksi (HPP) baru berhasil didaftarkan ke sistem master.';
    isModalOpen.value = false;
    
    // Reset form & reload histories
    form.value = { costPrice: 0, effectiveDate: new Date().toISOString().split('T')[0], notes: '' };
    await fetchCostHistories();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memperbarui catatan HPP.';
  } finally {
    isHistoryLoading.value = false;
  }
};

// Mencari detail informasi produk yang sedang aktif dipilih
const activeProduct = computed(() => {
  return products.value.find(p => p.id === selectedProductId.value);
});

onMounted(() => {
  fetchProducts();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <!-- Header Modul -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
          <i class="pi pi-history text-indigo-600"></i>
          Riwayat Harga Modal (HPP Master)
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Lacak fluktuasi harga pokok pembelian/produksi barang. Data ini menjadi acuan utama kalkulasi margin laba kotor retail.
        </p>
      </div>
      <button 
        @click="isModalOpen = true"
        :disabled="!selectedProductId"
        class="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-400 text-white text-sm font-semibold px-4 py-2.5 rounded-xl shadow-md transition-all duration-200 disabled:cursor-not-allowed"
      >
        <i class="pi pi-plus text-xs"></i>
        Suntik Harga Modal Baru
      </button>
    </div>

    <!-- Alert Panel Message -->
    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl text-sm font-medium">
      {{ successMessage }}
    </div>

    <!-- Selector Produk Utama -->
    <div class="bg-white p-5 rounded-xl border border-slate-200 shadow-sm space-y-4">
      <div class="max-w-md space-y-1.5">
        <label class="text-xs font-bold text-slate-400 uppercase tracking-wider">Pilih Item Komoditas Produk</label>
        <div class="relative">
          <select 
            v-model="selectedProductId" 
            @change="fetchCostHistories"
            class="w-full pl-3 pr-10 py-2.5 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all appearance-none font-semibold text-slate-800"
          >
            <option v-for="product in products" :key="product.id" :value="product.id">
              {{ product.name }} ({{ product.sku || 'No SKU' }})
            </option>
          </select>
          <span class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-slate-400">
            <i class="pi pi-chevron-down text-xs"></i>
          </span>
        </div>
      </div>

      <!-- Info Singkat Atribut Produk yang Dipilih -->
      <div v-if="activeProduct" class="pt-3 border-t border-slate-100 grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs font-medium">
        <div>
          <span class="text-slate-400 block">ID Sistem</span>
          <span class="text-slate-800 font-mono select-all">{{ activeProduct.id }}</span>
        </div>
        <div>
          <span class="text-slate-400 block">Stok Saat Ini</span>
          <span class="text-slate-800 font-bold">{{ activeProduct.qty || 0 }} Unit</span>
        </div>
        <div>
          <span class="text-slate-400 block">Harga Jual Retail</span>
          <span class="text-indigo-600 font-mono font-bold">{{ formatCurrency(activeProduct.selling_price || 0) }}</span>
        </div>
      </div>
    </div>

    <!-- Tabel Histori Log Perubahan Harga Pokok -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isHistoryLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Mengekstrak log HPP produk...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">Tanggal Efektif Berlaku</th>
              <th class="py-4 px-6">Catatan / Alasan Perubahan</th>
              <th class="py-4 px-6 text-right">Nilai Harga Modal (HPP)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="costHistories.length === 0 && !isHistoryLoading">
              <td colspan="3" class="py-12 text-center text-slate-400 font-normal">
                Belum ada rekam jejak perubahan HPP untuk produk ini.
              </td>
            </tr>
            <tr v-for="(history, index) in costHistories" :key="index" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6 font-semibold text-slate-900">
                {{ formatDate(history.effective_date || history.created_at) }}
                <span v-if="index === 0" class="ml-2 px-2 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] rounded-full">Aktif Saat Ini</span>
              </td>
              <td class="py-4 px-6 text-slate-500 font-normal max-w-sm truncate" title="Penyesuaian Sistem">
                Penyesuaian HPP Baru
              </td>
              <td class="py-4 px-6 text-right font-mono font-bold text-slate-900">
                {{ formatCurrency(history.average_cost) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL POPUP ADD COST HISTORY -->
    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-200 overflow-hidden transform transition-all">
        <div class="p-6 border-b border-slate-100 flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-900">Suntik Harga Pokok Produksi Baru</h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600">
            <i class="pi pi-times"></i>
          </button>
        </div>
        
        <form @submit.prevent="handleSubmitCost" class="p-6 space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Nilai Harga Modal Baru (Rupiah)</label>
            <input 
              v-model.number="form.costPrice"
              type="number"
              required
              min="1"
              placeholder="Masukkan nominal modal bersih..."
              class="w-full px-3 py-2 font-mono border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Tanggal Mulai Berlaku</label>
            <input 
              v-model="form.effectiveDate"
              type="date"
              required
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Keterangan Penyesuaian</label>
            <textarea 
              v-model="form.notes"
              rows="3"
              placeholder="Contoh: Kenaikan harga beli dari supplier utama akibat inflasi material vendor..."
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all resize-none"
            ></textarea>
          </div>

          <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-100 mt-6">
            <button type="button" @click="isModalOpen = false" class="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 rounded-lg border border-slate-200 transition-colors">
              Batal
            </button>
            <button type="submit" class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-md transition-colors">
              Simpan & Terapkan HPP
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>