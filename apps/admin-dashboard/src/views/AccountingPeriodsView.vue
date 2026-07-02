<template>
  <div class="p-1 animate-fade-in space-y-6">
    <div
      class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-outline-variant/20 pb-6"
    >
      <div>
        <h2 class="font-headline-lg text-headline-lg text-primary font-bold tracking-tight">
          Manajemen Periode Akuntansi & Lock System
        </h2>
        <p class="text-body-md text-on-surface-variant mt-1">
          Otoritas Finansial Owner: Pembukaan kalender fiskal pembukuan, pembekuan ayat jurnal, dan
          pencegahan manipulasi data transaksi masa lalu.
        </p>
      </div>
      <div>
        <button
          @click="openCreatePeriodModal"
          class="flex items-center gap-2 bg-primary text-on-primary px-5 py-2.5 rounded-xl shadow-md hover:bg-primary-container hover:text-on-primary-container active:scale-95 transition-all text-label-md font-bold"
        >
          <span class="material-symbols-outlined text-md">calendar_add_on</span>
          Buka Periode Baru
        </button>
      </div>
    </div>

    <div
      v-if="isLoading"
      class="flex flex-col items-center justify-center min-h-[350px] text-primary"
    >
      <span class="animate-spin material-symbols-outlined text-4xl mb-2">sync</span>
      <p class="text-label-md font-bold tracking-wider">
        Menyelaraskan Kalender Fiskal Korporat...
      </p>
    </div>

    <div v-else class="space-y-6">
      <div
        class="bg-surface rounded-2xl border border-outline-variant/30 shadow-sm overflow-hidden"
      >
        <div
          class="p-5 border-b border-outline-variant/20 bg-surface-container-lowest flex items-center justify-between"
        >
          <h3 class="font-title-medium text-primary font-bold flex items-center gap-2">
            <span class="material-symbols-outlined text-md">date_range</span> Kalender Siklus
            Pembukuan (accounting_periods)
          </h3>
          <span
            class="text-[10px] font-mono bg-surface-container-high px-3 py-1 rounded-full font-bold text-on-surface-variant"
          >
            Target Schema: finance_db
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr
                class="bg-surface-container-low text-on-surface-variant uppercase text-[10px] tracking-widest font-bold border-b border-outline-variant/20"
              >
                <th class="px-6 py-4">Nama Periode / Tahun</th>
                <th class="px-6 py-4">Tanggal Mulai</th>
                <th class="px-6 py-4">Tanggal Selesai</th>
                <th class="px-6 py-4 text-center">Status Buku</th>
                <th class="px-6 py-4 text-center">Otoritas Lock Jurnal</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-outline-variant/10 text-body-sm font-mono text-xs">
              <tr v-if="periodsList.length === 0">
                <td colspan="5" class="px-6 py-8 text-center font-sans text-on-surface-variant">
                  Belum ada siklus periode akuntansi yang terdaftar.
                </td>
              </tr>
              <tr
                v-for="period in periodsList"
                :key="period.id"
                class="hover:bg-surface-container-low/30 transition-colors"
              >
                <td class="px-6 py-4 font-bold text-on-surface text-sm font-sans">
                  {{ period.name }}
                </td>
                <td class="px-6 py-4 text-on-surface-variant">{{ period.start_date }}</td>
                <td class="px-6 py-4 text-on-surface-variant">{{ period.end_date }}</td>
                <td class="px-6 py-4 text-center font-sans">
                  <span
                    :class="
                      period.is_closed ? 'bg-error/10 text-error' : 'bg-success/10 text-success'
                    "
                    class="px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider"
                  >
                    {{ period.is_closed ? 'Closed / Permanen' : 'Active / Open' }}
                  </span>
                </td>
                <td class="px-6 py-4 text-center font-sans">
                  <button
                    v-if="!period.is_closed"
                    @click="openLockAuthorizationModal(period)"
                    class="inline-flex items-center gap-1 bg-surface border border-error/40 text-error px-3 py-1.5 rounded-lg hover:bg-error/5 active:scale-95 transition-all text-label-sm font-bold"
                  >
                    <span class="material-symbols-outlined text-sm">lock</span>
                    Tutup Buku & Bekukan
                  </button>
                  <span
                    v-else
                    class="text-outline text-xs italic font-mono flex items-center justify-center gap-1"
                    ><span class="material-symbols-outlined text-xs">verified</span> Archived</span
                  >
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
    >
      <div
        class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl relative border border-outline-variant/20"
      >
        <button
          @click="showCreateModal = false"
          class="absolute top-6 right-6 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1"
        >
          <span class="material-symbols-outlined">close</span>
        </button>
        <h3 class="text-title-lg font-bold text-primary mb-6 flex items-center gap-1.5">
          <span class="material-symbols-outlined">calendar_add_on</span> Inisialisasi Siklus Buku
          Baru
        </h3>

        <form @submit.prevent="executePeriodCreation" class="space-y-4">
          <div>
            <label
              class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block"
              >Nama Penanda Periode *</label
            >
            <input
              v-model="createForm.name"
              type="text"
              required
              placeholder="Contoh: Periode Juli 2026"
              class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:outline-none"
            />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label
                class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block"
                >Tanggal Mulai *</label
              >
              <input
                v-model="createForm.start_date"
                type="date"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:outline-none"
              />
            </div>
            <div>
              <label
                class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block"
                >Tanggal Selesai *</label
              >
              <input
                v-model="createForm.end_date"
                type="date"
                required
                class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:outline-none"
              />
            </div>
          </div>
          <div class="flex justify-end gap-2 pt-4 border-t border-outline-variant/20">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 text-label-md font-bold text-on-surface-variant"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="px-5 py-2 bg-primary text-on-primary rounded-lg text-label-md font-bold shadow-md"
            >
              Buka Kalender Fiskal
            </button>
          </div>
        </form>
      </div>
    </div>

    <div
      v-if="showLockModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
    >
      <div
        class="bg-surface rounded-2xl w-full max-w-md p-6 shadow-xl relative border border-outline-variant/20"
      >
        <button
          @click="showLockModal = false"
          class="absolute top-6 right-6 text-on-surface-variant hover:bg-surface-variant/20 rounded-full p-1"
        >
          <span class="material-symbols-outlined">close</span>
        </button>
        <h3 class="text-title-lg font-bold text-error mb-2 flex items-center gap-1.5">
          <span class="material-symbols-outlined">lock</span> Pembekuan Mutasi Jurnal Absolut
        </h3>
        <p class="text-xs text-on-surface-variant font-bold uppercase mt-1 mb-6">
          Target: {{ selectedPeriod?.name }}
        </p>

        <form @submit.prevent="executePeriodLockFreeze" class="space-y-4">
          <div
            class="p-3.5 bg-error/5 border border-error/20 rounded-xl text-xs text-error leading-relaxed font-sans"
          >
            <span class="font-black uppercase tracking-wider block mb-1"
              >⚠️ Peringatan Kritis Keuangan</span
            >
            Tindakan ini akan mengunci seluruh entri jurnal dari unit kasir ritel maupun denda sewa
            secara permanen. Setelah dikunci, data transaksi pada rentang tanggal ini tidak akan
            bisa dimodifikasi kembali demi integritas audit perpajakan.
          </div>
          <div>
            <label
              class="text-[10px] font-black text-on-surface-variant uppercase tracking-wider block"
              >Sebab Kebijakan Penutupan Buku *</label
            >
            <textarea
              v-model="lockForm.reason"
              required
              rows="3"
              placeholder="Sebutkan Berita Acara, contoh: Rekonsiliasi semesteran berimbang mutlak..."
              class="w-full mt-1 bg-surface-container border border-outline-variant/50 rounded-lg px-3 py-2 text-body-sm focus:outline-none"
            ></textarea>
          </div>
          <div class="flex justify-end gap-2 pt-4 border-t border-outline-variant/20">
            <button
              type="button"
              @click="showLockModal = false"
              class="px-4 py-2 text-label-md font-bold text-on-surface-variant"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="px-6 py-2 bg-error text-on-error rounded-lg text-label-md font-bold shadow-md"
            >
              Kunci & Amankan Buku Besar
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { financeApi } from '@frontend/api/finance'

interface AccountingPeriodItem {
  id: string
  name: string
  start_date: string
  end_date: string
  is_closed: boolean
}

const isLoading = ref(true)
const isSubmitting = ref(false)
const showCreateModal = ref(false)
const showLockModal = ref(false)

const periodsList = ref<AccountingPeriodItem[]>([])
const selectedPeriod = ref<AccountingPeriodItem | null>(null)

const createForm = ref({ name: '', start_date: '2026-07-01', end_date: '2026-07-31' })
const lockForm = ref({ reason: '' })

const loadFiscalsManifest = async () => {
  isLoading.value = true
  try {
    // Memanggil fungsi eksternal dari financeApi pembaca rute GET /api/finance/periods
    const response = await financeApi.getPeriods()
    periodsList.value = Array.isArray(response) ? response : (response as any).data || []
  } catch (err) {
    console.error('Gagal mengambil kalender fiskal akuntansi:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadFiscalsManifest()
})

const openCreatePeriodModal = () => {
  createForm.value = { name: '', start_date: '2026-07-01', end_date: '2026-07-31' }
  showCreateModal.value = true
}

const executePeriodCreation = async () => {
  isSubmitting.value = true
  try {
    // Mengubah format YYYY-MM-DD menjadi RFC3339 yang dipahami oleh golang time.Time
    const payload = {
      name: createForm.value.name,
      start_date: createForm.value.start_date + 'T00:00:00Z',
      end_date: createForm.value.end_date + 'T23:59:59Z',
    }
    await financeApi.createPeriod(payload)
    showCreateModal.value = false
    await loadFiscalsManifest()
  } catch (err) {
    alert('Gagal mendefinisikan kalender pembukuan fiskal baru.')
  } finally {
    isSubmitting.value = false
  }
}

const openLockAuthorizationModal = (period: AccountingPeriodItem) => {
  selectedPeriod.value = period
  lockForm.value.reason = ''
  showLockModal.value = true
}

const executePeriodLockFreeze = async () => {
  if (!selectedPeriod.value) return
  isSubmitting.value = true
  try {
    // Memanggil metode lockPeriod dari instance financeApi yang mengarah ke POST /api/finance/period-locks
    await financeApi.lockPeriod(selectedPeriod.value.id, lockForm.value.reason)
    showLockModal.value = false
    await loadFiscalsManifest()
  } catch (err) {
    alert('Otorisasi pembekuan siklus buku gagal dilakukan.')
  } finally {
    isSubmitting.value = false
  }
}
</script>
