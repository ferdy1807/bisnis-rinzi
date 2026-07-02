import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import { syncTokenFromUrl } from '@frontend/services/token';

import { setupPrimeVue } from '@frontend/plugins/primevue';
import '@frontend/styles/main.css';
import 'primeicons/primeicons.css';

// Serap token SSO dari query param (?accessToken=...) SEBELUM router guard berjalan
syncTokenFromUrl();

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

setupPrimeVue(app);

app.mount('#app');