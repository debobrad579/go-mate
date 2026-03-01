package games

func GetGamesList() []GameListItem {
	items := make([]GameListItem, 0, len(registry.rooms))
	for id, room := range registry.rooms {
		items = append(items, GameListItem{
			ID:          id,
			White:       room.safeGetWhite(),
			Black:       room.safeGetBlack(),
			TimeControl: room.Game.TimeControl,
		})
	}

	return items
}
