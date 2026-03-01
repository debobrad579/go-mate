package games

import (
	"encoding/json"
	"time"
)

func (gr *GameRoom) getThinkTime() int {
	if gr.Game.White == nil || gr.Game.Black == nil {
		return 0
	}
	return int(time.Since(gr.turnStart).Milliseconds())
}

func (gr *GameRoom) safeGetWhite() Player {
	white := Player{Name: "White"}
	if gr.Game.White != nil {
		white = Player{
			ID:   gr.Game.White.ID,
			Name: gr.Game.White.Name,
		}
	}

	return white
}

func (gr *GameRoom) safeGetBlack() Player {
	black := Player{Name: "Black"}
	if gr.Game.Black != nil {
		black = Player{
			ID:   gr.Game.Black.ID,
			Name: gr.Game.Black.Name,
		}
	}

	return black
}

func (gr *GameRoom) getGameData() ([]byte, error) {
	return json.Marshal(GameData{
		Moves:       gr.Game.Moves,
		TimeControl: gr.Game.TimeControl,
		ThinkTime:   gr.getThinkTime(),
		Result:      "*",
		White:       gr.safeGetWhite(),
		Black:       gr.safeGetBlack(),
	})
}
