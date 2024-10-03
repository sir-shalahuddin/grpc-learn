package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/pkg/response"
	pb "github.com/sir-shalahuddin/grpc-learn/bookservice/proto/authservice"
)

type AuthService interface {
	GetUserByID(ctx context.Context, userID string) (*pb.User, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

// AuthMiddleware struct holds the user service and JWT config
type authMiddleware struct {
	authService AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(authService AuthService) *authMiddleware {
	return &authMiddleware{authService: authService}
}

// AuthMiddleware provides JWT validation and role-based access control
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
		userIDStr, err := h.authService.ValidateToken(context.Background(), tokenString)
		if err != nil {
			return response.HandleError(c, nil, "Invalid or expired JWT", fiber.StatusUnauthorized)
		}

		// Retrieve the user
		user, err := h.authService.GetUserByID(context.Background(), userIDStr)
		if err != nil {
			return response.HandleError(c, nil, "Failed to fetch user detail", fiber.StatusInternalServerError)
		}

		// Check if the user has the required role
		if !hasAccess(user.Role, allowedRoles) {
			return response.HandleError(c, nil, "Access forbidden: insufficient permissions", fiber.StatusForbidden)
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return response.HandleError(c, nil, "Failed to parse userID", fiber.StatusInternalServerError)

		}

		c.Locals("id", userID) // Pass the user to the next handler
		return c.Next()
	}
}

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
