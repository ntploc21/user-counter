package services

import (
	"context"
	"locntp-user-counter/internal/models"
	"locntp-user-counter/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username string) (*models.User, error) {
	return s.repo.CreateUser(ctx, &models.CreateUserRequest{Username: username})
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserCounter(ctx context.Context, id uint) (int64, bool, error) {
	return s.repo.GetUserCount(ctx, id)
}

func (s *UserService) IncrementCounter(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.IncrementUserCount(ctx, id, 1)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteUser(ctx, id)
}
