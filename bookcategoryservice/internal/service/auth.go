package service

import (
	"context"
	"fmt"

	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/pkg/auth"
	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/auth"
)

type AuthRepository interface {
	GetUserByID(ctx context.Context, userID string) (*pb.User, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type authService struct {
	authRepo  AuthRepository
	jwtSecret string
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(authRepo AuthRepository, jwtSecret string) *authService {
	return &authService{
		authRepo:  authRepo,
		jwtSecret: jwtSecret,
	}
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
	claims, err := auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return "", auth.ErrInvalidToken
	}
	return claims.String(), nil
}
