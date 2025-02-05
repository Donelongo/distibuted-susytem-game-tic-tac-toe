package main

import (
	"encoding/json"
	"go-xox-grpc-ai/internal/game"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GameSession struct {
	Game            *game.GameBoard
	Clients         []*websocket.Conn
	PlayerSymbols   map[*websocket.Conn]string
	Disconnected    map[string]string // Track disconnected clients by their symbol
	mu              sync.Mutex
}

var sessions = make(map[string]*GameSession)
var sessionsMu sync.Mutex

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/check", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		sessionID := r.URL.Query().Get("session")
		if sessionID == "" {
			conn.Close()
			return
		}

		sessionsMu.Lock()
		session, exists := sessions[sessionID]
		if !exists {
			session = &GameSession{
				Game:          game.NewGame(),
				Clients:       []*websocket.Conn{},
				PlayerSymbols: make(map[*websocket.Conn]string),
				Disconnected:  make(map[string]string),
			}
			sessions[sessionID] = session
		}
		session.mu.Lock()
		session.Clients = append(session.Clients, conn)
		if len(session.PlayerSymbols) == 0 {
			session.PlayerSymbols[conn] = "X"
		} else if len(session.PlayerSymbols) == 1 {
			session.PlayerSymbols[conn] = "O"
		} else {
			// Reassign symbol if reconnecting
			for symbol, _ := range session.Disconnected {
				session.PlayerSymbols[conn] = symbol
				delete(session.Disconnected, symbol)
				break
			}
		}
		session.mu.Unlock()
		sessionsMu.Unlock()

		conn.WriteJSON(map[string]interface{}{
			"action":        "assign",
			"symbol":        session.PlayerSymbols[conn],
			"board":         session.Game.GetBoard(),
			"currentPlayer": session.Game.GetCurrentPlayer(),
			"finished":      session.Game.IsFinished(),
			"winner":        session.Game.GetWinner(),
		})

		for {
			var msg struct {
				Action        string   `json:"action"`
				Symbol        string   `json:"symbol"`
				Board         []string `json:"board"`
				CurrentPlayer string   `json:"currentPlayer"`
				Finished      bool     `json:"finished"`
				Winner        string   `json:"winner"`
			}
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			session.mu.Lock()
			if msg.Action == "reset" {
				session.Game = game.NewGame()
				for client := range session.PlayerSymbols {
					session.PlayerSymbols[client] = ""
				}
				for _, client := range session.Clients {
					client.WriteJSON(map[string]string{
						"action":        "start",
						"currentPlayer": session.Game.GetCurrentPlayer(),
					})
				}
			} else if msg.Action == "update" {
				session.Game.SetBoard(msg.Board)
				session.Game.SwitchCurrentPlayer()
				if session.Game.CheckGameFinished() {
					session.Game.GetWinner()
				}
			}
			for _, client := range session.Clients {
				client.WriteJSON(msg)
			}
			session.mu.Unlock()
		}

		// Handle disconnection
		session.mu.Lock()
		symbol := session.PlayerSymbols[conn]
		session.Disconnected[symbol] = symbol
		delete(session.PlayerSymbols, conn)
		session.mu.Unlock()

		conn.Close()
	})

	http.ListenAndServe(":8080", nil)
}
