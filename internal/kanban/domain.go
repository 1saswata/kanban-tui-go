package kanban

import (
	"time"

	"github.com/google/uuid"
)

type Status string

type Task struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      Status    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TaskStore interface {
	GetTasks() []Task
	CreateTask(Task) error
	UpdateTaskStatus(uuid.UUID, Status) error
	DeleteTask(uuid.UUID) error
}
