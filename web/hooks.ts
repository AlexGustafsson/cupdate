import { Dispatch, SetStateAction, useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

export function useFilter(): [string[], Dispatch<SetStateAction<string[]>>] {
  const [searchParams, setSearchParams] = useSearchParams()

  const [filter, setFilter] = useState(
    (searchParams.get('filter') || '').split(',').filter((x) => x.length > 0)
  )

  useEffect(() => {
    setSearchParams((current) => {
      if (filter.length === 0) {
        current.delete('filter')
      } else {
        current.set('filter', filter.join(','))
      }
      return current
    })
  }, [filter])

  return [filter, setFilter]
}

export function useSort(): [
  string | undefined,
  Dispatch<SetStateAction<string | undefined>>,
  'asc' | 'desc',
  Dispatch<SetStateAction<'asc' | 'desc'>>,
] {
  const [searchParams, setSearchParams] = useSearchParams()

  // TODO: Will treat an empty string differently from a missing value
  const [property, setProperty] = useState(
    searchParams.get('sort') || undefined
  )
  const [order, setOrder] = useState<'asc' | 'desc'>(
    searchParams.has('asc') ? 'asc' : searchParams.has('desc') ? 'desc' : 'asc'
  )

  useEffect(() => {
    setSearchParams((current) => {
      if (!property) {
        current.delete('sort')
      } else {
        current.set('sort', property)
      }
      return current
    })
  }, [property])

  useEffect(() => {
    setSearchParams((current) => {
      current.delete('asc')
      current.delete('desc')
      if (order === 'asc') {
        current.set('asc', '')
      } else if (order === 'desc') {
        current.set('desc', '')
      }
      return current
    })
  }, [order])

  return [property, setProperty, order, setOrder]
}
