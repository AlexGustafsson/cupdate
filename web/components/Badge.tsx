import { HTMLAttributes, type JSX } from 'react'

export type BadgeProps = {
  label: string
  color?: string | { light: string; dark: string }
  disabled?: boolean
}

export function Badge({
  label,
  color,
  disabled,
  className,
  ...rest
}: Omit<HTMLAttributes<HTMLSpanElement>, 'color'> & BadgeProps): JSX.Element {
  let backgroundColor = '#CC5889'
  if (typeof color === 'string') {
    backgroundColor = color
  } else if (color !== undefined) {
    backgroundColor = `light-dark(${color.light}, ${color.dark})`
  }

  return (
    <span
      {...rest}
      className={
        `rounded-md px-1 sm:px-2 py-1 text-xs text-nowrap text-white dark:text-[#dddddd] h-fit ${disabled ? 'opacity-50 hover:opacity-70' : ''}` +
        className
      }
      style={{ backgroundColor }}
    >
      {label}
    </span>
  )
}
