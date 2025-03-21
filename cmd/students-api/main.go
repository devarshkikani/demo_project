package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devarshkikani/demo_project/internal/config"
	"github.com/devarshkikani/demo_project/internal/http/handlers/student"
	"github.com/devarshkikani/demo_project/internal/storage/sqlite"
)

func main() {

	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initalized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))

	router.HandleFunc("GET /api/students/{id}", student.GetStudentById(storage))

	router.HandleFunc("GET /api/students", student.GetList(storage))

	router.HandleFunc("POST /api/students/update", student.UpdateStudent(storage))

	router.HandleFunc("POST /api/students/delete/{id}", student.DeleteStudent(storage))

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started %s", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to load server")
		}
	}()

	<-done

	slog.Info("Sutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		slog.Error("Failed to shutdon server", slog.String("error", err.Error()))
	}

	slog.Info("Server shotdown succeully")

}
