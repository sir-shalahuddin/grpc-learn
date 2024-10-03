package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/handler"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/proto/authservice"
	pb "github.com/sir-shalahuddin/grpc-learn/bookservice/proto/categoryservice"
)

// RegisterRoutes sets up the Fiber routes for user management
func RegisterRoutes(app *fiber.App, db *sql.DB, authSvc authservice.AuthServiceClient, ctgSvc pb.BookCategoryServiceClient) {
	txRepo := repository.NewTxRepository(db)

	bookRepo := repository.NewBookRepository(db)
	ctgRepo := repository.NewCategoryRepository(ctgSvc)
	bookService := service.NewBookService(bookRepo, ctgRepo)
	bookHandler := handler.NewBookHandler(bookService)

	authRepo := repository.NewAuthRepository(authSvc)
	authService := service.NewAuthService(authRepo)
	authMiddleware := handler.NewAuthMiddleware(authService)

	borrowingRecordRepo := repository.NewBorrowingRecordRepository(db)
	borrowingRecordService := service.NewBorrowingRecordService(borrowingRecordRepo, txRepo, bookRepo)
	borrowingRecordHandler := handler.NewBorrowingRecordHandler(borrowingRecordService)

	books := app.Group("/books")

	books.Post("/:id/borrow", authMiddleware.Protected("user"), borrowingRecordHandler.BorrowBook)
	books.Put("/:book_id/records/:record_id", authMiddleware.Protected("user"), borrowingRecordHandler.ReturnBook)
	books.Get("/records", authMiddleware.Protected("user"), borrowingRecordHandler.ListBorrowingRecords)

	books.Post("/", authMiddleware.Protected("librarian"), bookHandler.AddBook)
	books.Get("/:id", bookHandler.GetBookByID)
	books.Put("/:id", authMiddleware.Protected("librarian"), bookHandler.UpdateBook)
	books.Delete("/:id", authMiddleware.Protected("librarian"), bookHandler.DeleteBook)
	books.Get("/", bookHandler.ListBooks)

}
