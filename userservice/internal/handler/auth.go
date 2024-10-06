package handler

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/service"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/auth"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error)
}

type authHandler struct {
	authService AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService AuthService) *authHandler {
	validate := validator.New()
	validate.RegisterValidation("password", ValidatePassword)

	return &authHandler{
		authService: authService,
		validate:    validate,
	}
}

// Register handles user registration.
// @Summary Register a new user
// @Description Registers a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body dto.RegisterRequest true "Register Request"
// @Success 200 {object} response.Response "Registration successful"
// @Failure 400 {object} response.ErrorMessage "Invalid request payload"
// @Failure 409 {object} response.ErrorMessage "Email already exists"
// @Failure 500 {object} response.ErrorMessage "Failed to register user"
// @Router /auth/register [post]
func (h *authHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			switch e.Field() {
			case "Password":
				return response.HandleError(c, err, "invalid password format", fiber.StatusBadRequest)
			case "Email":
				return response.HandleError(c, err, "invalid email format", fiber.StatusBadRequest)
			default:
				return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
			}
		}
	}

	err := h.authService.Register(context.Background(), req)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			return response.HandleError(c, err, "", fiber.StatusConflict)
		}
		log.Printf("internal error: failed to register user: %v", err)
		return response.HandleError(c, err, "failed to register user", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "registration successful", nil, fiber.StatusOK)
}

// Login handles user login and returns JWT tokens.
// @Summary Login user
// @Description Logs in a user and returns JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "Login Request"
// @Success 200 {object} response.Response "Login successful"
// @Failure 400 {object} response.ErrorMessage "Invalid request payload"
// @Failure 401 {object} response.ErrorMessage "Invalid credentials"
// @Failure 500 {object} response.ErrorMessage "Failed to login user"
// @Router /auth/login [post]
func (h *authHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			switch e.Field() {
			case "Password":
				return response.HandleError(c, err, "invalid password format", fiber.StatusBadRequest)
			case "Email":
				return response.HandleError(c, err, "invalid email format", fiber.StatusBadRequest)
			default:
				return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
			}
		}
	}

	res, err := h.authService.Login(context.Background(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return response.HandleError(c, err, "", fiber.StatusUnauthorized)
		}
		log.Printf("internal error: failed to login: %v", err)
		return response.HandleError(c, err, "failed to login", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "login successful", res, fiber.StatusOK)
}

// RefreshToken handles token refresh and returns a new JWT access token.
// @Summary Refresh JWT Token
// @Description Refreshes a JWT token and returns a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshTokenRequest body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} response.Response "Token refreshed successfully"
// @Failure 400 {object} response.ErrorMessage "Invalid request payload"
// @Failure 401 {object} response.ErrorMessage "Invalid token"
// @Failure 500 {object} response.ErrorMessage "Failed to refresh token"
// @Router /auth/refresh-token [post]
func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, err, "invalid request payload", fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return response.HandleError(c, err, "invalid payload", fiber.StatusBadRequest)
	}

	res, err := h.authService.RefreshToken(context.Background(), req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return response.HandleError(c, err, "", fiber.StatusUnauthorized)
		}
		log.Printf("internal error: failed to refresh token: %v", err)
		return response.HandleError(c, err, "failed to refresh token", fiber.StatusInternalServerError)
	}

	return response.HandleSuccess(c, "token refreshed successfully", res, fiber.StatusOK)
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasNumber && hasUppercase && hasLowercase && hasSpecialChar
}
