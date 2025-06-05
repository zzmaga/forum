package handlers

import (
	"forum/models"
	"html/template"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
		return
	}

	// POST
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := models.RegisterUser(email, username, password)
	if err != nil {
		http.Error(w, "Ошибка регистрации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
