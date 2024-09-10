package repository

import (
	"context"
	"fmt"

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

	fmt.Println("getuser by id : ", resp)

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
		return "", fmt.Errorf("failed to validate token: %w", err)
	}
	fmt.Println("validate token : ", resp)
	return resp.UserId, nil
}
