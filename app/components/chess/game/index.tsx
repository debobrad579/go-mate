import {
  ChevronFirst,
  ChevronLast,
  ChevronLeft,
  ChevronRight,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Chessboard, type ChessboardProps } from "@/components/chess/board"
import type { Game } from "@/types/chess"
import { Clock } from "./clock"
import { MovesTable } from "./moves-table"
import { useChessGame } from "./use-chess-game"
import { MovesList } from "./moves-list"

type ChessGameProps = {
  gameData: Game
  onMove: ChessboardProps["onMove"]
}

export function ChessGame({
  gameData: { white, black, moves, result, thinkTime },
  onMove,
}: ChessGameProps) {
  const {
    game,
    undoCount,
    tick,
    previousMove,
    mouseOverBoard,
    moveSoundRef,
    captureSoundRef,
    loadPgn,
    reset,
    undoMove,
    redoMove,
    setUndoCount,
  } = useChessGame({ moves, result, thinkTime })

  function handleWhiteMoveClick(index: number) {
    loadPgn(moves.slice(0, index * 2 + 1))
    setUndoCount(moves.length - index * 2 - 1)
  }

  function handleBlackMoveClick(index: number) {
    loadPgn(moves.slice(0, index * 2 + 2))
    setUndoCount(moves.length - index * 2 - 2)
  }

  return (
    <div className="@container">
      <audio ref={moveSoundRef} src={"/static/audio/move.mp3"} />
      <audio ref={captureSoundRef} src={"/static/audio/capture.mp3"} />
      <div
        className="flex flex-col @lg:flex-row gap-2"
        onMouseEnter={() => (mouseOverBoard.current = true)}
        onMouseLeave={() => (mouseOverBoard.current = false)}
      >
        <div className="flex-1 space-y-2">
          <div>
            <Clock
              className="bg-gray-800 text-white"
              undoCount={undoCount}
              turn={game.turn() === "b"}
              result={
                { "0-1": "win", "1-0": "loss", "1/2-1/2": "draw", "*": "*" }[
                  result
                ] as "win" | "loss" | "draw" | "*"
              }
              thinkTime={(thinkTime ?? 0) + tick}
              player={black}
            />
            <Chessboard
              fen={game.fen()}
              previousMove={
                previousMove
                  ? { from: previousMove.from, to: previousMove.to }
                  : undefined
              }
              check={game.in_check() ? game.turn() : undefined}
              onMove={onMove}
              draggablePieces={game.turn()}
            />
            <Clock
              className="bg-gray-200 text-black"
              undoCount={undoCount}
              turn={game.turn() === "w"}
              result={
                { "1-0": "win", "0-1": "loss", "1/2-1/2": "draw", "*": "*" }[
                  result
                ] as "win" | "loss" | "draw" | "*"
              }
              thinkTime={(thinkTime ?? 0) + tick}
              player={white}
            />
          </div>
          <div className="flex gap-2">
            <Button
              className="w-full"
              onClick={reset}
              disabled={undoCount === moves.length}
            >
              <ChevronFirst />
            </Button>
            <Button
              className="w-full"
              onClick={undoMove}
              disabled={undoCount === moves.length}
            >
              <ChevronLeft />
            </Button>
            <Button
              className="w-full"
              onClick={redoMove}
              disabled={undoCount === 0}
            >
              <ChevronRight />
            </Button>
            <Button
              className="w-full"
              onClick={() => {
                loadPgn()
                setUndoCount(0)
              }}
              disabled={undoCount === 0}
            >
              <ChevronLast />
            </Button>
          </div>
        </div>
        <MovesTable
          moves={moves}
          result={result}
          undoCount={undoCount}
          onWhiteMoveClick={handleWhiteMoveClick}
          onBlackMoveClick={handleBlackMoveClick}
        />
        <MovesList
          moves={moves}
          result={result}
          undoCount={undoCount}
          onWhiteMoveClick={handleWhiteMoveClick}
          onBlackMoveClick={handleBlackMoveClick}
        />
      </div>
    </div>
  )
}
