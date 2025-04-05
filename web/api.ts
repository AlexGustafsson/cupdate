import { useCallback, useEffect, useMemo, useState } from 'react'

import type { Vulnerability } from './lib/osv/osv'
import { type Tag, tagByName } from './tags'

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
  vulnerabilities: number
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

export interface Graph {
  edges: Record<string, Record<string, boolean>>
  nodes: Record<string, GraphNode>
}

export interface GraphNode {
  domain: string
  type: string
  name: string
  labels?: Record<string, string>
  internalLabels?: Record<string, string>
}

export interface ImageScorecard {
  reportUrl: string
  score: number
  risk: 'critical' | 'high' | 'medium' | 'low'
  generatedAt: string
}

export interface ImageProvenance {
  buildInfo: ProvenanceBuildInfo[]
}

export interface ProvenanceBuildInfo {
  imageDigest: string
  architecture?: string
  architectureVariant?: string
  operatingSystem?: string
  source?: string
  sourceRevision?: string
  buildStartedOn?: string
  buildFinishedOn?: string
  dockerfile?: string
}

export interface ImageSBOM {
  sbom: SBOM[]
}

export interface SBOM {
  imageDigest: string
  type: 'spdx'
  sbom: string
  architecture?: string
  architectureVariant?: string
  operatingSystem?: string
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
          value: value.map((x) => {
            const tag = tagByName(x) || {}
            return { ...tag, name: x }
          }),
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
): [Result<ImagePage>, URLSearchParams, () => void] {
  const [result, setResult] = useState<Result<ImagePage>>({ status: 'idle' })

  const searchParams = useMemo(() => {
    const searchParams = new URLSearchParams()
    if (options?.tags !== undefined) {
      for (const tag of options.tags) {
        searchParams.append('tag', tag)
      }
    }

    if (options?.sort !== undefined) {
      searchParams.set('sort', options.sort)
    }

    if (options?.order !== undefined) {
      searchParams.set('order', options.order)
    }

    if (options?.page !== undefined) {
      // Page index starts at 1
      searchParams.set('page', (options.page + 1).toString())
    }

    if (options?.limit !== undefined) {
      searchParams.set('limit', options.limit.toString())
    }

    if (options?.query !== undefined) {
      searchParams.set('query', options.query)
    }

    return searchParams
  }, [
    options?.tags,
    options?.sort,
    options?.order,
    options?.page,
    options?.limit,
    options?.query,
  ])

  const update = useCallback(() => {
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/images?${searchParams.toString()}`
    )
      .then((res) => {
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [searchParams])

  useEffect(() => {
    update()
  }, [update])

  return [result, searchParams, update]
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
export function useImageScorecard(
  reference: string
): [Result<ImageScorecard | null>, () => void] {
  const [result, setResult] = useState<Result<ImageScorecard | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/image/scorecard?${query.toString()}`
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
export function useImageProvenance(
  reference: string
): [Result<ImageProvenance | null>, () => void] {
  const [result, setResult] = useState<Result<ImageProvenance | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/image/provenance?${query.toString()}`
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
export function useImageSBOM(
  reference: string
): [Result<ImageSBOM | null>, () => void] {
  const [result, setResult] = useState<Result<ImageSBOM | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(`${import.meta.env.VITE_API_ENDPOINT}/image/sbom?${query.toString()}`)
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
export function useImageVulnerabilities(
  reference: string
): [Result<Vulnerability[] | null>, () => void] {
  const [result, setResult] = useState<Result<Vulnerability[] | null>>({
    status: 'idle',
  })

  const update = useCallback(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env.VITE_API_ENDPOINT}/image/vulnerabilities?${query.toString()}`
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
      .then((value) => {
        setResult({
          status: 'resolved',
          value: value.vulnerabilities,
        })
      })
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
