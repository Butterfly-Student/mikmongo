import { createContext, useCallback, useContext, useEffect, useState } from "react"

export type Theme = "dark" | "light" | "system"

interface ThemeContextValue {
  theme: Theme
  resolvedTheme: "dark" | "light"
  setTheme: (theme: Theme) => void
  toggleTheme: () => void
}

const ThemeContext = createContext<ThemeContextValue | null>(null)

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window === "undefined") return "system"
    return (localStorage.getItem("theme") as Theme) ?? "system"
  })

  const getResolved = useCallback((t: Theme): "dark" | "light" => {
    if (t === "system") {
      return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
    }
    return t
  }, [])

  const [resolvedTheme, setResolvedTheme] = useState<"dark" | "light">(() => getResolved(theme))

  const applyTheme = useCallback(
    (t: Theme) => {
      const resolved = getResolved(t)
      const root = document.documentElement
      if (resolved === "dark") {
        root.classList.add("dark")
      } else {
        root.classList.remove("dark")
      }
      setResolvedTheme(resolved)
    },
    [getResolved]
  )

  const setTheme = useCallback(
    (t: Theme) => {
      setThemeState(t)
      if (t === "system") {
        localStorage.removeItem("theme")
      } else {
        localStorage.setItem("theme", t)
      }
      applyTheme(t)
    },
    [applyTheme]
  )

  const toggleTheme = useCallback(() => {
    const next = resolvedTheme === "dark" ? "light" : "dark"
    setTheme(next)
  }, [resolvedTheme, setTheme])

  // Sync with system preference changes when theme === "system"
  useEffect(() => {
    if (theme !== "system") return
    const mq = window.matchMedia("(prefers-color-scheme: dark)")
    const handler = () => applyTheme("system")
    mq.addEventListener("change", handler)
    return () => mq.removeEventListener("change", handler)
  }, [theme, applyTheme])

  // Apply on mount to sync with any server-rendered state
  useEffect(() => {
    applyTheme(theme)
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <ThemeContext.Provider value={{ theme, resolvedTheme, setTheme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  )
}

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext)
  if (!ctx) {
    throw new Error("useTheme must be used inside <ThemeProvider>")
  }
  return ctx
}
