// packages/frontend/api/base.ts
import { apiClient } from '../services/http';
import type { AxiosInstance } from 'axios';

export abstract class BaseApi {
    protected http: AxiosInstance;

    constructor() {
        this.http = apiClient;
    }
}