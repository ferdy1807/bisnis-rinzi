// packages/frontend/api/finance.ts
import { BaseApi } from './base';
import type {
    ChartOfAccount,
    AccountingPeriod,
    DailyClosing,
    ReconciliationLog,
    JournalEntry,
    LedgerReportResponse,
    IncomeStatementResponse,
    BalanceSheetResponse,
    CorporateDashboardAnalytics
} from '../types/finance';

export class FinanceApi extends BaseApi {
    // --- 1. CORPORATE ANALYTICS DASHBOARD ---
    public async getDashboardSummary(): Promise<CorporateDashboardAnalytics> {
        const response: any = await this.http.get('/api/finance/analytics/dashboard');
        return response.data?.data ?? response.data ?? {};
    }

    // --- 2. CHART OF ACCOUNTS (COA) ---
    public async getAccounts(): Promise<ChartOfAccount[]> {
        const response: any = await this.http.get('/api/finance/accounts');
        if (Array.isArray(response.data)) return response.data;
        if (response.data && Array.isArray(response.data.data)) return response.data.data;
        return [];
    }

    public async createAccount(payload: {
        account_code: string;
        account_name: string;
        account_group: string;
        normal_balance: 'DEBIT' | 'CREDIT';
    }): Promise<ChartOfAccount> {
        const response: any = await this.http.post('/api/finance/accounts', payload);
        return response.data?.data ?? response.data;
    }

    public async getAccountById(id: string): Promise<ChartOfAccount> {
        const response: any = await this.http.get(`/api/finance/accounts/${id}`);
        return response.data?.data ?? response.data;
    }

    public async updateAccount(id: string, payload: Partial<ChartOfAccount>): Promise<any> {
        const response: any = await this.http.put(`/api/finance/accounts/${id}`, payload);
        return response.data;
    }

    public async deleteAccount(id: string): Promise<any> {
        const response: any = await this.http.delete(`/api/finance/accounts/${id}`);
        return response.data;
    }

    // --- 3. ACCOUNTING PERIODS & LOCKS ---
    public async getPeriods(): Promise<AccountingPeriod[]> {
        const response: any = await this.http.get('/api/finance/accounting-periods');
        if (Array.isArray(response.data)) return response.data;
        if (response.data && Array.isArray(response.data.data)) return response.data.data;
        return [];
    }

    public async createPeriod(payload: { name: string; start_date: string; end_date: string }): Promise<any> {
        const response: any = await this.http.post('/api/finance/accounting-periods', payload);
        return response.data;
    }

    public async lockPeriod(periodId: string, reason: string): Promise<any> {
        const response: any = await this.http.post('/api/finance/period-locks', {
            period_id: periodId,
            lock_reason: reason
        });
        return response.data;
    }

    // --- 4. DAILY CLOSINGS & RECONCILIATION ---
    public async getDailyClosings(): Promise<DailyClosing[]> {
        const response: any = await this.http.get('/api/finance/daily-closings');
        if (Array.isArray(response.data)) return response.data;
        if (response.data && Array.isArray(response.data.data)) return response.data.data;
        return [];
    }

    public async getClosingDetail(id: string): Promise<DailyClosing> {
        const response: any = await this.http.get(`/api/finance/daily-closings/${id}`);
        return response.data?.data ?? response.data;
    }

    public async getReconciliationLogs(): Promise<ReconciliationLog[]> {
        const response: any = await this.http.get('/api/finance/reconciliation/logs');
        if (Array.isArray(response.data)) return response.data;
        if (response.data && Array.isArray(response.data.data)) return response.data.data;
        return [];
    }

    // --- 5. JOURNAL ENTRIES (GENERAL LEDGER) ---
    public async getJournalEntries(): Promise<JournalEntry[]> {
        const response: any = await this.http.get('/api/finance/journal-entries');
        if (Array.isArray(response.data)) return response.data;
        if (response.data && Array.isArray(response.data.data)) return response.data.data;
        return [];
    }

    public async createManualJournal(payload: {
        journal_id: string;
        reference_number: string;
        narration: string;
        details: { account_id: string; debit_amount: number; credit_amount: number }[];
    }): Promise<any> {
        const response: any = await this.http.post('/api/finance/journal-entries', payload);
        return response.data;
    }

    // --- 6. FINANCIAL STANDARD REPORTS ---
    public async getGeneralLedgerReport(accountCode: string, startDate: string, endDate: string): Promise<LedgerReportResponse> {
        const response: any = await this.http.get('/api/finance/reports/general-ledger', {
            params: { account_code: accountCode, start_date: startDate, end_date: endDate }
        });
        return response.data?.data ?? response.data;
    }

    public async getProfitLossReport(): Promise<IncomeStatementResponse> {
        const response: any = await this.http.get('/api/finance/reports/profit-loss');
        return response.data?.data ?? response.data;
    }

    public async getCashFlowReport(): Promise<any> {
        const response: any = await this.http.get('/api/finance/reports/cash-flow');
        return response.data?.data ?? response.data;
    }

    public async getBalanceSheetReport(): Promise<BalanceSheetResponse> {
        const response: any = await this.http.get('/api/finance/reports/balance-sheet');
        return response.data?.data ?? response.data;
    }

    // --- 7. RENTAL DAMAGES ---
    public async getRentalDamages(): Promise<any[]> {
        const response: any = await this.http.get('/api/rental/returns/damages');
        if (Array.isArray(response.data)) return response.data;
        if (response.data && Array.isArray(response.data.data)) return response.data.data;
        return [];
    }

    public async settleRentalDamage(id: string, payload: { payment_action: string; audit_notes: string }): Promise<any> {
        const response: any = await this.http.post(`/api/rental/returns/damages/${id}/settle`, payload);
        return response.data?.data ?? response.data;
    }

    // --- 8. DYNAMIC ANALYTICS DISPATCHER ---
    public async getAnalyticsData(action: string): Promise<any> {
        const response: any = await this.http.get(`/api/finance/analytics/${action}`);
        return response.data?.data ?? response.data;
    }

    // Helper Generator URL Unduh Dokumen Resmi
    public getExportUrl(format: 'pdf' | 'excel' | 'csv'): string {
        const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
        return `${baseUrl}/api/finance/reports/export/${format}`;
    }

    public async uploadDailyIncomeReport(date: string, pdfBlob: Blob, filename: string): Promise<any> {
        const formData = new FormData();
        formData.append('pdf', pdfBlob, filename);
        return this.http.post('/api/finance/daily-incomes/upload-pdf', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
    }

    // --- 9. DAILY INCOMES (AGGREGATED FROM JOURNALS) ---
    public async getDailyIncomes(): Promise<any[]> {
        const [journals, accounts] = await Promise.all([
            this.getJournalEntries(),
            this.getAccounts()
        ]);
        
        const accountMap = new Map();
        accounts.forEach(acc => accountMap.set(acc.id, acc.account_code));

        const incomes: any[] = [];
        
        journals.forEach((entry: any) => {
            let totalIncome = 0;
            
            if (entry.details && Array.isArray(entry.details)) {
                entry.details.forEach((detail: any) => {
                    const accCode = accountMap.get(detail.account_id) || '';
                    // 410000 = Pendapatan Retail, 420000 = Pendapatan Sewa, 421000 = Pendapatan Denda Sewa
                    if (detail.credit_amount > 0 && (accCode === '410000' || accCode === '420000' || accCode === '421000')) {
                       totalIncome += detail.credit_amount;
                    }
                });
            }
            
            if (totalIncome > 0) {
                incomes.push({
                    id: entry.id,
                    transaction_type: 'DEPOSIT',
                    reference_type: 'JOURNAL',
                    reference_id: entry.reference_number,
                    amount: totalIncome,
                    notes: entry.narration || '-',
                    created_by: 'Finance System',
                    created_at: entry.entry_date
                });
            }
        });
        
        return incomes;
    }
}

export const financeApi = new FinanceApi();