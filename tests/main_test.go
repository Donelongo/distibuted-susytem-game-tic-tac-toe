package main

import (
	"bytes"
	"encoding/json"
	"go-xox-grpc-ai/internal/game"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestCheckGameStatus tests the game status API
func TestCheckGameStatus(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"board": []string{"X", "O", "X", "O", "X", "O", "X", "O", "X"},
	})

	req := httptest.NewRequest("POST", "/api/check", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Board []string `json:"board"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		g := game.NewGame()
		g.SetBoard(req.Board)
		finished := g.CheckGameFinished()
		winner := g.GetWinner()

		resp := struct {
			Finished bool   `json:"finished"`
			Winner   string `json:"winner"`
		}{
			Finished: finished,
			Winner:   winner,
		}
		json.NewEncoder(w).Encode(resp)
	})

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Finished bool   `json:"finished"`
		Winner   string `json:"winner"`
	}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.True(t, resp.Finished)
	assert.Equal(t, game.PLAYER_X, resp.Winner)
}

// TestWebSocketConnection tests a WebSocket connection
func TestWebSocketConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil) // No extra headers
		if err != nil {
			http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
			return
		}
		defer conn.Close()

		msg := map[string]string{
			"action": "assign",
			"symbol": "X",
		}
		err = conn.WriteJSON(msg)
		assert.NoError(t, err)
	}))

	defer server.Close()

	url := "ws" + server.URL[len("http"):] + "/ws?session=testsession"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	var response map[string]string
	err = ws.ReadJSON(&response)
	assert.NoError(t, err)
	assert.Equal(t, "assign", response["action"])
	assert.Equal(t, "X", response["symbol"])
}
