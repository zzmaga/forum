package handlers

import (
	"database/sql"
	"forum/internal/database"
	"forum/internal/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetPosts()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	categories, err := database.GetCategories()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	userID, _ := GetUserIDFromSession(r)

	data := map[string]interface{}{
		"Posts":      posts,
		"Categories": categories,
		"UserID":     userID,
	}

	template.RenderTemplate(w, "index.html", data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetRegisterHandler(w, r)
	case http.MethodPost:
		PostRegisterHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetLoginHandler(w, r)
	case http.MethodPost:
		PostLoginHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
	template.RenderTemplate(w, "register.html", nil)
}

func GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	template.RenderTemplate(w, "login.html", nil)
}

func PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		log.Println("form parse error:", err)
		return
	}

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	log.Println("username:", username)
	log.Println("email:", email)
	log.Println("password:", password)

	if username == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Проверка на уникальность email
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

	// Проверка на уникальность username
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}

	// Сохранение в БД
	_, err = database.DB.Exec("INSERT INTO users(username, email, password) VALUES (?, ?, ?)", username, email, string(hash))
	if err != nil {
		http.Error(w, "Database insert error", http.StatusInternalServerError)
		log.Println("insert error:", err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		log.Println("form parse error:", err)
		return
	}

	email := r.Form.Get("email")
	passwordFromForm := r.Form.Get("password")

	if email == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	} else if passwordFromForm == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	var id int
	var password string
	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&id, &password)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordFromForm))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Генерируем UUID сессии
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // сессия живёт 1 день

	// Сохраняем в БД
	_, err = database.DB.Exec(
		"INSERT INTO sessions(id, user_id, expired_at) VALUES (?, ?, ?)",
		sessionID, id, expiresAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})

	// Перенаправляем на главную
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
