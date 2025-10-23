package main

import (
	"fmt"
	"net/http"

	"Puissance4-Web/internal/handlers"
)

func main() {
	// Routes principales
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/game/init", handlers.InitHandler)
	http.HandleFunc("/game/play", handlers.PlayHandler)
	http.HandleFunc("/game/end", handlers.EndHandler)
	http.HandleFunc("/game/scoreboard", handlers.ScoreboardHandler)

	// Fichiers statiques (CSS, JS, images…)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("🌐 Serveur lancé sur http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}
