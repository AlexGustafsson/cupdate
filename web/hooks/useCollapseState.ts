import { useCallback, useEffect, useState } from 'react'

export function useCollapseState(): [Set<string>, (id: string) => void] {
  const [collapsed, setCollapsed] = useState<Set<string>>(() => {
    // TODO: Remove in future version. Done to clean up local storage
    for (const key of Object.keys(localStorage)) {
      if (key.startsWith('cupdate-card-state-')) {
        localStorage.removeItem(key)
      }
    }

    let collapsed = new Set<string>()
    const item = localStorage.getItem('cupdate-collapsed-cards')
    if (item) {
      try {
        const items = JSON.parse(item)
        if (Array.isArray(items)) {
          collapsed = new Set(items.filter((x) => typeof x === 'string'))
        }
      } catch {
        localStorage.removeItem('cupdate-collapsed-cards')
      }
    }

    return collapsed
  })

  useEffect(() => {
    localStorage.setItem(
      'cupdate-collapsed-cards',
      JSON.stringify(Array.from(collapsed))
    )
  }, [collapsed])

  const toggleCollapsed = useCallback((id: string) => {
    setCollapsed((current) =>
      current.has(id)
        ? new Set([...current]).difference(new Set([id]))
        : new Set([...current, id])
    )
  }, [])

  return [collapsed, toggleCollapsed]
}
