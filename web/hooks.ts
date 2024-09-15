import { Dispatch, SetStateAction, useCallback, useMemo } from 'react'
import { useSearchParams } from 'react-router-dom'

export function useFilter(): [string[], Dispatch<SetStateAction<string[]>>] {
  const [searchParams, setSearchParams] = useSearchParams()

  const filter = useMemo(() => {
    return (searchParams.get('tags') || '')
      .split(',')
      .filter((x) => x.length > 0)
  }, [searchParams])

  const setFilter = useCallback(
    (s: string[] | ((current: string[]) => string[])) => {
      setSearchParams((current) => {
        if (typeof s === 'function') {
          s = s(
            (searchParams.get('tags') || '')
              .split(',')
              .filter((x) => x.length > 0)
          )
        }
        if (!s) {
          current.delete('tags')
        } else {
          current.set('tags', s.join(','))
        }
        return current
      })
    },
    [setSearchParams]
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
    searchParams.get('order') == 'desc'
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
