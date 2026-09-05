import {
  createContext,
  type JSX,
  type PropsWithChildren,
  use,
  useCallback,
  useEffect,
  useState,
} from 'react'

import themes from './themes'

const defaultTheme = themes[0].name

interface ThemeContextType {
  theme: string
  setTheme: (theme: string) => void
}

const ThemeContext = createContext<ThemeContextType>({
  theme: defaultTheme,
  setTheme: () => {},
})

export function ThemeProvider({ children }: PropsWithChildren): JSX.Element {
  const [themeName, setThemeName] = useState<string>(
    () => localStorage.getItem('theme') || defaultTheme
  )

  const theme = themes.find((x) => x.name === themeName)

  const setTheme = useCallback((theme: string) => {
    setThemeName(theme)
    localStorage.setItem('theme', theme)
  }, [])

  useEffect(() => {
    if (!theme) {
      return
    }

    const load = theme.load ? theme.load() : Promise.resolve()
    load.then(() => {
      const root = window.document.documentElement
      root.classList.remove(
        ...Array.from(root.classList).filter((x) => x.startsWith('theme-'))
      )
      root.classList.add(`theme-${theme.name}`)
      localStorage.setItem('theme', theme.name)
    })
  }, [theme])

  return (
    <ThemeContext
      value={{
        theme: themeName,
        setTheme,
      }}
    >
      {children}
    </ThemeContext>
  )
}

export function useTheme(): [string, (theme: string) => void] {
  const context = use(ThemeContext)

  return [context.theme, context.setTheme]
}
