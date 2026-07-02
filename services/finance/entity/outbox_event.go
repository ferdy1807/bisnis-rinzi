package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

// FinanceInboxEvent merepresentasikan struktur data event dari layanan luar
// yang masuk ke dalam antrean pemrosesan jurnal keuangan.
type FinanceInboxEvent struct {
	ID            string    `json:"id"`
	AggregateType string    `json:"aggregate_type"` // Contoh: "RETAIL_ORDER", "SHIFT_SESSION", "RENTAL_RESERVATION", "INVENTORY_STOCK"
	AggregateID   string    `json:"aggregate_id"`   // ID Dokumen / Invoice asal
	EventType     string    `json:"event_type"`     // Contoh: "ORDER_COMPLETED", "SHIFT_CLOSED", "RENTAL_CONFIRMED", "PRODUCT_RETURN_PROCESSED"
	Payload       []byte    `json:"payload"`        // Data mentah terkompresi JSON
	CreatedAt     time.Time `json:"created_at"`
}

// =========================================================================
// 🛒 DEFENISI PAYLOAD OPERASIONAL: PORTAL-TOKO (Cash, POS, Inventory)
// =========================================================================

// StoreShiftPayload menangkap pembukaan/penutupan laci kasir untuk penjurnalan saldo kas
type StoreShiftPayload struct {
	SessionID     string    `json:"session_id"`
	UserKaryawan  string    `json:"user_id"`
	OpeningCash   float64   `json:"opening_cash"`
	TotalIncome   float64   `json:"total_income"`    // Total Penerimaan POS
	ManualDeposit float64   `json:"manual_deposit"`  // Kas Masuk Ekstra (Manual)
	TotalExpense  float64   `json:"total_expense"`   // Total Kas Keluar (Expenses)
	ActualCash    float64   `json:"actual_cash"`
	ExpectedCash  float64   `json:"expected_cash"`
	Discrepancy  float64   `json:"discrepancy"`  // Selisih (Surplus/Minus) -> Masuk Akun Beban/Pendapatan Selisih Kas
	StatusShift  string    `json:"status_shift"` // "OPEN" atau "CLOSED"
	Timestamp    time.Time `json:"timestamp"`
}

// StoreOrderPayload menangkap invoice penjualan ritel lunas di kasir
type StoreOrderPayload struct {
	InvoiceNumber string    `json:"invoice_number"`
	CashierID     string    `json:"cashier_id"`
	TotalAmount   float64   `json:"total_amount"`   // Masuk akun: Debet Kas / QRIS, Kredit Pendapatan Ritel
	TaxAmount     float64   `json:"tax_amount"`     // Masuk akun: Kredit Utang Pajak (PPN)
	PaymentMethod string    `json:"payment_method"` // "TUNAI", "QRIS", "DEBIT"
	CreatedAt     time.Time `json:"created_at"`
}

// InventoryStockPayload menangkap mutasi penyesuaian/restock barang gudang dari Supplier
type InventoryStockPayload struct {
	ProductID     string    `json:"product_id"`
	MutationType  string    `json:"mutation_type"` // "INCOMING_SUPPLIER", "STOCK_ADJUSTMENT"
	Quantity      float64   `json:"qty"`
	UnitCost      float64   `json:"unit_cost"`      // Nilai harga beli modal HPP barang
	TotalCost     float64   `json:"total_cost"`     // Kebutuhan Akun: Debet Persediaan Barang Dagang, Kredit Kas/Utang Dagang
	ReferenceCode string    `json:"reference_code"` // Nomor PO Supplier / ID Gudang
	Timestamp     time.Time `json:"timestamp"`
}

// =========================================================================
// 📦 DEFENISI PAYLOAD OPERASIONAL: PORTAL-SEWA (Rental)
// =========================================================================

// RentalBookingPayload menangkap transaksi awal pemesanan hantaran / boks seserahan
type RentalBookingPayload struct {
	ReservationID string    `json:"reservation_id"`
	CustomerID    string    `json:"customer_id"`
	BoxCode       string    `json:"box_code"`
	RentalPrice   float64   `json:"rental_price"`   // Masuk akun: Kredit Pendapatan Sewa
	DepositAmount float64   `json:"deposit_amount"` // Masuk akun: Debet Kas, Kredit Utang Jaminan Pelanggan (Fisik)
	PenaltyAmount float64   `json:"penalty_amount"` // Masuk akun: Kredit Pendapatan Denda Kerusakan (jika ada)
	StatusBooking string    `json:"status_booking"` // "CONFIRMED", "RETURNED", "DAMAGED"
	DeadlineDate  time.Time `json:"deadline_date"`
}

// RentalReturnPayload menangkap penyelesaian akhir & pengembalian mika box hantaran
type RentalReturnPayload struct {
	ReturnID         string    `json:"return_id"`
	ReservationID    string    `json:"reservation_id"`
	LateDays         int       `json:"late_days"`
	TotalLateFees    float64   `json:"total_late_fees"`   // Pendapatan denda keterlambatan sewa
	TotalDamageFees  float64   `json:"total_damage_fees"` // Pendapatan denda kerusakan akrilik
	RemainingPayment float64   `json:"remaining_payment"` // Pelunasan piutang sewa (Debet Kas, Kredit Piutang/Pendapatan)
	GrandTotalPaid   float64   `json:"grand_total_paid"`  // Total uang fisik tunai yang diserahkan pelanggan saat ini
	Timestamp        time.Time `json:"timestamp"`
}

// =========================================================================
// 🧠 UNMARSHAL HELPER METHODS (Fungsi Pengurai Payload Otomatis & Adaptif)
// =========================================================================

// resolveRawPayload menangani deteksi double-string encoding akibat rest-api dispatching secara aman
func resolveRawPayload(raw []byte) []byte {
	if len(raw) == 0 {
		return raw
	}
	// Jika payload diawali tanda kutip ganda, artinya ini adalah raw string JSON terenkapsulasi
	if raw[0] == '"' {
		var unescapedString string
		if err := json.Unmarshal(raw, &unescapedString); err == nil {
			return []byte(unescapedString)
		}
	}
	return raw
}

// ToStoreShift mengekstrak payload jika event berasal dari pengelolaan laci kasir toko
func (e *FinanceInboxEvent) ToStoreShift() (*StoreShiftPayload, error) {
	if e.AggregateType != "SHIFT_SESSION" {
		return nil, fmt.Errorf("finance_entity: kegagalan tipe data, %s bukan bertipe SHIFT_SESSION", e.AggregateType)
	}
	var p StoreShiftPayload
	cleanedPayload := resolveRawPayload(e.Payload)
	if err := json.Unmarshal(cleanedPayload, &p); err != nil {
		return nil, fmt.Errorf("finance_entity: gagal mengurai SHIFT_SESSION: %w", err)
	}
	return &p, nil
}

// ToStoreOrder mengekstrak payload jika event berasal dari transaksi invoice retail kasir POS
func (e *FinanceInboxEvent) ToStoreOrder() (*StoreOrderPayload, error) {
	if e.AggregateType != "RETAIL_ORDER" {
		return nil, fmt.Errorf("finance_entity: kegagalan tipe data, %s bukan bertipe RETAIL_ORDER", e.AggregateType)
	}
	var p StoreOrderPayload
	cleanedPayload := resolveRawPayload(e.Payload)
	if err := json.Unmarshal(cleanedPayload, &p); err != nil {
		return nil, fmt.Errorf("finance_entity: gagal mengurai RETAIL_ORDER: %w", err)
	}
	return &p, nil
}

// ToRentalBooking mengekstrak payload jika event berasal dari transaksi penyewaan boks hantaran
func (e *FinanceInboxEvent) ToRentalBooking() (*RentalBookingPayload, error) {
	if e.AggregateType != "RENTAL_RESERVATION" {
		return nil, fmt.Errorf("finance_entity: kegagalan tipe data, %s bukan bertipe RENTAL_RESERVATION", e.AggregateType)
	}
	var p RentalBookingPayload
	cleanedPayload := resolveRawPayload(e.Payload)
	if err := json.Unmarshal(cleanedPayload, &p); err != nil {
		return nil, fmt.Errorf("finance_entity: gagal mengurai RENTAL_RESERVATION: %w", err)
	}
	return &p, nil
}

// ToRentalReturn mengekstrak payload pengembalian fisik boks, penyelesaian biaya denda, dan pelunasan sewa
func (e *FinanceInboxEvent) ToRentalReturn() (*RentalReturnPayload, error) {
	// Mendukung rute agregasi reservasi rental
	if e.AggregateType != "RENTAL_RESERVATION" && e.AggregateType != "RESERVATION" {
		return nil, fmt.Errorf("finance_entity: kegagalan tipe data, %s tidak didukung oleh handler return", e.AggregateType)
	}
	var p RentalReturnPayload
	cleanedPayload := resolveRawPayload(e.Payload)
	if err := json.Unmarshal(cleanedPayload, &p); err != nil {
		return nil, fmt.Errorf("finance_entity: gagal mengurai RENTAL_RETURN: %w", err)
	}
	return &p, nil
}

// ToInventoryStock mengekstrak payload jika event berasal dari penyesuaian stok masuk barang dagang gudang
func (e *FinanceInboxEvent) ToInventoryStock() (*InventoryStockPayload, error) {
	if e.AggregateType != "INVENTORY_STOCK" && e.AggregateType != "STOCK_MOVEMENT" {
		return nil, fmt.Errorf("finance_entity: kegagalan tipe data, %s bukan bertipe INVENTORY_STOCK", e.AggregateType)
	}
	var p InventoryStockPayload
	cleanedPayload := resolveRawPayload(e.Payload)
	if err := json.Unmarshal(cleanedPayload, &p); err != nil {
		return nil, fmt.Errorf("finance_entity: gagal mengurai INVENTORY_STOCK: %w", err)
	}
	return &p, nil
}
