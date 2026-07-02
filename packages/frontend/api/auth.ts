// packages/frontend/api/auth.ts
import { BaseApi } from './base';
import { getRefreshToken } from '../services/token';
import type { LoginRequest, LoginResponse, User, AuditLog, Role } from '../types/auth'; // Asumsi tipe data sudah dibuat nanti

export class AuthApi extends BaseApi {
    public async login(data: LoginRequest): Promise<LoginResponse> {
        const response = await this.http.post<LoginResponse>('/api/auth/login', data);
        return response.data;
    }

    public async createUser(data: { username: string; password?: string; full_name: string; role: string }): Promise<User> {
        const response: any = await this.http.post('/api/auth/register', data);
        return response.data?.data ?? response.data;
    }

    public async logout(): Promise<void> {
        const refresh_token = getRefreshToken();
        await this.http.post('/api/auth/logout', { refresh_token });
    }

    public async getMe(): Promise<User> {
        const response = await this.http.get<User>('/api/auth/me');
        return response.data;
    }

    public async updateProfile(userId: string, data: { full_name?: string, role?: string }): Promise<void> {
        await this.http.put(`/api/auth/users/${userId}`, data);
    }

    public async updatePassword(userId: string, data: { old_password?: string, new_password?: string }): Promise<void> {
        await this.http.put(`/api/auth/users/${userId}/password`, data);
    }

    public async register(data: { username: string; password: string; full_name: string; role: string }): Promise<void> {
        await this.http.post('/api/auth/register', { ...data, is_active: false });
    }

    public async deleteUser(userId: string): Promise<void> {
        await this.http.delete(`/api/auth/users/${userId}`);
    }

    public async forceLogoutUser(userId: string): Promise<void> {
        await this.http.delete(`/api/auth/users/${userId}/sessions`);
    }

    public async getUsers(): Promise<User[]> {
        const response: any = await this.http.get('/api/auth/users');
        if (Array.isArray(response.data)) return response.data;
        return response.data?.data || [];
    }

    public async getRoles(): Promise<Role[]> {
        const response: any = await this.http.get('/api/auth/roles');
        if (Array.isArray(response.data)) return response.data;
        return response.data?.data || [];
    }

    public async getAuditLogs(): Promise<AuditLog[]> {
        const response: any = await this.http.get('/api/auth/audit-logs');
        if (Array.isArray(response.data)) return response.data;
        return response.data?.data || [];
    }
}

export const authApi = new AuthApi();