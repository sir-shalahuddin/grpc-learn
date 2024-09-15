package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetUser struct {
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      string    `json:"role"`
}

type UpdateUserRoles struct {
	Role string `json:"role" validate:"required"`
}
