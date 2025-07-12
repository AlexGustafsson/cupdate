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

const MAX_SCALE = 1.0
const MIN_SCALE = 0.4

export function Surface({
  children,
}: PropsWithChildren<Record<never, never>>): JSX.Element {
  const surfaceRef = useRef<HTMLDivElement>(null)
  const contentRef = useRef<HTMLDivElement>(null)
  const controlsRef = useRef<HTMLDivElement>(null)

  const [offset, setOffset] = useState<{ x: number; y: number }>({ x: 0, y: 0 })
  const [scale, setScale] = useState<number>(1.0)
  const [isDragging, setIsDragging] = useState(false)
  const [hasInteracted, setHasInteracted] = useState(false)

  const scaleToFit =
    surfaceRef.current && contentRef.current && controlsRef.current
      ? Math.min(
          surfaceRef.current.offsetWidth / contentRef.current.offsetWidth,
          (surfaceRef.current.offsetHeight -
            controlsRef.current.offsetHeight * 2) /
            contentRef.current.offsetHeight
        )
      : MIN_SCALE

  const onPointerMove = useCallback((e: PointerEvent) => {
    setOffset((current) => ({
      x: current.x + e.movementX,
      y: current.y + e.movementY,
    }))
    setHasInteracted(true)
  }, [])
  const onPointerUp = useCallback(() => {
    document.removeEventListener('pointermove', onPointerMove)

    setIsDragging(false)
    setHasInteracted(true)
  }, [onPointerMove])
  const onPointerDown = useCallback(
    (e: PointerEvent) => {
      if (e.buttons !== 1) {
        return
      }

      document.addEventListener('pointermove', onPointerMove)
      document.addEventListener('pointerup', onPointerUp, { once: true })

      setIsDragging(true)
      setHasInteracted(true)
    },
    [onPointerMove, onPointerUp]
  )

  const onZoom = useCallback(
    (delta: number) => {
      setScale((current) =>
        Math.min(Math.max(current + delta * 0.1, scaleToFit), MAX_SCALE)
      )
      setHasInteracted(true)
    },
    [scaleToFit]
  )

  const onCenter = useCallback(() => {
    if (!contentRef.current || !surfaceRef.current) {
      return
    }

    const surfaceWidth = surfaceRef.current.offsetWidth
    const surfaceHeight = surfaceRef.current.offsetHeight

    const contentWidth = contentRef.current.offsetWidth
    const contentHeight = contentRef.current.offsetHeight

    setScale(scaleToFit)
    setOffset({
      x: surfaceWidth / 2 - contentWidth / 2,
      y: surfaceHeight / 2 - contentHeight / 2,
    })
    setHasInteracted(false)
  }, [scaleToFit])

  const onWheel = useCallback(
    (e: WheelEvent) => {
      setScale((current) =>
        Math.min(Math.max(current - e.deltaY * 0.001, scaleToFit), MAX_SCALE)
      )
      setHasInteracted(true)
      e.preventDefault()
    },
    [scaleToFit]
  )

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
  const onContentResize = useCallback(() => {
    if (!hasInteracted) {
      onCenter()
    }
  }, [hasInteracted, onCenter])

  useEffect(() => {
    const observer = new ResizeObserver(onContentResize)
    if (contentRef.current) {
      observer.observe(contentRef.current)
      window.addEventListener('resize', onContentResize)
    }
    return () => {
      observer.disconnect()
      window.removeEventListener('resize', onContentResize)
    }
  }, [onContentResize])
  return (
    <div ref={surfaceRef} className="relative w-full h-full select-none">
      <div
        className={`overflow-hidden w-full h-full ${isDragging ? 'cursor-grabbing' : 'cursor-grab'} touch-none`}
      >
        <div
          ref={contentRef}
          className="w-fit h-fit"
          style={{
            transform: `translate3d(${offset.x}px, ${offset.y}px, 0) scale(${scale}, ${scale})`,
          }}
        >
          {children}
        </div>
      </div>
      <div
        ref={controlsRef}
        className="absolute right-0 top-0 gap-x-2 flex flex-row pb-4"
      >
        <div className="rounded-sm bg-white dark:bg-[#1e1e1e] shadow-md border border-[#e5e5e5] dark:border-[#333333] hover:bg-[#f5f5f5] dark:hover:bg-[#333333]">
          <button
            type="button"
            className="cursor-pointer p-2"
            onClick={() => onCenter()}
          >
            <FluentFullScreenMaximize16Regular />
          </button>
        </div>
        <div className="flex flex-row divide-x divide-[#e5e5e5] dark:divide-[#333333] rounded-sm bg-white dark:bg-[#1e1e1e] shadow-md border border-[#e5e5e5] dark:border-[#333333]">
          <button
            type="button"
            disabled={Math.abs(scale - scaleToFit) < 0.001}
            className="p-2 cursor-pointer disabled:cursor-not-allowed disabled:bg-[#f5f5f5] dark:disabled:bg-[#333333] hover:bg-[#f5f5f5] dark:hover:bg-[#333333]"
            onClick={() => onZoom(-1)}
          >
            <FluentSubtract16Regular />
          </button>
          <button
            type="button"
            disabled={scale === MAX_SCALE}
            className="p-2 cursor-pointer disabled:cursor-not-allowed disabled:bg-[#f5f5f5] dark:disabled:bg-[#333333] hover:bg-[#f5f5f5] dark:hover:bg-[#333333]"
            onClick={() => onZoom(1)}
          >
            <FluentAdd16Regular />
          </button>
        </div>
      </div>
    </div>
  )
}
