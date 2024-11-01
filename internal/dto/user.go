package dto

import (
	"time"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

// GET /users/{userID}
type UserResponse struct {
	ID          uuid.UUID      `json:"id"`
	Email       string         `json:"email"`
	AccountType model.UserType `json:"accountType"`
	CreatedAt   time.Time      `json:"createdAt"`
}

// PUT /users/{userID}
type UpdateUserRequest struct {
	Email    string         `json:"email,omitempty" validate:"omitempty,email"`
	Password string         `json:"password,omitempty" validate:"omitempty,min=8"`
	Role     model.UserType `json:"role,omitempty" validate:"omitempty,oneof=regular admin"`
}
