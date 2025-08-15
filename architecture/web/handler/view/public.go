package view

import (
	"log"
	"net/http"
)

func NewView(templatesDir string) *View {
	return &View{templatesDir: templatesDir}
}

func (v *View) ExecuteTemplate(w http.ResponseWriter, pg interface{}, names ...string) {
	tmpl, err := v.getTemplate(names...)
	if err != nil {
		log.Printf("m.newView: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main", pg)
	if err != nil {
		log.Printf("tmpl.ExecuteTemplate: %v", err)
		return
	}
}
