package main

import (
	"context"
	"errors"
	"forum/handlers"
	"forum/models"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	models.InitDB("./forum.db")
	process, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	assets := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assets))
	mux.HandleFunc("GET /", handlers.IndexHandler)

	mux.HandleFunc("/register", handlers.RegisterHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: assets,
	}

	go func() {
		slog.Info("http server listen on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			stop()
		}
	}()

	<-process.Done()

	slog.Info("received interrupt signal")

	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("server stopped")
	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
