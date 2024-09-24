package dto

import "github.com/google/uuid"

type CreateBookCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateBookCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UpdateBookCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}
