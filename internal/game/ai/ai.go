package ai

import (
	"fmt"
	"go-xox-grpc-ai/internal/game"
	"math/rand"
	"time"
)

var Difficulty string

func StartGame(difficulty string) {
	Difficulty = difficulty

	fmt.Printf("Want to be %s or %s? ", game.PLAYER_X, game.PLAYER_O)
	var player string
	for player != game.PLAYER_X && player != game.PLAYER_O {
		fmt.Print("Choice: ")
		fmt.Scanln(&player)
	}

	currentGame := game.NewGame()
	currentGame.Start()

	var aiFunction func()

	if difficulty == "1" {
		aiFunction = aiRandomPlay(currentGame)
	} else if difficulty == "2" {
		aiFunction = aiBestOrRandom(currentGame)
	} else if difficulty == "3" {
		aiFunction = aiPlay(currentGame)
	}

	for !currentGame.CheckGameFinished() {
		if currentGame.GetCurrentPlayer() == player {
			currentGame.PlayRound()
		} else {
			aiFunction()
		}
		currentGame.RenderBoard()
	}

	winner := currentGame.GetWinner()
	if winner == game.TIE {
		fmt.Println("Tie!")
	} else {
		fmt.Println("Player " + winner + " Won!")
	}

}

func aiBestOrRandom(g *game.GameBoard) func() {
	return func() {
		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(2)
		if random == 0 {
			aiRandomPlay(g)()
		} else {
			aiPlay(g)()
		}
	}
}

func aiRandomPlay(g *game.GameBoard) func() {
	return func() {
		var position int = -1
		for position == -1 {
			rand.Seed(time.Now().UnixNano())
			randomPos := rand.Intn(9)
			if g.GetBoard()[randomPos] == game.EMPTY {
				position = randomPos
			}
		}

		g.SetBoardValue(position, g.GetCurrentPlayer())
		g.SwitchCurrentPlayer()
	}
}

func aiPlay(g *game.GameBoard) func() {
	return func() {
		pos := findBestMove(g.GetBoard(), g.GetCurrentPlayer())

		g.SetBoardValue(pos, g.GetCurrentPlayer())
		g.SwitchCurrentPlayer()
	}
}

func findBestMove(board []string, currentPlayer string) int {
	bestVal := -1000
	bestMove := -1

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pos := i*3 + j%3
			if board[pos] == game.EMPTY {
				board[pos] = currentPlayer

				moveValue := miniMax(board, 0, false, currentPlayer)

				board[pos] = game.EMPTY

				if moveValue > bestVal {
					bestMove = pos
					bestVal = moveValue
				}
			}
		}
	}

	return bestMove
}

func isMovesLeft(board []string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pos := i*3 + j%3
			if board[pos] == game.EMPTY {
				return true
			}
		}
	}
	return false
}

func miniMax(board []string, depth int, isMax bool, currentPlayer string) int {
	score := heuristic(board, currentPlayer) // evaluate heuristic result

	if score == 10 || score == -10 {
		return score
	}

	if !isMovesLeft(board) {
		return 0
	}

	if isMax {
		best := -1000

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				pos := i*3 + j%3
				if board[pos] == game.EMPTY {
					board[pos] = currentPlayer
					best = max(best, miniMax(board, depth+1, !isMax, currentPlayer))
					board[pos] = game.EMPTY
				}
			}
		}
		return best
	} else {
		best := 1000

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				pos := i*3 + j%3
				if board[pos] == game.EMPTY {
					if currentPlayer == game.PLAYER_O {
						board[pos] = game.PLAYER_X
					} else {
						board[pos] = game.PLAYER_O
					}
					best = min(best, miniMax(board, depth+1, !isMax, currentPlayer))
					board[pos] = game.EMPTY
				}
			}
		}
		return best
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func heuristic(board []string, currentPlayer string) int {
	opponent := game.PLAYER_X
	if currentPlayer == game.PLAYER_X {
		opponent = game.PLAYER_O
	}

	for i := 0; i < 3; i++ {
		if board[i*3] == board[i*3+1] && board[i*3+1] == board[i*3+2] {
			if board[i*3] == currentPlayer {
				return 10
			} else if board[i*3] == opponent {
				return -10
			}
		}
	}

	for i := 0; i < 3; i++ {
		if board[i] == board[3+i] && board[3+i] == board[6+i] {
			if board[i] == currentPlayer {
				return 10
			} else if board[i] == opponent {
				return -10
			}
		}
	}

	if board[0] == board[4] && board[4] == board[8] {
		if board[0] == currentPlayer {
			return 10
		} else if board[0] == opponent {
			return -10
		}
	}

	if board[2] == board[4] && board[4] == board[6] {
		if board[2] == currentPlayer {
			return 10
		} else if board[2] == opponent {
			return -10
		}
	}

	return 0
}
