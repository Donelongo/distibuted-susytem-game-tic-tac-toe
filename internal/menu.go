package internal

import (
	"fmt"
	"go-xox-grpc-ai/internal/game/ai"
	"go-xox-grpc-ai/internal/game/online"
	"go-xox-grpc-ai/internal/utils"
)

func MainMenu() {
	fmt.Println("Welcome to Tic Tac Toe!")
	fmt.Println("Please select your choice:")
	fmt.Println("1) Play against computer")
	fmt.Println("2) Play against your friend")
	var choice = utils.GetUserInput("1", "2")

	if choice == "1" {
		playAgainstComputer()
	} else if choice == "2" {
		playAgainstPlayer()
	}
}

func playAgainstComputer() {
	fmt.Println()
	fmt.Println("1) Easy")
	fmt.Println("2) Medium")
	fmt.Println("3) Hard")
	var choice = utils.GetUserInput("1", "2", "3")

	ai.StartGame(choice)
}

func playAgainstPlayer() {
	fmt.Println()
	fmt.Println("1) Host a game")
	fmt.Println("2) Join a game")
	var choice = utils.GetUserInput("1", "2")

	if choice == "1" {
		server := online.NewServer()
		server.HostGame()
	} else if choice == "2" {
		var host string
		fmt.Print("Enter host ip:port to join the game: ")
		fmt.Scanln(&host)

		client := online.NewClient()
		client.JoinGame(host)
	}
}
