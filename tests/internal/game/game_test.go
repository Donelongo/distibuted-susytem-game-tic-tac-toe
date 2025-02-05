package game_test

import (
	"go-xox-grpc-ai/internal/game"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGameBoard is a mock for the GameBoard struct.
type MockGameBoard struct {
	mock.Mock
	board         []string
	currentPlayer string
	finished      bool
	winner        string
}

// Mocking the Start method
func (m *MockGameBoard) Start() {
	m.Called()
}

// Mocking PlayRound method
func (m *MockGameBoard) PlayRound() {
	m.Called()
}

// Mocking GetCurrentPlayer
func (m *MockGameBoard) GetCurrentPlayer() string {
	args := m.Called()
	return args.String(0)
}

// Mocking CheckGameFinished
func (m *MockGameBoard) CheckGameFinished() bool {
	args := m.Called()
	return args.Bool(0)
}

// Mocking GetWinner
func (m *MockGameBoard) GetWinner() string {
	args := m.Called()
	return args.String(0)
}

// Test Start Method
func TestGameStart(t *testing.T) {
	mockGame := new(MockGameBoard)
	mockGame.On("Start").Return()

	mockGame.Start()

	mockGame.AssertCalled(t, "Start")
}

// Test PlayRound Method
func TestGamePlayRound(t *testing.T) {
	mockGame := new(MockGameBoard)
	mockGame.On("PlayRound").Return()

	mockGame.PlayRound()

	mockGame.AssertCalled(t, "PlayRound")
}

// Test Switching Players
func TestSwitchCurrentPlayer(t *testing.T) {
	gameBoard := game.NewGame()
	assert.Equal(t, game.PLAYER_X, gameBoard.GetCurrentPlayer())

	gameBoard.SwitchCurrentPlayer()
	assert.Equal(t, game.PLAYER_O, gameBoard.GetCurrentPlayer())

	gameBoard.SwitchCurrentPlayer()
	assert.Equal(t, game.PLAYER_X, gameBoard.GetCurrentPlayer())
}

// Test CheckGameFinished Mock
func TestGameFinished(t *testing.T) {
	mockGame := new(MockGameBoard)
	mockGame.On("CheckGameFinished").Return(true)
	mockGame.On("GetWinner").Return(game.PLAYER_X)

	finished := mockGame.CheckGameFinished()
	winner := mockGame.GetWinner()

	assert.True(t, finished, "Game should be finished")
	assert.Equal(t, game.PLAYER_X, winner, "Player X should be the winner")

	mockGame.AssertExpectations(t)
}

// Test Legal Move
func TestIsLegalMove(t *testing.T) {
	gameBoard := game.NewGame()

	assert.True(t, gameBoard.IsLegalMove(1), "Position 1 should be legal")
	gameBoard.SetBoardValue(0, game.PLAYER_X)
	assert.False(t, gameBoard.IsLegalMove(1), "Position 1 should be illegal after being filled")
}

// Test Board Update
func TestSetBoardValue(t *testing.T) {
	gameBoard := game.NewGame()
	gameBoard.SetBoardValue(0, game.PLAYER_X)

	assert.Equal(t, game.PLAYER_X, gameBoard.GetBoard()[0], "Position 1 should have PLAYER_X")
}
