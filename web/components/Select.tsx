import type { JSX, PropsWithChildren, SelectHTMLAttributes } from 'react'

export function Select({
  children,
  className,
  ...rest
}: PropsWithChildren<SelectHTMLAttributes<HTMLSelectElement>>): JSX.Element {
  return (
    <div className={`select relative ${className || ''}`}>
      <select
        {...rest}
        className="pl-3 pr-8 py-2 text-sm cursor-pointer w-full"
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
