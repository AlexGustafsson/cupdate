import { PropsWithChildren, SelectHTMLAttributes } from 'react'

export function Select({
  children,
  ...rest
}: PropsWithChildren<SelectHTMLAttributes<HTMLSelectElement>>): JSX.Element {
  return (
    <div className="relative border border-gray-200 rounded transition-colors duration-300 ease focus:outline-none focus:border-gray-300 hover:border-gray-300 shadow-sm focus:shadow-md bg-white hover:bg-[#fafafa]">
      <select
        {...rest}
        className="pl-3 pr-8 py-2 bg-transparent placeholder:text-gray-300 text-sm appearance-none cursor-pointer"
      >
        {children}
      </select>
      <svg
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
