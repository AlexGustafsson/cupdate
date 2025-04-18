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

export interface Filter {
  tags: string[]
  operator?: 'and' | 'or'
}

export function useFilter(): [Filter, Dispatch<SetStateAction<Filter>>] {
  const [searchParams, setSearchParams] = useSearchParams()

  const filter = useMemo(() => {
    const tagop = searchParams.get('tagop')
    return {
      tags: searchParams.getAll('tag'),
      operator: (tagop === 'and' || tagop === 'or' ? tagop : undefined) as
        | 'and'
        | 'or'
        | undefined,
    }
  }, [searchParams])

  const setFilter = useCallback(
    (s: Filter | ((current: Filter) => Filter)) => {
      setSearchParams((current) => {
        if (typeof s === 'function') {
          const tagop = searchParams.get('tagop')
          s = s({
            tags: searchParams.getAll('tag'),
            operator: tagop === 'and' || tagop === 'or' ? tagop : undefined,
          })
        }
        current.delete('tag')
        current.delete('tagop')
        if (s) {
          for (const tag of s.tags) {
            current.append('tag', tag)
          }
          if (s.operator) {
            current.set('tagop', s.operator)
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

export function useLayout(): [
  'list' | 'grid',
  Dispatch<SetStateAction<'list' | 'grid'>>,
] {
  const [layout, setLayout] = useState<'list' | 'grid'>(() => {
    const value = localStorage.getItem('layout')
    if (value === 'list' || value === 'grid') {
      return value
    }

    return 'list'
  })

  // Store layout in local storage
  useEffect(() => {
    localStorage.setItem('layout', layout)
  }, [layout])

  return [layout, setLayout]
}
