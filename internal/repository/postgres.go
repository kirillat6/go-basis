package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kirillat6/go-basis/internal/errs"
	"github.com/kirillat6/go-basis/internal/task"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, title string) (*task.Task, error) {
	query := `
	INSERT INTO tasks (title, completed)
	VALUES ($1, $2)
	RETURNING id, title, completed;
	`

	t := &task.Task{}
	err := r.db.QueryRow(ctx, query, title, false).Scan(
		&t.ID,
		&t.Title,
		&t.Completed,
	)

	if err != nil {
		return nil, fmt.Errorf("Не удалось добавить задачу: %w", err)
	}

	return t, nil

}

func (r *TaskRepository) DeleteTask(ctx context.Context, id int) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1 
	`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Не удалось удалить задачу: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("Не удалось найти задачу: %w", errs.ErrNotFound)
	}

	return nil
}

func (r *TaskRepository) GetTasks(ctx context.Context) ([]task.Task, error) {
	query := `
	SELECT id, title, completed FROM tasks
	ORDER BY id ASC;
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []task.Task{}, fmt.Errorf("Не удалось найти задачи: %w", err)
	}
	defer rows.Close()
	tasks := []task.Task{}
	for rows.Next() {
		t := task.Task{}
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Completed,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) GetTask(ctx context.Context, id int) (task.Task, error) {
	query := `
	SELECT id, title, completed FROM tasks
	WHERE id = $1;
	`
	currentTask := task.Task{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&currentTask.ID,
		&currentTask.Title,
		&currentTask.Completed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return task.Task{}, fmt.Errorf("Задача не найдена: %w", errs.ErrNotFound)
		}
		return task.Task{}, err
	}
	return currentTask, nil
}

func (r *TaskRepository) ChangeTask(ctx context.Context, id int, isComplete *bool, title *string) error {
	if isComplete == nil && title == nil {
		return nil
	}

	query := `UPDATE tasks SET `
	args := []any{}
	argID := 1

	if title != nil {
		query += fmt.Sprintf("title = $%d, ", argID)
		args = append(args, *title)
		argID++
	}

	if isComplete != nil {
		query += fmt.Sprintf("completed = $%d, ", argID)
		args = append(args, *isComplete)
		argID++
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)
	cmdTag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("Задача не найдена: %w", errs.ErrNotFound)
	}
	return nil
}
