import { useEffect, useState } from "react"

export function useEventSource<T>(url: string) {
  const [data, setData] = useState<T | null>(null)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    const source = new EventSource(url)

    source.onmessage = (event) => {
      try {
        setData(JSON.parse(event.data))
        setError(null)
      } catch (e) {
        setError(
          e instanceof Error ? e : new Error("failed to parse event data"),
        )
      }
    }

    source.onerror = () => {
      setError(new Error("connection lost, reconnecting..."))
    }

    return () => source.close()
  }, [url])

  return { data, error }
}
