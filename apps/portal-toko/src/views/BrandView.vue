<template>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary">Merek & Satuan Produk</h2>
        <p class="text-body-md text-on-surface-variant">
          Kelola daftar referensi merek dan satuan barang.
        </p>
      </div>
      <div class="flex gap-3">
        <button
          @click="openAddBrandModal"
          class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2 rounded-xl shadow-md hover:bg-primary-container hover:text-on-primary-container transition-colors"
        >
          <span class="material-symbols-outlined">add</span>
          <span class="font-bold text-label-md">Tambah Merek</span>
        </button>
        <button
          @click="openAddUnitModal"
          class="flex items-center gap-2 bg-secondary text-on-secondary px-4 py-2 rounded-xl shadow-md hover:bg-secondary-container hover:text-on-secondary-container transition-colors"
        >
          <span class="material-symbols-outlined">add</span>
          <span class="font-bold text-label-md">Tambah Satuan</span>
        </button>
      </div>
    </div>

    <!-- 2 Column Layout (Mirror Kanan Kiri) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- KIRI: Data Table Merek -->
      <div
        class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]"
      >
        <div
          class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col sm:flex-row justify-between items-center gap-4"
        >
          <h3 class="font-headline-md text-primary flex items-center gap-2">
            <span class="material-symbols-outlined">branding_watermark</span> Daftar Merek
          </h3>
        </div>

        <div class="overflow-x-auto flex-1">
          <table class="w-full text-left">
            <thead>
              <tr
                class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
              >
                <th class="px-6 py-4">Kode</th>
                <th class="px-6 py-4">Nama Merek</th>
                <th class="px-6 py-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-outline-variant/10">
              <tr v-if="isLoadingBrand">
                <td colspan="3" class="px-6 py-12 text-center text-primary">
                  <span class="animate-spin material-symbols-outlined text-4xl">sync</span>
                </td>
              </tr>
              <tr v-else-if="brands.length === 0">
                <td colspan="3" class="px-6 py-12 text-center text-on-surface-variant">
                  <p>Merek tidak ditemukan.</p>
                </td>
              </tr>
              <tr
                v-else
                v-for="b in brands"
                :key="b.id"
                class="hover:bg-surface-container-low/50 transition-colors group"
              >
                <td class="px-6 py-4 font-label-sm text-primary">{{ b.code }}</td>
                <td class="px-6 py-4 text-body-sm font-bold text-on-surface">{{ b.name }}</td>
                <td class="px-6 py-4 text-center space-x-2">
                  <button
                    @click="openEditBrandModal(b)"
                    class="p-2 text-primary hover:bg-primary/10 rounded-full transition-colors inline-flex items-center justify-center"
                    title="Edit Merek"
                  >
                    <span class="material-symbols-outlined text-sm">edit</span>
                  </button>
                  <button
                    @click="confirmDeleteBrand(b)"
                    class="p-2 text-error hover:bg-error/10 rounded-full transition-colors inline-flex items-center justify-center"
                    title="Hapus Merek"
                  >
                    <span class="material-symbols-outlined text-sm">delete</span>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- KANAN: Data Table Satuan -->
      <div
        class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]"
      >
        <div
          class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col sm:flex-row justify-between items-center gap-4"
        >
          <h3 class="font-headline-md text-secondary flex items-center gap-2">
            <span class="material-symbols-outlined">straighten</span> Daftar Satuan
          </h3>
        </div>

        <div class="overflow-x-auto flex-1">
          <table class="w-full text-left">
            <thead>
              <tr
                class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
              >
                <th class="px-6 py-4">Kode</th>
                <th class="px-6 py-4">Nama Satuan</th>
                <th class="px-6 py-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-outline-variant/10">
              <tr v-if="isLoadingUnit">
                <td colspan="3" class="px-6 py-12 text-center text-secondary">
                  <span class="animate-spin material-symbols-outlined text-4xl">sync</span>
                </td>
              </tr>
              <tr v-else-if="units.length === 0">
                <td colspan="3" class="px-6 py-12 text-center text-on-surface-variant">
                  <p>Satuan tidak ditemukan.</p>
                </td>
              </tr>
              <tr
                v-else
                v-for="u in units"
                :key="u.id"
                class="hover:bg-surface-container-low/50 transition-colors group"
              >
                <td class="px-6 py-4 font-label-sm text-secondary">{{ u.code }}</td>
                <td class="px-6 py-4 text-body-sm font-bold text-on-surface">{{ u.name }}</td>
                <td class="px-6 py-4 text-center space-x-2">
                  <button
                    @click="openEditUnitModal(u)"
                    class="p-2 text-secondary hover:bg-secondary/10 rounded-full transition-colors inline-flex items-center justify-center"
                    title="Edit Satuan"
                  >
                    <span class="material-symbols-outlined text-sm">edit</span>
                  </button>
                  <button
                    @click="confirmDeleteUnit(u)"
                    class="p-2 text-error hover:bg-error/10 rounded-full transition-colors inline-flex items-center justify-center"
                    title="Hapus Satuan"
                  >
                    <span class="material-symbols-outlined text-sm">delete</span>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modal Form Brand -->
    <div
      v-if="showBrandModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
    >
      <div
        class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl transform transition-all relative flex flex-col"
      >
        <button
          @click="closeBrandModal"
          class="absolute top-6 right-6 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1 transition-colors z-10"
        >
          <span class="material-symbols-outlined">close</span>
        </button>

        <h3 class="text-title-lg font-bold text-primary mb-6 flex items-center gap-2 flex-shrink-0">
          <span class="material-symbols-outlined">{{
            isEditBrandMode ? 'edit_square' : 'add_box'
          }}</span>
          {{ isEditBrandMode ? 'Edit Merek' : 'Tambah Merek Baru' }}
        </h3>

        <form @submit.prevent="submitBrandForm">
          <div class="space-y-4">
            <div>
              <label
                class="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider block"
                >Kode Merek *</label
              >
              <input
                v-model="brandForm.code"
                type="text"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-primary focus:outline-none"
                placeholder="Contoh: SONY"
              />
            </div>

            <div>
              <label
                class="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider block"
                >Nama Merek *</label
              >
              <input
                v-model="brandForm.name"
                type="text"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-primary focus:outline-none"
                placeholder="Contoh: Sony"
              />
            </div>

            <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-outline-variant/20">
              <button
                type="button"
                @click="closeBrandModal"
                class="px-4 py-2 text-label-md font-medium text-on-surface-variant hover:bg-surface-variant/20 rounded-lg transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                class="px-6 py-2 bg-primary text-on-primary rounded-lg font-label-md shadow-md hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all flex items-center gap-2 disabled:opacity-70"
                :disabled="isSubmittingBrand"
              >
                <span
                  v-if="isSubmittingBrand"
                  class="animate-spin material-symbols-outlined text-[18px]"
                  >sync</span
                >
                <span v-else class="material-symbols-outlined text-[18px]">check</span>
                Simpan
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal Form Unit -->
    <div
      v-if="showUnitModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
    >
      <div
        class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl transform transition-all relative flex flex-col"
      >
        <button
          @click="closeUnitModal"
          class="absolute top-6 right-6 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1 transition-colors z-10"
        >
          <span class="material-symbols-outlined">close</span>
        </button>

        <h3
          class="text-title-lg font-bold text-secondary mb-6 flex items-center gap-2 flex-shrink-0"
        >
          <span class="material-symbols-outlined">{{
            isEditUnitMode ? 'edit_square' : 'add_box'
          }}</span>
          {{ isEditUnitMode ? 'Edit Satuan' : 'Tambah Satuan Baru' }}
        </h3>

        <form @submit.prevent="submitUnitForm">
          <div class="space-y-4">
            <div>
              <label
                class="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider block"
                >Kode Satuan *</label
              >
              <input
                v-model="unitForm.code"
                type="text"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-secondary focus:outline-none"
                placeholder="Contoh: PCS"
              />
            </div>

            <div>
              <label
                class="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider block"
                >Nama Satuan *</label
              >
              <input
                v-model="unitForm.name"
                type="text"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-secondary focus:outline-none"
                placeholder="Contoh: Pieces"
              />
            </div>

            <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-outline-variant/20">
              <button
                type="button"
                @click="closeUnitModal"
                class="px-4 py-2 text-label-md font-medium text-on-surface-variant hover:bg-surface-variant/20 rounded-lg transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                class="px-6 py-2 bg-secondary text-on-secondary rounded-lg font-label-md shadow-md hover:bg-secondary-container hover:text-on-secondary-container active:scale-95 transition-all flex items-center gap-2 disabled:opacity-70"
                :disabled="isSubmittingUnit"
              >
                <span
                  v-if="isSubmittingUnit"
                  class="animate-spin material-symbols-outlined text-[18px]"
                  >sync</span
                >
                <span v-else class="material-symbols-outlined text-[18px]">check</span>
                Simpan
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DashboardLayout from '@frontend/layouts/DashboardLayout.vue'
import { inventoryApi } from '@frontend/api/inventory'
import type { Brand, Unit } from '@frontend/types/inventory'

// --- STATE BRAND ---
const isLoadingBrand = ref(true)
const isSubmittingBrand = ref(false)
const showBrandModal = ref(false)
const isEditBrandMode = ref(false)
const brands = ref<Brand[]>([])
const defaultBrandForm = () => ({ id: '', code: '', name: '' })
const brandForm = ref(defaultBrandForm())

// --- STATE UNIT ---
const isLoadingUnit = ref(true)
const isSubmittingUnit = ref(false)
const showUnitModal = ref(false)
const isEditUnitMode = ref(false)
const units = ref<Unit[]>([])
const defaultUnitForm = () => ({ id: '', code: '', name: '' })
const unitForm = ref(defaultUnitForm())

const loadBrands = async () => {
  isLoadingBrand.value = true
  try {
    const res = await inventoryApi.getBrands()
    brands.value = res
  } catch (error: any) {
    console.error('Gagal memuat merek:', error)
  } finally {
    isLoadingBrand.value = false
  }
}

const loadUnits = async () => {
  isLoadingUnit.value = true
  try {
    const res = await inventoryApi.getUnits()
    units.value = res
  } catch (error: any) {
    console.error('Gagal memuat satuan:', error)
  } finally {
    isLoadingUnit.value = false
  }
}

onMounted(() => {
  loadBrands()
  loadUnits()
})

// --- ACTIONS BRAND ---
const openAddBrandModal = () => {
  brandForm.value = defaultBrandForm()
  isEditBrandMode.value = false
  showBrandModal.value = true
}

const openEditBrandModal = (brand: Brand) => {
  brandForm.value = { ...brand }
  isEditBrandMode.value = true
  showBrandModal.value = true
}

const closeBrandModal = () => {
  showBrandModal.value = false
  brandForm.value = defaultBrandForm()
}

const confirmDeleteBrand = async (brand: Brand) => {
  if (confirm(`Apakah Anda yakin ingin menghapus merek ${brand.name}?`)) {
    try {
      await inventoryApi.deleteBrand(brand.id)
      await loadBrands()
    } catch (error: any) {
      alert('Gagal menghapus merek: ' + (error.response?.data?.message || error.message))
    }
  }
}

const submitBrandForm = async () => {
  isSubmittingBrand.value = true
  try {
    if (isEditBrandMode.value && brandForm.value.id) {
      await inventoryApi.updateBrand(brandForm.value.id, {
        code: brandForm.value.code,
        name: brandForm.value.name,
      })
    } else {
      await inventoryApi.createBrand({
        code: brandForm.value.code,
        name: brandForm.value.name,
      })
    }
    closeBrandModal()
    await loadBrands()
  } catch (error: any) {
    let errorMessage = error.response?.data?.message || error.message
    if (errorMessage.includes('brands_code_key')) {
      errorMessage = 'Kode merek sudah digunakan. Silakan gunakan kode lain.'
    }
    alert('Gagal menyimpan merek: ' + errorMessage)
  } finally {
    isSubmittingBrand.value = false
  }
}

// --- ACTIONS UNIT ---
const openAddUnitModal = () => {
  unitForm.value = defaultUnitForm()
  isEditUnitMode.value = false
  showUnitModal.value = true
}

const openEditUnitModal = (unit: Unit) => {
  unitForm.value = { ...unit }
  isEditUnitMode.value = true
  showUnitModal.value = true
}

const closeUnitModal = () => {
  showUnitModal.value = false
  unitForm.value = defaultUnitForm()
}

const confirmDeleteUnit = async (unit: Unit) => {
  if (confirm(`Apakah Anda yakin ingin menghapus satuan ${unit.name}?`)) {
    try {
      await inventoryApi.deleteUnit(unit.id)
      await loadUnits()
    } catch (error: any) {
      alert('Gagal menghapus satuan: ' + (error.response?.data?.message || error.message))
    }
  }
}

const submitUnitForm = async () => {
  isSubmittingUnit.value = true
  try {
    if (isEditUnitMode.value && unitForm.value.id) {
      await inventoryApi.updateUnit(unitForm.value.id, {
        code: unitForm.value.code,
        name: unitForm.value.name,
      })
    } else {
      await inventoryApi.createUnit({
        code: unitForm.value.code,
        name: unitForm.value.name,
      })
    }
    closeUnitModal()
    await loadUnits()
  } catch (error: any) {
    let errorMessage = error.response?.data?.message || error.message
    if (errorMessage.includes('units_code_key')) {
      errorMessage = 'Kode satuan sudah digunakan. Silakan gunakan kode lain.'
    }
    alert('Gagal menyimpan satuan: ' + errorMessage)
  } finally {
    isSubmittingUnit.value = false
  }
}
</script>
