package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/pkg/response"
)

type BorrowingRecordService interface {
	BorrowBook(ctx context.Context, bookID, userID uuid.UUID, dueDate *time.Time) error
	ReturnBook(ctx context.Context, bookID uuid.UUID) error
	ListBorrowingRecordsByBookID(ctx context.Context, bookID uuid.UUID) ([]models.BorrowingRecord, error)
}

type borrowingRecordHandler struct {
	service BorrowingRecordService
}

func NewBorrowingRecordHandler(service BorrowingRecordService) *borrowingRecordHandler {
	return &borrowingRecordHandler{service: service}
}

// BorrowBook handles the borrowing of a book by a user.
func (h *borrowingRecordHandler) BorrowBook(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("id").(string)
	if !ok {
		return response.HandleError(c, errors.New("user ID not found"), "invalid user", fiber.StatusUnauthorized)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.HandleError(c, err, "invalid user ID", fiber.StatusBadRequest)
	}

	// Parse request body
	type Request struct {
		BookID  uuid.UUID  `json:"book_id"`
		DueDate *time.Time `json:"due_date,omitempty"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Call the service to borrow a book
	err = h.service.BorrowBook(c.Context(), req.BookID, userID, req.DueDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book borrowed successfully",
	})
}

// ReturnBook handles the return of a borrowed book using the book ID.
func (h *borrowingRecordHandler) ReturnBook(c *fiber.Ctx) error {
	bookID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	err = h.service.ReturnBook(c.Context(), bookID)
	if err != nil {
		if errors.Is(err, service.ErrBorrowingRecordNotFound) {
			return response.HandleError(c, err, "borrowing record not found", fiber.StatusNotFound)
		}
		return response.HandleError(c, err, "failed to return book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book returned successfully", nil, fiber.StatusOK)
}

// ListBorrowingRecords lists all borrowing records for a given book ID.
func (h *borrowingRecordHandler) ListBorrowingRecords(c *fiber.Ctx) error {
	bookID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	records, err := h.service.ListBorrowingRecordsByBookID(c.Context(), bookID)
	if err != nil {
		return response.HandleError(c, err, "failed to list borrowing records", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "list of borrowing records", records, fiber.StatusOK)
}
