package model

import (
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	UserTypeRegular UserType = "regular"
	UserTypeAdmin   UserType = "admin"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	PasswordHash string    `json:"-"`
	AccountType  UserType  `json:"accountType"`
	CreatedAt    time.Time `json:"createdAt"`
}
