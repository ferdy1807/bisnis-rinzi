import { createRouter, createWebHistory } from 'vue-router';
import { authApi } from '@frontend/api/auth';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Root redirect
    {
      path: '/',
      redirect: '/dashboard'
    },

    // --- 1. DASHBOARD ---
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },

    // --- 2. POS & KASIR ---
    {
      path: '/pos',
      name: 'pos',
      component: () => import('../views/PosView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import('../views/RiwayatTransaksi.vue'),
      meta: { requiresAuth: true }
    },

    // --- 3. SHIFT & LAPORAN ---
    {
      path: '/shifts',
      name: 'shifts',
      component: () => import('../views/ShiftView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/daily-report',
      name: 'daily-report',
      component: () => import('../views/DailyReport.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/report/shift',
      name: 'shift-report',
      component: () => import('../views/ShiftReport.vue'),
      meta: { requiresAuth: true }
    },

    // --- 4. MASTER BARANG ---
    {
      path: '/products',
      name: 'products',
      component: () => import('../views/ProductView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/products/detail/:id',
      name: 'product-detail',
      component: () => import('../views/ProductDetailView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import('../views/CategoryView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/brands',
      name: 'brands',
      component: () => import('../views/BrandView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/units',
      name: 'units',
      component: () => import('../views/UnitView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/stock-movement',
      name: 'stock-movement',
      component: () => import('../views/StockMovement.vue'),
      meta: { requiresAuth: true }
    },

    // --- 5. KEUANGAN ---
    {
      path: '/cash-in',
      name: 'cash-in',
      component: () => import('../views/KasMasuk.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/cash-out',
      name: 'cash-out',
      component: () => import('../views/KasKeluar.vue'),
      meta: { requiresAuth: true }
    },

    // --- CATCH-ALL: Redirect ke dashboard ---
    {
      path: '/:pathMatch(.*)*',
      redirect: '/dashboard'
    }
  ],
});

// Navigation Guard Penegak Validitas Akses ACL Multi-Schema
router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth) {
    try {
      const user = await authApi.getMe();
      // Kunci proteksi: Hanya akun dengan role CASHIER yang diizinkan masuk ke portal toko ini
      if (user && user.role === 'CASHIER') {
        next();
      } else {
        // Alihkan ke halaman SSO eksternal jika role tidak sesuai aturan regulasi
        window.location.href = import.meta.env.VITE_SSO_URL || 'http://localhost:5173/forbidden';
      }
    } catch (err) {
      // Bersihkan sesi dan lempar balik ke gerbang Login SSO jika token kedaluwarsa/cacat
      window.location.href = import.meta.env.VITE_SSO_URL || 'http://localhost:5173/login';
    }
  } else {
    next();
  }
});

export default router;