package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ---------- TEMPLATES ----------

var templates = template.Must(template.ParseGlob("templates/*.html"))

// ---------- TYPES MÉTIER ----------

type Player struct {
	Name  string
	Color string
}

type GameResult struct {
	Player1 string
	Player2 string
	Winner  string
	Date    time.Time
	Turns   int
}

type CurrentGame struct {
	Player1 Player
	Player2 Player
	Started time.Time
	Turns   int
	// Tu pourras ajouter ici : plateau, joueur courant, etc.
}

// ---------- DONNÉES GLOBALES (SIMPLIFIÉES) ----------

var scoreboard []GameResult
var currentGame *CurrentGame

// ---------- STRUCTURE GÉNÉRALE DES DONNÉES PAGE ----------

type PageData struct {
	Title        string
	ErrorCode    int
	ErrorMessage string
	CurrentGame  *CurrentGame
	Scoreboard   []GameResult
}

// ---------- OUTILS ----------

var nameRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{3,20}$`)

func normalizeName(s string) string {
	return strings.TrimSpace(s)
}

func normalizeColor(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "rouge" || s == "jaune" {
		return s
	}
	return ""
}

func redirectToError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	u := "/error?code=" + url.QueryEscape((http.StatusText(code))) +
		"&msg=" + url.QueryEscape(msg)
	http.Redirect(w, r, u, http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data PageData) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("Erreur template:", err)
		// Dernier recours : simple texte (on évite boucle /error infinie)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Erreur interne du serveur."))
	}
}

// ---------- HANDLERS ----------

// GET /
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		redirectToError(w, r, http.StatusNotFound, "Page introuvable.")
		return
	}

	data := PageData{
		Title: "Power'4 Web — Accueil",
	}
	renderTemplate(w, r, "home.html", data)
}

// GET /game/init
func gameInitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée.")
		return
	}

	data := PageData{
		Title: "Power'4 Web — Initialisation",
	}
	renderTemplate(w, r, "game_init.html", data)
}

// POST /game/init/traitement
func gameInitTraitementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		redirectToError(w, r, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée.")
		return
	}

	// Utilisation de FormValue + normalisation
	p1Name := normalizeName(r.FormValue("player1_name"))
	p2Name := normalizeName(r.FormValue("player2_name"))
	p1Color := normalizeColor(r.FormValue("player1_color"))
	p2Color := normalizeColor(r.FormValue("player2_color"))

	// Validation des pseudos
	if !nameRegex.MatchString(p1Name) || !nameRegex.MatchString(p2Name) {
		redirectToError(w, r, http.StatusBadRequest, "Les pseudos doivent faire 3 à 20 caractères (lettres, chiffres, -, _).")
		return
	}

	// Validation des couleurs
	if p1Color == "" || p2Color == "" {
		redirectToError(w, r, http.StatusBadRequest, "Couleur de jeton invalide.")
		return
	}
	if p1Color == p2Color {
		redirectToError(w, r, http.StatusBadRequest, "Les deux joueurs doivent avoir une couleur différente.")
		return
	}

	// Initialisation de la partie
	currentGame = &CurrentGame{
		Player1: Player{Name: p1Name, Color: p1Color},
		Player2: Player{Name: p2Name, Color: p2Color},
		Started: time.Now(),
		Turns:   0,
	}

	// Redirection vers la page de jeu
	http.Redirect(w, r, "/game/play", http.StatusSeeOther)
}

// GET /game/play
func gamePlayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée.")
		return
	}

	if currentGame == nil {
		redirectToError(w, r, http.StatusBadRequest, "Aucune partie en cours. Veuillez initialiser une nouvelle partie.")
		return
	}

	data := PageData{
		Title:       "Power'4 Web — Partie en cours",
		CurrentGame: currentGame,
	}
	renderTemplate(w, r, "game_play.html", data)
}

// POST /game/play/traitement
func gamePlayTraitementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		redirectToError(w, r, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée.")
		return
	}

	if currentGame == nil {
		redirectToError(w, r, http.StatusBadRequest, "Aucune partie en cours.")
		return
	}

	// Exemple : récupération de la colonne jouée
	colStr := strings.TrimSpace(r.FormValue("column"))
	if colStr == "" {
		redirectToError(w, r, http.StatusBadRequest, "Aucune colonne sélectionnée.")
		return
	}

	// TODO : convertir colStr en int, vérifier les bornes, appliquer le coup sur le plateau,
	// vérifier victoire/égalité, incrémenter currentGame.Turns, etc.

	// TODO : si victoire ou égalité → remplir un GameResult, l’ajouter à scoreboard, puis rediriger vers /game/end
	// sinon → rediriger vers /game/play

	// Pour l’instant, simple retour vers la page de jeu
	http.Redirect(w, r, "/game/play", http.StatusSeeOther)
}

// GET /game/end
func gameEndHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée.")
		return
	}

	data := PageData{
		Title: "Power'4 Web — Fin de partie",
		// Tu pourras passer ici le résultat de la dernière partie
	}
	renderTemplate(w, r, "game_end.html", data)
}

// GET /game/scoreboard
func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée.")
		return
	}

	data := PageData{
		Title:      "Power'4 Web — Scoreboard",
		Scoreboard: scoreboard,
	}
	renderTemplate(w, r, "scoreboard.html", data)
}

// GET /error
func errorHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := strings.TrimSpace(r.FormValue("code")) // ex : "Not Found"
	msg := strings.TrimSpace(r.FormValue("msg"))

	data := PageData{
		Title:        "Erreur — Power'4 Web",
		ErrorCode:    0, // si tu veux, tu peux parser un vrai code numérique
		ErrorMessage: msg,
	}

	renderTemplate(w, r, "error.html", data)
}

// ---------- MAIN ----------

func main() {
	// Dossier assets exposé via /static/
	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", assets))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game/init", gameInitHandler)
	http.HandleFunc("/game/init/traitement", gameInitTraitementHandler)
	http.HandleFunc("/game/play", gamePlayHandler)
	http.HandleFunc("/game/play/traitement", gamePlayTraitementHandler)
	http.HandleFunc("/game/end", gameEndHandler)
	http.HandleFunc("/game/scoreboard", scoreboardHandler)
	http.HandleFunc("/error", errorHandler)

	log.Println("Serveur démarré sur http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
