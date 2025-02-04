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
	Board           []string
	CurrentPlayer   string
	Finished        bool
	Winner          string
	PreviousStarter string
	Scores          map[string]int
	Clients         []*websocket.Conn
	PlayerSymbols   map[*websocket.Conn]string
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
				Board:           make([]string, 9),
				CurrentPlayer:   "X",
				PreviousStarter: "X",
				Scores:          map[string]int{"X": 0, "O": 0},
				Clients:         []*websocket.Conn{},
				PlayerSymbols:   make(map[*websocket.Conn]string),
			}
			sessions[sessionID] = session
		}
		session.mu.Lock()
		session.Clients = append(session.Clients, conn)
		if len(session.PlayerSymbols) == 0 {
			session.PlayerSymbols[conn] = "X"
		} else if len(session.PlayerSymbols) == 1 {
			session.PlayerSymbols[conn] = "O"
		}
		session.mu.Unlock()
		sessionsMu.Unlock()

		conn.WriteJSON(map[string]interface{}{
			"action":        "assign",
			"symbol":        session.PlayerSymbols[conn],
			"board":         session.Board,
			"currentPlayer": session.CurrentPlayer,
			"finished":      session.Finished,
			"winner":        session.Winner,
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
				session.Board = make([]string, 9)
				if session.Winner == "-" {
					if session.PreviousStarter == "X" {
						session.CurrentPlayer = "O"
					} else {
						session.CurrentPlayer = "X"
					}
				} else {
					session.CurrentPlayer = session.Winner
				}
				session.PreviousStarter = session.CurrentPlayer
				session.Finished = false
				session.Winner = ""
				for client := range session.PlayerSymbols {
					session.PlayerSymbols[client] = ""
				}
				for _, client := range session.Clients {
					client.WriteJSON(map[string]string{
						"action":        "start",
						"currentPlayer": session.CurrentPlayer,
					})
				}
			} else if msg.Action == "update" {
				session.Board = msg.Board
				session.CurrentPlayer = msg.CurrentPlayer
				session.Finished = msg.Finished
				session.Winner = msg.Winner
				if msg.Winner != "-" {
					session.Scores[msg.Winner]++
				}
			}
			for _, client := range session.Clients {
				client.WriteJSON(msg)
			}
			session.mu.Unlock()
		}
		conn.Close()
	})

	http.ListenAndServe(":8080", nil)
}

//! the issue that i have to fix is that after recoveringa  and if player one is the one on turn u can see the game updating for player 2 as well but player two can play.
