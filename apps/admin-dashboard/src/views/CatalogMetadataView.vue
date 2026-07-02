<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { inventoryApi } from '@frontend/api/inventory';
import type { Category, Brand, Unit } from '@frontend/types/inventory';

// --- STATE MANAGEMENT ---
type MetadataTab = 'categories' | 'brands' | 'units';
const activeTab = ref<MetadataTab>('categories');
const isLoading = ref(false);
const isModalOpen = ref(false);
const errorMessage = ref('');

// Data Arrays
const categories = ref<Category[]>([]);
const brands = ref<Brand[]>([]);
const units = ref<Unit[]>([]);

// Form CRUD State
const isEditMode = ref(false);
const currentId = ref<string | null>(null);
const form = ref({
  code: '',
  name: ''
});

// --- FETCH DATA OPERATIONS ---
const fetchTabData = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    if (activeTab.value === 'categories') {
      categories.value = await inventoryApi.getCategories();
    } else if (activeTab.value === 'brands') {
      brands.value = await inventoryApi.getBrands();
    } else if (activeTab.value === 'units') {
      units.value = await inventoryApi.getUnits();
    }
  } catch (error: any) {
    errorMessage.value = error?.message || `Gagal memuat data master ${activeTab.value}.`;
  } finally {
    isLoading.value = false;
  }
};

const changeTab = (tab: MetadataTab) => {
  activeTab.value = tab;
  fetchTabData();
};

// --- MODAL CONTROLS ---
const openCreateModal = () => {
  isEditMode.value = false;
  currentId.value = null;
  form.value = { code: '', name: '' };
  isModalOpen.value = true;
};

const openEditModal = (item: Category | Brand | Unit) => {
  isEditMode.value = true;
  currentId.value = item.id || null;
  form.value = {
    code: item.code || '',
    name: item.name || ''
  };
  isModalOpen.value = true;
};

// --- SUBMIT (CREATE / UPDATE) ---
const handleSubmit = async () => {
  if (!form.value.code || !form.value.name) {
    errorMessage.value = 'Kode dan Nama wajib diisi.';
    return;
  }

  isLoading.value = true;
  try {
    const payload = { code: form.value.code.toUpperCase(), name: form.value.name };

    if (activeTab.value === 'categories') {
      if (isEditMode.value && currentId.value) {
        await inventoryApi.updateCategory(currentId.value, payload);
      } else {
        await inventoryApi.createCategory(payload);
      }
    } else if (activeTab.value === 'brands') {
      if (isEditMode.value && currentId.value) {
        await inventoryApi.updateBrand(currentId.value, payload);
      } else {
        await inventoryApi.createBrand(payload);
      }
    } else if (activeTab.value === 'units') {
      if (isEditMode.value && currentId.value) {
        await inventoryApi.updateUnit(currentId.value, payload);
      } else {
        await inventoryApi.createUnit(payload);
      }
    }

    await fetchTabData();
    isModalOpen.value = false;
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal menyimpan perubahan master data.';
  } finally {
    isLoading.value = false;
  }
};

// --- DELETE ACTION ---
const handleDelete = async (id: string) => {
  if (!confirm(`Apakah Anda yakin ingin menghapus entitas ini? Tindakan ini dapat merusak relasi integritas katalog produk.`)) return;

  isLoading.value = true;
  try {
    if (activeTab.value === 'categories') {
      await inventoryApi.deleteCategory(id);
    } else if (activeTab.value === 'brands') {
      await inventoryApi.deleteBrand(id);
    } else if (activeTab.value === 'units') {
      await inventoryApi.deleteUnit(id);
    }
    await fetchTabData();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal menghapus entitas master.';
  } finally {
    isLoading.value = false;
  }
};

// Computed list data untuk mempermudah render tabel dinamis
const currentList = computed(() => {
  if (activeTab.value === 'categories') return categories.value;
  if (activeTab.value === 'brands') return brands.value;
  return units.value;
});

onMounted(() => {
  fetchTabData();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
          <i class="pi pi-bookmarks text-indigo-600"></i>
          Metadata & Atribut Produk Master
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Otoritas pengelolaan data dasar Kategori, Merek (Brand), dan Satuan Unit untuk menyelaraskan katalog stok inventaris.
        </p>
      </div>
      <button 
        @click="openCreateModal"
        class="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-semibold px-4 py-2.5 rounded-xl shadow-md transition-all duration-200"
      >
        <i class="pi pi-plus text-xs"></i>
        Tambah Atribut Baru
      </button>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl flex items-start gap-3">
      <i class="pi pi-exclamation-circle text-lg mt-0.5"></i>
      <span class="text-sm font-medium">{{ errorMessage }}</span>
    </div>

    <div class="flex border-b border-slate-200 space-x-2">
      <button
        @click="changeTab('categories')"
        class="px-5 py-3 text-sm font-semibold border-b-2 transition-all flex items-center gap-2"
        :class="activeTab === 'categories' ? 'border-indigo-600 text-indigo-600 font-bold' : 'border-transparent text-slate-500 hover:text-slate-800'"
      >
        <i class="pi pi-folder"></i> Kategori Barang
      </button>
      <button
        @click="changeTab('brands')"
        class="px-5 py-3 text-sm font-semibold border-b-2 transition-all flex items-center gap-2"
        :class="activeTab === 'brands' ? 'border-indigo-600 text-indigo-600 font-bold' : 'border-transparent text-slate-500 hover:text-slate-800'"
      >
        <i class="pi pi-tag"></i> Merek / Brand
      </button>
      <button
        @click="changeTab('units')"
        class="px-5 py-3 text-sm font-semibold border-b-2 transition-all flex items-center gap-2"
        :class="activeTab === 'units' ? 'border-indigo-600 text-indigo-600 font-bold' : 'border-transparent text-slate-500 hover:text-slate-800'"
      >
        <i class="pi pi-box"></i> Satuan Unit
      </button>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Sinkronisasi data master...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">ID Sistem</th>
              <th class="py-4 px-6">Kode Unik Atribut</th>
              <th class="py-4 px-6">Nama Atribut</th>
              <th class="py-4 px-6 text-right">Aksi Manajemen</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-if="currentList.length === 0 && !isLoading">
              <td colspan="4" class="py-12 text-center text-slate-400 font-normal">
                Belum ada data master untuk sub-katalog ini. Klik "Tambah Atribut Baru" untuk mendaftarkannya.
              </td>
            </tr>
            <tr v-for="item in currentList" :key="item.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6 font-mono text-xs text-slate-400 select-all">{{ item.id }}</td>
              <td class="py-4 px-6"><span class="px-2.5 py-1 bg-slate-100 text-slate-700 font-mono text-xs rounded-md border">{{ item.code }}</span></td>
              <td class="py-4 px-6 text-slate-900 font-semibold">{{ item.name }}</td>
              <td class="py-4 px-6 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button @click="openEditModal(item)" class="p-2 text-slate-500 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors">
                    <i class="pi pi-pencil text-xs"></i>
                  </button>
                  <button @click="handleDelete(item.id || '')" class="p-2 text-slate-500 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors">
                    <i class="pi pi-trash text-xs"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-200 overflow-hidden transform transition-all animate-fade-in">
        <div class="p-6 border-b border-slate-100 flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-900 uppercase tracking-wide">
            {{ isEditMode ? 'Ubah Atribut' : 'Daftar Atribut' }} - {{ activeTab.toUpperCase() }}
          </h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600">
            <i class="pi pi-times"></i>
          </button>
        </div>
        
        <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Kode Atribut (Harus Unik)</label>
            <input 
              v-model="form.code"
              type="text"
              required
              placeholder="Contoh: ATK, SONY, PCS"
              class="w-full px-3 py-2 font-mono uppercase border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"/>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Atribut Resmi</label>
            <input 
              v-model="form.name"
              type="text"
              required
              placeholder="Contoh: Alat Tulis Kantor, Sony Corp, Pieces"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"/>
          </div>

          <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-100 mt-6">
            <button type="button" @click="isModalOpen = false" class="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 rounded-lg border border-slate-200 transition-colors">
              Batal
            </button>
            <button type="submit" :disabled="isLoading" class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-md transition-colors disabled:opacity-50">
              <span>{{ isLoading ? 'Menyimpan...' : 'Simpan Master Atribut' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>