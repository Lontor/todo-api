package dto

import (
	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

// POST /users/{userID}/tasks
type CreateTaskRequest struct {
	Description string           `json:"description" validate:"required,min=10,max=200"`
	Status      model.TaskStatus `json:"status,omitempty" validate:"omitempty,oneof='to do' 'in progress' 'done'"`
	UserID      uuid.UUID        `json:"userID" validate:"required"`
}

// PATCH /users/{userID}/tasks/{taskID}
type UpdateTaskRequest struct {
	Description string           `json:"description,omitempty" validate:"omitempty,min=10,max=200"`
	Status      model.TaskStatus `json:"status,omitempty" validate:"omitempty,oneof='to do' 'in progress' 'done'"`
	UserID      uuid.UUID        `json:"userID" validate:"required"`
	TaskID      uuid.UUID        `json:"taskID" validate:"required"`
}
