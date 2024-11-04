package dto

import (
	"time"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

// POST /auth/login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// POST /auth/register
type RegisterRequest struct {
	Email    string         `json:"email" validate:"required,email"`
	Password string         `json:"password" validate:"required,min=8"`
	Role     model.UserType `json:"role,omitempty" validate:"omitempty,oneof=regular admin"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"userID"`
	ExpiresAt time.Time `json:"expiresAt"`
}
