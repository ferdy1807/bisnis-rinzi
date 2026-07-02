<template>
   <div class="flex flex-col h-[calc(100vh-8rem)]">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-4">
        <div>
          <h2 class="font-headline-lg text-headline-lg text-primary">Terminal Kasir</h2>
          <p class="text-body-sm text-on-surface-variant">
            Pilih produk atau scan barcode untuk menambahkan ke keranjang.
          </p>
        </div>
        <div class="flex gap-2 w-full md:w-auto">
          <div class="relative w-full md:w-64">
            <span
              class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant"
              >barcode_scanner</span>
            <input
              v-model="searchQuery"
              class="w-full bg-surface-container-highest border-none rounded-xl pl-10 pr-4 py-2.5 text-body-md focus:ring-2 focus:ring-primary focus:outline-none transition-all"
              placeholder="Scan atau cari nama..."
              type="text"
            />
          </div>
        </div>
      </div>

      <div class="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-6 overflow-hidden">
        <div
          class="lg:col-span-8 flex flex-col bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
        >
          <div
            class="p-4 border-b border-outline-variant/20 flex gap-2 overflow-x-auto hide-scrollbar"
          >
            <button
              @click="selectedCategory = null"
              :class="[
                'px-4 py-1.5 rounded-lg text-label-sm font-bold whitespace-nowrap shadow-sm transition-colors',
                selectedCategory === null
                  ? 'bg-primary text-on-primary'
                  : 'bg-surface-container hover:bg-surface-container-high text-on-surface-variant',
              ]"
            >
              Semua
            </button>
            <button
              v-for="cat in categories"
              :key="cat.id"
              @click="selectedCategory = cat.id"
              :class="[
                'px-4 py-1.5 rounded-lg text-label-sm font-bold whitespace-nowrap transition-colors',
                selectedCategory === cat.id
                  ? 'bg-primary text-on-primary shadow-sm'
                  : 'bg-surface-container hover:bg-surface-container-high text-on-surface-variant',
              ]"
            >
              {{ cat.name }}
            </button>
          </div>

          <div class="flex-1 p-4 overflow-y-auto custom-scrollbar">
            <div
              v-if="isLoadingProducts"
              class="flex justify-center items-center h-full text-outline-variant"
            >
              <span class="animate-spin material-symbols-outlined text-4xl">sync</span>
            </div>
            <div
              v-else-if="filteredProducts.length === 0"
              class="flex justify-center items-center h-full text-on-surface-variant"
            >
              Produk tidak ditemukan.
            </div>
            <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4">
              <div
                v-for="product in filteredProducts"
                :key="product.id"
                @click="cartStore.addItem(product)"
                class="bg-white border border-outline-variant/30 rounded-xl p-3 cursor-pointer hover:border-primary hover:shadow-md transition-all active:scale-95 flex flex-col h-full"
              >
                <div
                  class="aspect-square bg-surface-container-low rounded-lg mb-3 flex items-center justify-center overflow-hidden relative"
                >
                  <span
                    v-if="!product.image"
                    class="material-symbols-outlined text-4xl text-outline-variant"
                    >image</span
                  >
                  <img v-else :src="product.image" class="w-full h-full object-cover" />
                  <div
                    class="absolute top-1 right-1 bg-primary text-on-primary text-[10px] px-1.5 py-0.5 rounded font-bold"
                  >
                    Sisa {{ product.qty ?? 0 }}
                  </div>
                </div>
                <h4 class="font-label-sm text-on-surface leading-tight mb-1 flex-1">
                  {{ product.name }}
                </h4>
                <p class="text-primary font-bold text-label-md">
                  Rp {{ product.selling_price.toLocaleString('id-ID') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div
          class="lg:col-span-4 flex flex-col bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
        >
          <div
            class="p-4 border-b border-outline-variant/20 bg-primary/5 flex justify-between items-center"
          >
            <h3 class="font-headline-md text-primary flex items-center gap-2">
              <span class="material-symbols-outlined">shopping_cart</span> Keranjang
            </h3>
            <span class="bg-primary text-on-primary text-xs px-2 py-1 rounded-full font-bold"
              >{{ cartStore.items.length }} Item</span
            >
          </div>

          <div class="flex-1 p-2 overflow-y-auto custom-scrollbar">
            <div
              v-if="cartStore.items.length === 0"
              class="h-full flex flex-col items-center justify-center text-outline"
            >
              <span class="material-symbols-outlined text-5xl mb-2 opacity-50"
                >remove_shopping_cart</span
              >
              <p class="text-label-sm">Keranjang masih kosong</p>
            </div>

            <div v-else class="space-y-2">
              <div
                v-for="item in cartStore.items"
                :key="item.product_id"
                class="bg-white p-3 border border-outline-variant/30 rounded-xl flex flex-col gap-2"
              >
                <div class="flex justify-between items-start">
                  <p class="font-label-sm text-on-surface truncate pr-2">{{ item.product_name }}</p>
                  <button
                    @click="cartStore.removeItem(item.product_id)"
                    class="text-error hover:bg-error/10 p-1 rounded transition-colors shrink-0"
                  >
                    <span class="material-symbols-outlined text-[16px]">delete</span>
                  </button>
                </div>
                <div class="flex justify-between items-center">
                  <p class="text-primary font-bold text-label-sm">
                    Rp {{ item.subtotal.toLocaleString('id-ID') }}
                  </p>

                  <div
                    class="flex items-center bg-surface-container rounded-lg border border-outline-variant/50 overflow-hidden"
                  >
                    <button
                      @click="cartStore.updateItemQty(item.product_id, item.qty - 1)"
                      class="w-7 h-7 flex items-center justify-center hover:bg-outline-variant/30 text-on-surface-variant transition-colors"
                    >
                      <span class="material-symbols-outlined text-[16px]">remove</span>
                    </button>
                    <input
                      :value="item.qty"
                      readonly
                      class="w-8 h-7 text-center bg-transparent border-none text-label-sm font-bold focus:ring-0 p-0"
                    />
                    <button
                      @click="cartStore.updateItemQty(item.product_id, item.qty + 1)"
                      class="w-7 h-7 flex items-center justify-center hover:bg-outline-variant/30 text-on-surface-variant transition-colors"
                    >
                      <span class="material-symbols-outlined text-[16px]">add</span>
                    </button>
                  </div>
                </div>
                <div
                  class="flex justify-between items-center mt-1 pt-2 border-t border-outline-variant/20"
                >
                  <span class="text-[11px] text-on-surface-variant font-medium">Diskon Item:</span>
                  <div class="relative w-24">
                    <span class="absolute left-1.5 top-1/2 -translate-y-1/2 text-error text-[10px]"
                      >-Rp</span
                    >
                    <input
                      type="number"
                      :value="item.discount"
                      @change="
                        (e) =>
                          cartStore.updateItemDiscount(
                            item.product_id,
                            Number((e.target as HTMLInputElement).value),
                          )
                      "
                      class="w-full pl-7 pr-1.5 py-0.5 text-right text-error bg-surface-container-lowest rounded border border-outline-variant/50 focus:border-primary focus:ring-1 focus:ring-primary outline-none text-[11px]"
                      min="0"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="p-4 bg-surface-container-lowest border-t border-outline-variant/30">
            <div class="space-y-1 mb-4">
              <div class="flex justify-between text-label-sm text-on-surface-variant">
                <span>Subtotal</span>
                <span>Rp {{ cartStore.subtotal.toLocaleString('id-ID') }}</span>
              </div>
              <!-- <div class="flex justify-between items-center text-label-sm text-on-surface-variant">
                <span>Diskon (Rp)</span>
                <div class="relative">
                  <span class="absolute left-2 top-1/2 -translate-y-1/2 text-error">- Rp</span>
                  <input 
                    type="number" 
                    v-model.number="cartStore.discountAmount" 
                    class="w-32 pl-10 pr-2 py-1 text-right text-error bg-surface-container-lowest rounded-lg border border-outline-variant/50 focus:border-primary focus:ring-1 focus:ring-primary outline-none"
                    min="0"
                  />
                </div>
              </div>
              <div class="flex justify-between text-label-sm text-on-surface-variant">
                <span>Pajak (PPN)</span>
                <span>Rp {{ cartStore.tax.toLocaleString('id-ID') }}</span>
              </div> -->
              <div
                class="flex justify-between font-headline-md text-primary mt-2 pt-2 border-t border-outline-variant/30"
              >
                <span>Total</span>
                <span>Rp {{ cartStore.grandTotal.toLocaleString('id-ID') }}</span>
              </div>
            </div>

            <div class="flex gap-2">
              <button
                @click="cartStore.clearCart()"
                class="p-3 border border-error text-error rounded-xl hover:bg-error/10 transition-colors flex items-center justify-center"
                title="Batalkan Transaksi"
              >
                <span class="material-symbols-outlined">close</span> Reset
              </button>
              <button
                @click="showCheckoutModal = true"
                class="flex-1 bg-primary text-on-primary py-3 rounded-xl font-label-md shadow-lg hover:bg-primary-container hover:-translate-y-1 active:scale-95 transition-all flex items-center justify-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed"
                :disabled="cartStore.items.length === 0 || cartStore.isCheckoutLoading"
              >
                <span class="material-symbols-outlined">payments</span> Bayar Sekarang
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Checkout -->
    <div
      v-if="showCheckoutModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div class="bg-surface rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div
          class="p-4 border-b border-outline-variant/30 flex justify-between items-center bg-primary text-on-primary"
        >
          <h3 class="font-headline-sm">Pembayaran</h3>
          <button
            @click="showCheckoutModal = false"
            class="hover:bg-white/20 rounded-full p-1 transition-colors"
          >
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div class="text-center">
            <p class="text-label-md text-on-surface-variant mb-1">Total Tagihan</p>
            <p class="font-headline-lg text-primary">
              Rp {{ cartStore.grandTotal.toLocaleString('id-ID') }}
            </p>
          </div>

          <div>
            <label class="block text-label-sm font-bold text-on-surface mb-2"
              >Metode Pembayaran</label
            >
            <select
              v-model="paymentMethod"
              class="w-full bg-surface-container border border-outline-variant/50 rounded-xl px-4 py-3 text-body-md focus:ring-2 focus:ring-primary focus:outline-none"
            >
              <option value="CASH">Tunai (Cash)</option>
              <option value="QRIS">QRIS</option>
              <option value="TRANSFER">Transfer Bank</option>
            </select>
          </div>

          <div v-if="paymentMethod === 'CASH'">
            <label class="block text-label-sm font-bold text-on-surface mb-2"
              >Nominal Uang Diterima</label
            >
            <div class="relative mb-2">
              <span
                class="absolute left-4 top-1/2 -translate-y-1/2 text-on-surface-variant font-bold"
                >Rp</span
              >
              <input
                type="number"
                v-model.number="amountPaid"
                class="w-full bg-surface-container border border-outline-variant/50 rounded-xl pl-12 pr-4 py-3 text-headline-sm focus:ring-2 focus:ring-primary focus:outline-none"
                placeholder="0"
                min="0"
              />
            </div>
            <div class="flex gap-2 flex-wrap">
              <button
                @click="setAmountPaid(cartStore.grandTotal)"
                class="px-3 py-1.5 bg-surface-container-high hover:bg-primary hover:text-on-primary rounded-lg text-label-sm border border-outline-variant/30 transition-colors"
              >
                Uang Pas
              </button>
              <button
                @click="setAmountPaid(50000)"
                class="px-3 py-1.5 bg-surface-container-high hover:bg-primary hover:text-on-primary rounded-lg text-label-sm border border-outline-variant/30 transition-colors"
              >
                50.000
              </button>
              <button
                @click="setAmountPaid(100000)"
                class="px-3 py-1.5 bg-surface-container-high hover:bg-primary hover:text-on-primary rounded-lg text-label-sm border border-outline-variant/30 transition-colors"
              >
                100.000
              </button>
            </div>

            <div
              class="mt-4 p-3 bg-surface-container-lowest border border-outline-variant/30 rounded-xl flex justify-between items-center"
            >
              <span class="text-label-md text-on-surface-variant">Kembalian</span>
              <span
                :class="['font-headline-sm', changeAmount >= 0 ? 'text-success' : 'text-error']"
              >
                Rp {{ changeAmount >= 0 ? changeAmount.toLocaleString('id-ID') : '0' }}
              </span>
            </div>
          </div>
        </div>
        <div class="p-4 border-t border-outline-variant/30 bg-surface-container-lowest">
          <button
            @click="processCheckout"
            :disabled="
              cartStore.isCheckoutLoading ||
              (paymentMethod === 'CASH' && (amountPaid || 0) < cartStore.grandTotal)
            "
            class="w-full bg-primary text-on-primary py-3 rounded-xl font-label-md shadow-md hover:bg-primary-container hover:-translate-y-0.5 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <template v-if="cartStore.isCheckoutLoading">
              <span class="animate-spin material-symbols-outlined">sync</span> Memproses...
            </template>
            <template v-else>
              <span class="material-symbols-outlined">check_circle</span> Konfirmasi Pembayaran
            </template>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Modal Preview Struk -->
    <div
      v-if="showReceiptModal && completedSale"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
    >
      <div class="bg-surface rounded-2xl shadow-xl w-full max-w-sm overflow-hidden flex flex-col max-h-[90vh]">
        <div class="p-4 border-b border-outline-variant/30 flex justify-between items-center bg-surface-container">
          <h3 class="font-bold text-on-surface">Preview Struk</h3>
          <button @click="closeReceiptModal" class="hover:bg-outline-variant/20 rounded-full p-1 transition-colors">
            <span class="material-symbols-outlined text-on-surface-variant">close</span>
          </button>
        </div>
        
        <!-- Receipt Content -->
        <div class="p-6 bg-gray-100 flex justify-center overflow-y-auto custom-scrollbar flex-1">
          <div id="receipt-print-area" class="bg-white p-4 text-black w-[80mm] min-h-[100mm] shadow-sm text-[12px] font-mono leading-tight">
            <!-- Header -->
            <div class="text-center mb-4 border-b border-dashed border-gray-400 pb-2">
              <h2 class="font-bold text-[16px] uppercase tracking-wider">Bisnis Rinzi</h2>
              <p>Toko Rinzi</p>
              <p>Jl. KH M Noor, Pedamaran VI</p>
            </div>
            
            <!-- Meta -->
            <div class="mb-2 space-y-0.5">
              <div class="flex justify-between"><span>No:</span><span>{{ completedSale.invoice_number }}</span></div>
              <div class="flex justify-between"><span>Tgl:</span><span>{{ new Date(completedSale.transaction_date).toLocaleString('id-ID') }}</span></div>
              <div class="flex justify-between"><span>Ksr:</span><span>{{ cashierName }}</span></div>
            </div>
            
            <div class="border-t border-dashed border-gray-400 my-2"></div>
            
            <!-- Items -->
            <div class="space-y-2 mb-2">
              <div v-for="item in completedItems" :key="item.product_id">
                <div class="font-bold whitespace-normal">{{ item.product_name }}</div>
                <div class="flex justify-between">
                  <span>{{ item.qty }} x {{ item.unit_price.toLocaleString('id-ID') }}</span>
                  <span>{{ (item.qty * item.unit_price).toLocaleString('id-ID') }}</span>
                </div>
                <div v-if="item.discount > 0" class="flex justify-between text-gray-600">
                  <span>Diskon Item</span>
                  <span>-{{ item.discount.toLocaleString('id-ID') }}</span>
                </div>
              </div>
            </div>
            
            <div class="border-t border-dashed border-gray-400 my-2"></div>
            
            <!-- Totals -->
            <div class="space-y-1 font-bold">
              <div class="flex justify-between">
                <span>Subtotal</span>
                <span>{{ completedSale.subtotal.toLocaleString('id-ID') }}</span>
              </div>
              <div v-if="completedSale.discount > 0" class="flex justify-between font-normal">
                <span>Diskon Transaksi</span>
                <span>-{{ completedSale.discount.toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between text-[14px] mt-1 pt-1 border-t border-gray-300">
                <span>TOTAL</span>
                <span>{{ completedSale.total.toLocaleString('id-ID') }}</span>
              </div>
            </div>
            
            <div class="border-t border-dashed border-gray-400 my-2"></div>
            
            <!-- Payment -->
            <div class="space-y-1">
              <div class="flex justify-between">
                <span>Bayar ({{ completedSale.payment_method }})</span>
                <span>{{ (completedSale.amount_paid || completedSale.total).toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between">
                <span>Kembali</span>
                <span>{{ (completedSale.change_amount || 0).toLocaleString('id-ID') }}</span>
              </div>
            </div>
            
            <div class="border-t border-dashed border-gray-400 my-2"></div>
            
            <!-- Footer -->
            <div class="text-center mt-4 space-y-1">
              <p class="font-bold">TERIMA KASIH</p>
              <p class="text-[10px]">Barang yang dibeli tidak dapat ditukar/dikembalikan.</p>
            </div>
          </div>
        </div>
        
        <!-- Actions -->
        <div class="p-4 border-t border-outline-variant/30 flex gap-3 bg-surface">
          <button @click="closeReceiptModal" class="flex-1 px-4 py-2 border border-outline-variant rounded-xl text-on-surface-variant font-bold hover:bg-surface-container transition-colors">
            Tutup
          </button>
          <button @click="printReceipt" class="flex-1 px-4 py-2 bg-primary text-on-primary rounded-xl font-bold hover:bg-primary-container hover:text-on-primary-container shadow-md flex items-center justify-center gap-2 transition-colors">
            <span class="material-symbols-outlined text-sm">print</span> Cetak Struk
          </button>
        </div>
      </div>
    </div>
</template>

<style>
@media print {
  @page {
    margin: 0;
    size: 80mm auto;
  }
  
  * {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
  }
  
  body * {
    visibility: hidden;
  }
  
  #receipt-print-area, #receipt-print-area * {
    visibility: visible;
  }
  
  #receipt-print-area {
    position: absolute;
    left: 0;
    top: 0;
    width: 80mm !important;
    margin: 0 !important;
    padding: 4mm !important;
    box-shadow: none !important;
    background-color: white !important;
    color: black !important;
  }
}
</style>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'  
import { useCartStore } from '@frontend/stores/chart'
import type { Product } from '@frontend/types/inventory'
import { inventoryApi } from '@frontend/api/inventory'
import { posApi } from '@frontend/api/pos'
import { useAuthStore } from '@frontend/stores/auth'
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'

const authStore = useAuthStore()
const cartStore = useCartStore()
const searchQuery = ref('')
const selectedCategory = ref<string | null>(null)

const showCheckoutModal = ref(false)
const paymentMethod = ref('CASH')
const amountPaid = ref<number | null>(null)

const changeAmount = computed(() => {
  return (amountPaid.value || 0) - cartStore.grandTotal
})

const setAmountPaid = (amount: number) => {
  amountPaid.value = amount
}

const cashierName = computed(() => {
  return authStore.user?.full_name || authStore.user?.username || 'Kasir'
})

import type { Sale, SaleItem } from '@frontend/types/pos'
import type { Category } from '@frontend/types/inventory'

const currentSale = ref<Sale | null>(null)
const currentItems = ref<SaleItem[]>([])

// Data Produk dari API
const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const isLoadingProducts = ref<boolean>(false)

onMounted(async () => {
  isLoadingProducts.value = true
  try {
    const [fetchedProducts, fetchedCategories] = await Promise.all([
      inventoryApi.getProducts(),
      inventoryApi.getCategories(),
    ])

    if (fetchedProducts && fetchedProducts.length > 0) {
      products.value = fetchedProducts

      setTimeout(async () => {
        for (const p of products.value) {
          try {
            const mediaList = await inventoryApi.getProductMedia(p.id)
            if (mediaList && mediaList.length > 0) {
              p.image = inventoryApi.getMediaUrl(mediaList[0].id)
            }
          } catch (e) {
            console.error(`Failed to fetch media for product ${p.id}`, e)
          }
        }
      }, 0)
    }
    if (fetchedCategories && fetchedCategories.length > 0) {
      categories.value = fetchedCategories
    }
  } catch (error) {
    console.error('Gagal mengambil data dari API:', error)
  } finally {
    isLoadingProducts.value = false
  }
})

// Fungsi Pencarian Produk
const filteredProducts = computed(() => {
  let result = products.value

  if (selectedCategory.value) {
    result = result.filter((p) => p.category_id === selectedCategory.value)
  }

  if (searchQuery.value) {
    result = result.filter((p) => p.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
  }

  return result
})

// Fungsi Generate PDF Natively
const generateAndUploadInvoice = async (sale: Sale, items: SaleItem[], cashier: string) => {
  try {
    const doc = new jsPDF('p', 'mm', 'a4')
    const pageWidth = doc.internal.pageSize.getWidth()

    // Konfigurasi Header
    doc.setFillColor(15, 23, 42) // slate-900
    doc.rect(0, 0, pageWidth, 40, 'F')

    // Teks Header Kiri
    doc.setTextColor(255, 255, 255)
    doc.setFontSize(24)
    doc.setFont('helvetica', 'bold')
    doc.text('INVOICE', 15, 18)

    doc.setFontSize(12)
    doc.setFont('helvetica', 'normal')
    doc.text('Toko Bisnis Rinzi', 15, 28)

    // Teks Header Kanan
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text(`No. ${sale.invoice_number}`, pageWidth - 15, 16, { align: 'right' })

    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    const d = new Date(sale.transaction_date || new Date())
    const dateStr = d.toLocaleString('id-ID', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
    doc.text(dateStr, pageWidth - 15, 24, { align: 'right' })

    doc.setFont('helvetica', 'bold')
    doc.text(`Kasir: ${cashier}`, pageWidth - 15, 32, { align: 'right' })

    // Persiapkan data tabel
    const tableBody = items.map((item) => [
      item.product_name,
      `${item.qty} ${item.unit_code}`,
      `Rp ${item.unit_price.toLocaleString('id-ID')}`,
      item.discount > 0 ? `-Rp ${item.discount.toLocaleString('id-ID')}` : '-',
      `Rp ${item.subtotal.toLocaleString('id-ID')}`,
    ])

    // Render Tabel
    autoTable(doc, {
      startY: 50,
      head: [['Item', 'Qty', 'Harga Satuan', 'Diskon', 'Subtotal']],
      body: tableBody,
      theme: 'grid',
      headStyles: { fillColor: [248, 250, 252], textColor: [100, 116, 139], fontStyle: 'bold' },
      styles: { font: 'helvetica', fontSize: 10, cellPadding: 4, textColor: [30, 41, 59] },
      columnStyles: {
        0: { cellWidth: 60 },
        1: { halign: 'center' },
        2: { halign: 'right' },
        3: { halign: 'right', textColor: [220, 38, 38] }, // merah untuk diskon
        4: { halign: 'right', fontStyle: 'bold' },
      },
    })

    // Summary Box
    // @ts-ignore
    const finalY = (doc as any).lastAutoTable?.finalY || 50
    const summaryX = pageWidth - 80

    doc.setFillColor(248, 250, 252)
    // Increase box height to fit change amount
    doc.roundedRect(summaryX, finalY + 10, 65, sale.payment_method === 'CASH' ? 65 : 45, 3, 3, 'FD')

    doc.setFontSize(10)
    doc.setTextColor(100, 116, 139)
    doc.text('Subtotal', summaryX + 5, finalY + 18)
    doc.setFont('helvetica', 'bold')
    doc.setTextColor(30, 41, 59)
    doc.text(`Rp ${sale.subtotal.toLocaleString('id-ID')}`, summaryX + 60, finalY + 18, {
      align: 'right',
    })

    if (sale.discount > 0) {
      doc.setFont('helvetica', 'normal')
      doc.setTextColor(220, 38, 38)
      doc.text('Diskon Tambahan', summaryX + 5, finalY + 26)
      doc.setFont('helvetica', 'bold')
      doc.text(`- Rp ${sale.discount.toLocaleString('id-ID')}`, summaryX + 60, finalY + 26, {
        align: 'right',
      })
    }

    doc.setDrawColor(226, 232, 240)
    doc.line(summaryX + 5, finalY + 32, summaryX + 60, finalY + 32)

    doc.setFontSize(12)
    doc.setTextColor(37, 99, 235) // primary blue
    doc.text('Total Akhir', summaryX + 5, finalY + 40)
    doc.text(`Rp ${sale.total.toLocaleString('id-ID')}`, summaryX + 60, finalY + 40, {
      align: 'right',
    })

    if (sale.payment_method === 'CASH') {
      doc.setFontSize(10)
      doc.setTextColor(100, 116, 139)
      doc.text('Tunai (Diterima)', summaryX + 5, finalY + 50)
      doc.setTextColor(30, 41, 59)
      doc.text(
        `Rp ${(sale.amount_paid || sale.total).toLocaleString('id-ID')}`,
        summaryX + 60,
        finalY + 50,
        { align: 'right' },
      )

      doc.setTextColor(100, 116, 139)
      doc.text('Kembali', summaryX + 5, finalY + 58)
      doc.setTextColor(22, 163, 74) // green
      doc.setFont('helvetica', 'bold')
      doc.text(
        `Rp ${(sale.change_amount || 0).toLocaleString('id-ID')}`,
        summaryX + 60,
        finalY + 58,
        { align: 'right' },
      )

      doc.setFont('helvetica', 'normal')
      doc.setTextColor(100, 116, 139)
      doc.text('Pembayaran:', summaryX + 5, finalY + 68)
      doc.setTextColor(30, 41, 59)
      doc.setFont('helvetica', 'bold')
      doc.text(sale.payment_method, summaryX + 60, finalY + 68, { align: 'right' })
    } else {
      doc.setFontSize(10)
      doc.setTextColor(100, 116, 139)
      doc.setFont('helvetica', 'normal')
      doc.text('Pembayaran:', summaryX + 5, finalY + 50)
      doc.setTextColor(22, 163, 74) // success green
      doc.setFont('helvetica', 'bold')
      doc.text(sale.payment_method, summaryX + 60, finalY + 50, { align: 'right' })
    }

    // Footer
    const footerY = finalY + 80
    doc.setDrawColor(226, 232, 240)
    doc.line(15, footerY, pageWidth - 15, footerY)

    doc.setFontSize(11)
    doc.setTextColor(30, 41, 59)
    doc.setFont('helvetica', 'bold')
    doc.text('Terima kasih telah berbelanja di Toko Bisnis Rinzi', pageWidth / 2, footerY + 10, {
      align: 'center',
    })

    doc.setFontSize(9)
    doc.setTextColor(100, 116, 139)
    doc.setFont('helvetica', 'normal')
    doc.text(
      'Barang yang sudah dibeli tidak dapat ditukar atau dikembalikan.',
      pageWidth / 2,
      footerY + 15,
      { align: 'center' },
    )

    // Dapatkan Blob
    const pdfBlob = doc.output('blob')

    // Upload
    const result = await posApi.uploadInvoice(sale.id, pdfBlob)
    console.log('Upload result:', result)
    alert(
      `File PDF berhasil dibuat dan dikirim ke MinIO!\nTautan: ${result.data?.invoice_url || result.invoice_url || 'URL tidak ditemukan'}`,
    )
  } catch (e: any) {
    console.error('Gagal membuat dan upload invoice PDF:', e)
    alert(`Gagal membuat dan upload invoice PDF: ${e.message || e}`)
  }
}

// Fungsi Checkout
const processCheckout = async () => {
  // Simpan data keranjang sebelum di-clear oleh store
  const savedItems = JSON.parse(JSON.stringify(cartStore.items))
  const savedDiscount = cartStore.discountAmount
  const savedSubtotal = cartStore.subtotal
  const savedTotal = cartStore.grandTotal
  const savedAmountPaid = amountPaid.value || savedTotal
  const savedChangeAmount = changeAmount.value > 0 ? changeAmount.value : 0
  const method = paymentMethod.value

  try {
    const res = await cartStore.checkout(method, savedAmountPaid)
    if (res && res.invoice_number) {
      // Proses invoice PDF dan Upload di belakang layar
      const dummySale = {
        id: res.id,
        invoice_number: res.invoice_number,
        transaction_date: res.transaction_date || new Date().toISOString(),
        subtotal: savedSubtotal,
        discount: savedDiscount,
        total: savedTotal,
        amount_paid: savedAmountPaid,
        change_amount: savedChangeAmount,
        payment_method: method,
      }
      
      // Simpan data untuk dirender ke struk
      completedSale.value = dummySale as Sale
      completedItems.value = savedItems

      // Tampilkan modal struk
      showCheckoutModal.value = false
      showReceiptModal.value = true

      // Jangan await (fire and forget)
      generateAndUploadInvoice(dummySale as Sale, savedItems, cashierName.value)
      
      amountPaid.value = null
    } else {
      alert('Pembayaran Berhasil disimpan secara offline/antrean.')
      showCheckoutModal.value = false
    }
  } catch (error) {
    alert('Terjadi kesalahan saat memproses pembayaran.')
  }
}

// --- LOGIKA RECEIPT (STRUK) ---
const showReceiptModal = ref(false)
const completedSale = ref<Sale | null>(null)
const completedItems = ref<SaleItem[]>([])

const closeReceiptModal = () => {
  showReceiptModal.value = false
  completedSale.value = null
  completedItems.value = []
}

const printReceipt = () => {
  window.print()
}
</script>
