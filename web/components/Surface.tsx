import {
  type JSX,
  type PropsWithChildren,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react'
import { FluentFullScreenMaximize16Regular } from './icons/fluent-full-screen-maximize-16-regular'
import { FluentAdd16Regular } from './icons/fluent-plus-16-regular'
import { FluentSubtract16Regular } from './icons/fluent-subtract-16-regular'

export function Surface({
  children,
}: PropsWithChildren<Record<never, never>>): JSX.Element {
  const surfaceRef = useRef<HTMLDivElement>(null)
  const contentRef = useRef<HTMLDivElement>(null)

  const [offset, setOffset] = useState<{ x: number; y: number }>({ x: 0, y: 0 })
  const [scale, setScale] = useState<number>(1.0)
  const [isDragging, setIsDragging] = useState(false)

  const onPointerMove = useCallback((e: PointerEvent) => {
    setOffset((current) => ({
      x: current.x + e.movementX,
      y: current.y + e.movementY,
    }))
  }, [])
  const onPointerUp = useCallback(
    (e: PointerEvent) => {
      document.removeEventListener('pointermove', onPointerMove)

      setIsDragging(false)
    },
    [onPointerMove]
  )
  const onPointerDown = useCallback(
    (e: PointerEvent) => {
      if (e.buttons !== 1) {
        return
      }

      document.addEventListener('pointermove', onPointerMove)
      document.addEventListener('pointerup', onPointerUp)

      setIsDragging(true)
    },
    [onPointerMove, onPointerUp]
  )

  const onZoom = useCallback((delta: number) => {
    setScale((current) => Math.min(Math.max(current + delta * 0.1, 0.4), 1))
  }, [])

  const onCenter = useCallback(() => {
    if (!contentRef.current || !surfaceRef.current) {
      return
    }

    const surfaceWidth = surfaceRef.current.offsetWidth
    const surfaceHeight = surfaceRef.current.offsetHeight

    const contentWidth = contentRef.current.offsetWidth
    const contentHeight = contentRef.current.offsetHeight

    const scale = Math.min(
      surfaceWidth / contentWidth,
      surfaceHeight / contentHeight
    )

    setScale(Math.min(Math.max(scale, 0.3), 1))
    setOffset({
      x: surfaceWidth / 2 - contentWidth / 2,
      y: surfaceHeight / 2 - contentHeight / 2,
    })
  }, [])

  const onWheel = useCallback((e: WheelEvent) => {
    setScale((current) =>
      Math.min(Math.max(current - e.deltaY * 0.001, 0.3), 1)
    )
    e.preventDefault()
  }, [])

  useEffect(() => {
    surfaceRef.current?.addEventListener('pointerdown', onPointerDown)
    surfaceRef.current?.addEventListener('wheel', onWheel)

    return () => {
      surfaceRef.current?.removeEventListener('pointerdown', onPointerDown)
      surfaceRef.current?.removeEventListener('wheel', onWheel)
    }
  }, [onPointerDown, onWheel])

  // Center on child re-size (which should only happen a couple of times the
  // first few renders)
  const onContentResize = useCallback(
    (entries: ResizeObserverEntry[], observer: ResizeObserver) => {
      onCenter()
    },
    [onCenter]
  )

  useEffect(() => {
    const observer = new ResizeObserver(onContentResize)
    if (contentRef.current) {
      observer.observe(contentRef.current)
    }
    return () => {
      observer.disconnect()
    }
  }, [onContentResize])

  return (
    <div
      ref={surfaceRef}
      className={`relative w-full h-full overflow-hidden select-none ${isDragging ? 'cursor-grabbing' : 'cursor-grab'} touch-none`}
    >
      <div className="absolute left-0 bottom-0 m-2 rounded bg-white dark:bg-[#1e1e1e] flex flex-col p-2 z-50 shadow-md gap-y-2">
        <button type="button" onClick={() => onZoom(1)}>
          <FluentAdd16Regular />
        </button>
        <button type="button" onClick={() => onZoom(-1)}>
          <FluentSubtract16Regular />
        </button>
        <button type="button" onClick={() => onCenter()}>
          <FluentFullScreenMaximize16Regular />
        </button>
      </div>
      <div
        ref={contentRef}
        className="pointer-events-none w-fit h-fit"
        style={{
          transform: `translate3d(${offset.x}px, ${offset.y}px, 0) scale(${scale}, ${scale})`,
        }}
      >
        {children}
      </div>
    </div>
  )
}
