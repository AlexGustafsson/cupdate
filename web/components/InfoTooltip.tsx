import { PropsWithChildren } from 'react'

import { FluentInfo16Regular } from './icons/fluent-info-16-regular'

export function InfoTooltip({
  children,
  icon,
}: PropsWithChildren<{ icon?: JSX.Element }>): JSX.Element {
  return (
    <span className="inline-block align-middle relative group ml-1">
      {icon || <FluentInfo16Regular />}
      <span
        role="tooltip"
        className="absolute invisible group-hover:visible bottom-full p-2 left-2/4 -translate-x-2/4 z-50"
      >
        <div className="p-2 bg-white dark:bg-[#292929] border-solid border-[1px] border-[#d9d9d9] dark:border-[#454545] rounded w-60 text-xs text-left font-normal shadow-around">
          {children}
        </div>
      </span>
    </span>
  )
}