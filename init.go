package main

import (
	"fmt"
)

const (
	rows    = 6
	columns = 7
	empty   = '.'
)

type Game struct {
	board    [rows][columns]rune
	players  [2]rune
	current  int
	turns    int
	gameOver bool
}

func main() {
	game := initializeGame()

	for !game.gameOver {
		displayBoard(&game)
		fmt.Printf("Tour du joueur %c. Choisissez une colonne (0 à %d) : ", game.players[game.current], columns-1)

		var col int
		_, err := fmt.Scan(&col)
		if err != nil {
			fmt.Println("Entrée invalide. Veuillez entrer un numéro de colonne.")
			continue
		}

		if !dropDisc(&game, col) {
			continue
		}

		if checkVictoryLines(&game) || checkVictoryColumns(&game) {
			displayBoard(&game)
			fmt.Printf("Le joueur %c a gagné !\n", game.players[game.current])
			game.gameOver = true
			break
		}

		if checkDraw(&game) {
			displayBoard(&game)
			fmt.Println("Match nul !")
			game.gameOver = true
			break
		}

		switchPlayer(&game) // FT5
	}
}
