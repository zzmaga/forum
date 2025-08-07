package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func SwitchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/register":
		RegisterHandler(w, r)
	case "/login":
		LoginHandler(w, r)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "", 500)
		log.Fatal("500")
	}
	tmpl.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "", 500)
		}
		tmpl.Execute(w, nil)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "", 500)
		}
		tmpl.Execute(w, nil)
		return
	}
}
