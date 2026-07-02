// packages/frontend/services/http.ts
import axios from 'axios';
import { getToken, setToken, getRefreshToken, setRefreshToken, clearToken, clearRefreshToken } from './token';

// Instance utama Axios menuju API Gateway
export const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request Interceptor: Sisipkan Access Token
apiClient.interceptors.request.use(
    (config) => {
        const token = getToken();
        if (token && config.headers) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Response Interceptor: Tangani 401 dan Refresh Token Flow
let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });
    failedQueue = [];
};

apiClient.interceptors.response.use(
    (response) => {
        // Unwrap otomatis jika Go backend membungkus respons dengan struktur { success, message, data }
        if (response.data && typeof response.data === 'object' && 'success' in response.data && 'data' in response.data) {
            response.data = response.data.data;
        }
        return response;
    },
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            if (isRefreshing) {
                return new Promise(function (resolve, reject) {
                    failedQueue.push({ resolve, reject });
                }).then(token => {
                    originalRequest.headers.Authorization = 'Bearer ' + token;
                    return apiClient(originalRequest);
                }).catch(err => Promise.reject(err));
            }

            originalRequest._retry = true;
            isRefreshing = true;

            const refreshToken = getRefreshToken();
            if (!refreshToken) {
                // Logout user jika tidak ada refresh token
                clearToken();
                clearRefreshToken();
                const ssoUrl = import.meta.env.VITE_SSO_URL || 'http://localhost:5173';
                const currentUrl = encodeURIComponent(window.location.href);
                window.location.href = `${ssoUrl}?redirect=${currentUrl}`;
                return Promise.reject(error);
            }

            try {
                const { data } = await axios.post(`${apiClient.defaults.baseURL}/api/auth/refresh-token`, {
                    refresh_token: refreshToken
                });

                const actualData = data.data ? data.data : data;
                const newAccessToken = actualData.access_token;
                setToken(newAccessToken);

                if (actualData.refresh_token) {
                    setRefreshToken(actualData.refresh_token);
                }

                processQueue(null, newAccessToken);
                originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
                return apiClient(originalRequest);
            } catch (err) {
                processQueue(err, null);
                clearToken();
                clearRefreshToken();
                const ssoUrl = import.meta.env.VITE_SSO_URL || 'http://localhost:5173';
                const currentUrl = encodeURIComponent(window.location.href);
                window.location.href = `${ssoUrl}?redirect=${currentUrl}`;
                return Promise.reject(err);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);