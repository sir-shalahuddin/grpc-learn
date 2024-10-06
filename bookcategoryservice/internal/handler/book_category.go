package handler

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/dto"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/models"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/pkg/response"
)

type bookCategoryService interface {
	CreateCategory(ctx context.Context, req dto.CreateBookCategoryRequest) (dto.CreateBookCategoryResponse, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error)
	GetAllCategories(ctx context.Context) ([]models.BookCategory, error)
	UpdateCategory(ctx context.Context, category *models.BookCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type bookCategoryHandler struct {
	service  bookCategoryService
	validate *validator.Validate
}

// NewBookCategoryHandler creates a new instance of BookCategoryHandler.
func NewBookCategoryHandler(service bookCategoryService) *bookCategoryHandler {
	validate := validator.New()

	return &bookCategoryHandler{
		service:  service,
		validate: validate}
}

// CreateCategory handles the creation of a new book category.
// @Summary Create a new book category
// @Description Create a new book category with the provided details
// @Tags BookCategory
// @Accept json
// @Produce json
// @Param category body dto.CreateBookCategoryRequest true "Category details"
// @Success 201 {object} response.Response "success to create category"
// @Failure 400 {object} response.ErrorMessage "invalid payload"
// @Failure 409 {object} response.ErrorMessage "category already exists"
// @Failure 500 {object} response.ErrorMessage "failed to create book category"
// @Router /categories [post]
// @Security BearerAuth
func (h *bookCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateBookCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
	}

	res, err := h.service.CreateCategory(context.Background(), req)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateCategory) {
			return response.HandleError(c, err, "", fiber.StatusConflict)
		}
		return response.HandleError(c, err, "failed to create book category", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "success to create category", res, fiber.StatusCreated)
}

// GetCategoryByID retrieves a book category by its ID.
// @Summary Retrieve a book category by ID
// @Description Get the details of a book category by its ID
// @Tags BookCategory
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response "success retrieve category"
// @Failure 400 {object} response.ErrorMessage "invalid category ID"
// @Failure 404 {object} response.ErrorMessage "category not found"
// @Failure 500 {object} response.ErrorMessage "failed to retrieve category"
// @Router /categories/{id} [get]
func (h *bookCategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return response.HandleError(c, err, "invalid category ID", fiber.StatusBadRequest)
	}

	category, err := h.service.GetCategoryByID(context.Background(), id)
	if err != nil {
		return response.HandleError(c, err, "failed to retrieve category", fiber.StatusInternalServerError)
	}

	if category == nil {
		return response.HandleError(c, err, "category not found", fiber.StatusNotFound)
	}

	return response.HandleSuccess(c, "success retrieve category", category, fiber.StatusOK)
}

// GetAllCategories retrieves all book categories.
// @Summary Retrieve all book categories
// @Description Get a list of all book categories
// @Tags BookCategory
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "success to retrieve categories"
// @Failure 500 {object} response.ErrorMessage "failed to retrieve categories"
// @Router /categories [get]
func (h *bookCategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories(context.Background())
	if err != nil {
		return response.HandleError(c, err, "failed to retrieve categories", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "success to retrive categories", categories, fiber.StatusOK)
}

// UpdateCategory updates an existing book category by its ID.
// @Summary Update a book category
// @Description Update the details of a book category by its ID
// @Tags BookCategory
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body dto.UpdateBookCategoryRequest true "Updated category details"
// @Success 200 {object} response.Response "category updated successfully"
// @Failure 400 {object} response.ErrorMessage "invalid category ID or payload"
// @Failure 500 {object} response.ErrorMessage "failed to update category"
// @Router /categories/{id} [put]
// @Security BearerAuth
func (h *bookCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return response.HandleError(c, err, "invalid category ID", fiber.StatusBadRequest)
	}

	var req dto.UpdateBookCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
	}

	category := &models.BookCategory{
		ID:   id,
		Name: req.Name,
	}

	if err := h.service.UpdateCategory(context.Background(), category); err != nil {
		return response.HandleError(c, err, "failed to update category", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "category updated successfully", nil, fiber.StatusOK)
}

// DeleteCategory deletes a book category by its ID.
// @Summary Delete a book category
// @Description Remove a book category by its ID
// @Tags BookCategory
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response "category deleted successfully"
// @Failure 400 {object} response.ErrorMessage "invalid category ID"
// @Failure 500 {object} response.ErrorMessage "failed to delete category"
// @Router /categories/{id} [delete]
// @Security BearerAuth
func (h *bookCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return response.HandleError(c, err, "invalid category ID", fiber.StatusBadRequest)
	}

	if err := h.service.DeleteCategory(context.Background(), id); err != nil {
		return response.HandleError(c, err, "failed to update category", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "category deleted successfully", nil, fiber.StatusOK)
}
