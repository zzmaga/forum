package routes

import (
	"forum/handlers"
	"net/http"
)

func GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	return mux
}
