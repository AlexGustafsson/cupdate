import {
  type JSX,
  type PropsWithChildren,
  useEffect,
  useRef,
  useState,
} from 'react'

export type MenuProps = {
  icon: JSX.Element
}

export function Menu({
  icon,
  children,
}: PropsWithChildren<MenuProps>): JSX.Element {
  const openRef = useRef<HTMLButtonElement>(null)
  const [isOpen, setIsOpen] = useState<boolean>(false)

  useEffect(() => {
    if (isOpen) {
      const handle = (e: MouseEvent) => {
        if (
          e.target === openRef.current ||
          openRef.current?.contains(e.target as Node)
        ) {
          return
        }

        setIsOpen(false)
      }
      document.addEventListener('click', handle)
      return () => document.removeEventListener('click', handle)
    }
  }, [isOpen])

  return (
    <div>
      <button
        type="button"
        ref={openRef}
        onClick={() => setIsOpen(true)}
        className="menu-button"
      >
        {icon}
      </button>
      {isOpen && <ul className="menu-container">{children}</ul>}
    </div>
  )
}
