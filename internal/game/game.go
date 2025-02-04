package game

import (
	"fmt"
	"go-xox-grpc-ai/internal/utils"
	"strconv"
)

var PLAYER_X = "X"
var PLAYER_O = "O"
var EMPTY = ""
var TIE = "-"

type GameBoard struct {
	board         []string
	currentPlayer string
	finished      bool
	winner        string
}

func NewGame() *GameBoard {
	return &GameBoard{
		board:         make([]string, 9),
		currentPlayer: PLAYER_X,
	}
}

func (g *GameBoard) Start() {
	fmt.Println("==========================")
	fmt.Println("====== Game Started ======")
	fmt.Println("==========================")
	fmt.Println()
	g.RenderBoard()
}

func (g *GameBoard) RenderBoard() {
	fmt.Println()
	fmt.Printf(" %s | %s | %s \n", g.renderCell(0), g.renderCell(1), g.renderCell(2))
	fmt.Println("----------------------")
	fmt.Printf(" %s | %s | %s \n", g.renderCell(3), g.renderCell(4), g.renderCell(5))
	fmt.Println("----------------------")
	fmt.Printf(" %s | %s | %s \n", g.renderCell(6), g.renderCell(7), g.renderCell(8))
	fmt.Println()
	fmt.Printf("==========================")
	fmt.Println()
}

func (g *GameBoard) renderCell(position int) string {
	if g.board[position] == EMPTY {
		return fmt.Sprintf("[%d]", position+1)
	} else {
		return " " + g.board[position] + " "
	}
}

func (g *GameBoard) PlayRound() {
	var position int
	for position == 0 {
		fmt.Print("Position: ")
		var _choice = utils.GetUserInput("1", "2", "3", "4", "5", "6", "7", "8", "9")
		pos, err := strconv.Atoi(_choice)
		if err != nil || pos < 1 || pos > 9 {
			continue
		}
		if g.board[pos-1] == EMPTY {
			position = pos
		}
	}

	g.board[position-1] = g.currentPlayer
	g.SwitchCurrentPlayer()
}

func (g *GameBoard) GetCurrentPlayer() string {
	return g.currentPlayer
}

func (g *GameBoard) SwitchCurrentPlayer() {
	if g.currentPlayer == PLAYER_X {
		g.currentPlayer = PLAYER_O
	} else {
		g.currentPlayer = PLAYER_X
	}
}

func (g *GameBoard) GetBoard() []string {
	return g.board
}

func (g *GameBoard) SetBoard(newBoard []string) {
	g.board = newBoard
}

// position: from 1 to 9 == [1,9]
func (g *GameBoard) IsLegalMove(position int) bool {
	if position < 1 || position > 9 {
		return false
	}
	return g.board[position-1] == EMPTY
}

func (g *GameBoard) SetBoardValue(pos int, value string) {
	g.board[pos] = value
}

var cellPoints = []int{8, 1, 6, 3, 5, 7, 4, 9, 2}

func (g *GameBoard) CheckGameFinished() bool {
	filledCells := 0
	x_points := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	o_points := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	for pos, val := range g.board {
		if val == PLAYER_X {
			x_points[pos/3][pos%3] = cellPoints[pos]
			filledCells++
		} else if val == PLAYER_O {
			o_points[pos/3][pos%3] = cellPoints[pos]
			filledCells++
		}
	}

	if checkBoardPoints(x_points) {
		g.finished = true
		g.winner = PLAYER_X
		return true
	} else if checkBoardPoints(o_points) {
		g.finished = true
		g.winner = PLAYER_O
		return true
	} else if filledCells == 9 {
		g.finished = true
		g.winner = TIE
		return true
	}

	return false
}

func checkBoardPoints(p [][]int) bool {
	for i := 0; i < 3; i++ {
		if p[i][0]+p[i][1]+p[i][2] >= 15 {
			return true
		}
	}

	for i := 0; i < 3; i++ {
		if p[0][i]+p[1][i]+p[2][i] >= 15 {
			return true
		}
	}

	if p[0][0]+p[1][1]+p[2][2] >= 15 {
		return true
	}
	if p[0][2]+p[1][1]+p[2][0] >= 15 {
		return true
	}
	return false
}

func (g *GameBoard) GetWinner() string {
	return g.winner
}

func (g *GameBoard) IsFinished() bool {
	return g.finished
}
