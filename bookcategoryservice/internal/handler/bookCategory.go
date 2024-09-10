package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/models"
)

type bookCategoryService interface {
	CreateCategory(ctx context.Context, name string) (uuid.UUID, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error)
	GetAllCategories(ctx context.Context) ([]models.BookCategory, error)
	UpdateCategory(ctx context.Context, category *models.BookCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type bookCategoryHandler struct {
	service bookCategoryService
}

// NewBookCategoryHandler creates a new instance of BookCategoryHandler.
func NewBookCategoryHandler(service bookCategoryService) *bookCategoryHandler {
	return &bookCategoryHandler{service: service}
}

func (h *bookCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name cannot be empty",
		})
	}

	id, err := h.service.CreateCategory(c.Context(), req.Name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":   id,
		"name": req.Name,
	})
}

func (h *bookCategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	category, err := h.service.GetCategoryByID(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve category",
		})
	}

	if category == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.JSON(category)
}

func (h *bookCategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve categories",
		})
	}

	return c.JSON(categories)
}

func (h *bookCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name cannot be empty",
		})
	}

	category := &models.BookCategory{
		ID:   id,
		Name: req.Name,
	}

	if err := h.service.UpdateCategory(c.Context(), category); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category updated successfully",
	})
}

func (h *bookCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	if err := h.service.DeleteCategory(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
