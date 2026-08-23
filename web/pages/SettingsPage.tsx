import { useEffect } from 'react'
import type { JSX } from 'react/jsx-runtime'
import { FluentBug16Regular } from '../components/icons/fluent-bug-16-regular'
import { FluentPaintBrush16Regular } from '../components/icons/fluent-paint-brush-16-regular'
import { ThemePreview } from '../components/icons/theme-preview'
import { useLoadTheme, useTheme, useThemes } from '../ThemeProvider'
import { Card } from './image-page/Card'

export function SettingsPage(): JSX.Element {
  const [currentTheme, setTheme] = useTheme()
  const loadTheme = useLoadTheme()
  const themes = useThemes()

  useEffect(() => {
    Promise.all(themes.map(loadTheme))
  }, [loadTheme, themes])

  return (
    <div className="flex flex-col items-center w-full pt-2 pb-10 px-2">
      <main className="min-w-[200px] max-w-[800px] w-full box-border space-y-6 mt-6">
        <Card
          persistenceKey="theming"
          tabs={[
            {
              icon: <FluentPaintBrush16Regular />,
              label: 'Theming',
              content: (
                <>
                  <div className="grid justify-center grid-cols-3">
                    {themes.map((theme) => (
                      <button
                        type="button"
                        key={theme.name}
                        onClick={() => setTheme(theme.name)}
                        className="flex flex-col items-center flex-shrink-1 cursor-pointer"
                      >
                        <div
                          className={`bg-surface-2-bg p-2 m-2 rounded-sm border ${theme.name === currentTheme ? 'border-accent ' : 'border-surface-1-stroke'}`}
                        >
                          <div className={`theme-${theme.name}`}>
                            <ThemePreview />
                          </div>
                        </div>
                        <p className="whitespace-pre">{theme.displayName}</p>
                        <p className="whitespace-pre text-xs">
                          {theme.description}
                        </p>
                      </button>
                    ))}
                  </div>
                </>
              ),
            },
          ]}
        />

        <Card
          persistenceKey="build"
          tabs={[
            {
              icon: <FluentBug16Regular />,
              label: 'Build',
              content: (
                <p>
                  Cupdate version:{' '}
                  {import.meta.env.VITE_CUPDATE_VERSION || 'development build'}.
                </p>
              ),
            },
          ]}
        />
      </main>
    </div>
  )
}
