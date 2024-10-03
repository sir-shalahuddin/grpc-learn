package dto

import (
	"time"
)

type BorrowBookRequest struct {
	DueDate *time.Time `json:"due_date"`
}
