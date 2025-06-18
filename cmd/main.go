package main

import (
	"context"
	"errors"
	"forum/handlers"
	"forum/models"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	models.InitDB()
	process, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	assets := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assets))
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/post", handlers.Post)
	mux.HandleFunc("/comment", handlers.Comment)
	mux.HandleFunc("/like", handlers.Like)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
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
}
