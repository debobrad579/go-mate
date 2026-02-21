export type Game = {
  moves: string[]
  result: "1-0" | "0-1" | "1/2-1/2" | "*"
  thinkTime?: number
  white: { name: string; elo: string; timestamps: number[]; title?: string }
  black: { name: string; elo: string; timestamps: number[]; title?: string }
}
