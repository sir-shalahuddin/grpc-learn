package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	router "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/routes"
	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/auth"
)

func StartRESTServer(db *sql.DB, authSvc pb.AuthServiceClient, port string) {
	app := fiber.New()

	app.Use(cors.New())

	router.RegisterRoutes(app, db, authSvc)
	err := (app.Listen(":" + port))
	if err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}
}
