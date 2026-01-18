package services

import (
	"context"
	"fmt"

	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/repositories"
)

type CategoryService interface {
	ListByUser(ctx context.Context, clerkID string, limit, offset int) ([]model.Category, error)
	Create(ctx context.Context, clerkID string, category *model.Category) (*model.Category, error)
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
	userService  UserService
}

func NewCategoryService(categoryRepo repositories.CategoryRepository, userService UserService) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		userService:  userService,
	}
}

func (s *categoryService) ListByUser(ctx context.Context, clerkID string, limit, offset int) ([]model.Category, error) {
	userID, err := s.userService.GetInternalIDByClerkID(ctx, clerkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal user ID for clerk %s: %w", clerkID, err)
	}
	return s.categoryRepo.ListByUser(ctx, userID, limit, offset)
}

func (s *categoryService) Create(ctx context.Context, clerkID string, category *model.Category) (*model.Category, error) {
	userID, err := s.userService.GetInternalIDByClerkID(ctx, clerkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal user ID for clerk %s: %w", clerkID, err)
	}
	category.UserID = userID

	// Add business logic here if needed (e.g., check for duplicate category names)
	return s.categoryRepo.Create(ctx, category)
}
