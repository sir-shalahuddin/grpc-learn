package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	router "github.com/sir-shalahuddin/grpc-learn/bookservice/internal/routes"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/proto/authservice"
)

func StartRESTServer(db *sql.DB, authSvc authservice.AuthServiceClient, port string) {
	app := fiber.New()

	app.Use(cors.New())

	router.RegisterRoutes(app, db, authSvc)
	err := (app.Listen(":" + port))
	if err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}
}
