package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.gohtml", "templates/layout.gohtml"))
	tmpl.ExecuteTemplate(w, "layout", nil)
}
