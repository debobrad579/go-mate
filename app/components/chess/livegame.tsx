import { Game } from "@/types/chess"
import { ChessGame } from "./game"
import useWebSocket, { ReadyState } from "react-use-websocket"
import { useUser } from "@/user-context"
import { useEffect } from "react"

export const sampleGame: Game = {
  moves: ["e4", "e5", "Bc4", "Nc6", "Qh5", "Nf6", "Qxf7#"],
  result: "1-0",
  white: {
    name: "AliceBot",
    elo: "1850",
    title: "WFM",
    timestamps: [5, 12, 20, 28],
  },
  black: {
    name: "StockfishTest",
    elo: "1900",
    timestamps: [6, 15, 24],
  },
}

export function LiveGame() {
  const { user } = useUser()

  const { sendMessage, lastMessage, readyState } = useWebSocket("/ws")

  useEffect(() => {
    if (lastMessage != null) {
      console.log(`New message: ${lastMessage.data}`)
    }
  }, [lastMessage])

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState]

  return (
    <div>
      {connectionStatus}
      <br />
      {user != null ? <h1>User: {user.Name}</h1> : <h1>No user logged in</h1>}
      <ChessGame
        gameData={sampleGame}
        onMove={(from, to) => {
          sendMessage(`from ${from} to ${to}`)
        }}
      />
    </div>
  )
}
