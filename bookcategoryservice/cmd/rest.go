package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	router "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/routes"
	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/auth"
)

func StartRESTServer(db *sql.DB, authSvc pb.AuthServiceClient, port string, jwtSecret string) {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	router.RegisterRoutes(app, db, authSvc, jwtSecret)
	err := (app.Listen(":" + port))
	if err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}
}
