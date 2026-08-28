import {
  type JSX,
  type PropsWithChildren,
  useEffect,
  useRef,
  useState,
} from 'react'

import type { Filter } from '../hooks'
import { compareTags, type Tag } from '../tags'
import { Badge } from './Badge'

const IOS = [
  'iPad Simulator',
  'iPhone Simulator',
  'iPod Simulator',
  'iPad',
  'iPhone',
  'iPod',
].includes(navigator.platform)

export function TagSelect({
  tags,
  filter,
  onChange,
  className,
}: PropsWithChildren<{
  tags: Tag[]
  filter: Filter
  onChange: React.Dispatch<React.SetStateAction<Filter>>
  className?: string
}>): JSX.Element {
  const menuRef = useRef<HTMLDivElement>(null)
  const selectRef = useRef<HTMLDivElement>(null)

  const [isOpen, setIsOpen] = useState(false)

  useEffect(() => {
    if (isOpen) {
      const handler = (e: MouseEvent) => {
        if (!menuRef.current || !selectRef.current) {
          return
        }

        const outsideOfMenu =
          e.offsetX < menuRef.current.offsetLeft ||
          e.offsetX >
            menuRef.current.offsetLeft + menuRef.current.offsetWidth ||
          e.offsetY < menuRef.current.offsetTop ||
          e.offsetY > menuRef.current.offsetTop + menuRef.current.offsetHeight

        const outsideOfSelect =
          e.offsetX < selectRef.current.clientLeft ||
          e.offsetX >
            selectRef.current.clientLeft + selectRef.current.clientWidth ||
          e.offsetY < selectRef.current.clientTop ||
          e.offsetY >
            selectRef.current.clientTop + selectRef.current.clientHeight

        if (outsideOfMenu && outsideOfSelect) {
          setIsOpen(false)
        }
      }
      document.addEventListener('mousedown', handler)
      return () => document.removeEventListener('mousedown', handler)
    }
  }, [isOpen])

  useEffect(() => {
    if (isOpen) {
      const handler = (e: KeyboardEvent) => {
        if (e.key === 'Escape') {
          setIsOpen(false)
        }
      }
      document.addEventListener('keydown', handler)
      return () => document.removeEventListener('keydown', handler)
    }
  }, [isOpen])

  // Use the nice native multi-select input on iOS
  if (IOS) {
    return (
      <div
        className={`relative border border-surface-1-stroke rounded-sm transition-colors focus:border-surface-1-stroke hover:border-surface-1-stroke shadow-xs focus:shadow-md bg-surface-1-bg ${className || ''}`}
      >
        <select
          multiple
          className="pl-3 pr-8 py-2 text-sm cursor-pointer appearance-none w-full"
          value={filter.tags}
          onChange={(e) =>
            onChange({
              tags: Array.from(
                e.target.selectedOptions,
                (option) => option.value
              ),
            })
          }
        >
          {tags
            .toSorted((a, b) => compareTags(a.name, b.name))
            .map((x) => (
              <option key={x.name} value={x.name}>
                {x.name}
              </option>
            ))}
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

  return (
    <div
      ref={selectRef}
      onMouseDown={() => setIsOpen(true)}
      onKeyDown={(e) => {
        if (e.key === ' ' || e.key === 'enter') {
          setIsOpen(true)
          e.preventDefault()
          return false
        }
      }}
      role="menu"
      tabIndex={0}
      className={`pl-3 pr-8 py-2 relative border border-surface-1-stroke rounded-sm transition-colors focus:border-surface-1-stroke hover:border-surface-1-stroke shadow-xs focus:shadow-md bg-surface-1-bg cursor-pointer ${className || ''}`}
    >
      <p className="text-sm">
        {filter.tags.length > 0 ? `${filter.tags.length} selected` : 'Tags'}
      </p>
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
      {isOpen && (
        <div
          ref={menuRef}
          className="absolute -top-4 -left-4 p-2 z-150 text-black dark:text-[#dddddd]"
        >
          <div className="flex max-h-64 overflow-y-auto flex-col gap-y-2 py-2 px-3 pr-6 bg-surface-2-bg border border-surface-2-stroke rounded-lg w-max shadow">
            {tags
              .toSorted((a, b) =>
                compareTags(
                  a.name,
                  b.name,
                  filter.tags.includes(a.name),
                  filter.tags.includes(b.name)
                )
              )
              .map((x) => (
                <label
                  key={x.name}
                  className="cursor-pointer flex shrink-0 items-center"
                >
                  <input
                    role="menuitemcheckbox"
                    aria-checked={filter.tags.includes(x.name)}
                    type="checkbox"
                    checked={filter.tags.includes(x.name)}
                    onChange={(e) =>
                      onChange((current) =>
                        e.target.checked
                          ? { tags: [...current.tags, x.name] }
                          : { tags: current.tags.filter((y) => y !== x.name) }
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
