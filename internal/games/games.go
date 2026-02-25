package games

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/debobrad579/chessgo/internal/chess"
)

type gameRoom struct {
	game      *chess.Game
	whiteConn *websocket.Conn
	blackConn *websocket.Conn
	mu        sync.Mutex
	broadcast chan struct{}
}

type gamesRegistry struct {
	mu    sync.Mutex
	rooms map[uuid.UUID]*gameRoom
}

var registry = gamesRegistry{
	rooms: make(map[uuid.UUID]*gameRoom),
}
