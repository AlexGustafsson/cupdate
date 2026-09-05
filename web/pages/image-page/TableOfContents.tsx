import type { JSX } from 'react'
import { FluentChevronDown16Regular } from '../../components/icons/fluent-chevron-down-16-regular'
import { FluentChevronUp16Regular } from '../../components/icons/fluent-chevron-up-16-regular'

export type TableOfContentsProps = {
  className?: string
  onClick: (id: string) => void
  onToggleCollapsed: (id: string) => void
  items: {
    id: string
    icon: JSX.Element
    label: string
    collapsed?: boolean
    disabled?: boolean
  }[]
}

export function TableOfContents({
  className,
  items,
  onClick,
  onToggleCollapsed,
}: TableOfContentsProps): JSX.Element {
  return (
    <aside
      className={`w-min rounded-lg bg-surface-1-bg shadow p-4 pt-2 ${className ?? ''}`}
    >
      <p className="font-semibold text-sm whitespace-pre text-center">
        Table of contents
      </p>
      <ul className="space-y-1 my-2">
        {items.map(({ id, icon, label, collapsed, disabled }) => (
          <li key={id} className="flex items-center">
            <button
              type="button"
              role="tab"
              disabled={disabled}
              onClick={() => onClick(id)}
              className="w-full btn-flat flex gap-x-1 items-center disabled:bg-[initial] disabled:text-surface-1-fg-disabled font-semibold text-sm p-1"
              title={disabled ? 'Unavailable' : undefined}
            >
              {icon}
              <p className="flex-grow text-left whitespace-pre">{label}</p>
            </button>
            {/* Collapse */}
            {!disabled && collapsed !== undefined && (
              <button
                type="button"
                className="btn-flat btn-medium btn-square disabled:bg-[initial]"
                disabled={disabled}
                onClick={() => onToggleCollapsed(id)}
                title={collapsed ? 'Expand' : 'Collapse'}
              >
                {collapsed ? (
                  <FluentChevronDown16Regular />
                ) : (
                  <FluentChevronUp16Regular />
                )}
              </button>
            )}
          </li>
        ))}
      </ul>
    </aside>
  )
}
