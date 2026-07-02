// apps/login-sso/src/main.ts
import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';

// Mengimpor Router dan Guard
import { router } from './router';
import { setupRouterGuard } from './router/guard';

// Mengimpor fondasi dari packages/frontend
import { setupPrimeVue } from '@frontend/plugins/primevue';
import '@frontend/styles/main.css';
import 'primeicons/primeicons.css'; // Icon PrimeVue

const app = createApp(App);
const pinia = createPinia();

// Urutan registrasi sangat penting
app.use(pinia);
app.use(router);

// Mengaktifkan pengecekan hak akses SSO
setupRouterGuard(router);

// Inisialisasi tema global PrimeVue + Tailwind
setupPrimeVue(app);

app.mount('#app');