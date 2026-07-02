package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

// CashOutboxEvent merepresentasikan baris arsip dari tabel outbox_events di dalam skema database cash_db.
type CashOutboxEvent struct {
	ID            string    `json:"id"`
	AggregateType string    `json:"aggregate_type"` // Contoh: "CASH_SESSION" atau "EXPENSE"
	AggregateID   string    `json:"aggregate_id"`   // Berisi SessionID atau ExpenseID
	EventType     string    `json:"event_type"`     // Contoh: "CASHIER_SHIFT_CLOSED" atau "EXPENSE_RECORDED"
	Payload       []byte    `json:"payload"`        // Data biner JSONB asli dari database lokal
	Status        string    `json:"status"`         // "PENDING", "SUCCESS", "FAILED"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// LocalCashShiftPayload mencerminkan bentuk fisik payload JSON asli yang disimpan oleh cash_service saat kasir menutup laci.
type LocalCashShiftPayload struct {
	CashierID       string    `json:"cashier_id"`
	SessionID       string    `json:"session_id"`
	OpeningCash     float64   `json:"opening_cash"`    // Saldo Awal Laci
	TotalIncome     float64   `json:"total_income"`    // Total Penerimaan POS
	ManualDeposit   float64   `json:"manual_deposit"`  // Kas Masuk Ekstra (Manual)
	TotalExpense    float64   `json:"total_expense"`   // Total Kas Keluar (Expenses)
	ActualCash      float64   `json:"actual_cash"`     // Uang fisik yang dihitung nyata di laci
	ExpectedCash    float64   `json:"expected_cash"`   // Uang sistem yang seharusnya ada berdasarkan transaksi POS
	DifferenceCash  float64   `json:"difference_cash"` // Nilai selisih tekor/surplus (e.g., -1000)
	ClosedTimestamp time.Time `json:"closed_timestamp"`
}

// LocalExpensePayload memetakan pengeluaran kas operasional toko (misal: beli sapu, air galon) untuk kebutuhan jurnal beban keuangan.
type LocalExpensePayload struct {
	ExpenseID   string    `json:"expense_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // Contoh: "PERLENGKAPAN_TOKO", "BIAYA_OPERASIONAL"
	Timestamp   time.Time `json:"timestamp"`
}

// UnmarshalLocalShift mengurai slice byte mentah pangkalan data lokal menjadi bentuk struct shift kasir yang valid.
func UnmarshalLocalShift(raw []byte) (*LocalCashShiftPayload, error) {
	var p LocalCashShiftPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("cash_entity: gagal mengekstrak data shift lokal: %w", err)
	}
	return &p, nil
}

// UnmarshalLocalExpense mengurai slice byte mentah menjadi objek struct pengeluaran biaya tunai retail.
func UnmarshalLocalExpense(raw []byte) (*LocalExpensePayload, error) {
	var p LocalExpensePayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("cash_entity: gagal mengekstrak data biaya lokal: %w", err)
	}
	return &p, nil
}
