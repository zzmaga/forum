package handlers

import (
	"database/sql"
	"forum/internal/database"
	internal "forum/internal/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	internal.RenderTemplate(w, "index.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetRegisterHandler(w, r)
	case http.MethodPost:
		PostRegisterHandler(w, r)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetLoginHandler(w, r)
	case http.MethodPost:
		PostLoginHandler(w, r)
	}
}

func GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
	internal.RenderTemplate(w, "register.html", nil)
}

func GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	internal.RenderTemplate(w, "login.html", nil)
}

func PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		log.Println("form parse error:", err)
		return
	}

	username := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	// log.Printf("Register: name=%s, email=%s", name, email) // временно, для отладки

	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already taken", http.StatusBadRequest)
		return
	}

	// hashing password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}

	// Saving user
	_, err = database.DB.Exec("INSERT INTO users(username, email, password_hash) VALUES (?, ?, ?)", username, email, string(hash))
	if err != nil {
		http.Error(w, "Database insert error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		log.Println("form parse error:", err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	// log.Printf("Login: email=%s", email)
	var id int
	var passwordHash string
	err := database.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", email).Scan(&id, &passwordHash)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Здесь потом добавим cookie-сессию
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
