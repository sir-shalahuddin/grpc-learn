package handler

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

type AdminService interface {
	ListUsers(ctx context.Context) ([]*models.User, error)
	UpdateUserRoles(ctx context.Context, userID uuid.UUID, roles string) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type adminHandler struct {
	adminService AdminService
}

func NewAdminHandler(adminService AdminService) *adminHandler {
	return &adminHandler{adminService: adminService}
}

func (h *adminHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.adminService.ListUsers(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list users")
	}

	return response.HandleSuccess(c, "retrieve list user success", users, fiber.StatusOK)
}

func (h *adminHandler) UpdateUserRoles(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "failed to parse id : invalid request params", fiber.StatusBadRequest)
	}

	var req struct {
		Roles string `json:"roles"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	if err := h.adminService.UpdateUserRoles(context.Background(), userID, req.Roles); err != nil {
		// fmt.Println(err)
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return response.HandleError(c, err, "", fiber.StatusForbidden)
		}
		return response.HandleError(c, err, "failed to update user roles", fiber.StatusInternalServerError)
	}
	return response.HandleSuccess(c, "update user role success", nil, fiber.StatusOK)
}

func (h *adminHandler) DeleteUser(c *fiber.Ctx) error {

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "failed to parse id : invalid request params", fiber.StatusBadRequest)
	}

	if err := h.adminService.DeleteUser(context.Background(), userID); err != nil {
		return response.HandleError(c, err, "failed to delete user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "delete user success", nil, fiber.StatusOK)
}
