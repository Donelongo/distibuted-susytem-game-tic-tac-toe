# Golang TicTacToe Game - With AI or Multiplayer(gRPC)

---

## About the project

The main purpose of this project is to learn gRPC with Golang by practicing.

I implemented a TicTacToe game in Golang. It has a simple CLI interface and a gRPC server. The gRPC server can be used to play the game with multiple players. As an extra, the game also has an AI that can play against you.

---

## Getting Started

#### Prerequisites

- Golang 1.19 or higher

#### How to run

1. Clone the repository
2. Run `go run cmd/main.go`
3. Select `Play against computer` or `Play against your friend`

![](https://i.ibb.co/1MmdR36/Screenshot-2023-02-15-at-01-53-00.png)

---

## Play Against Computer - Artificial Intelligence (AI) Part

The AI is implemented using the Minimax algorithm.

When you choice `Play against computer`, you have to select difficulty for AI.

1. `Easy`: AI will play randomly

   ![](https://s3.gifyu.com/images/ezgif.com-video-to-gif5.gif)
2. `Medium`: AI will play randomly or use minimax algorithm
3. `Hard`: AI will use minimax algorithm

   ![](https://s3.gifyu.com/images/ezgif.com-video-to-gif6.gif)

---

## Play Against Your Friend - gRPC Part

When you choice `Play against your friend`, you can host a game or join a game.

1. `Host`: You will be the host of the game. You have to wait for your friend to join the game.
2. `Join`: You will be the guest of the game. You have to enter the host's IP address to join the game.

![](https://s9.gifyu.com/images/ezgif.com-video-to-gif4.gif)

---

## Issues

- [ ] The game doesn't work properly when you play against your friend. The game ends when someone wins but you can not do anything in that point in cli, you need to interrupt and run again. `Play Again` or `Return Menu` options may be added.

---

## License

Distributed under the GPL License. See `LICENSE` for more information.

---

## Contribution

Any contributions you make are greatly appreciated.
