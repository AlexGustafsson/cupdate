import { HTMLAttributes } from 'react'

export type BadgeProps = {
  label: string
  color: string
  disabled?: boolean
}

export function Badge({
  label,
  color,
  disabled,
  className,
  ...rest
}: HTMLAttributes<HTMLSpanElement> & BadgeProps): JSX.Element {
  return (
    <span
      {...rest}
      className={
        `rounded-full px-2 py-1 text-xs text-nowrap ${color} ${disabled ? 'opacity-50' : ''} m-1 ` +
        className
      }
    >
      {label}
    </span>
  )
}
