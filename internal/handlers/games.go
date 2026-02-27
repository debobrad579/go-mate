package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/debobrad579/chessgo/internal/games"
)

func (cfg *Config) NewGameHandler(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := games.New(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (cfg *Config) ConnectToGameHandler(w http.ResponseWriter, r *http.Request) {
	gameIDStr := r.PathValue("gameID")

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "invalid game ID", http.StatusBadRequest)
		return
	}

	user, err := cfg.getUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	room, err := games.GetGameRoom(gameID)
	if err != nil {
		http.Error(w, "game room not found", http.StatusNotFound)
		return
	}

	conn, playerColor := room.Connect(w, r, user)
	if conn == nil {
		return
	}
	defer room.Disconnect(conn, playerColor)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message: ", err)
			return
		}

		if err = room.MakeMove(message, playerColor); err != nil {
			return
		}
	}
}

func (cfg *Config) GamesListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := games.Subscribe()
	defer games.Unsubscribe(ch)

	sendSnapshot(w, flusher)

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			sendSnapshot(w, flusher)
		}
	}
}

func sendSnapshot(w http.ResponseWriter, flusher http.Flusher) {
	data, err := json.Marshal(games.GetGamesList())
	if err != nil {
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}
