export interface Role {
    code: string;
    name: string;
    dashboard_url: string;
    created_at?: string;
    updated_at?: string;
}

export interface User {
    id: string;
    username: string;
    full_name: string;
    role: string;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
    deleted_at?: string;
}

export interface LoginRequest {
    username: string;
    password?: string;
}

export interface LoginResponse {
    access_token: string;
    refresh_token: string;
    expires_in?: number;
}

export interface AuditLog {
    id: string;
    user_id: string;
    user_name?: string;
    action?: string;
    entity_name?: string;
    entity_id?: string;
    old_data?: Record<string, any>;
    new_data?: Record<string, any>;
    created_at?: string;
    updated_at?: string;
}

export interface RefreshToken {
    id: string;
    user_id: string;
    token: string;
    expires_at: string;
    device_info?: string;
    ip_address?: string;
    created_at?: string;
    updated_at?: string;
    revoked_at?: string;
}