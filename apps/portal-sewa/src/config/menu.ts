// apps/portal-sewa/src/config/menu.ts

export interface MenuLeaf {
  label: string;
  icon?: string;
  to: string;
}

export interface MenuGroup {
  label: string;
  icon: string;
  children: MenuLeaf[];
}

export interface MenuRoot {
  label: string;
  icon: string;
  to: string;
}

export type MenuItem = MenuRoot | MenuGroup;

export const rentalMenuItems: MenuItem[] = [
  {
    label: 'Dashboard',
    icon: 'monitoring',
    to: '/'
  },
  {
    label: 'Pemesanan & Invoice',
    icon: 'receipt',
    children: [
      { label: 'Pesanan Mendatang', icon: 'upcoming', to: '/upcoming' },
      { label: 'Buat Reservasi', icon: 'add_circle', to: '/reservations-create' },
      { label: 'Riwayat Reservasi', icon: 'history', to: '/reservations-list' }
    ]
  },
  {
    label: 'Titipan & Pengembalian',
    icon: 'flowsheet',
    children: [
      { label: 'Penerimaan Titipan', icon: 'move_to_inbox', to: '/workflow-contents-received' },
      { label: 'Pengembalian & Denda', icon: 'assignment_return', to: '/returns' }
    ]
  },
  {
    label: 'Informasi Katalog',
    icon: 'menu_book',
    children: [
      { label: 'Kalender Ketersediaan', icon: 'calendar_month', to: '/availability-calendar' },
      { label: 'Katalog Box', icon: 'inventory_2', to: '/catalog' }
    ]
  }
];