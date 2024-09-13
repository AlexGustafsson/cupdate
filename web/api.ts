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
  images?: number
  outdated?: number
  pods?: number
}

export interface PaginationMetadata {
  total: number
  page: number
  size: number
  next?: string
  previous?: string
}

export interface Image {
  name: string
  currentVersion: string
  latestVersion: string
  tags: string[]
  links: ImageLink[]
  image?: string
}

export interface ImageDescription {
  html?: string
}

export interface ImageReleaseNotes {
  title: string
  html?: string
}

export interface ImageLink {
  type: string
  url: string
}

export interface Graph {
  root: GraphNode
}

export interface GraphNode {
  domain: string
  type: string
  name: string
  parents: GraphNode[]
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
  sort?: string
  asc?: boolean
  desc?: boolean
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
    if (options?.asc !== undefined) {
      query.set('asc', '')
    }
    if (options?.desc !== undefined) {
      query.set('desc', '')
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
  }, [options])

  return result
}

// TODO: Add query parameters
export function useImage(name: string, version: string): Result<Image | null> {
  const [result, setResult] = useState<Result<Image | null>>({ status: 'idle' })

  useEffect(() => {
    const query = new URLSearchParams({ name, version })
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
  }, [name, version])

  return result
}

// TODO: Add query parameters
export function useImageDescription(
  name: string,
  version: string
): Result<ImageDescription | null> {
  const [result, setResult] = useState<Result<ImageDescription | null>>({
    status: 'idle',
  })

  useEffect(() => {
    const query = new URLSearchParams({ name, version })
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
  }, [name, version])

  return result
}

// TODO: Add query parameters
export function useImageReleaseNotes(
  name: string,
  version: string
): Result<ImageReleaseNotes | null> {
  const [result, setResult] = useState<Result<ImageReleaseNotes | null>>({
    status: 'idle',
  })

  useEffect(() => {
    const query = new URLSearchParams({ name, version })
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
  }, [name, version])

  return result
}

// TODO: Add query parameters
export function useImageGraph(
  name: string,
  version: string
): Result<Graph | null> {
  const [result, setResult] = useState<Result<Graph | null>>({
    status: 'idle',
  })

  useEffect(() => {
    const query = new URLSearchParams({ name, version })
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
  }, [name, version])

  return result
}
