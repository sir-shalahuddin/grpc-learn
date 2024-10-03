package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddBookRequest struct {
	Title         string     `json:"title" validate:"required"`
	Author        string     `json:"author" validate:"required"`
	ISBN          string     `json:"isbn"`
	PublishedDate *time.Time `json:"published_date"`
	CategoryID    uuid.UUID  `json:"category_id"`
	Stock         int        `json:"stock"`
}

type GetBookResponse struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	ISBN          string     `json:"isbn"`
	PublishedDate *time.Time `json:"published_date"`
	Category      string     `json:"category"`
	Stock         int        `json:"stock"`
}

type UpdateBookRequest struct {
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	ISBN          string     `json:"isbn"`
	PublishedDate *time.Time `json:"published_date"`
	CategoryID    uuid.UUID  `json:"category_id"`
	Stock         int        `json:"stock"`
}
