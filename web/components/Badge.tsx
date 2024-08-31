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
  ...rest
}: HTMLAttributes<HTMLSpanElement> & BadgeProps): JSX.Element {
  return (
    <span
      key={label}
      className={`rounded-full px-2 py-1 text-xs text-nowrap cursor-pointer ${color} ${disabled ? 'opacity-50' : ''} m-1`}
      {...rest}
    >
      {label}
    </span>
  )
}
