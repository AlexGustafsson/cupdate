import { SimpleIconsRss } from './components/icons/simple-icon-rss'
import { Dashboard } from './pages/Dashboard'

export function App(): JSX.Element {
  return (
    <>
      <div className="absolute top-0 right-0 p-5">
        <a target="_blank" href="/feed.rss">
          <SimpleIconsRss className="text-orange-400" />
        </a>
      </div>
      <Dashboard />
    </>
  )
}
