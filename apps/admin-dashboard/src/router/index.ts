// apps/admin-dashboard/src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import { authApi } from '@frontend/api/auth';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/analytics'
    },
    // 1. Dashboard Intelijen Utama
    {
      path: '/analytics',
      name: 'analytics',
      component: () => import('../views/DashboardAnalyticsView.vue'),
      meta: { requiresAuth: true }
    },
    // 2. Manajemen Keuangan
    {
      path: '/closing-audit',
      name: 'closing-audit',
      component: () => import('../views/FinancialClosingAuditView.vue'), // Menangani getDailyClosings & getShifts
      meta: { requiresAuth: true }
    },
    {
      path: '/finance/expense-categories',
      name: 'finance-expense-categories',
      component: () => import('../views/ExpenseCategoriesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/finance/internal-incomes',
      name: 'finance-internal-incomes',
      component: () => import('../views/InternalIncomesView.vue'), // Menangani getInternalIncomes
      meta: { requiresAuth: true }
    },
    {
      path: '/ledger',
      name: 'ledger',
      component: () => import('../views/GlobalLedgerView.vue'), // Menangani getJournalEntries & getGeneralLedgerReport
      meta: { requiresAuth: true, hideMenu: false }
    },
    {
      path: '/finance/coa',
      name: 'finance-coa',
      component: () => import('../views/CoaManagementView.vue'), // Menangani getAccounts
      meta: { requiresAuth: true, hideMenu: false }
    },
    {
      path: '/finance/periods',
      name: 'finance-periods',
      component: () => import('../views/AccountingPeriodsView.vue'), // Menangani getPeriods & lockPeriod
      meta: { requiresAuth: true, hideMenu: false }
    },
    {
      path: '/reconciliation-logs',
      name: 'reconciliation-logs',
      component: () => import('../views/ReconciliationLogsView.vue'), // Menangani getReconciliationLogs
      meta: { requiresAuth: true, hideMenu: false }
    },
    {
      path: '/finance/reports',
      name: 'finance-reports',
      component: () => import('../views/FinancialReportsView.vue'), // Menangani ProfitLoss & BalanceSheet
      meta: { requiresAuth: true, hideMenu: false }
    },
    {
      path: '/finance/export',
      name: 'finance-export',
      component: () => import('../views/FinanceExportView.vue'), // Menangani getExportUrl
      meta: { requiresAuth: true, hideMenu: false }
    },
    // 3. Kontrol Inventaris Master
    {
      path: '/catalog',
      name: 'catalog',
      component: () => import('../views/CatalogMasterView.vue'), // Menangani getProducts & updateStockThreshold
      meta: { requiresAuth: true }
    },
    {
      path: '/catalog/metadata',
      name: 'catalog-metadata',
      component: () => import('../views/CatalogMetadataView.vue'), // Menangani getCategories & getBrands
      meta: { requiresAuth: true }
    },
    {
      path: '/catalog/cost-histories',
      name: 'catalog-cost-histories',
      component: () => import('../views/ProductCostHistoriesView.vue'), // Menangani getCostHistories & addCostHistory
      meta: { requiresAuth: true }
    },
    {
      path: '/catalog/stock-adjustments',
      name: 'catalog-stock-adjustments',
      component: () => import('../views/StockAdjustmentsView.vue'), // Menangani adjustStock & getLowStockProducts
      meta: { requiresAuth: true }
    },
    {
      path: '/catalog/movements',
      name: 'catalog-movements',
      component: () => import('../views/StockMovementsView.vue'), // Menangani getStockMovements
      meta: { requiresAuth: true }
    },
    {
      path: '/catalog/pos-sales',
      name: 'catalog-pos-sales',
      component: () => import('../views/PosSalesHistoryView.vue'), // Menangani getSalesHistory & getTopProducts
      meta: { requiresAuth: true }
    },
    // 4. Otoritas Logistik & Audit Rental
    {
      path: '/rental/reservations',
      name: 'rental-reservations',
      component: () => import('../views/RentalReservationsView.vue'), // Menangani getAllReservations, getActive, dll.
      meta: { requiresAuth: true }
    },
    {
      path: '/rental/deposits',
      name: 'rental-deposits',
      component: () => import('../views/RentalDepositsView.vue'), // Menangani saveDepositItems
      meta: { requiresAuth: true }
    },
    {
      path: '/rental/reversals',
      name: 'rental-reversals',
      component: () => import('../views/RentalReversalsView.vue'), // Menangani rollbackPickupReservation, undoReady, cancelReservation
      meta: { requiresAuth: true }
    },
    {
      path: '/rental-damages',
      name: 'rental-damages',
      component: () => import('../views/RentalDamagesView.vue'), // Menangani getRentalDamages & settleRentalDamage
      meta: { requiresAuth: true }
    },
    {
      path: '/rental/availability',
      name: 'rental-availability',
      component: () => import('../views/RentalAvailabilityView.vue'), // Menangani checkAvailability
      meta: { requiresAuth: true }
    },
    // 5. Otoritas Keamanan & Aturan Sistem
    {
      path: '/users',
      name: 'users',
      component: () => import('../views/UserManagementView.vue'), // Menangani register staf baru
      meta: { requiresAuth: true }
    },
    {
      path: '/audit-trail',
      name: 'audit-trail',
      component: () => import('../views/AuditTrailView.vue'), // Menangani getMe dan log internal via /api/auth/audit-logs
      meta: { requiresAuth: true }
    },
    {
      path: '/system/dynamic-analytics',
      name: 'dynamic-analytics',
      component: () => import('../views/DynamicAnalyticsView.vue'), // Menangani getAnalyticsData
      meta: { requiresAuth: true }
    },
    // Fallback Wildcard Route
    {
      path: '/:pathMatch(.*)*',
      redirect: '/analytics'
    }
  ]
});

/**
 * Navigation Guard — Validasi token via authApi.getMe() langsung.
 * Pendekatan ini benar karena syncTokenFromUrl() di main.ts sudah
 * menjamin token tersimpan di memory SEBELUM guard ini berjalan.
 * Tidak perlu authStore karena user belum ter-fetch saat guard pertama kali jalan.
 */
router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth) {
    try {
      const user = await authApi.getMe();
      // Kunci proteksi: Hanya akun dengan role OWNER yang diizinkan masuk
      if (user && user.role === 'OWNER') {
        next();
      } else {
        // Role tidak sesuai → arahkan ke halaman forbidden SSO
        window.location.href = import.meta.env.VITE_SSO_URL
          ? `${import.meta.env.VITE_SSO_URL}/forbidden`
          : 'http://localhost:5173/forbidden';
      }
    } catch (err) {
      // Token tidak valid / expired → kembalikan ke login SSO
      window.location.href = import.meta.env.VITE_SSO_URL || 'http://localhost:5173';
    }
  } else {
    next();
  }
});

export default router;