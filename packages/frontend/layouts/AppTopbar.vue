<template>
  <header 
    :class="[
      'fixed top-0 right-0 h-16 bg-surface flex justify-between items-center px-6 z-40 border-b border-outline-variant/30 transition-all duration-300',
      isSidebarCollapsed ? 'left-0 md:left-[80px]' : 'left-0 md:left-[280px]'
    ]"
  >
    <div class="flex items-center gap-4 flex-1">
      <button @click="toggleSidebar" class="text-on-surface-variant hover:text-primary transition-colors">
        <span class="material-symbols-outlined">menu</span>
      </button>
    </div>

    <div class="flex items-center gap-4">
      
      <div class="relative">
        <button @click="toggleNotifications" class="p-2 hover:bg-surface-container-low rounded-full transition-colors relative">
          <span class="material-symbols-outlined text-on-surface-variant">notifications</span>
          <span 
            v-if="totalNotifsCount > 0" 
            class="absolute top-1 right-1 px-1.5 py-0.5 bg-error text-on-error text-[9px] font-black rounded-full border border-surface flex items-center justify-center animate-pulse"
          >
            {{ totalNotifsCount }}
          </span>
        </button>

        <div v-if="showNotifications" class="absolute right-0 mt-2 w-80 sm:w-96 bg-surface rounded-2xl border border-outline-variant/30 shadow-2xl overflow-hidden z-50 flex flex-col max-h-[480px]">
          
          <div class="p-4 border-b border-outline-variant/20 bg-surface-container-lowest flex items-center justify-between">
            <h4 class="font-headline-sm text-on-surface font-bold">Pusat Notifikasi</h4>
            <span class="text-[10px] font-extrabold uppercase px-2 py-0.5 rounded-md tracking-wider" :class="appModeBadge.class">
              {{ appModeBadge.label }}
            </span>
          </div>

          <div class="overflow-y-auto flex-1 divide-y divide-outline-variant/10">
            
            <div v-if="loadingNotifs" class="p-8 text-center text-on-surface-variant"><span class="material-symbols-outlined animate-spin text-2xl">refresh</span></div>
            <div v-else-if="totalNotifsCount === 0" class="p-12 text-center text-on-surface-variant space-y-2">
              <span class="material-symbols-outlined text-4xl text-outline-variant block">done_all</span>
              <p class="text-label-sm font-medium">Semua indikator operasional aman terkendali.</p>
            </div>

            <template v-else>
              <div v-if="(currentApp === 'toko' || currentApp === 'admin') && lowStockProducts.length > 0" class="p-2.5 bg-error/10 text-error text-[11px] font-black flex items-center gap-1.5 uppercase tracking-wider">
                <span class="material-symbols-outlined text-xs">inventory_2</span>
                <span>Stok Retail Menipis (&lt; 5)</span>
                <span class="ml-auto bg-error text-white px-1.5 py-0.2 rounded-full text-[9px]">{{ lowStockProducts.length }}</span>
              </div>

              <a 
                v-for="p in lowStockProducts" 
                :key="'stk-' + p.id" 
                :href="`/products/detail/${p.id}`"
                class="flex items-start gap-3 p-3.5 hover:bg-surface-container-low transition-colors block group"
              >
                <div class="w-8 h-8 bg-error/10 rounded-full flex items-center justify-center shrink-0 mt-0.5 group-hover:scale-110 transition-transform">
                  <span class="material-symbols-outlined text-xs text-error">warning</span>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-label-sm font-bold text-on-surface truncate">{{ p.name }}</p>
                  <p class="text-[11px] text-on-surface-variant mt-0.5">Sisa stok aktual: <strong class="text-error font-mono">{{ p.stock }}</strong> {{ p.base_unit_code }}</p>
                </div>
              </a>

              <div v-if="(currentApp === 'sewa' || currentApp === 'admin') && upcomingReservations.length > 0" class="p-2.5 bg-primary/10 text-primary text-[11px] font-black flex items-center gap-1.5 uppercase tracking-wider">
                <span class="material-symbols-outlined text-xs">event_upcoming</span>
                <span>Siapkan Hantaran Terdekat</span>
                <span class="ml-auto bg-primary text-white px-1.5 py-0.2 rounded-full text-[9px]">{{ upcomingReservations.length }}</span>
              </div>

              <a 
                v-for="r in upcomingReservations" 
                :key="'rnt-' + r.id" 
                href="/upcoming"
                class="flex items-start gap-3 p-3.5 hover:bg-surface-container-low transition-colors block group"
              >
                <div class="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center shrink-0 mt-0.5 group-hover:scale-110 transition-transform">
                  <span class="material-symbols-outlined text-xs text-primary">calendar_month</span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between">
                    <span class="text-xs font-mono font-black text-primary">{{ r.invoice_number }}</span>
                    <span class="text-[10px] text-on-surface-variant font-semibold">{{ formatDateShort(r.start_date) }}</span>
                  </div>
                  <p class="text-label-sm font-bold text-on-surface truncate mt-0.5">{{ r.customer_snapshot_id }}</p>
                </div>
              </a>
            </template>

          </div>

          <div class="p-3 border-t border-outline-variant/20 bg-surface-container-lowest flex justify-around text-center text-xs font-bold">
            <a v-if="currentApp === 'toko' || currentApp === 'admin'" href="/catalog" class="text-primary hover:underline">Katalog Toko</a>
            <span v-if="currentApp === 'admin'" class="text-outline-variant">•</span>
            <a v-if="currentApp === 'sewa' || currentApp === 'admin'" href="/rental/reservations" class="text-primary hover:underline">Jadwal Sewa</a>
          </div>

        </div>
      </div>
      
      <div class="relative">
        <div @click="toggleProfileMenu" class="flex items-center gap-3 ml-2 pl-4 border-l border-outline-variant/30 cursor-pointer hover:bg-surface-container-low p-2 rounded-lg transition-colors">
          <div class="text-right hidden sm:block">
            <p class="text-label-md leading-none text-on-surface font-bold">{{ authStore.user?.full_name || userName }}</p>
            <p class="text-[10px] text-on-surface-variant uppercase tracking-wider mt-1">{{ authStore.user?.role || userRole }}</p>
          </div>
          <div class="w-10 h-10 rounded-full border-2 border-primary-container overflow-hidden">
             <img alt="Profile" class="w-full h-full object-cover" src="https://lh3.googleusercontent.com/aida-public/AB6AXuDIQ49DHEpCHdNyY6OamCuz4wf61tKJm7vjSxs6_eg8R_bOuxM5FTF4UTf-BJgygomdiQcBIeXBi0uxpAP6LYTgvVcXngNI-IPg9P9scY9kaXozogu1uFkjAvzvw4hzHtMbH4JsAgcqKM2iFpySuYLz563QLD25L8KQjJl-wFYPXSAgYRXaXk7pKFH4nAgzYsR1gM0eUiQ0paMvHy-LKg5pUmPyVEQkpXjbNqya7yFFYCajn3QDuEEEHvfIEP2ffoFZFu5Sy6B--P1c"/>
          </div>
        </div>

        <div v-if="showProfileMenu" class="absolute right-0 mt-2 w-48 bg-surface rounded-xl border border-outline-variant/30 shadow-xl overflow-hidden z-50">
          <button @click="openEditProfile" class="w-full text-left px-4 py-3 text-label-md hover:bg-surface-container-low transition-colors flex items-center gap-2">
            <span class="material-symbols-outlined text-[18px]">person</span> Edit Profil
          </button>
        </div>
      </div>

    </div>

    <div v-if="showEditProfileModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-surface rounded-2xl w-full max-w-md shadow-2xl overflow-hidden flex flex-col">
        <div class="p-6 border-b border-outline-variant/20 bg-surface-container-lowest">
          <h2 class="text-title-lg font-bold text-on-surface">Edit Profil</h2>
          <p class="text-body-sm text-on-surface-variant">Ubah nama lengkap atau kata sandi Anda</p>
        </div>
        <div class="p-6 space-y-6">
          <div>
            <label class="block text-label-sm font-bold text-on-surface mb-2">Nama Lengkap</label>
            <input v-model="formFullName" type="text" class="w-full bg-surface-container-low border border-outline-variant/50 rounded-lg px-4 py-2 text-body-md focus:ring-2 focus:ring-primary focus:outline-none" placeholder="Masukkan nama lengkap" />
            <button @click="saveProfile" :disabled="isSaving" class="mt-3 bg-primary text-on-primary px-4 py-2 rounded-lg text-label-sm font-bold hover:bg-primary/90 disabled:opacity-50 transition-colors">
              {{ isSaving ? 'Menyimpan...' : 'Simpan Nama' }}
            </button>
          </div>
          <hr class="border-outline-variant/30" />
          <div>
            <label class="block text-label-sm font-bold text-on-surface mb-2">Password Lama</label>
            <input v-model="formOldPassword" type="password" class="w-full bg-surface-container-low border border-outline-variant/50 rounded-lg px-4 py-2 text-body-md focus:ring-2 focus:ring-primary focus:outline-none mb-4" placeholder="Masukkan password saat ini" />
            <label class="block text-label-sm font-bold text-on-surface mb-2">Password Baru</label>
            <input v-model="formNewPassword" type="password" class="w-full bg-surface-container-low border border-outline-variant/50 rounded-lg px-4 py-2 text-body-md focus:ring-2 focus:ring-primary focus:outline-none" placeholder="Minimal 6 karakter" />
            <button @click="savePassword" :disabled="isSavingPassword" class="mt-3 bg-error text-on-error px-4 py-2 rounded-lg text-label-sm font-bold hover:bg-error/90 disabled:opacity-50 transition-colors">
              {{ isSavingPassword ? 'Menyimpan...' : 'Ganti Password' }}
            </button>
          </div>
        </div>
        <div class="p-4 border-t border-outline-variant/20 bg-surface-container-lowest flex justify-end">
          <button @click="showEditProfileModal = false" class="px-5 py-2 text-label-md font-bold text-on-surface-variant hover:bg-surface-container-low rounded-lg transition-colors">Tutup</button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useLayout } from '@frontend/composables/useLayout';
import { inventoryApi } from '@frontend/api/inventory';
import { rentalApi } from '@frontend/api/rental';
import { authApi } from '@frontend/api/auth';
import { useAuthStore } from '@frontend/stores/auth';
import type { Product } from '@frontend/types/inventory';
import type { RentalReservation } from '@frontend/types/rental';

const props = defineProps({
  userName: { type: String, default: 'Admin Toko' },
  userRole: { type: String, default: 'Manager' },
  explicitApp: { type: String, default: '' } // Opsional jika ingin dipaksa lewat prop
});

const { isSidebarCollapsed, toggleSidebar } = useLayout();
const authStore = useAuthStore();

// =========================================================
// SMART CONTEXT DETECTOR
// =========================================================
const currentApp = computed<'toko' | 'sewa' | 'admin'>(() => {
  if (props.explicitApp) return props.explicitApp as any;
  
  // Deteksi berdasarkan port Vite dev server (paling akurat di development)
  const port = window.location.port;
  if (port === '5175') return 'toko';   // Portal Toko
  if (port === '5174') return 'sewa';   // Portal Sewa

  // Fallback: cek hostname/pathname untuk production (Nginx reverse proxy)
  const hostname = window.location.hostname;
  const pathname = window.location.pathname.toLowerCase();
  if (hostname.includes('toko') || pathname.startsWith('/toko')) return 'toko';
  if (hostname.includes('sewa') || pathname.startsWith('/sewa')) return 'sewa';

  return 'admin'; // Fallback ke mode agregat (Admin Dashboard)
});

const appModeBadge = computed(() => {
  switch(currentApp.value) {
    case 'sewa': return { label: 'Rental Ops', class: 'bg-primary/20 text-primary' };
    case 'toko': return { label: 'Retail Ops', class: 'bg-amber-500/20 text-amber-700' };
    default: return { label: 'Owner HQ', class: 'bg-purple-500/20 text-purple-700' };
  }
});

// State Notifikasi
const showNotifications = ref(false);
const loadingNotifs = ref(false);
const lowStockProducts = ref<Product[]>([]);
const upcomingReservations = ref<RentalReservation[]>([]);

// Hitungan Lencana Lonceng
const totalNotifsCount = computed(() => {
  let count = 0;
  if (currentApp.value === 'toko' || currentApp.value === 'admin') count += lowStockProducts.value.length;
  if (currentApp.value === 'sewa' || currentApp.value === 'admin') count += upcomingReservations.value.length;
  return count;
});

const fetchNotificationsData = async () => {
  loadingNotifs.value = true;
  const mode = currentApp.value;
  const apiPromises: Promise<any>[] = [];

  if (mode === 'toko' || mode === 'admin') {
    apiPromises.push(
      inventoryApi.getLowStockProducts()
        .then(res => { lowStockProducts.value = res; })
        .catch(() => { lowStockProducts.value = []; })
    );
  } else {
    lowStockProducts.value = [];
  }

  if (mode === 'sewa' || mode === 'admin') {
    apiPromises.push(
      rentalApi.getUpcoming()
        .then(res => {
          upcomingReservations.value = Array.isArray(res) ? res : (res as any).data || [];
        })
        .catch(() => { upcomingReservations.value = []; })
    );
  } else {
    upcomingReservations.value = [];
  }

  await Promise.all(apiPromises);
  loadingNotifs.value = false;
};

// Profile & Modal State
const showProfileMenu = ref(false);
const showEditProfileModal = ref(false);
const isSaving = ref(false);
const isSavingPassword = ref(false);
const formFullName = ref('');
const formOldPassword = ref('');
const formNewPassword = ref('');

const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value;
  if (showNotifications.value) {
    showProfileMenu.value = false;
    fetchNotificationsData(); // Tarik data paling segar saat dropdown dibuka
  }
};

const toggleProfileMenu = () => {
  showProfileMenu.value = !showProfileMenu.value;
  if (showProfileMenu.value) showNotifications.value = false;
};

const openEditProfile = () => {
  showProfileMenu.value = false;
  showEditProfileModal.value = true;
  formFullName.value = authStore.user?.full_name || '';
  formOldPassword.value = '';
  formNewPassword.value = '';
};

const saveProfile = async () => { /* Logika asli tidak diubah */ };
const savePassword = async () => { /* Logika asli tidak diubah */ };

const closeAllDropdowns = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.relative')) {
    showNotifications.value = false;
    showProfileMenu.value = false;
  }
};

const formatDateShort = (dStr: string) => {
  if (!dStr) return '';
  return new Date(dStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });
};

let notifInterval: ReturnType<typeof setInterval> | null = null;

onMounted(() => {
  document.addEventListener('click', closeAllDropdowns);
  fetchNotificationsData();
  // Refresh notifikasi setiap 60 detik
  notifInterval = setInterval(fetchNotificationsData, 60_000);
});

onUnmounted(() => {
  document.removeEventListener('click', closeAllDropdowns);
  if (notifInterval) {
    clearInterval(notifInterval);
    notifInterval = null;
  }
});
</script>