package server

import (
	"net/http"

	"github.com/kirillat6/go-basis/internal/config"
	"github.com/kirillat6/go-basis/internal/handler"
	"github.com/kirillat6/go-basis/internal/middlewares"
	"github.com/kirillat6/go-basis/internal/repository"
)



func New(repo *repository.TaskRepository, cfg config.Config) *http.Server{
	mux := http.NewServeMux()
	handlers := handler.NewHandler(repo)
	mux.HandleFunc("GET /tasks", handlers.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", handlers.GetTask)
	mux.HandleFunc("POST /tasks", handlers.CreateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handlers.DeleteTask)
	mux.HandleFunc("PATCH /tasks/{id}", handlers.ChangeTask)

	wrappedMux := middlewares.Logger(mux)

	server := &http.Server{
		Addr: cfg.Port,
		Handler: wrappedMux,
	}
	return server
}

