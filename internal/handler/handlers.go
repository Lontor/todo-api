package handler

import (
	"net/http"
)

type APIHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request) // GET /health

	Login(w http.ResponseWriter, r *http.Request)    // POST /auth/login
	Register(w http.ResponseWriter, r *http.Request) // POST /auth/register

	GetUsers(w http.ResponseWriter, r *http.Request)   // GET /users
	GetUser(w http.ResponseWriter, r *http.Request)    // GET /users/{userID}
	UpdateUser(w http.ResponseWriter, r *http.Request) // PUT /users/{userID}
	DeleteUser(w http.ResponseWriter, r *http.Request) // DELETE /users/{userID}

	GetUserTasks(w http.ResponseWriter, r *http.Request)   // GET /users/{userID}/tasks
	CreateUserTask(w http.ResponseWriter, r *http.Request) // POST /users/{userID}/tasks
	GetUserTask(w http.ResponseWriter, r *http.Request)    // GET /users/{userID}/tasks/{taskID}
	UpdateUserTask(w http.ResponseWriter, r *http.Request) // PATCH /users/{userID}/tasks/{taskID}
	DeleteUserTask(w http.ResponseWriter, r *http.Request) // DELETE /users/{userID}/tasks/{taskID}
}
