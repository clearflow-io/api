package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/table"
)

type ExpenseRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Expense, error)
	Create(ctx context.Context, expense *model.Expense) (*model.Expense, error)
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Expense, error) {
	query := table.Expense.SELECT(
		table.Expense.AllColumns,
	).FROM(
		table.Expense,
	).WHERE(
		table.Expense.UserID.EQ(postgres.UUID(userID)),
	).ORDER_BY(
		table.Expense.CreatedAt.DESC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	var dest []model.Expense
	err := query.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *expenseRepository) Create(ctx context.Context, expense *model.Expense) (*model.Expense, error) {
	query := table.Expense.INSERT(
		table.Expense.UserID,
		table.Expense.Amount,
		table.Expense.Description,
		table.Expense.PurchaseDate,
		table.Expense.BillDate,
		table.Expense.CategoryID,
	).VALUES(
		expense.UserID,
		expense.Amount,
		expense.Description,
		expense.PurchaseDate,
		expense.BillDate,
		expense.CategoryID,
	).RETURNING(table.Expense.AllColumns)

	err := query.QueryContext(ctx, r.db, expense)
	if err != nil {
		return nil, err
	}

	return expense, nil
}
