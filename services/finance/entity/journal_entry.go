package entity

import (
	"errors"
	"time"
)

type JournalEntry struct {
	ID                 string     `json:"id"`
	JournalID          string     `json:"journal_id"`
	AccountingPeriodID string     `json:"accounting_period_id"`
	ReferenceNumber    string     `json:"reference_number"` // e.g., Invoice Number atau Sesi Kasir ID
	EntryDate          time.Time  `json:"entry_date"`
	Narration          string     `json:"narration"` // Deskripsi transaksi keseluruhan
	IsPosted           bool       `json:"is_posted"` // Status posting ke Buku Besar (Ledger)
	PostedAt           *time.Time `json:"posted_at,omitempty"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	Details            []*JournalEntryDetail `json:"details,omitempty"`
}

func (je *JournalEntry) Validate() error {
	if je.JournalID == "" {
		return errors.New("entri jurnal harus dikaitkan dengan master jurnal yang valid")
	}
	if je.AccountingPeriodID == "" {
		return errors.New("entri jurnal harus terikat ke periode akuntansi yang aktif")
	}
	if je.ReferenceNumber == "" {
		return errors.New("nomor referensi dokumen asal wajib dicantumkan")
	}
	return nil
}
