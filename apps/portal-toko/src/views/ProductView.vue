<template>
  <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6 border-b border-outline-variant/20 pb-4">
    <div>
      <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">Katalog Produk</h2>
      <p class="text-body-md text-on-surface-variant">
        Kelola daftar produk, stok, dan harga barang ritel toko.
      </p>
    </div>
    <div>
      <button
        @click="openAddStockModal"
        class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2.5 rounded-xl shadow-md hover:bg-primary-container hover:text-on-primary-container transition-all cursor-pointer"
      >
        <span class="material-symbols-outlined">inventory</span>
        <span class="font-bold text-label-md">Tambah Stok (Inbound)</span>
      </button>
    </div>
  </div>

  <div class="grid grid-cols-1 gap-6">
    <div class="space-y-8">
      <div class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]">
        
        <div class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col lg:flex-row justify-between items-center gap-4">
          <div class="flex flex-col sm:flex-row sm:items-center gap-4 w-full lg:w-auto">
            <h3 class="font-headline-md text-primary flex items-center gap-2 shrink-0 font-bold">
              <span class="material-symbols-outlined">list_alt</span> Daftar Produk
            </h3>
            
            <div class="flex items-center gap-2 text-xs font-mono font-black">
              <span class="px-3 py-1.5 bg-primary/10 text-primary rounded-xl border border-primary/20">
                SKU MASTER: {{ products.length }} ITEM
              </span>
              <span class="px-3 py-1.5 bg-secondary/10 text-secondary rounded-xl border border-secondary/20">
                TOTAL STOK: {{ totalPhysicalStock.toLocaleString('id-ID') }} UNIT
              </span>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row gap-2 w-full lg:w-auto">
            <div class="relative w-full sm:w-64">
              <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-lg">search</span>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Cari SKU atau Nama..."
                class="w-full pl-10 pr-4 py-2 bg-surface border border-outline-variant/50 rounded-lg text-body-sm focus:ring-2 focus:ring-primary focus:outline-none"
              />
            </div>
            <div class="relative w-full sm:w-48">
              <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-lg">barcode_scanner</span>
              <input
                v-model="barcodeQuery"
                @keyup.enter="searchByBarcode"
                type="text"
                placeholder="Scan Barcode & Enter"
                class="w-full pl-10 pr-4 py-2 bg-surface border border-outline-variant/50 rounded-lg text-body-sm focus:ring-2 focus:ring-primary focus:outline-none"
              />
            </div>
          </div>
        </div>

        <div class="overflow-x-auto flex-1">
          <table class="w-full text-left">
            <thead>
              <tr class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-black select-none border-b">
                <th class="px-6 py-4">Produk</th>
                <th class="px-6 py-4">SKU</th>
                
                <th @click="toggleSortCategory" class="px-6 py-4 cursor-pointer hover:bg-surface-container-high transition-colors group">
                  <div class="flex items-center gap-1">
                    <span>Kategori</span>
                    <span class="material-symbols-outlined text-sm text-outline group-hover:text-primary">
                      {{ sortOrder === 'asc' ? 'arrow_upward' : sortOrder === 'desc' ? 'arrow_downward' : 'swap_vert' }}
                    </span>
                  </div>
                </th>
                
                <th class="px-6 py-4">Harga Jual</th>
                <th class="px-6 py-4">Stok</th>
                <th class="px-6 py-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-outline-variant/10 text-xs font-semibold">
              <tr v-if="isLoading">
                <td colspan="6" class="px-6 py-12 text-center text-primary">
                  <span class="animate-spin material-symbols-outlined text-4xl">sync</span>
                </td>
              </tr>
              <tr v-else-if="filteredProducts.length === 0">
                <td colspan="6" class="px-6 py-12 text-center text-on-surface-variant font-medium italic">
                  Produk tidak ditemukan di dalam database inventory_db.
                </td>
              </tr>
              <tr
                v-else
                v-for="p in paginatedProducts"
                :key="p.id"
                class="hover:bg-surface-container-low/50 transition-colors group"
              >
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-lg border border-outline-variant/30 overflow-hidden bg-surface-container flex-shrink-0">
                      <img v-if="p.image" :src="p.image" alt="Product Image" class="w-full h-full object-cover" />
                      <span v-else class="material-symbols-outlined w-full h-full flex items-center justify-center text-on-surface-variant/50 text-xl">inventory_2</span>
                    </div>
                    <div class="font-bold text-on-surface text-sm">{{ p.name }}</div>
                  </div>
                </td>
                <td class="px-6 py-4 font-mono text-primary font-bold">{{ p.sku }}</td>
                <td class="px-6 py-4 text-on-surface-variant font-bold">
                  {{ getCategoryName(p.category_id) }}
                </td>
                <td class="px-6 py-4 font-mono font-bold text-sm">
                  Rp {{ p.selling_price?.toLocaleString('id-ID') }}
                </td>
                <td class="px-6 py-4">
                  <span :class="['px-2 py-1 rounded-md font-mono font-bold text-xs', (p.qty || 0) < 10 ? 'bg-error/10 text-error' : 'bg-primary/10 text-primary']">
                    {{ p.qty || 0 }} {{ p.base_unit_code }}
                  </span>
                </td>
                <td class="px-6 py-4 text-center">
                  <router-link :to="`/products/detail/${p.id}`" class="p-2 text-primary hover:bg-primary/10 rounded-full transition-colors inline-flex items-center justify-center" title="Edit/Detail Produk">
                    <span class="material-symbols-outlined text-sm">edit</span>
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div v-if="filteredProducts.length > 0" class="p-4 border-t border-outline-variant/20 bg-surface-container-lowest">
          <Paginator :rows="rows" :totalRecords="filteredProducts.length" :first="first" @page="onPageChange" :rowsPerPageOptions="[5, 10, 20]"></Paginator>
        </div>
      </div>
    </div>
  </div>

  <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
    <div class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl relative flex flex-col border border-outline-variant/30">
      <button @click="showAddModal = false" class="absolute top-5 right-5 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1 transition-colors border-none bg-transparent cursor-pointer">
        <span class="material-symbols-outlined">close</span>
      </button>

      <h3 class="text-title-lg font-black text-primary mb-5 flex items-center gap-2 shrink-0">
        <span class="material-symbols-outlined">inventory</span> Tambah Stok Masuk
      </h3>

      <form @submit.prevent="submitForm" class="space-y-4">
        <div class="space-y-1.5 relative">
          <label class="text-[11px] font-black text-on-surface-variant uppercase tracking-wider block">
            Pilih Produk (Cari Nama / SKU) *
          </label>
          
          <div class="relative">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-md">search</span>
            <input 
              v-model="modalSearchQuery"
              type="text"
              placeholder="Ketik kata kunci nama atau SKU..."
              class="w-full pl-9 pr-4 py-2 bg-surface-container border border-outline-variant/50 rounded-lg text-body-sm focus:ring-1 focus:ring-primary focus:outline-none font-semibold text-slate-800"
              @focus="isDropdownOpen = true"
            />
          </div>

          <div v-if="isDropdownOpen && filteredModalProducts.length > 0" class="absolute left-0 right-0 mt-1 max-h-40 bg-surface border border-outline-variant shadow-lg rounded-xl overflow-y-auto z-20 divide-y divide-slate-100">
            <div 
              v-for="prd in filteredModalProducts" 
              :key="prd.id"
              @click="selectProductForStock(prd)"
              class="p-2.5 hover:bg-primary/5 cursor-pointer transition-colors text-xs font-semibold flex flex-col gap-0.5"
            >
              <div class="text-slate-900 font-bold">{{ prd.name }}</div>
              <div class="text-primary font-mono text-[10px] uppercase tracking-wider">SKU: {{ prd.sku }} | Stok Saat Ini: {{ prd.qty || 0 }}</div>
            </div>
          </div>
          
          <div v-if="selectedProductLabel" class="mt-2 p-2.5 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-xl text-xs font-bold flex items-center justify-between">
            <div class="min-w-0 flex-1 pr-2">
              <span class="block text-slate-500 text-[10px] uppercase tracking-wider font-mono">Produk Terpilih Valid:</span>
              <span class="block text-slate-800 font-black truncate">{{ selectedProductLabel }}</span>
            </div>
            <span class="material-symbols-outlined text-emerald-600 text-sm">verified</span>
          </div>
        </div>

        <div>
          <label class="text-[11px] font-black text-on-surface-variant uppercase tracking-wider block mb-1">
            Jumlah Stok Masuk (+) *
          </label>
          <input v-model.number="stockForm.new_qty" type="number" min="1" required class="w-full bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-primary focus:outline-none font-mono font-bold" placeholder="Minimal 1" />
          <p class="text-[10px] text-error font-bold mt-1">
            * Hanya untuk penambahan stok fisik. Tidak bisa mengurangi stok!
          </p>
        </div>

        <div>
          <label class="text-[11px] font-black text-on-surface-variant uppercase tracking-wider block mb-1">
            Keterangan / Referensi
          </label>
          <input v-model="stockForm.reference" type="text" required class="w-full bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-primary focus:outline-none font-semibold text-slate-700" placeholder="Contoh: Titipan Barang / DO-001" />
        </div>

        <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-outline-variant/20">
          <button type="button" @click="showAddModal = false" class="px-4 py-2 text-xs font-bold text-on-surface-variant hover:bg-surface-variant/20 rounded-lg transition-colors border-none bg-transparent cursor-pointer">
            Batal
          </button>
          <button type="submit" class="px-5 py-2 bg-primary text-on-primary rounded-xl font-bold text-xs shadow-md hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all flex items-center gap-2 disabled:opacity-70 cursor-pointer" :disabled="isSubmitting || !stockForm.product_id">
            <span v-if="isSubmitting" class="animate-spin material-symbols-outlined text-[16px]">sync</span>
            <span v-else class="material-symbols-outlined text-[16px]">add</span>
            Tambah Stok
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { inventoryApi } from '@frontend/api/inventory'
import { cashApi } from '@frontend/api/cash'
import type {
  Product,
  Category,
  Brand,
  Unit,
  StockMovement,
} from '@frontend/types/inventory'
import type { CashierSession } from '@frontend/types/cash'
import Paginator from 'primevue/paginator'

// --- CORE COMPONENT STATES ---
const isLoading = ref(true)
const isSubmitting = ref(false)
const showAddModal = ref(false)

const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const brands = ref<Brand[]>([])
const units = ref<Unit[]>([])
const lowStockProducts = ref<Product[]>([])
const movements = ref<StockMovement[]>([])

const currentSession = ref<CashierSession | null>(null)
const currentTime = ref('')
const shiftDuration = ref('')
let timerInterval: number | null = null

// --- SEARCH & SORT LOGIC STATES ---
const searchQuery = ref('')
const barcodeQuery = ref('')
const sortOrder = ref<'none' | 'asc' | 'desc'>('none') // State pengurutan kategori

// Search states inside stock modal
const modalSearchQuery = ref('')
const isDropdownOpen = ref(false)
const selectedProductLabel = ref('')

// Pagination states
const first = ref(0)
const rows = ref(5)

const onPageChange = (event: any) => {
  first.value = event.first
  rows.value = event.rows
}

watch(searchQuery, () => {
  first.value = 0
})

// Perbarui Real-time jam digital operasional toko
const updateClock = () => {
  const now = new Date()
  currentTime.value = now.toLocaleString('id-ID', {
    weekday: 'long', day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit',
  })

  if (currentSession.value && currentSession.value.open_time) {
    const start = new Date(currentSession.value.open_time).getTime()
    const diff = Math.floor((now.getTime() - start) / 1000)
    const h = Math.floor(diff / 3600)
    const m = Math.floor((diff % 3600) / 60)
    const s = diff % 60
    shiftDuration.value = `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  } else {
    shiftDuration.value = '00:00:00'
  }
}

// Perhitungan Akumulasi KPI Fisik Berjalan
const totalPhysicalStock = computed(() => {
  return products.value.reduce((sum, p) => sum + (p.qty || 0), 0)
})

// TUGAS 3: Manajemen Algoritma Pengurutan (Sort) Berdasarkan Kategori
const toggleSortCategory = () => {
  if (sortOrder.value === 'none') sortOrder.value = 'asc'
  else if (sortOrder.value === 'asc') sortOrder.value = 'desc'
  else sortOrder.value = 'none'
}

const getCategoryName = (id: string) => {
  const cat = categories.value.find((c) => c.id === id)
  return cat ? cat.name : '-'
}

// Pemfilteran Pencarian & Pengurutan Berbasis Kategori secara Sempurna
const filteredProducts = computed(() => {
  let result = [...products.value]

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase().trim()
    result = result.filter(
      (p) => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q),
    )
  }

  if (sortOrder.value !== 'none') {
    result.sort((a, b) => {
      const nameA = getCategoryName(a.category_id).toLowerCase()
      const nameB = getCategoryName(b.category_id).toLowerCase()
      if (sortOrder.value === 'asc') return nameA.localeCompare(nameB)
      return nameB.localeCompare(nameA)
    })
  }

  return result
})

const paginatedProducts = computed(() => {
  return filteredProducts.value.slice(first.value, first.value + rows.value)
})

// TUGAS 4: Filter Dropdown Terbuka untuk Memudahkan Pencarian Inbound Stok Banyak
const filteredModalProducts = computed(() => {
  if (!modalSearchQuery.value) return products.value
  const q = modalSearchQuery.value.toLowerCase().trim()
  return products.value.filter(p => 
    p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q)
  )
})

const selectProductForStock = (prd: Product) => {
  stockForm.value.product_id = prd.id
  selectedProductLabel.value = `${prd.name} (${prd.sku})`
  modalSearchQuery.value = prd.name
  isDropdownOpen.value = false
}

// Inisialisasi Ulang Skema Form Tambah Stok
const defaultStockForm = () => ({
  product_id: '',
  new_qty: 1,
  reference: 'Penerimaan Stok Kasir'
})
const stockForm = ref(defaultStockForm())

const openAddStockModal = () => {
  stockForm.value = defaultStockForm()
  modalSearchQuery.value = ''
  selectedProductLabel.value = ''
  isDropdownOpen.value = false
  showAddModal.value = true
}

const loadData = async () => {
  isLoading.value = true
  try {
    const [prodsRes, catsRes, brandsRes, unitsRes, lowRes, movRes, sessionRes] = await Promise.all([
      inventoryApi.getProducts(),
      inventoryApi.getCategories(),
      inventoryApi.getBrands(),
      inventoryApi.getUnits(),
      inventoryApi.getLowStockProducts(),
      inventoryApi.getStockMovements(),
      cashApi.getCurrentSession().catch(() => null),
    ])
    products.value = prodsRes
    categories.value = catsRes
    brands.value = brandsRes
    units.value = unitsRes
    lowStockProducts.value = lowRes
    movements.value = movRes.slice(0, 5)
    currentSession.value = sessionRes

    setTimeout(async () => {
      for (const p of products.value) {
        try {
          const mediaList = await inventoryApi.getProductMedia(p.id)
          if (mediaList && mediaList.length > 0) {
            p.image = inventoryApi.getMediaUrl(mediaList[0].id)
          }
        } catch (e) {
          console.error(`Gagal mengambil gambar untuk ${p.id}`, e)
        }
      }
    }, 0)
  } catch (error: any) {
    console.error('Gagal memuat data:', error)
  } finally {
    isLoading.value = false
  }
}

const searchByBarcode = async () => {
  if (!barcodeQuery.value) return
  isLoading.value = true
  try {
    const product = await inventoryApi.getProductByBarcode(barcodeQuery.value)
    if (product) {
      products.value = [product]
      searchQuery.value = ''
      first.value = 0
    } else {
      products.value = []
    }
  } catch (error) {
    console.error(error)
    products.value = []
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadData()
  updateClock()
  timerInterval = window.setInterval(updateClock, 1000)
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

const submitForm = async () => {
  if (stockForm.value.new_qty < 1) {
    alert('Kuantitas penambahan stok minimal 1');
    return;
  }

  isSubmitting.value = true
  try {
    const selectedProduct = products.value.find(p => p.id === stockForm.value.product_id)
    const currentQty = selectedProduct ? (selectedProduct.qty || selectedProduct.stock || 0) : 0
    const finalQty = currentQty + stockForm.value.new_qty

    await inventoryApi.adjustStock({
      product_id: stockForm.value.product_id,
      new_qty: finalQty,
      reference: stockForm.value.reference
    })
    
    stockForm.value = defaultStockForm()
    showAddModal.value = false
    await loadData()
  } catch (error: any) {
    console.error('Gagal menyesuaikan stok:', error)
  } finally {
    isSubmitting.value = false
  }
}
</script>