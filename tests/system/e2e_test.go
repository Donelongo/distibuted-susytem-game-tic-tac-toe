package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// MockWebSocketHandler simulates WebSocket interactions for testing
func MockWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a WebSocket connection
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins for testing
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Unable to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Simulate the initial message from the server (e.g., player assignment)
	conn.WriteJSON(map[string]interface{}{"action": "assign"})

	// Simulate player moves
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// Handle player move and update board
		if msg["action"] == "update" {
			// Mock board update
			conn.WriteJSON(map[string]interface{}{
				"action": "update",
				"board":  msg["board"], // Mock updating the board
			})
		}
	}
}

func TestEndToEndGameFlow(t *testing.T) {
	// Start a mock WebSocket server
	server := httptest.NewServer(http.HandlerFunc(MockWebSocketHandler))
	defer server.Close()

	// Extract WebSocket URL
	wsURL := "ws" + server.URL[4:] + "/ws?session=test123"

	// Connect first player (Player 1)
	player1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer player1.Close()

	// Connect second player (Player 2)
	player2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer player2.Close()

	// Read initial assignments
	var msg map[string]interface{}
	err = player1.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, "assign", msg["action"])

	err = player2.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, "assign", msg["action"])

	// Player 1 makes a move
	move := map[string]interface{}{
		"action": "update",
		"board":  []string{"X", "", "", "", "", "", "", "", ""},
	}
	err = player1.WriteJSON(move)
	assert.NoError(t, err)

	// Verify board update on Player 2's side
	err = player2.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, "update", msg["action"])
	assert.Equal(t, "X", msg["board"].([]interface{})[0])

	// Simulate more moves
	move2 := map[string]interface{}{
		"action": "update",
		"board":  []string{"X", "O", "", "", "", "", "", "", ""},
	}
	err = player2.WriteJSON(move2)
	assert.NoError(t, err)

	err = player1.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, "O", msg["board"].([]interface{})[1])

	// Simulate game completion
	finalMove := map[string]interface{}{
		"action": "update",
		"board":  []string{"X", "O", "X", "O", "X", "O", "X", "", ""},
	}
	err = player1.WriteJSON(finalMove)
	assert.NoError(t, err)

	err = player2.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, true, msg["finished"])

	t.Log("End-to-end game flow test passed.")
}
