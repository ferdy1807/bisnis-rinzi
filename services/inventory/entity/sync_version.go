package entity

import "time"

// SyncVersion merepresentasikan tabel 'sync_versions' untuk manajemen sinkronisasi data PWA Offline
type SyncVersion struct {
	ID            int64     `json:"id"`             // BIGSERIAL
	EntityType    string    `json:"entity_type"`    // 'product', 'category', 'brand', 'unit'
	EntityID      string    `json:"entity_id"`      // UUID entitas yang berubah
	Operation     string    `json:"operation"`      // 'INSERT', 'UPDATE', 'DELETE'
	VersionNumber int64     `json:"version_number"` // Sequence monotonik global
	ChangedAt     time.Time `json:"changed_at"`
}
