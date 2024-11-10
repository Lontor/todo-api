package service

import (
	"context"
	"net/http"
	"time"

	"github.com/Lontor/todo-api/pkg/custom_errors"
	"github.com/Lontor/todo-api/pkg/utils"
	"github.com/go-playground/validator/v10"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/internal/repository"
	"github.com/google/uuid"
)

type taskService struct {
	r repository.TaskRepository
	v *validator.Validate
}

func NewTaskService(r repository.TaskRepository) TaskService {
	return &taskService{r, validator.New()}
}

func (s *taskService) CreateTask(ctx context.Context, data dto.CreateTaskRequest) error {
	if err := s.v.Struct(data); err != nil {
		return custom_errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}

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

	return utils.RepositoryErrorToHTTPError(s.r.Create(ctx, task))
}

func (s *taskService) GetTasksByUser(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error) {
	tokenUserID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	if role != model.UserTypeAdmin {
		if tokenUserID != userID {
			return nil, custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	if status == "" {
		tasks, err := s.r.GetByUserID(ctx, userID)
		return tasks, utils.RepositoryErrorToHTTPError(err)
	}

	return s.r.GetByUserIDAndStatus(ctx, userID, status)
}

func (s *taskService) GetTaskByID(ctx context.Context, id uuid.UUID) (model.Task, error) {
	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	task, err := s.r.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, utils.RepositoryErrorToHTTPError(err)
	}

	if role != model.UserTypeAdmin {
		if userID != task.UserID {
			return model.Task{}, custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, data dto.UpdateTaskRequest) error {
	if err := s.v.Struct(data); err != nil {
		return custom_errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	task, err := s.r.GetByID(ctx, data.TaskID)

	if err != nil {
		return utils.RepositoryErrorToHTTPError(err)
	}

	if role != model.UserTypeAdmin {
		if userID != task.UserID {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	if data.Description == "" && data.Status == "" {
		return custom_errors.NewHTTPError(http.StatusBadRequest, "no fields to update")
	}

	task.Description = data.Description
	task.Status = data.Status
	task.UpdatedAt = time.Now()

	return utils.RepositoryErrorToHTTPError(s.r.Update(ctx, task))
}

func (s *taskService) DeleteTask(ctx context.Context, id uuid.UUID, user uuid.UUID) error {
	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	task, err := s.r.GetByID(ctx, id)
	if err != nil {
		return utils.RepositoryErrorToHTTPError(err)
	}

	if task.UserID != user {
		return custom_errors.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if role != model.UserTypeAdmin {
		if userID != task.UserID {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return utils.RepositoryErrorToHTTPError(s.r.Delete(ctx, id))
}
