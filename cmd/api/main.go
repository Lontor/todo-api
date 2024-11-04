package main

import (
	"log"
	"net/http"

	"github.com/Lontor/todo-api/internal/handler"
	md "github.com/Lontor/todo-api/internal/middleware"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/internal/repository"
	"github.com/Lontor/todo-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open("host=localhost user=root dbname=todo_db password=34123413 port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(model.User{}, model.Task{})
	if err != nil {
		log.Fatalf("failed to migrate")
	}

	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)

	userService := service.NewUserService(userRepository)
	taskService := service.NewTaskService(taskRepository)

	handlers := handler.NewAPIHandler(userService, taskService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(md.JWTAuth("Y80zXN/dnoc14mdIpchh4ZXOdDAZfWulff5jEcCWHEc="))

	router.Get("/health", handlers.HealthCheck)

	router.Post("/auth/register", handlers.Register)
	router.Post("/auth/login", handlers.Login)

	router.Get("/users", handlers.GetUsers)
	router.Get("/users/{userID}", handlers.GetUser)
	router.Put("/users/{userID}", handlers.UpdateUser)
	router.Delete("/users/{userID}", handlers.DeleteUser)

	router.Get("/users/{userID}/tasks", handlers.GetUserTasks)
	router.Post("/users/{userID}/tasks", handlers.CreateUserTask)
	router.Get("/users/{userID}/tasks/{taskID}", handlers.GetUserTask)
	router.Patch("/users/{userID}/tasks/{taskID}", handlers.UpdateUserTask)
	router.Delete("/users/{userID}/tasks/{taskID}", handlers.DeleteUserTask)

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", router)
}
