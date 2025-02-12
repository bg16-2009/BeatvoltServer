package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func loadTemplates() {
	tpl = template.New("")
	_, err := tpl.ParseGlob("templates/**/*.html")
	_, err = tpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("error parsing templates: %w", err)
	}
}

func renderTemplate(name string, w http.ResponseWriter, data any) {
	err := tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
}
