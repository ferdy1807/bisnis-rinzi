<template>
  <Dialog 
    :visible="visible" 
    @update:visible="$emit('update:visible', $event)"
    modal 
    :style="{ width: '48rem' }" 
    class="font-sans"
    :pt="{
      root: { class: 'rounded-3xl shadow-2xl border-none overflow-hidden bg-white' },
      mask: { class: 'backdrop-blur-sm bg-slate-900/40' },
      header: { class: 'bg-[#f8f9fc] border-b border-[#e2e8f0] px-6 py-4' },
      content: { class: 'px-0 py-0' }
    }"
    @hide="closeModal"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="w-11 h-11 rounded-full bg-emerald-100 flex items-center justify-center">
          <i class="pi pi-box text-xl text-emerald-600"></i>
        </div>
        <div>
          <h2 class="text-lg font-black text-[#1a1c20] leading-tight">Pengelolaan Barang Titipan</h2>
          <p class="text-[11px] text-[#58627c] font-semibold mt-0.5">Pendataan serah terima fisik barang dari pengantin</p>
        </div>
      </div>
    </template>

    <div v-if="reservation" class="flex flex-col bg-slate-50/50">
      
      <!-- Info Banner (Tanda Terima) -->
      <div class="bg-gradient-to-r from-[#254582] to-[#1a3668] p-6 text-white flex justify-between items-center shadow-inner">
        <div>
          <span class="text-[10px] font-bold text-[#8ba6d9] uppercase tracking-widest block mb-1">Klien Penyewa</span>
          <strong class="text-2xl font-black block leading-none tracking-tight">{{ parseCustomer(reservation.customer_snapshot_id).name }}</strong>
          <span class="text-xs text-[#a3bced] font-mono mt-2.5 inline-block bg-white/10 px-2 py-0.5 rounded-md border border-white/10"><i class="pi pi-hashtag text-[10px]"></i> {{ reservation.invoice_number }}</span>
        </div>
      </div>

      <div class="p-6 space-y-6">
        
        <!-- DAFTAR BARANG YANG SUDAH DITERIMA -->
        <div class="space-y-3">
          <h3 class="text-sm font-black text-[#1a1c20] uppercase tracking-wider flex items-center gap-2 border-b border-gray-200 pb-2">
            <span class="w-7 h-7 rounded-lg bg-blue-100 text-blue-700 flex items-center justify-center"><i class="pi pi-list text-xs"></i></span>
            Riwayat Barang Masuk
          </h3>
          
          <div class="bg-white border border-gray-200/80 rounded-2xl overflow-hidden shadow-xs">
            <div v-if="!reservation.contents || reservation.contents.length === 0" class="p-3 m-3 bg-amber-50 border border-amber-200 rounded-xl flex items-center gap-3">
              <div class="w-8 h-8 bg-amber-100 text-amber-600 rounded-lg flex items-center justify-center shrink-0">
                <i class="pi pi-info-circle text-sm"></i>
              </div>
              <div>
                <h4 class="text-xs font-bold text-amber-900">Belum Ada Barang Tercatat</h4>
                <p class="text-[10px] text-amber-700/80">Riwayat masih kosong. Silakan tambahkan melalui formulir di bawah.</p>
              </div>
            </div>
            <DataTable v-else :value="reservation.contents" class="p-datatable-sm text-xs border-none">
              <Column field="item_name" header="Nama/Tipe" class="font-bold text-gray-800 w-[20%]"></Column>
              <Column field="quantity" header="Qty" class="w-[8%] text-center text-emerald-700 font-bold bg-emerald-50/50"></Column>
              <Column field="description" header="Rincian / Keterangan Lengkap"></Column>
              <Column field="condition_notes" header="Kondisi Tiba" class="italic text-gray-500 text-[11px]"></Column>
              <Column header="Waktu" class="text-right w-[15%] font-mono text-[10px] text-gray-400">
                <template #body="{ data }">
                  {{ new Date(data.created_at).toLocaleDateString('id-ID', { day:'numeric', month:'short' }) }}
                </template>
              </Column>
            </DataTable>
          </div>
        </div>

        <!-- FORM TAMBAH BARANG BARU -->
        <div class="space-y-4 bg-emerald-50/50 p-5 rounded-2xl border border-emerald-100 shadow-inner">
          <h3 class="text-sm font-black text-emerald-800 uppercase tracking-wider flex items-center gap-2">
            <span class="w-7 h-7 rounded-lg bg-emerald-200 text-emerald-800 flex items-center justify-center"><i class="pi pi-plus text-xs"></i></span>
            Catat Barang Baru
          </h3>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="md:col-span-2 flex flex-col gap-1.5">
              <label class="text-xs font-bold text-emerald-900">Nama / Tipe Barang <span class="text-red-500">*</span></label>
              <InputText 
                v-model="form.item_name" 
                placeholder="Cth: Perhiasan Emas / Kado Kosmetik" 
                class="rounded-xl border-emerald-200 w-full text-sm px-4 py-2.5 shadow-2xs focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all bg-white" 
              />
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-bold text-emerald-900">Jumlah <span class="text-red-500">*</span></label>
              <InputNumber 
                v-model="form.quantity" 
                placeholder="1" 
                class="w-full"
                inputClass="rounded-xl border-emerald-200 w-full text-sm px-4 py-2.5 shadow-2xs focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all bg-white"
                :min="1"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-emerald-900">Keterangan / Rincian Isi <span class="text-red-500">*</span></label>
            <InputText 
              v-model="form.item_description" 
              placeholder="Cth: 1 Set Kalung Emas & 2 Cincin (Tanpa Surat)" 
              class="rounded-xl border-emerald-200 w-full text-sm px-4 py-2.5 shadow-2xs focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all bg-white" 
            />
          </div>
          
          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-emerald-900">Kondisi Tiba (Cacat/Lecet/dll)</label>
            <Textarea 
              v-model="form.condition_notes" 
              rows="2" 
              placeholder="Cth: Kotak kayu cincin ujungnya gompal sedikit di pojok kiri..." 
              class="rounded-xl border-emerald-200 w-full text-sm px-4 py-2.5 shadow-2xs focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all bg-white" 
            />
          </div>

          <div class="pt-2">
            <Button 
              label="Simpan Barang" 
              icon="pi pi-check" 
              class="w-full bg-emerald-600 hover:bg-emerald-700 text-white border-none font-bold text-sm py-3.5 rounded-xl shadow-md transition-all" 
              :loading="submitting"
              :disabled="!form.item_name || !form.item_description || !form.quantity"
              @click="submitDeposit" 
            />
          </div>
        </div>

      </div>

      <!-- Tombol Aksi Bawah -->
      <div class="p-5 bg-white border-t border-gray-200 flex justify-between items-center rounded-b-3xl">
        <span class="text-xs text-gray-500 italic flex items-center gap-1"><i class="pi pi-check-circle text-emerald-500"></i> Pastikan rincian sudah sesuai</span>
        <Button label="Selesai & Tutup Formulir" icon="pi pi-arrow-right" iconPos="right" class="bg-[#1a1c20] hover:bg-[#2c2f36] text-white border-none font-bold text-xs py-2.5 px-6 rounded-xl shadow-lg transition-all" @click="closeModal" />
      </div>

    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRentalStore } from '@frontend/stores/rental';
import { rentalApi } from '@frontend/api/rental';
import { useToast } from 'primevue/usetoast';
import type { RentalReservation } from '@frontend/types/rental';

import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Textarea from 'primevue/textarea';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';

const props = defineProps<{
  visible: boolean;
  reservation: RentalReservation | null;
}>();

const emit = defineEmits(['update:visible', 'saved']);

const rentalStore = useRentalStore();
const toast = useToast();

const submitting = ref(false);
const form = ref({
  item_name: '',
  quantity: 1,
  item_description: '',
  condition_notes: ''
});

// Helper parsing nama
const parseCustomer = (str?: string) => {
  if (!str) return { name: 'Tanpa Nama', phone: '-' };
  const m = str.match(/^(.*?)\s*\((.*?)\)$/);
  return m ? { name: m[1]?.trim() || str, phone: m[2]?.trim() || '-' } : { name: str, phone: '-' };
};

const closeModal = () => {
  emit('update:visible', false);
};

const submitDeposit = async () => {
  if (!props.reservation) return;
  submitting.value = true;
  try {
    await rentalApi.saveDepositItems(props.reservation.id, {
      item_name: form.value.item_name,
      description: form.value.item_description,
      quantity: form.value.quantity,
      condition_notes: form.value.condition_notes
    });
    toast.add({ severity: 'success', summary: 'Sukses', detail: 'Barang titipan berhasil dicatat', life: 3000 });
    
    // Update local content list to react immediately
    const updated = await rentalApi.getReservationDetail(props.reservation.id);
    if (updated.data && updated.data.contents) {
       props.reservation.contents = updated.data.contents;
    }

    form.value = { item_name: '', quantity: 1, item_description: '', condition_notes: '' };
    
    emit('saved');
    
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Gagal Menyimpan', detail: error.response?.data?.message || 'Terjadi kesalahan sistem', life: 3000 });
  } finally {
    submitting.value = false;
  }
};
</script>
