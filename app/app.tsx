import React from "react"
import { createRoot } from "react-dom/client"
import { BrowserRouter, Route, Routes } from "react-router"
import { ThemeProvider } from "./components/themes/theme-provider"
import { ThemeToggle } from "./components/themes/theme-toggle"
import { PlayPage } from "./pages/play"
import { UserProvider } from "./user-context"

function App() {
  return (
    <UserProvider>
      <BrowserRouter basename="/app">
        <Routes>
          <Route path="/" element={<PlayPage />} />
          <Route path="/play" element={<PlayPage />} />
          <Route path="/2" element={<ThemeToggle />} />
          <Route
            path="/3"
            element={
              <div className="text-blue-500 text-4xl font-bold">Page 3</div>
            }
          />
        </Routes>
      </BrowserRouter>
    </UserProvider>
  )
}

createRoot(document.getElementById("app")!).render(
  <React.StrictMode>
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <App />
    </ThemeProvider>
  </React.StrictMode>,
)
