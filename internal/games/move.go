package games

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/debobrad579/chessgo/internal/chess"
)

func (gr *GameRoom) MakeMove(message []byte, color chess.Color) error {
	var move chess.Move
	if err := json.Unmarshal(message, &move); err != nil {
		return err
	}

	gr.mu.Lock()
	defer gr.mu.Unlock()

	if gr.Game.White == nil || gr.Game.Black == nil {
		return errors.New("game not started")
	}

	if (color == chess.White && gr.Game.Turn() == chess.Black) || (color == chess.Black && gr.Game.Turn() == chess.White) {
		return nil
	}

	if !gr.Game.IsMoveValid(move) {
		return nil
	}

	if color == chess.White {
		gr.whiteTime -= int(time.Since(gr.turnStart).Milliseconds()) - gr.Game.TimeControl.Increment
		move.Timestamp = gr.whiteTime
	} else {
		gr.blackTime -= int(time.Since(gr.turnStart).Milliseconds()) - gr.Game.TimeControl.Increment
		move.Timestamp = gr.blackTime
	}

	gr.Game.Move(move)
	gr.turnStart = time.Now()

	select {
	case gr.broadcast <- struct{}{}:
	default:
	}

	return nil
}
