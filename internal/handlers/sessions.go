package handlers

import (
	"database/sql"
	"forum/internal/database"
	"log"
	"net/http"
	"time"
)

func GetUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, err
	}

	var userID int
	var expiresAtStr string
	err = database.DB.QueryRow(
		"SELECT user_id, expired_at FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userID, &expiresAtStr)

	if err != nil {
		return 0, err
	}

	expiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		// Сессия просрочена
		return 0, sql.ErrNoRows
	}

	return userID, nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_, err := database.DB.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		if err != nil {
			log.Printf("Warning: failed to delete session: %v", err)
		}
	}

	// удаляем cookie у клиента
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
