# Tic Tac Toe

Welcome to the Tic Tac Toe game! This project allows you to play Tic Tac Toe against a friend locally.

## Features

- Play against a friend locally.
- Continue an existing game using a session ID.
- Navigate back to the home page using the "Home" button.
- Multiple game sessions can be conducted simultaneously.

## Getting Started

### Prerequisites

- Go 1.16 or later
- A web browser

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/tic-tac-toe.git
    cd tic-tac-toe
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

### Running the Game

1. Start the Go server:
    ```sh
    go run cmd/main.go
    ```

2. Open your web browser and navigate to `http://localhost:8080/`.

### How to Play

#### Main Menu

When you open the game in your browser, you will see the main menu with the following options:

- **New Game**: Start a new game session.
- **Continue Game**: Continue an existing game by entering the session ID.

#### New Game

1. Click the "New Game" button.
2. You will be redirected to the game page, and a new session ID will be generated and displayed on the page.
3. Play the game by clicking on the cells to place your symbol (X or O).

#### Continue Game

1. Click the "Continue Game" button.
2. Enter the session ID of the game you want to continue.
3. You will be redirected to the game page, and the game will continue from where it left off.

#### Home Button

- On the game page, you will see a "Home" button at the top left corner.
- Click the "Home" button to navigate back to the main menu.

### Project Structure

- `cmd/main.go`: The main entry point for the Go server.
- `internal/`: Contains the game logic and utilities.
- `web/`: Contains the HTML, CSS, and JavaScript files for the web interface.

### Contributing

Contributions are welcome! Please fork the repository and submit a pull request.
