package server

// import (
// 	"context"

// 	"github.com/google/uuid"
// 	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
// 	pb "github.com/sir-shalahuddin/grpc-learn/userservice/proto"
// )

// type AuthService interface {
// 	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
// 	ValidateToken(tokenStr string) (uuid.UUID, error)
// }

// type authServiceServer struct {
// 	pb.UnimplementedAuthServiceServer
// 	authService AuthService
// }

// func NewAuthServiceServer(authService AuthService) *authServiceServer {
// 	return &authServiceServer{authService: authService}
// }

// func (s *authServiceServer) GetUserByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
// 	userID, err := uuid.Parse(in.GetUserId())
// 	if err != nil {
// 		return nil, err
// 	}

// 	user, err := s.authService.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	userData := pb.User{
// 		UserId: in.GetUserId(),
// 		Email:  user.Email,
// 		Name:   user.Name,
// 		Role:   user.Role,
// 	}

// 	return &pb.GetUserByIDResponse{
// 		User: &userData,
// 	}, nil
// }

// func (s *authServiceServer) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
// 	userID, err := s.authService.ValidateToken(in.GetToken())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pb.ValidateTokenResponse{UserId: userID.String()}, nil

// }
