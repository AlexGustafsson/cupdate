import type {
  Graph,
  Image,
  ImageDescription,
  ImagePage,
  ImageProvenance,
  ImageReleaseNotes,
  ImageSBOM,
  ImageScorecard,
  PaginationMetadata,
  WorkflowRun,
} from './models'

import {
  type JSX,
  type PropsWithChildren,
  createContext,
  use,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'
import { type Tag, Tags, tagByName } from '../../tags'
import type { Vulnerability } from '../osv/osv'
import type { ApiClient, GetImagesOptions } from './client'

export type Result<T> =
  | { status: 'idle' }
  | { status: 'resolved'; value: T }
  | { status: 'rejected'; error: unknown }

const ApiContext = createContext<ApiClient>({} as ApiClient)

export function ApiProvider({
  children,
  client,
}: PropsWithChildren<{ client: ApiClient }>): JSX.Element {
  return <ApiContext value={client}>{children}</ApiContext>
}

export function useApiClient(): ApiClient {
  const client = use(ApiContext)
  return client
}

export function useTags(): [Result<Tag[]>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<Tag[]>>({ status: 'idle' })

  const update = useCallback(() => {
    client
      .getTags()
      .then((value: string[]) =>
        setResult({
          status: 'resolved',
          value: Array.from(
            new Set([...value, ...Tags.map((x) => x.name)])
          ).map((x) => {
            const tag = tagByName(x) || {}
            return { ...tag, name: x }
          }),
        })
      )
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client])

  useEffect(() => update(), [update])

  return [result, update]
}

export function useImages(
  options?: GetImagesOptions
): [Result<ImagePage>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<ImagePage>>({ status: 'idle' })

  const update = useCallback(() => {
    client
      .getImages(options)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, options])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImage(
  reference: string
): [Result<Image | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<Image | null>>({ status: 'idle' })

  const update = useCallback(() => {
    client
      .getImage(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImageDescription(
  reference: string
): [Result<ImageDescription | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<ImageDescription | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getImageDescription(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImageReleaseNotes(
  reference: string
): [Result<ImageReleaseNotes | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<ImageReleaseNotes | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getImageReleaseNotes(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImageGraph(
  reference: string
): [Result<Graph | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<Graph | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getImageGraph(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImageScorecard(
  reference: string
): [Result<ImageScorecard | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<ImageScorecard | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getImageScorecard(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImageProvenance(
  reference: string
): [Result<ImageProvenance | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<ImageProvenance | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getImageProvenance(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImageSBOM(
  reference: string
): [Result<ImageSBOM | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<ImageSBOM | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getImageSBOM(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useImageVulnerabilities(
  reference: string
): [Result<Vulnerability[] | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<Vulnerability[] | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getImageVulnerabilities(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useLatestWorkflowRun(
  reference: string
): [Result<WorkflowRun | null>, () => void] {
  const client = useApiClient()
  const [result, setResult] = useState<Result<WorkflowRun | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    client
      .getLatestImageWorkflow(reference)
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [client, reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function useScheduleScan(): (reference: string) => Promise<void> {
  const client = useApiClient()

  return client.scheduleImageScan.bind(client)
}

export function usePagination<T extends { pagination: PaginationMetadata }>(
  page: T | undefined,
  searchParams: URLSearchParams
): (
  | {
      label: string
      current: boolean
      index: undefined
    }
  | { label: string; current: boolean; index: number; href: string }
)[] {
  const pages = useMemo(() => {
    if (!page) {
      return []
    }

    const hrefForPage = (page: number) => {
      const params = new URLSearchParams(searchParams)
      params.set('page', page.toString())
      return `/?${params.toString()}`
    }

    // Page index in the API starts at 1.
    const pageIndex = Math.max(page.pagination.page - 1, 0)

    const totalPages = Math.ceil(page.pagination.total / page.pagination.size)

    // Try to keep 9 pages displayed at all times, with 4 pages allocated
    // for previous pages and 5 for next pages
    let rangeStart = Math.max(0, pageIndex - 4)
    const rangeEnd = Math.min(
      totalPages,
      pageIndex + 9 - (pageIndex - rangeStart)
    )
    rangeStart = Math.max(0, pageIndex - 9 - (pageIndex - rangeEnd))
    const range = rangeEnd - rangeStart

    const pages: (
      | {
          label: string
          current: boolean
          index: undefined
        }
      | { label: string; current: boolean; index: number; href: string }
    )[] = new Array(range).fill('').map((_, i) => ({
      label: (rangeStart + i + 1).toString(),
      current: rangeStart + i === pageIndex,
      index: rangeStart + i,
      href: hrefForPage(rangeStart + i + 1),
    }))

    if (pages[8]?.index && pages[8].index < totalPages - 1) {
      pages[7] = { index: undefined, label: '...', current: false }
      pages[8] = {
        label: totalPages.toString(),
        current: false,
        index: totalPages - 1,
        href: hrefForPage(totalPages),
      }
    }

    if (pages[0]?.index && pages[0].index > 0) {
      pages[1] = { index: undefined, label: '...', current: false }
      pages[0] = {
        label: '1',
        current: false,
        index: 0,
        href: hrefForPage(0),
      }
    }

    return pages
  }, [page, searchParams])

  return pages
}

declare global {
  interface Window {
    pushManager?: PushManager
  }
}

export function usePushNotifications(): [
  boolean,
  PushSubscription | null,
  () => void,
  () => void,
] {
  const supported = window.pushManager !== undefined
  const [subscription, setSubscription] = useState<PushSubscription | null>(
    null
  )

  useEffect(() => {
    if (!window.pushManager) {
      return
    }

    window.pushManager
      .getSubscription()
      .then((subscription) => {
        console.log(subscription)
        setSubscription(subscription)
      })
      .catch((error) => {
        console.error('Failed to get current subscription', error)
      })
  }, [])

  const subscribe = useCallback(() => {
    if (!window.pushManager) {
      return Promise.reject(new Error('not supported'))
    }

    const options: PushSubscriptionOptionsInit = {
      applicationServerKey:
        'BOoKX5pU-ZMFDknU6ny_iRLuIiqXXykMTumqFaPEIAGATCi3unaV6qtUUcuCL5teYK00__BDd3CBuNiQi9yrn-c',
      userVisibleOnly: true,
    }

    window.pushManager
      .subscribe(options)
      .then((subscription) => {
        setSubscription(subscription)
        console.log(subscription)
      })
      .catch((error) => {
        console.error('Failed to subscribe to notifications', error)
      })
  }, [])

  const unsubscribe = useCallback(() => {
    if (!window.pushManager) {
      return Promise.reject(new Error('not supported'))
    }

    subscription
      ?.unsubscribe()
      .then(() => {
        window.pushManager!.getSubscription().then(setSubscription)
      })
      .catch((error) => {
        console.error('Failed to unsubscribe to notifications', error)
      })
  }, [subscription])

  return [supported, subscription, subscribe, unsubscribe]
}
