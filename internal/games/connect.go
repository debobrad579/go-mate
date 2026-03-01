package games

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/debobrad579/chessgo/internal/database"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (gr *GameRoom) Connect(w http.ResponseWriter, r *http.Request, user *database.User) (*websocket.Conn, chess.Color) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade websocket", http.StatusBadRequest)
		return nil, chess.White
	}

	gr.mu.Lock()

	if gr.whiteConn != nil && gr.blackConn != nil {
		gr.mu.Unlock()
		return nil, chess.White
		// TODO: Add spectating
	}

	playerColor := chess.White
	if gr.Game.White != nil && gr.Game.White.ID == user.ID {
		gr.whiteConn = conn
	} else if gr.Game.Black != nil && gr.Game.Black.ID == user.ID {
		gr.blackConn = conn
		playerColor = chess.Black
	} else if gr.whiteConn == nil {
		gr.Game.White = user
		gr.whiteConn = conn
		if gr.blackConn != nil {
			gr.turnStart = time.Now()
			registry.notifySubscribers()
		}
	} else if gr.blackConn == nil {
		gr.Game.Black = user
		gr.blackConn = conn
		if gr.whiteConn != nil {
			gr.turnStart = time.Now()
			registry.notifySubscribers()
		}
		playerColor = chess.Black
	}

	if gr.Game.White != nil && gr.Game.Black != nil {
		gr.ThinkTime = int(time.Since(gr.turnStart).Milliseconds())
	} else {
		gr.ThinkTime = 0
	}

	gr.mu.Unlock()

	gr.broadcast <- struct{}{}

	return conn, playerColor
}

func (gr *GameRoom) Disconnect(conn *websocket.Conn, playerColor chess.Color) {
	conn.Close()
	gr.mu.Lock()

	if playerColor == chess.White {
		gr.whiteConn = nil
	} else {
		gr.blackConn = nil
	}

	if gr.whiteConn == nil && gr.blackConn == nil {
		close(gr.broadcast)
		gr.mu.Unlock()
		registry.mu.Lock()
		delete(registry.rooms, gr.ID)
		registry.mu.Unlock()
		registry.notifySubscribers()
		return
	}

	gr.mu.Unlock()
}
