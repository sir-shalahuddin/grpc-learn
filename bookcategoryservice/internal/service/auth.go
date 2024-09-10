package service

import (
	"context"
	"fmt"

	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/auth"
)

type AuthRepository interface {
	GetUserByID(ctx context.Context, userID string) (*pb.User, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type authService struct {
	authRepo AuthRepository
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(authRepo AuthRepository) *authService {
	return &authService{authRepo: authRepo}
}

// GetUserByID retrieves a user by their ID.
func (s *authService) GetUserByID(ctx context.Context, userID string) (*pb.User, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// ValidateToken validates a JWT token and retrieves the associated user.
func (s *authService) ValidateToken(ctx context.Context, token string) (string, error) {
	userID, err := s.authRepo.ValidateToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	return userID, nil
}
