package main

import (
    "html/template"
    "log"
    "net/http"
    "path/filepath"
)

var templates *template.Template

func init() {
    // Parse tous les templates dans templates/*.gohtml
    var err error
    templates, err = template.ParseGlob(filepath.Join("templates", "*.gohtml"))
    if err != nil {
        log.Fatalf("Erreur lors du parsing des templates : %v", err)
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "home.gohtml", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    // Route vers page d'accueil
    http.HandleFunc("/", homeHandler)

    // Servir fichiers statiques (CSS, images)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    log.Println("Serveur démarré sur http://localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatalf("Erreur serveur : %v", err)
    }
}
