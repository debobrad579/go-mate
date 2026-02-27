package games

import (
	"github.com/debobrad579/chessgo/internal/chess"
	"github.com/google/uuid"
)

type GameListItem struct {
	ID   uuid.UUID   `json:"id"`
	Game *chess.Game `json:"game"`
}

func GetGamesList() []GameListItem {
	var items []GameListItem

	for key, value := range registry.rooms {
		items = append(items, GameListItem{ID: key, Game: value.game})
	}

	return items
}
