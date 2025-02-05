package online_test

import (
	"go-xox-grpc-ai/internal/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullGameAIvsAI(t *testing.T) {
	gameBoard := game.NewGame()
	gameBoard.Start()

	currentPlayer := game.PLAYER_X
	moveCount := 0

	// Simulate moves until the board is full
	for moveCount < 9 {
		for i := range gameBoard.GetBoard() {
			if gameBoard.GetBoard()[i] == game.EMPTY {
				gameBoard.GetBoard()[i] = currentPlayer // Simulating move
				break
			}
		}
		moveCount++

		// Alternate players
		if currentPlayer == game.PLAYER_X {
			currentPlayer = game.PLAYER_O
		} else {
			currentPlayer = game.PLAYER_X
		}
	}

	// Debugging: Print the final board state
	t.Logf("Final Board State: %+v", gameBoard.GetBoard())

	result := gameBoard.GetWinner()
	t.Logf("Winner Detected: %s", result)

	validResults := []string{game.PLAYER_X, game.PLAYER_O}

	assert.Contains(t, validResults, result, "Game should end with a win")
}
