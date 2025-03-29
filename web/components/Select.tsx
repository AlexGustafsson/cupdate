import type { JSX, PropsWithChildren, SelectHTMLAttributes } from 'react'

export function Select({
  children,
  ...rest
}: PropsWithChildren<SelectHTMLAttributes<HTMLSelectElement>>): JSX.Element {
  return (
    <div className="relative border border-[#e5e5e5] dark:border-[#333333] rounded-sm transition-colors focus:border-[#f0f0f0] dark:focus:border-[#333333] hover:border-[#f0f0f0] dark:hover:border-[#333333] shadow-xs focus:shadow-md bg-white dark:bg-[#1e1e1e] dark:hover:bg-[#262626]">
      <select
        {...rest}
        className="pl-3 pr-8 py-2 text-sm cursor-pointer appearance-none focus:bg-[#f5f5f5] dark:focus:bg-[#262626]"
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
