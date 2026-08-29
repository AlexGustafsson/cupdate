import { type JSX, useEffect, useState } from 'react'
import { FluentChevronDown16Regular } from '../../components/icons/fluent-chevron-down-16-regular'
import { FluentChevronUp16Regular } from '../../components/icons/fluent-chevron-up-16-regular'
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
        >
          <button
            type="button"
            className="btn-flat btn-small btn-square ml-1"
            tabIndex={0}
          >
            <FluentOpen16Regular />
          </button>
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
      className={`flex-shrink-0 px-1 py-2 border-b-2 ${active ? 'border-accent' : 'border-transparent'}`}
    >
      <button
        type="button"
        role="tab"
        onClick={onClick}
        disabled={disabled}
        className="btn-flat disabled:cursor-default disabled:bg-[initial] font-semibold text-sm p-1"
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
    <div className="rounded-lg bg-surface-1-bg shadow">
      {/* Header */}
      <div
        className={`sticky z-50 top-[64px] bg-surface-1-bg flex items-center w-full ${showContent ? 'rounded-t-lg border-b border-surface-1-stroke mb-2' : 'rounded-lg'}`}
      >
        {/* Tabs */}
        <div className="flex items-center flex-grow px-2 max-w-full overflow-auto">
          {tabs.map((tab, i) => (
            <Tab
              key={tab.label}
              icon={tab.icon}
              label={tab.label}
              action={tab.action}
              disabled={showContent ? i === selectedTabIndex : false}
              active={tabs.length > 1 && i === selectedTabIndex && showContent}
              onClick={() => {
                setSelectedTabIndex(i)
                setShowContent(true)
              }}
            />
          ))}
        </div>

        {/* Collapse */}
        <button
          type="button"
          onClick={() => setShowContent((current) => !current)}
          className="btn-flat btn-medium btn-square mr-2"
          tabIndex={0}
        >
          {showContent ? (
            <FluentChevronUp16Regular />
          ) : (
            <FluentChevronDown16Regular />
          )}
        </button>
      </div>

      {/* Content */}
      {showContent && (
        <div className="p-4 pt-2">{tabs[selectedTabIndex].content}</div>
      )}
    </div>
  )
}
