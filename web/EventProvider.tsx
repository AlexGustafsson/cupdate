import {
  type DependencyList,
  type JSX,
  type PropsWithChildren,
  createContext,
  use,
  useCallback,
  useEffect,
  useRef,
} from 'react'

export type ImageEvent = {
  type: 'imageUpdated' | 'imageProcessed'
  reference: string
}

export type Event = ImageEvent

type EventHandler = (e: Event) => void

interface EventContextType {
  addCallback: (handler: EventHandler) => void
  removeCallback: (handler: EventHandler) => void
}

const EventContext = createContext<EventContextType>({
  addCallback: () => {},
  removeCallback: () => {},
})

export function EventProvider({ children }: PropsWithChildren): JSX.Element {
  const callbacksRef = useRef<EventHandler[]>([])

  const addCallback = useCallback((handler: EventHandler) => {
    callbacksRef.current = [
      ...callbacksRef.current.filter((x) => x !== handler),
      handler,
    ]
  }, [])

  const removeCallback = useCallback((handler: EventHandler) => {
    callbacksRef.current = [
      ...callbacksRef.current.filter((x) => x !== handler),
    ]
  }, [])

  const onMessage = useCallback((e: MessageEvent) => {
    try {
      const data = JSON.parse(e.data)
      if (
        data !== null &&
        typeof data === 'object' &&
        typeof data.type === 'string'
      ) {
        for (const callback of callbacksRef.current) {
          callback(data as Event)
        }
      }
    } catch {
      // Do nothing
    }
  }, [])

  useEffect(() => {
    const eventSource = new EventSource(
      `${import.meta.env.VITE_API_ENDPOINT}/events`
    )

    eventSource.addEventListener('message', onMessage)

    return () => eventSource.close()
  }, [onMessage])

  return (
    <EventContext value={{ addCallback, removeCallback }}>
      {children}
    </EventContext>
  )
}

export function useEvents(callback: EventHandler, deps: DependencyList) {
  const context = use(EventContext)

  const memoizedCallback = useCallback(callback, deps)

  useEffect(() => {
    context.addCallback(memoizedCallback)
    return () => context.removeCallback(memoizedCallback)
  }, [context.addCallback, context.removeCallback, memoizedCallback])
}
