// apps/portal-toko/src/stores/finance.ts
import { defineStore } from 'pinia';
import { financeApi } from '@frontend/api/finance';
import type { AccountingPeriod, CorporateDashboardAnalytics } from '@frontend/types/finance';

export const useFinanceStore = defineStore('finance', {
    state: () => ({
        activePeriod: null as AccountingPeriod | null,
        dashboardData: null as CorporateDashboardAnalytics | null,
        loading: false,
        startDate: '2026-01-01', // Default sesuai kebutuhan audit 2026
        endDate: '2026-12-31'
    }),

    actions: {
        async fetchDashboardSummary() {
            this.loading = true;
            try {
                this.dashboardData = await financeApi.getDashboardSummary();
            } finally {
                this.loading = false;
            }
        },

        async loadActivePeriod() {
            const periods = await financeApi.getPeriods();
            // Ambil periode yang is_closed = false
            this.activePeriod = periods.find(p => !p.is_closed) || null;
        },

        setPeriodRange(start: string, end: string) {
            this.startDate = start;
            this.endDate = end;
        }
    }
});