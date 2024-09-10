package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/sir-shalahuddin/grpc-learn/bookservice/proto/authservice"
)

type authRepository struct {
	grpc authservice.AuthServiceClient
}

func NewAuthRepository(grpc authservice.AuthServiceClient) *authRepository {
	return &authRepository{grpc: grpc}
}

// GetUserByID retrieves a user by their ID using the AuthService gRPC client.
func (r *authRepository) GetUserByID(ctx context.Context, userID string) (*authservice.User, error) {
	// Prepare the gRPC request
	req := &authservice.GetUserByIDRequest{
		UserId: userID,
	}

	// Call the gRPC service
	resp, err := r.grpc.GetUserByID(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return resp.User, nil
}

// ValidateToken validates a JWT token using the AuthService gRPC client and returns the associated user ID if valid.
func (r *authRepository) ValidateToken(ctx context.Context, token string) (string, error) {
	// Prepare the gRPC request
	req := &authservice.ValidateTokenRequest{
		Token: token,
	}

	// Call the gRPC service
	resp, err := r.grpc.ValidateToken(ctx, req)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	return resp.UserId, nil
}
