package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/response"
)

// AuthService interface defines methods for authentication and authorization
type AuthMiddlewareService interface {
	ValidateToken(tokenStr string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// authMiddleware struct holds the AuthService instance
type authMiddleware struct {
	service AuthMiddlewareService
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(service AuthMiddlewareService) *authMiddleware {
	return &authMiddleware{service: service}
}

// Protected provides JWT validation and role-based access control
func (h *authMiddleware) Protected(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.HandleError(c, nil, "Missing or malformed JWT", fiber.StatusUnauthorized)
		}

		// Extract token from Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return response.HandleError(c, nil, "Missing or malformed JWT", fiber.StatusUnauthorized)
		}

		// Parse and validate the JWT token
		userID, err := h.service.ValidateToken(tokenString)
		if err != nil {
			return response.HandleError(c, err, "Invalid or expired JWT", fiber.StatusUnauthorized)
		}

		// Retrieve the user
		user, err := h.service.GetUserByID(context.Background(), userID)
		if err != nil {
			return response.HandleError(c, err, "Failed to fetch user details", fiber.StatusInternalServerError)
		}

		// Check if the user has the required role
		if !hasAccess(user.Role, allowedRoles) {
			return response.HandleError(c, nil, "Access forbidden: insufficient permissions", fiber.StatusForbidden)
		}

		// Pass the user ID to the next handler
		c.Locals("id", userID)
		return c.Next()
	}
}

// hasAccess checks if the user's role is in the allowedRoles list
func hasAccess(userRole string, allowedRoles []string) bool {
	if len(allowedRoles) == 0 {
		allowedRoles = append(allowedRoles, "user")
	}
	if userRole == "super admin" {
		return true
	}

	for _, role := range allowedRoles {
		if role == userRole {
			return true
		}
	}
	return false
}
