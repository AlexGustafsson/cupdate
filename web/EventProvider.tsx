import {
  type JSX,
  type PropsWithChildren,
  createContext,
  use,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react'

export type ImageEvent = { type: 'imageUpdated'; reference: string }

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

export function useEvents(callback: EventHandler) {
  const context = use(EventContext)

  useEffect(() => {
    context.addCallback(callback)
    return () => context.removeCallback(callback)
  }, [context.addCallback, context.removeCallback, callback])
}
