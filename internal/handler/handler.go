package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kirillat6/go-basis/internal/errs"
	"github.com/kirillat6/go-basis/internal/repository"
	"github.com/kirillat6/go-basis/internal/task"
)

type Handler struct {
	repo *repository.TaskRepository
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewHandler(repo *repository.TaskRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		tasks, err := h.repo.GetTasks(r.Context())
		if err != nil {
			writeError(w, "Произошла ошибка при получении задач",http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			writeError(w, "Произошла ошибка на сервере", http.StatusInternalServerError)
			return
		}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetId(r)
		if handleError(w, err) {
			return
		} 
		task, err := h.repo.GetTask(r.Context(), id)
		if handleError(w, err) {
			return
		}
		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			writeError(w, "Проблема кодировки", http.StatusInternalServerError)
			return
		}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req task.TaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, "Некорректный формат JSON",http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Title) == "" {
			writeError(w, "Поле title не может быть пустым",http.StatusBadRequest)
			return
		}

		task, err := h.repo.CreateTask(r.Context(), req.Title)
		if err != nil {
			writeError(w, "Произошла ошибка при добавлении задачи",http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(task) 
		if err != nil {
			writeError(w, "Проблема кодировки",http.StatusInternalServerError)
			return
		}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetId(r)
		if errors.Is(err, errs.ErrNotFound) {
			writeError(w, "Задача с таким id не найден!", http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrBadRequest) {
			writeError(w, "ID должен быть цифрой!", http.StatusBadRequest)
    		return
		}
		if err != nil {
			writeError(w, "Проблемы с сервером!!", http.StatusInternalServerError)
			return
		}
		err = h.repo.DeleteTask(r.Context(), id)
		if errors.Is(err, errs.ErrNotFound) {
			writeError(w, "Задача с таким id не найден!", http.StatusNotFound)
			return
		}
		if err != nil {
			writeError(w, "Проблемы с сервером!!", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ChangeTask(w http.ResponseWriter, r *http.Request) {
	var req = task.TaskPatchRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, "Некорректный формат JSON", http.StatusBadRequest)
			return
		}

		var titlePtr *string
		if req.Title != nil {
			trimmed := strings.TrimSpace(*req.Title)
			if trimmed != "" {
				titlePtr  = &trimmed
			} else {
				writeError(w, "Поле title не может быть пустым", http.StatusBadRequest)
				return
			}
		}
		id, err := GetId(r)
		if errors.Is(err, errs.ErrNotFound) {
			writeError(w, "Задача с таким id не найден!", http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrBadRequest) {
			writeError(w, "ID должен быть цифрой!", http.StatusBadRequest)
    		return
		}
		if err != nil {
			writeError(w, "Проблемы с сервером!", http.StatusInternalServerError)
			return
		} 
		err = h.repo.ChangeTask(r.Context(), id, req.Completed, titlePtr)
		if errors.Is(err, errs.ErrNotFound) {
			writeError(w, "Задача с таким id не найдена!", http.StatusNotFound)
			return
		}
		if err != nil {
			writeError(w, "Проблемы с сервером!", http.StatusInternalServerError)
			return
		} 
		w.WriteHeader(http.StatusOK)
}

func writeError(w http.ResponseWriter, message string, status int){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{Message: message}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
	}
}

func handleError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, errs.ErrNotFound) {
		writeError(w, "Задача с таким id не найден!",http.StatusNotFound)
		return true
	} else if errors.Is(err, errs.ErrBadRequest) {
		writeError(w, "ID должен быть цифрой!",http.StatusBadRequest)
		return true
	} 
	if err != nil {
		writeError(w, "Проблемы с сервером!", http.StatusInternalServerError)
		return true
	} 
	return false
}

func GetId(r *http.Request) (int, error) {
	strId := r.PathValue("id")
	if strId == "" {
		return 0, fmt.Errorf("Задача с таким id не найдена: %w", errs.ErrNotFound)
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		return 0, fmt.Errorf("ID должен быть цифрой: %w", errs.ErrBadRequest)
	}
	return id, nil
}
