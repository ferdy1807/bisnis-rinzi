<template>
  <div>
    <!-- Mobile overlay -->
    <div 
      v-if="!isSidebarCollapsed" 
      @click="isSidebarCollapsed = true"
      class="fixed inset-0 bg-black/50 z-40 md:hidden"
    ></div>

    <aside 
      :class="[
        'fixed left-0 top-0 h-full bg-surface-container-low text-on-surface flex flex-col py-6 border-r border-outline-variant/30 shadow-sm z-50 transition-all duration-300',
        isSidebarCollapsed ? '-translate-x-full md:translate-x-0 w-[280px] md:w-[80px] px-2' : 'translate-x-0 w-[280px] px-4'
      ]"
    >
      <!-- Logo & Branding -->
      <div :class="['mb-8 flex items-center gap-3', isSidebarCollapsed ? 'justify-center md:px-0' : 'px-4']">
        <div class="w-10 h-10 bg-primary rounded-xl flex items-center justify-center text-on-primary font-bold text-headline-md shrink-0">BR</div>
        <div v-if="!isSidebarCollapsed" class="overflow-hidden">
        <h1 class="font-headline-md text-headline-md leading-none text-on-surface whitespace-nowrap">{{ appTitle }}</h1>
        <p class="text-label-sm text-on-surface-variant/70 whitespace-nowrap">{{ appSubtitle }}</p>
      </div>
    </div>

    <!-- Navigation Menu (Accordion & Filtered Hide Menu) -->
    <nav class="flex-1 space-y-2 custom-scrollbar overflow-y-auto pr-1">
      <div v-for="(item, index) in visibleMenuItems" :key="index" class="group">
        
        <!-- Parent Menu Button -->
        <button 
          @click="toggleMenu(index, item)"
          :class="[
            'w-full flex items-center gap-3 text-on-surface-variant py-3 rounded-lg transition-colors font-label-md text-label-md hover:bg-surface-variant/20 cursor-pointer',
            isSidebarCollapsed ? 'justify-center px-0' : 'px-4',
            activeMenu === index ? 'bg-primary/10 text-primary font-bold border-r-4 border-primary' : ''
          ]"
        >
          <span 
            class="material-symbols-outlined" 
            :style="activeMenu === index ? 'font-variation-settings: \'FILL\' 1;' : ''"
          >
            {{ item.icon }}
          </span>
          <span v-if="!isSidebarCollapsed" class="font-medium flex-1 text-left">{{ item.label }}</span>
          
          <span 
            v-if="!isSidebarCollapsed && item.children && item.children.length > 0" 
            :class="['material-symbols-outlined text-[18px] transition-transform duration-300', activeMenu === index ? 'rotate-180' : '']"
          >
            expand_more
          </span>
        </button>

        <!-- Submenu Container -->
        <div 
          v-if="item.children && item.children.length > 0 && !isSidebarCollapsed" 
          v-show="activeMenu === index"
          class="overflow-hidden transition-all duration-300 mt-1"
        >
          <div class="ml-6 pl-6 space-y-1 py-1 relative">
            <div class="absolute left-0 top-0 bottom-0 w-[1.5px] bg-outline-variant/40"></div>
            <router-link 
              v-for="(sub, sIndex) in item.children" 
              :key="sIndex"
              :to="sub.to || '#'"
              class="block py-2 px-3 rounded-md text-label-sm transition-all"
              :class="[
                /* PERBAIKAN 1: Menggunakan pencocokan string ketat === demi validitas animasi visual tunggal */
                $route.path === sub.to
                  ? 'text-primary font-bold bg-primary/10'
                  : 'text-on-surface-variant hover:text-primary hover:bg-primary/5'
              ]"
            >
              {{ sub.label }}
            </router-link>
          </div>
        </div>

      </div>
    </nav>

    <!-- Sidebar Footer (Logout) -->
    <div :class="['pt-6 mt-4 border-t border-outline-variant/30 space-y-4', isSidebarCollapsed ? 'px-0' : 'px-4']">
      <button 
        @click="$emit('logout')"
        :class="[
          'w-full flex items-center gap-3 text-error py-2 hover:bg-error/10 transition-colors font-label-md text-label-md rounded-lg cursor-pointer',
          isSidebarCollapsed ? 'justify-center px-0' : 'px-4'
        ]"
      >
        <span class="material-symbols-outlined">logout</span>
        <span v-if="!isSidebarCollapsed" class="font-medium">Logout</span>
      </button>
    </div>
  </aside>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useLayout } from '@frontend/composables/useLayout';

// Mendefinisikan tipe data struktur menu[cite: 3]
export interface MenuItem {
  label: string;
  icon?: string;
  to?: string;
  children?: MenuItem[];
}

const props = defineProps({
  appTitle: { type: String, default: 'Bisnis Rinzi' },
  appSubtitle: { type: String, default: 'Modul Sistem' },
  menuItems: { type: Array as () => MenuItem[], required: true }
});

defineEmits(['logout']);

const { isSidebarCollapsed } = useLayout();
const activeMenu = ref<number | null>(null);
const router = useRouter();
const route = useRoute();

// --- PERBAIKAN 2: DYNAMIC FILTER HIDE MENU BERBASIS ROUTER META ---
const visibleMenuItems = computed(() => {
  return props.menuItems
    .map(item => {
      // Membuat salinan objek agar tidak merusak data props asli secara mutasi langsung
      const copy = { ...item };

      // Jika menu memiliki submenu, filter berdasarkan meta.hideMenu dari file router
      if (copy.children) {
        copy.children = copy.children.filter(child => {
          const resolvedRoute = router.resolve(child.to || '#');
          return !resolvedRoute?.meta?.hideMenu;
        });
      }

      return copy;
    })
    .filter(item => {
      // Cek apakah menu induk itu sendiri di-set hideMenu di router
      if (item.to) {
        const resolvedRoute = router.resolve(item.to);
        if (resolvedRoute?.meta?.hideMenu) return false;
      }

      // Jika itu menu grup beranak, sembunyikan jika seluruh anaknya habis terfilter
      if (item.children && item.children.length === 0) {
        return false;
      }

      return true;
    });
});

// --- PERBAIKAN 3: REKONSILIASI PENCOCOKAN KETAT UNTUK ACCORDION SYNC ---
const syncActiveMenu = () => {
  activeMenu.value = null; // Reset ke posisi tertutup default sebelum pencarian ketat
  
  for (let i = 0; i < visibleMenuItems.value.length; i++) {
    const item = visibleMenuItems.value[i];
    
    // 1. Cek kecocokan absolut rute induk tunggal
    if (item.to && route.path === item.to) {
      activeMenu.value = i;
      return;
    }

    // 2. Cek kecocokan absolut di dalam daftar anak (submenu)
    if (item.children) {
      const hasActiveChild = item.children.some(child => child.to === route.path);
      if (hasActiveChild) {
        activeMenu.value = i;
        return;
      }
    }
  }
};

onMounted(() => {
  syncActiveMenu();
});

watch(() => route.path, () => {
  syncActiveMenu();
});

function toggleMenu(index: number, item: MenuItem) {
  // Jika tombol menu adalah root navigasi tunggal (tidak bercabang)
  if (item.to && (!item.children || item.children.length === 0)) {
    router.push(item.to);
    return;
  }

  // Jika kondisi sidebar sedang ciut (collapsed), abaikan pembukaan accordion dropdown
  if (isSidebarCollapsed.value) {
    return;
  }
  
  // Toggle buka-tutup kontainer menu
  activeMenu.value = activeMenu.value === index ? null : index;
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}
</style>