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
  image?: string
}

export interface ImageDescription {
  html?: string
}

export interface ImageReleaseNotes {
  title: string
  html?: string
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

export function useImages(): Result<ImagePage> {
  const [result, setResult] = useState<Result<ImagePage>>({ status: 'idle' })

  useEffect(() => {
    fetch(`${import.meta.env['VITE_API_ENDPOINT']}/images`)
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

// TODO: Add query parameters
export function useImage(): Result<Image> {
  const [result, setResult] = useState<Result<Image>>({ status: 'idle' })

  useEffect(() => {
    fetch(`${import.meta.env['VITE_API_ENDPOINT']}/image`)
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

// TODO: Add query parameters
export function useImageDescription(): Result<ImageDescription | undefined> {
  const [result, setResult] = useState<Result<ImageDescription | undefined>>({
    status: 'idle',
  })

  useEffect(() => {
    fetch(`${import.meta.env['VITE_API_ENDPOINT']}/image/description`)
      .then((res) => {
        if (res.status === 404) {
          setResult({ status: 'resolved', value: undefined })
        } else if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [])

  return result
}

// TODO: Add query parameters
export function useImageReleaseNotes(): Result<ImageReleaseNotes | undefined> {
  const [result, setResult] = useState<Result<ImageReleaseNotes | undefined>>({
    status: 'idle',
  })

  useEffect(() => {
    fetch(`${import.meta.env['VITE_API_ENDPOINT']}/image/release-notes`)
      .then((res) => {
        if (res.status === 404) {
          setResult({ status: 'resolved', value: undefined })
        } else if (res.status !== 200) {
          throw new Error(`unexpected status code ${res.status}`)
        }

        return res.json()
      })
      .then((value) => setResult({ status: 'resolved', value }))
      .catch((error) => setResult({ status: 'rejected', error }))
  }, [])

  return result
}
