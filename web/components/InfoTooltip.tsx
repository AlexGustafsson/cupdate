import type { JSX, PropsWithChildren } from 'react'

import { FluentInfo16Regular } from './icons/fluent-info-16-regular'

export function InfoTooltip({
  children,
  icon,
  className,
}: PropsWithChildren<{ icon?: JSX.Element; className?: string }>): JSX.Element {
  return (
    <span
      className={`inline-block align-middle relative group/tooltip ${className ?? ''}`}
    >
      {icon || <FluentInfo16Regular />}
      <span
        role="tooltip"
        className="starting:opacity-0 transition-opacity absolute hidden group-hover/tooltip:block bottom-full p-2 left-2/4 -translate-x-2/4 z-200"
      >
        <div className="p-2 bg-surface-2-bg border-solid border-[1px] border-surface-2-stroke rounded-sm w-60 text-xs text-left font-normal shadow-around">
          {children}
        </div>
      </span>
    </span>
  )
}
