package main

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sir-shalahuddin/grpc-learn/userservice/config"
	_ "github.com/sir-shalahuddin/grpc-learn/userservice/docs"
	db "github.com/sir-shalahuddin/grpc-learn/userservice/pkg/database"
)

// @title User Service API
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host user-rest.sirlearn.my.id
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @Login
// @description API for managing users and profiles
// @description Login as Super Admin or Librarian with the following credentials:
// @description - **Super Admin**: `superadmin@mail.com`, Password: `Password123!`
// @description - **Librarian**: `librarian@mail.com`, Password: `Password123!`

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
