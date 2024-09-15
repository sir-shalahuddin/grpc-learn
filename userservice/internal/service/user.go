package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (dto.GetProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return dto.GetProfileResponse{}, fmt.Errorf("service: failed to get user by ID: %w", err)
	}

	if user == nil {
		return dto.GetProfileResponse{}, ErrUserNotFound
	}

	return dto.GetProfileResponse{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) error {

	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil && existingUser.UserID != userID {
		return ErrDuplicateEmail
	}

	user := models.User{
		UserID: userID,
		Email:  req.Email,
		Name:   req.Name,
	}

	if err := s.repo.UpdateUser(ctx, &user); err != nil {
		return fmt.Errorf("service: failed to update user: %w", err)
	}

	return nil
}
