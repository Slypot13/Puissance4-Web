package handlers

import (
	"html/template"
	"net/http"
)

func EndHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/end.gohtml", "templates/layout.gohtml"))
	tmpl.ExecuteTemplate(w, "layout", nil)
}
