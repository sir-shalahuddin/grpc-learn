package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

type AuthService interface {
	Register(ctx context.Context, email, password, name string) error
	Login(ctx context.Context, email, password string) (map[string]string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}

type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *authHandler {
	return &authHandler{authService: authService}
}

// Register handles user registration.
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
		return response.HandleError(c, err, "failed to register user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "registration successful", nil, fiber.StatusOK)
}

// Login handles user login and returns JWT tokens.
func (h *authHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	tokens, err := h.authService.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		return response.HandleError(c, err, "failed to login", fiber.StatusUnauthorized)
	}

	return response.HandleSuccess(c, "login successful", tokens, fiber.StatusOK)
}

// RefreshToken handles token refresh and returns a new JWT access token.
func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	newToken, err := h.authService.RefreshToken(context.Background(), req.Token)
	if err != nil {
		return response.HandleError(c, err, "failed to refresh token", fiber.StatusUnauthorized)
	}

	return response.HandleSuccess(c, "token refreshed successfully", map[string]string{"access_token": newToken}, fiber.StatusOK)
}
