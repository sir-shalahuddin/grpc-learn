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
	ReturnBook(ctx context.Context, recordID uuid.UUID) error
	ListBorrowingRecordsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.BorrowingRecord, error)
}

type borrowingRecordHandler struct {
	service BorrowingRecordService
}

func NewBorrowingRecordHandler(service BorrowingRecordService) *borrowingRecordHandler {
	return &borrowingRecordHandler{service: service}
}

func (h *borrowingRecordHandler) BorrowBook(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uuid.UUID)
	if !ok {
		return response.HandleError(c, errors.New("user ID not found"), "invalid user", fiber.StatusUnauthorized)
	}

	bookID, err := uuid.Parse(c.Params("book_id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	var req struct {
		DueDate *time.Time `json:"due_date"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	err = h.service.BorrowBook(c.Context(), bookID, userID, req.DueDate)
	if err != nil {
		if errors.Is(err, service.ErrBookUnavailable) {
			return response.HandleError(c, err, "book is not available", fiber.StatusConflict)
		}
		return response.HandleError(c, err, "failed to borrow book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book borrowed successfully", nil, fiber.StatusOK)
}

func (h *borrowingRecordHandler) ReturnBook(c *fiber.Ctx) error {
	recordID, err := uuid.Parse(c.Params("record_id"))
	if err != nil {
		return response.HandleError(c, err, "invalid record ID", fiber.StatusBadRequest)
	}

	err = h.service.ReturnBook(c.Context(), recordID)
	if err != nil {
		if errors.Is(err, service.ErrBorrowingRecordNotFound) {
			return response.HandleError(c, err, "borrowing record not found", fiber.StatusNotFound)
		}
		return response.HandleError(c, err, "failed to return book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book returned successfully", nil, fiber.StatusOK)
}

func (h *borrowingRecordHandler) ListBorrowingRecords(c *fiber.Ctx) error {
	bookID, err := uuid.Parse(c.Query("book_id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	records, err := h.service.ListBorrowingRecordsByBookID(c.Context(), bookID)
	if err != nil {
		return response.HandleError(c, err, "failed to list borrowing records", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "list of borrowing records", records, fiber.StatusOK)
}
