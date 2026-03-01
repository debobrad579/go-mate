import { CreateGameButton } from "@/components/chess/CreateGameButton"
import { GameList } from "@/components/chess/GameList"

export default function HomePage() {
  return (
    <>
      <CreateGameButton />
      <GameList />
    </>
  )
}
