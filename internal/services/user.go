package services

import (
	"context"

	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/internal/repositories"
)

type UserService interface {
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	return s.userRepo.List(ctx, limit, offset)
}

func (s *userService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return s.userRepo.Create(ctx, user)
}
