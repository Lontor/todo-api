package dto

import (
	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

// POST /users/{userID}/tasks
type CreateTaskRequest struct {
	Description string    `json:"description" validate:"required,min=3,max=200"`
	UserID      uuid.UUID `json:"-" validate:"required"`
}

// PATCH /tasks/{taskID}
type UpdateTaskRequest struct {
	Description string           `json:"description,omitempty" validate:"omitempty,min=3,max=200"`
	Status      model.TaskStatus `json:"status,omitempty" validate:"omitempty,oneof='to do' 'in progress' 'done'"`
	TaskID      uuid.UUID        `json:"-" validate:"required"`
}
