export interface OutboxEvent {
    id: string;
    aggregate_type: string;
    aggregate_id: string;
    event_type: string;
    payload: Record<string, any>;
    status: 'PENDING' | 'PROCESSED' | 'FAILED';
    error_message?: string;
    created_at?: string;
    updated_at?: string;
    processed_at?: string;
}

export interface SyncVersion {
    id: string | number; // bigserial
    entity_type: string;
    entity_id: string;
    operation: 'INSERT' | 'UPDATE' | 'DELETE';
    version_number: number;
    changed_at?: string;
}
