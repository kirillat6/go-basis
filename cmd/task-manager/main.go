package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kirillat6/go-basis/internal/task"
)

func main() {
	tm := task.NewTaskManager()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		tasks := tm.GetTasks()
		err := json.NewEncoder(w).Encode(tasks)
		if err != nil {
			http.Error(w, "{\"message\": \"Произошла ошибка на сервере\"}", http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request){
		strId := r.PathValue("id")
		if strId == "" {
			http.Error(w, "{\"message\": \"Задача с таким id не найден!\"}", http.StatusNotFound)
			return
		}
		id, err := strconv.Atoi(strId)
		if err != nil {
			http.Error(w, "{\"message\": \"ID должен быть цифрой!\"}", http.StatusBadRequest)
			return
		}
		task, err := tm.GetTask(id)
		if err != nil {
			http.Error(w, "{\"message\": \"Задача с таким id не найден!\"}", http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, "{\"message\":\"Проблема кодировки\"}", http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request){
		var req task.TaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "{\"message\": \"Некорректный формат JSON\"}", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Title) == "" {
			http.Error(w, "{\"message\": \"Поле title не может быть пустым\"}", http.StatusBadRequest)
			return
		}

		task := tm.CreateTask(req.Title)
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(task) 
		if err != nil {
			http.Error(w, "{\"message\":\"Проблема кодировки\"}", http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request){
		strId := r.PathValue("id")
		if strId == "" {
			http.Error(w, "{\"message\": \"Задача с таким id не найден!\"}", http.StatusNotFound)
			return
		}
		id, err := strconv.Atoi(strId)
		if err != nil {
			http.Error(w, "{\"message\": \"ID должен быть цифрой!\"}", http.StatusBadRequest)
			return
		}
		err = tm.DeleteTask(id)
		if err != nil {
			http.Error(w, "{\"message\": \"Задача с таким id не найден!\"}", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("PATCH /tasks/{id}", func(w http.ResponseWriter, r *http.Request){
		var req = task.TaskPatchRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "{\"message\": \"Некорректный формат JSON\"}", http.StatusBadRequest)
			return
		}
		strId := r.PathValue("id")
		if strId == "" {
			http.Error(w, "{\"message\": \"Задача с таким id не найден!\"}", http.StatusNotFound)
			return
		}
		id, err := strconv.Atoi(strId)
		if err != nil {
			http.Error(w, "{\"message\": \"ID должен быть цифрой!\"}", http.StatusBadRequest)
			return
		}

		var titlePtr *string
		if req.Title != nil {
			trimmed := strings.TrimSpace(*req.Title)
			titlePtr  = &trimmed
		}

		err = tm.ChangeTask(id, req.Completed, titlePtr)
		if err != nil {
			http.Error(w, "{\"message\": \"Задача с таким id не найден!\"}", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Сервер завершил работу с ошибкой!")
	}
}
