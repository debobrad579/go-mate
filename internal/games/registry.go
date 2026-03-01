package games

import (
	"errors"
	"sync"
	"time"

	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameRoom struct {
	ID             uuid.UUID   `json:"id"`
	Game           *chess.Game `json:"game"`
	whiteConn      *websocket.Conn
	blackConn      *websocket.Conn
	whiteTime      int
	blackTime      int
	mu             sync.Mutex
	broadcast      chan struct{}
	turnStart      time.Time
	spectatorConns map[uuid.UUID]*websocket.Conn
}

func GetGameRoom(gameID uuid.UUID) (*GameRoom, error) {
	registry.mu.Lock()
	room, ok := registry.rooms[gameID]
	registry.mu.Unlock()

	if !ok {
		return nil, errors.New("game room not found")
	}

	return room, nil
}

type gamesRegistry struct {
	mu          sync.Mutex
	rooms       map[uuid.UUID]*GameRoom
	subscribers map[chan struct{}]struct{}
}

func (reg *gamesRegistry) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	reg.mu.Lock()
	reg.subscribers[ch] = struct{}{}
	reg.mu.Unlock()
	return ch
}

func (reg *gamesRegistry) unsubscribe(ch chan struct{}) {
	reg.mu.Lock()
	delete(reg.subscribers, ch)
	reg.mu.Unlock()
}

func (reg *gamesRegistry) notifySubscribers() {
	reg.mu.Lock()
	for ch := range reg.subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
	reg.mu.Unlock()
}

var registry = gamesRegistry{
	rooms:       make(map[uuid.UUID]*GameRoom),
	subscribers: make(map[chan struct{}]struct{}),
}

func Subscribe() chan struct{} {
	return registry.subscribe()
}

func Unsubscribe(ch chan struct{}) {
	registry.unsubscribe(ch)
}
