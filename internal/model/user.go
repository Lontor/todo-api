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
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AccountType  UserType  `json:"AccountType"`
	CreatedAt    time.Time `json:"createdAt"`
}
