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
