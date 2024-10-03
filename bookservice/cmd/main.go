package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/config"
	grpcclient "github.com/sir-shalahuddin/grpc-learn/bookservice/pkg/datasource/grpc"
	db "github.com/sir-shalahuddin/grpc-learn/bookservice/pkg/datasource/postgres"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no env provided")
	}

	AppConfig := config.AppConfig{
		// GRPCPort: config.GetEnv("GRPC_PORT"),
		RESTPort: config.GetEnv("REST_PORT"),
	}
	DBConfig := config.DBConfig{
		Host: config.GetEnv("DB_HOST"),
		Port: config.GetEnv("DB_PORT"),
		User: config.GetEnv("DB_USER"),
		Pass: config.GetEnv("DB_PASS"),
		Name: config.GetEnv("DB_NAME"),
	}

	GRPCConfig := config.GRPCConfig{
		AuthAddress:     config.GetEnv("AUTH_ADDRESS"),
		CategoryAddress: config.GetEnv("CTG_ADDRESS"),
	}

	db, err := db.NewDB(DBConfig)
	if err != nil {
		panic(err)
	}

	authClients, err := grpcclient.NewAuthClients(GRPCConfig.AuthAddress)
	if err != nil {
		panic(err)
	}

	categoryClients, err := grpcclient.NewCategoryClients(GRPCConfig.CategoryAddress)
	if err != nil {
		panic(err)
	}

	// Start REST server in a separate goroutine
	StartRESTServer(db, authClients, categoryClients, AppConfig.RESTPort)

	// Start gRPC server in a separate goroutine
	// go StartGRPCServer(db, AppConfig.GRPCPort, JWTConfig.Secret)

	// Handle termination signals to gracefully shutdown servers
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// <-sigs

	// log.Println("Shutting down servers...")
}
