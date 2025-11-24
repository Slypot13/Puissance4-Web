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
	Name, Color string
}

type GameResult struct {
	Player1, Player2, Winner string
	Date                     time.Time
	Turns                    int
}

type CurrentGame struct {
	Player1, Player2 Player
	Started          time.Time
	Turns            int
	CurrentPlayer    int
	Grid             [][]int
}

type PageData struct {
	Title        string
	CurrentGame  *CurrentGame
	Scoreboard   []GameResult
	LastResult   *GameResult
	ErrorCode    int
	ErrorMessage string
}

var (
	currentGame *CurrentGame
	scoreboard  []GameResult
)

var nameRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{3,20}$`)

/* ---------- helpers ---------- */
func render(w http.ResponseWriter, tpl string, data PageData) {
	if err := templates.ExecuteTemplate(w, tpl, data); err != nil {
		log.Println("template error:", err)
		http.Error(w, "Erreur serveur", 500)
	}
}

func redirectErr(w http.ResponseWriter, r *http.Request, code int, msg string) {
	u := "/error?code=" + url.QueryEscape(strconv.Itoa(code)) + "&msg=" + url.QueryEscape(msg)
	http.Redirect(w, r, u, http.StatusSeeOther)
}

func mustMethod(w http.ResponseWriter, r *http.Request, want string) bool {
	if r.Method != want {
		redirectErr(w, r, 405, "Méthode HTTP non autorisée")
		return false
	}
	return true
}

func normalizeColor(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "rouge" || s == "jaune" {
		return s
	}
	return ""
}

func boardFull(g [][]int) bool {
	for c := 0; c < 7; c++ {
		if g[0][c] == 0 {
			return false
		}
	}
	return true
}

func checkWin(grid [][]int, p int) bool {
	for r := 0; r < 6; r++ {
		for c := 0; c < 7; c++ {
			if c+3 < 7 && grid[r][c] == p && grid[r][c+1] == p && grid[r][c+2] == p && grid[r][c+3] == p {
				return true
			}
			if r+3 < 6 && grid[r][c] == p && grid[r+1][c] == p && grid[r+2][c] == p && grid[r+3][c] == p {
				return true
			}
			if r+3 < 6 && c+3 < 7 && grid[r][c] == p && grid[r+1][c+1] == p && grid[r+2][c+2] == p && grid[r+3][c+3] == p {
				return true
			}
			if r-3 >= 0 && c+3 < 7 && grid[r][c] == p && grid[r-1][c+1] == p && grid[r-2][c+2] == p && grid[r-3][c+3] == p {
				return true
			}
		}
	}
	return false
}

/* ---------- handlers ---------- */
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		redirectErr(w, r, 404, "Page introuvable")
		return
	}
	render(w, "home.html", PageData{Title: "Power'4 Web — Accueil"})
}

func gameInitHandler(w http.ResponseWriter, r *http.Request) {
	if !mustMethod(w, r, http.MethodGet) {
		return
	}
	render(w, "game_init.html", PageData{Title: "Power'4 Web — Initialisation"})
}

func gameInitTraitementHandler(w http.ResponseWriter, r *http.Request) {
	if !mustMethod(w, r, http.MethodPost) {
		return
	}

	p1 := strings.TrimSpace(r.FormValue("player1_name"))
	p2 := strings.TrimSpace(r.FormValue("player2_name"))
	c1 := normalizeColor(r.FormValue("player1_color"))
	c2 := normalizeColor(r.FormValue("player2_color"))

	if !nameRegex.MatchString(p1) || !nameRegex.MatchString(p2) {
		redirectErr(w, r, 400, "Pseudo invalide (3-20 caractères)")
		return
	}
	if c1 == "" || c2 == "" || c1 == c2 {
		redirectErr(w, r, 400, "Couleurs invalides ou identiques")
		return
	}

	g := make([][]int, 6)
	for i := range g {
		g[i] = make([]int, 7)
	}

	currentGame = &CurrentGame{
		Player1:       Player{Name: p1, Color: c1},
		Player2:       Player{Name: p2, Color: c2},
		Started:       time.Now(),
		Turns:         0,
		CurrentPlayer: 1,
		Grid:          g,
	}

	http.Redirect(w, r, "/game/play", http.StatusSeeOther)
}

func gamePlayHandler(w http.ResponseWriter, r *http.Request) {
	if !mustMethod(w, r, http.MethodGet) {
		return
	}
	if currentGame == nil {
		redirectErr(w, r, 400, "Aucune partie en cours")
		return
	}
	render(w, "game_play.html", PageData{Title: "Power'4 Web — Partie en cours", CurrentGame: currentGame})
}

func gamePlayTraitementHandler(w http.ResponseWriter, r *http.Request) {
	if !mustMethod(w, r, http.MethodPost) {
		return
	}
	if currentGame == nil {
		redirectErr(w, r, 400, "Aucune partie en cours")
		return
	}

	col, err := strconv.Atoi(strings.TrimSpace(r.FormValue("column")))
	if err != nil || col < 1 || col > 7 {
		redirectErr(w, r, 400, "Colonne invalide")
		return
	}
	col--

	placed := false
	for row := 5; row >= 0; row-- {
		if currentGame.Grid[row][col] == 0 {
			currentGame.Grid[row][col] = currentGame.CurrentPlayer
			currentGame.Turns++
			placed = true
			break
		}
	}
	if !placed {
		redirectErr(w, r, 400, "Colonne pleine")
		return
	}

	// check win/draw
	winner := ""
	if checkWin(currentGame.Grid, currentGame.CurrentPlayer) {
		if currentGame.CurrentPlayer == 1 {
			winner = currentGame.Player1.Name
		} else {
			winner = currentGame.Player2.Name
		}
		scoreboard = append(scoreboard, GameResult{Player1: currentGame.Player1.Name, Player2: currentGame.Player2.Name, Winner: winner, Date: time.Now(), Turns: currentGame.Turns})
		http.Redirect(w, r, "/game/end", 303)
		return
	}
	if boardFull(currentGame.Grid) {
		scoreboard = append(scoreboard, GameResult{Player1: currentGame.Player1.Name, Player2: currentGame.Player2.Name, Winner: "", Date: time.Now(), Turns: currentGame.Turns})
		http.Redirect(w, r, "/game/end", 303)
		return
	}

	currentGame.CurrentPlayer = 3 - currentGame.CurrentPlayer
	http.Redirect(w, r, "/game/play", 303)
}

func gameEndHandler(w http.ResponseWriter, r *http.Request) {
	if !mustMethod(w, r, http.MethodGet) {
		return
	}
	var last *GameResult
	if len(scoreboard) > 0 {
		last = &scoreboard[len(scoreboard)-1]
	}
	render(w, "game_end.html", PageData{Title: "Power'4 Web — Fin de partie", LastResult: last})
}

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	if !mustMethod(w, r, http.MethodGet) {
		return
	}
	render(w, "scoreboard.html", PageData{Title: "Power'4 Web — Scoreboard", Scoreboard: scoreboard})
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	code, _ := strconv.Atoi(r.FormValue("code"))
	render(w, "error.html", PageData{Title: "Power'4 Web — Erreur", ErrorCode: code, ErrorMessage: r.FormValue("msg")})
}

/* ---------- main ---------- */
func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game/init", gameInitHandler)
	http.HandleFunc("/game/init/traitement", gameInitTraitementHandler)
	http.HandleFunc("/game/play", gamePlayHandler)
	http.HandleFunc("/game/play/traitement", gamePlayTraitementHandler)
	http.HandleFunc("/game/end", gameEndHandler)
	http.HandleFunc("/game/scoreboard", scoreboardHandler)
	http.HandleFunc("/error", errorHandler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
