package handler

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (dto.GetProfileResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) error
}

type userHandler struct {
	userService UserService
	validate    *validator.Validate
}

func NewUserHandler(userService UserService) *userHandler {
	validate := validator.New()
	validate.RegisterValidation("password", ValidatePassword)

	return &userHandler{
		userService: userService,
		validate:    validate}
}

// GetProfile retrieves the profile of the currently authenticated user.
// @Summary Get user profile
// @Description Retrieves the profile of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.GetProfileResponse} "Profile retrieved successfully"
// @Failure 400 {object} response.ErrorMessage "Invalid user ID"
// @Failure 404 {object} response.ErrorMessage "User not found"
// @Failure 500 {object} response.ErrorMessage "Failed to retrieve user"
// @Router /profile [get]
// @Security BearerAuth
func (h *userHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uuid.UUID)
	if !ok {
		return response.HandleError(c, fmt.Errorf("invalid user ID"), "failed to retrieve user", fiber.StatusInternalServerError)
	}

	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.HandleError(c, err, "user not found", fiber.StatusNotFound)
		}
		return response.HandleError(c, err, "failed to retrieve user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "get profile success", user, fiber.StatusOK)
}

// UpdateProfile updates the profile of the currently authenticated user.
// @Summary Update user profile
// @Description Updates the profile of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param updateProfileRequest body dto.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} response.Response "Profile updated successfully"
// @Failure 400 {object} response.ErrorMessage "Invalid request payload"
// @Failure 409 {object} response.ErrorMessage "Duplicate email"
// @Failure 500 {object} response.ErrorMessage "Failed to update profile"
// @Router /profile [put]
// @Security BearerAuth
func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uuid.UUID)
	if !ok {
		log.Println("invalid user ID")
		return response.HandleError(c, fmt.Errorf("invalid user ID"), "failed to update profile", fiber.StatusInternalServerError)
	}

	var req dto.UpdateProfileRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
	}

	err := h.userService.UpdateUser(c.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			return response.HandleError(c, err, "", fiber.StatusConflict)
		}
		log.Printf("internal error: failed to update user: %v", err)
		return response.HandleError(c, err, "failed to update user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "update profile success", nil, fiber.StatusOK)
}
