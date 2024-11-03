package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/pkg/custom_errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Get(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]model.User); ok {
		return users, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	data := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     model.UserTypeRegular,
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "role", model.UserTypeAdmin)

		mockRepo.On("Create", ctx, mock.Anything).Return(nil)

		err := userService.CreateUser(ctx, data)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not admin", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "role", model.UserTypeRegular)

		err := userService.CreateUser(ctx, data)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("permission denied - invalid role", func(t *testing.T) {
		data.Role = model.UserTypeAdmin
		ctx := context.WithValue(context.Background(), "role", model.UserTypeRegular)

		err := userService.CreateUser(ctx, data)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "role", model.UserTypeAdmin)

		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		err := userService.CreateUser(ctx, data)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("invalid data", func(t *testing.T) {
		data.Email = "invalid-email"
		ctx := context.WithValue(context.Background(), "role", model.UserTypeAdmin)

		err := userService.CreateUser(ctx, data)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		assert.Contains(t, err.Error(), "validation")
	})
}

func TestGetUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "role", model.UserTypeAdmin)
		expectedUsers := []model.User{
			{ID: uuid.New(), Email: "user1@example.com", AccountType: model.UserTypeRegular},
			{ID: uuid.New(), Email: "user2@example.com", AccountType: model.UserTypeAdmin},
		}

		mockRepo.On("Get", ctx).Return(expectedUsers, nil)

		users, err := userService.GetUsers(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockRepo.AssertCalled(t, "Get", ctx)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not admin", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "role", model.UserTypeRegular)

		users, err := userService.GetUsers(ctx)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
		assert.Nil(t, users)
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "role", model.UserTypeAdmin)

		mockRepo.On("Get", ctx).Return(nil, errors.New("db error"))

		users, err := userService.GetUsers(ctx)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		assert.Nil(t, users)
		mockRepo.AssertCalled(t, "Get", ctx)
		mockRepo.ExpectedCalls = nil
	})
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userID := uuid.New()
	otherUserID := uuid.New()
	user := model.User{
		ID:          userID,
		Email:       "test@example.com",
		AccountType: model.UserTypeRegular,
	}

	t.Run("success - admin", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("GetByID", ctx, userID).Return(user, nil)

		result, err := userService.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertCalled(t, "GetByID", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("success - owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("GetByID", ctx, userID).Return(user, nil)

		result, err := userService.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertCalled(t, "GetByID", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		_, err := userService.GetUserByID(ctx, otherUserID)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("GetByID", ctx, userID).Return(model.User{}, errors.New("db error"))

		_, err := userService.GetUserByID(ctx, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertCalled(t, "GetByID", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userID := uuid.New()
	data := dto.UpdateUserRequest{
		UserID:   userID,
		Email:    "newemail@example.com",
		Password: "newpassword",
		Role:     model.UserTypeRegular,
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("Update", ctx, mock.Anything).Return(nil)

		err := userService.UpdateUser(ctx, data)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Update", ctx, mock.Anything)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", uuid.New())
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		err := userService.UpdateUser(ctx, data)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("permission denied - invalid role", func(t *testing.T) {
		data.Role = model.UserTypeAdmin
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		err := userService.UpdateUser(ctx, data)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("no fields to update", func(t *testing.T) {
		errData := data
		errData.Email = ""
		errData.Password = ""
		errData.Role = ""
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		err := userService.UpdateUser(ctx, errData)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		assert.EqualError(t, err, "no fields to update")
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("db error"))

		err := userService.UpdateUser(ctx, data)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertCalled(t, "Update", ctx, mock.Anything)
		mockRepo.ExpectedCalls = nil
	})
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("success - admin", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("Delete", ctx, userID).Return(nil)

		err := userService.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("success - owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		mockRepo.On("Delete", ctx, userID).Return(nil)

		err := userService.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("permission denied - not owner", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeRegular)

		err := userService.DeleteUser(ctx, otherUserID)

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		assert.EqualError(t, err, "permission denied")
	})

	t.Run("repository error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("Delete", ctx, userID).Return(errors.New("db error"))

		err := userService.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertCalled(t, "Delete", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", userID)
		ctx = context.WithValue(ctx, "role", model.UserTypeAdmin)

		mockRepo.On("Delete", ctx, userID).Return(errors.New("no user found with id"))

		err := userService.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "no user found with id")
		mockRepo.AssertCalled(t, "Delete", ctx, userID)
		mockRepo.ExpectedCalls = nil
	})
}

func TestAuthenticateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 14)
	user := model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(passwordHash),
		AccountType:  model.UserTypeRegular,
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

		authResponse, err := userService.AuthenticateUser(ctx, "test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, authResponse.Token)
		assert.True(t, authResponse.ExpiresAt.After(time.Now()))
		mockRepo.AssertCalled(t, "GetByEmail", ctx, "test@example.com")
		mockRepo.ExpectedCalls = nil
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("GetByEmail", ctx, "test@example.com").Return(model.User{}, errors.New("user not found"))

		authResponse, err := userService.AuthenticateUser(ctx, "test@example.com", "password123")

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
		assert.EqualError(t, err, "user not found")
		assert.Empty(t, authResponse.Token)
		mockRepo.AssertCalled(t, "GetByEmail", ctx, "test@example.com")
		mockRepo.ExpectedCalls = nil
	})

	t.Run("invalid credentials", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

		authResponse, err := userService.AuthenticateUser(ctx, "test@example.com", "wrongpassword")

		var httpErr *custom_errors.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
		assert.EqualError(t, err, "invalid credentials")
		assert.Empty(t, authResponse.Token)
		mockRepo.AssertCalled(t, "GetByEmail", ctx, "test@example.com")
		mockRepo.ExpectedCalls = nil
	})
}
