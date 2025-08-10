package main

import (
	"forum/internal/database"
	"forum/internal/handlers"
	"log"
	"net/http"
)

func main() {
	database.InitDB("forum.db")

	mux := http.NewServeMux()

	// Статические файлы
	fs := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// routes
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/create-post", handlers.CreatePostHandler)
	mux.HandleFunc("/post", handlers.ViewPostHandler)
	mux.HandleFunc("/comment", handlers.CreateCommentHandler)
	mux.HandleFunc("/like", handlers.LikeHandler)
	mux.HandleFunc("/filter", handlers.FilterPostsHandler)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
