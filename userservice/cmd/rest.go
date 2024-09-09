package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sir-shalahuddin/grpc-learn/userservice/config"
	router "github.com/sir-shalahuddin/grpc-learn/userservice/internal/routes"
)

func StartRESTServer(db *sql.DB, appConfig config.AppConfig, jwtConfig config.JWTConfig) {
	app := fiber.New()

	app.Use(cors.New())

	router.RegisterRoutes(app, db, jwtConfig)
	err := (app.Listen(":" + appConfig.Port))
	if err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}
}
