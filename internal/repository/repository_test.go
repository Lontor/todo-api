package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var userRepo UserRepository
var taskRepo TaskRepository

func getTestUser() model.User {
	user := model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		AccountType:  model.UserTypeRegular,
		CreatedAt:    time.Now(),
	}
	return user
}

func getTestTask() model.Task {
	task := model.Task{
		ID:          uuid.New(),
		Description: fmt.Sprintf("Task %s", uuid.New()),
		Status:      model.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return task
}

func TestMain(m *testing.M) {
	var err error

	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err := db.AutoMigrate(&model.User{}, &model.Task{}); err != nil {
		panic("failed to migrate database")
	}

	userRepo = NewUserRepository(db)
	taskRepo = NewTaskRepository(db)

	m.Run()
}

func clearDB() {
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		panic("failed to clear database")
	}
	if err := db.Exec("DELETE FROM tasks").Error; err != nil {
		panic("failed to clear database")
	}
}
