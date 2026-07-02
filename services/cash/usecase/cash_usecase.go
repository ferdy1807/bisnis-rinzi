package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"bisnis-rinzi/packages/backend/database/minio"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/packages/backend/utils"
	"bisnis-rinzi/services/cash/dto"
	"bisnis-rinzi/services/cash/entity"
	"bisnis-rinzi/services/cash/repository"

	minioSDK "github.com/minio/minio-go/v7"
)

type cashUseCase struct {
	cashRepo    repository.CashRepository
	minioClient *minio.MinioClient
	mediaBucket string
}

func NewCashUseCase(repo repository.CashRepository, mc *minio.MinioClient, bucket string) CashUseCase {
	return &cashUseCase{cashRepo: repo, minioClient: mc, mediaBucket: bucket}
}

func (u *cashUseCase) OpenSession(ctx context.Context, cashierID string, input dto.OpenSessionRequest) (string, error) {
	existing, err := u.cashRepo.FindActiveSessionByCashierID(ctx, cashierID)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", errors.New("gagal membuka shift baru: Anda memiliki sesi kasir yang masih aktif berjalan")
	}

	now := time.Now()
	session := &entity.CashierSession{
		ID:           utils.GenerateUUIDv4(),
		CashierID:    cashierID,
		OpeningCash:  input.OpeningCash,
		ExpectedCash: &input.OpeningCash,
		Status:       "OPEN",
		OpenTime:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := session.ValidateOpen(); err != nil {
		return "", err
	}

	if err := u.cashRepo.SaveSession(ctx, session); err != nil {
		return "", err
	}
	return session.ID, nil
}

func (u *cashUseCase) CloseSession(ctx context.Context, cashierID string, input dto.CloseSessionRequest) error {
	session, err := u.cashRepo.FindActiveSessionByCashierID(ctx, cashierID)
	if err != nil || session == nil {
		return errors.New("tidak ada sesi kasir aktif yang dapat ditutup")
	}

	session.CalculateClosure(input.ActualCash)

	// Dapatkan summary untuk melengkapi payload
	summary, _ := u.GetShiftSummary(ctx, session.ID)
	var totalIncome, manualDeposit, totalExpense float64
	if summary != nil {
		totalIncome = summary.TotalIncome
		totalExpense = summary.TotalExpense
		manualDeposit = summary.TotalDeposit - summary.TotalIncome
	}

	reportPayload := map[string]interface{}{
		"session_id":       session.ID,
		"cashier_id":       session.CashierID,
		"opening_cash":     session.OpeningCash,
		"total_income":     totalIncome,
		"manual_deposit":   manualDeposit,
		"total_expense":    totalExpense,
		"expected_cash":    session.ExpectedCash,
		"actual_cash":      session.ActualCash,
		"difference_cash":  session.Difference,
		"closed_timestamp": session.CloseTime,
	}

	event, err := outbox.CreateEvent("CASH_SESSION", session.ID, "CASHIER_SHIFT_CLOSED", reportPayload)
	if err != nil {
		return fmt.Errorf("gagal merakit event laporan keuangan: %w", err)
	}

	return u.cashRepo.CloseSessionTx(ctx, session, event)
}

func (u *cashUseCase) CreateCashTransaction(ctx context.Context, cashierID string, input dto.CashTransactionRequest) error {
	session, err := u.cashRepo.FindActiveSessionByCashierID(ctx, cashierID)
	if err != nil || session == nil {
		return errors.New("transaksi ditolak: wajib membuka sesi kasir terlebih dahulu")
	}

	now := time.Now()
	trx := &entity.CashTransaction{
		ID:              utils.GenerateUUIDv4(),
		SessionID:       session.ID,
		TransactionType: input.TransactionType,
		ReferenceType:   "MANUAL",
		Amount:          input.Amount,
		CreatedBy:       cashierID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if input.Notes != "" {
		trx.Notes = &input.Notes
	}

	currentBalance := *session.ExpectedCash
	switch input.TransactionType {
	case "DEPOSIT":
		currentBalance += input.Amount
	case "WITHDRAWAL":
		if currentBalance < input.Amount {
			return errors.New("saldo kas di laci tidak mencukupi untuk melakukan penarikan uang")
		}
		currentBalance -= input.Amount
	default:
		return errors.New("tipe transaksi kas tidak valid")
	}

	session.ExpectedCash = &currentBalance
	session.UpdatedAt = now

	// Gunakan operasi atomik: simpan transaksi + update saldo session dalam 1 DB transaction
	return u.cashRepo.SaveCashTransactionAndUpdateSessionTx(ctx, trx, session)
}

func (u *cashUseCase) RecordInternalIncome(ctx context.Context, input dto.InternalIncomeRequest) error {
	session, err := u.cashRepo.FindLatestOpenSession(ctx)
	if err != nil {
		return fmt.Errorf("gagal mencari sesi kasir: %v", err)
	}
	if session == nil {
		return errors.New("tidak ada sesi kasir yang sedang buka saat ini, tidak dapat mencatat kas masuk")
	}

	now := time.Now()
	trx := &entity.CashTransaction{
		ID:              utils.GenerateUUIDv4(),
		SessionID:       session.ID,
		TransactionType: "DEPOSIT",
		ReferenceType:   input.Source,
		Amount:          input.Amount,
		Notes:           &input.Description,
		CreatedBy:       session.CashierID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	currentBalance := *session.ExpectedCash
	currentBalance += input.Amount
	session.ExpectedCash = &currentBalance
	session.UpdatedAt = now

	// Gunakan operasi atomik: simpan kas masuk + update saldo session dalam 1 DB transaction
	return u.cashRepo.SaveCashTransactionAndUpdateSessionTx(ctx, trx, session)
}

func (u *cashUseCase) CreateCategory(ctx context.Context, input dto.CreateExpenseCategoryRequest) error {
	now := time.Now()
	cat := &entity.ExpenseCategory{
		ID:        utils.GenerateUUIDv4(),
		Code:      input.Code,
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return u.cashRepo.SaveExpenseCategory(ctx, cat)
}

func (u *cashUseCase) AddExpense(ctx context.Context, cashierID string, input dto.CreateExpenseRequest) error {
	session, err := u.cashRepo.FindActiveSessionByCashierID(ctx, cashierID)
	if err != nil || session == nil {
		return errors.New("sesi kasir tidak aktif")
	}

	if *session.ExpectedCash < input.Amount {
		return errors.New("dana laci tidak cukup untuk membiayai pengeluaran operasional ini")
	}

	now := time.Now()
	exp := &entity.Expense{
		ID:          utils.GenerateUUIDv4(),
		ExpenseDate: now,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Description: input.Description,
		CreatedBy:   cashierID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Buat transaksi mutasi
	notes := "Pengeluaran kas: " + input.Description
	trx := &entity.CashTransaction{
		ID:              utils.GenerateUUIDv4(),
		SessionID:       session.ID,
		TransactionType: "WITHDRAWAL",
		ReferenceType:   "EXPENSE",
		ReferenceID:     &exp.ID,
		Amount:          input.Amount,
		Notes:           &notes,
		CreatedBy:       cashierID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	newBalance := *session.ExpectedCash - input.Amount
	session.ExpectedCash = &newBalance
	session.UpdatedAt = now

	// Gunakan operasi atomik: simpan expense + mutasi kas + update saldo session dalam 1 DB transaction
	return u.cashRepo.SaveExpenseAndUpdateSessionTx(ctx, exp, trx, session)
}

func (u *cashUseCase) GetCurrentSession(ctx context.Context, cashierID string) (*entity.CashierSession, error) {
	return u.cashRepo.FindActiveSessionByCashierID(ctx, cashierID)
}

func (u *cashUseCase) GetCategories(ctx context.Context) ([]*entity.ExpenseCategory, error) {
	return u.cashRepo.FindAllExpenseCategories(ctx)
}

func (u *cashUseCase) GetAllTransactions(ctx context.Context) ([]*entity.CashTransaction, error) {
	return u.cashRepo.FindAllTransactions(ctx)
}

func (u *cashUseCase) UpdateCategory(ctx context.Context, id string, input dto.CreateExpenseCategoryRequest) error {
	cat, err := u.cashRepo.FindExpenseCategoryByID(ctx, id)
	if err != nil || cat == nil {
		return errors.New("kategori pengeluaran tidak ditemukan")
	}
	cat.Code = input.Code
	cat.Name = input.Name
	cat.UpdatedAt = time.Now()
	return u.cashRepo.UpdateExpenseCategory(ctx, cat)
}

func (u *cashUseCase) DeleteCategory(ctx context.Context, id string) error {
	cat, err := u.cashRepo.FindExpenseCategoryByID(ctx, id)
	if err != nil || cat == nil {
		return errors.New("kategori pengeluaran tidak ditemukan")
	}
	return u.cashRepo.DeleteExpenseCategory(ctx, id)
}

func (u *cashUseCase) GetCategoryByID(ctx context.Context, id string) (*entity.ExpenseCategory, error) {
	return u.cashRepo.FindExpenseCategoryByID(ctx, id)
}

func (u *cashUseCase) GetAllExpenses(ctx context.Context) ([]*entity.Expense, error) {
	return u.cashRepo.FindAllExpenses(ctx)
}

func (u *cashUseCase) GetExpenseByID(ctx context.Context, id string) (*entity.Expense, error) {
	return u.cashRepo.FindExpenseByID(ctx, id)
}

func (u *cashUseCase) UploadShiftReport(ctx context.Context, id string, file io.Reader, fileSize int64, filename string) (*entity.CashierSession, error) {
	session, err := u.cashRepo.FindSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("sesi kasir tidak ditemukan")
	}

	err = u.minioClient.CreateBucketIfNotExist(ctx, u.mediaBucket, "us-east-1")
	if err != nil {
		return nil, fmt.Errorf("gagal verifikasi minio cluster: %w", err)
	}

	_ = u.minioClient.MakeBucketPublic(ctx, u.mediaBucket)

	now := time.Now()
	objectName := fmt.Sprintf("%s/%s-%s-%s", now.Format("2006"), now.Format("01"), now.Format("02"), filename)

	_, err = u.minioClient.Client.PutObject(ctx, u.mediaBucket, objectName, file, fileSize, minioSDK.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return nil, fmt.Errorf("gagal upload pdf ke minio: %w", err)
	}

	url := fmt.Sprintf("http://localhost:9000/%s/%s", u.mediaBucket, objectName)
	session.ReceiptURL = &url

	err = u.cashRepo.UpdateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *cashUseCase) GetAllShifts(ctx context.Context) ([]*entity.CashierSession, error) {
	return u.cashRepo.FindAllSessions(ctx)
}

func (u *cashUseCase) GetInternalIncomes(ctx context.Context) ([]*entity.CashTransaction, error) {
	return u.cashRepo.FindInternalIncomes(ctx)
}

func (u *cashUseCase) GetShiftByID(ctx context.Context, id string) (*entity.CashierSession, error) {
	return u.cashRepo.FindSessionByID(ctx, id)
}

func (u *cashUseCase) GetShiftSummary(ctx context.Context, id string) (*dto.ShiftSummaryResponse, error) {
	session, err := u.cashRepo.FindSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("sesi shift tidak ditemukan")
	}

	transactions, _ := u.cashRepo.FindTransactionsBySession(ctx, id)

	var totalDeposit, totalWithdrawal float64
	var totalIncome, totalExpense float64
	var validCashIn float64

	for _, t := range transactions {
		switch t.TransactionType {
		case "DEPOSIT":
			totalDeposit += t.Amount
			if t.ReferenceType == "POS_SALE" || t.ReferenceType == "MANUAL" {
				validCashIn += t.Amount
			}
			if t.ReferenceType == "POS_SALE" {
				totalIncome += t.Amount
			}
		case "WITHDRAWAL":
			totalWithdrawal += t.Amount
			if t.ReferenceType == "EXPENSE" {
				totalExpense += t.Amount
			}
		}
	}

	expectedBalance := session.OpeningCash + validCashIn - totalWithdrawal

	resp := &dto.ShiftSummaryResponse{
		SessionID:       session.ID,
		CashierID:       session.CashierID,
		OpeningCash:     session.OpeningCash,
		TotalIncome:     totalIncome,
		TotalExpense:    totalExpense,
		TotalDeposit:    totalDeposit,
		TotalWithdrawal: totalWithdrawal,
		ExpectedCash:    expectedBalance,
		Status:          session.Status,
	}

	if session.ActualCash != nil {
		resp.ActualCash = *session.ActualCash
	}
	if session.Difference != nil {
		resp.Difference = *session.Difference
	}

	return resp, nil
}

func (u *cashUseCase) GetTransactionByID(ctx context.Context, id string) (*entity.CashTransaction, error) {
	return u.cashRepo.FindTransactionByID(ctx, id)
}

func (u *cashUseCase) UpdateExpense(ctx context.Context, id string, input dto.CreateExpenseRequest) error {
	exp, err := u.cashRepo.FindExpenseByID(ctx, id)
	if err != nil || exp == nil {
		return errors.New("data pengeluaran tidak ditemukan")
	}

	exp.CategoryID = input.CategoryID
	exp.Amount = input.Amount
	exp.Description = input.Description
	exp.UpdatedAt = time.Now()

	return u.cashRepo.UpdateExpense(ctx, exp)
}

func (u *cashUseCase) DeleteExpense(ctx context.Context, id string) error {
	exp, err := u.cashRepo.FindExpenseByID(ctx, id)
	if err != nil || exp == nil {
		return errors.New("data pengeluaran tidak ditemukan")
	}

	return u.cashRepo.DeleteExpense(ctx, id)
}
