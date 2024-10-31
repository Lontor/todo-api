package repository

import (
	"context"
	"fmt"

	"github.com/Lontor/todo-api/internal/model"
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
	return r.db.Create(&task).Error
}

func (r *taskRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where(&model.Task{UserID: userID}).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where(&model.Task{UserID: userID, Status: status}).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Task, error) {
	var task model.Task
	err := r.db.Where(&model.Task{ID: id}).First(&task).Error
	return task, err
}

func (r *taskRepository) Update(ctx context.Context, id uuid.UUID, description string, status model.TaskStatus) error {
	updates := map[string]interface{}{}

	if description != "" {
		updates["description"] = description
	}
	if status != "" {
		updates["status"] = status
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result := r.db.Model(&model.Task{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no task found with id %s", id)
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(&model.Task{}, id).Error
}
