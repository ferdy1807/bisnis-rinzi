<template>
  <div class="space-y-8">
    <div class="flex flex-col xl:flex-row xl:items-center justify-between gap-6 border-b border-outline-variant/20 pb-6">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">Katalog Master Komprehensif</h2>
        <p class="text-body-md text-on-surface-variant mt-1">
          Konsol manajemen agregasi inventaris terpadu: Sinkronisasi data real-time.
        </p>
      </div>
      <div class="flex flex-wrap items-center gap-4 bg-surface-container-low p-3 rounded-2xl border border-outline-variant/30 shadow-sm">
        <button
          @click="openCreateModal"
          class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2.5 rounded-xl shadow-sm hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all text-label-md font-bold cursor-pointer"
        >
          <span class="material-symbols-outlined text-md">add_box</span>
          Tambah Produk Baru
        </button>
        <button
          @click="syncAllCatalogs(true)"
          :disabled="isLoading"
          class="flex items-center gap-2 bg-surface text-on-surface border border-outline-variant px-4 py-2.5 rounded-xl shadow-sm hover:bg-surface-container-high active:scale-95 transition-all text-label-md font-bold disabled:opacity-60 cursor-pointer"
        >
          <span class="material-symbols-outlined text-md" :class="isLoading ? 'animate-spin' : ''">sync</span>
          Sinkronisasi Katalog Master
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium flex items-center gap-2 animate-fade-in">
      <span class="material-symbols-outlined">error</span>
      <span>{{ errorMessage }}</span>
    </div>

    <div class="flex border-b border-outline-variant/30 gap-6 text-title-small font-bold">
      <button 
        @click="activeTab = 'ALL'" 
        class="pb-3 border-b-2 transition-all px-2 cursor-pointer"
        :class="activeTab === 'ALL' ? 'border-primary text-primary font-black' : 'border-transparent text-on-surface-variant hover:text-on-surface'"
      >
        Semua Inventaris ({{ allItemsCount }})
      </button>
      <button 
        @click="activeTab = 'RETAIL'" 
        class="pb-3 border-b-2 transition-all px-2 cursor-pointer"
        :class="activeTab === 'RETAIL' ? 'border-primary text-primary font-black' : 'border-transparent text-on-surface-variant hover:text-on-surface'"
      >
        Produk POS Ritel ({{ retailItemsCount }})
      </button>
      <button 
        @click="activeTab = 'RENTAL'" 
        class="pb-3 border-b-2 transition-all px-2 cursor-pointer"
        :class="activeTab === 'RENTAL' ? 'border-primary text-primary font-black' : 'border-transparent text-on-surface-variant hover:text-on-surface'"
      >
        Armada Rental ({{ rentalItemsCount }})
      </button>
    </div>

    <div class="bg-surface border border-outline-variant/30 rounded-2xl p-4 shadow-sm flex flex-col sm:flex-row gap-4 items-center justify-between">
      <div class="relative w-full sm:w-80">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-on-surface-variant">
          <span class="material-symbols-outlined text-md">search</span>
        </span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari berdasarkan nama, SKU, atau kode..."
          class="w-full pl-10 pr-4 py-2 bg-surface-container-lowest border border-outline-variant/60 rounded-xl text-body-medium outline-none focus:ring-2 focus:ring-primary/20 transition-all text-on-surface"
        />
      </div>
      <div class="text-[11px] font-black uppercase tracking-wider text-on-surface-variant font-mono">
        Menampilkan: {{ filteredCatalogList.length }} Entri Valid
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[300px] bg-surface/50 rounded-2xl border border-outline-variant/20">
      <span class="animate-spin material-symbols-outlined text-5xl text-primary mb-3">sync</span>
      <p class="text-title-small font-bold text-primary animate-pulse">Menghubungkan ke inventory_db & rental_db...</p>
    </div>

    <div v-else class="bg-surface border border-outline-variant/30 rounded-2xl shadow-sm overflow-hidden flex flex-col">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-body-sm border-collapse">
          <thead>
            <tr class="bg-surface-container-lowest text-on-surface-variant text-[11px] font-bold tracking-wider uppercase border-b border-outline-variant/30">
              <th class="py-4 px-6 rounded-tl-2xl">Kode</th>
              <th class="py-4 px-6">Identitas Aset (Nama)</th>
              <th class="py-4 px-6 text-center">Tipe Sistem</th>
              <th class="py-4 px-6 text-center">Status Gudang / Fisik</th>
              <th class="py-4 px-6 text-center">Preview Gambar</th>
              <th class="py-4 px-6 text-right">Tarif Jual / Sewa</th>
              <th class="py-4 px-6 text-center rounded-tr-2xl">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/10 text-xs font-semibold text-on-surface">
            <tr v-if="filteredCatalogList.length === 0">
              <td colspan="8" class="px-6 py-12 text-center font-medium text-on-surface-variant italic">
                Tidak ada entri katalog master yang sesuai dengan parameter pencarian saat ini.
              </td>
            </tr>
            <tr 
              v-for="item in filteredCatalogList" 
              :key="item.global_id" 
              class="hover:bg-surface-container-low/40 transition-colors">
              <td class="px-6 py-4 font-mono text-xs text-on-surface-variant select-all">
                {{ item.code_or_sku }}
              </td>
              <td class="px-6 py-4">
                <div class="font-bold text-on-surface text-sm">{{ item.name }}</div>
                <div class="text-[10px] text-outline font-normal mt-0.5">ID: {{ item.id }}</div>
              </td>
              <td class="px-6 py-4 text-center">
                <span 
                  class="px-2.5 py-0.5 rounded text-[9px] font-black tracking-widest uppercase border"
                  :class="item.type === 'RETAIL' 
                    ? 'bg-primary/10 text-primary border-primary/20' 
                    : 'bg-secondary/10 text-secondary border-secondary/20'">
                  {{ item.type }}
                </span>
              </td>
              <td class="py-3.5 px-6 text-center">
                <span class="inline-flex items-center gap-1 text-[11px] font-bold uppercase tracking-widest px-2.5 py-1 rounded-md"
                  :class="item.stock > 10 ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'">
                  <span class="material-symbols-outlined text-[14px]">inventory_2</span>
                  {{ formatNumber(item.stock) }} Unit
                </span>
              </td>
              <td class="py-3.5 px-6 text-center">
                <div v-if="item.media_url" class="flex justify-center">
                  <img :src="item.media_url" alt="Preview" class="h-12 w-12 object-cover rounded-lg border border-outline-variant/30 shadow-sm" />
                </div>
                <div v-else class="flex justify-center">
                  <div class="h-12 w-12 rounded-lg bg-surface-container flex items-center justify-center border border-outline-variant/30">
                    <span class="material-symbols-outlined text-outline text-xl">image_not_supported</span>
                  </div>
                </div>
              </td>
              <td class="py-3.5 px-6 text-right">
                <span class="text-sm font-bold font-mono tracking-tight text-on-surface">Rp {{ item.price.toLocaleString('id-ID') }}</span>
              </td>
              <td class="py-3.5 px-6 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button @click="openEditModal(item)" class="text-primary hover:text-primary-container p-1 rounded-full hover:bg-surface-container-high transition-colors" title="Edit">
                    <span class="material-symbols-outlined text-sm">edit</span>
                  </button>
                  <button @click="confirmDelete(item)" class="text-error hover:text-error-container p-1 rounded-full hover:bg-surface-container-high transition-colors" title="Hapus">
                    <span class="material-symbols-outlined text-sm">delete</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="isModalOpen" class="fixed inset-0 bg-black/60 backdrop-blur-xs z-50 flex items-center justify-center p-4 animate-fade-in">
      <div class="bg-surface border border-outline-variant/30 rounded-3xl w-full max-w-2xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
        <div class="p-6 border-b border-outline-variant/20 flex items-center justify-between bg-surface-container-low">
          <div class="flex justify-between items-center w-full">
            <h3 class="text-title-lg font-bold text-on-surface">{{ isEditMode ? 'Edit' : 'Tambah' }} Entri Katalog {{ formType === 'RETAIL' ? 'POS Ritel' : 'Armada Rental' }}</h3>
            <button @click="isModalOpen = false" class="text-on-surface-variant hover:text-on-surface rounded-full p-1 transition-colors">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>
        </div>

        <form @submit.prevent="handleFormSubmit" class="p-6 space-y-5 overflow-y-auto flex-1 custom-scrollbar">
          <div v-if="!isEditMode" class="flex gap-4 mb-4 p-1.5 bg-surface-container-low rounded-2xl w-max border border-outline-variant/30">
            <label 
              @click="formType = 'RETAIL'"
              :class="formType === 'RETAIL' ? 'bg-primary text-on-primary shadow-sm' : 'text-on-surface-variant'"
              class="px-4 py-2 rounded-xl flex items-center gap-2 cursor-pointer select-none font-bold text-sm transition-all"
            >
              <span class="material-symbols-outlined text-sm">{{ formType === 'RETAIL' ? 'radio_button_checked' : 'radio_button_unchecked' }}</span>
              <span>POS Ritel</span>
            </label>
            <label 
              @click="formType = 'RENTAL'"
              :class="formType === 'RENTAL' ? 'bg-secondary text-on-secondary shadow-sm' : 'text-on-surface-variant'"
              class="px-4 py-2 rounded-xl flex items-center gap-2 cursor-pointer select-none font-bold text-sm transition-all"
            >
              <span class="material-symbols-outlined text-sm">{{ formType === 'RENTAL' ? 'radio_button_checked' : 'radio_button_unchecked' }}</span>
              <span>Rental</span>
            </label>
          </div>

          <div v-if="formType === 'RETAIL'" class="space-y-4 animate-fade-in">
            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Kategori Barang *</label>
                <select v-model="formRetail.categoryIndex" required class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-semibold outline-none focus:ring-2 focus:ring-primary/20">
                  <option value="" disabled>Pilih Kategori</option>
                  <option v-for="(cat, idx) in mockCategories" :key="cat.id" :value="idx">{{ cat.name }} [{{ cat.code }}]</option>
                </select>
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Merek / Brand *</label>
                <select v-model="formRetail.brandIndex" required class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-semibold outline-none focus:ring-2 focus:ring-primary/20">
                  <option value="" disabled>Pilih Brand</option>
                  <option v-for="(b, idx) in mockBrands" :key="b.id" :value="idx">{{ b.name }} [{{ b.code }}]</option>
                </select>
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Nama Item Produk *</label>
              <input v-model="formRetail.name" type="text" required placeholder="Contoh: Lampu LED Tubles 12W" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary/20 font-semibold" />
            </div>

            <div v-if="!isEditMode" class="bg-surface-container-low p-3.5 rounded-xl border border-outline-variant/40 flex items-center justify-between">
              <div>
                <span class="text-[10px] font-black text-on-surface-variant uppercase tracking-widest block">Generated SKU Real-time</span>
                <span class="text-xs font-mono font-black text-primary mt-1 block tracking-wider">{{ generatedSku }}</span>
              </div>
              <span class="text-[9px] bg-primary/10 text-primary font-bold px-2 py-0.5 rounded font-mono border border-primary/20">AUTO</span>
            </div>

            <div class="grid grid-cols-3 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Satuan Basis *</label>
                <select v-model="formRetail.baseUnitCode" required class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-semibold outline-none focus:ring-2 focus:ring-primary/20">
                  <option value="" disabled>Pilih Satuan</option>
                  <option v-for="unit in mockUnits" :key="unit.id" :value="unit.code">{{ unit.code }} - {{ unit.name }}</option>
                </select>
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Harga Pokok / Modal (Rp) *</label>
                <input v-model.number="formRetail.costPrice" type="number" required min="0" placeholder="0" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-mono font-bold" />
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Harga Jual (Rp) *</label>
                <input v-model.number="formRetail.sellingPrice" type="number" required min="0" placeholder="0" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-mono font-bold" />
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Stok Awal Gudang *</label>
                <input v-model.number="formRetail.initialQty" type="number" required min="0" placeholder="0" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-mono font-bold" />
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Foto Produk Toko</label>
              <div class="border-2 border-dashed border-outline-variant/60 rounded-xl p-4 text-center bg-surface-container-lowest relative group flex flex-col items-center justify-center">
                <input type="file" :required="!isEditMode" accept="image/*" @change="handleFileChange" class="absolute inset-0 opacity-0 cursor-pointer z-10" />
                <img v-if="previewUrl" :src="previewUrl" class="h-32 object-contain mb-2 rounded-lg border border-outline-variant/30" />
                <span v-else class="material-symbols-outlined text-3xl text-outline mb-1 block group-hover:text-primary transition-colors">add_photo_alternate</span>
                <span class="text-xs font-semibold text-on-surface-variant block truncate">{{ selectedFile ? selectedFile.name : 'Pilih berkas gambar dari perangkat...' }}</span>
              </div>
            </div>
          </div>

          <div v-if="formType === 'RENTAL'" class="space-y-4 animate-fade-in">
            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Kategori Sewa *</label>
                <select v-model="formRental.categoryId" required class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-semibold outline-none focus:ring-2 focus:ring-secondary/20">
                  <option value="" disabled>Pilih Kategori Rental</option>
                  <option v-for="rc in mockRentalCategories" :key="rc.id" :value="rc.id">{{ rc.name }}</option>
                </select>
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Nama Unit Rental *</label>
              <input v-model="formRental.name" type="text" required placeholder="Contoh: Boks Hantaran Akrilik Gold Premium XL" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-sm outline-none focus:ring-2 focus:ring-secondary/20 font-semibold" />
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Deskripsi Unit</label>
              <textarea v-model="formRental.description" rows="2" placeholder="Tulis rincian fisik..." class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-medium resize-none outline-none focus:ring-2 focus:ring-secondary/20"></textarea>
            </div>

            <div class="grid grid-cols-3 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Tarif Sewa (Rp) *</label>
                <input v-model.number="formRental.rentalPrice" type="number" required min="0" placeholder="0" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-mono font-bold" />
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Kuantitas Armada *</label>
                <input v-model.number="formRental.quantityAvailable" type="number" required min="0" placeholder="0" class="w-full px-3 py-2 bg-surface-container-lowest border border-outline-variant rounded-xl text-xs font-mono font-bold" />
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-black text-on-surface-variant uppercase tracking-wider">Foto Unit Sewa</label>
              <div class="border-2 border-dashed border-outline-variant/60 rounded-xl p-4 text-center bg-surface-container-lowest relative group flex flex-col items-center justify-center">
                <input type="file" :required="!isEditMode" accept="image/*" @change="handleFileChange" class="absolute inset-0 opacity-0 cursor-pointer z-10" />
                <img v-if="previewUrl" :src="previewUrl" class="h-32 object-contain mb-2 rounded-lg border border-outline-variant/30" />
                <span v-else class="material-symbols-outlined text-3xl text-outline mb-1 block group-hover:text-secondary transition-colors">add_photo_alternate</span>
                <span class="text-xs font-semibold text-on-surface-variant block truncate">{{ selectedFile ? selectedFile.name : 'Pilih berkas gambar...' }}</span>
              </div>
            </div>
          </div>

          <div class="flex justify-end pt-4 border-t border-outline-variant/20 mt-4">
            <button type="submit" :disabled="isActionLoading" class="bg-primary text-on-primary px-6 py-2 rounded-xl text-label-md font-bold hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all shadow-sm disabled:opacity-50">
              {{ isActionLoading ? 'Menyimpan...' : (isEditMode ? 'Simpan Perubahan' : 'Simpan Produk') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useInventoryStore } from '@frontend/stores/inventory';
import { useRentalStore } from '@frontend/stores/rental';
import { inventoryApi } from '@frontend/api/inventory';
import { rentalApi } from '@frontend/api/rental';

// --- INITIALIZE STORES ---
const inventoryStore = useInventoryStore();
const rentalStore = useRentalStore();

const { products: retailProducts } = storeToRefs(inventoryStore);
const { products: rentalProducts } = storeToRefs(rentalStore);

// --- LOCAL SYSTEM STATE ---
const isLoading = ref(false);
const isActionLoading = ref(false);
const isModalOpen = ref(false);
const formType = ref<'RETAIL' | 'RENTAL'>('RETAIL');
const isEditMode = ref(false);
const editingId = ref<string | null>(null);
const activeTab = ref<'ALL' | 'RETAIL' | 'RENTAL'>('ALL');
const searchQuery = ref('');
const errorMessage = ref('');

// Media upload references
const selectedFile = ref<File | null>(null);
const fileExtension = ref('jpg');
const previewUrl = ref<string | null>(null);

// --- FALLBACK METADATA MASTER FROM SQL DEFINITION ---
const mockCategories = ref<any[]>([]);
const mockBrands = ref<any[]>([]);
const mockRentalCategories = ref<any[]>([]);
const mockUnits = ref<any[]>([]);

const fetchReferenceData = async () => {
  try {
    const [cats, brands, rCats, units] = await Promise.all([
      inventoryApi.getCategories(),
      inventoryApi.getBrands(),
      rentalApi.getCategories(),
      inventoryApi.getUnits()
    ]);
    mockCategories.value = cats || [];
    mockBrands.value = brands || [];
    mockRentalCategories.value = (rCats as any)?.data || rCats || [];
    mockUnits.value = units || [];
  } catch (error) {
    console.error('Failed to load reference data', error);
  }
};

// --- FORM DATA REFS ---
const formRetail = ref({
  categoryIndex: '' as string | number,
  brandIndex: '' as string | number,
  name: '',
  baseUnitCode: 'PCS',
  costPrice: 0,
  sellingPrice: 0,
  initialQty: 10
});

const formRental = ref({
  categoryId: '',
  code: '',
  name: '',
  description: '',
  rentalPrice: 50000,
  quantityAvailable: 5
});

// --- REAL-TIME AUTOMATED SKU GENERATOR (FOR PRODUCTS) ---
const generatedSku = computed(() => {
  if (formRetail.value.categoryIndex === '' || formRetail.value.brandIndex === '') {
    return 'SKU-AUTO-WAITING';
  }
  const catCode = mockCategories.value[Number(formRetail.value.categoryIndex)].code;
  const brandCode = mockBrands.value[Number(formRetail.value.brandIndex)].code;
  
  // Gunakan nama produk utuh, ganti spasi dengan _, dan jadikan kapital
  const nameSanitized = formRetail.value.name.trim().replace(/\s+/g, '_').replace(/[^a-zA-Z0-9_]/g, '').toUpperCase();
  const nameSuffix = nameSanitized.length > 0 ? nameSanitized : 'PRD';

  return `${catCode}-${brandCode}-${nameSuffix}`;
});

// Reset file saat tipe form berpindah
watch(formType, () => {
  selectedFile.value = null;
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
    previewUrl.value = null;
  }
});

const openCreateModal = () => {
  isEditMode.value = false;
  editingId.value = null;
  isModalOpen.value = true;
  // Reset state form
  formRetail.value = { categoryIndex: '', brandIndex: '', name: '', baseUnitCode: 'PCS', costPrice: 0, sellingPrice: 0, initialQty: 10 };
  formRental.value = { categoryId: '', code: '', name: '', description: '', rentalPrice: 50000, quantityAvailable: 5 };
  selectedFile.value = null;
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
    previewUrl.value = null;
  }
};

const openEditModal = (item: any) => {
  isEditMode.value = true;
  editingId.value = item.id;
  formType.value = item.type;
  isModalOpen.value = true;
  previewUrl.value = item.media_url || null;
  selectedFile.value = null;

  if (item.type === 'RETAIL') {
    const cIdx = mockCategories.value.findIndex(c => c.id === item.category_id);
    const bIdx = item.brand_id ? mockBrands.value.findIndex((b: any) => b.id === item.brand_id) : -1;
    formRetail.value = {
      categoryIndex: cIdx >= 0 ? String(cIdx) : '',
      brandIndex: bIdx >= 0 ? String(bIdx) : '',
      name: item.name,
      baseUnitCode: item.base_unit_code || 'PCS',
      costPrice: item.cost_price || 0,
      sellingPrice: item.price,
      initialQty: item.stock
    };
  } else {
    formRental.value = {
      categoryId: item.category_id,
      code: item.code_or_sku,
      name: item.name,
      description: item.description || '',
      rentalPrice: item.price,
      quantityAvailable: item.stock
    };
  }
};

const confirmDelete = async (item: any) => {
  if (confirm(`Apakah Anda yakin ingin menghapus ${item.name}?`)) {
    try {
      if (item.type === 'RETAIL') {
        await inventoryApi.deleteProduct(item.id);
      } else {
        await rentalApi.deleteProduct(item.id);
      }
      syncAllCatalogs();
    } catch (error) {
      console.error('Gagal menghapus produk', error);
      alert('Gagal menghapus produk.');
    }
  }
};

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    selectedFile.value = target.files[0];
    const ext = target.files[0].name.split('.').pop();
    fileExtension.value = ext ? ext.toLowerCase() : 'jpg';
    if (previewUrl.value) URL.revokeObjectURL(previewUrl.value);
    previewUrl.value = URL.createObjectURL(target.files[0]);
  }
};

// --- TRANSACTION WORKFLOW DISPATCHER ---
const handleFormSubmit = async () => {
  errorMessage.value = '';
  isActionLoading.value = true;

  try {
    if (formType.value === 'RETAIL') {
      // 1. Validasi SKU
      if (generatedSku.value === 'SKU-AUTO-WAITING') throw new Error('Silakan pilih kategori dan brand terlebih dahulu.');
      
      const categoryId = mockCategories.value[Number(formRetail.value.categoryIndex)].id;
      const brandId = mockBrands.value[Number(formRetail.value.brandIndex)].id;
      const targetFileName = `${generatedSku.value}.${fileExtension.value}`;

      // Mutasi store inventory via API
      const payload = {
        sku: isEditMode.value && editingId.value ? undefined : generatedSku.value, // Do not modify SKU on edit
        category_id: categoryId,
        brand_id: brandId,
        name: formRetail.value.name,
        base_unit_code: formRetail.value.baseUnitCode,
        cost_price: formRetail.value.costPrice,
        selling_price: formRetail.value.sellingPrice,
        initial_qty: formRetail.value.initialQty,
        is_active: true
      };
      
      let prodId = editingId.value;
      if (isEditMode.value && prodId) {
        await inventoryApi.updateProduct(prodId, payload as any);
      } else {
        const createdProd = await inventoryApi.createProduct(payload as any);
        prodId = createdProd?.id || (createdProd as any)?.data?.id;
      }

      if (selectedFile.value && prodId) {
        const formData = new FormData();
        formData.append('media', selectedFile.value, targetFileName);
        await inventoryApi.uploadProductMedia(prodId, formData);
      }
      
    } else {
      // Skema Rental Product
      const selectedCat = mockRentalCategories.value.find(rc => rc.id === formRental.value.categoryId);
      if (!selectedCat) throw new Error('Kategori sewa tidak valid.');

      // Mutasi store rental via API
      const rentalNameSanitized = formRental.value.name.trim().replace(/\s+/g, '_');
      const rentalCode = `${selectedCat.code}_${rentalNameSanitized}`;
      const payload = {
         category_id: formRental.value.categoryId,
         code: isEditMode.value ? formRental.value.code : rentalCode,
         name: formRental.value.name,
         description: formRental.value.description,
         rental_price: formRental.value.rentalPrice,
         quantity_available: formRental.value.quantityAvailable,
         is_active: true
      };
      
      if (isEditMode.value && editingId.value) {
         await rentalApi.updateProduct(editingId.value, payload);
         if (selectedFile.value) {
            const targetRentalFileName = `${editingId.value}.${fileExtension.value}`;
            await rentalApi.uploadProductMedia(editingId.value, selectedFile.value, targetRentalFileName);
         }
      } else {
         await rentalApi.createProduct(payload);

         if (selectedFile.value) {
            const resProds: any = await rentalApi.getProducts();
            const allProds = resProds?.data || resProds || [];
            const createdProd = allProds.find((p: any) => p.code === rentalCode);
            
            if (createdProd && createdProd.id) {
               const sanitizedRentalName = createdProd.id;
               const targetRentalFileName = `${sanitizedRentalName}.${fileExtension.value}`;
               await rentalApi.uploadProductMedia(createdProd.id, selectedFile.value, targetRentalFileName);
            }
         }
      }
    }

    // Refresh data katalog master setelah berhasil mutasi
    await syncAllCatalogs();
    isModalOpen.value = false;
  } catch (err: any) {
    errorMessage.value = err?.message || 'Gagal mendaftarkan aset baru ke kluster database.';
  } finally {
    isActionLoading.value = false;
  }
};

// --- DATA NORMALIZATION MAPPERS ---
const unifiedCatalog = computed(() => {
  const retailList = (retailProducts.value || []).map((p: any) => ({
    global_id: `retail-${p.id}`,
    id: p.id,
    type: 'RETAIL',
    code_or_sku: p.sku,
    name: p.name,
    category_id: p.category_id,
    brand_id: p.brand_id,
    base_unit_code: p.base_unit_code,
    stock: p.qty || p.stock?.qty || 0,
    cost_price: Number(p.cost_price || 0),
    price: Number(p.selling_price || 0),
    media_url: p.media && p.media.length > 0 ? inventoryApi.getMediaUrl(p.media[0].id) : null
  }));

  const rentalList = (rentalProducts.value || []).map((rp: any) => ({
    global_id: `rental-${rp.id}`,
    id: rp.id,
    type: 'RENTAL',
    code_or_sku: rp.code,
    name: rp.name,
    description: rp.description,
    category_id: rp.category_id,
    stock: rp.quantity_available,
    price: Number(rp.rental_price || 0),
    media_url: rp.object_name ? `http://localhost:9000/foto-produk-sewa/${rp.object_name}` : null
  }));

  return [...retailList, ...rentalList].sort((a, b) => a.name.localeCompare(b.name));
});

const allItemsCount = computed(() => unifiedCatalog.value.length);
const retailItemsCount = computed(() => unifiedCatalog.value.filter(i => i.type === 'RETAIL').length);
const rentalItemsCount = computed(() => unifiedCatalog.value.filter(i => i.type === 'RENTAL').length);

const filteredCatalogList = computed(() => {
  let result = unifiedCatalog.value;
  if (activeTab.value === 'RETAIL') result = result.filter(i => i.type === 'RETAIL');
  if (activeTab.value === 'RENTAL') result = result.filter(i => i.type === 'RENTAL');

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase().trim();
    result = result.filter(i => 
      i.name?.toLowerCase().includes(query) ||
      i.code_or_sku?.toLowerCase().includes(query) ||
      i.id?.toLowerCase().includes(query)
    );
  }
  return result;
});

const syncAllCatalogs = async (forceRefresh: boolean = false) => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    await Promise.all([
      inventoryStore.fetchProducts(forceRefresh),
      rentalStore.fetchProducts()
    ]);
  } catch (err: any) {
    errorMessage.value = err?.message || 'Gagal menyinkronkan agregasi data katalog master.';
  } finally {
    isLoading.value = false;
  }
};

const formatNumber = (value: number | undefined | null): string => {
  return new Intl.NumberFormat('id-ID', { minimumFractionDigits: 0 }).format(value ?? 0);
};

onMounted(() => {
  syncAllCatalogs();
  fetchReferenceData();
});
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
  height: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}
</style>