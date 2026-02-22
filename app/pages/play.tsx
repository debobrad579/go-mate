import { ChessGame } from "@/components/chess/game"
import type { Game } from "@/types/chess"
import { useUser } from "@/user-context"

export const sampleGame: Game = {
  moves: ["e4", "e5", "Bc4", "Nc6", "Qh5", "Nf6", "Qxf7#"],
  result: "1-0",
  white: {
    name: "AliceBot",
    elo: "1850",
    title: "WFM",
    timestamps: [
      5, // after 1.e4
      12, // after 2.Bc4
      20, // after 3.Qh5
      28, // after 4.Qxf7#
    ],
  },
  black: {
    name: "StockfishTest",
    elo: "1900",
    timestamps: [
      6, // after 1...e5
      15, // after 2...Nc6
      24, // after 3...Nf6
    ],
  },
}

export function PlayPage() {
  const { user } = useUser()

  return (
    <div className="h-screen w-screen flex items-center justify-center overflow-hidden">
      <div className="aspect-square w-[min(100vw,100vh)] p-4">
        {user != null && <h1>{user.Name}</h1>}
        <ChessGame gameData={sampleGame} />
      </div>
    </div>
  )
}
