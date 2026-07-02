import { ref } from 'vue';

// State global reaktif untuk layout (hanya di memori klien)
const isSidebarCollapsed = ref(false);

export function useLayout() {
    function toggleSidebar() {
        isSidebarCollapsed.value = !isSidebarCollapsed.value;
    }

    return {
        isSidebarCollapsed,
        toggleSidebar
    };
}