<template>
    <div class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]">
      <div class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col sm:flex-row justify-between items-center gap-4">
        <h3 class="font-headline-md text-primary flex items-center gap-2">
          <span class="material-symbols-outlined">sync_alt</span> Histori Mutasi
        </h3>
        <!-- Search bar -->
        <div class="relative w-full sm:w-80">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-lg">search</span>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari Nama Produk, ID, atau Referensi..."
            class="w-full pl-10 pr-4 py-2 bg-surface border border-outline-variant/50 rounded-lg text-body-sm focus:ring-2 focus:ring-primary focus:outline-none"
          />
        </div>
      </div>

      <div class="overflow-x-auto flex-1">
        <table class="w-full text-left">
          <thead>
            <tr class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold">
              <th class="px-6 py-4">Waktu (Terbaru)</th>
              <th class="px-6 py-4">Nama Produk</th>
              <th class="px-6 py-4">Produk ID</th>
              <th class="px-6 py-4 text-right">Kuantitas</th>
              <th class="px-6 py-4">Tipe Mutasi</th>
              <th class="px-6 py-4">Keterangan (Referensi)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/20">
            <tr v-if="isLoading" class="text-center">
              <td colspan="6" class="py-8 text-on-surface-variant">Memuat data mutasi stok...</td>
            </tr>
            <tr v-else-if="filteredMovements.length === 0" class="text-center">
              <td colspan="6" class="py-8 text-on-surface-variant">Belum ada mutasi stok ditemukan.</td>
            </tr>
            <tr v-for="m in paginatedMovements" :key="m.id" class="hover:bg-surface-container-lowest/50 transition-colors">
              <td class="px-6 py-4 text-body-sm whitespace-nowrap">
                {{ m.created_at ? new Date(m.created_at).toLocaleString('id-ID') : '-' }}
              </td>
              <td class="px-6 py-4 font-medium text-body-sm">
                {{ getProductName(m.product_id) }}
              </td>
              <td class="px-6 py-4 font-mono text-xs text-on-surface-variant">
                {{ m.product_id.split('-')[0] }}...
              </td>
              <td class="px-6 py-4 text-right">
                <span :class="['font-bold', m.qty > 0 ? 'text-primary' : 'text-error']">
                  {{ m.qty > 0 ? '+' : '' }}{{ m.qty }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span :class="[
                  'px-2 py-1 rounded-md text-[10px] font-bold flex items-center gap-1 w-fit',
                  m.movement_type === 'SALE' ? 'bg-error/10 text-error' : 
                  m.movement_type === 'INITIAL' ? 'bg-primary/10 text-primary' : 
                  'bg-tertiary/10 text-tertiary'
                ]">
                  <span class="material-symbols-outlined text-[14px]">
                    {{ m.movement_type === 'SALE' ? 'trending_down' : m.movement_type === 'INITIAL' ? 'add_box' : 'edit_square' }}
                  </span>
                  {{ m.movement_type }}
                </span>
              </td>
              <td class="px-6 py-4 text-body-sm text-on-surface-variant">
                {{ m.reference || '-' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- Pagination -->
      <div v-if="filteredMovements.length > 0" class="p-4 border-t border-outline-variant/20 bg-surface-container-lowest">
        <Paginator 
          :rows="rows" 
          :totalRecords="filteredMovements.length" 
          :first="first"
          @page="onPageChange"
          :rowsPerPageOptions="[10, 20, 50]"
        ></Paginator>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { inventoryApi } from '@frontend/api/inventory';
import type { StockMovement, Product } from '@frontend/types/inventory';
import Paginator from 'primevue/paginator';

const movements = ref<StockMovement[]>([]);
const products = ref<Product[]>([]);
const isLoading = ref(true);
const searchQuery = ref('');

// Pagination states
const first = ref(0);
const rows = ref(6);

const fetchData = async () => {
  isLoading.value = true;
  try {
    const [movRes, prodRes] = await Promise.all([
      inventoryApi.getStockMovements(),
      inventoryApi.getProducts()
    ]);
    movements.value = movRes;
    products.value = prodRes;
  } catch (error) {
    console.error("Gagal mengambil pergerakan stok atau produk", error);
  } finally {
    isLoading.value = false;
  }
};

const getProductName = (id: string) => {
  const p = products.value.find(x => x.id === id);
  return p ? p.name : 'Unknown Product';
};

const filteredMovements = computed(() => {
  let list = movements.value.slice().sort((a, b) => new Date(b.created_at || 0).getTime() - new Date(a.created_at || 0).getTime());
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(m => 
      (m.reference && m.reference.toLowerCase().includes(q)) || 
      m.product_id.toLowerCase().includes(q) ||
      getProductName(m.product_id).toLowerCase().includes(q)
    );
  }
  return list;
});

const paginatedMovements = computed(() => {
  return filteredMovements.value.slice(first.value, first.value + rows.value);
});

const onPageChange = (event: any) => {
  first.value = event.first;
  rows.value = event.rows;
};

// Reset pagination when search query changes
import { watch } from 'vue';
watch(searchQuery, () => {
  first.value = 0;
});

onMounted(() => {
  fetchData();
});
</script>
