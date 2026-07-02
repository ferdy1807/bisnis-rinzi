// packages/frontend/services/token.ts

let accessTokenMemory: string | null = null;

export const getToken = () => accessTokenMemory;
export const setToken = (token: string) => { accessTokenMemory = token; };
export const clearToken = () => { accessTokenMemory = null; };

export const getRefreshToken = () => localStorage.getItem('refresh_token');
export const setRefreshToken = (token: string) => localStorage.setItem('refresh_token', token);
export const clearRefreshToken = () => localStorage.removeItem('refresh_token');

/**
 * Menangkap token dari URL dan membersihkannya agar URL kembali bersih
 */
export const syncTokenFromUrl = () => {
    const urlParams = new URLSearchParams(window.location.search);
    const access = urlParams.get('accessToken');
    const refresh = urlParams.get('refreshToken');

    if (access) setToken(access);
    if (refresh) setRefreshToken(refresh);

    // Bersihkan URL dari token agar tidak terekspos di history browser
    if (access || refresh) {
        window.history.replaceState({}, document.title, window.location.pathname);
    }
};