import { type JSX, useEffect, useState } from 'react'
import { FluentChevronDown20Regular } from '../../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../../components/icons/fluent-chevron-up-20-regular'
import { FluentOpen16Regular } from '../../components/icons/fluent-open-16-regular'

export type TabActionProps = {
  action: TabAction
}

export function TabAction({ action }: TabActionProps): JSX.Element | undefined {
  switch (action.type) {
    case 'external-link':
      return (
        <a
          target="_blank"
          rel="noreferrer"
          href={action.href}
          title={action.title}
          className="ml-1 rounded p-1 focus:bg-[#f5f5f5] dark:focus:bg-[#262626] hover:bg-[#f5f5f5] dark:hover:bg-[#262626] cursor-pointer"
          tabIndex={0}
        >
          <FluentOpen16Regular />
        </a>
      )
    default:
      return undefined
  }
}

export type TabProps = {
  label: string
  icon?: JSX.Element
  action?: TabAction
  disabled?: boolean
  active?: boolean
  onClick?: React.MouseEventHandler<HTMLButtonElement> | undefined
}

export function Tab({
  icon,
  label,
  action,
  disabled,
  active,
  onClick,
}: TabProps): JSX.Element {
  return (
    <div
      role="tablist"
      className={`flex-shrink-0 px-1 py-2 border-b-2 ${active ? 'border-blue-400 dark:border-blue-700' : 'border-transparent'}`}
    >
      <button
        type="button"
        role="tab"
        onClick={onClick}
        disabled={disabled}
        className="flex items-center font-semibold text-sm rounded p-1 enabled:hover:bg-[#f5f5f5] enabled:focus:bg-[#f5f5f5] dark:enabled:hover:bg-[#262626] dark:enabled:focus:bg-[#262626] enabled:cursor-pointer"
        tabIndex={0}
      >
        {icon}
        <p className={icon ? 'ml-2' : ''}>{label}</p>
        {action && <TabAction action={action} />}
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
  action?: TabAction
  content: JSX.Element
}

type TabAction = {
  type: 'external-link'
  href: string
  title?: string
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
        className={`sticky z-50 top-[64px] bg-white dark:bg-[#1e1e1e] flex items-center w-full ${showContent ? 'rounded-t-lg border-b border-[#e5e5e5] dark:border-[#333333] mb-2' : 'rounded-lg'}`}
      >
        {/* Tabs */}
        <div className="flex items-center flex-grow px-2 max-w-full overflow-auto">
          {tabs.map((tab, i) => (
            <Tab
              key={tab.label}
              icon={tab.icon}
              label={tab.label}
              action={tab.action}
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
            className="flex items-center gap-x-2 font-semibold text-sm rounded p-1 enabled:focus:bg-[#f5f5f5] dark:enabled:focus:bg-[#262626] enabled:hover:bg-[#f5f5f5] dark:enabled:hover:bg-[#262626] enabled:cursor-pointer"
            tabIndex={0}
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
