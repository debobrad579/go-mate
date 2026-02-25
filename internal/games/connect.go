package games

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/debobrad579/chessgo/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectToGame(w http.ResponseWriter, r *http.Request, gameID uuid.UUID, user *database.User) error {
	registry.mu.Lock()
	room, ok := registry.rooms[gameID]
	registry.mu.Unlock()

	if !ok {
		return errors.New("game not found")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.New("upgrade error")
	}

	room.mu.Lock()

	if room.whiteConn != nil && room.blackConn != nil {
		room.mu.Unlock()
		conn.Close()
		return errors.New("game in progress")
	}

	isWhite := false
	if room.whiteConn == nil {
		room.whiteConn = conn
		if user != nil {
			room.game.White.Name = user.Name
		}
		isWhite = true
	} else {
		room.blackConn = conn
		if user != nil {
			room.game.Black.Name = user.Name
		}
	}

	room.mu.Unlock()

	defer func() {
		conn.Close()
		room.mu.Lock()
		if isWhite {
			room.whiteConn = nil
		} else {
			room.blackConn = nil
		}

		if room.whiteConn == nil && room.blackConn == nil {
			close(room.broadcast)
			room.mu.Unlock()
			registry.mu.Lock()
			delete(registry.rooms, gameID)
			registry.mu.Unlock()
			return
		}

		room.mu.Unlock()
	}()

	room.broadcast <- struct{}{}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var move chess.Move
		if err := json.Unmarshal(message, &move); err != nil {
			continue
		}

		room.mu.Lock()

		if (isWhite && room.game.Turn() == chess.Black) || (!isWhite && room.game.Turn() == chess.White) {
			room.mu.Unlock()
			continue
		}

		if !room.game.IsMoveValid(move) {
			room.mu.Unlock()
			continue
		}

		room.game.Move(move)
		room.mu.Unlock()

		select {
		case room.broadcast <- struct{}{}:
		default:
		}
	}
}
