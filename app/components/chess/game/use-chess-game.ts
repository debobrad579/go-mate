import { useEffect, useMemo, useRef, useState } from "react"
import { Chess } from "chess.js"
import { useEventListener } from "@/hooks/useEventListener"
import { convertMovesToPgn } from "./utils"

export function useChessGame({
  moves,
  result,
  thinkTime,
}: {
  moves: string[]
  result: string
  thinkTime?: number | null
}) {
  const [game, setGame] = useState(new Chess())
  const [undoCount, setUndoCount] = useState(0)
  const [tick, setTick] = useState(0)
  const [firstLoad, setFirstLoad] = useState(true)
  const mouseOverBoard = useRef(false)
  const moveSoundRef = useRef<HTMLAudioElement>(null)
  const captureSoundRef = useRef<HTMLAudioElement>(null)
  const gameLength = useMemo(() => game.history().length, [game])

  useEffect(() => {
    if (thinkTime == null || result !== "*") return
    const start = Date.now()
    const interval = setInterval(() => {
      setTick(Math.floor((Date.now() - start) / 1000))
    }, 1000)
    return () => clearInterval(interval)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [thinkTime, result, moves.length])

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

    const soundRef = move.includes("x") ? captureSoundRef : moveSoundRef
    if (soundRef.current != null) {
      soundRef.current.currentTime = 0
      soundRef.current.play()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [gameLength, moves.length])

  useEventListener("keydown", (e: KeyboardEvent) => {
    if (!mouseOverBoard.current) return
    const actions: Record<string, () => void> = {
      ArrowLeft: undoMove,
      ArrowRight: redoMove,
      ArrowUp: reset,
      ArrowDown: () => {
        loadPgn()
        setUndoCount(0)
      },
    }
    if (e.key in actions) {
      e.preventDefault()
      actions[e.key]()
    }
  })

  function reset() {
    const gameCopy = { ...game }
    gameCopy.reset()
    setGame(gameCopy)
    setUndoCount(moves.length)
  }

  function loadPgn(movesArr: string[] = moves) {
    const gameCopy = { ...game }
    gameCopy.load_pgn(convertMovesToPgn(movesArr))
    setGame(gameCopy)
  }

  function undoMove() {
    if (undoCount === moves.length) return
    const gameCopy = { ...game }
    gameCopy.undo()
    setGame(gameCopy)
    setUndoCount((prev) => prev + 1)
  }

  function redoMove() {
    if (undoCount === 0) return
    const move = moves.at(-undoCount)
    if (!move) return
    const gameCopy = { ...game }
    gameCopy.load_pgn(
      convertMovesToPgn(
        undoCount === 1
          ? moves
          : moves.slice(0, moves.length - (undoCount - 1)),
      ),
    )
    setGame(gameCopy)
    setUndoCount((prev) => prev - 1)
  }

  const previousMove = game.history({ verbose: true }).at(-1)

  return {
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
  }
}
