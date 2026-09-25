package task

import "fmt"

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TaskManager struct {
	tasks []Task
	nextID int 
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: []Task{},
		nextID: 0,
	}
}

type TaskRequest struct {
	Title 	  string `json:"title"`
}
type TaskPatchRequest struct {
	Title 	  *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func (h *TaskManager) CreateTask(title string) Task {
	h.nextID++
	newTask := Task{
		ID:        h.nextID,
		Title:     title,
		Completed: false,
	}
	h.tasks = append(h.tasks, newTask)
	return newTask
}

func (h *TaskManager) DeleteTask(id int) error {
		for i := range h.tasks {
		if h.tasks[i].ID == id {
			h.tasks = append(h.tasks[:i], h.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Задача не найдена!")
}

func (h *TaskManager) ChangeTask(id int, isComplete *bool, title *string) error {
	for i := range h.tasks {
		if h.tasks[i].ID == id {
			if isComplete != nil {
				h.tasks[i].Completed = *isComplete
			}
			
			if title != nil {
				h.tasks[i].Title = *title
			}
			
			return nil 
		}
	}
	
	return fmt.Errorf("Задача не найдена!")
}

func (h *TaskManager) GetTasks() []Task {
	return h.tasks
}

func (h *TaskManager) GetTask(id int) (Task, error) {
	for _, task := range h.tasks{
		if task.ID == id {
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("Задача не найдена!")
}
