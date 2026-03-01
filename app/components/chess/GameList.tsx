import { GameRoom } from "@/types/chess"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useEventSource } from "@/hooks/useEventSource"
import { formatTimeControl } from "@/lib/formatters"
import { useNavigate } from "react-router"
import { Timer } from "lucide-react"

export function GameList() {
  const { data } = useEventSource<GameRoom[]>("/games")
  const navigate = useNavigate()

  return (
    <div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Player</TableHead>
            <TableHead>
              <Timer />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data?.map((room) => (
            <TableRow
              key={room.id}
              onClick={() =>
                (room.game.white == null || room.game.black == null) &&
                navigate(`/live/${room.id}`)
              }
              className={
                room.game.white != null && room.game.black != null
                  ? "opacity-50"
                  : "cursor-pointer"
              }
            >
              <TableCell>
                {room.game.white != null
                  ? room.game.white.name
                  : room.game.black != null
                    ? room.game.black.name
                    : "Anonymous"}
              </TableCell>
              <TableCell>{formatTimeControl(room.game.time_control)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
