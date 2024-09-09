package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

type AuthService interface {
	Register(ctx context.Context, email, password, name string) error
	Login(ctx context.Context, email, password string) (map[string]string, error) // Returns JWT token
	RefreshToken(ctx context.Context, token string) (string, error)               // Returns new JWT token
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ValidateToken(tokenStr string) (uuid.UUID, error)
	HasAccess(userRole string, allowedRoles []string) bool
}

type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	err := h.authService.Register(context.Background(), req.Email, req.Password, req.Name)
	if err != nil {
		// log.Printf(err.Error())
		return response.HandleError(c, err, "failed to register user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "register success", nil, fiber.StatusOK)
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	token, err := h.authService.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		return response.HandleError(c, err, "failed to login", fiber.StatusUnauthorized)
	}

	return response.HandleSuccess(c, "login successfull", token, fiber.StatusOK)
}

func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	newToken, err := h.authService.RefreshToken(context.Background(), req.Token)
	if err != nil {
		log.Println(err.Error())
		return response.HandleError(c, err, "failed to refresh token", fiber.StatusUnauthorized)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"message":      "login successful",
		"access_token": newToken,
	})
}
