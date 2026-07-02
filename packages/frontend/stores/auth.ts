import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { authApi } from '../api/auth';
import {
    setToken,
    setRefreshToken,
    clearToken,
    clearRefreshToken,
    getToken,
    getRefreshToken
} from '../services/token';
import type { LoginRequest, User } from '../types/auth';

export const useAuthStore = defineStore('auth', () => {
    // State global untuk menyimpan data sesi pengguna aktif
    const user = ref<User | null>(null);
    const token = ref<string | null>(getToken());
    const refreshToken = ref<string | null>(getRefreshToken());
    const loading = ref<boolean>(false);

    // Getters untuk validasi status login dan peran sistem
    const isAuthenticated = computed(() => !!token.value);
    const role = computed(() => user.value?.role || null);

    /**
     * Mengirimkan permintaan login ke backend melalui API Gateway
     */
    async function login(payload: LoginRequest) {
        loading.value = true;
        try {
            const response = await authApi.login(payload);

            // Bypass type-checking TypeScript dengan 'any' untuk mengekstrak pembungkus .data (jika ada)
            const responseData: any = response;
            const actualData = responseData.data ? responseData.data : responseData;

            token.value = actualData.access_token;
            refreshToken.value = actualData.refresh_token;

            if (!token.value) {
                throw new Error("Token tidak ditemukan dalam respons backend.");
            }

            // Pastikan nilainya bukan null sebelum disimpan ke localStorage
            if (token.value) setToken(token.value);
            if (refreshToken.value) setRefreshToken(refreshToken.value);

            await fetchMe();
        } catch (error) {
            clearAuth();
            throw error;
        } finally {
            loading.value = false;
        }
    }

    /**
     * Mengambil data profil terperinci milik pengguna yang sedang aktif
     */
    async function fetchMe() {
        try {
            const userData: any = await authApi.getMe();
            user.value = userData.data ? userData.data : userData;
        } catch (error) {
            clearAuth();
            throw error;
        }
    }

    /**
     * Mengakhiri sesi pengguna dan membersihkan data autentikasi
     */
    async function logout() {
        loading.value = true;
        try {
            await authApi.logout();
        } catch (error) {
            console.error('Logout API error:', error);
        } finally {
            clearAuth();
            loading.value = false;
            // Mengarahkan kembali ke halaman gerbang login SSO
            window.location.href = 'http://localhost:5173/';
        }
    }

    /**
     * Menghapus seluruh data sesi dari state internal dan storage perangkat
     */
    function clearAuth() {
        user.value = null;
        token.value = null;
        refreshToken.value = null;
        clearToken();
        clearRefreshToken();
    }

    /**
     * Memeriksa apakah pengguna memiliki salah satu peran yang diizinkan
     */
    function hasRole(allowedRoles: string | string[]): boolean {
        if (!user.value) return false;
        if (Array.isArray(allowedRoles)) {
            return allowedRoles.includes(user.value.role);
        }
        return user.value.role === allowedRoles;
    }

    /**
     * Memeriksa hak akses atau batasan permission tertentu di dalam sistem
     */
    function hasPermission(permission: string): boolean {
        // Sesuai aturan, peran OWNER memiliki kendali mutlak atas seluruh fitur
        if (user.value?.role === 'OWNER') return true;
        return false;
    }

    /**
     * Mengambil token dari URL (untuk SSO lintas port lokal) dan memuat user
     */
    async function initializeSession() {
        const urlParams = new URLSearchParams(window.location.search);
        const urlAccess = urlParams.get('accessToken');
        const urlRefresh = urlParams.get('refreshToken');

        if (urlAccess && urlRefresh) {
            token.value = urlAccess;
            refreshToken.value = urlRefresh;
            setToken(urlAccess);
            setRefreshToken(urlRefresh);

            // Bersihkan URL agar token tidak terlihat
            window.history.replaceState({}, document.title, window.location.pathname);
        }

        if (token.value || getRefreshToken()) {
            try {
                await fetchMe();
            } catch (e) {
                // Token tidak valid atau kedaluwarsa
                clearAuth();
            }
        }
    }

    return {
        user,
        token,
        refreshToken,
        loading,
        isAuthenticated,
        role,
        login,
        logout,
        fetchMe,
        hasRole,
        hasPermission,
        clearAuth,
        initializeSession
    };
});