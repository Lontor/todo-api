package service

import (
	"context"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

type TaskService interface {
	CreateTask(ctx context.Context, data dto.CreateTaskRequest) error
	GetTasksByUser(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (model.Task, error)
	UpdateTask(ctx context.Context, data dto.UpdateTaskRequest) error
	DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
