package repository

import (
	"context"
	"testing"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := userRepo.Create(context.Background(), user)

	assert.NoError(t, err)

	var createdUser model.User
	err = db.First(&createdUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.PasswordHash, createdUser.PasswordHash)
	assert.Equal(t, user.AccountType, createdUser.AccountType)
	assert.Equal(t, user.CreatedAt.Equal(createdUser.CreatedAt), true)
}

func TestCreateUser_WithExistingEmai(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := userRepo.Create(context.Background(), user)

	assert.NoError(t, err)
	var createdUser model.User
	err = db.First(&createdUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, createdUser.Email)

	user.ID = uuid.New()
	err = userRepo.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")

	err = db.First(&createdUser, user.ID).Error
	assert.Error(t, err)
}

func TestCreateUser_WithEmptyRequiredFields(t *testing.T) {
	clearDB()
	user := getTestUser()

	user.Email = ""
	err := userRepo.Create(context.Background(), user)

	assert.Error(t, err)
	err = db.First(&user, user.ID).Error
	assert.Error(t, err)

	user = getTestUser()
	user.PasswordHash = ""
	err = userRepo.Create(context.Background(), user)

	assert.Error(t, err)
	err = db.First(&user, user.ID).Error
	assert.Error(t, err)

	user = getTestUser()
	user.AccountType = ""
	err = userRepo.Create(context.Background(), user)

	assert.Error(t, err)
	err = db.First(&user, user.ID).Error
	assert.Error(t, err)
}

func TestGetUsers(t *testing.T) {
	clearDB()
	users := make([]model.User, 3)

	for i := range users {
		users[i] = getTestUser()

		err := db.Create(&users[i]).Error
		require.NoError(t, err)
	}

	receivedUsers, err := userRepo.Get(context.Background())
	assert.NoError(t, err)
	require.Equal(t, len(users), len(receivedUsers))

	for i, user := range receivedUsers {
		assert.NoError(t, err)
		assert.Equal(t, users[i].ID, user.ID)
		assert.Equal(t, users[i].Email, user.Email)
		assert.Equal(t, users[i].PasswordHash, user.PasswordHash)
		assert.Equal(t, users[i].AccountType, user.AccountType)
		assert.Equal(t, users[i].CreatedAt.Equal(user.CreatedAt), true)
	}
}

func TestGetUserByEmail(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := db.Create(&user).Error
	require.NoError(t, err)

	checkUser, err := userRepo.GetByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, checkUser.ID)
	assert.Equal(t, user.Email, checkUser.Email)
	assert.Equal(t, user.PasswordHash, checkUser.PasswordHash)
	assert.Equal(t, user.AccountType, checkUser.AccountType)
	assert.Equal(t, user.CreatedAt.Equal(checkUser.CreatedAt), true)

}

func TestGetUserByEmail_NotFound(t *testing.T) {
	clearDB()
	user := getTestUser()

	_, err := userRepo.GetByEmail(context.Background(), user.Email)
	assert.Error(t, err)
}

func TestGetUserByID(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := db.Create(&user).Error
	require.NoError(t, err)

	checkUser, err := userRepo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)

	assert.Equal(t, user.ID, checkUser.ID)
	assert.Equal(t, user.Email, checkUser.Email)
	assert.Equal(t, user.PasswordHash, checkUser.PasswordHash)
	assert.Equal(t, user.AccountType, checkUser.AccountType)
	assert.Equal(t, user.CreatedAt.Equal(checkUser.CreatedAt), true)
}

func TestGetUserByID_NotFound(t *testing.T) {
	clearDB()

	_, err := userRepo.GetByID(context.Background(), uuid.New())
	assert.Error(t, err)
}

func TestUpdateUser_Success(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := db.Create(&user).Error
	require.NoError(t, err)

	updatedUser := user
	updatedUser.Email = "updated@example.com"
	updatedUser.PasswordHash = "new_hash"
	updatedUser.AccountType = model.UserTypeAdmin

	err = userRepo.Update(context.Background(), updatedUser)
	assert.NoError(t, err)

	var checkUser model.User
	err = db.Where("id = ?", user.ID).First(&checkUser).Error
	require.NoError(t, err)
	assert.Equal(t, updatedUser.Email, checkUser.Email)
	assert.Equal(t, updatedUser.PasswordHash, checkUser.PasswordHash)
	assert.Equal(t, updatedUser.AccountType, checkUser.AccountType)
	assert.Equal(t, user.CreatedAt.Equal(checkUser.CreatedAt), true)
}

func TestUpdateUser_NotFound(t *testing.T) {
	clearDB()

	nonExistentUser := getTestUser()
	nonExistentUser.ID = uuid.New()
	nonExistentUser.Email = "doesnotexist@example.com"

	err := userRepo.Update(context.Background(), nonExistentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no user found with id")
}

func TestUpdateUser_PartialUpdate(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := db.Create(&user).Error
	require.NoError(t, err)

	partialUpdate := user
	partialUpdate.Email = "partialupdate@example.com"

	err = userRepo.Update(context.Background(), partialUpdate)
	assert.NoError(t, err)

	var checkUser model.User
	err = db.Where("id = ?", user.ID).First(&checkUser).Error
	require.NoError(t, err)
	assert.Equal(t, partialUpdate.Email, checkUser.Email)
	assert.Equal(t, user.PasswordHash, checkUser.PasswordHash)
	assert.Equal(t, user.AccountType, checkUser.AccountType)
	assert.Equal(t, user.CreatedAt.Equal(checkUser.CreatedAt), true)
}

func TestUpdate_NoFieldsToUpdate(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := db.Create(&user).Error
	require.NoError(t, err)

	noUpdate := user
	noUpdate.Email = ""
	noUpdate.PasswordHash = ""
	noUpdate.AccountType = ""

	err = userRepo.Update(context.Background(), noUpdate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fields to update")
}

func TestDelete_Success(t *testing.T) {
	clearDB()
	user := getTestUser()

	err := db.Create(&user).Error
	require.NoError(t, err)

	err = userRepo.Delete(context.Background(), user.ID)
	assert.NoError(t, err)

	_, err = userRepo.GetByID(context.Background(), user.ID)
	assert.Error(t, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	clearDB()

	nonExistentID := uuid.New()
	err := userRepo.Delete(context.Background(), nonExistentID)
	assert.Error(t, err)
}
