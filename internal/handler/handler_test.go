package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kirillat6/go-basis/internal/task"
)

type fakeRepository struct {
	fakeCreateTask func(ctx context.Context, title string) (*task.Task, error)
	fakeDeleteTask func(ctx context.Context, id int) error
	fakeGetTask    func(ctx context.Context, id int) (task.Task, error)
	fakeGetTasks   func(ctx context.Context) ([]task.Task, error)
	fakeChangeTask func(ctx context.Context, id int, completed *bool, title *string) error
}

func (f *fakeRepository) CreateTask(ctx context.Context, title string) (*task.Task, error) {
	return f.fakeCreateTask(ctx, title)
}

func (f *fakeRepository) DeleteTask(ctx context.Context, id int) error {
	return f.fakeDeleteTask(ctx, id)
}

func (f *fakeRepository) GetTask(ctx context.Context, id int) (task.Task, error) {
	return f.fakeGetTask(ctx, id)
}

func (f *fakeRepository) GetTasks(ctx context.Context) ([]task.Task, error) {
	return f.fakeGetTasks(ctx)
}

func (f *fakeRepository) ChangeTask(ctx context.Context, id int, completed *bool, title *string) error {
	return f.fakeChangeTask(ctx, id, completed, title)
}

func TestGetTask(t *testing.T) {
	fakeRepo := &fakeRepository{
		fakeGetTask: func(ctx context.Context, id int) (task.Task, error) {
			return task.Task{
				ID:        42,
				Title:     "Test task",
				Completed: false,
			}, nil
		},
	}

	h := NewHandler(fakeRepo)

	r := httptest.NewRequest(http.MethodGet, "/tasks/42", nil)
	r.SetPathValue("id", "42")

	w := httptest.NewRecorder()

	h.GetTask(w, r)

	var got task.Task

	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Errorf("Не удалось декодировать JSON: %v", err)
		return
	}

	if got.ID != 42 {
		t.Errorf("Пришёл неожиданный ID. Ожидалось: %v. Пришло: %v",
			42,
			got.ID,
		)
		return
	}

	if got.Title != "Test task" {
		t.Errorf("Пришёл неожиданный Title. Ожидалось: %v. Пришло: %v",
			"Test task",
			got.Title,
		)
		return
	}

	if got.Completed {
		t.Errorf("Пришёл неожиданный Completed. Ожидалось: %v. Пришло: %v",
			false,
			got.Completed,
		)
		return
	}

	if w.Code != http.StatusOK {
		t.Errorf("Пришёл не соответствующий код. Ожидалось: %v. Пришло: %v",
			http.StatusOK,
			w.Code,
		)
		return
	}

}
