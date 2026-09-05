import { type JSX, useState } from 'react'
import { FluentChevronDown16Regular } from './icons/fluent-chevron-down-16-regular'
import { FluentChevronUp16Regular } from './icons/fluent-chevron-up-16-regular'
import { FluentOpen16Regular } from './icons/fluent-open-16-regular'

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
      className={`flex items-center flex-shrink-0 px-1 py-2 border-b-2 ${active ? 'border-accent' : 'border-transparent'}`}
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
      </button>
      {action && <TabAction action={action} />}
    </div>
  )
}

export type CardProps = {
  id?: string
  tabs: Tab[]
  collapsed: boolean
  onToggleCollapsed: () => void
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

export function Card({
  id,
  tabs,
  collapsed,
  onToggleCollapsed,
}: CardProps): JSX.Element {
  const [selectedTabIndex, setSelectedTabIndex] = useState(0)

  return (
    <div id={id} className="rounded-lg bg-surface-1-bg shadow scroll-mt-[64px]">
      {/* Header */}
      <div
        className={`sticky z-50 top-[64px] bg-surface-1-bg flex items-center w-full ${collapsed ? 'rounded-lg' : 'rounded-t-lg border-b border-surface-1-stroke mb-2'}`}
      >
        {/* Tabs */}
        <div className="flex items-center flex-grow px-2 max-w-full overflow-auto">
          {tabs.map((tab, i) => (
            <Tab
              key={tab.label}
              icon={tab.icon}
              label={tab.label}
              action={tab.action}
              disabled={collapsed ? false : i === selectedTabIndex}
              active={tabs.length > 1 && i === selectedTabIndex && !collapsed}
              onClick={() => {
                setSelectedTabIndex(i)
                onToggleCollapsed()
              }}
            />
          ))}
        </div>

        {/* Collapse */}
        <button
          type="button"
          onClick={() => onToggleCollapsed()}
          className="btn-flat btn-medium btn-square mr-2"
          tabIndex={0}
          title={collapsed ? 'Expand' : 'Collapse'}
        >
          {collapsed ? (
            <FluentChevronDown16Regular />
          ) : (
            <FluentChevronUp16Regular />
          )}
        </button>
      </div>

      {/* Content */}
      {!collapsed && (
        <div className="p-4 pt-2">{tabs[selectedTabIndex].content}</div>
      )}
    </div>
  )
}
