// apps/login-sso/src/router/guard.ts
import { Router } from 'vue-router';
import { useAuthStore } from '@frontend/stores/auth';

export function setupRouterGuard(router: Router) {
    router.beforeEach(async (to, from, next) => {
        const authStore = useAuthStore();

        // Jika pengguna sudah memiliki token dan mencoba mengakses halaman login
        if (authStore.isAuthenticated && to.name === 'Login') {
            const role = authStore.role;
            const urlParams = `?accessToken=${authStore.token}&refreshToken=${authStore.refreshToken}`;

            // Redirect berdasarkan role menggunakan environment variables
            if (role === 'OWNER') {
                window.location.href = (import.meta.env.VITE_ADMIN_URL || 'http://localhost:5176') + urlParams;
            } else if (role === 'CASHIER') {
                window.location.href = (import.meta.env.VITE_TOKO_URL || 'http://localhost:5175') + urlParams;
            } else if (role === 'PEGAWAI') {
                window.location.href = (import.meta.env.VITE_SEWA_URL || 'http://localhost:5174') + urlParams;
            } else {
                next();
            }
            return;
        }

        next();
    });
}