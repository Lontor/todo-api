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
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;not null;" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"not null" json:"-" validate:"required"`
	AccountType  UserType  `gorm:"not null" json:"accountType" validate:"required"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`
}
