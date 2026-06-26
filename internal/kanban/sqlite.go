package kanban

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database - %w", err)
	}
	sqlString := `CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`
	_, err = conn.Exec(sqlString)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: conn}, nil
}

func (ss *SQLiteStore) GetTasks(ctx context.Context) ([]Task, error) {
	sqlString := `SELECT id, title,	description, status, created_at, updated_at
		FROM tasks;`
	rows, err := ss.db.QueryContext(ctx, sqlString)
	if err != nil {
		return nil, fmt.Errorf("get tasks query failed - %w", err)
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status,
			&t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan task failed - %w", err)
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in row iteration - %w", err)
	}
	return tasks, nil
}

func (ss *SQLiteStore) CreateTask(ctx context.Context, t Task) error {
	_, err := ss.db.ExecContext(ctx, `INSERT INTO tasks
	(id, title,	description, status, created_at, updated_at) 
	values(?, ?, ?, ?, ?, ?)`, t.ID, t.Title, t.Description, t.Status,
		t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}

func (ss *SQLiteStore) UpdateTaskStatus(ctx context.Context, id uuid.UUID, s Status) error {
	res, err := ss.db.ExecContext(ctx, `UPDATE tasks SET
	status = ?, updated_at = ? WHERE id = ? `, s, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update - %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows effected - %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (ss *SQLiteStore) DeleteTask(ctx context.Context, id uuid.UUID) error {
	res, err := ss.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete - %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows effected - %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
