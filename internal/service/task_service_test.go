package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/pkg/custom_errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error) {
	args := m.Called(ctx, userID, status)
	if tasks, ok := args.Get(0).([]model.Task); ok {
		return tasks, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	userID := uuid.New()
	data := dto.CreateTaskRequest{
		UserID:      userID,
		Description: "Test task description",
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("Create", ctx, mock.Anything).Return(nil)

		err := taskService.CreateTask(ctx, data)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", uuid.New())
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		err := taskService.CreateTask(ctx, data)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		err := taskService.CreateTask(ctx, data)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)
		mockRepo.ExpectedCalls = nil
	})
}

func TestGetTasksByUser(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	userID := uuid.New()
	task := model.Task{
		ID:          uuid.New(),
		UserID:      userID,
		Description: "Sample Task",
		Status:      model.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks := []model.Task{task}

	t.Run("success - without status filter", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByUserID", ctx, userID).Return(tasks, nil)

		result, err := taskService.GetTasksByUser(ctx, userID, "")

		assert.NoError(t, err)
		assert.Equal(t, tasks, result)
		mockRepo.AssertCalled(t, "GetByUserID", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("success - with status filter", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByUserIDAndStatus", ctx, userID, model.StatusTodo).Return(tasks, nil)

		result, err := taskService.GetTasksByUser(ctx, userID, model.StatusTodo)

		assert.NoError(t, err)
		assert.Equal(t, tasks, result)
		mockRepo.AssertCalled(t, "GetByUserIDAndStatus", ctx, userID, model.StatusTodo)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", uuid.New())
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		result, err := taskService.GetTasksByUser(ctx, userID, model.StatusTodo)

		assert.Nil(t, result)
		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByUserIDAndStatus", ctx, userID, model.StatusTodo).Return(nil, fmt.Errorf("user with id %s not found", userID))

		result, err := taskService.GetTasksByUser(ctx, userID, model.StatusTodo)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "user with id "+userID.String()+" not found")
		mockRepo.AssertCalled(t, "GetByUserIDAndStatus", ctx, userID, model.StatusTodo)
		mockRepo.ExpectedCalls = nil
	})
}

func TestGetTaskByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	userID := uuid.New()
	task := model.Task{
		ID:          uuid.New(),
		UserID:      userID,
		Description: "Sample Task",
		Status:      model.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, task.ID).Return(task, nil)

		result, err := taskService.GetTaskByID(ctx, task.ID, task.UserID)

		assert.NoError(t, err)
		assert.Equal(t, task, result)
		mockRepo.AssertCalled(t, "GetByID", ctx, task.ID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", uuid.New())
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, task.ID).Return(task, nil)

		result, err := taskService.GetTaskByID(ctx, task.ID, task.UserID)

		assert.Equal(t, model.Task{}, result)
		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("task not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, task.ID).Return(model.Task{}, fmt.Errorf("record not found"))

		result, err := taskService.GetTaskByID(ctx, task.ID, task.UserID)

		assert.Equal(t, model.Task{}, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "record not found")
		mockRepo.AssertCalled(t, "GetByID", ctx, task.ID)
		mockRepo.ExpectedCalls = nil
	})
}

func TestUpdateTask(t *testing.T) {
	userID := uuid.New()
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskService := NewTaskService(mockRepo)

		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		data := dto.UpdateTaskRequest{
			TaskID:      taskID,
			UserID:      userID,
			Description: "Updated task description",
			Status:      model.StatusInProgress,
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(nil)

		err := taskService.UpdateTask(ctx, data)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Update", ctx, mock.Anything)
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskService := NewTaskService(mockRepo)

		ctx := context.WithValue(context.Background(), "userID", uuid.New())
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		data := dto.UpdateTaskRequest{
			TaskID:      taskID,
			UserID:      userID,
			Description: "Updated task description",
			Status:      model.StatusInProgress,
		}

		err := taskService.UpdateTask(ctx, data)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("no fields to update", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskService := NewTaskService(mockRepo)

		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		data := dto.UpdateTaskRequest{
			TaskID:      taskID,
			UserID:      userID,
			Description: "",
			Status:      "",
		}

		err := taskService.UpdateTask(ctx, data)

		assert.Error(t, err)
		assert.EqualError(t, err, "no fields to update")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskService := NewTaskService(mockRepo)

		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		data := dto.UpdateTaskRequest{
			TaskID:      taskID,
			UserID:      userID,
			Description: "Updated task description",
			Status:      model.StatusInProgress,
		}

		mockRepo.On("Update", ctx, mock.Anything).Return(fmt.Errorf("db error"))

		err := taskService.UpdateTask(ctx, data)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertCalled(t, "Update", ctx, mock.Anything)
	})
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	userID := uuid.New()
	task := model.Task{
		ID:          uuid.New(),
		UserID:      userID,
		Description: "Sample Task",
		Status:      model.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("success - admin", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", uuid.New())
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("GetByID", ctx, task.ID).Return(task, nil)
		mockRepo.On("Delete", ctx, task.ID).Return(nil)

		err := taskService.DeleteTask(ctx, task.ID, task.UserID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", ctx, task.ID)
		mockRepo.AssertCalled(t, "Delete", ctx, task.ID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("success - owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, task.ID).Return(task, nil)
		mockRepo.On("Delete", ctx, task.ID).Return(nil)

		err := taskService.DeleteTask(ctx, task.ID, task.UserID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "GetByID", ctx, task.ID)
		mockRepo.AssertCalled(t, "Delete", ctx, task.ID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", uuid.New())
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, task.ID).Return(task, nil)

		err := taskService.DeleteTask(ctx, task.ID, task.UserID)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("task not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, task.ID).Return(model.Task{}, fmt.Errorf("no task found with id %s", task.ID))

		err := taskService.DeleteTask(ctx, task.ID, task.UserID)

		assert.Error(t, err)
		assert.EqualError(t, err, "no task found with id "+task.ID.String())
		mockRepo.AssertCalled(t, "GetByID", ctx, task.ID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("repository error on delete", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, task.ID).Return(task, nil)
		mockRepo.On("Delete", ctx, task.ID).Return(fmt.Errorf("db error"))

		err := taskService.DeleteTask(ctx, task.ID, task.UserID)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertCalled(t, "GetByID", ctx, task.ID)
		mockRepo.AssertCalled(t, "Delete", ctx, task.ID)
		mockRepo.ExpectedCalls = nil
	})
}
