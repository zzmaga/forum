package view

import (
	"fmt"
	"html/template"
)

func (v *View) getTemplate(names ...string) (*template.Template, error) {
	paths := []string{v.templatesDir + "/main.html", v.templatesDir + "/navbar.html", v.templatesDir + "/alert.html"}
	for _, name := range names {
		paths = append(paths, v.templatesDir+"/"+name)
	}

	t, err := template.ParseFiles(paths...)
	if err != nil {
		return nil, fmt.Errorf("template.ParseFiles: %w", err)
	}
	return t, nil
}
