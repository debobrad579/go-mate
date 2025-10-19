const theme = localStorage.getItem("ui-theme")

if (
  theme === "dark"
  || (
    (theme == null || theme === "system")
    && window.matchMedia
    && window.matchMedia('(prefers-color-scheme: dark)').matches
  )
) {
  document.documentElement.className = "dark"
}
