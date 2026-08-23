import { type JSX, useEffect, useState } from 'react'
import { FluentDismiss16Regular } from './icons/fluent-dismiss-16-regular'
import { FluentInfo16Regular } from './icons/fluent-info-16-regular'

export function DemoWarning(): JSX.Element | null {
  const [showDemoWarning, setShowDemoWarning] = useState(() => {
    const item = sessionStorage.getItem('showDemoWarning')
    if (item) {
      return item === 'true'
    }
    return import.meta.env.VITE_DEMO_MODE === 'true'
  })

  useEffect(() => {
    sessionStorage.setItem(
      'showDemoWarning',
      showDemoWarning ? 'true' : 'false'
    )
  }, [showDemoWarning])

  if (!showDemoWarning) {
    return null
  }

  return (
    <div className="max-w-[600px] flex text-xs gap-x-2 border-1 surface-0-bg border border-surface-0-stroke rounded p-2">
      <div>
        <FluentInfo16Regular className="flex-shrink-0" />
      </div>
      <p>
        <span className="font-semibold">Demo mode</span> Cupdate is in read-only
        demo mode. GitHub Pages may fail to route some pages.
      </p>
      <div>
        <button
          type="button"
          className="flex-shrink-0 cursor-pointer"
          onClick={() => setShowDemoWarning(false)}
        >
          <FluentDismiss16Regular />
        </button>
      </div>
    </div>
  )
}
