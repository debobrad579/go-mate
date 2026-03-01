package games

import (
	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/google/uuid"
)

type PlayerRole int

const (
	White PlayerRole = iota
	Black
	Spectator
)

type Player struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GameData struct {
	Moves       []chess.Move      `json:"moves"`
	TimeControl chess.TimeControl `json:"time_control"`
	ThinkTime   int               `json:"think_time"`
	Result      string            `json:"result"`
	White       Player            `json:"white"`
	Black       Player            `json:"black"`
}

type GameListItem struct {
	ID          uuid.UUID         `json:"id"`
	White       Player            `json:"white"`
	Black       Player            `json:"black"`
	TimeControl chess.TimeControl `json:"time_control"`
}
