import type { JSX } from 'react'
import { useMountedDuration } from '../../hooks/useMountedDuration'

export function ImageSkeleton(): JSX.Element | null {
  const mountedDuration = useMountedDuration()

  if (mountedDuration < 500) {
    return null
  }

  return (
    <div className="w-full transition-discrete opacity-0 starting:opacity-0 animate-pulse relative">
      <div className="flex flex-col max-w-[800px] w-full px-4 fixed left-1/2 -translate-x-1/2">
        <div className="mt-6 rounded w-[90px] h-[90px] bg-gray-200 self-center" />
        <div className="mt-4 rounded w-2/5 h-6 bg-gray-200 self-center" />
        <div className="mt-4 rounded w-1/5 h-6 bg-gray-200 self-center" />
        <div className="mt-4 rounded w-3/5 h-6 bg-gray-200 self-center" />
        <div className="mt-4 rounded w-2/5 h-6 bg-gray-200 self-center" />
      </div>
    </div>
  )
}
