package games

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/debobrad579/chessgo/internal/chess"
)

func (gr *GameRoom) MakeMove(message []byte, playerRole PlayerRole) error {
	if playerRole == Spectator {
		return errors.New("cannot make move as spectator")
	}

	var move chess.Move
	if err := json.Unmarshal(message, &move); err != nil {
		return err
	}

	gr.mu.Lock()
	defer gr.mu.Unlock()

	if gr.Game.White == nil || gr.Game.Black == nil {
		return errors.New("game not started")
	}

	if (playerRole == White) != (gr.Game.Turn() == chess.White) {
		return errors.New("not your turn")
	}

	if !gr.Game.IsMoveValid(move) {
		return errors.New("invalid move")
	}

	if playerRole == White {
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
