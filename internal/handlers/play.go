package handlers

import (
	"html/template"
	"net/http"
)

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/play.gohtml", "templates/layout.gohtml"))
	tmpl.ExecuteTemplate(w, "layout", nil)
}
