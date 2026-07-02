// apps/admin-dashboard/src/config/menu.ts

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

export const adminDashboardMenuItems: MenuItem[] = [
    // --- 1. DASHBOARD INTELIJEN & ANALITIK ---
    {
        label: 'Dashboard Intelijen',
        icon: 'Wallet', // CorporateDashboardAnalytics
        to: '/analytics'
    },

    // --- 2. MANAJEMEN KEUANGAN (FINANCE_DB, CASH_DB & POS_DB) ---
    {
        label: 'Manajemen Keuangan',
        icon: 'fact_check',
        children: [
            { label: 'Kliring Tutup Shift', icon: 'account_balance', to: '/closing-audit' }, // getDailyClosings() & getShifts()
            { label: 'Kategori Pengeluaran Toko', icon: 'money_off', to: '/finance/expense-categories' },
            { label: 'Transaksi Pendapatan Internal', icon: 'input', to: '/finance/internal-incomes' },
            { label: 'Buku Besar & Jurnal', icon: 'menu_book', to: '/ledger' }, // getJournalEntries()
            { label: 'Bagan Akun (COA Master)', icon: 'account_tree', to: '/finance/coa' }, // getAccounts()
            { label: 'Kunci Periode Akuntansi', icon: 'lock_clock', to: '/finance/periods' }, // lockPeriod()
            { label: 'Log Hasil Rekonsiliasi', icon: 'rule', to: '/reconciliation-logs' },
            { label: 'Laporan Finansial Standar', icon: 'analytics', to: '/finance/reports' }, // ProfitLoss, BalanceSheet
            { label: 'Ekspor Dokumen Resmi', icon: 'download', to: '/finance/export' } // getExportUrl()
        ]
    },

    // --- 3. KONTROL INVENTARIS MASTER (INVENTORY_DB & POS_DB) ---
    {
        label: 'Kontrol Inventaris Master',
        icon: 'inventory_2',
        children: [
            { label: 'Katalog Produk Gabungan', icon: 'layers', to: '/catalog' }, // getProducts() & updateStockThreshold()
            { label: 'Kategori & Merek Barang', icon: 'bookmarks', to: '/catalog/metadata' },
            { label: 'Riwayat Harga Modal (HPP)', icon: 'history_toggle_off', to: '/catalog/cost-histories' }, // <--- BARU: getCostHistories() & addCostHistory()
            { label: 'Penyesuaian & Aturan Stok', icon: 'tune', to: '/catalog/stock-adjustments' }, // adjustStock()
            { label: 'Riwayat Mutasi Barang', icon: 'swap_horiz', to: '/catalog/movements' }, // getStockMovements()
            { label: 'Riwayat Transaksi POS Ritel', icon: 'receipt_long', to: '/catalog/pos-sales' } // getSalesHistory()
        ]
    },

    // --- 4. OTORITAS LOGISTIK & AUDIT RENTAL (RENTAL_DB) ---
    {
        label: 'Otoritas & Logistik Rental',
        icon: 'supervised_user_circle',
        children: [
            { label: 'Kontrak & Reservasi Aktif', icon: 'calendar_today', to: '/rental/reservations' }, // getAllReservations()
            { label: 'Gudang Jaminan Konsumen', icon: 'lock', to: '/rental/deposits' }, // <--- BARU: saveDepositItems()
            { label: 'Otoritas Pembatalan (Rollback)', icon: 'published_with_changes', to: '/rental/reversals' }, // <--- BARU: rollbackPickupReservation(), undoReady(), cancelReservation()
            { label: 'Audit Kerusakan Rental', icon: 'gavel', to: '/rental-damages' }, // getRentalDamages() & settleRentalDamage()
            { label: 'Cek Ketersediaan Unit', icon: 'event_available', to: '/rental/availability' } // checkAvailability()
        ]
    },

    // --- 5. OTORITAS KEAMANAN & SISTEM DISPATCHER ---
    {
        label: 'Keamanan & Aturan Sistem',
        icon: 'admin_panel_settings',
        children: [
            { label: 'Manajemen Akun Staf', icon: 'manage_accounts', to: '/users' }, // register staf baru
            { label: 'Jejak Log Audit (Trail)', icon: 'assignment_ind', to: '/audit-trail' }, // /api/auth/audit-logs
            { label: 'Dispatcher Analitik Dinamis', icon: 'dynamic_feed', to: '/system/dynamic-analytics' } // getAnalyticsData()
        ]
    }
];