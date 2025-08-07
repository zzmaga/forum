package main

import (
	"forum/routes"
	"log"
	"net/http"
	// _ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := routes.GetRoutes()

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
