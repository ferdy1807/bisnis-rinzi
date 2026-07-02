<template>
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary">Satuan Unit</h2>
        <p class="text-body-md text-on-surface-variant">Kelola daftar satuan pengukuran (unit) untuk inventaris Anda.</p>
      </div>
      <button @click="openCreateModal" class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2 rounded-xl shadow-md hover:bg-primary/90 transition-colors">
        <span class="material-symbols-outlined">add</span>
        <span class="text-label-md font-bold">Tambah Satuan</span>
      </button>
    </div>

    <div class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden flex flex-col min-h-[500px]">
      <div class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex justify-between items-center">
        <h3 class="font-headline-md text-primary flex items-center gap-2">
          <span class="material-symbols-outlined">straighten</span> Daftar Satuan
        </h3>
      </div>

      <div class="overflow-x-auto flex-1">
        <table class="w-full text-left">
          <thead>
            <tr class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold">
              <th class="px-6 py-4">Kode Satuan</th>
              <th class="px-6 py-4">Nama Satuan</th>
              <th class="px-6 py-4 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/20">
            <tr v-if="isLoading" class="text-center">
              <td colspan="3" class="py-8 text-on-surface-variant">Memuat data...</td>
            </tr>
            <tr v-else-if="units.length === 0" class="text-center">
              <td colspan="3" class="py-8 text-on-surface-variant">Belum ada satuan yang ditambahkan.</td>
            </tr>
            <tr v-for="unit in units" :key="unit.id" class="hover:bg-surface-container-lowest/50 transition-colors">
              <td class="px-6 py-4 font-mono text-sm">{{ unit.code }}</td>
              <td class="px-6 py-4 font-medium">{{ unit.name }}</td>
              <td class="px-6 py-4 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button @click="openEditModal(unit)" class="p-2 text-primary hover:bg-primary/10 rounded-full transition-colors" title="Edit">
                    <span class="material-symbols-outlined text-sm">edit</span>
                  </button>
                  <button @click="confirmDelete(unit)" class="p-2 text-error hover:bg-error/10 rounded-full transition-colors" title="Hapus">
                    <span class="material-symbols-outlined text-sm">delete</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl transform transition-all">
        <h3 class="text-title-lg font-bold text-primary mb-4">{{ isEditing ? 'Edit Satuan' : 'Tambah Satuan Baru' }}</h3>
        
        <form @submit.prevent="saveUnit" class="space-y-4">
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1">Kode Satuan</label>
            <input v-model="formData.code" type="text" required class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-primary focus:ring-1 focus:ring-primary outline-none bg-surface-container-lowest" placeholder="Contoh: PCS" />
          </div>
          <div>
            <label class="block text-label-sm text-on-surface-variant mb-1">Nama Satuan</label>
            <input v-model="formData.name" type="text" required class="w-full px-4 py-2 rounded-lg border border-outline-variant focus:border-primary focus:ring-1 focus:ring-primary outline-none bg-surface-container-lowest" placeholder="Contoh: Pieces" />
          </div>

          <div class="flex justify-end gap-3 mt-8">
            <button type="button" @click="closeModal" class="px-4 py-2 text-label-md font-medium text-on-surface-variant hover:bg-surface-variant/20 rounded-lg transition-colors">
              Batal
            </button>
            <button type="submit" :disabled="isSaving" class="px-4 py-2 text-label-md font-medium bg-primary text-on-primary rounded-lg shadow hover:bg-primary/90 transition-colors disabled:opacity-70">
              {{ isSaving ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { inventoryApi } from '@frontend/api/inventory';
import type { Unit } from '@frontend/types/inventory';

const units = ref<Unit[]>([]);
const isLoading = ref(true);

const showModal = ref(false);
const isEditing = ref(false);
const isSaving = ref(false);
const formData = ref({ id: '', code: '', name: '' });

const fetchUnits = async () => {
  isLoading.value = true;
  try {
    units.value = await inventoryApi.getUnits();
  } catch (error) {
    console.error("Gagal mengambil satuan", error);
    alert("Gagal memuat daftar satuan.");
  } finally {
    isLoading.value = false;
  }
};

const openCreateModal = () => {
  isEditing.value = false;
  formData.value = { id: '', code: '', name: '' };
  showModal.value = true;
};

const openEditModal = (unit: Unit) => {
  isEditing.value = true;
  formData.value = { id: unit.id, code: unit.code, name: unit.name };
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
};

const saveUnit = async () => {
  isSaving.value = true;
  try {
    if (isEditing.value) {
      await inventoryApi.updateUnit(formData.value.id, {
        code: formData.value.code,
        name: formData.value.name
      });
    } else {
      await inventoryApi.createUnit({
        code: formData.value.code,
        name: formData.value.name
      });
    }
    closeModal();
    await fetchUnits();
  } catch (error) {
    console.error("Gagal menyimpan satuan", error);
    alert("Gagal menyimpan data satuan. Periksa koneksi atau input Anda.");
  } finally {
    isSaving.value = false;
  }
};

const confirmDelete = async (unit: Unit) => {
  if (confirm(`Apakah Anda yakin ingin menghapus satuan "${unit.name}"?`)) {
    try {
      await inventoryApi.deleteUnit(unit.id);
      await fetchUnits();
    } catch (error) {
      console.error("Gagal menghapus satuan", error);
      alert("Gagal menghapus satuan. Satuan ini mungkin masih digunakan oleh produk.");
    }
  }
};

onMounted(() => {
  fetchUnits();
});
</script>
