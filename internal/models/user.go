package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"                    json:"id"`
	Username  string    `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Counter   int64     `gorm:"default:0;not null"            json:"counter"`
	CreatedAt time.Time `                                     json:"created_at"`
	UpdatedAt time.Time `                                     json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Request/Response DTOs
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
}

type CreateUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Counter  int64  `json:"counter"`
}

type GetCountResponse struct {
	UserID  uint   `json:"user_id"`
	Counter int64  `json:"counter"`
	Source  string `json:"source"` // "utils" or "database"
}

type IncrementResponse struct {
	UserID     uint  `json:"user_id"`
	NewCounter int64 `json:"new_counter"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
