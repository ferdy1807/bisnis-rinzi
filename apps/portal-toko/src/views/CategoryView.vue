<template>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary">Kategori Produk</h2>
        <p class="text-body-md text-on-surface-variant">
          Kelola daftar kategori produk untuk pengelompokan item.
        </p>
      </div>
      <div>
        <button
          @click="openAddModal"
          class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2 rounded-xl shadow-md hover:bg-primary-container hover:text-on-primary-container transition-colors"
        >
          <span class="material-symbols-outlined">add</span>
          <span class="font-bold text-label-md">Tambah Kategori</span>
        </button>
      </div>
    </div>

    <!-- Data Table -->
    <div
      class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]"
    >
      <div
        class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col sm:flex-row justify-between items-center gap-4"
      >
        <h3 class="font-headline-md text-primary flex items-center gap-2">
          <span class="material-symbols-outlined">category</span> Daftar Kategori
        </h3>
      </div>

      <div class="overflow-x-auto flex-1">
        <table class="w-full text-left">
          <thead>
            <tr
              class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold"
            >
              <th class="px-6 py-4">Kode Kategori</th>
              <th class="px-6 py-4">Nama Kategori</th>
              <th class="px-6 py-4 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/10">
            <tr v-if="isLoading">
              <td colspan="3" class="px-6 py-12 text-center text-primary">
                <span class="animate-spin material-symbols-outlined text-4xl">sync</span>
              </td>
            </tr>
            <tr v-else-if="categories.length === 0">
              <td colspan="3" class="px-6 py-12 text-center text-on-surface-variant">
                <p>Kategori tidak ditemukan.</p>
              </td>
            </tr>
            <tr
              v-else
              v-for="c in categories"
              :key="c.id"
              class="hover:bg-surface-container-low/50 transition-colors group"
            >
              <td class="px-6 py-4 font-label-sm text-primary">{{ c.code }}</td>
              <td class="px-6 py-4 text-body-sm font-bold text-on-surface">{{ c.name }}</td>
              <td class="px-6 py-4 text-center space-x-2">
                <button
                  @click="openEditModal(c)"
                  class="p-2 text-primary hover:bg-primary/10 rounded-full transition-colors inline-flex items-center justify-center"
                  title="Edit Kategori"
                >
                  <span class="material-symbols-outlined text-sm">edit</span>
                </button>
                <button
                  @click="confirmDelete(c)"
                  class="p-2 text-error hover:bg-error/10 rounded-full transition-colors inline-flex items-center justify-center"
                  title="Hapus Kategori"
                >
                  <span class="material-symbols-outlined text-sm">delete</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form (Add/Edit) -->
    <div
      v-if="showFormModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
    >
      <div
        class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl transform transition-all relative flex flex-col"
      >
        <button
          @click="closeModal"
          class="absolute top-6 right-6 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1 transition-colors z-10"
        >
          <span class="material-symbols-outlined">close</span>
        </button>

        <h3 class="text-title-lg font-bold text-primary mb-6 flex items-center gap-2 flex-shrink-0">
          <span class="material-symbols-outlined">{{
            isEditMode ? 'edit_square' : 'add_box'
          }}</span>
          {{ isEditMode ? 'Edit Kategori' : 'Tambah Kategori Baru' }}
        </h3>

        <form @submit.prevent="submitForm">
          <div class="space-y-4">
            <div>
              <label
                class="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider block"
                >Kode Kategori *</label
              >
              <input
                v-model="form.code"
                type="text"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-primary focus:outline-none"
                placeholder="Contoh: ELEK"
              />
            </div>

            <div>
              <label
                class="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider block"
                >Nama Kategori *</label
              >
              <input
                v-model="form.name"
                type="text"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:ring-1 focus:ring-primary focus:outline-none"
                placeholder="Contoh: Elektronik"
              />
            </div>

            <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-outline-variant/20">
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 text-label-md font-medium text-on-surface-variant hover:bg-surface-variant/20 rounded-lg transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                class="px-6 py-2 bg-primary text-on-primary rounded-lg font-label-md shadow-md hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all flex items-center gap-2 disabled:opacity-70"
                :disabled="isSubmitting"
              >
                <span v-if="isSubmitting" class="animate-spin material-symbols-outlined text-[18px]"
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
import { inventoryApi } from '@frontend/api/inventory'
import type { Category } from '@frontend/types/inventory'

const isLoading = ref(true)
const isSubmitting = ref(false)
const showFormModal = ref(false)
const isEditMode = ref(false)

const categories = ref<Category[]>([])

const defaultForm = () => ({ id: '', code: '', name: '' })
const form = ref(defaultForm())

const loadData = async () => {
  isLoading.value = true
  try {
    const res = await inventoryApi.getCategories()
    categories.value = res
  } catch (error: any) {
    console.error('Gagal memuat kategori:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadData()
})

const openAddModal = () => {
  form.value = defaultForm()
  isEditMode.value = false
  showFormModal.value = true
}

const openEditModal = (category: Category) => {
  form.value = { ...category }
  isEditMode.value = true
  showFormModal.value = true
}

const closeModal = () => {
  showFormModal.value = false
  form.value = defaultForm()
}

const confirmDelete = async (category: Category) => {
  if (confirm(`Apakah Anda yakin ingin menghapus kategori ${category.name}?`)) {
    try {
      await inventoryApi.deleteCategory(category.id)
      await loadData()
    } catch (error: any) {
      alert('Gagal menghapus kategori: ' + (error.response?.data?.message || error.message))
    }
  }
}

const submitForm = async () => {
  isSubmitting.value = true
  try {
    if (isEditMode.value && form.value.id) {
      await inventoryApi.updateCategory(form.value.id, {
        code: form.value.code,
        name: form.value.name,
      })
    } else {
      await inventoryApi.createCategory({
        code: form.value.code,
        name: form.value.name,
      })
    }
    closeModal()
    await loadData()
  } catch (error: any) {
    alert('Gagal menyimpan kategori: ' + (error.response?.data?.message || error.message))
  } finally {
    isSubmitting.value = false
  }
}
</script>
