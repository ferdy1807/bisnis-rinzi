package entity

import (
	"errors"
	"time"
)

type CashierSession struct {
	ID           string     `json:"id"`
	CashierID    string     `json:"cashier_id"`
	OpeningCash  float64    `json:"opening_cash"`
	ExpectedCash *float64   `json:"expected_cash,omitempty"`
	ActualCash   *float64   `json:"actual_cash,omitempty"`
	Difference   *float64   `json:"difference,omitempty"`
	Status       string     `json:"status"` // "OPEN", "CLOSED"
	ReceiptURL   *string    `json:"receipt_url,omitempty"`
	OpenTime     time.Time  `json:"open_time"`
	CloseTime    *time.Time `json:"close_time,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	TotalManualIncome float64    `json:"total_manual_income,omitempty"`
}

func (cs *CashierSession) ValidateOpen() error {
	if cs.CashierID == "" {
		return errors.New("cashier_id wajib ditentukan untuk membuka shift")
	}
	if cs.OpeningCash < 0 {
		return errors.New("modal awal pembukaan laci kas tidak boleh negatif")
	}
	if cs.Status != "OPEN" {
		return errors.New("status sesi baru harus bernilai OPEN")
	}
	return nil
}

func (cs *CashierSession) CalculateClosure(actual float64) {
	cs.ActualCash = &actual
	expected := 0.0
	if cs.ExpectedCash != nil {
		expected = *cs.ExpectedCash
	}
	diff := actual - expected
	cs.Difference = &diff
	cs.Status = "CLOSED"
	now := time.Now()
	cs.CloseTime = &now
	cs.UpdatedAt = now
}
