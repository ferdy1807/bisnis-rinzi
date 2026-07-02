<template>
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div class="flex items-center gap-4">
        <router-link
          to="/products"
          class="p-2 bg-surface hover:bg-surface-variant rounded-full text-on-surface-variant transition-colors border border-outline-variant/30"
        >
          <span class="material-symbols-outlined">arrow_back</span>
        </router-link>
        <div>
          <h2 class="font-headline-lg text-headline-lg text-primary">Detail Produk</h2>
          <p class="text-body-md text-on-surface-variant">
            Informasi lengkap terkait SKU dan master data produk.
          </p>
        </div>
      </div>
      <div v-if="product" class="flex items-center gap-2">
        <button
          @click="showAddStockModal = true"
          class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2 rounded-xl shadow-md hover:bg-primary/90 transition-colors"
        >
          <span class="material-symbols-outlined">add_box</span>
          <span class="text-label-md font-bold">Tambah Stok</span>
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="flex justify-center py-20 text-primary">
      <span class="animate-spin material-symbols-outlined text-5xl">sync</span>
    </div>

    <div
      v-else-if="!product"
      class="bg-surface rounded-2xl border border-outline-variant/30 p-12 text-center shadow-sm">
      <span class="material-symbols-outlined text-6xl text-on-surface-variant mb-4"
        >inventory_2</span
      >
      <h3 class="text-headline-md font-bold mb-2">Produk Tidak Ditemukan</h3>
      <p class="text-on-surface-variant mb-6">
        Data produk dengan ID tersebut tidak tersedia di sistem.
      </p>
      <router-link
        to="/products"
        class="bg-primary text-on-primary px-6 py-2.5 rounded-xl font-medium inline-block"
        >Kembali ke Katalog</router-link
      >
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Kolom Kiri: Informasi Utama -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Identitas Produk -->
        <div
          class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
        >
          <div class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest">
            <h3 class="font-headline-md text-primary flex items-center gap-2">
              <span class="material-symbols-outlined">badge</span> Identitas Produk
            </h3>
          </div>
          <div class="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <p class="text-label-sm text-on-surface-variant mb-1">Nama Produk</p>
              <p class="text-body-lg font-bold">{{ product.name }}</p>
            </div>
            <div>
              <p class="text-label-sm text-on-surface-variant mb-1">Kode SKU</p>
              <p class="text-body-lg font-mono">{{ product.sku || '-' }}</p>
            </div>
            <div>
              <p class="text-label-sm text-on-surface-variant mb-1">Barcode Utama</p>
              <div class="flex items-center gap-2">
                <span class="material-symbols-outlined text-on-surface-variant text-sm"
                  >barcode</span
                >
                <p class="text-body-md font-mono">{{ product.barcode || '-' }}</p>
              </div>
            </div>
            <div>
              <p class="text-label-sm text-on-surface-variant mb-1">Kategori ID</p>
              <p class="text-body-md font-mono">{{ product.category_id || '-' }}</p>
            </div>
            <div>
              <p class="text-label-sm text-on-surface-variant mb-1">Brand ID</p>
              <p class="text-body-md font-mono">{{ product.brand_id || '-' }}</p>
            </div>
          </div>
        </div>

        <!-- Harga -->
        <div
          class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
        >
          <div class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest">
            <h3 class="font-headline-md text-primary flex items-center gap-2">
              <span class="material-symbols-outlined">payments</span> Informasi Harga
            </h3>
          </div>
          <div class="p-6">
            <div class="bg-primary/5 p-4 rounded-xl border border-primary/20 text-center">
              <p class="text-label-sm text-primary font-bold mb-2">Harga Jual (Selling Price)</p>
              <p class="text-title-lg font-bold text-primary">
                Rp {{ product.selling_price?.toLocaleString('id-ID') || 0 }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Kolom Kanan: Stok & Status -->
      <div class="space-y-6">
        <!-- Kartu Stok Aktual -->
        <div
          class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden text-center relative"
        >
          <div
            :class="[
              'h-2 w-full absolute top-0 left-0',
              (product.qty || 0) <= 0 ? 'bg-error' : 'bg-primary',
            ]"
          ></div>
          <div class="p-8">
            <p
              class="text-label-md text-on-surface-variant uppercase tracking-widest font-bold mb-2"
            >
              Kuantitas Aktual
            </p>
            <h2
              :class="[
                'text-display-lg font-bold leading-none mb-2',
                (product.qty || 0) <= 0 ? 'text-error' : 'text-primary',
              ]"
            >
              {{ product.qty || 0 }}
            </h2>
            <p class="text-title-md text-on-surface-variant">
              {{ product.base_unit_code || 'Unit' }}
            </p>
          </div>
          <div
            v-if="(product.qty || 0) <= 0"
            class="bg-error/10 py-3 px-4 text-error text-label-md font-bold flex justify-center items-center gap-2"
          >
            <span class="material-symbols-outlined text-sm">warning</span> Stok Habis!
          </div>
        </div>

        <!-- Detail Meta -->
        <div
          class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
        >
          <div class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest">
            <h3 class="font-headline-md text-primary flex items-center gap-2">
              <span class="material-symbols-outlined">info</span> Informasi Tambahan
            </h3>
          </div>
          <div class="p-4 space-y-4">
            <div class="flex justify-between items-center pb-1">
              <span class="text-body-sm text-on-surface-variant">Status Aktif</span>
              <span
                :class="[
                  'px-2 py-1 rounded text-xs font-bold',
                  product.is_active ? 'bg-primary/20 text-primary' : 'bg-error/20 text-error',
                ]"
              >
                {{ product.is_active ? 'AKTIF' : 'NONAKTIF' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Galeri Media Produk -->
        <div
          class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
        >
          <div
            class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex justify-between items-center"
          >
            <h3 class="font-headline-md text-primary flex items-center gap-2">
              <span class="material-symbols-outlined">imagesmode</span> Media Produk
            </h3>
            <button
              @click="triggerUpload"
              :disabled="isUploading"
              class="text-label-sm font-bold bg-primary/10 text-primary px-3 py-1.5 rounded-lg hover:bg-primary/20 transition-colors flex items-center gap-1 disabled:opacity-50"
            >
              <span v-if="isUploading" class="animate-spin material-symbols-outlined text-sm"
                >sync</span
              >
              <span v-else class="material-symbols-outlined text-sm">upload</span>
              {{ mediaList.length > 0 ? 'Ubah Foto' : 'Tambah Foto' }}
            </button>
            <input
              type="file"
              ref="fileInput"
              class="hidden"
              accept="image/*"
              @change="handleFileUpload"
            />
          </div>
          <div class="p-4">
            <div
              v-if="mediaList.length === 0"
              class="text-center py-6 text-on-surface-variant text-label-sm border-2 border-dashed border-outline-variant/30 rounded-xl"
            >
              Belum ada foto/media untuk produk ini.
            </div>
            <div v-else class="grid grid-cols-2 gap-4">
              <div
                v-for="media in mediaList"
                :key="media.id"
                class="aspect-square rounded-xl border border-outline-variant/30 overflow-hidden bg-surface-container-lowest flex items-center justify-center relative group"
              >
                <img
                  :src="inventoryApi.getMediaUrl(media.id)"
                  alt="Product Media"
                  class="w-full h-full object-cover"
                />
                <div
                  class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                >
                  <span
                    class="text-white text-xs font-bold px-2 py-1 bg-black/50 rounded-lg backdrop-blur-sm"
                    >{{ media.original_file_name || 'Gambar' }}</span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Riwayat Mutasi Stok -->
    <div
      v-if="product"
      class="mt-6 bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
    >
      <div class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest">
        <h3 class="font-headline-md text-primary flex items-center gap-2">
          <span class="material-symbols-outlined">history</span> Riwayat Mutasi Stok
        </h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr
              class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
            >
              <th class="px-6 py-4">Waktu</th>
              <th class="px-6 py-4">Tipe Mutasi</th>
              <th class="px-6 py-4">Referensi</th>
              <th class="px-6 py-4 text-right">Kuantitas</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/20">
            <tr v-if="movements.length === 0">
              <td colspan="4" class="px-6 py-8 text-center text-on-surface-variant text-label-sm">
                Belum ada riwayat pergerakan stok.
              </td>
            </tr>
            <tr
              v-for="m in movements"
              :key="m.id"
              class="hover:bg-surface-container-lowest/50 transition-colors"
            >
              <td class="px-6 py-4 text-sm">{{ formatDate(m.created_at) }}</td>
              <td class="px-6 py-4">
                <span
                  :class="[
                    'px-2 py-1 rounded text-xs font-bold',
                    m.movement_type === 'IN'
                      ? 'bg-primary/10 text-primary'
                      : m.movement_type === 'OUT'
                        ? 'bg-error/10 text-error'
                        : 'bg-tertiary/10 text-tertiary',
                  ]"
                >
                  {{ m.movement_type }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">{{ m.reference || '-' }}</td>
              <td
                :class="[
                  'px-6 py-4 text-right font-bold text-sm',
                  m.movement_type === 'IN'
                    ? 'text-primary'
                    : m.movement_type === 'OUT'
                      ? 'text-error'
                      : 'text-on-surface',
                ]"
              >
                {{ m.movement_type === 'IN' ? '+' : m.movement_type === 'OUT' ? '-' : ''
                }}{{ m.qty }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Tambah Stok -->
    <div
      v-if="showAddStockModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    >
      <div class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl transform transition-all">
        <h3 class="text-title-lg font-bold text-primary mb-4 flex items-center gap-2">
          <span class="material-symbols-outlined">add_box</span> Tambah Stok
        </h3>

        <form @submit.prevent="submitAddStock" class="space-y-4">
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1"
              >Jumlah Ditambahkan</label
            >
            <div class="relative">
              <input
                v-model.number="stockForm.addedQty"
                type="number"
                min="1"
                required
                class="w-full pl-4 pr-12 py-2 rounded-lg border border-outline-variant focus:border-primary focus:ring-1 focus:ring-primary outline-none bg-surface-container-lowest"
              />
              <span
                class="absolute right-4 top-1/2 -translate-y-1/2 text-on-surface-variant text-sm"
                >{{ product?.base_unit_code }}</span
              >
            </div>
          </div>
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1"
              >Keterangan / Referensi</label
            >
            <input
              v-model="stockForm.reference"
              type="text"
              required
              class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-primary focus:ring-1 focus:ring-primary outline-none bg-surface-container-lowest"
              placeholder="Contoh: Pembelian dari Supplier X"
            />
          </div>

          <div class="flex justify-end gap-3 mt-8">
            <button
              type="button"
              @click="showAddStockModal = false"
              class="px-4 py-2 text-label-md font-medium text-on-surface-variant hover:bg-surface-variant/20 rounded-lg transition-colors"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSubmittingStock"
              class="px-4 py-2 text-label-md font-medium bg-primary text-on-primary rounded-lg shadow hover:bg-primary/90 transition-colors disabled:opacity-70 flex items-center gap-2"
            >
              <span v-if="isSubmittingStock" class="animate-spin material-symbols-outlined text-sm"
                >sync</span
              >
              {{ isSubmittingStock ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { inventoryApi } from '@frontend/api/inventory'
import type { Product, ProductMedia, StockMovement } from '@frontend/types/inventory'

const route = useRoute()
const isLoading = ref(true)
const product = ref<Product | null>(null)
const mediaList = ref<ProductMedia[]>([])
const movements = ref<StockMovement[]>([])

const showAddStockModal = ref(false)
const isSubmittingStock = ref(false)
const stockForm = ref({
  addedQty: 1,
  reference: '',
})

// --- STATE UPLOAD PHOTO ---
const fileInput = ref<HTMLInputElement | null>(null)
const isUploading = ref(false)

const triggerUpload = () => {
  fileInput.value?.click()
}

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0 || !product.value) return

  const file = target.files[0]
  isUploading.value = true

  try {
    // 1. Delete old media if exists
    if (mediaList.value.length > 0) {
      for (const media of mediaList.value) {
        await inventoryApi.deleteProductMedia(product.value.id, media.id)
      }
    }

    // 2. Rename file to product name
    const ext = file.name.split('.').pop() || 'jpg'
    const newFileName = `${product.value.name.replace(/[^a-zA-Z0-9]/g, '_')}.${ext}`
    const renamedFile = new File([file], newFileName, { type: file.type })

    // 3. Upload new media
    const formData = new FormData()
    formData.append('media', renamedFile)
    await inventoryApi.uploadProductMedia(product.value.id, formData)

    // 4. Reload media list
    mediaList.value = await inventoryApi.getProductMedia(product.value.id)
  } catch (e) {
    console.error('Gagal mengunggah foto:', e)
    alert('Terjadi kesalahan saat mengunggah foto produk.')
  } finally {
    isUploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

const loadProduct = async () => {
  isLoading.value = true
  try {
    const id = route.params.id as string
    if (id) {
      product.value = await inventoryApi.getProductById(id)

      try {
        const allMovements = await inventoryApi.getStockMovements()
        movements.value = allMovements
          .filter((m) => m.product_id === id)
          .sort(
            (a, b) => new Date(b.created_at || 0).getTime() - new Date(a.created_at || 0).getTime(),
          )
      } catch (e) {
        console.error('Gagal mengambil riwayat stok:', e)
      }

      try {
        mediaList.value = await inventoryApi.getProductMedia(id)
      } catch (e) {
        console.error('Gagal mengambil media produk:', e)
      }
    }
  } catch (error) {
    console.error('Gagal mengambil data produk:', error)
  } finally {
    isLoading.value = false
  }
}

const submitAddStock = async () => {
  if (!product.value) return
  isSubmittingStock.value = true
  try {
    const newQty = (product.value.qty || product.value.stock || 0) + stockForm.value.addedQty
    await inventoryApi.adjustStock({
      product_id: product.value.id,
      new_qty: newQty,
      reference: stockForm.value.reference,
    })

    // Reset form and modal
    showAddStockModal.value = false
    stockForm.value = { addedQty: 1, reference: '' }

    // Reload product and movements
    await loadProduct()
  } catch (e) {
    console.error('Gagal menambah stok:', e)
    alert('Terjadi kesalahan saat menambah stok.')
  } finally {
    isSubmittingStock.value = false
  }
}

const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  const d = new Date(dateString)
  return d.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  loadProduct()
})
</script>
