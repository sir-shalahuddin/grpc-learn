package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/handler"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/service"
	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/auth"
)

// RegisterRoutes sets up the Fiber routes for user management
func RegisterRoutes(app *fiber.App, db *sql.DB, authSvc pb.AuthServiceClient, jwtSecret string) {

	bookcategoryRepo := repository.NewBookCategoryRepository(db)
	bookcategoryService := service.NewBookCategoryService(bookcategoryRepo)
	bookcategoryHandler := handler.NewBookCategoryHandler(bookcategoryService)

	authRepo := repository.NewAuthRepository(authSvc)
	authService := service.NewAuthService(authRepo, jwtSecret)
	authMiddleware := handler.NewAuthMiddleware(authService)

	books := app.Group("/categories")

	// documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Define the routes
	books.Post("/", authMiddleware.Protected("librarian"), bookcategoryHandler.CreateCategory)
	books.Get("/:id", bookcategoryHandler.GetCategoryByID)
	books.Put("/:id", authMiddleware.Protected("librarian"), bookcategoryHandler.UpdateCategory)
	books.Delete("/:id", authMiddleware.Protected("librarian"), bookcategoryHandler.DeleteCategory)
	books.Get("/", bookcategoryHandler.GetAllCategories)
}
