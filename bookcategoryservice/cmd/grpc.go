package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/setup"
	"google.golang.org/grpc"
)

// StartGRPCServer initializes and starts the gRPC server
func StartGRPCServer(db *sql.DB, port string) {

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register AuthService routes
	setup.GRPCServer(grpcServer, db)
	log.Printf("gRPC server listening on %s", ":"+port)

	// Start the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
