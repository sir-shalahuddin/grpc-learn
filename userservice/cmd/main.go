package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
		Host: config.GetEnv("DB_HOST"),
		Port: config.GetEnv("DB_PORT"),
		User: config.GetEnv("DB_USER"),
		Pass: config.GetEnv("DB_PASS"),
		Name: config.GetEnv("DB_NAME"),
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
	if mode == "rest" || mode == "hybrid" {
		go StartRESTServer(db, AppConfig.RESTPort, JWTConfig.Secret)
	}

	// Start gRPC server in a separate goroutine
	if mode == "grpc" || mode == "hybrid" {
		go StartGRPCServer(db, AppConfig.GRPCPort, JWTConfig.Secret)
	}

	// Handle termination signals to gracefully shutdown servers
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down servers...")
}
