import React from "react"
import ReactDOM from "react-dom/client"
import { BrowserRouter, Route, Routes } from "react-router"

function App() {
  return (
    <BrowserRouter basename="/app">
      <Routes>
        <Route path="/" element={<div className="text-blue-500 text-4xl font-bold">Index</div>} />
        <Route path="/2" element={<div className="text-blue-500 text-4xl font-bold">Page 2</div>} />
        <Route path="/3" element={<div className="text-blue-500 text-4xl font-bold">Page 3</div>} />
      </Routes>
    </BrowserRouter>
  )
}

ReactDOM.createRoot(document.getElementById("app")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)

