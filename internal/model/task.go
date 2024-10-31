package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "to do"
	StatusInProgress TaskStatus = "in progress"
	StatusDone       TaskStatus = "Done"
)

type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;not null;" json:"id" validate:"required"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" json:"userId" validate:"required"`
	Description string     `gorm:"not null" json:"description" validate:"required,min=10,max=200"`
	Status      TaskStatus `gorm:"not null" json:"status" validate:"required"`
	CreatedAt   time.Time  `gorm:"not null" json:"createdAt" validate:"required"`
	UpdatedAt   time.Time  `gorm:"not null" json:"updatedAt" validate:"required"`
}
