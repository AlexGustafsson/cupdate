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
import { useApiClient } from './lib/api/ApiProvider'
import type { WebPushSubscription } from './lib/api/models'
import { webPushSubscriptionDigest } from './lib/api/util'

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

// Polyfill until proper typing support exists
declare global {
  interface Window {
    pushManager?: PushManager
  }
}

type Result<T, E = Error> =
  | { status: 'loading' }
  | { status: 'resolved'; value: T }
  | { status: 'rejected'; error: E }

export function useWebPushSubscription(): [
  boolean,
  Result<PushSubscription | null>,
  Result<boolean>,
  () => void,
  () => void,
] {
  const apiClient = useApiClient()

  const supported = window.pushManager !== undefined

  const [subscription, setSubscription] = useState<
    Result<PushSubscription | null>
  >({
    status: 'loading',
  })

  const [synced, setSynced] = useState<Result<boolean>>({ status: 'loading' })

  // Get the user agent's current subscription on load
  useEffect(() => {
    if (!window.pushManager) {
      return
    }

    window.pushManager
      .getSubscription()
      .then((subscription) => {
        setSubscription({
          status: 'resolved',
          value: subscription,
        })
      })
      .catch((error) => {
        setSubscription({ status: 'rejected', error })
        console.error('failed to set subscription', error)
      })
  }, [])

  // Keep synced state up-to-date
  useEffect(() => {
    if (subscription.status === 'loading') {
      setSynced({ status: 'loading' })
      return
    } else if (subscription.status === 'rejected') {
      setSynced({
        status: 'rejected',
        error: new Error('cannot check state when subscription status failed'),
      })
      return
    }

    if (subscription.value === null) {
      // Assume synced as we have no way of knowing if the server has the
      // subscription or not. In practice, the subscription will likely start
      // bouncing when the user agent no longer has it - meaning the server
      // should drop the subscription soon enough
      setSynced({ status: 'resolved', value: true })
      return
    }

    webPushSubscriptionDigest(subscription.value.toJSON())
      .then((digest) => {
        apiClient.checkWebPushSubscription(digest).then((ok) => {
          setSynced({ status: 'resolved', value: ok })
        })
      })
      .catch((error) => {
        setSynced({ status: 'rejected', error })
        console.error('failed to set sync status', error)
      })
  }, [subscription, apiClient])

  const subscribe = useCallback(() => {
    apiClient
      .getWebPushServerKey()
      .then((applicationServerKey) =>
        window.pushManager
          ?.subscribe({
            applicationServerKey,
            userVisibleOnly: true,
          })
          .then((subscription) =>
            apiClient
              .createWebPushSubscription(
                subscription.toJSON() as WebPushSubscription
              )
              .then(() => subscription)
          )
          .then((subscription) => {
            setSubscription({ status: 'resolved', value: subscription })
          })
      )
      .catch((error) => {
        console.error('failed to subscribe', error)
      })
  }, [apiClient])

  const unsubscribe = useCallback(() => {
    window.pushManager
      ?.getSubscription()
      .then((subscription) => {
        if (!subscription) {
          return
        }

        return webPushSubscriptionDigest(subscription.toJSON())
          .then((digest) =>
            subscription
              .unsubscribe()
              .then((ok) => [ok, digest] as [boolean, string])
          )
          .then(([ok, digest]) => {
            if (!ok) {
              throw new Error('Failed to unsubscribe')
            }

            return apiClient.deleteWebPushSubscription(digest)
          })
      })
      .catch((error) => {
        console.error('failed to unsubscribe', error)
      })
  }, [apiClient])

  return [supported, subscription, synced, subscribe, unsubscribe]
}
