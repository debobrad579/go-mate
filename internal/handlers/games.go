package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/debobrad579/chessgo/internal/games"
)

func (cfg *Config) NewGameHandler(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	games.ConnectToGame(w, r, gameID, user)
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
