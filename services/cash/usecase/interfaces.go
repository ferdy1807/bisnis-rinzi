package usecase

import (
	"bisnis-rinzi/services/cash/dto"
	"bisnis-rinzi/services/cash/entity"
	"context"
	"io"
)

type CashUseCase interface {
	// Sessions
	OpenSession(ctx context.Context, cashierID string, input dto.OpenSessionRequest) (string, error)
	CloseSession(ctx context.Context, cashierID string, input dto.CloseSessionRequest) error
	GetCurrentSession(ctx context.Context, cashierID string) (*entity.CashierSession, error)
	GetAllShifts(ctx context.Context) ([]*entity.CashierSession, error)
	GetShiftByID(ctx context.Context, id string) (*entity.CashierSession, error)
	GetShiftSummary(ctx context.Context, id string) (*dto.ShiftSummaryResponse, error)
	UploadShiftReport(ctx context.Context, id string, file io.Reader, fileSize int64, filename string) (*entity.CashierSession, error)

	// Transactions
	CreateCashTransaction(ctx context.Context, cashierID string, input dto.CashTransactionRequest) error
	RecordInternalIncome(ctx context.Context, input dto.InternalIncomeRequest) error
	GetAllTransactions(ctx context.Context) ([]*entity.CashTransaction, error)
	GetInternalIncomes(ctx context.Context) ([]*entity.CashTransaction, error)
	GetTransactionByID(ctx context.Context, id string) (*entity.CashTransaction, error)

	// Expense Categories
	CreateCategory(ctx context.Context, input dto.CreateExpenseCategoryRequest) error
	UpdateCategory(ctx context.Context, id string, input dto.CreateExpenseCategoryRequest) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategories(ctx context.Context) ([]*entity.ExpenseCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*entity.ExpenseCategory, error)

	// Expenses
	AddExpense(ctx context.Context, cashierID string, input dto.CreateExpenseRequest) error
	UpdateExpense(ctx context.Context, id string, input dto.CreateExpenseRequest) error
	DeleteExpense(ctx context.Context, id string) error
	GetAllExpenses(ctx context.Context) ([]*entity.Expense, error)
	GetExpenseByID(ctx context.Context, id string) (*entity.Expense, error)
}
