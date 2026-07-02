package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

// InventoryOutboxEvent merepresentasikan struktur baris data outbox yang tersimpan di dalam tabel outbox_events pada inventory_db.
type InventoryOutboxEvent struct {
	ID            string    `json:"id"`
	AggregateType string    `json:"aggregate_type"` // Contoh: "INVENTORY_STOCK" atau "STOCK_MOVEMENT"
	AggregateID   string    `json:"aggregate_id"`   // Berisi ProductID atau ReferenceCode PO Supplier
	EventType     string    `json:"event_type"`     // Contoh: "INCOMING_SUPPLIER" atau "STOCK_ADJUSTED"
	Payload       []byte    `json:"payload"`        // Data biner JSONB mutasi stok logistik gudang
	Status        string    `json:"status"`         // "PENDING", "SUCCESS", "FAILED"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// LocalStockMutationPayload mencerminkan bentuk fisik payload JSON asli yang disimpan oleh inventory_service saat terjadi pasokan barang masuk.
type LocalStockMutationPayload struct {
	ProductID     string    `json:"product_id"`
	MutationType  string    `json:"mutation_type"`  // "INCOMING_SUPPLIER" atau "STOCK_ADJUSTMENT"
	Quantity      float64   `json:"qty"`            // Jumlah kuantitas fisik barang dagangan yang dimutasi
	UnitCost      float64   `json:"unit_cost"`      // Harga beli per satuan barang dari supplier (Nilai Pokok Modal)
	TotalCost     float64   `json:"total_cost"`     // Total pengeluaran modal (Kuantitas * UnitCost) untuk aset persediaan
	ReferenceCode string    `json:"reference_code"` // Nomor Purchase Order (PO) Supplier atau dokumen surat jalan gudang
	Timestamp     time.Time `json:"timestamp"`
}

// UnmarshalLocalStock mengurai slice byte mentah pangkalan data lokal menjadi bentuk struct mutasi stok gudang yang valid.
func UnmarshalLocalStock(raw []byte) (*LocalStockMutationPayload, error) {
	var p LocalStockMutationPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("inventory_entity: gagal mengekstrak payload mutasi stok lokal: %w", err)
	}
	return &p, nil
}
