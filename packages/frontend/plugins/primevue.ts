// packages/frontend/plugins/primevue.ts
import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';
import type { App } from 'vue';

export const setupPrimeVue = (app: App) => {
    app.use(PrimeVue, {
        theme: {
            preset: Aura,
            options: {
                darkModeSelector: '.dark',
                cssLayer: {
                    name: 'primevue',
                    order: 'tailwind-base, primevue, tailwind-utilities'
                }
            }
        },
        ripple: true
    });

    // Mendaftarkan service PrimeVue yang sering digunakan
    app.use(ToastService);
    app.use(ConfirmationService);
};