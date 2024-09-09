package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)
	UpdateUserRoles(ctx context.Context, userID uuid.UUID, roles string) error
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
		return nil, err
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
		return err
	}
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
