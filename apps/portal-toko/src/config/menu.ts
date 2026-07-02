// apps/portal-toko/src/config/menu.ts

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

export const tokoMenuItems: MenuItem[] = [
  {
    label: 'Dashboard',
    icon: 'dashboard',
    to: '/dashboard'
  },
  {
    label: 'POS & Kasir',
    icon: 'point_of_sale',
    children: [
      { label: 'Keranjang Belanja', icon: 'shopping_cart', to: '/pos' },
      { label: 'Riwayat Transaksi', icon: 'receipt_long', to: '/transactions' }
    ]
  },
  {
    label: 'Shift & Laporan',
    icon: 'payments',
    children: [
      { label: 'Shift Kasir', icon: 'schedule', to: '/shifts' },
      { label: 'Laporan Sesi Shift', icon: 'today', to: '/daily-report' },
      { label: 'Laporan Harian Kasir', icon: 'summarize', to: '/report/shift' }
    ]
  },
  {
    label: 'Master Barang',
    icon: 'inventory_2',
    children: [
      { label: 'Katalog Produk', icon: 'category', to: '/products' },
      { label: 'Kategori', icon: 'label', to: '/categories' },
      { label: 'Merek', icon: 'verified', to: '/brands' },
      { label: 'Satuan', icon: 'straighten', to: '/units' },
      { label: 'Pergerakan Stok', icon: 'swap_vert', to: '/stock-movement' }
    ]
  },
  {
    label: 'Keuangan',
    icon: 'account_balance_wallet',
    children: [
      { label: 'Kas Masuk', icon: 'arrow_downward', to: '/cash-in' },
      { label: 'Kas Keluar', icon: 'arrow_upward', to: '/cash-out' }
    ]
  }
];
