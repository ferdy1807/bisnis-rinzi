<template>
  <DashboardLayout
    :menu-items="adminDashboardMenuItems"
    app-title="Admin Dashboard"
    app-subtitle="Manajemen Bisnis"
  >
    <RouterView />
  </DashboardLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import DashboardLayout from '@frontend/layouts/DashboardLayout.vue';
import { adminDashboardMenuItems } from './config/menu';
import { useAuthStore } from '@frontend/stores/auth';
import { syncTokenFromUrl } from '@frontend/services/token';

const authStore = useAuthStore();

onMounted(async () => {
  syncTokenFromUrl();
  await authStore.initializeSession();
});
</script>

<style>
.material-symbols-outlined {
  font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
}
</style>