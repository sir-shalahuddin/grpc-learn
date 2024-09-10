package handler

import (
	"context"
	"errors"
	"log"

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
		log.Printf("Error listing users: %v", err)
		return response.HandleError(c, err, "Failed to list users", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "Retrieve list of users successful", users, fiber.StatusOK)
}

func (h *adminHandler) UpdateUserRoles(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "Invalid request parameters", fiber.StatusBadRequest)
	}

	var req struct {
		Roles string `json:"roles"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "Invalid request payload", fiber.StatusBadRequest)
	}

	err = h.adminService.UpdateUserRoles(context.Background(), userID, req.Roles)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return response.HandleError(c, err, "Insufficient permissions", fiber.StatusForbidden)
		}
		log.Printf("Error updating user roles: %v", err)
		return response.HandleError(c, err, "Failed to update user roles", fiber.StatusInternalServerError)
	}
	return response.HandleSuccess(c, "Update user role successful", nil, fiber.StatusOK)
}

func (h *adminHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.HandleError(c, err, "Invalid request parameters", fiber.StatusBadRequest)
	}

	err = h.adminService.DeleteUser(context.Background(), userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return response.HandleError(c, err, "Failed to delete user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "Delete user successful", nil, fiber.StatusOK)
}
