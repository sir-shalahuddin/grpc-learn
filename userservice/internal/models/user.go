package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID
	Name      string    
	Email     string   
	Password  string    
	CreatedAt time.Time 
	UpdatedAt time.Time 
	Role      string    
}
