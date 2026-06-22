package kanban

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusTodo  Status = "TODO"
	StatusDoing Status = "DOING"
	StatusDone  Status = "DONE"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      Status    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TaskStore interface {
	GetTasks(ctx context.Context) ([]Task, error)
	CreateTask(ctx context.Context, t Task) error
	UpdateTaskStatus(ctx context.Context, id uuid.UUID, s Status) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}
