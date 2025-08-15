package template

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// RenderTemplate отображает HTML шаблон с данными
func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	// Путь к HTML файлам
	tmplPath := filepath.Join("ui", "html", templateName)
	
	// Парсим шаблон
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Устанавливаем заголовок
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Выполняем шаблон
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
} 