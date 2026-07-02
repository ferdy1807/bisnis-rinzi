package entity

import "time"

type JournalEntryDetail struct {
	ID             string    `json:"id"`
	JournalEntryID string    `json:"journal_entry_id"`
	AccountID      string    `json:"account_id"`
	DebitAmount    float64   `json:"debit_amount"`
	CreditAmount   float64   `json:"credit_amount"`
	CreatedAt      time.Time `json:"created_at"`
}
