package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/debobrad579/chessgo/internal/chess"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan *chess.Game, 256)

var (
	mu   sync.Mutex
	game *chess.Game
)

func init() {
	game = &chess.Game{
		State:  chess.NewGameState(),
		Moves:  []chess.Move{},
		Result: "*",
		White:  chess.Player{Name: "Brady DeBoer", Elo: 1733},
		Black:  chess.Player{Name: "Lee Hendon", Elo: 1735},
	}
}

func WebsocketsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	clients[conn] = true
	defer func() {
		conn.Close()
		delete(clients, conn)
	}()

	data, err := json.Marshal(game)
	if err != nil {
		fmt.Println("Error encoding initial game:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		fmt.Println("Error sending initial game state:", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			sendData(conn, game)
			break
		}

		var move chess.Move
		if err = json.Unmarshal(message, &move); err != nil {
			sendData(conn, game)
			break
		}

		if !game.IsMoveValid(move) {
			fmt.Println("Invalid move:", move)
			sendData(conn, game)
			continue
		}

		game.Move(move)

		broadcast <- game
	}
}

func sendData(conn *websocket.Conn, game *chess.Game) error {
	data, err := json.Marshal(game)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

func HandleBroadcasts() {
	for {
		game := <-broadcast

		data, err := json.Marshal(game)
		if err != nil {
			continue
		}

		for conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				conn.Close()
				delete(clients, conn)
			}
		}
	}
}
