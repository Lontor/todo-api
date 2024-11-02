package service

import (
	"context"
	"net/http"
	"time"

	"github.com/Lontor/todo-api/pkg/custom_errors"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/internal/repository"
	"github.com/google/uuid"
)

type taskService struct {
	r repository.TaskRepository
}

func NewTaskService(r *repository.TaskRepository) TaskService {
	return &taskService{*r}
}

func (s *taskService) CreateTask(ctx context.Context, data dto.CreateTaskRequest) error {
	role := ctx.Value("role").(model.UserType)

	if role != model.UserTypeAdmin {
		if ctx.Value("userID").(uuid.UUID) != data.UserID {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	now := time.Now()

	task := model.Task{
		ID:          uuid.New(),
		UserID:      data.UserID,
		Description: data.Description,
		Status:      model.StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.r.Create(ctx, task)
}

func (s *taskService) GetTasksByUser(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error) {
	tokenUserID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	if role != model.UserTypeAdmin {
		if tokenUserID != userID {
			return nil, custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return s.r.GetByUserID(ctx, userID)
}

func (s *taskService) GetTaskByID(ctx context.Context, id uuid.UUID) (model.Task, error) {
	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	task, err := s.r.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, err
	}

	if role != model.UserTypeAdmin {
		if userID != task.UserID {
			return model.Task{}, custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, data dto.UpdateTaskRequest) error {
	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	if role != model.UserTypeAdmin {
		if userID != data.UserID {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return s.r.Update(ctx, model.Task{
		ID:          data.TaskID,
		Description: *data.Description,
		Status:      *data.Status,
		UpdatedAt:   time.Now(),
	})
}

func (s *taskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	task, err := s.r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if role != model.UserTypeAdmin {
		if userID != task.UserID {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return s.r.Delete(ctx, id)
}
