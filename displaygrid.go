package main

import "fmt"

func displayBoard(game *Game) {
	fmt.Println("\nEtat actuel de la grille :")
	for col := 0; col < columns; col++ {
		fmt.Printf("%d ", col)
	}
	fmt.Println()
	for _, row := range game.board {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Println()
	}
	fmt.Println()
}
