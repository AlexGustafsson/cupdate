import type { JSX } from 'react'

export type ToastProps = {
  title: string
  body: string

  onPrimaryAction?: () => void
  primaryAction: string

  onSecondaryAction?: () => void
  secondaryAction?: string
}

export function Toast({
  title,
  body,
  onPrimaryAction,
  primaryAction,
  onSecondaryAction,
  secondaryAction,
}: ToastProps): JSX.Element {
  return (
    <div className="starting:opacity-0 transition-opacity rounded-sm bg-surface-2-bg p-3 shadow-md text-sm max-w-[300px] border border-surface-2-stroke">
      <p className="font-semibold">{title}</p>
      <p className="mt-1">{body}</p>
      <div className="flex items-center justify-end gap-x-2 mt-2">
        {secondaryAction && (
          <button
            type="button"
            className="btn-flat p-1 text-accent hover:underline cursor-pointer"
            onClick={() => onSecondaryAction?.()}
          >
            {secondaryAction}
          </button>
        )}
        <button
          type="button"
          className="btn-flat p-1 text-accent hover:underline cursor-pointer"
          onClick={() => onPrimaryAction?.()}
        >
          {primaryAction}
        </button>
      </div>
    </div>
  )
}
