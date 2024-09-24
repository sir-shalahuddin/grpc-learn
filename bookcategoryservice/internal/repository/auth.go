package repository

import (
	"context"
	"fmt"
	"log"

	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/auth"
)

type authRepository struct {
	grpc pb.AuthServiceClient
}

func NewAuthRepository(grpc pb.AuthServiceClient) *authRepository {
	return &authRepository{grpc: grpc}
}

// GetUserByID retrieves a user by their ID using the AuthService gRPC client.
func (r *authRepository) GetUserByID(ctx context.Context, userID string) (*pb.User, error) {
	// Prepare the gRPC request
	req := &pb.GetUserByIDRequest{
		UserId: userID,
	}

	// Call the gRPC service
	resp, err := r.grpc.GetUserByID(ctx, req)
	if err != nil {
		log.Printf("[AuthRepository - GetUserByID] Error calling GetUserByID: %v", err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return resp.User, nil
}

// ValidateToken validates a JWT token using the AuthService gRPC client and returns the associated user ID if valid.
func (r *authRepository) ValidateToken(ctx context.Context, token string) (string, error) {
	// Prepare the gRPC request
	req := &pb.ValidateTokenRequest{
		Token: token,
	}

	// Call the gRPC service
	resp, err := r.grpc.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("[AuthRepository - ValidateToken] Error validating token: %v", err)
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	return resp.UserId, nil
}
