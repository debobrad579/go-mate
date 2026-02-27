package games

import (
	"maps"
	"slices"
)

func GetGamesList() []*GameRoom {
	return slices.Collect(maps.Values(registry.rooms))
}
