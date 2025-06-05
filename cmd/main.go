package main

import (
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	models.InitDB("./forum.db")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Forum is working"))
	})
	mux.HandleFunc("/register", handlers.RegisterHandler)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
