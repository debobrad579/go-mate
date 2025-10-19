import React from "react"
import { createRoot } from "react-dom/client"
import { BrowserRouter, Route, Routes } from "react-router"
import { Chessboard } from "./components/chessboard/chessboard"
import { ThemeProvider } from "./components/themes/theme-provider"
import { ThemeToggle } from "./components/themes/theme-toggle"

function App() {
  return (
    <BrowserRouter basename="/app">
      <Routes>
        <Route path="/" element={<Chessboard fen="rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" onMove={(from: string, to: string) => {
          console.log("From: " + from)
          console.log("To: " + to)
        }} />} />
        <Route path="/2" element={<ThemeToggle />} />
        <Route path="/3" element={<div className="text-blue-500 text-4xl font-bold">Page 3</div>} />
      </Routes>
    </BrowserRouter>
  )
}

createRoot(document.getElementById("app")!).render(
  <React.StrictMode>
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <App />
    </ThemeProvider>
  </React.StrictMode>
)

