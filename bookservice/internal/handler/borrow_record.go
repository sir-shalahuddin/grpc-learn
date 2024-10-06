package handler

import (
	"context"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/pkg/response"
)

type BorrowingRecordService interface {
	BorrowBook(ctx context.Context, req dto.BorrowBookRequest, bookID, userID uuid.UUID) error
	ReturnBook(ctx context.Context, bookID, record_id uuid.UUID) error
	ListBorrowingRecords(ctx context.Context, queries map[string]string, userID uuid.UUID) ([]models.BorrowingRecord, error)
}

type borrowingRecordHandler struct {
	service BorrowingRecordService
}

func NewBorrowingRecordHandler(service BorrowingRecordService) *borrowingRecordHandler {
	return &borrowingRecordHandler{service: service}
}

// BorrowBook godoc
// @Summary Borrow a book
// @Description User borrows a book by providing necessary details
// @Tags Borrowing
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param BorrowBookRequest body dto.BorrowBookRequest true "Borrow Book Request"
// @Success 201 {object} response.Response "Book successfully borrowed"
// @Failure 400 {object} response.ErrorMessage "Invalid book ID or request payload"
// @Failure 500 {object} response.ErrorMessage "Failed to borrow book"
// @Security BearerAuth
// @Router /books/{id}/borrow [post]
func (h *borrowingRecordHandler) BorrowBook(c *fiber.Ctx) error {
	userID := c.Locals("id").(uuid.UUID)

	bookID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	var req dto.BorrowBookRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println(err)
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	if err := h.service.BorrowBook(c.Context(), req, bookID, userID); err != nil {
		log.Println(err)
		return response.HandleError(c, err, "failed to borrow book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book successfully borrowed", nil, fiber.StatusCreated)
}

// ReturnBook godoc
// @Summary Return a borrowed book
// @Description User returns a borrowed book by providing the book and record IDs
// @Tags Borrowing
// @Accept json
// @Produce json
// @Param book_id path string true "Book ID"
// @Param record_id path string true "Borrowing Record ID"
// @Success 200 {object} response.Response "Book returned successfully"
// @Failure 400 {object} response.ErrorMessage "Invalid book or record ID"
// @Failure 404 {object} response.ErrorMessage "Borrowing record not found"
// @Failure 500 {object} response.ErrorMessage "Failed to return book"
// @Security BearerAuth
// @Router /books/{book_id}/records/{record_id} [put]
func (h *borrowingRecordHandler) ReturnBook(c *fiber.Ctx) error {
	bookID, err := uuid.Parse(c.Params("book_id"))
	if err != nil {
		return response.HandleError(c, err, "invalid book ID", fiber.StatusBadRequest)
	}

	recordID, err := uuid.Parse(c.Params("record_id"))
	if err != nil {
		return response.HandleError(c, err, "invalid record ID", fiber.StatusBadRequest)
	}

	err = h.service.ReturnBook(c.Context(), bookID, recordID)
	if err != nil {
		if errors.Is(err, service.ErrBorrowingRecordNotFound) {
			return response.HandleError(c, err, "borrowing record not found", fiber.StatusNotFound)
		}
		log.Println(err)
		return response.HandleError(c, err, "failed to return book", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "book returned successfully", nil, fiber.StatusOK)
}

// ListBorrowingRecords godoc
// @Summary List borrowing records
// @Description Retrieve a list of borrowing records for the user
// @Tags Borrowing
// @Produce json
// @Param limit query int false "Limit the number of records returned"
// @Param offset query int false "Skip a number of records for pagination"
// @Success 200 {object} response.Response "List of borrowing records"
// @Failure 500 {object} response.ErrorMessage "Failed to list borrowing records"
// @Security BearerAuth
// @Router /books/records [get]
func (h *borrowingRecordHandler) ListBorrowingRecords(c *fiber.Ctx) error {
	queries := c.Queries()

	userID := c.Locals("id").(uuid.UUID)

	records, err := h.service.ListBorrowingRecords(c.Context(), queries, userID)
	if err != nil {
		log.Println(err)
		return response.HandleError(c, err, "failed to list borrowing records", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "list of borrowing records", records, fiber.StatusOK)
}
