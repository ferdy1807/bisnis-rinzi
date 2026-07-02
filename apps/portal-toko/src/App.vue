<template>
  <DashboardLayout
    :menu-items="tokoMenuItems"
    app-title="Portal Toko"
    app-subtitle="Internal Staff"
  >
    <RouterView />
  </DashboardLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import DashboardLayout from '@frontend/layouts/DashboardLayout.vue';
import { tokoMenuItems } from './config/menu';
import { useAuthStore } from '@frontend/stores/auth';
import { syncTokenFromUrl } from '@frontend/services/token';

const authStore = useAuthStore();

onMounted(async () => {
  // 1. Tangkap token dari URL jika ada
  syncTokenFromUrl();
  
  // 2. Baru inisialisasi session (yang akan mengecek getToken())
  await authStore.initializeSession();
});
</script>