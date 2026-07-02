<template>
  <div :class="['min-h-screen bg-background font-body-md transition-all duration-300', isSidebarCollapsed ? 'sidebar-collapsed' : '']">
    
    <AppSidebar 
      :app-title="appTitle" 
      :app-subtitle="appSubtitle" 
      :menu-items="menuItems"
      @logout="handleLogout"
    />
    
    <AppTopbar 
      :user-name="userName" 
      :user-role="userRole" 
    />

    <!-- Main Content Canvas -->
    <main 
      :class="[
        'pt-24 pb-24 md:pb-8 pr-4 md:pr-gutter-desktop transition-all duration-300 min-h-screen',
        isSidebarCollapsed ? 'md:ml-[80px] pl-4 md:pl-6' : 'md:ml-[280px] pl-4 md:pl-6'
      ]"
    >
      <slot/>
    </main>

  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import AppSidebar from './AppSidebar.vue';
import AppTopbar from './AppTopbar.vue';
import { useLayout } from '@frontend/composables/useLayout';
import type { MenuItem } from './AppSidebar.vue';

// Mengambil fungsi logout dari auth store
import { useAuthStore } from '@frontend/stores/auth';

defineProps({
  appTitle: { type: String, default: 'Bisnis Rinzi' },
  appSubtitle: { type: String, default: 'Dashboard' },
  menuItems: { type: Array as () => MenuItem[], required: true },
  userName: { type: String, default: 'Admin' },
  userRole: { type: String, default: 'Staf' }
});

const { isSidebarCollapsed } = useLayout();
const authStore = useAuthStore();

const handleLogout = async () => {
  await authStore.logout();
};

// Menutup sidebar otomatis di layar kecil
const checkScreenSize = () => {
  if (window.innerWidth < 1024) {
    isSidebarCollapsed.value = true;
  } else {
    isSidebarCollapsed.value = false;
  }
};

onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
});
</script>