package repository

import (
	"context"
	"fmt"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/go-playground/validator/v10"
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
	validate := validator.New()
	err := validate.Struct(task)
	if err != nil {
		return err
	}

	var user model.User
	if err := r.db.First(&user, task.UserID).Error; err != nil {
		return fmt.Errorf("user with ID %s does not exist", task.UserID)
	}

	return r.db.Create(&task).Error
}

func (r *taskRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user with id %s not found", userID)
	}

	var tasks []model.Task
	err = r.db.Where(&model.Task{UserID: userID}).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status model.TaskStatus) ([]model.Task, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user with id %s not found", userID)
	}

	var tasks []model.Task
	err = r.db.Where(&model.Task{UserID: userID, Status: status}).Find(&tasks).Error
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
	result := r.db.Delete(&model.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no task found with id %s", id)
	}
	return nil
}
