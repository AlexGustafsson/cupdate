import {
  createContext,
  type JSX,
  type PropsWithChildren,
  use,
  useCallback,
  useEffect,
  useState,
} from 'react'

export interface Theme {
  name: string
  displayName: string
  description: string
  file: string
}

const systemTheme: Theme = {
  name: 'system',
  displayName: 'System',
  description: "Cupdate's default auto theme",
  file: '', // Bundled
}

const bundledThemes: Theme[] = [
  systemTheme,
  {
    name: 'light',
    displayName: 'Light',
    description: "Cupdate's default light theme",
    file: '', // Bundled
  },
  {
    name: 'dark',
    displayName: 'Dark',
    description: "Cupdate's default dark theme",
    file: '', // Bundled
  },
]

interface ThemeContextType {
  theme: string
  setTheme: (theme: string) => void

  themes: Theme[]

  loadedThemes: string[]
  setLoadedThemes: React.Dispatch<React.SetStateAction<string[]>>
}

const ThemeContext = createContext<ThemeContextType>({
  theme: systemTheme.name,
  setTheme: () => {},

  themes: [...bundledThemes],

  loadedThemes: bundledThemes.map((x) => x.name),
  setLoadedThemes: () => {},
})

export function ThemeProvider({ children }: PropsWithChildren): JSX.Element {
  const [themeName, setThemeName] = useState<string>(
    () => localStorage.getItem('theme') || systemTheme.name
  )

  const loadTheme = useLoadTheme()

  const [themes, setThemes] = useState<Theme[]>(() => [...bundledThemes])
  const [loadedThemes, setLoadedThemes] = useState<string[]>(() =>
    bundledThemes.map((x) => x.name)
  )

  const theme = themes.find((x) => x.name === themeName)

  const setTheme = useCallback((theme: string) => {
    setThemeName(theme)
    localStorage.setItem('theme', theme)
  }, [])

  useEffect(() => {
    fetch('/assets/themes/index.json')
      .then((res) => res.json())
      .then((index: Theme[]) => {
        setThemes([...bundledThemes, ...index])
      })
  }, [])

  useEffect(() => {
    if (!theme) {
      return
    }

    loadTheme(theme).then(() => {
      const root = window.document.documentElement
      root.classList.remove(
        ...Array.from(root.classList).filter((x) => x.startsWith('theme-'))
      )
      root.classList.add(`theme-${theme.name}`)
      localStorage.setItem('theme', theme.name)
    })
  }, [theme, loadTheme])

  return (
    <ThemeContext
      value={{
        theme: themeName,
        themes,
        setTheme,
        loadedThemes,
        setLoadedThemes,
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

export function useThemes(): Theme[] {
  const context = use(ThemeContext)

  return context.themes
}

export function useLoadTheme(): (theme: Theme) => Promise<void> {
  const context = use(ThemeContext)

  return async (theme: Theme) => {
    if (context.loadedThemes.includes(theme.name)) {
      return
    }

    const res = await fetch(`/assets/themes/${theme.file}`)
    const content = await res.text()

    const sheet = new CSSStyleSheet()
    sheet.replace(content)

    document.adoptedStyleSheets.push(sheet)
    context.setLoadedThemes((current) => [...current, theme.name])
  }
}
