# Golang TicTacToe Game - Multiplayer(WebSocket)

---

## About the project

The main purpose of this project is to learn WebSocket and Golang by practicing.

I implemented a TicTacToe game in Golang. It has a simple web interface and a WebSocket server. The WebSocket server can be used to play the game with multiple players.

---

## Getting Started

#### Prerequisites

- Golang 1.19 or higher

#### How to run

1. Clone the repository
2. Run `go run cmd/main.go`

#### Running Multiple Clients

To run multiple clients, use different session IDs in the URL:

1. Open `http://localhost:8080/?session=session1` in one browser tab
2. Open `http://localhost:8080/?session=session2` in another browser tab

Each session ID represents a different game session.

---

## Play Against Your Friend - WebSocket Part

When you choose `Play against your friend`, you can host a game or join a game.

1. `Host`: You will be the host of the game. You have to wait for your friend to join the game.
2. `Join`: You will be the guest of the game. You have to enter the host's IP address to join the game.

---

## Issues

- [ ] The game doesn't work properly when you play against your friend. The game ends when someone wins but you can not do anything at that point in the web interface, you need to refresh the page. `Play Again` or `Return Menu` options may be added.

---

## License

Distributed under the GPL License. See `LICENSE` for more information.

---

## Contribution

Any contributions you make are greatly appreciated.
