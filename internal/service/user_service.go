package service

import (
	"context"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, data dto.RegisterRequest) error
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
	UpdateUser(ctx context.Context, data dto.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	AuthenticateUser(ctx context.Context, email, password string) (dto.AuthResponse, error)
}
