package games

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/debobrad579/chessgo/internal/database"
)

func New(user *database.User, color chess.Color, timeControl chess.TimeControl) ([]byte, error) {
	type returnVals struct {
		GameID uuid.UUID `json:"game_id"`
	}

	game := chess.Game{
		State:       chess.NewGameState(),
		Moves:       []chess.Move{},
		Result:      "*",
		TimeControl: timeControl,
	}

	if color == chess.White {
		game.White = user
	} else {
		game.Black = user
	}

	room := GameRoom{
		ID:             uuid.New(),
		Game:           &game,
		broadcast:      make(chan struct{}),
		whiteTime:      game.TimeControl.Base,
		blackTime:      game.TimeControl.Base,
		spectatorConns: make(map[uuid.UUID]*websocket.Conn),
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
