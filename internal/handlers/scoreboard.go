package handlers

import (
	"html/template"
	"net/http"
)

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/scoreboard.gohtml", "templates/layout.gohtml"))
	tmpl.ExecuteTemplate(w, "layout", nil)
}
