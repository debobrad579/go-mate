package games

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/debobrad579/chessgo/internal/database"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (gr *GameRoom) Connect(w http.ResponseWriter, r *http.Request, user *database.User) (*websocket.Conn, PlayerRole) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade websocket", http.StatusBadRequest)
		return nil, Spectator
	}

	gr.mu.Lock()
	defer gr.mu.Unlock()

	playerRole := gr.assignRole(conn, user)

	if playerRole != Spectator {
		if gr.Game.White != nil && gr.Game.Black != nil {
		}

		select {
		case gr.broadcast <- struct{}{}:
		default:
		}
	}

	return conn, playerRole
}

func (gr *GameRoom) Disconnect(user *database.User) {
	gr.mu.Lock()

	switch {
	case gr.spectatorConns[user.ID] != nil:
		gr.spectatorConns[user.ID].Close()
		delete(gr.spectatorConns, user.ID)
	case gr.Game.White != nil && gr.whiteConn != nil && user.ID == gr.Game.White.ID:
		gr.whiteConn.Close()
		gr.whiteConn = nil
	case gr.Game.Black != nil && gr.blackConn != nil && user.ID == gr.Game.Black.ID:
		gr.blackConn.Close()
		gr.blackConn = nil
	}

	roomEmpty := gr.whiteConn == nil && gr.blackConn == nil

	if roomEmpty {
		close(gr.broadcast)
	}

	gr.mu.Unlock()

	if roomEmpty {
		registry.mu.Lock()
		delete(registry.rooms, gr.ID)
		registry.mu.Unlock()
		registry.notifySubscribers()
	}
}
