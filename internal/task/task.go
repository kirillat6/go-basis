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


func (h *TaskManager) CreateTask(title string) {
	h.nextID++
	newTask := Task{
		ID:        h.nextID,
		Title:     title,
		Completed: false,
	}
	h.tasks = append(h.tasks, newTask)
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

func (h *TaskManager) CompleteTask(id int) error{
	for i := range h.tasks {
		if h.tasks[i].ID == id {
			h.tasks[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("Задача не найдена!")
}

func (h *TaskManager) GetTasks() []Task {
	return h.tasks
}
