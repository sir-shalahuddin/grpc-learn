package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user by ID: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("service: user not found")
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, name, email string) error {
	user := &models.User{
		UserID: userID,
		Name:   name,
		Email:  email,
	}

	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("service: failed to update user: %w", err)
	}
	return nil
}
