import { useEffect, useMemo, useRef, useState } from "react"
import { Chess } from "chess.js"
import {
  ChevronFirst,
  ChevronLast,
  ChevronLeft,
  ChevronRight,
} from "lucide-react"
import { useEventListener } from "@/hooks/useEventListener"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableFooter,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { cn } from "@/lib/utils"
import { MoveCell } from "./move-cell"
import { Chessboard } from "./chessboard"
import { Game } from "@/types/chess"
import { Clock } from "./clock"

export function ChessGame({
  gameData: { white, black, moves, result, thinkTime },
}: {
  gameData: Game
}) {
  const [game, setGame] = useState(new Chess())
  const [undoCount, setUndoCount] = useState(0)
  const [tick, setTick] = useState(0)
  const [firstLoad, setFirstLoad] = useState(true)
  const mouseOverBoard = useRef(false)
  const tableScrollAreaRef = useRef<HTMLDivElement>(null)
  const listScrollAreaRef = useRef<HTMLDivElement>(null)
  const moveSoundRef = useRef<HTMLAudioElement>(null)
  const captureSoundRef = useRef<HTMLAudioElement>(null)
  const gameLength = useMemo(() => game.history().length, [game])

  useEffect(() => {
    if (thinkTime == null || result !== "*") return

    const start = Date.now()

    const interval = setInterval(() => {
      const elapsed = Math.floor((Date.now() - start) / 1000)
      setTick(elapsed)
    }, 1000)

    return () => clearInterval(interval)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [thinkTime, result, moves.length])

  useEventListener("keydown", (e: KeyboardEvent) => {
    if (!mouseOverBoard.current) return

    if (e.key === "ArrowLeft") {
      e.preventDefault()
      undoMove()
    }
    if (e.key === "ArrowRight") {
      e.preventDefault()
      redoMove()
    }
    if (e.key === "ArrowUp") {
      e.preventDefault()
      reset()
    }
    if (e.key === "ArrowDown") {
      e.preventDefault()
      loadPgn()
      setUndoCount(0)
    }
  })

  useEffect(() => {
    loadPgn()
    setUndoCount(0)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [moves.length])

  useEffect(() => {
    const move = game.history().at(-1)

    if (
      !navigator.userActivation.hasBeenActive ||
      !move ||
      !gameLength ||
      !moveSoundRef.current ||
      !captureSoundRef.current
    )
      return

    if (firstLoad) return setFirstLoad(false)

    if (move.includes("x")) {
      captureSoundRef.current.currentTime = 0
      captureSoundRef.current.play()
    } else {
      moveSoundRef.current.currentTime = 0
      moveSoundRef.current.play()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [gameLength, moves.length])

  function reset() {
    const gameCopy = { ...game }
    gameCopy.reset()
    setGame(gameCopy)
    setUndoCount(moves.length)
  }

  function loadPgn(movesArr: string[] | undefined = moves) {
    const gameCopy = { ...game }
    gameCopy.load_pgn(convertMovesToPgn(movesArr))
    setGame(gameCopy)
  }

  function undoMove() {
    if (undoCount === moves.length) return

    const gameCopy = { ...game }
    gameCopy.undo()
    setGame(gameCopy)
    setUndoCount((prevCount) => prevCount + 1)
  }

  function redoMove() {
    if (undoCount === 0) return
    const move = moves.at(-undoCount)
    if (!move) return

    const gameCopy = { ...game }
    if (undoCount === 1) {
      gameCopy.load_pgn(convertMovesToPgn(moves))
    } else {
      gameCopy.load_pgn(
        convertMovesToPgn(moves.slice(0, length - (undoCount - 1))),
      )
    }

    setGame(gameCopy)
    setUndoCount((prevCount) => prevCount - 1)
  }

  const previousMove = game.history({ verbose: true }).at(-1)

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
        <ScrollArea
          ref={tableScrollAreaRef}
          className="w-[10.5rem] pr-1 aspect-video hidden @lg:block"
        >
          <Table>
            <TableHeader>
              <TableRow className="text-muted-foreground">
                <MoveCell
                  active={undoCount === moves.length}
                  undoCount={undoCount}
                  scrollAreaRef={tableScrollAreaRef}
                  isTableCell
                  noStyles
                >
                  No.
                </MoveCell>
                <TableCell>White</TableCell>
                <TableCell>Black</TableCell>
              </TableRow>
            </TableHeader>
            <TableBody>
              {getMoveNumberArrays(moves).map((moveSet, index) => (
                <TableRow key={index}>
                  <TableCell>{index + 1}.</TableCell>
                  <MoveCell
                    onClick={() => {
                      loadPgn(moves.slice(0, index * 2 + 1))
                      setUndoCount(moves.length - index * 2 - 1)
                    }}
                    active={undoCount === moves.length - index * 2 - 1}
                    undoCount={undoCount}
                    scrollAreaRef={tableScrollAreaRef}
                    isTableCell
                  >
                    {moveSet[0]}
                  </MoveCell>
                  <MoveCell
                    onClick={() => {
                      loadPgn(moves.slice(0, index * 2 + 2))
                      setUndoCount(moves.length - index * 2 - 2)
                    }}
                    active={undoCount === moves.length - index * 2 - 2}
                    undoCount={undoCount}
                    scrollAreaRef={tableScrollAreaRef}
                    isTableCell
                  >
                    {moveSet[1]}
                  </MoveCell>
                </TableRow>
              ))}
            </TableBody>
            <TableFooter>
              <TableRow>
                <TableCell className="font-bold text-right">
                  {result.split("-")[0]}
                </TableCell>
                <TableCell className="font-bold text-center">-</TableCell>
                <TableCell className="font-bold">
                  {result.split("-")[1]}
                </TableCell>
              </TableRow>
            </TableFooter>
          </Table>
        </ScrollArea>
        <ScrollArea
          ref={listScrollAreaRef}
          className="@lg:hidden w-full pb-2 text-nowrap"
        >
          <MoveCell
            active={undoCount === moves.length}
            undoCount={undoCount}
            scrollAreaRef={listScrollAreaRef}
            noStyles
          />
          <div className={cn("flex gap-4 w-full")}>
            {getMoveNumberArrays(moves).map((moveSet, index) => (
              <div key={index} className="flex gap-2">
                <div>{index + 1}.</div>
                <MoveCell
                  onClick={() => {
                    loadPgn(moves.slice(0, index * 2 + 1))
                    setUndoCount(moves.length - index * 2 - 1)
                  }}
                  active={undoCount === moves.length - index * 2 - 1}
                  undoCount={undoCount}
                  scrollAreaRef={listScrollAreaRef}
                >
                  {moveSet[0]}
                </MoveCell>
                <MoveCell
                  onClick={() => {
                    loadPgn(moves.slice(0, index * 2 + 2))
                    setUndoCount(moves.length - index * 2 - 2)
                  }}
                  active={undoCount === moves.length - index * 2 - 2}
                  undoCount={undoCount}
                  scrollAreaRef={listScrollAreaRef}
                >
                  {moveSet[1]}
                </MoveCell>
              </div>
            ))}
            <div className="font-bold whitespace-nowrap">
              {result.replace("-", " - ")}
            </div>
          </div>
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </div>
    </div>
  )
}

function convertMovesToPgn(moves: string[]): string {
  let pgn = ""
  for (let i = 0; i < moves.length; i++) {
    const move = moves[i]
    if (i % 2 === 0) {
      pgn += `${Math.floor(i / 2) + 1}. ${move} `
    } else {
      pgn += `${move} `
    }
  }
  return pgn.trim()
}

function getMoveNumberArrays(arr: string[]): [string, string][] {
  if (arr.length === 0) {
    return []
  }

  const moveSet: [string, string] = [arr[0] || "", arr[1] || ""]

  return [moveSet, ...getMoveNumberArrays(arr.slice(2))]
}
