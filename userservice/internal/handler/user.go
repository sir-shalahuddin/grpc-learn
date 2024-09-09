package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, name, email string) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("id").(uuid.UUID)
	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return response.HandleError(c, err, "failed to retrive user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "get profile success", user, fiber.StatusOK)
}

func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("id").(uuid.UUID)

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	err := h.userService.UpdateUser(c.Context(), userID, req.Name, req.Email)
	if err != nil {
		return response.HandleError(c, err, "failed to retrive user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "update profile success", nil, fiber.StatusOK)
}
