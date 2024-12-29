import {
  type Dispatch,
  type SetStateAction,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react'
import { useSearchParams } from 'react-router-dom'

export function useFilter(): [string[], Dispatch<SetStateAction<string[]>>] {
  const [searchParams, setSearchParams] = useSearchParams()

  const filter = useMemo(() => {
    return searchParams.getAll('tag')
  }, [searchParams])

  const setFilter = useCallback(
    (s: string[] | ((current: string[]) => string[])) => {
      setSearchParams((current) => {
        if (typeof s === 'function') {
          s = s(searchParams.getAll('tag'))
        }
        current.delete('tag')
        if (s) {
          for (const tag of s) {
            current.append('tag', tag)
          }
        }
        return current
      })
    },
    [searchParams, setSearchParams]
  )

  return [filter, setFilter]
}

export function useSort(): [
  string | undefined,
  Dispatch<SetStateAction<string | undefined>>,
  'asc' | 'desc' | undefined,
  Dispatch<SetStateAction<'asc' | 'desc' | undefined>>,
] {
  const [searchParams, setSearchParams] = useSearchParams()

  const property = searchParams.get('sort') || undefined
  const order =
    searchParams.get('order') === 'desc'
      ? 'desc'
      : searchParams.get('order') === 'asc'
        ? 'asc'
        : undefined

  const setProperty = useCallback(
    (s?: string | ((current: string | undefined) => string | undefined)) => {
      setSearchParams((current) => {
        if (typeof s === 'function') {
          s = s(current.get('sort') || undefined)
        }
        if (!s) {
          current.delete('sort')
        } else {
          current.set('sort', s)
        }
        return current
      })
    },
    [setSearchParams]
  )

  const setOrder = useCallback(
    (
      s?:
        | 'asc'
        | 'desc'
        | ((current: 'asc' | 'desc' | undefined) => 'asc' | 'desc' | undefined)
    ) => {
      setSearchParams((current) => {
        if (typeof s === 'function') {
          s = s(
            current.get('order') === 'asc'
              ? 'asc'
              : current.get('order') === 'desc'
                ? 'desc'
                : undefined
          )
        }
        if (!s) {
          current.delete('order')
        } else {
          current.set('order', s)
        }
        return current
      })
    },
    [setSearchParams]
  )

  return [property, setProperty, order, setOrder]
}

export function useDebouncedEffect(
  effect: React.EffectCallback,
  deps?: React.DependencyList
) {
  const timeoutRef = useRef<ReturnType<typeof setTimeout>>(null)

  // biome-ignore lint/correctness/useExhaustiveDependencies: effect should not be changed
  useEffect(() => {
    if (timeoutRef.current !== null) {
      clearTimeout(timeoutRef.current)
      timeoutRef.current = null
    }

    timeoutRef.current = setTimeout(() => {
      effect()
      timeoutRef.current = null
    }, 200)
  }, deps)
}

export function useQuery(): [
  string | undefined,
  Dispatch<SetStateAction<string | undefined>>,
] {
  const [searchParams, setSearchParams] = useSearchParams()

  const query = useMemo(() => {
    return searchParams.get('query') || undefined
  }, [searchParams])

  const setQuery = useCallback(
    (s?: string | ((current: string | undefined) => string | undefined)) => {
      setSearchParams((current) => {
        if (typeof s === 'function') {
          s = s(current.get('query') || undefined)
        }
        if (!s) {
          current.delete('query')
        } else {
          current.set('query', s)
        }
        return current
      })
    },
    [setSearchParams]
  )

  return [query, setQuery]
}
