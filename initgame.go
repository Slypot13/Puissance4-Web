package main

func initialiezGame() Game {
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
