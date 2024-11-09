package repository

import (
	"context"
	"fmt"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/pkg/custom_errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user model.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to create user: %v", err))
	}
	return nil
}

func (r *userRepository) Get(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to fetch users: %v", err))
	}
	return users, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{Email: email}).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, custom_errors.NewRepositoryError(
				custom_errors.CodeNotFound,
				fmt.Sprintf("user with email %s not found", email))
		}
		return user, custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to fetch user with email %s: %v", email, err))
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{ID: id}).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, custom_errors.NewRepositoryError(
				custom_errors.CodeNotFound,
				fmt.Sprintf("user with ID %s not found", id))
		}
		return user, custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to fetch user with ID %s: %v", id, err))
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user model.User) error {
	result := r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to update user with ID %s: %v", user.ID, result.Error))
	}
	if result.RowsAffected == 0 {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeNotFound,
			fmt.Sprintf("no user found with id %s", user.ID))
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeDBError,
			fmt.Sprintf("failed to delete user with ID %s: %v", id, result.Error))
	}
	if result.RowsAffected == 0 {
		return custom_errors.NewRepositoryError(
			custom_errors.CodeNotFound,
			fmt.Sprintf("no user found with id %s", id))
	}
	return nil
}
