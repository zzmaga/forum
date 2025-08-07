package handlers

import (
	internal "forum/internal/template"
	"log"
	"net/http"
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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		log.Println("form parse error:", err)
		return
	}

	// name := r.Form.Get("name")
	// email := r.Form.Get("email")
	// password := r.Form.Get("password")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	internal.RenderTemplate(w, "login.html", nil)
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		log.Println("form parse error:", err)
		return
	}

	// email := r.Form.Get("email")
	// password := r.Form.Get("password")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	internal.RenderTemplate(w, "index.html", nil)
}
