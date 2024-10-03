package main

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sir-shalahuddin/grpc-learn/userservice/config"
	db "github.com/sir-shalahuddin/grpc-learn/userservice/pkg/database"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no env provided")
	}

	AppConfig := config.AppConfig{
		GRPCPort: config.GetEnv("GRPC_PORT"),
		RESTPort: config.GetEnv("REST_PORT"),
		Mode:     config.GetEnv("SERVER_MODE"),
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

	JWTConfig := config.JWTConfig{
		Secret: config.GetEnv("JWT_SECRET"),
	}

	db, err := db.NewDB(DBConfig)
	if err != nil {
		panic(err)
	}

	mode := strings.ToLower(AppConfig.Mode)
	// Start REST server in a separate goroutine
	if mode == "rest" {
		StartRESTServer(db, AppConfig.RESTPort, JWTConfig.Secret)
	}

	// Start gRPC server in a separate goroutine
	if mode == "grpc" {
		StartGRPCServer(db, AppConfig.GRPCPort, JWTConfig.Secret)
	}
}
