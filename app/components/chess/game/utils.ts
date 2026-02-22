export function convertMovesToPgn(moves: string[]): string {
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

export function getMoveNumberArrays(arr: string[]): [string, string][] {
  if (arr.length === 0) {
    return []
  }

  const moveSet: [string, string] = [arr[0] || "", arr[1] || ""]

  return [moveSet, ...getMoveNumberArrays(arr.slice(2))]
}
