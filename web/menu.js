document.getElementById("new-game-button").addEventListener("click", () => {
    window.location.href = "game.html?session=new";
});

document.getElementById("continue-game-button").addEventListener("click", () => {
    const sessionID = prompt("Enter your session ID:");
    if (sessionID) {
        window.location.href = `game.html?session=${sessionID}`;
    }
});
