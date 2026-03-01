package games

import (
	"github.com/gorilla/websocket"
)

func (gr *GameRoom) runBroadcastLoop() {
	for range gr.broadcast {
		gr.mu.Lock()

		data, err := gr.getGameData()
		if err != nil {
			gr.mu.Unlock()
			continue
		}

		white := gr.whiteConn
		black := gr.blackConn

		gr.mu.Unlock()

		if white != nil {
			if err := white.WriteMessage(websocket.TextMessage, data); err != nil {
				white.Close()
				gr.mu.Lock()
				if gr.whiteConn == white {
					gr.whiteConn = nil
				}
				gr.mu.Unlock()
			}
		}

		if black != nil {
			if err := black.WriteMessage(websocket.TextMessage, data); err != nil {
				black.Close()
				gr.mu.Lock()
				if gr.blackConn == black {
					gr.blackConn = nil
				}
				gr.mu.Unlock()
			}
		}

		for _, conn := range gr.spectatorConns {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				conn.Close()
			}
		}
	}
}
