package repository

import (
	"context"
	"fmt"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/pkg/custom_errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) Create(ctx context.Context, task model.Task) error {
	if err := r.db.Create(&task).Error; err != nil {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to create task: %v", err))
	}
	return nil
}

func (r *taskRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where(&model.Task{UserID: userID}).Find(&tasks).Error
	if err != nil {
		return nil, custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to fetch tasks for user %s: %v", userID, err))
	}
	return tasks, nil
}

func (r *taskRepository) GetByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where(&model.Task{UserID: userID, Status: status}).Find(&tasks).Error
	if err != nil {
		return nil, custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to fetch tasks for user %s with status %s: %v", userID, status, err))
	}
	return tasks, nil
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Task, error) {
	var task model.Task
	err := r.db.Where(&model.Task{ID: id}).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return task, custom_errors.NewRepositoryError(
				custom_errors.CodeNotFound,
				fmt.Sprintf("task with ID %s not found", id))
		}
		return task, custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to fetch task with ID %s: %v", id, err))
	}
	return task, nil
}

func (r *taskRepository) Update(ctx context.Context, task model.Task) error {
	result := r.db.Model(&model.Task{}).Where("id = ?", task.ID).Updates(task)
	if result.Error != nil {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to update task with ID %s: %v", task.ID, result.Error))
	}
	if result.RowsAffected == 0 {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeNotFound,
			fmt.Sprintf("no task found with id %s", task.ID))
	}
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&model.Task{}, id)
	if result.Error != nil {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to delete task with ID %s: %v", id, result.Error))
	}
	if result.RowsAffected == 0 {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeNotFound,
			fmt.Sprintf("no task found with id %s", id))
	}
	return nil
}
