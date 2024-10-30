package repository

import (
	"context"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Task, error)
	UpdateDescription(ctx context.Context, id uuid.UUID, description string) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.TaskStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
