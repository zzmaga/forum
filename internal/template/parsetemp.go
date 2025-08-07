package templates

import (
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func RenderTemplate(w http.ResponseWriter, name string, data any) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Ошибка шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}
