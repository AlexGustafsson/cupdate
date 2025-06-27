import { type JSX, useCallback, useLayoutEffect } from 'react'
import { Link, Route, Routes, useLocation } from 'react-router-dom'

import { EventProvider } from './EventProvider'
import { InfoTooltip } from './components/InfoTooltip'
import { Menu } from './components/Menu'
import { FluentAlert24Regular } from './components/icons/fluent-alert-24-regular'
import { FluentAlertBadge24Regular } from './components/icons/fluent-alert-badge-24-regular'
import { FluentArrowLeft24Regular } from './components/icons/fluent-arrow-left-24-regular'
import { FluentOpen16Regular } from './components/icons/fluent-open-16-regular'
import { FluentWarning16Filled } from './components/icons/fluent-warning-16-filled'
import { useWebPushSubscription } from './hooks'
import { DEFAULT_RSS_ENDPOINT } from './lib/api/api-client'
import { Dashboard } from './pages/Dashboard'
import { ImagePage } from './pages/ImagePage'

export function App(): JSX.Element {
  const location = useLocation()

  // biome-ignore lint/correctness/useExhaustiveDependencies: Trigger on change
  useLayoutEffect(() => {
    document.documentElement.scrollTo({ top: 0, left: 0, behavior: 'instant' })
  }, [location.pathname, location.search])

  const [
    webPushSupported,
    webPushSubscription,
    webPushSynced,
    subscribe,
    unsubscribe,
  ] = useWebPushSubscription()

  const webPushOk =
    webPushSubscription.status !== 'rejected' &&
    webPushSynced.status !== 'rejected'

  const toggleWebPushSubscription = useCallback(() => {
    if (webPushSubscription.status !== 'resolved') {
      return
    }

    if (webPushSubscription.value) {
      unsubscribe()
    } else {
      subscribe()
    }
  }, [webPushSubscription, subscribe, unsubscribe])

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
          <Menu
            icon={
              webPushOk ? (
                <FluentAlert24Regular />
              ) : (
                <FluentAlertBadge24Regular />
              )
            }
          >
            {webPushSupported && (
              <li
                className="flex items-center"
                onClick={toggleWebPushSubscription}
              >
                {webPushSubscription.status === 'resolved'
                  ? webPushSubscription.value !== null
                    ? 'Unsubscribe'
                    : 'Subscribe'
                  : ''}
                {webPushSubscription.status === 'rejected' ? (
                  <InfoTooltip
                    icon={<FluentWarning16Filled className="text-red-400" />}
                  >
                    {webPushSubscription.error.toString()}
                  </InfoTooltip>
                ) : webPushSynced.status === 'rejected' ? (
                  <InfoTooltip
                    icon={<FluentWarning16Filled className="text-red-400" />}
                  >
                    {webPushSynced.error.toString()}
                  </InfoTooltip>
                ) : undefined}
              </li>
            )}
            <a target="_blank" href={DEFAULT_RSS_ENDPOINT} rel="noreferrer">
              <li className="flex items-center gap-x-2">
                RSS feed <FluentOpen16Regular />
              </li>
            </a>
          </Menu>
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
