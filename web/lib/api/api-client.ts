import type { Vulnerability } from '../osv/osv'
import type { GetImagesOptions, ApiClient as IApiClient } from './client'
import type {
  Graph,
  Image,
  ImageDescription,
  ImagePage,
  ImageProvenance,
  ImageReleaseNotes,
  ImageSBOM,
  ImageScorecard,
  WebPushSubscription,
  WorkflowRun,
} from './models'

export const DEFAULT_API_ENDPOINT = import.meta.env.VITE_API_ENDPOINT

export const DEFAULT_RSS_ENDPOINT = `${import.meta.env.VITE_API_ENDPOINT}/feed.rss`

export class ApiClient implements IApiClient {
  #endpoint: string

  constructor(endpoint: string) {
    this.#endpoint = endpoint
  }

  async getTags(): Promise<string[]> {
    const res = await fetch(`${this.#endpoint}/tags`)

    if (res.status !== 200) {
      throw new Error(`unexpected status code ${res.status}`)
    }

    return res.json()
  }

  async getImages(options?: GetImagesOptions): Promise<ImagePage> {
    const searchParams = new URLSearchParams()

    if (options?.tags !== undefined) {
      for (const tag of options.tags) {
        searchParams.append('tag', tag)
      }
    }

    if (options?.tagop !== undefined) {
      searchParams.append('tagop', options.tagop)
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

    const res = await fetch(
      `${this.#endpoint}/images?${searchParams.toString()}`,
      {
        headers: {
          accept: 'application/json',
        },
      }
    )

    if (res.status !== 200) {
      throw new Error(`unexpected status code ${res.status}`)
    }

    return res.json()
  }

  async #getResource<T>(path: string, reference: string): Promise<T | null> {
    const query = new URLSearchParams({ reference })

    const res = await fetch(
      `${this.#endpoint}${path}${query === undefined ? '' : `?${query.toString()}`}`,
      {
        headers: {
          accept: 'application/json',
        },
      }
    )

    if (res.status === 404) {
      return null
    } else if (res.status !== 200) {
      throw new Error(`unexpected status code ${res.status}`)
    }

    return res.json()
  }

  getImage(reference: string): Promise<Image | null> {
    return this.#getResource('/image', reference)
  }

  getImageDescription(reference: string): Promise<ImageDescription | null> {
    return this.#getResource('/image/description', reference)
  }

  getImageReleaseNotes(reference: string): Promise<ImageReleaseNotes | null> {
    return this.#getResource('/image/release-notes', reference)
  }

  getImageGraph(reference: string): Promise<Graph | null> {
    return this.#getResource('/image/graph', reference)
  }

  getImageScorecard(reference: string): Promise<ImageScorecard | null> {
    return this.#getResource('/image/scorecard', reference)
  }

  getImageProvenance(reference: string): Promise<ImageProvenance | null> {
    return this.#getResource('/image/provenance', reference)
  }

  getImageSBOM(reference: string): Promise<ImageSBOM | null> {
    return this.#getResource('/image/sbom', reference)
  }

  getImageVulnerabilities(reference: string): Promise<Vulnerability[] | null> {
    return this.#getResource<{ vulnerabilities: Vulnerability[] }>(
      '/image/vulnerabilities',
      reference
    ).then((x) => x?.vulnerabilities || null)
  }

  getLatestImageWorkflow(reference: string): Promise<WorkflowRun | null> {
    return this.#getResource('/image/workflows/latest', reference)
  }

  async getWebPushServerKey(): Promise<string> {
    const res = await fetch(`${this.#endpoint}/webpush/key`, {
      method: 'GET',
      headers: {
        accept: 'application/json',
      },
    })

    if (res.status !== 200) {
      throw new Error(`unexpected status code ${res.status}`)
    }

    const body = await res.json()
    return body.key
  }

  async createWebPushSubscription(
    subscription: WebPushSubscription
  ): Promise<void> {
    const res = await fetch(`${this.#endpoint}/webpush/key`, {
      method: 'POST',
      headers: {
        accept: 'application/json',
        'content-type': 'application/json',
      },
      body: JSON.stringify(subscription),
    })

    if (res.status === 409) {
      throw new Error('subscription already exists')
    } else if (res.status !== 201) {
      throw new Error(`unexpected status code ${res.status}`)
    }
  }

  async deleteWebPushSubscription(digest: string): Promise<void> {
    const res = await fetch(
      `${this.#endpoint}/webpush/subscription/${encodeURIComponent(digest)}`,
      {
        method: 'DELETE',
      }
    )

    if (res.status === 404) {
      throw new Error('subscription does not exist')
    } else if (res.status !== 201) {
      throw new Error(`unexpected status code ${res.status}`)
    }
  }

  async checkWebPushSubscription(digest: string): Promise<boolean> {
    const res = await fetch(
      `${this.#endpoint}/webpush/subscription/${encodeURIComponent(digest)}`,
      {
        method: 'HEAD',
      }
    )

    switch (res.status) {
      case 200:
        return true
      case 404:
        return false
      default:
        throw new Error(`unexpected status code ${res.status}`)
    }
  }

  async scheduleImageScan(reference: string): Promise<void> {
    const query = new URLSearchParams({ reference })

    const res = await fetch(
      `${this.#endpoint}/image/scans?${query.toString()}`,
      {
        method: 'POST',
      }
    )

    if (res.status !== 202) {
      throw new Error(`unexpected status - ${res.status}`)
    }
  }

  async dump(): Promise<void> {
    const tags = await this.getTags()

    const pages = [await this.getImages()]

    // Arbitrarily limit to 2 pages
    for (let i = 1; i <= 2; i++) {
      const totalPages = Math.ceil(
        pages[0].pagination.total / pages[0].pagination.size
      )

      if (i > totalPages) {
        break
      }

      pages.push(await this.getImages({ page: i }))
    }

    const resources: Record<string, unknown> = {}

    // The loop is a naive attempt to lower memory usage a bit
    for (let i = 0; i < pages.length; i++) {
      for (let j = 0; j < pages[i].images.length; j++) {
        const reference = pages[i].images[j].reference
        resources[`${reference}`] = Object.fromEntries(
          await Promise.all(
            [
              this.getImage,
              this.getImageDescription,
              this.getImageReleaseNotes,
              this.getImageGraph,
              this.getImageScorecard,
              this.getImageProvenance,
              this.getImageSBOM,
              this.getImageVulnerabilities,
              this.getLatestImageWorkflow,
            ].map((x) =>
              x
                .bind(this)(reference)
                .then((y) => [x.name, y])
            )
          )
        )
      }
    }

    const content = JSON.stringify({ tags, pages, resources })

    const blob = new Blob([content], { type: 'application/json' })

    const a = document.createElement('a')
    a.download = 'cupdate-dump.json'
    a.href = URL.createObjectURL(blob)
    a.click()
  }
}
