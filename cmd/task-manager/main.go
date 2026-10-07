package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kirillat6/go-basis/internal/config"
	"github.com/kirillat6/go-basis/internal/repository"
	"github.com/kirillat6/go-basis/internal/server"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Println(".env файл не найден, используются переменные окружения")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка конфигурации:", err)
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к БД")
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("Ошибка подключения к БД")
	}
	repo := repository.NewTaskRepository(pool)

	serv := server.New(repo, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		err := serv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Сервер завершил работу...")
		} else {
			log.Fatal("Ошибка сервера:", err)
		}
	}()
	<-stop
	log.Println("Получен сигнал завершения. Останавливаем сервер...")

	shutdownCtx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	err = serv.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("Ошибка при остановку сервера: %v", err)
	}
}
