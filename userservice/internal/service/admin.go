package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrInsufficientPermissions = errors.New("insufficient permissions to action to super admin")
)

type AdminRepository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context) ([]dto.GetUser, error)
	UpdateUserRoles(ctx context.Context, userID uuid.UUID, roles string) error
}

type adminService struct {
	repo AdminRepository
}

func NewAdminService(repo AdminRepository) *adminService {
	return &adminService{
		repo: repo,
	}
}

// ListUsers retrieves all users from the repository.
func (s *adminService) ListUsers(ctx context.Context) ([]dto.GetUser, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUserRoles updates the roles of a user, ensuring the user is found and
// preventing promotion to "super admin" without sufficient permissions.
func (s *adminService) UpdateUserRoles(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRoles) error {
	// Prevent promotion to "super admin" by unauthorized users
	if req.Role == "super admin" {
		return ErrInsufficientPermissions
	}

	// Retrieve the user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Update the user's roles
	err = s.repo.UpdateUserRoles(ctx, userID, req.Role)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser checks if a user exists and then deletes them.
func (s *adminService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if user.Role == "super admin" {
		return ErrInsufficientPermissions
	}
	

	// Delete the user
	err = s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
