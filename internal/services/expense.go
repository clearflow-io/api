package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/repositories"
	"github.com/igorschechtel/clearflow-backend/internal/utils"
)

type ExpenseService interface {
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Expense, error)
	Create(ctx context.Context, expense *model.Expense) (*model.Expense, error)
}

type expenseService struct {
	expenseRepo  repositories.ExpenseRepository
	categoryRepo repositories.CategoryRepository
}

func NewExpenseService(expenseRepo repositories.ExpenseRepository, categoryRepo repositories.CategoryRepository) ExpenseService {
	return &expenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *expenseService) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Expense, error) {
	return s.expenseRepo.ListByUser(ctx, userID, limit, offset)
}

func (s *expenseService) Create(ctx context.Context, expense *model.Expense) (*model.Expense, error) {
	// Business Logic: If a category is provided, verify it exists and belongs to the user
	if expense.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *expense.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, utils.ErrNotFound
		}
		if category.UserID != expense.UserID {
			return nil, utils.ErrForbidden
		}
	}

	return s.expenseRepo.Create(ctx, expense)
}
