package games

import (
	"encoding/json"

	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/google/uuid"
)

func New() ([]byte, error) {
	type returnVals struct {
		GameID uuid.UUID `json:"game_id"`
	}

	game := chess.Game{
		State:  chess.NewGameState(),
		Moves:  []chess.Move{},
		Result: "*",
		White:  chess.Player{Name: "White", Elo: 1600},
		Black:  chess.Player{Name: "Black", Elo: 1600},
	}

	room := gameRoom{
		game:      &game,
		broadcast: make(chan struct{}),
	}

	gameID := uuid.New()

	registry.mu.Lock()
	registry.rooms[gameID] = &room
	registry.mu.Unlock()

	data, err := json.Marshal(returnVals{gameID})
	if err != nil {
		return nil, err
	}

	go room.runBroadcastLoop()

	return data, nil
}
