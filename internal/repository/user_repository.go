package repository

import (
	"context"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	Get(ctx context.Context) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
