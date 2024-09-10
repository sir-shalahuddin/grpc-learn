package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/handler"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/service"
	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/auth"
)

// RegisterRoutes sets up the Fiber routes for user management
func RegisterRoutes(app *fiber.App, db *sql.DB, authSvc pb.AuthServiceClient) {

	bookcategoryRepo := repository.NewBookCategoryRepository(db)
	bookcategoryService := service.NewBookCategoryService(bookcategoryRepo)
	bookcategoryHandler := handler.NewBookCategoryHandler(bookcategoryService)

	authRepo := repository.NewAuthRepository(authSvc)
	authService := service.NewAuthService(authRepo)
	authMiddleware := handler.NewAuthMiddleware(authService)

	books := app.Group("/categories")

	// Define the routes
	books.Post("/", authMiddleware.Protected("librarian"), bookcategoryHandler.CreateCategory)
	books.Get("/:id", bookcategoryHandler.GetCategoryByID)
	books.Put("/:id", authMiddleware.Protected("librarian"), bookcategoryHandler.UpdateCategory)
	books.Delete("/:id", authMiddleware.Protected("librarian"), bookcategoryHandler.DeleteCategory)
	books.Get("/", bookcategoryHandler.GetAllCategories)
}
