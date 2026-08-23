import type { HTMLAttributes, JSX } from 'react'
import type { Color } from '../tags'

export type BadgeProps = {
  label: string
  color?: Color
  disabled?: boolean
}

export function Badge({
  label,
  disabled,
  color,
  className,
  ...rest
}: HTMLAttributes<HTMLSpanElement> & BadgeProps): JSX.Element {
  return (
    <span
      {...rest}
      className={`rounded-md px-1 sm:px-2 py-1 text-xs text-nowrap h-fit ${disabled ? 'opacity-50 hover:opacity-70' : ''}${className}`}
      style={{
        backgroundColor: `var(--color-${color})`,
        color: `contrast-color(var(--color-${color}))`,
      }}
    >
      {label}
    </span>
  )
}
