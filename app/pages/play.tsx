import { LiveGame } from "@/components/chess/livegame"

export function PlayPage() {
  return (
    <div className="h-screen w-screen flex items-center justify-center overflow-hidden">
      <div className="aspect-square w-[min(100vw,100vh)] p-4">
        <LiveGame />
      </div>
    </div>
  )
}
