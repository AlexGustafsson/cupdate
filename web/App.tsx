import { type JSX, useLayoutEffect } from 'react'
import { Link, Route, Routes, useLocation } from 'react-router-dom'

import { EventProvider } from './EventProvider'
import { FluentArrowLeft24Regular } from './components/icons/fluent-arrow-left-24-regular'
import { SimpleIconsRss } from './components/icons/simple-icons-rss'
import { DEFAULT_RSS_ENDPOINT } from './lib/api/api-client'
import { Dashboard } from './pages/Dashboard'
import { ImagePage } from './pages/ImagePage'

export function App(): JSX.Element {
  const location = useLocation()

  // biome-ignore lint/correctness/useExhaustiveDependencies: Trigger on change
  useLayoutEffect(() => {
    document.documentElement.scrollTo({ top: 0, left: 0, behavior: 'instant' })
  }, [location.pathname, location.search])

  return (
    <>
      <div className="fixed top-0 left-0 h-[64px] w-full grid grid-cols-3 items-center shadow-sm bg-white/90 dark:bg-[#1e1e1e]/70 z-250 backdrop-blur-md">
        <div className="justify-self-start ml-5">
          {location.pathname !== '/' && (
            <Link
              to={typeof location.state === 'string' ? location.state : '/'}
            >
              <FluentArrowLeft24Regular />
            </Link>
          )}
        </div>
        <div className="justify-self-center">
          <Link to="/">
            <h1 className="text-xl font-medium">Cupdate</h1>
          </Link>
        </div>
        <div className="justify-self-end mr-5">
          <a target="_blank" href={DEFAULT_RSS_ENDPOINT} rel="noreferrer">
            <SimpleIconsRss className="text-orange-400" />
          </a>
        </div>
      </div>
      <main className="pt-[64px]">
        <EventProvider>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/image" element={<ImagePage />} />
          </Routes>
        </EventProvider>
      </main>
    </>
  )
}
