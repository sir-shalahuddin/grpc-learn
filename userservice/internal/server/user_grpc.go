package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
	pb "github.com/sir-shalahuddin/grpc-learn/userservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ValidateToken(tokenStr string) (uuid.UUID, error)
}

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
	authService AuthService
}

func NewAuthServiceServer(authService AuthService) *authServiceServer {
	return &authServiceServer{authService: authService}
}

func (s *authServiceServer) GetUserByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID format: %v", err)
	}

	user, err := s.authService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user by ID: %v", err)
	}

	userData := &pb.User{
		UserId: user.UserID.String(),
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
	}

	return &pb.GetUserByIDResponse{
		User: userData,
	}, nil
}

func (s *authServiceServer) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, err := s.authService.ValidateToken(in.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid or expired token: %v", err)
	}

	return &pb.ValidateTokenResponse{UserId: userID.String()}, nil
}
