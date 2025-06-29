package handlers

import (
	"fmt"
	"forum/models"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RegisterGetHandler(w, r)
	}
	if r.Method == http.MethodPost {
		RegisterPostHandler(w, r)
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}

func RegisterGetHandler(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmple.Execute(w, nil)
	return
}

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = models.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("Registering user:", username, email, password)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}
