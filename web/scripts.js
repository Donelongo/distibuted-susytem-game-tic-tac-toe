document.addEventListener("DOMContentLoaded", () => {
    const gameBoard = document.getElementById("game-board");
    const status = document.getElementById("status");
    const playAgainButton = document.getElementById("play-again-button");
    const homeButton = document.getElementById("home-button");
    const playerXScore = document.getElementById("player-x-score");
    const playerOScore = document.getElementById("player-o-score");
    const sessionIdElement = document.getElementById("session-id");

    playAgainButton.addEventListener("click", () => {
        socket.send(JSON.stringify({ action: "reset" }));
    });

    homeButton.addEventListener("click", () => {
        window.location.href = "index.html";
    });

    let currentPlayer = "X";
    let board = Array(9).fill("");
    const sessionID = new URLSearchParams(window.location.search).get("session");
    let socket;

    if (sessionID === "new") {
        const newSessionID = Date.now().toString();
        socket = new WebSocket(`ws://localhost:8080/ws?session=${newSessionID}`);
        sessionIdElement.textContent = newSessionID;
    } else {
        socket = new WebSocket(`ws://localhost:8080/ws?session=${sessionID}`);
        sessionIdElement.textContent = sessionID;
    }

    let playerSymbol = "";
    let scores = { X: 0, O: 0 };

    socket.onopen = () => {
        socket.send(JSON.stringify({ action: "join" }));
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.action === "update") {
            board = data.board;
            currentPlayer = data.currentPlayer;
            renderBoard();
            if (data.finished) {
                if (data.winner === "-") {
                    status.textContent = "It's a tie!";
                } else if (data.winner === playerSymbol) {
                    status.textContent = "You win!";
                    scores[playerSymbol]++;
                } else {
                    status.textContent = "You lose!";
                    scores[data.winner]++;
                }
                updateScoreboard();
                gameBoard.style.pointerEvents = "none";
                playAgainButton.style.display = "block";
            }
        } else if (data.action === "reset") {
            board = Array(9).fill("");
            gameBoard.style.pointerEvents = "auto";
            status.textContent = `You are Player ${playerSymbol}`;
            playAgainButton.style.display = "none";
            renderBoard();
        } else if (data.action === "assign") {
            playerSymbol = data.symbol;
            board = data.board;
            currentPlayer = data.currentPlayer;
            status.textContent = `You are Player ${playerSymbol}`;
            renderBoard();
        } else if (data.action === "start") {
            currentPlayer = data.currentPlayer;
            if (currentPlayer === playerSymbol) {
                status.textContent = `You are Player ${playerSymbol} and you start the game!`;
            } else {
                status.textContent = `You are Player ${playerSymbol}. Player ${currentPlayer} starts the game.`;
            }
            renderBoard();
        }
    };

    function renderBoard() {
        gameBoard.innerHTML = "";
        board.forEach((cell, index) => {
            const cellElement = document.createElement("div");
            cellElement.classList.add("cell");
            if (cell === "X") {
                cellElement.classList.add("x");
            } else if (cell === "O") {
                cellElement.classList.add("o");
            }
            cellElement.textContent = cell;
            cellElement.addEventListener("click", () => handleCellClick(index));
            gameBoard.appendChild(cellElement);
        });
    }

    function handleCellClick(index) {
        if (board[index] === "" && currentPlayer === playerSymbol) {
            board[index] = currentPlayer;
            currentPlayer = currentPlayer === "X" ? "O" : "X";
            status.textContent = `You are Player ${playerSymbol}`;
            renderBoard();
            checkGameStatus();
        }
    }

    function checkGameStatus() {
        fetch("/api/check", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ board })
        })
        .then(response => response.json())
        .then(data => {
            socket.send(JSON.stringify({
                action: "update",
                board: board,
                currentPlayer: currentPlayer,
                finished: data.finished,
                winner: data.winner
            }));
        });
    }

    function updateScoreboard() {
        playerXScore.textContent = `Player X: ${scores.X}`;
        playerOScore.textContent = `Player O: ${scores.O}`;
    }

    renderBoard();
});
