package repository

import (
	"bisnis-rinzi/packages/backend/outbox"
	"bisnis-rinzi/services/cash/entity"
	"context"
)

type CashRepository interface {
	// Cashier Sessions
	SaveSession(ctx context.Context, session *entity.CashierSession) error
	UpdateSession(ctx context.Context, session *entity.CashierSession) error
	FindActiveSessionByCashierID(ctx context.Context, cashierID string) (*entity.CashierSession, error)
	FindLatestOpenSession(ctx context.Context) (*entity.CashierSession, error)
	FindSessionByID(ctx context.Context, id string) (*entity.CashierSession, error)
	FindAllSessions(ctx context.Context) ([]*entity.CashierSession, error)
	CloseSessionTx(ctx context.Context, session *entity.CashierSession, event *outbox.Event) error

	// Cash Transactions (Mutasi Laci)
	SaveCashTransaction(ctx context.Context, tx *entity.CashTransaction) error
	// Atomic: Menyimpan transaksi kas dan memperbarui saldo session dalam 1 DB transaction
	SaveCashTransactionAndUpdateSessionTx(ctx context.Context, trx *entity.CashTransaction, session *entity.CashierSession) error
	FindAllTransactions(ctx context.Context) ([]*entity.CashTransaction, error)
	FindInternalIncomes(ctx context.Context) ([]*entity.CashTransaction, error)
	FindTransactionByID(ctx context.Context, id string) (*entity.CashTransaction, error)
	FindTransactionsBySession(ctx context.Context, sessionID string) ([]*entity.CashTransaction, error)

	// Expense Categories
	SaveExpenseCategory(ctx context.Context, cat *entity.ExpenseCategory) error
	UpdateExpenseCategory(ctx context.Context, cat *entity.ExpenseCategory) error
	DeleteExpenseCategory(ctx context.Context, id string) error
	FindAllExpenseCategories(ctx context.Context) ([]*entity.ExpenseCategory, error)
	FindExpenseCategoryByID(ctx context.Context, id string) (*entity.ExpenseCategory, error)

	// Expenses
	SaveExpenseTx(ctx context.Context, exp *entity.Expense, txRecord *entity.CashTransaction) error
	// Atomic: Menyimpan pengeluaran, transaksi kas, dan memperbarui saldo session dalam 1 DB transaction
	SaveExpenseAndUpdateSessionTx(ctx context.Context, exp *entity.Expense, txRecord *entity.CashTransaction, session *entity.CashierSession) error
	UpdateExpense(ctx context.Context, exp *entity.Expense) error
	DeleteExpense(ctx context.Context, id string) error
	FindAllExpenses(ctx context.Context) ([]*entity.Expense, error)
	FindExpenseByID(ctx context.Context, id string) (*entity.Expense, error)
}
