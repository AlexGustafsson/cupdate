import type { JSX, PropsWithChildren, SelectHTMLAttributes } from 'react'

export function Select({
  children,
  className,
  ...rest
}: PropsWithChildren<SelectHTMLAttributes<HTMLSelectElement>>): JSX.Element {
  return (
    <div
      className={`relative border border-surface-1-stroke rounded-sm transition-colors focus:border-surface-1-stroke hover:border-surface-1-stroke shadow-xs focus:shadow-md bg-surface-1-bg ${className || ''}`}
    >
      <select
        {...rest}
        className="pl-3 pr-8 py-2 text-sm cursor-pointer appearance-none focus:bg-surface-2-bg w-full"
      >
        {children}
      </select>
      <svg
        role="img"
        aria-label="icon"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        strokeWidth="1.2"
        stroke="currentColor"
        className="h-5 w-5 ml-1 absolute top-2.5 right-2.5 pointer-events-none"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M8.25 15 12 18.75 15.75 15m-7.5-6L12 5.25 15.75 9"
        />
      </svg>
    </div>
  )
}
