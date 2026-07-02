// packages/frontend/api/cash.ts
import { BaseApi } from './base';
import type { CashierSession, ShiftSummary, CashTransaction, ExpenseCategory, Expense } from '@frontend/types/cash';

export class CashApi extends BaseApi {
    public async getCurrentSession(): Promise<CashierSession | null> {
        const response: any = await this.http.get('/api/cash/shifts/current');
        if (response.data?.data === null) return null;
        return response.data?.data ?? response.data ?? null;
    }

    public async openSession(openingCash: number): Promise<CashierSession> {
        const response: any = await this.http.post('/api/cash/shifts/open', {
            opening_cash: openingCash
        });
        return response.data?.data ?? response.data;
    }

    public async closeSession(actualCash: number): Promise<CashierSession> {
        const response: any = await this.http.post('/api/cash/shifts/close', {
            actual_cash: actualCash
        });
        return response.data?.data ?? response.data;
    }

    public async getShiftSummary(sessionId: string): Promise<ShiftSummary> {
        const response: any = await this.http.get(`/api/cash/shifts/${sessionId}/summary`);
        return response.data?.data ?? response.data;
    }

    public async uploadShiftReport(sessionId: string, pdfBlob: Blob, filename: string): Promise<any> {
        const formData = new FormData();
        formData.append('report', pdfBlob, filename);

        const response = await this.http.post(`/api/cash/shifts/${sessionId}/report`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        return response.data;
    }

    public async getShifts(): Promise<CashierSession[]> {
        const response: any = await this.http.get('/api/cash/shifts');
        return response.data?.data ?? response.data ?? [];
    }

    public async getTransactions(): Promise<CashTransaction[]> {
        const response: any = await this.http.get('/api/cash/transactions');
        return response.data?.data ?? response.data ?? [];
    }

    public async createTransaction(payload: { transaction_type: string, reference_type?: string, amount: number, notes: string }): Promise<CashTransaction> {
        const response: any = await this.http.post('/api/cash/transactions', payload);
        return response.data?.data ?? response.data;
    }

    public async getInternalIncomes(): Promise<any[]> {
        const response: any = await this.http.get('/api/cash/internal-incomes');
        return response.data?.data ?? response.data ?? [];
    }

    public async createInternalIncome(payload: any): Promise<any> {
        const response: any = await this.http.post('/api/cash/internal-incomes', payload);
        return response.data?.data ?? response.data;
    }

    public async getExpenseCategories(): Promise<ExpenseCategory[]> {
        const response: any = await this.http.get('/api/cash/expense-categories');
        const data = response.data?.data ?? response.data ?? [];
        return data.filter((c: any) => !c.code.toUpperCase().startsWith('INC-'));
    }

    public async getExpenses(): Promise<Expense[]> {
        const response: any = await this.http.get('/api/cash/expenses');
        return response.data?.data ?? response.data ?? [];
    }

    public async createExpense(payload: { category_id: string, amount: number, description: string }): Promise<Expense> {
        const response: any = await this.http.post('/api/cash/expenses', payload);
        return response.data?.data ?? response.data;
    }

    public async createExpenseCategory(payload: { code: string; name: string }): Promise<any> {
        const response: any = await this.http.post('/api/cash/expense-categories', payload);
        return response.data;
    }

    public async getOtherIncomeCategories(): Promise<any[]> {
        const response: any = await this.http.get('/api/cash/expense-categories');
        const data = response.data?.data ?? response.data ?? [];
        return data.filter((c: any) => c.code.toUpperCase().startsWith('INC-'));
    }

    public async createOtherIncomeCategory(payload: { code: string; name: string }): Promise<any> {
        // Reuse the expense-categories endpoint to store income categories
        const response: any = await this.http.post('/api/cash/expense-categories', payload);
        return response.data;
    }

    public async updateExpenseCategory(id: string, payload: { name: string; description?: string }): Promise<any> {
        const response: any = await this.http.put(`/api/cash/expense-categories/${id}`, payload);
        return response.data;
    }

    public async deleteExpenseCategory(id: string): Promise<any> {
        const response: any = await this.http.delete(`/api/cash/expense-categories/${id}`);
        return response.data;
    }
}

export const cashApi = new CashApi();
