package repository

import (
	"context"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task model.Task) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error)
	GetByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Task, error)
	Update(ctx context.Context, id uuid.UUID, description string, status model.TaskStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
