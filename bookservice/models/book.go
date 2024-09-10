package models

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID  `json:"id"`             // Unique identifier for the book
	Title         string     `json:"title"`          // Title of the book
	Author        string     `json:"author"`         // Author of the book
	ISBN          string     `json:"isbn"`           // ISBN number of the book
	PublishedDate *time.Time `json:"published_date"` // Date when the book was published
	CategoryID    uuid.UUID  `json:"category_id"`    // ID of the category the book belongs to
	CategoryName  string     `json:"category_name"`  // Name of the category the book belongs to
	Stock         int        `json:"stock"`          // Number of copies available
	AddedBy       uuid.UUID  `json:"added_by"`       // User ID of the librarian who added the book
	CreatedAt     *time.Time `json:"created_at"`     // Timestamp when the book was created
	UpdatedAt     *time.Time `json:"updated_at"`     // Timestamp when the book was last updated
	Version       int        `json:"version"`
}

// BorrowingRecord represents a record of a book borrowed by a user.
type BorrowingRecord struct {
	ID         uuid.UUID  `json:"id"`          // Unique identifier for the borrowing record
	BookID     uuid.UUID  `json:"book_id"`     // ID of the borrowed book
	UserID     uuid.UUID  `json:"user_id"`     // ID of the user who borrowed the book
	BorrowedAt time.Time  `json:"borrowed_at"` // Timestamp when the book was borrowed
	DueDate    *time.Time `json:"due_date"`    // Due date for returning the borrowed book
	ReturnedAt time.Time  `json:"returned_at"` // Timestamp when the book was returned (if applicable)
}
