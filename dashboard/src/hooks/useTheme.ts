// src/hooks/useTheme.ts — dark/light mode toggle (Plan 01-03 uses this)
import { useEffect, useState } from "react"

type Theme = "dark" | "light" | "system"

export function useTheme() {
  const [theme, setTheme] = useState<Theme>(() => {
    return (localStorage.getItem("theme") as Theme) ?? "system"
  })

  useEffect(() => {
    const root = document.documentElement
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches

    if (theme === "dark" || (theme === "system" && prefersDark)) {
      root.classList.add("dark")
    } else {
      root.classList.remove("dark")
    }

    if (theme === "system") {
      localStorage.removeItem("theme")
    } else {
      localStorage.setItem("theme", theme)
    }
  }, [theme])

  return { theme, setTheme }
}
