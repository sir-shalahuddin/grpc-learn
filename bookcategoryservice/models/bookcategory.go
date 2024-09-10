package models

import "github.com/google/uuid"

type BookCategory struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
