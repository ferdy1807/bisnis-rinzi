// apps/portal-sewa/src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import { authApi } from '@frontend/api/auth';
import HomeView from '../views/HomeView.vue'; // Kanban Board

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // --- 1. DASHBOARD ---
    {
      path: '/',
      name: 'dashboard-kanban',
      component: HomeView,
      meta: { requiresAuth: true }
    },
    {
      path: '/overdue',
      name: 'overdue-reservations',
      component: () => import('../views/OverdueView.vue'),
      meta: { requiresAuth: true }
    },

    // --- 2. PEMESANAN & INVOICE ---
    {
      path: '/upcoming',
      name: 'upcoming-reservations',
      component: () => import('../views/UpcomingView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/reservations-create',
      name: 'create-reservation',
      component: () => import('../views/CreateReservationView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/reservations-list',
      name: 'reservation-list',
      component: () => import('../views/ReservationListview.vue'),
      meta: { requiresAuth: true }
    },

    // --- 3. ALUR KERJA (WORKFLOW) ---
    {
      path: '/workflow-contents-received',
      name: 'workflow-contents-received',
      component: () => import('../views/ContentsReceivedView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/returns',
      name: 'returns-and-penalties',
      component: () => import('../views/ReturnsView.vue'),
      meta: { requiresAuth: true }
    },

    // --- 4. INFORMASI KATALOG ---
    {
      path: '/availability-calendar',
      name: 'availability-calendar',
      component: () => import('../views/AvailabilityCalendarView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/catalog',
      name: 'catalog-readonly',
      component: () => import('../views/CatalogView.vue'),
      meta: { requiresAuth: true }
    },

    // --- CATCH-ALL: Redirect ke dashboard ---
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
});

// Navigation Guard Penegak Validitas Akses ACL Multi-Schema
router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth) {
    try {
      const user = await authApi.getMe();
      // Kunci proteksi: Hanya akun dengan role PEGAWAI yang diizinkan masuk ke portal ini
      if (user && user.role === 'PEGAWAI') {
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