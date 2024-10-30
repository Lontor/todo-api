package repository

import (
	"context"
	"database/sql"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/google/uuid"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `INSERT INTO tasks (user_id, description, status, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query, task.UserID, task.Description, task.Status)
	return err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	var task model.Task
	query := `SELECT id, user_is, description, status, created_at, updated_at FROM tasks WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.UserID, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Task, error) {
	var tasks []*model.Task
	query := `SELECT id, user_id, description, status, created_at, updated_at FROM tasks WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (r *taskRepository) UpdateDescription(ctx context.Context, id uuid.UUID, description string) error {
	query := `UPDATE tasks SET description = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, description, id)
	return err
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.TaskStatus) error {
	query := `UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
