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
] {
  const [searchParams, setSearchParams] = useSearchParams()

  // TODO: Will treat an empty string differently from a missing value
  const [sort, setSort] = useState(searchParams.get('sort') || undefined)

  useEffect(() => {
    setSearchParams((current) => {
      if (!sort) {
        current.delete('sort')
      } else {
        current.set('sort', sort)
      }
      return current
    })
  }, [sort])

  return [sort, setSort]
}
