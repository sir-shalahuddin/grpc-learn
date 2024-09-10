package grpcclient

import (
	"fmt"
	"log"

	pb "github.com/sir-shalahuddin/grpc-learn/bookservice/proto/authservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClients(address string) (pb.AuthServiceClient, error) {
	conn, err := grpc.NewClient(
		address, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	// Create a new AuthService client for the server.
	client := pb.NewAuthServiceClient(conn)
	log.Printf("Connected to gRPC server at %s", address)

	return client, nil
}
