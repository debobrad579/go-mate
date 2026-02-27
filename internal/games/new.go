package games

import (
	"encoding/json"

	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/debobrad579/chessgo/internal/database"
	"github.com/google/uuid"
)

func New(user *database.User) ([]byte, error) {
	type returnVals struct {
		GameID uuid.UUID `json:"game_id"`
	}

	game := chess.Game{
		State:  chess.NewGameState(),
		Moves:  []chess.Move{},
		Result: "*",
		TimeControl: chess.TimeControl{
			Base:      3 * 60 * 1000,
			Increment: 2 * 1000,
		},
		White: user,
	}

	room := GameRoom{
		ID:        uuid.New(),
		Game:      &game,
		broadcast: make(chan struct{}),
		whiteTime: game.TimeControl.Base,
		blackTime: game.TimeControl.Base,
	}

	data, err := json.Marshal(returnVals{room.ID})
	if err != nil {
		return nil, err
	}

	registry.mu.Lock()
	registry.rooms[room.ID] = &room
	registry.mu.Unlock()

	registry.notifySubscribers()

	go room.runBroadcastLoop()

	return data, nil
}
