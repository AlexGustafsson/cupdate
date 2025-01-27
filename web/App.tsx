import { type JSX, useCallback, useState } from 'react'
import { NavLink, Route, Routes, useLocation } from 'react-router-dom'

import { type Event, RSSFeedEndpoint, useEvents } from './api'
import { Toast } from './components/Toast'
import { FluentArrowLeft24Regular } from './components/icons/fluent-arrow-left-24-regular'
import { FluentWeatherMoon24Regular } from './components/icons/fluent-weather-moon-24-regular'
import { FluentWeatherSunny24Regular } from './components/icons/fluent-weather-sunny-24-regular'
import { SimpleIconsRss } from './components/icons/simple-icons-rss'
import { useTheme } from './hooks/useTheme'
import { Dashboard } from './pages/Dashboard'
import { ImagePage } from './pages/ImagePage'

export function App(): JSX.Element {
  const location = useLocation()

  const [isUpdateAvailable, setIsUpdateAvailable] = useState(false)

  const [theme, setTheme] = useTheme()

  const onEvent = useCallback((e: Event) => {
    switch (e.type) {
      case 'imageUpdated':
        setIsUpdateAvailable(true)
    }
  }, [])

  useEvents(onEvent)
  return (
    <>
      <div className="fixed top-0 left-0 h-[64px] w-full grid grid-cols-3 items-center shadow-sm bg-white dark:bg-[#1e1e1e] z-50">
        <div className="justify-self-start ml-5">
          {location.pathname !== '/' && (
            <NavLink to="/">
              <FluentArrowLeft24Regular />
            </NavLink>
          )}
        </div>
        <div className="justify-self-center">
          <NavLink to="/">
            <h1 className="text-xl font-medium">Cupdate</h1>
          </NavLink>
        </div>
        <div className="justify-self-end mr-5 flex items-center gap-x-2">
          <button
            type="button"
            onClick={() =>
              setTheme((current) =>
                current === 'light'
                  ? 'dark'
                  : current === undefined
                    ? 'light'
                    : undefined
              )
            }
          >
            {theme === 'light' ? (
              <FluentWeatherSunny24Regular />
            ) : theme === 'dark' ? (
              <FluentWeatherMoon24Regular />
            ) : (
              <FluentWeatherSunny24Regular />
            )}
          </button>
          <a target="_blank" href={RSSFeedEndpoint} rel="noreferrer">
            <SimpleIconsRss className="text-orange-400" />
          </a>
        </div>
      </div>
      <div className="fixed bottom-0 right-0 p-4 z-50">
        {isUpdateAvailable && (
          <Toast
            title="New data available"
            body="One or more images have been updated. Reload to view the latest data."
            secondaryAction="Dismiss"
            onSecondaryAction={() => setIsUpdateAvailable(false)}
            primaryAction="Reload"
            onPrimaryAction={() =>
              window.location.replace(window.location.href)
            }
          />
        )}
      </div>
      <main className="pt-[64px]">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/image" element={<ImagePage />} />
        </Routes>
      </main>
    </>
  )
}
