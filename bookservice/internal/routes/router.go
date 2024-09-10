package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/handler"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/proto/authservice"
)

// RegisterRoutes sets up the Fiber routes for user management
func RegisterRoutes(app *fiber.App, db *sql.DB, authSvc authservice.AuthServiceClient) {

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	authRepo := repository.NewAuthRepository(authSvc)
	authService := service.NewAuthService(authRepo)
	authMiddleware := handler.NewAuthMiddleware(authService)

	books := app.Group("/books")

	// Define the routes
	books.Post("/", authMiddleware.Protected("librarian"), bookHandler.CreateBook)      // Create a new book
	books.Get("/:id", bookHandler.GetBookByID)                                          // Get a book by ID
	books.Put("/:id", authMiddleware.Protected("librarian"), bookHandler.UpdateBook)    // Update a book by ID
	books.Delete("/:id", authMiddleware.Protected("librarian"), bookHandler.DeleteBook) // Delete a book by ID
	books.Get("/", bookHandler.ListBooks)
}
