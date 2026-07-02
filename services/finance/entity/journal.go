package entity

import (
	"errors"
	"time"
)

type Journal struct {
	ID          string    `json:"id"`
	JournalCode string    `json:"journal_code"` // e.g., "GJ" (General Journal), "SJ" (Sales Journal)
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (j *Journal) Validate() error {
	if j.JournalCode == "" || j.Name == "" {
		return errors.New("kode dan nama jurnal wajib ditentukan")
	}
	return nil
}
