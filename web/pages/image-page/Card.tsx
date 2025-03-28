import { type JSX, useEffect, useState } from 'react'
import { FluentChevronDown20Regular } from '../../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../../components/icons/fluent-chevron-up-20-regular'

export type TabProps = {
  label: string
  icon?: JSX.Element
  disabled?: boolean
  active?: boolean
  onClick?: React.MouseEventHandler<HTMLButtonElement> | undefined
}

export function Tab({
  icon,
  label,
  disabled,
  active,
  onClick,
}: TabProps): JSX.Element {
  return (
    <div
      className={`flex-shrink-0 px-1 py-2 border-b-2 ${active ? 'border-blue-400 dark:border-blue-700' : 'border-transparent'}`}
    >
      <button
        type="button"
        onClick={onClick}
        disabled={disabled}
        className="flex items-center gap-x-2 font-semibold text-sm rounded p-1 enabled:hover:bg-[#f5f5f5] dark:enabled:hover:bg-[#262626] enabled:cursor-pointer"
      >
        {icon}
        <p>{label}</p>
      </button>
    </div>
  )
}

export type CardProps = {
  /** Local storage key used to persist the state of the card. */
  persistenceKey?: string
  tabs: Tab[]
}

export interface Tab {
  icon?: JSX.Element
  label: string
  content: JSX.Element
}

export function Card({ persistenceKey, tabs }: CardProps): JSX.Element {
  const [showContent, setShowContent] = useState(() => {
    // Try to load persisted state
    if (persistenceKey) {
      const item = localStorage.getItem(`cupdate-card-state-${persistenceKey}`)
      if (item === 'false') {
        return false
      }
    }

    return true
  })

  const [selectedTabIndex, setSelectedTabIndex] = useState(0)

  // Persist state
  useEffect(() => {
    if (persistenceKey) {
      localStorage.setItem(
        `cupdate-card-state-${persistenceKey}`,
        showContent ? 'true' : 'false'
      )
    }
  }, [persistenceKey, showContent])

  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] shadow">
      {/* Header */}
      <div
        className={`sticky top-[64px] bg-white dark:bg-[#1e1e1e] flex items-center w-full ${showContent ? 'rounded-t-lg border-b border-[#e5e5e5] dark:border-[#333333] mb-2' : 'rounded-lg'}`}
      >
        {/* Tabs */}
        <div className="flex items-center flex-grow px-2 max-w-full overflow-auto">
          {tabs.map((tab, i) => (
            <Tab
              key={tab.label}
              icon={tab.icon}
              label={tab.label}
              disabled={
                i === selectedTabIndex || tabs.length === 1 || !showContent
              }
              active={tabs.length > 1 && i === selectedTabIndex && showContent}
              onClick={() => setSelectedTabIndex(i)}
            />
          ))}
        </div>

        {/* Controls */}
        <div className="flex items-center p-2">
          <button
            type="button"
            onClick={() => setShowContent((current) => !current)}
            className="flex items-center gap-x-2 font-semibold text-sm rounded p-1 enabled:hover:bg-[#f5f5f5] dark:enabled:hover:bg-[#262626] enabled:cursor-pointer"
          >
            {showContent ? (
              <FluentChevronUp20Regular />
            ) : (
              <FluentChevronDown20Regular />
            )}
          </button>
        </div>
      </div>

      {/* Content */}
      {showContent && (
        <div className="p-4 pt-2">{tabs[selectedTabIndex].content}</div>
      )}
    </div>
  )
}
