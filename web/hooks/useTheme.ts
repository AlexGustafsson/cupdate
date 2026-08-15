import { useEffect, useState } from 'react'

export const themes = ['system', 'light', 'dark'] as const
export type Theme = (typeof themes)[number]

const defaultTheme: Theme = 'system'

export function useTheme(): [Theme, (theme: Theme) => void] {
  const [theme, setTheme] = useState<Theme>(() => {
    const stored = (localStorage.getItem('theme') || '') as Theme
    if (themes.includes(stored)) {
      return stored
    } else {
      return defaultTheme
    }
  })

  useEffect(() => {
    const root = window.document.documentElement
    root.classList.remove(...themes.map((x) => `theme-${x}`))
    root.classList.add(`theme-${theme}`)

    localStorage.setItem('theme', theme)
  }, [theme])

  return [theme, setTheme]
}
