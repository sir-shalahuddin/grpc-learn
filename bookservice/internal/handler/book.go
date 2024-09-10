package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/pkg/response"
)

type BookService interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBookByID(ctx context.Context, id uuid.UUID) (*models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id uuid.UUID) error
	ListBooks(ctx context.Context, title, author, category string) ([]*models.Book, error)
}

type bookHandler struct {
	bookService BookService
}

func NewBookHandler(bookService BookService) *bookHandler {
	return &bookHandler{bookService: bookService}
}

func (h *bookHandler) CreateBook(c *fiber.Ctx) error {
	userID := c.Locals("id").(uuid.UUID)

	var req struct {
		Title         string     `json:"title"`
		Author        string     `json:"author"`
		ISBN          string     `json:"isbn"`
		PublishedDate *time.Time `json:"published_date"`
		CategoryID    uuid.UUID  `json:"category_id"`
		Stock         int        `json:"stock"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	book := &models.Book{
		ID:            uuid.New(),
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		PublishedDate: req.PublishedDate,
		CategoryID:    req.CategoryID,
		Stock:         req.Stock,
		AddedBy:       userID,
	}

	err := h.bookService.CreateBook(c.Context(), book)
	if err != nil {
		return response.HandleError(c, err, "failed to create book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book created successfully", book, fiber.StatusCreated)
}

func (h *bookHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	book, err := h.bookService.GetBookByID(c.Context(), id)
	if err != nil {
		return response.HandleError(c, err, "failed to retrieve book", fiber.StatusInternalServerError)
	}
	if book == nil {
		return response.HandleError(c, nil, "book not found", fiber.StatusNotFound)
	}

	return response.HandleSuccess(c, "book retrieved successfully", book, fiber.StatusOK)
}

func (h *bookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	var req struct {
		Title         string     `json:"title"`
		Author        string     `json:"author"`
		ISBN          string     `json:"isbn"`
		PublishedDate *time.Time `json:"published_date"`
		CategoryID    uuid.UUID  `json:"category_id"`
		Stock         int        `json:"stock"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	book := &models.Book{
		ID:            id,
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		PublishedDate: req.PublishedDate,
		CategoryID:    req.CategoryID,
		Stock:         req.Stock,
	}

	err = h.bookService.UpdateBook(c.Context(), book)
	if err != nil {
		return response.HandleError(c, err, "failed to update book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book updated successfully", book, fiber.StatusOK)
}

func (h *bookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	err = h.bookService.DeleteBook(c.Context(), id)
	if err != nil {
		return response.HandleError(c, err, "failed to delete book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book deleted successfully", nil, fiber.StatusOK)
}

func (h *bookHandler) ListBooks(c *fiber.Ctx) error {
	title := c.Query("title")
	author := c.Query("author")
	category := c.Query("category")

	books, err := h.bookService.ListBooks(c.Context(), title, author, category)
	if err != nil {
		return response.HandleError(c, err, "failed to retrieve books", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "books retrieved successfully", books, fiber.StatusOK)
}
