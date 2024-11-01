package dto

import (
	"github.com/Lontor/todo-api/internal/model"
)

// PUT /users/{userID}
type UpdateUserRequest struct {
	Email    string         `json:"email,omitempty" validate:"omitempty,email"`
	Password string         `json:"password,omitempty" validate:"omitempty,min=8"`
	Role     model.UserType `json:"role,omitempty" validate:"omitempty,oneof=regular admin"`
}
