<template>
  <div class="p-1 animate-fade-in space-y-6">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-outline-variant/20 pb-6">
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">Buku Besar Global & Jurnal Umum</h2>
        <p class="text-body-md text-on-surface-variant mt-1">
          Otoritas Pemeriksaan Akuntansi Owner: Transparansi entri berpasangan (Double-Entry Log), audit trail jurnal penyesuaian, dan pelaporan keuangan konsolidasi.
        </p>
      </div>
      <div class="flex flex-wrap items-center gap-3">
        <button
          @click="openManualJournalModal"
          class="flex items-center gap-2 bg-primary text-on-primary px-4 py-2.5 rounded-xl shadow-md hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all text-label-md font-bold">
          <span class="material-symbols-outlined text-md">edit_note</span>
          Buat Jurnal Manual
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[400px] text-primary">
      <span class="animate-spin material-symbols-outlined text-4xl mb-2">sync</span>
      <p class="text-label-md font-bold tracking-wider">Memetakan Jejak Transaksi Jurnal Berpasangan...</p>
    </div>

    <div v-else class="space-y-6">
      <div class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden">
        <div class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <h3 class="font-title-medium text-primary font-bold flex items-center gap-2">
            <span class="material-symbols-outlined text-md">menu_book</span> Transaksi Jurnal Berhitung (journal_entries)
          </h3>
          <div class="flex gap-2">
            <span class="text-xs bg-secondary/10 text-secondary px-3 py-1 rounded-full font-bold font-mono border border-secondary/10">
              Target Skema: finance_db
            </span>
          </div>
        </div>
        
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse border border-outline-variant/30 table-fixed min-w-[900px]">
            <thead>
              <tr class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold border-b border-outline-variant/30 divide-x divide-outline-variant/30">
                <th class="px-4 py-3.5 w-[15%]">Tanggal / No Ref</th>
                <th class="px-4 py-3.5 w-[30%]">Keterangan Narasi Jurnal</th>
                <th class="px-4 py-3.5 w-[25%]">Kode Akun Pokok (COA)</th>
                <th class="px-4 py-3.5 w-[10%] text-right">Debet</th>
                <th class="px-4 py-3.5 w-[10%] text-right">Kredit</th>
                <th class="px-4 py-3.5 w-[10%] text-center">Status</th>
              </tr>
            </thead>
            <tbody class="text-body-sm divide-y divide-outline-variant/30 font-mono text-xs">
              <template v-if="journalEntries.length === 0">
                <tr>
                  <td colspan="6" class="px-6 py-8 text-center font-sans text-on-surface-variant">Belum ada ayat jurnal akuntansi terdaftar.</td>
                </tr>
              </template>
              
              <template v-else v-for="entry in journalEntries" :key="entry.id">
                <tr 
                  v-for="(detail, index) in (entry.details || [])" 
                  :key="detail.id"
                  class="divide-x divide-outline-variant/20 hover:bg-surface-container-low/10 transition-colors"
                  :class="{'border-t-2 border-outline-variant/40 bg-surface-container-lowest/40': index === 0}"
                >
                  <td 
                    v-if="index === 0" 
                    :rowspan="(entry.details || []).length" 
                    class="px-4 py-3 align-top bg-surface font-sans"
                  >
                    <div class="text-on-surface font-bold text-sm">{{ formatDate(entry.entry_date) }}</div>
                    <div class="text-primary font-normal text-xs mt-1 break-all tracking-tight leading-normal">
                      {{ entry.reference_number }}
                    </div>
                  </td>

                  <td 
                    v-if="index === 0" 
                    :rowspan="(entry.details || []).length" 
                    class="px-4 py-3 align-top bg-surface font-sans text-sm text-on-surface-variant break-words whitespace-pre-wrap leading-relaxed"
                  >
                    {{ entry.narration || 'Tidak ada catatan narasi transaksi.' }}
                  </td>

                  <td class="px-4 py-3 align-middle text-on-surface text-xs font-medium">
                    <div class="flex items-center gap-2">
                      <span 
                        class="w-1.5 h-1.5 rounded-full flex-shrink-0" 
                        :class="Number(detail.debit_amount) > 0 ? 'bg-primary' : 'bg-secondary ml-4'"
                      ></span>
                      <span class="truncate" :title="getAccountMetadata(detail.account_id)">
                        {{ getAccountMetadata(detail.account_id) }}
                      </span>
                    </div>
                  </td>

                  <td class="px-4 py-3 align-middle text-right font-bold text-on-surface tabular-nums">
                    {{ Number(detail.debit_amount) > 0 ? 'Rp ' + formatCurrency(detail.debit_amount) : '-' }}
                  </td>

                  <td class="px-4 py-3 align-middle text-right font-bold text-on-surface tabular-nums">
                    {{ Number(detail.credit_amount) > 0 ? 'Rp ' + formatCurrency(detail.credit_amount) : '-' }}
                  </td>

                  <td 
                    v-if="index === 0" 
                    :rowspan="(entry.details || []).length" 
                    class="px-4 py-3 align-middle text-center bg-surface font-sans"
                  >
                    <span 
                      :class="entry.is_posted ? 'bg-success/10 text-success border border-success/20' : 'bg-warning/10 text-warning border border-warning/20'" 
                      class="px-2 py-0.5 rounded text-[9px] font-bold tracking-widest uppercase inline-block"
                    >
                      {{ entry.is_posted ? 'POSTED' : 'UNPOSTED' }}
                    </span>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div v-if="showJournalModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="bg-surface rounded-2xl w-full max-w-2xl p-6 shadow-xl relative border border-outline-variant/20 max-h-[90vh] overflow-y-auto">
        <button @click="showJournalModal = false" class="absolute top-6 right-6 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1">
          <span class="material-symbols-outlined">close</span>
        </button>
        
        <h3 class="text-title-lg font-bold text-primary mb-2 flex items-center gap-1.5">
          <span class="material-symbols-outlined">edit_note</span> Entri Jurnal Umum Penyesuaian Manual
        </h3>
        <p class="text-xs text-on-surface-variant mb-6">Otoritas Owner: Penyusunan ayat jurnal ganda penyeimbang kas lintas-divisi.</p>
        
        <form @submit.prevent="submitManualJournalEntry" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block">Nomor Referensi Dokumen *</label>
              <input v-model="journalForm.reference_number" type="text" required placeholder="Contoh: JVM-2026-001" class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:outline-none" />
            </div>
            <div>
              <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block">Tanggal Entri Transaksi *</label>
              <input v-model="journalForm.entry_date" type="date" required class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:outline-none" />
            </div>
          </div>
          <div>
            <label class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block">Narasi Keterangan Jurnal *</label>
            <textarea v-model="journalForm.narration" required rows="2" placeholder="Contoh: Penyesuaian koreksi selisih kas fisik akibat mika pecah" class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:outline-none resize-none"></textarea>
          </div>

          <div class="border-t border-outline-variant/20 pt-4 space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-xs font-black uppercase tracking-wider text-primary">Komponen Alokasi Baris Jurnal (Details)</span>
              <button type="button" @click="addDetailRow" class="text-xs font-bold text-primary flex items-center gap-1">
                <span class="material-symbols-outlined text-xs">add</span> Tambah Baris Akun
              </button>
            </div>

            <div v-for="(row, index) in journalForm.details" :key="index" class="grid grid-cols-1 sm:grid-cols-4 gap-3 bg-surface-container-low p-3 rounded-xl border border-outline-variant/10 relative">
              <div class="sm:col-span-2">
                <label class="text-[9px] font-bold text-on-surface-variant uppercase">Pilih Akun COA *</label>
                <select v-model="row.account_id" required class="w-full mt-0.5 bg-surface border border-outline-variant/50 rounded-lg px-2 py-1.5 text-xs focus:outline-none">
                  <option value="" disabled>Pilih Akun Pokok</option>
                  <option v-for="acc in accountsList" :key="acc.id" :value="acc.id">[{{ acc.account_code }}] {{ acc.account_name }}</option>
                </select>
              </div>
              <div>
                <label class="text-[9px] font-bold text-on-surface-variant uppercase">Debet (Rp)</label>
                <input v-model.number="row.debit_amount" type="number" min="0" step="any" class="w-full mt-0.5 bg-surface border border-outline-variant/50 rounded-lg px-2 py-1.5 text-xs font-mono focus:outline-none" />
              </div>
              <div class="flex items-end gap-2">
                <div class="flex-1">
                  <label class="text-[9px] font-bold text-on-surface-variant uppercase">Kredit (Rp)</label>
                  <input v-model.number="row.credit_amount" type="number" min="0" step="any" class="w-full mt-0.5 bg-surface border border-outline-variant/50 rounded-lg px-2 py-1.5 text-xs font-mono focus:outline-none" />
                </div>
                <button v-if="journalForm.details.length > 2" type="button" @click="removeDetailRow(index)" class="text-error mb-1 p-1 hover:bg-error/10 rounded">
                  <span class="material-symbols-outlined text-sm">delete</span>
                </button>
              </div>
            </div>
          </div>

          <div class="bg-surface-container-high p-3 rounded-xl flex items-center justify-between text-xs font-mono">
            <div>Total Debet: <span class="font-bold">Rp {{ formatCurrency(totalDebetCalculated) }}</span></div>
            <div>Total Kredit: <span class="font-bold">Rp {{ formatCurrency(totalKreditCalculated) }}</span></div>
            <div :class="isJournalBalanced ? 'text-teal-600' : 'text-error'" class="font-bold uppercase tracking-widest">
              {{ isJournalBalanced ? 'Balanced' : 'Unbalanced' }}
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-4 border-t border-outline-variant/20">
            <button type="button" :disabled="isSubmitting" @click="showJournalModal = false" class="px-4 py-2 text-label-md font-bold text-on-surface-variant rounded-xl hover:bg-surface-variant/20 transition-colors">Batal</button>
            <button type="submit" :disabled="!isJournalBalanced || isSubmitting" class="px-6 py-2 bg-primary text-on-primary rounded-xl text-label-md font-bold shadow-md disabled:opacity-50 transition-all">
              <span v-if="isSubmitting" class="animate-spin material-symbols-outlined text-sm mr-1">sync</span>
              Posting Jurnal Penyesuaian
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { financeApi } from '@frontend/api/finance';
import type { JournalEntry, ChartOfAccount } from '@frontend/types/finance';

const isLoading = ref(true);
const isSubmitting = ref(false);
const showJournalModal = ref(false);

const journalEntries = ref<JournalEntry[]>([]);
const accountsList = ref<ChartOfAccount[]>([]);
const accountsMap = ref<Record<string, ChartOfAccount>>({});

const getTodayDateString = (): string => {
  const today = new Date();
  const year = today.getFullYear();
  const month = String(today.getMonth() + 1).padStart(2, '0');
  const day = String(today.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const journalForm = ref({
  journal_id: '11de91a0-da11-4c91-b621-fe7369b168c7',
  reference_number: '',
  entry_date: getTodayDateString(),
  narration: '',
  details: [
    { account_id: '', debit_amount: 0, credit_amount: 0 },
    { account_id: '', debit_amount: 0, credit_amount: 0 }
  ]
});

const loadLedgerArchitecture = async () => {
  isLoading.value = true;
  try {
    const [entriesRes, accountsRes] = await Promise.all([
      financeApi.getJournalEntries(),
      financeApi.getAccounts()
    ]);

    let cleanAccounts = [];
    if (Array.isArray(accountsRes)) {
      cleanAccounts = accountsRes;
    } else if (accountsRes && Array.isArray((accountsRes as any).data)) {
      cleanAccounts = (accountsRes as any).data;
    } else if (accountsRes && (accountsRes as any).data && Array.isArray((accountsRes as any).data.data)) {
      cleanAccounts = (accountsRes as any).data.data;
    }
    
    accountsList.value = cleanAccounts;
    
    accountsMap.value = {};
    cleanAccounts.forEach((acc: ChartOfAccount) => {
      accountsMap.value[acc.id] = acc;
    });

    if (Array.isArray(entriesRes)) {
      journalEntries.value = entriesRes;
    } else if (entriesRes && Array.isArray((entriesRes as any).data)) {
      journalEntries.value = (entriesRes as any).data;
    } else {
      journalEntries.value = [];
    }

  } catch (err) {
    console.error('Gagal mengambil manifes pembukuan buku besar dari API gateway:', err);
    journalEntries.value = [];
    accountsList.value = [];
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadLedgerArchitecture();
});

const getAccountMetadata = (accountId: string): string => {
  const acc = accountsMap.value[accountId];
  return acc ? `[${acc.account_code}] ${acc.account_name}` : 'Memuat Referensi COA...';
};

const totalDebetCalculated = computed(() => journalForm.value.details.reduce((sum, item) => sum + (Number(item.debit_amount) || 0), 0));
const totalKreditCalculated = computed(() => journalForm.value.details.reduce((sum, item) => sum + (Number(item.credit_amount) || 0), 0));
const isJournalBalanced = computed(() => {
  const debet = totalDebetCalculated.value;
  const kredit = totalKreditCalculated.value;
  return debet > 0 && Math.abs(debet - kredit) < 0.001;
});

const openManualJournalModal = () => {
  journalForm.value = {
    journal_id: '11de91a0-da11-4c91-b621-fe7369b168c7',
    reference_number: '',
    entry_date: getTodayDateString(),
    narration: '',
    details: [
      { account_id: '', debit_amount: 0, credit_amount: 0 },
      { account_id: '', debit_amount: 0, credit_amount: 0 }
    ]
  };
  showJournalModal.value = true;
};

const addDetailRow = () => {
  journalForm.value.details.push({ account_id: '', debit_amount: 0, credit_amount: 0 });
};

const removeDetailRow = (index: number) => {
  journalForm.value.details.splice(index, 1);
};

const submitManualJournalEntry = async () => {
  if (!isJournalBalanced.value || isSubmitting.value) return;
  isSubmitting.value = true;
  try {
    const sanitizedDetails = journalForm.value.details.map(item => ({
      account_id: item.account_id,
      debit_amount: Number(item.debit_amount) || 0,
      credit_amount: Number(item.credit_amount) || 0
    }));

    const payload = {
      ...journalForm.value,
      details: sanitizedDetails
    };

    await financeApi.createManualJournal(payload);
    showJournalModal.value = false;
    await loadLedgerArchitecture();
  } catch (err: any) {
    const apiError = err.response?.data?.message || err.message || 'Koneksi API bermasalah';
    alert('Gagal mengeksekusi posting ayat jurnal manual: ' + apiError);
  } finally {
    isSubmitting.value = false;
  }
};

const formatCurrency = (val: number | string): string => {
  const numericVal = typeof val === 'string' ? Number(val) : val;
  if (isNaN(numericVal)) return '0';
  return new Intl.NumberFormat('id-ID', { minimumFractionDigits: 0 }).format(numericVal);
};

const formatDate = (dateStr: string | Date): string => {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  if (isNaN(d.getTime())) return String(dateStr);
  return d.toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' });
};
</script>