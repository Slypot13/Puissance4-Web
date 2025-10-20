package main

import (
	"fmt"
)

const (
	rows    = 6
	columns = 7
	empty   = '.'
)

type Game
 struct {
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

		if checkVictoryLines(&game) || checkVictoryColumns(&game) || checkVictoryDiagonals(&game) {
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

		switchPlayer(&game)
	}
}

func initializeGame() Game {
	var game Game
	for i := range game.board {
		for j := range game.board[i] {
			game.board[i][j] = empty
		}
	}
	game.players[0] = 'X'
	game.players[1] = 'O'
	game.current = 0
	game.turns = 0
	game.gameOver = false
	return game
}

func displayBoard(game *Game) {
	fmt.Println("\nÉtat actuel de la grille :")
	for col := 0; col < columns; col++ {
		fmt.Printf(" %d", col)
	}
	fmt.Println()
	for _, row := range game.board {
		for _, cell := range row {
			fmt.Printf(" %c", cell)
		}
		fmt.Println()
	}
	fmt.Println()
}

func dropDisc(game *Game, col int) bool {
	if col < 0 || col >= columns {
		fmt.Println("Colonne invalide. Essayez de nouveau.")
		return false
	}
	for row := rows - 1; row >= 0; row-- {
		if game.board[row][col] == empty {
			game.board[row][col] = game.players[game.current]
			return true
		}
	}
	fmt.Println("Colonne pleine. Choisissez une autre colonne.")
	return false
}

func checkVictoryLines(game *Game) bool {
	symbol := game.players[game.current]
	for row := 0; row < rows; row++ {
		count := 0
		for col := 0; col < columns; col++ {
			if game.board[row][col] == symbol {
				count++
				if count == 4 {
					return true
				}
			} else {
				count = 0
			}
		}
	}
	return false
}

func checkVictoryColumns(game *Game) bool {
	symbol := game.players[game.current]
	for col := 0; col < columns; col++ {
		count := 0
		for row := 0; row < rows; row++ {
			if game.board[row][col] == symbol {
				count++
				if count == 4 {
					return true
				}
			} else {
				count = 0
			}
		}
	}
	return false
}

func checkVictoryDiagonals(game *Game) bool {
	symbol := game.players[game.current]
	// Diagonale descendante (\)
	for row := 0; row < rows-3; row++ {
		for col := 0; col < columns-3; col++ {
			if game.board[row][col] == symbol &&
				game.board[row+1][col+1] == symbol &&
				game.board[row+2][col+2] == symbol &&
				game.board[row+3][col+3] == symbol {
				return true
			}
		}
	}
	// Diagonale montante (/)
	for row := 3; row < rows; row++ {
		for col := 0; col < columns-3; col++ {
			if game.board[row][col] == symbol &&
				game.board[row-1][col+1] == symbol &&
				game.board[row-2][col+2] == symbol &&
				game.board[row-3][col+3] == symbol {
				return true
			}
		}
	}
	return false
}

func checkDraw(game *Game) bool {
	for col := 0; col < columns; col++ {
		if game.board[0][col] == empty {
			return false
		}
	}
	return true
}

func switchPlayer(game *Game) {
	game.current = 1 - game.current
	game.turns++
}
