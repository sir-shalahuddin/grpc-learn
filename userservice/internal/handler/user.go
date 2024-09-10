package handler

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, name, email string) error
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *userHandler {
	return &userHandler{userService: userService}
}

// GetProfile retrieves the profile of the currently authenticated user.
func (h *userHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uuid.UUID)
	if !ok {
		return response.HandleError(c, fmt.Errorf("invalid user ID"), "failed to retrieve user", fiber.StatusInternalServerError)
	}

	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return response.HandleError(c, err, "failed to retrieve user", fiber.StatusInternalServerError)
	}

	if user == nil {
		return response.HandleError(c, fmt.Errorf("user not found"), "user not found", fiber.StatusNotFound)
	}

	return response.HandleSuccess(c, "get profile success", user, fiber.StatusOK)
}

// UpdateProfile updates the profile of the currently authenticated user.
func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uuid.UUID)
	if !ok {
		return response.HandleError(c, fmt.Errorf("invalid user ID"), "failed to update profile", fiber.StatusInternalServerError)
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	err := h.userService.UpdateUser(c.Context(), userID, req.Name, req.Email)
	if err != nil {
		return response.HandleError(c, err, "failed to update user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "update profile success", nil, fiber.StatusOK)
}
