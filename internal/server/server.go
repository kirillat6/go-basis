package server

import (
	"net/http"

	"github.com/kirillat6/go-basis/internal/config"
	"github.com/kirillat6/go-basis/internal/handler"
	"github.com/kirillat6/go-basis/internal/middlewares"
)

func New(repo handler.TaskRepositoryInterface, cfg config.Config) *http.Server {
	mux := http.NewServeMux()
	h := handler.NewHandler(repo)
	mux.HandleFunc("GET /tasks", h.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", h.GetTask)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
	mux.HandleFunc("PATCH /tasks/{id}", h.ChangeTask)

	wrappedMux := middlewares.Logger(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: wrappedMux,
	}
	return server
}
