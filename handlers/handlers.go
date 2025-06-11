package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error while opening files", http.StatusInternalServerError)
		log.Fatal("Error while opening file: ", err)
	}

	tmp.Execute()
}
