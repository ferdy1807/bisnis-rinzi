package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	dbMinio "bisnis-rinzi/packages/backend/database/minio"
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/packages/backend/utils"
	"bisnis-rinzi/services/finance/dto"
	"bisnis-rinzi/services/finance/entity"
	"bisnis-rinzi/services/finance/repository"
)

type financeUseCase struct {
	financeRepo repository.FinanceRepository
	minioClient *dbMinio.MinioClient
	bucketName  string
}

func NewFinanceUseCase(repo repository.FinanceRepository, minioClient *dbMinio.MinioClient, bucketName string) FinanceUseCase {
	if minioClient != nil {
		minioClient.CreateBucketIfNotExist(context.Background(), bucketName, "us-east-1")
		minioClient.MakeBucketPublic(context.Background(), bucketName)
	}
	return &financeUseCase{
		financeRepo: repo,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (u *financeUseCase) UploadDailyIncomeReport(ctx context.Context, filename string, fileBytes []byte) (string, error) {
	if u.minioClient == nil {
		return "", errors.New("minio client is not configured")
	}
	url, err := u.minioClient.UploadFile(ctx, u.bucketName, filename, fileBytes, "application/pdf")
	if err != nil {
		return "", fmt.Errorf("failed to upload report to minio: %w", err)
	}
	return url, nil
}

func (u *financeUseCase) CreateAccount(ctx context.Context, input dto.CreateCOARequest) error {
	existing, err := u.financeRepo.FindCOAByCode(ctx, input.AccountCode)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("nomor kode akun perkiraan sudah terdaftar")
	}

	coa := &entity.ChartOfAccount{
		ID:             utils.GenerateUUIDv4(),
		AccountCode:    input.AccountCode,
		AccountName:    input.AccountName,
		AccountGroup:   input.AccountGroup,
		NormalBalance:  input.NormalBalance,
		CurrentBalance: 0.0,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := coa.Validate(); err != nil {
		return err
	}
	return u.financeRepo.SaveCOA(ctx, coa)
}

func (u *financeUseCase) OpenAccountingPeriod(ctx context.Context, input dto.OpenPeriodRequest) error {
	active, err := u.financeRepo.FindActivePeriod(ctx)
	if err != nil {
		return err
	}
	if active != nil {
		return errors.New("gagal membuka periode baru: terdapat periode akuntansi berjalan yang belum ditutup")
	}

	p := &entity.AccountingPeriod{
		ID:        utils.GenerateUUIDv4(),
		Name:      input.Name,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		IsClosed:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := p.Validate(); err != nil {
		return err
	}
	return u.financeRepo.SavePeriod(ctx, p)
}

func (u *financeUseCase) RecordManualJournal(ctx context.Context, input dto.CreateJournalEntryRequest) error {
	period, err := u.financeRepo.FindActivePeriod(ctx)
	if err != nil || period == nil {
		return errors.New("entri jurnal ditolak: tidak ada periode akuntansi aktif")
	}

	// Validasi Balance (Keseimbangan Debet Kredit)
	var totalDebit, totalCredit float64
	var details []*entity.JournalEntryDetail
	entryID := utils.GenerateUUIDv4()
	now := time.Now()

	for _, d := range input.Details {
		totalDebit += d.DebitAmount
		totalCredit += d.CreditAmount

		details = append(details, &entity.JournalEntryDetail{
			ID:             utils.GenerateUUIDv4(),
			JournalEntryID: entryID,
			AccountID:      d.AccountID,
			DebitAmount:    d.DebitAmount,
			CreditAmount:   d.CreditAmount,
			CreatedAt:      now,
		})
	}

	if totalDebit != totalCredit {
		return fmt.Errorf("jurnal entri tidak seimbang (unbalanced): Total Debet Rp%.2f tidak sama dengan Total Kredit Rp%.2f", totalDebit, totalCredit)
	}

	// Cek atau buat journal category
	journalCode := input.JournalID
	journal, err := u.financeRepo.FindJournalByCode(ctx, journalCode)
	if err != nil {
		return err
	}
	if journal == nil {
		journal = &entity.Journal{
			ID:          utils.GenerateUUIDv4(),
			JournalCode: journalCode,
			Name:        "Jurnal " + journalCode,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := u.financeRepo.SaveJournal(ctx, journal); err != nil {
			return err
		}
	}

	entry := &entity.JournalEntry{
		ID:                 entryID,
		JournalID:          journal.ID,
		AccountingPeriodID: period.ID,
		ReferenceNumber:    input.ReferenceNumber,
		EntryDate:          now,
		Narration:          input.Narration,
		IsPosted:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return u.financeRepo.SaveJournalEntryTx(ctx, entry, details)
}

func (u *financeUseCase) ProcessClosingAndReconciliation(ctx context.Context, cashierID string, input dto.ProcessDailyClosingRequest) error {
	// 1. Hitung total nominal buku arus kas otomatis berdasarkan sinkronisasi outbox yang masuk
	salesSys, _ := u.financeRepo.SumTransactionAmountByType(ctx, "SALES", input.ClosingDate)
	rentalSys, _ := u.financeRepo.SumTransactionAmountByType(ctx, "RENTAL", input.ClosingDate)
	expenseSys, _ := u.financeRepo.SumTransactionAmountByType(ctx, "EXPENSE", input.ClosingDate)

	// Cek apakah sudah ada closing hari ini
	existingClosing, err := u.financeRepo.FindDailyClosingByDate(ctx, input.ClosingDate)
	if err != nil {
		return err
	}
	if existingClosing != nil {
		return errors.New("tutup buku untuk hari ini sudah dilakukan sebelumnya")
	}

	closingID := utils.GenerateUUIDv4()
	now := time.Now()
	var logs []*entity.ReconciliationLog

	// 2. Petakan data kliring rekonsiliasi kasir fisik vs sistem
	for _, r := range input.Reconciliations {
		var sysAmt float64
		if r.TargetSystem == "POS" {
			sysAmt = salesSys
		}
		if r.TargetSystem == "RENTAL" {
			sysAmt = rentalSys
		}
		if r.TargetSystem == "CASH" {
			sysAmt = expenseSys
		}

		discrepancy := r.ActualAmount - sysAmt
		logs = append(logs, &entity.ReconciliationLog{
			ID:             utils.GenerateUUIDv4(),
			DailyClosingID: closingID,
			TargetSystem:   r.TargetSystem,
			SystemAmount:   sysAmt,
			ActualAmount:   r.ActualAmount,
			Discrepancy:    discrepancy,
			Notes:          r.Notes,
			ReconciledBy:   cashierID,
			CreatedAt:      now,
		})
	}

	netFlow := (salesSys + rentalSys) - expenseSys
	dcEntity := &entity.DailyClosing{
		ID:                closingID,
		ClosingDate:       input.ClosingDate,
		TotalSalesRetail:  salesSys,
		TotalRentalIncome: rentalSys,
		TotalExpenses:     expenseSys,
		NetCashFlow:       netFlow,
		IsReconciled:      true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// 3. Rakit outbox untuk pelaporan pajak/sinkronisasi cabang jika dibutuhkan di masa depan
	outboxPayload := map[string]interface{}{"daily_closing_id": closingID, "net_cash_flow": netFlow}
	event, err := outbox.CreateEvent("FINANCE_CLOSING", closingID, "DAILY_CLOSING_COMPLETED", outboxPayload)
	if err != nil {
		return err
	}

	return u.financeRepo.SaveDailyClosingTx(ctx, dcEntity, logs, event)
}

func (u *financeUseCase) GenerateLedgerReport(ctx context.Context, accountCode string, start, end time.Time) (*dto.LedgerReportResponse, error) {
	coa, err := u.financeRepo.FindCOAByCode(ctx, accountCode)
	if err != nil || coa == nil {
		return nil, errors.New("akun perkiraan tidak ditemukan")
	}

	lines, err := u.financeRepo.GetLedgerLines(ctx, coa.ID, start, end)
	if err != nil {
		return nil, err
	}

	// Urutkan running balance berdasarkan posisi saldo normal COA
	var currentBal float64 = 0.0 // Anggap saldo awal periode berjalan 0 demi simplisitas matematika
	for idx := range lines {
		if coa.NormalBalance == "DEBIT" {
			currentBal += lines[idx].Debit - lines[idx].Credit
		} else {
			currentBal += lines[idx].Credit - lines[idx].Debit
		}
		lines[idx].RunningBalance = currentBal
	}

	return &dto.LedgerReportResponse{
		AccountCode:    coa.AccountCode,
		AccountName:    coa.AccountName,
		InitialBalance: 0.0,
		FinalBalance:   currentBal,
		Lines:          lines,
	}, nil
}

func (u *financeUseCase) GenerateTrialBalance(ctx context.Context) ([]dto.TrialBalanceLine, error) {
	coas, _ := u.financeRepo.FindAllCOA(ctx)
	var lines []dto.TrialBalanceLine
	for _, coa := range coas {
		var deb, cred float64
		if coa.NormalBalance == "DEBIT" {
			deb = coa.CurrentBalance
		} else {
			cred = coa.CurrentBalance
		}
		lines = append(lines, dto.TrialBalanceLine{AccountCode: coa.AccountCode, AccountName: coa.AccountName, Debit: deb, Credit: cred})
	}
	return lines, nil
}

// ------------------------------------------
// NEW IMPLEMENTATIONS
// ------------------------------------------

func (u *financeUseCase) GetAccountByID(ctx context.Context, id string) (*entity.ChartOfAccount, error) {
	return u.financeRepo.FindCOAByID(ctx, id)
}

func (u *financeUseCase) UpdateAccount(ctx context.Context, id string, req dto.CreateCOARequest) error {
	acc, err := u.financeRepo.FindCOAByID(ctx, id)
	if err != nil || acc == nil {
		return errors.New("akun tidak ditemukan")
	}
	acc.AccountName = req.AccountName
	acc.NormalBalance = req.NormalBalance
	acc.IsActive = req.IsActive
	acc.UpdatedAt = time.Now()
	return u.financeRepo.UpdateCOA(ctx, acc)
}

func (u *financeUseCase) DeleteAccount(ctx context.Context, id string) error {
	return u.financeRepo.DeleteCOA(ctx, id)
}

func (u *financeUseCase) GetAllJournals(ctx context.Context) ([]*entity.JournalEntry, error) {
	entries, err := u.financeRepo.FindAllJournalEntries(ctx)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (u *financeUseCase) GetJournalByID(ctx context.Context, id string) (*entity.JournalEntry, error) {
	entry, err := u.financeRepo.FindJournalEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	details, err := u.financeRepo.FindJournalEntryDetails(ctx, id)
	if err == nil {
		entry.Details = details
	}
	return entry, nil
}

func (u *financeUseCase) GenerateCashFlowStatement(ctx context.Context) (map[string]interface{}, error) {
	closings, _ := u.financeRepo.FindAllDailyClosings(ctx)

	var operating, investing, financing, net float64
	for _, c := range closings {
		// Operating: Uang masuk dari penjualan & sewa & pendapatan lain, dikurangi pengeluaran operasional
		operating += c.TotalSalesRetail + c.TotalRentalIncome + c.TotalOtherIncome - c.TotalExpenses
		net += c.NetCashFlow
	}

	// Untuk bisnis MVP ini, seluruh mutasi dianggap operasional, investing/financing diset 0.
	res := map[string]interface{}{
		"operating_activities": operating,
		"investing_activities": investing,
		"financing_activities": financing,
		"net_cash_flow":        net,
	}
	return res, nil
}

func (u *financeUseCase) GetAnalyticsData(ctx context.Context, action string) (map[string]interface{}, error) {
	// Simulasi agregator analytics dari Database
	return u.financeRepo.GetAnalyticsData(ctx, action)
}

func (u *financeUseCase) RecordAutoJournalSystem(ctx context.Context, sourceModule string, refID string, lines []entity.JournalEntryDetail) error {
	// Akan diimplementasikan lebih lanjut untuk integrasi worker
	return nil
}

func (u *financeUseCase) GenerateIncomeStatement(ctx context.Context) (*dto.IncomeStatementResponse, error) {
	rev, hpp, exp, _ := u.financeRepo.GetTotalMonthlyAnalytics(ctx)

	items := map[string]float64{
		"Pendapatan Ritel & Sewa": rev,
		"Harga Pokok Penjualan":   hpp,
		"Beban Operasional":       exp,
	}

	grossProfit := rev - hpp
	netIncome := grossProfit - exp

	return &dto.IncomeStatementResponse{
		PeriodName:   "Periode Berjalan Terkonsolidasi",
		TotalRevenue: rev,
		TotalCOGS:    hpp,
		GrossProfit:  grossProfit,
		TotalExpense: exp,
		NetIncome:    netIncome,
		Items:        items,
	}, nil
}

func (u *financeUseCase) GenerateBalanceSheet(ctx context.Context) (*dto.BalanceSheetResponse, error) {
	// Untuk MVP, Balance Sheet dihitung berdasarkan net income dan mutasi kas
	// Total Aset = Kas & Setara Kas (dari net cash flow closings) + Aset Lancar lainnya
	closings, _ := u.financeRepo.FindAllDailyClosings(ctx)
	var cash float64
	for _, c := range closings {
		cash += c.NetCashFlow
	}

	// Ekuitas disetarakan dengan aset untuk menyimbangkan (Simplified MVP)
	items := map[string]float64{
		"Kas Tunai & Setara Kas": cash,
		"Modal & Laba Ditahan":   cash,
	}

	return &dto.BalanceSheetResponse{
		PeriodName:  "Neraca Terkonsolidasi MVP",
		TotalAssets: cash,
		TotalLiab:   0.0,
		TotalEquity: cash,
		Items:       items,
	}, nil
}

func (u *financeUseCase) GetCorporateDashboardAnalytics(ctx context.Context, start, end time.Time) (map[string]interface{}, error) {
	metrics, _ := u.financeRepo.GetDashboardMetrics(ctx)

	year := time.Now().Year()
	monthlyTrend, _ := u.financeRepo.GetMonthlyRevenueTrend(ctx, year)
	rentals, _ := u.financeRepo.GetRentalAnalytics(ctx, start, end)
	topProd, _ := u.financeRepo.GetTopProductsAnalytics(ctx, 5)
	catCont, _ := u.financeRepo.GetCategoryContributionAnalytics(ctx)
	stocks, _ := u.financeRepo.GetStockValueAnalytics(ctx)
	cashiers, _ := u.financeRepo.GetCashierPerformanceAnalytics(ctx)

	metrics["sales_trend"] = monthlyTrend
	metrics["rental_trend"] = rentals
	metrics["top_products"] = topProd
	metrics["category_contribution"] = catCont
	metrics["inventory_valuation"] = stocks
	metrics["cashier_performance"] = cashiers

	return metrics, nil
}

func (u *financeUseCase) GetAccounts(ctx context.Context) ([]*entity.ChartOfAccount, error) {
	return u.financeRepo.FindAllCOA(ctx)
}

func (u *financeUseCase) GetAllPeriods(ctx context.Context) ([]*entity.AccountingPeriod, error) {
	return u.financeRepo.FindAllPeriods(ctx)
}

func (u *financeUseCase) DeleteAccountingPeriod(ctx context.Context, id string) error {
	return u.financeRepo.DeletePeriod(ctx, id)
}

func (u *financeUseCase) UpdateAccountingPeriod(ctx context.Context, id string, req dto.OpenPeriodRequest) error {
	p, err := u.financeRepo.FindPeriodByID(ctx, id)
	if err != nil || p == nil {
		return errors.New("periode akuntansi tidak ditemukan")
	}
	p.Name = req.Name
	p.StartDate = req.StartDate
	p.EndDate = req.EndDate
	p.UpdatedAt = time.Now()

	return u.financeRepo.UpdatePeriod(ctx, p)
}

func (u *financeUseCase) LockPeriod(ctx context.Context, id string, userID string, reason string) error {
	p, err := u.financeRepo.FindPeriodByID(ctx, id)
	if err != nil || p == nil {
		return errors.New("periode akuntansi tidak ditemukan")
	}
	if p.IsClosed {
		return errors.New("periode tersebut sudah berada dalam status terkunci")
	}

	lockEntity := &entity.PeriodLock{
		ID:                 utils.GenerateUUIDv4(),
		AccountingPeriodID: id,
		LockedBy:           userID,
		LockReason:         reason,
		CreatedAt:          time.Now(),
	}
	return u.financeRepo.ClosePeriodTx(ctx, id, lockEntity)
}

func (u *financeUseCase) GetAllDailyClosings(ctx context.Context) ([]*entity.DailyClosing, error) {
	return u.financeRepo.FindAllDailyClosings(ctx)
}

func (u *financeUseCase) GetDailyClosingDetail(ctx context.Context, id string) (*entity.DailyClosing, error) {
	return u.financeRepo.FindDailyClosingByID(ctx, id)
}

func (u *financeUseCase) GetReconciliationLogs(ctx context.Context, closingID string) ([]*entity.ReconciliationLog, error) {
	return u.financeRepo.FindReconciliationLogsByClosingID(ctx, closingID)
}

// File: services/finance/usecase/finance_usecase.go

func (u *financeUseCase) EnqueueInboxEvent(ctx context.Context, req dto.JournalIncomingRequest) error {
	// 1. Transformasikan DTO menjadi objek standard library outbox.Event
	evt := &outbox.Event{
		ID:            req.ID,
		AggregateType: req.AggregateType,
		AggregateID:   req.AggregateID,
		EventType:     req.EventType,
		Payload:       []byte(req.Payload),
		Status:        outbox.StatusPending,
		CreatedAt:     req.CreatedAt,
		UpdatedAt:     time.Now().UTC(),
	}

	// 2. Buka transaksi database secara aman via internal database connection milik UseCase/Repo Anda
	// Catatan: Sesuaikan u.repo.GetPool() atau u.dbClient tergantung bagaimana Anda menyimpan pool di usecase
	tx, err := u.financeRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("usecase: gagal menginisialisasi transaksi db: %w", err)
	}
	defer tx.Rollback(ctx)

	// 3. Simpan row ke tabel outbox_events menggunakan transaksi yang sama
	if err := outbox.SaveEventTx(ctx, tx, evt); err != nil {
		return fmt.Errorf("usecase: gagal menyimpan data outbox ke tabel lokal: %w", err)
	}

	// 4. Komit transaksi database
	return tx.Commit(ctx)
}
