package entity

import (
	"errors"
	"strings"
	"time"
)

type ChartOfAccount struct {
	ID             string     `json:"id"`
	AccountCode    string     `json:"account_code"` // e.g., "11110" (Kas Utama), "41100" (Pendapatan Retail)
	AccountName    string     `json:"account_name"`
	AccountGroup   string     `json:"account_group"`  // "ASSET", "LIABILITY", "EQUITY", "REVENUE", "EXPENSE"
	NormalBalance  string     `json:"normal_balance"` // "DEBIT" or "CREDIT"
	CurrentBalance float64    `json:"current_balance"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func (coa *ChartOfAccount) Validate() error {
	if strings.TrimSpace(coa.AccountCode) == "" {
		return errors.New("kode akun (account code) tidak boleh kosong")
	}
	if strings.TrimSpace(coa.AccountName) == "" {
		return errors.New("nama akun (account name) tidak boleh kosong")
	}
	group := strings.ToUpper(coa.AccountGroup)
	if group != "ASSET" && group != "LIABILITY" && group != "EQUITY" && group != "REVENUE" && group != "EXPENSE" {
		return errors.New("kelompok akun harus berupa ASSET, LIABILITY, EQUITY, REVENUE, atau EXPENSE")
	}
	balance := strings.ToUpper(coa.NormalBalance)
	if balance != "DEBIT" && balance != "CREDIT" {
		return errors.New("saldo normal akun harus bernilai DEBIT atau CREDIT")
	}
	return nil
}
