import { useCallback, useEffect, useMemo, useState } from 'react'

import { type Tag, TagsByName } from './tags'

export interface ImagePage {
  images: Image[]
  summary: ImagePageSummary
  pagination: PaginationMetadata
}

export interface ImagePageSummary {
  images: number
  outdated: number
  vulnerable: number
  processing: number
  failed: number
}

export interface PaginationMetadata {
  total: number
  /** Page index. Starts at 1. */
  page: number
  size: number
  next?: string
  previous?: string
}

export interface Image {
  reference: string
  created?: string
  latestReference?: string
  latestCreated?: string
  description?: string
  tags: string[]
  links: ImageLink[]
  vulnerabilities: ImageVulnerability[]
  image?: string
  lastModified: string
}

export interface ImageDescription {
  html?: string
  markdown?: string
}

export interface ImageReleaseNotes {
  title: string
  html?: string
  released?: string
}

export interface ImageLink {
  type: string
  url: string
}

export interface ImageVulnerability {
  severity: string
  authority: string
  links: string[]
  description?: string
}

export interface Graph {
  edges: Record<string, Record<string, boolean>>
  nodes: Record<string, GraphNode>
}

export interface GraphNode {
  domain: string
  type: string
  name: string
}

export interface WorkflowRun {
  jobs: JobRun[]
  traceId?: string
}

export type JobRun = {
  jobId?: string
  jobName?: string
  dependsOn: string[]
  steps: StepRun[]
} & (
  | {
      result: 'succeeded' | 'failed'
      started: string
      duration: number
    }
  | { result: 'skipped' }
)

export interface StepRun {
  stepName?: string
  result: 'succeeded' | 'skipped' | 'failed'
  error?: string
  duration?: number
}

export type Result<T> =
  | { status: 'idle' }
  | { status: 'resolved'; value: T }
  | { status: 'rejected'; error: unknown }

export function useTags(): [Result<Tag[]>, () => void] {
  const [result, setResult] = useState<Result<Tag[]>>({ status: 'idle' })

  const update = useCallback(() => {
    fetch(`${import.meta.env.VITE_API_ENDPOINT}/tags`)
      .then((res) => {
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value: string[]) =>
        setResult({
          status: 'resolved',
          value: value.map((x) => TagsByName[x] || { name: x }),
        })
      )
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [])

  useEffect(() => update(), [update])

  return [result, update]
}

interface UseImagesProps {
  tags?: string[]
  // Technically only "reference" | undefined, but let's be lax for now as we
  // otherwise would have to parse whatever query parameter we got and handle
  // errors
  sort?: string
  order?: 'asc' | 'desc'
  page?: number
  limit?: number
  query?: string
}

export function useImages(
  options?: UseImagesProps
): [Result<ImagePage>, () => void] {
  const [result, setResult] = useState<Result<ImagePage>>({ status: 'idle' })

  const update = useCallback(() => {
    const query = new URLSearchParams()
    if (options?.tags !== undefined) {
      for (const tag of options.tags) {
        query.append('tag', tag)
      }
    }
    if (options?.sort !== undefined) {
      query.set('sort', options.sort)
    }
    if (options?.order !== undefined) {
      query.set('order', options.order)
    }
    if (options?.page !== undefined) {
      // Page index starts at 1
      query.set('page', (options.page + 1).toString())
    }
    if (options?.limit !== undefined) {
      query.set('limit', options.limit.toString())
    }
    if (options?.query !== undefined) {
      query.set('query', options.query)
    }

    fetch(`${import.meta.env.VITE_API_ENDPOINT}/images?${query.toString()}`)
      .then((res) => {
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [
    options?.tags,
    options?.sort,
    options?.order,
    options?.page,
    options?.limit,
    options?.query,
  ])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

// TODO: Add query parameters
export function useImage(
  reference: string
): [Result<Image | null>, () => void] {
  const [result, setResult] = useState<Result<Image | null>>({ status: 'idle' })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(`${import.meta.env.VITE_API_ENDPOINT}/image?${query.toString()}`)
      .then((res) => {
        if (res.status === 404) {
          return null
        } else if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

// TODO: Add query parameters
export function useImageDescription(
  reference: string
): [Result<ImageDescription | null>, () => void] {
  const [result, setResult] = useState<Result<ImageDescription | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/image/description?${query.toString()}`
    )
      .then((res) => {
        if (res.status === 404) {
          return null
        }
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

// TODO: Add query parameters
export function useImageReleaseNotes(
  reference: string
): [Result<ImageReleaseNotes | null>, () => void] {
  const [result, setResult] = useState<Result<ImageReleaseNotes | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/image/release-notes?${query.toString()}`
    )
      .then((res) => {
        if (res.status === 404) {
          return null
        }
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

// TODO: Add query parameters
export function useImageGraph(
  reference: string
): [Result<Graph | null>, () => void] {
  const [result, setResult] = useState<Result<Graph | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/image/graph?${query.toString()}`
    )
      .then((res) => {
        if (res.status === 404) {
          return null
        }
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

// TODO: Add query parameters
export function useLatestWorkflowRun(
  reference: string
): [Result<WorkflowRun | null>, () => void] {
  const [result, setResult] = useState<Result<WorkflowRun | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/image/workflows/latest?${query.toString()}`
    )
      .then((res) => {
        if (res.status === 404) {
          return null
        }
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [reference])

  useEffect(() => {
    update()
  }, [update])

  return [result, update]
}

export function usePagination<T extends { pagination: PaginationMetadata }>(
  page: T | undefined
): { index: number | undefined; label: string; current: boolean }[] {
  const pages = useMemo(() => {
    if (!page) {
      return []
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

    const pages: {
      index: number | undefined
      label: string
      current: boolean
    }[] = new Array(range).fill('').map((_, i) => ({
      index: rangeStart + i,
      label: (rangeStart + i + 1).toString(),
      current: rangeStart + i === pageIndex,
    }))

    if (pages[8]?.index && pages[8].index < totalPages - 1) {
      pages[7] = { index: undefined, label: '...', current: false }
      pages[8] = {
        index: totalPages - 1,
        label: totalPages.toString(),
        current: false,
      }
    }

    if (pages[0]?.index && pages[0].index > 0) {
      pages[1] = { index: undefined, label: '...', current: false }
      pages[0] = {
        index: 0,
        label: '1',
        current: false,
      }
    }

    return pages
  }, [page])

  return pages
}

export async function scheduleScan(reference: string): Promise<void> {
  const query = new URLSearchParams({ reference })
  const res = await fetch(
    `${import.meta.env.VITE_API_ENDPOINT}/image/scans?${query.toString()}`,
    {
      method: 'POST',
    }
  )

  if (res.status !== 202) {
    throw new Error(`unexpected status - ${res.status}`)
  }
}

export const RSSFeedEndpoint = `${import.meta.env.VITE_API_ENDPOINT}/feed.rss`
