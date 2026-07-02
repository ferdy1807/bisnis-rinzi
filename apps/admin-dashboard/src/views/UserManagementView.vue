<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { authApi } from '@frontend/api/auth';

// --- STATE MANAGEMENT ---
const users = ref<any[]>([]);
const isLoading = ref(false);
const isSubmitLoading = ref(false);
const isModalOpen = ref(false);
const isRoleModalOpen = ref(false);
const selectedUser = ref<any>(null);
const errorMessage = ref('');
const successMessage = ref('');
const searchQuery = ref('');

const roleForm = ref({
  role: ''
});

// Form State Registrasi Staf Baru (/api/auth/register)
const form = ref({
  fullname: '',
  username: '',
  password: '',
  role: 'PEGAWAI' // Default role sekunder staf lapangan
});

// --- FETCH DATA STAFF ---
const fetchUsers = async () => {
  isLoading.value = true;
  errorMessage.value = '';
  try {
    const response = await authApi.getUsers();
    if (response && response.length > 0) {
      users.value = response;
    } else {
      // Fallback data simulasi lokal jika dispatcher mengembalikan metrik murni
      users.value = [
        { id: 'usr-01', full_name: 'Ferdy Perdana', username: 'ferdy.owner', email: 'owner@toko.com', role: 'OWNER', status: 'ACTIVE' },
        { id: 'usr-02', full_name: 'Kasir Utama', username: 'cashier.toko', email: 'cashier@toko.com', role: 'CASHIER', status: 'ACTIVE' },
        { id: 'usr-03', full_name: 'Staf Logistik', username: 'staff.rental', email: 'staff@toko.com', role: 'PEGAWAI', status: 'ACTIVE' }
      ];
    }
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memuat manifes akun staf.';
  } finally {
    isLoading.value = false;
  }
};

// --- REGISTER NEW STAFF (POST /api/auth/register) ---
const handleRegisterStaff = async () => {
  if (!form.value.fullname || !form.value.username|| !form.value.password ) {
    errorMessage.value = 'Nama, Username, dan Password wajib diisi.';
    return;
  }

  isSubmitLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';

  try {
    const payload = {
      full_name: form.value.fullname,
      username: form.value.username,
      password: form.value.password,
      role: form.value.role
    };

    await authApi.createUser(payload);
    
    successMessage.value = `Akun staf baru dengan wewenang ${form.value.role} berhasil didaftarkan secara sah.`;
    isModalOpen.value = false;
    
    // Reset form
    form.value = { fullname: '', username: '', password: '', role: 'PEGAWAI' };
    await fetchUsers();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal meregistrasikan akun operator baru.';
  } finally {
    isSubmitLoading.value = false;
  }
};

// --- ACTIONS ---
const handleDeleteUser = async (user: any) => {
  if (confirm(`Apakah Anda yakin ingin menghapus staf ${user.full_name || user.name}?`)) {
    try {
      await authApi.deleteUser(user.id);
      successMessage.value = `Staf ${user.full_name || user.name} berhasil dihapus.`;
      await fetchUsers();
    } catch (error: any) {
      errorMessage.value = error?.message || 'Gagal menghapus akun staf.';
    }
  }
};

const handleForceLogout = async (user: any) => {
  if (confirm(`Paksa keluar (force logout) staf ${user.full_name || user.name} dari sistem?`)) {
    try {
      await authApi.forceLogoutUser(user.id);
      successMessage.value = `Sesi staf ${user.full_name || user.name} berhasil dihentikan.`;
    } catch (error: any) {
      errorMessage.value = error?.message || 'Gagal melakukan force logout.';
    }
  }
};

const openRoleModal = (user: any) => {
  selectedUser.value = user;
  roleForm.value.role = user.role;
  isRoleModalOpen.value = true;
};

const handleChangeRole = async () => {
  if (!selectedUser.value) return;
  isSubmitLoading.value = true;
  errorMessage.value = '';
  successMessage.value = '';
  try {
    await authApi.updateProfile(selectedUser.value.id, { role: roleForm.value.role });
    successMessage.value = `Hak akses untuk ${selectedUser.value.full_name || selectedUser.value.name} berhasil diperbarui.`;
    isRoleModalOpen.value = false;
    await fetchUsers();
  } catch (error: any) {
    errorMessage.value = error?.message || 'Gagal memperbarui hak akses.';
  } finally {
    isSubmitLoading.value = false;
  }
};

// --- FILTER SEARCH ---
const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value;
  const query = searchQuery.value.toLowerCase();
  return users.value.filter(u => 
    u.full_name?.toLowerCase().includes(query) ||
    u.name?.toLowerCase().includes(query) ||
    u.username?.toLowerCase().includes(query) ||
    u.role?.toLowerCase().includes(query)
  );
});

onMounted(() => {
  fetchUsers();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-200 pb-5">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
          <i class="pi pi-admin-panel-settings text-indigo-600"></i>
          Otoritas Otorisasi & Manajemen Akun Staf
        </h1>
        <p class="text-sm text-slate-500 mt-1">
          Kendali penuh Owner dalam meregistrasikan operator portal, membatasi hak akses ruang lingkup database, serta menjaga kredensial penanggung jawab shift.
        </p>
      </div>
      <button 
        @click="isModalOpen = true"
        class="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-semibold px-4 py-2.5 rounded-xl shadow-md transition-all duration-200"
      >
        <i class="pi pi-user-plus text-xs"></i>
        Daftarkan Staf Baru
      </button>
    </div>

    <div v-if="errorMessage" class="bg-rose-50 border border-rose-200 text-rose-700 p-4 rounded-xl text-sm font-medium">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="bg-emerald-50 border border-emerald-200 text-emerald-700 p-4 rounded-xl text-sm font-medium">
      {{ successMessage }}
    </div>

    <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center gap-3 justify-between">
      <div class="relative w-full sm:w-80">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <i class="pi pi-search text-sm"></i>
        </span>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari berdasarkan nama, username, atau role..."
          class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
        />
      </div>
      
      <div class="text-xs font-semibold text-slate-400 uppercase tracking-wider">
        Total Operator Aktif: <span class="text-slate-800 font-bold font-mono">{{ filteredUsers.length }} Staf</span>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden relative">
      <div v-if="isLoading" class="absolute inset-0 bg-white/70 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center">
        <i class="pi pi-spin pi-spinner text-2xl text-indigo-600 mb-2"></i>
        <span class="text-xs font-semibold text-slate-500">Mengekstrak berkas kredensial operator...</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 text-slate-400 text-[11px] font-bold tracking-wider uppercase border-b border-slate-200">
              <th class="py-4 px-6">Username</th>
              <th class="py-4 px-6">ID Pengguna</th>
              <th class="py-4 px-6">Nama Lengkap Resmi</th>
              <th class="py-4 px-6 text-center">Grup Peran (ROLE)</th>
              <th class="py-4 px-6 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm font-medium text-slate-700">
            <tr v-for="user in filteredUsers" :key="user.id" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900 flex items-center gap-2">
                  <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
                  {{ user.username }}
                </div>
              </td>
              <td class="py-4 px-6 font-mono text-xs text-slate-400 select-all">{{ user.id }}</td>
              <td class="py-4 px-6">
                <div class="font-bold text-slate-900 flex items-center gap-2">
                  <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
                  {{ user.full_name }}
                </div>
              </td>
              <td class="py-4 px-6">
                <span 
                  class="px-2.5 py-1 text-[10px] font-extrabold rounded-md border tracking-wider uppercase"
                  :class="{
                    'bg-indigo-50 text-indigo-700 border-indigo-200': user.role === 'OWNER',
                    'bg-emerald-50 text-emerald-700 border-emerald-200': user.role === 'CASHIER',
                    'bg-amber-50 text-amber-700 border-amber-200': user.role === 'PEGAWAI'
                  }">{{ user.role }}
                </span>
              </td>
              <td class="py-4 px-6">
                <div class="flex items-center justify-center gap-2">
                  <button 
                    @click="openRoleModal(user)"
                    class="p-1.5 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                    title="Ganti Role"
                  >
                    <i class="pi pi-shield"></i>
                  </button>
                  <button 
                    @click="handleForceLogout(user)"
                    class="p-1.5 text-amber-600 hover:bg-amber-50 rounded-lg transition-colors"
                    title="Force Logout"
                  >
                    <i class="pi pi-sign-out"></i>
                  </button>
                  <button 
                    @click="handleDeleteUser(user)"
                    class="p-1.5 text-rose-600 hover:bg-rose-50 rounded-lg transition-colors"
                    title="Hapus Staf"
                  >
                    <i class="pi pi-trash"></i>
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
          <h3 class="text-base font-bold text-slate-900">Registrasi Otoritas Akun Operator</h3>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-slate-600"><i class="pi pi-times"></i></button>
        </div>
        
        <form @submit.prevent="handleRegisterStaff" class="p-6 space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Lengkap Sesuai KTP</label>
            <input 
              v-model="form.fullname"
              type="text"
              required
              placeholder="Contoh: Ahmad Subagja"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Username Unik</label>
              <input 
                v-model="form.username"
                type="text"
                required
                placeholder="ahmad.staff"
                class="w-full px-3 py-2 font-mono border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all"
              />
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Kata Sandi Akses (Password)</label>
            <input 
              v-model="form.password"
              type="password"
              required
              minlength="6"
              placeholder="Minimal 6 karakter rahasia..."
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all"
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Grup Hak Akses (Role Operator)</label>
            <select 
              v-model="form.role"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 transition-all font-semibold text-slate-700"
            >
              <option value="CASHIER">CASHIER (Akses: Kasir, POS, Inventaris)</option>
              <option value="PEGAWAI">PEGAWAI (Akses: Rental Unit)</option>
            </select>
          </div>

          <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-100 mt-6">
            <button type="button" @click="isModalOpen = false" class="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 rounded-lg border border-slate-200 transition-colors">
              Batal
            </button>
            <button 
              type="submit" 
              :disabled="isSubmitLoading"
              class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-md transition-colors disabled:opacity-50"
            >
              <span>{{ isSubmitLoading ? 'Mendaftarkan Akun...' : 'Daftarkan Anggota Staf' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL GANTI ROLE -->
    <div v-if="isRoleModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md border border-slate-200 overflow-hidden transform transition-all animate-fade-in">
        <div class="p-6 border-b border-slate-100 flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-900">Ubah Grup Hak Akses</h3>
          <button @click="isRoleModalOpen = false" class="text-slate-400 hover:text-slate-600"><i class="pi pi-times"></i></button>
        </div>
        
        <form @submit.prevent="handleChangeRole" class="p-6 space-y-4">
          <div v-if="selectedUser" class="bg-slate-50 p-4 rounded-lg border border-slate-200 space-y-2 mb-2">
            <div class="text-xs text-slate-500 uppercase tracking-wider font-bold">Informasi Akun</div>
            <div class="flex flex-col gap-1">
              <span class="text-sm font-semibold text-slate-900">{{ selectedUser.full_name || selectedUser.name }}</span>
              <span class="text-xs font-mono text-slate-600">{{ selectedUser.username }}</span>
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-500 uppercase tracking-wider">Grup Hak Akses Baru</label>
            <select 
              v-model="roleForm.role"
              class="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-slate-50 outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 transition-all font-semibold text-slate-700"
            >
              <option value="OWNER">OWNER (Akses Penuh)</option>
              <option value="CASHIER">CASHIER (Akses: Kasir, POS, Inventaris)</option>
              <option value="PEGAWAI">PEGAWAI (Akses: Rental Unit)</option>
            </select>
          </div>

          <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-100 mt-6">
            <button type="button" @click="isRoleModalOpen = false" class="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 rounded-lg border border-slate-200 transition-colors">
              Batal
            </button>
            <button 
              type="submit" 
              :disabled="isSubmitLoading"
              class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-md transition-colors disabled:opacity-50"
            >
              <span>{{ isSubmitLoading ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>