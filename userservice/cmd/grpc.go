package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/sir-shalahuddin/grpc-learn/userservice/config"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/server"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/service"
	pb "github.com/sir-shalahuddin/grpc-learn/userservice/proto"
	"google.golang.org/grpc"
)

// StartGRPCServer initializes and starts the gRPC server
func StartGRPCServer(db *sql.DB, jwtConfig config.JWTConfig) {

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtConfig.Secret)
	authServer := server.NewAuthServiceServer(authService)

	lis, err := net.Listen("tcp", "localhost:3001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register AuthService routes
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	log.Printf("gRPC server listening on %s", ":3001")

	// Start the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
