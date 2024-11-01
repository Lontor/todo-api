package dto

import "github.com/Lontor/todo-api/internal/model"

// POST /users/{userID}/tasks
type CreateTaskRequest struct {
	Description string           `json:"description" validate:"required,min=10,max=200"`
	Status      model.TaskStatus `json:"status,omitempty" validate:"omitempty,oneof='to do' 'in progress' 'done'"`
}

// PATCH /users/{userID}/tasks/{taskID}
type UpdateTaskRequest struct {
	Description *string           `json:"description,omitempty" validate:"omitempty,min=10,max=200"`
	Status      *model.TaskStatus `json:"status,omitempty" validate:"omitempty,oneof='to do' 'in progress' 'done'"`
}
