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
		Host:                   config.GetEnv("DB_HOST"),
		Port:                   config.GetEnv("DB_PORT"),
		User:                   config.GetEnv("DB_USER"),
		Pass:                   config.GetEnv("DB_PASS"),
		Name:                   config.GetEnv("DB_NAME"),
		InstanceConnectionName: config.GetEnv("INSTANCE_CONNECTION_NAME"),
		UseUnixSocket:          config.GetEnvAsBool("USE_UNIX_SOCKET"),
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

}
