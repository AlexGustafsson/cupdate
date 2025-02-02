import { useEffect, useRef, useState } from 'react'

export function useMountedDuration(): number {
  const startedRef = useRef(performance.now())
  const [duration, setDuration] = useState(0)

  useEffect(() => {
    const callback = () => {
      setDuration(performance.now() - startedRef.current)
    }
    const handle = setInterval(callback, 1000)
    return () => clearInterval(handle)
  }, [])

  return duration
}
