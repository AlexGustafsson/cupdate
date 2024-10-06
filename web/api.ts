import { useEffect, useState } from 'react'

export interface Tag {
  name: string
  description: string
  color: string
}

export interface ImagePage {
  images: Image[]
  summary: ImagePageSummary
  pagination: PaginationMetadata
}

export interface ImagePageSummary {
  images: number
  outdated: number
  pods: number
}

export interface PaginationMetadata {
  total: number
  page: number
  size: number
  next?: string
  previous?: string
}

export interface Image {
  reference: string
  latestReference: string
  description?: string
  tags: string[]
  links: ImageLink[]
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
}

export type Result<T> =
  | { status: 'idle' }
  | { status: 'resolved'; value: T }
  | { status: 'rejected'; error: unknown }

export function useTags(): Result<Tag[]> {
  const [result, setResult] = useState<Result<Tag[]>>({ status: 'idle' })

  useEffect(() => {
    fetch(`${import.meta.env['VITE_API_ENDPOINT']}/tags`)
      .then((res) => {
        if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [])

  return result
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
}

export function useImages(options?: UseImagesProps): Result<ImagePage> {
  const [result, setResult] = useState<Result<ImagePage>>({ status: 'idle' })

  useEffect(() => {
    const query = new URLSearchParams()
    if (options?.tags !== undefined) {
      query.set('tags', options.tags.join(','))
    }
    if (options?.sort !== undefined) {
      query.set('sort', options.sort)
    }
    if (options?.order !== undefined) {
      query.set('order', options.order)
    }
    if (options?.page !== undefined) {
      query.set('page', options.page.toString())
    }
    if (options?.limit !== undefined) {
      query.set('limit', options.limit.toString())
    }

    fetch(`${import.meta.env['VITE_API_ENDPOINT']}/images?${query.toString()}`)
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
  ])

  return result
}

// TODO: Add query parameters
export function useImage(reference: string): Result<Image | null> {
  const [result, setResult] = useState<Result<Image | null>>({ status: 'idle' })

  useEffect(() => {
    const query = new URLSearchParams({ reference })
    fetch(`${import.meta.env['VITE_API_ENDPOINT']}/image?${query.toString()}`)
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

  return result
}

// TODO: Add query parameters
export function useImageDescription(
  reference: string
): Result<ImageDescription | null> {
  const [result, setResult] = useState<Result<ImageDescription | null>>({
    status: 'idle',
  })

  useEffect(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env['VITE_API_ENDPOINT']}/image/description?${query.toString()}`
    )
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

  return result
}

// TODO: Add query parameters
export function useImageReleaseNotes(
  reference: string
): Result<ImageReleaseNotes | null> {
  const [result, setResult] = useState<Result<ImageReleaseNotes | null>>({
    status: 'idle',
  })

  useEffect(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env['VITE_API_ENDPOINT']}/image/release-notes?${query.toString()}`
    )
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

  return result
}

// TODO: Add query parameters
export function useImageGraph(reference: string): Result<Graph | null> {
  const [result, setResult] = useState<Result<Graph | null>>({
    status: 'idle',
  })

  useEffect(() => {
    const query = new URLSearchParams({ reference })
    fetch(
      `${import.meta.env['VITE_API_ENDPOINT']}/image/graph?${query.toString()}`
    )
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

  return result
}
