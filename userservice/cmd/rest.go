package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	router "github.com/sir-shalahuddin/grpc-learn/userservice/internal/routes"
)

func StartRESTServer(db *sql.DB, port string, jwtSecret string) {
	app := fiber.New()

	app.Use(cors.New())

	router.RegisterRoutes(app, db, jwtSecret)
	err := (app.Listen(":" + port))
	if err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}
}
