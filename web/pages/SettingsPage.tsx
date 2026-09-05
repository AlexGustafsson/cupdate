import { useCallback, useEffect } from 'react'
import type { JSX } from 'react/jsx-runtime'
import { Card } from '../components/Card'
import { FluentBug16Regular } from '../components/icons/fluent-bug-16-regular'
import { FluentPaintBrush16Regular } from '../components/icons/fluent-paint-brush-16-regular'
import { ThemePreview } from '../components/icons/theme-preview'
import { useCollapseState } from '../hooks/useCollapseState'
import { useLoadTheme, useTheme, useThemes } from '../ThemeProvider'
import { TableOfContents } from './image-page/TableOfContents'

export function SettingsPage(): JSX.Element {
  const [currentTheme, setTheme] = useTheme()
  const loadTheme = useLoadTheme()
  const themes = useThemes()

  useEffect(() => {
    Promise.all(themes.map(loadTheme))
  }, [loadTheme, themes])

  const scrollIntoView = useCallback((id: string) => {
    const element = document.getElementById(id)
    element?.scrollIntoView({ block: 'start', behavior: 'smooth' })
  }, [])

  const [collapsed, toggleCollapsed] = useCollapseState()

  return (
    <div className="flex flex-col items-center w-full pt-2 pb-10 px-2">
      <div className="w-full mt-6 grid grid-cols-1 lg:grid-cols-[1fr_max(200px,min(100%,800px))_1fr]">
        {/* Left column */}
        <div className="pr-6 hidden lg:block"></div>
        {/* Big center column */}
        <div className="flex flex-col w-full space-y-6">
          {/* Theming */}
          <Card
            id="theming"
            collapsed={collapsed.has('theming')}
            onToggleCollapsed={() => toggleCollapsed('theming')}
            tabs={[
              {
                icon: <FluentPaintBrush16Regular />,
                label: 'Theming',
                content: (
                  <>
                    <div className="grid justify-center grid-cols-3 space-x-2 space-y-2">
                      {themes.map((theme) => (
                        <button
                          type="button"
                          key={theme.name}
                          onClick={() => setTheme(theme.name)}
                          className="flex-col p-2"
                          disabled={theme.name === currentTheme}
                          tabIndex={0}
                        >
                          <div
                            className={`bg-surface-2-bg p-2 m-2 rounded-sm border border-surface-1-stroke`}
                          >
                            <div className={`theme-${theme.name} w-full`}>
                              <ThemePreview className="w-full" />
                            </div>
                          </div>
                          <p className="whitespace-pre">{theme.displayName}</p>
                          <p className="text-xs">{theme.description}</p>
                        </button>
                      ))}
                    </div>
                  </>
                ),
              },
            ]}
          />

          {/* Build info */}
          <Card
            id="build"
            collapsed={collapsed.has('build')}
            onToggleCollapsed={() => toggleCollapsed('build')}
            tabs={[
              {
                icon: <FluentBug16Regular />,
                label: 'Build',
                content: (
                  <p>
                    Cupdate version:{' '}
                    {import.meta.env.VITE_CUPDATE_VERSION ||
                      'development build'}
                    .
                  </p>
                ),
              },
            ]}
          />
        </div>
        {/* Right column */}
        <div className="pl-6 hidden lg:block">
          <div className="sticky top-[64px]">
            <TableOfContents
              onClick={scrollIntoView}
              onToggleCollapsed={toggleCollapsed}
              items={[
                {
                  id: 'theming',
                  icon: <FluentPaintBrush16Regular />,
                  label: 'Theming',
                  collapsed: collapsed.has('theming'),
                },
                {
                  id: 'build',
                  icon: <FluentBug16Regular />,
                  label: 'Build',
                  collapsed: collapsed.has('build'),
                },
              ]}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
