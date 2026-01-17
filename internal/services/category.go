package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/repositories"
)

type CategoryService interface {
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Category, error)
	Create(ctx context.Context, category *model.Category) (*model.Category, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Category, error) {
	return s.categoryRepo.ListByUser(ctx, userID, limit, offset)
}

func (s *categoryService) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	// Add business logic here if needed (e.g., check for duplicate category names)
	return s.categoryRepo.Create(ctx, category)
}
