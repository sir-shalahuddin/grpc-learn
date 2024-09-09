package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInsufficientPermissions = errors.New("insufficient permissions to promote to super admin")
)

type adminService struct {
	userRepo UserRepository
}

func NewAdminService(userRepo UserRepository) *adminService {
	return &adminService{
		userRepo: userRepo,
	}
}

func (s *adminService) ListUsers(ctx context.Context) ([]*models.User, error) {
	// List all users from the repository
	return s.userRepo.ListUsers(ctx)
}

func (s *adminService) UpdateUserRoles(ctx context.Context, userID uuid.UUID, roles string) error {
	// Retrieve the user

	if roles == "super admin"{
		return ErrInsufficientPermissions
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	err = s.userRepo.UpdateUserRoles(ctx, userID, roles)
	if err != nil {
		return err
	}

	return nil
}

func (s *adminService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Delete the user
	return s.userRepo.DeleteUser(ctx, userID)
}
