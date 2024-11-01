package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateTask_Success(t *testing.T) {
	clearDB()
	user := getTestUser()
	task := getTestTask()
	task.UserID = user.ID

	err := db.Create(&user).Error
	require.NoError(t, err)

	err = taskRepo.Create(context.Background(), task)
	assert.NoError(t, err)

	var createdTask model.Task
	err = db.First(&createdTask, task.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, task.Description, createdTask.Description)
	assert.Equal(t, task.UserID, createdTask.UserID)
	assert.Equal(t, task.Status, createdTask.Status)
	assert.Equal(t, task.UpdatedAt.Equal(createdTask.UpdatedAt), true)
	assert.Equal(t, task.CreatedAt.Equal(createdTask.CreatedAt), true)
}

func TestCreateTask_UserIDInvalid(t *testing.T) {
	clearDB()
	task := getTestTask()
	task.UserID = uuid.New()

	err := taskRepo.Create(context.Background(), task)
	assert.Error(t, err)
}

func TestCreateTask_WithEmptyRequiredFields(t *testing.T) {
	clearDB()
	user := getTestUser()
	task := getTestTask()
	task.UserID = user.ID

	err := db.Create(&user).Error
	require.NoError(t, err)

	errTask := task
	errTask.Description = ""
	err = taskRepo.Create(context.Background(), errTask)
	assert.Error(t, err)

	errTask = task
	errTask.Status = ""
	err = taskRepo.Create(context.Background(), errTask)
	assert.Error(t, err)
}

func TestGetTaskByUserIDAndStatus_UserNotFound(t *testing.T) {
	clearDB()

	nonExistentUserID := uuid.New()
	status := model.StatusInProgress

	receivedTasks, err := taskRepo.GetByUserIDAndStatus(context.Background(), nonExistentUserID, status)
	assert.Error(t, err)
	assert.Empty(t, receivedTasks)
}

func TestGetTaskByUserID_Success(t *testing.T) {
	clearDB()

	user := getTestUser()
	err := db.Create(&user).Error
	require.NoError(t, err)

	tasks := make([]model.Task, 3)

	for i := range tasks {
		task := getTestTask()
		task.UserID = user.ID
		tasks[i] = task

		err = db.Create(&tasks[i]).Error
		require.NoError(t, err)
	}

	receivedTasks, err := taskRepo.GetByUserID(context.Background(), user.ID)
	assert.NoError(t, err)
	require.Equal(t, len(tasks), len(receivedTasks))

	for i, task := range receivedTasks {
		assert.Equal(t, tasks[i].ID, task.ID)
		assert.Equal(t, tasks[i].UserID, task.UserID)
		assert.Equal(t, tasks[i].Description, task.Description)
		assert.Equal(t, tasks[i].Status, task.Status)
		assert.Equal(t, tasks[i].UpdatedAt.Equal(task.UpdatedAt), true)
		assert.Equal(t, tasks[i].CreatedAt.Equal(task.CreatedAt), true)
	}
}

func TestGetTaskByID_NotFound(t *testing.T) {
	clearDB()

	nonExistentTaskID := uuid.New()

	receivedTask, err := taskRepo.GetByID(context.Background(), nonExistentTaskID)
	assert.Error(t, err)
	assert.Empty(t, receivedTask)
}

func TestGetTaskByUserID_NotFound(t *testing.T) {
	clearDB()

	nonExistentUserID := uuid.New()

	receivedTasks, err := taskRepo.GetByUserID(context.Background(), nonExistentUserID)
	assert.Error(t, err)
	assert.Empty(t, receivedTasks)
}

func TestGetTaskByUserIDAndStatus_Success(t *testing.T) {
	clearDB()

	user := getTestUser()
	err := db.Create(&user).Error
	require.NoError(t, err)

	tasks := make([]model.Task, 3)
	status := model.StatusInProgress

	for i := range tasks {
		task := getTestTask()
		task.UserID = user.ID
		task.Status = status
		tasks[i] = task

		err = db.Create(&tasks[i]).Error
		require.NoError(t, err)
	}

	otherTask := getTestTask()
	otherTask.UserID = user.ID
	otherTask.Status = model.StatusDone
	err = db.Create(&otherTask).Error
	require.NoError(t, err)

	receivedTasks, err := taskRepo.GetByUserIDAndStatus(context.Background(), user.ID, status)
	assert.NoError(t, err)
	require.Equal(t, len(tasks), len(receivedTasks))

	for i, task := range receivedTasks {
		assert.Equal(t, tasks[i].ID, task.ID)
		assert.Equal(t, tasks[i].UserID, task.UserID)
		assert.Equal(t, tasks[i].Description, task.Description)
		assert.Equal(t, tasks[i].Status, task.Status)
		assert.Equal(t, tasks[i].UpdatedAt.Equal(task.UpdatedAt), true)
		assert.Equal(t, tasks[i].CreatedAt.Equal(task.CreatedAt), true)
	}
}

func TestGetTaskByID_Success(t *testing.T) {
	clearDB()

	user := getTestUser()
	err := db.Create(&user).Error
	require.NoError(t, err)

	task := getTestTask()
	task.UserID = user.ID
	err = db.Create(&task).Error
	require.NoError(t, err)

	receivedTask, err := taskRepo.GetByID(context.Background(), task.ID)
	assert.NoError(t, err)
	assert.Equal(t, task.ID, receivedTask.ID)
	assert.Equal(t, task.UserID, receivedTask.UserID)
	assert.Equal(t, task.Description, receivedTask.Description)
	assert.Equal(t, task.Status, receivedTask.Status)
	assert.Equal(t, task.UpdatedAt.Equal(receivedTask.UpdatedAt), true)
	assert.Equal(t, task.CreatedAt.Equal(receivedTask.CreatedAt), true)
}

func TestUpdateTask_Success(t *testing.T) {
	clearDB()

	user := getTestUser()
	err := db.Create(&user).Error
	require.NoError(t, err)

	task := getTestTask()
	task.UserID = user.ID
	err = db.Create(&task).Error
	require.NoError(t, err)

	newDescription := "Updated task description"
	newStatus := model.StatusInProgress

	err = taskRepo.Update(context.Background(), task.ID, newDescription, newStatus)
	assert.NoError(t, err)

	var updatedTask model.Task
	err = db.First(&updatedTask, task.ID).Error
	require.NoError(t, err)
	assert.Equal(t, newDescription, updatedTask.Description)
	assert.Equal(t, newStatus, updatedTask.Status)
}

func TestUpdateTask_PartialUpdate(t *testing.T) {
	clearDB()

	user := getTestUser()
	err := db.Create(&user).Error
	require.NoError(t, err)

	task := getTestTask()
	task.UserID = user.ID
	err = db.Create(&task).Error
	require.NoError(t, err)

	newDescription := "Updated description only"

	err = taskRepo.Update(context.Background(), task.ID, newDescription, "")
	assert.NoError(t, err)

	var updatedTask model.Task
	err = db.First(&updatedTask, task.ID).Error
	require.NoError(t, err)
	assert.Equal(t, newDescription, updatedTask.Description)
	assert.Equal(t, task.Status, updatedTask.Status)

	newStatus := model.StatusDone

	err = taskRepo.Update(context.Background(), task.ID, "", newStatus)
	assert.NoError(t, err)

	err = db.First(&updatedTask, task.ID).Error
	require.NoError(t, err)
	assert.Equal(t, newStatus, updatedTask.Status)
	assert.Equal(t, newDescription, updatedTask.Description)
}

func TestUpdateTask_NotFound(t *testing.T) {
	clearDB()

	nonExistentID := uuid.New()

	err := taskRepo.Update(context.Background(), nonExistentID, "Some description", model.StatusDone)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no task found with id")
}

func TestUpdateTask_NoFieldsToUpdate(t *testing.T) {
	clearDB()

	user := getTestUser()
	err := db.Create(&user).Error
	require.NoError(t, err)

	task := getTestTask()
	task.UserID = user.ID
	err = db.Create(&task).Error
	require.NoError(t, err)

	err = taskRepo.Update(context.Background(), task.ID, "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fields to update")
}

func TestDeleteTask_Success(t *testing.T) {
	clearDB()

	user := getTestUser()
	err := db.Create(&user).Error
	require.NoError(t, err)

	task := getTestTask()
	task.UserID = user.ID
	err = db.Create(&task).Error
	require.NoError(t, err)

	err = taskRepo.Delete(context.Background(), task.ID)
	assert.NoError(t, err)

	var deletedTask model.Task
	err = db.First(&deletedTask, task.ID).Error
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestDeleteTask_NotFound(t *testing.T) {
	clearDB()

	nonExistentID := uuid.New()

	err := taskRepo.Delete(context.Background(), nonExistentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no task found with id")
}
