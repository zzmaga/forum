package handlers

import (
	"forum/models"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		LoginGetHandler(w, r)
	}
	if r.Method == http.MethodPost {
		LoginPostHandler(w, r)
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}

func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	tmple, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmple.Execute(w, nil)
	return
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	var hashedPassword string
	var username string
	err := models.DB.QueryRow("SELECT username, password FROM users WHERE email = ?", email).Scan(&username, &hashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	}
	log.Println("Entering user", email, password)
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // сессия на 1 день

	// Получаем user_id
	var userID int
	err = models.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Сохраняем сессию в БД
	_, err = models.DB.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: expiresAt,
		Path:    "/",
	})

	log.Println("User logged in:", username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
