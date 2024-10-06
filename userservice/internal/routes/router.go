package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/handler"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/service"
)

// RegisterRoutes sets up the Fiber routes for user management
func RegisterRoutes(app *fiber.App, db *sql.DB, jwtSecret string) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	adminService := service.NewAdminService(userRepo)
	adminHandler := handler.NewAdminHandler(adminService)

	authMiddleware := handler.NewAuthMiddleware(authService)

	// documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Authentication routes
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh-token", authHandler.RefreshToken)

	// User routes
	profile := app.Group("/profile", authMiddleware.Protected())
	profile.Get("/", userHandler.GetProfile)
	profile.Put("/", userHandler.UpdateProfile)

	// Admin routes
	admin := app.Group("/admin", authMiddleware.Protected("super admin"))
	admin.Get("/users", adminHandler.ListUsers)
	admin.Put("/users/:id/roles", adminHandler.UpdateUserRoles)
	admin.Delete("/users/:id", adminHandler.DeleteUser)
}
