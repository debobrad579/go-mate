package games

import (
	"time"

	"github.com/debobrad579/chessgo/internal/database"
	"github.com/gorilla/websocket"
)

func (gr *GameRoom) assignRole(conn *websocket.Conn, user *database.User) PlayerRole {
	switch {
	case gr.Game.White != nil && gr.Game.White.ID == user.ID:
		gr.whiteConn = conn
		return White
	case gr.Game.Black != nil && gr.Game.Black.ID == user.ID:
		gr.blackConn = conn
		return Black
	case gr.whiteConn == nil:
		gr.Game.White = user
		gr.whiteConn = conn
		if gr.blackConn != nil {
			gr.turnStart = time.Now()
			registry.notifySubscribers()
		}
		return White
	case gr.blackConn == nil:
		gr.Game.Black = user
		gr.blackConn = conn
		if gr.whiteConn != nil {
			gr.turnStart = time.Now()
			registry.notifySubscribers()
		}
		return Black
	default:
		data, err := gr.getGameData()
		if err != nil {
			return Spectator
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.Close()
			return Spectator
		}
		gr.spectatorConns[user.ID] = conn
		return Spectator
	}
}
