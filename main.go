package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)


var templates = template.Must(template.ParseGlob("templates/*.html"))

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
	Player1       Player
	Player2       Player
	Started       time.Time
	Turns         int
	CurrentPlayer int       
	Grid          [][]int   
}

var (
	currentGame *CurrentGame
	scoreboard  []GameResult
)


type PageData struct {
	Title        string
	CurrentGame  *CurrentGame
	Scoreboard   []GameResult
	LastResult   *GameResult
	ErrorCode    int
	ErrorMessage string
}


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

func renderTemplate(w http.ResponseWriter, name string, data PageData) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("Erreur template:", err)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Erreur interne du serveur."))
	}
}

func redirectToError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	u := "/error?code=" + url.QueryEscape(strconv.Itoa(code)) +
		"&msg=" + url.QueryEscape(msg)
	http.Redirect(w, r, u, http.StatusSeeOther)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {

		redirectToError(w, r, 404, "Page introuvable.")
		return
	}

	data := PageData{
		Title: "Power'4 Web — Accueil",
	}
	renderTemplate(w, "home.html", data)
}

func gameInitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, 405, "Méthode HTTP non autorisée.")
		return
	}

	data := PageData{
		Title: "Power'4 Web — Initialisation",
	}
	renderTemplate(w, "game_init.html", data)
}

func gameInitTraitementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		redirectToError(w, r, 405, "Méthode HTTP non autorisée.")
		return
	}

	p1Name := normalizeName(r.FormValue("player1_name"))
	p2Name := normalizeName(r.FormValue("player2_name"))
	p1Color := normalizeColor(r.FormValue("player1_color"))
	p2Color := normalizeColor(r.FormValue("player2_color"))

	if !nameRegex.MatchString(p1Name) || !nameRegex.MatchString(p2Name) {
		redirectToError(w, r, 400, "Les pseudos doivent contenir entre 3 et 20 caractères (lettres, chiffres, -, _).")
		return
	}

	if p1Color == "" || p2Color == "" {
		redirectToError(w, r, 400, "Couleur de jeton invalide.")
		return
	}
	if p1Color == p2Color {
		redirectToError(w, r, 400, "Les deux joueurs doivent avoir une couleur différente.")
		return
	}

	currentGame = &CurrentGame{
		Player1: Player{
			Name:  p1Name,
			Color: p1Color,
		},
		Player2: Player{
			Name:  p2Name,
			Color: p2Color,
		},
		Started:       time.Now(),
		Turns:         0,
		CurrentPlayer: 1,
		Grid:          make([][]int, 6),
	}

	for i := range currentGame.Grid {
		currentGame.Grid[i] = make([]int, 7)
	}

	http.Redirect(w, r, "/game/play", http.StatusSeeOther)
}

func gamePlayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, 405, "Méthode HTTP non autorisée.")
		return
	}

	if currentGame == nil {
		redirectToError(w, r, 400, "Aucune partie en cours. Veuillez en démarrer une nouvelle.")
		return
	}

	data := PageData{
		Title:       "Power'4 Web — Partie en cours",
		CurrentGame: currentGame,
	}
	renderTemplate(w, "game_play.html", data)
}

func gamePlayTraitementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		redirectToError(w, r, 405, "Méthode HTTP non autorisée.")
		return
	}

	if currentGame == nil {
		redirectToError(w, r, 400, "Aucune partie en cours.")
		return
	}

	colStr := strings.TrimSpace(r.FormValue("column"))
	if colStr == "" {
		redirectToError(w, r, 400, "Aucune colonne sélectionnée.")
		return
	}

	col, err := strconv.Atoi(colStr)
	if err != nil || col < 1 || col > 7 {
		redirectToError(w, r, 400, "Colonne invalide.")
		return
	}
	col-- // 1-7 → 0-6

	placed := false
	for row := len(currentGame.Grid) - 1; row >= 0; row-- {
		if currentGame.Grid[row][col] == 0 {
			currentGame.Grid[row][col] = currentGame.CurrentPlayer
			currentGame.Turns++
			placed = true
			break
		}
	}

	if !placed {
		redirectToError(w, r, 400, "Cette colonne est déjà pleine.")
		return
	}

	if checkWin(currentGame.Grid, currentGame.CurrentPlayer) {
		winnerName := currentGame.Player1.Name
		if currentGame.CurrentPlayer == 2 {
			winnerName = currentGame.Player2.Name
		}

		result := GameResult{
			Player1: currentGame.Player1.Name,
			Player2: currentGame.Player2.Name,
			Winner:  winnerName,
			Date:    time.Now(),
			Turns:   currentGame.Turns,
		}
		scoreboard = append(scoreboard, result)

		http.Redirect(w, r, "/game/end", http.StatusSeeOther)
		return
	}

	if boardFull(currentGame.Grid) {
		result := GameResult{
			Player1: currentGame.Player1.Name,
			Player2: currentGame.Player2.Name,
			Winner:  "",
			Date:    time.Now(),
			Turns:   currentGame.Turns,
		}
		scoreboard = append(scoreboard, result)

		http.Redirect(w, r, "/game/end", http.StatusSeeOther)
		return
	}

	if currentGame.CurrentPlayer == 1 {
		currentGame.CurrentPlayer = 2
	} else {
		currentGame.CurrentPlayer = 1
	}

	http.Redirect(w, r, "/game/play", http.StatusSeeOther)
}

func gameEndHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, 405, "Méthode HTTP non autorisée.")
		return
	}

	var last *GameResult
	if len(scoreboard) > 0 {
		last = &scoreboard[len(scoreboard)-1]
	}

	data := PageData{
		Title:      "Power'4 Web — Fin de partie",
		LastResult: last,
	}
	renderTemplate(w, "game_end.html", data)
}

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		redirectToError(w, r, 405, "Méthode HTTP non autorisée.")
		return
	}

	data := PageData{
		Title:      "Power'4 Web — Scoreboard",
		Scoreboard: scoreboard,
	}
	renderTemplate(w, "scoreboard.html", data)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := strings.TrimSpace(r.FormValue("code"))
	msg := strings.TrimSpace(r.FormValue("msg"))

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		code = 0
	}

	data := PageData{
		Title:        "Power'4 Web — Erreur",
		ErrorCode:    code,
		ErrorMessage: msg,
	}
	renderTemplate(w, "error.html", data)
}

func boardFull(grid [][]int) bool {
	for col := 0; col < 7; col++ {
		if grid[0][col] == 0 {
			return false
		}
	}
	return true
}

func checkWin(grid [][]int, player int) bool {
	rows := len(grid)
	cols := len(grid[0])

	for r := 0; r < rows; r++ {
		for c := 0; c <= cols-4; c++ {
			if grid[r][c] == player &&
				grid[r][c+1] == player &&
				grid[r][c+2] == player &&
				grid[r][c+3] == player {
				return true
			}
		}
	}

	for c := 0; c < cols; c++ {
		for r := 0; r <= rows-4; r++ {
			if grid[r][c] == player &&
				grid[r+1][c] == player &&
				grid[r+2][c] == player &&
				grid[r+3][c] == player {
				return true
			}
		}
	}

	for r := 0; r <= rows-4; r++ {
		for c := 0; c <= cols-4; c++ {
			if grid[r][c] == player &&
				grid[r+1][c+1] == player &&
				grid[r+2][c+2] == player &&
				grid[r+3][c+3] == player {
				return true
			}
		}
	}

	for r := 3; r < rows; r++ {
		for c := 0; c <= cols-4; c++ {
			if grid[r][c] == player &&
				grid[r-1][c+1] == player &&
				grid[r-2][c+2] == player &&
				grid[r-3][c+3] == player {
				return true
			}
		}
	}

	return false
}

func main() {

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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
