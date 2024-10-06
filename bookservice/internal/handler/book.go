package handler

import (
	"context"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/pkg/response"
)

type BookService interface {
	AddBook(ctx context.Context, req dto.AddBookRequest, userID uuid.UUID) error
	GetBookByID(ctx context.Context, id uuid.UUID) (*dto.GetBookResponse, error)
	UpdateBook(ctx context.Context, req dto.UpdateBookRequest, bookID uuid.UUID) error
	DeleteBook(ctx context.Context, id uuid.UUID) error
	ListBooks(ctx context.Context, title, author, category string, page string) ([]*dto.GetBookResponse, error)
}

type bookHandler struct {
	bookService BookService
}

func NewBookHandler(bookService BookService) *bookHandler {
	return &bookHandler{bookService: bookService}
}

// AddBook godoc
// @Summary Add a new book
// @Description Librarian adds a new book
// @Tags Books
// @Accept json
// @Produce json
// @Param AddBookRequest body dto.AddBookRequest true "Add Book Request"
// @Success 201 {object} response.Response "Book successfully added"
// @Failure 400 {object} response.ErrorMessage "Invalid request payload"
// @Failure 409 {object} response.ErrorMessage "Duplicate book"
// @Failure 404 {object} response.ErrorMessage "Category not found"
// @Failure 500 {object} response.ErrorMessage "Internal server error"
// @Security BearerAuth
// @Router /books [post]
func (h *bookHandler) AddBook(c *fiber.Ctx) error {
	userID := c.Locals("id").(uuid.UUID)

	var req dto.AddBookRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	if err := h.bookService.AddBook(c.Context(), req, userID); err != nil {
		if errors.Is(err, service.ErrBookDuplicate) {
			return response.HandleError(c, err, "", fiber.StatusConflict)
		}
		if errors.Is(err, service.ErrCategoryNotFound) {
			return response.HandleError(c, err, "", fiber.StatusNotFound)
		}
		log.Println(err)
		return response.HandleError(c, err, "failed to add book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book successfully added", nil, fiber.StatusCreated)
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Retrieves a book by its ID
// @Tags Books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} response.Response "Book retrieved successfully"
// @Failure 400 {object} response.ErrorMessage "Invalid book ID"
// @Failure 404 {object} response.ErrorMessage "Book not found"
// @Failure 500 {object} response.ErrorMessage "Internal server error"
// @Router /books/{id} [get]
func (h *bookHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	book, err := h.bookService.GetBookByID(c.Context(), id)
	if err != nil {
		log.Println(err)
		return response.HandleError(c, err, "failed to retrieve book", fiber.StatusInternalServerError)
	}
	if book == nil {
		return response.HandleError(c, nil, "book not found", fiber.StatusNotFound)
	}

	return response.HandleSuccess(c, "book retrieved successfully", book, fiber.StatusOK)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Librarian updates an existing book
// @Tags Books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param UpdateBookRequest body dto.UpdateBookRequest true "Update Book Request"
// @Success 200 {object} response.Response "Book updated successfully"
// @Failure 400 {object} response.ErrorMessage "Invalid request payload or book ID"
// @Failure 500 {object} response.ErrorMessage "Internal server error"
// @Security BearerAuth
// @Router /books/{id} [put]
func (h *bookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	var req dto.UpdateBookRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	err = h.bookService.UpdateBook(c.Context(), req, id)
	if err != nil {
		log.Println(err)
		return response.HandleError(c, err, "failed to update book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book updated successfully", "", fiber.StatusOK)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Librarian deletes an existing book
// @Tags Books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} response.Response "Book deleted successfully"
// @Failure 400 {object} response.ErrorMessage "Invalid book ID"
// @Failure 500 {object} response.ErrorMessage "Internal server error"
// @Security BearerAuth
// @Router /books/{id} [delete]
func (h *bookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	err = h.bookService.DeleteBook(c.Context(), id)
	if err != nil {
		log.Println(err)
		return response.HandleError(c, err, "failed to delete book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book deleted successfully", nil, fiber.StatusOK)
}

// ListBooks godoc
// @Summary List books
// @Description Retrieves a list of books, optionally filtered by title, author, or category
// @Tags Books
// @Produce json
// @Param title query string false "Book title"
// @Param author query string false "Book author"
// @Param category query string false "Book category"
// @Param page query string false "Page number"
// @Success 200 {object} response.Response "Books retrieved successfully"
// @Failure 500 {object} response.ErrorMessage "Internal server error"
// @Router /books [get]
func (h *bookHandler) ListBooks(c *fiber.Ctx) error {
	title := c.Query("title")
	author := c.Query("author")
	category := c.Query("category")
	page := c.Query("page", "1") // Default to page 1 if not provided

	books, err := h.bookService.ListBooks(c.Context(), title, author, category, page)
	if err != nil {
		log.Println(err)
		return response.HandleError(c, err, "failed to retrieve books", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "books retrieved successfully", books, fiber.StatusOK)
}
