import type { JSX } from 'react'
import { useMountedDuration } from '../../hooks/useMountedDuration'

export function DashboardSkeleton(): JSX.Element | null {
  const mountedDuration = useMountedDuration()

  if (mountedDuration < 500) {
    return null
  }

  return (
    <div className="w-full transition-discrete opacity-0 starting:opacity-0 animate-pulse relative">
      <div className="flex flex-col max-w-[800px] w-full px-4 fixed left-1/2 -translate-x-1/2">
        <div className="mt-12 w-full h-10 rounded bg-gray-200 max-w-[500px] self-center" />
        <div className="mt-20 space-y-10">
          {[...Array(10)].map((_, i) => (
            <div className="flex" key={i.toString()}>
              <div className="rounded w-[48px] h-[48px] bg-gray-200" />
              <div className="ml-4 flex-grow space-y-2">
                <div className="rounded w-2/5 h-6 bg-gray-200" />
                <div className="rounded w-3/5 h-6 bg-gray-200" />
              </div>
              <div className="rounded w-1/5 h-6 bg-gray-200" />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
