package setup

import (
	"database/sql"

	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/server"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/service"
	pb "github.com/sir-shalahuddin/grpc-learn/userservice/proto"
	"google.golang.org/grpc"
)

func GRPCServer(grpc *grpc.Server, db *sql.DB, jwtSecret string) {
	// Initialize repositories, services, and servers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authServer := server.NewAuthServiceServer(authService)

	// Register AuthService routes
	pb.RegisterAuthServiceServer(grpc, authServer)

}
