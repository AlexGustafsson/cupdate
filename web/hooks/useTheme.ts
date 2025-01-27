import { type Dispatch, useEffect, useState } from 'react'

export function useTheme(): [
  'light' | 'dark' | undefined,
  Dispatch<React.SetStateAction<'light' | 'dark' | undefined>>,
] {
  const [theme, setTheme] = useState<'light' | 'dark' | undefined>()

  useEffect(() => {
    const item = localStorage.getItem('theme')
    if (item === 'light' || item === 'dark') {
      setTheme(item)
    }
  }, [])

  useEffect(() => {
    switch (theme) {
      case 'light':
        localStorage.setItem('theme', 'light')
        document.documentElement.setAttribute('data-theme', 'light')
        break
      case 'dark':
        localStorage.setItem('theme', 'dark')
        document.documentElement.setAttribute('data-theme', 'dark')
        break
      case undefined:
        localStorage.removeItem('theme')
        document.documentElement.removeAttribute('data-theme')
        break
    }
  }, [theme])

  return [theme, setTheme]
}
