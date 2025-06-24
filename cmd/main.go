package main

import (
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"
	"os"
)

func main() {
	models.InitDB()
	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
