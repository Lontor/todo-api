package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/internal/service"
	"github.com/Lontor/todo-api/pkg/custom_errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type apiHandler struct {
	userService service.UserService
	taskService service.TaskService
}

func NewAPIHandler(userService service.UserService, taskService service.TaskService) APIHandler {
	return &apiHandler{
		userService: userService,
		taskService: taskService,
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 500 {object} custom_errors.HTTPError
// @Router /health [get]
func (*apiHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{Status: "alive"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

// Login godoc
// @Summary User login
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /auth/login [post]
func (h *apiHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	response, err := h.userService.AuthenticateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Register godoc
// @Summary User registration
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param register body dto.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /auth/register [post]
func (h *apiHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{"message": "User registered successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users [get]
func (h *apiHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	response, err := h.userService.GetUsers(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.User
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID} [get]
func (h *apiHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	response, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param update body dto.UpdateUserRequest true "Update User Request"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID} [put]
func (h *apiHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	req.UserID = userID

	err = h.userService.UpdateUser(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "User updated successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID} [delete]
func (h *apiHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "User deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetUserTasks godoc
// @Summary Get tasks for a user
// @Description Get tasks for a user by their ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param filter query string false "Task status filter"
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.Task
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID}/tasks [get]
func (h *apiHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	status := model.TaskStatus(r.URL.Query().Get("filter"))

	response, err := h.taskService.GetTasksByUser(r.Context(), userID, status)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateUserTask godoc
// @Summary Create a new task for a user
// @Description Create a new task for a user by their ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param create body dto.CreateTaskRequest true "Create Task Request"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} map[string]string
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID}/tasks [post]
func (h *apiHandler) CreateUserTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	req.UserID = userID

	err = h.taskService.CreateTask(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{"message": "Task created successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetUserTask godoc
// @Summary Get a task for a user
// @Description Get a task for a user by their ID and task ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param taskID path string true "Task ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.Task
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID}/tasks/{taskID} [get]
func (h *apiHandler) GetUserTask(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		http.Error(w, "Invalid taskID format", http.StatusBadRequest)
		return
	}

	response, err := h.taskService.GetTaskByID(r.Context(), taskID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateUserTask godoc
// @Summary Update a task for a user
// @Description Update a task for a user by their ID and task ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param taskID path string true "Task ID"
// @Param Authorization header string true "Bearer token"
// @Param update body dto.UpdateTaskRequest true "Update Task Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID}/tasks/{taskID} [patch]
func (h *apiHandler) UpdateUserTask(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateTaskRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		http.Error(w, "Invalid taskID format", http.StatusBadRequest)
		return
	}

	req.UserID = userID
	req.TaskID = taskID

	err = h.taskService.UpdateTask(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Task updated successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteUserTask godoc
// @Summary Delete a task for a user
// @Description Delete a task for a user by their ID and task ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param taskID path string true "Task ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} custom_errors.HTTPError
// @Failure 500 {object} custom_errors.HTTPError
// @Router /users/{userID}/tasks/{taskID} [delete]
func (h *apiHandler) DeleteUserTask(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		http.Error(w, "Invalid taskID format", http.StatusBadRequest)
		return
	}

	err = h.taskService.DeleteTask(r.Context(), taskID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Task deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return err
	}
	return nil
}

func handleServiceError(w http.ResponseWriter, err error) {
	if httpErr, ok := err.(*custom_errors.HTTPError); ok {
		http.Error(w, httpErr.Message, httpErr.Code)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
