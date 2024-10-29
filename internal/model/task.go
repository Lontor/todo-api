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
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"stasus"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
