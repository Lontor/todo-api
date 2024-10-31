package repository

import (
	"context"
	"fmt"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/go-playground/validator/v10"
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
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return err
	}
	return r.db.Create(&user).Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{Email: email}).First(&user).Error
	return user, err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{ID: id}).First(&user).Error
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user model.User) error {
	updates := map[string]interface{}{}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.PasswordHash != "" {
		updates["password_hash"] = user.PasswordHash
	}
	if user.AccountType != "" {
		updates["account_type"] = user.AccountType
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result := r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %s", user.ID)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}
	return nil
}
