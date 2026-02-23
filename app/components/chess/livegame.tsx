import { ChessGame, type ChessGameHandle } from "./game"
import { useUser } from "@/user-context"
import { useEffect, useRef, useState, useCallback } from "react"
import type { Game } from "./types"

export const defaultGame: Game = {
  moves: [],
  result: "*",
  white: { name: "White", elo: "1500" },
  black: { name: "Black", elo: "1500" },
}

type ReadyState = "Connecting" | "Open" | "Closing" | "Closed"

export function LiveGame() {
  const [game, setGameData] = useState(defaultGame)
  const [readyState, setReadyState] = useState<ReadyState>("Connecting")
  const { user } = useUser()
  const wsRef = useRef<WebSocket | null>(null)
  const chessGameRef = useRef<ChessGameHandle>(null)

  useEffect(() => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws`)
    wsRef.current = ws

    ws.onopen = () => setReadyState("Open")
    ws.onclose = () => setReadyState("Closed")
    ws.onerror = () => setReadyState("Closed")
    ws.onmessage = (event) => {
      if (event.data) {
        const data: Game = JSON.parse(event.data)
        setGameData(data)
      }
    }

    return () => {
      setReadyState("Closing")
      ws.close()
    }
  }, [])

  const sendJsonMessage = useCallback((message: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    }
  }, [])

  return (
    <div>
      {readyState}
      <br />
      {user != null ? <h1>User: {user.Name}</h1> : <h1>No user logged in</h1>}
      <ChessGame
        ref={chessGameRef}
        gameData={game}
        onMove={(move) => {
          if (!chessGameRef.current?.makeMove(move)) return
          sendJsonMessage(move)
        }}
      />
    </div>
  )
}
