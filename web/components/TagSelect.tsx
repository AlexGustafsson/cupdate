import { PropsWithChildren, useRef, useState } from 'react'

import { Tag } from '../tags'
import { Badge } from './Badge'

export function TagSelect({
  tags,
  filter,
  onChange,
}: PropsWithChildren<{
  tags: Tag[]
  filter: string[]
  onChange: React.Dispatch<React.SetStateAction<string[]>>
}>): JSX.Element {
  const menuRef = useRef<HTMLDivElement>(null)

  const [isOpen, setIsOpen] = useState(false)

  return (
    <div
      ref={menuRef}
      onClick={() => setIsOpen(true)}
      onBlur={() => setIsOpen(false)}
      tabIndex={0}
      className="relative pl-3 pr-8 py-2 relative border border-gray-200 rounded transition-colors duration-300 ease focus:outline-none focus:border-gray-300 hover:border-gray-300 shadow-sm focus:shadow-md bg-white hover:bg-[#fafafa]"
    >
      <p className="text-sm appearance-none cursor-pointer">
        {filter.length > 0 ? `${filter.length} selected` : 'Tags'}
      </p>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        strokeWidth="1.2"
        stroke="currentColor"
        className="h-5 w-5 ml-1 absolute top-2.5 right-2.5 text-gray-700 pointer-events-none"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M8.25 15 12 18.75 15.75 15m-7.5-6L12 5.25 15.75 9"
        />
      </svg>
      {isOpen && (
        <div className="absolute group-hover:visible -top-4 -left-4 p-2 z-50 text-black dark:text-[#dddddd]">
          <div className="flex flex-col gap-y-2 py-2 px-3 pr-6 bg-white dark:bg-[#292929] border-solid border-[1px] border-[#d0d0d0]/95 dark:border-[#454545] rounded-lg w-max shadow">
            {tags.map((x) => (
              <label className="cursor-pointer">
                <input
                  type="checkbox"
                  checked={filter.includes(x.name)}
                  onChange={(e) =>
                    onChange((current) =>
                      e.target.checked
                        ? [...current, x.name]
                        : current.filter((y) => y !== x.name)
                    )
                  }
                  className="scale-125 cursor-pointer"
                />
                <Badge
                  title={x.description}
                  label={x.name}
                  color={x.color}
                  className="ml-2"
                />
              </label>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
