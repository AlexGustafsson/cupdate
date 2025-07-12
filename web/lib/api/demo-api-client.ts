import type { Vulnerability } from '../osv/osv'

import type { ApiClient, GetImagesOptions } from './client'

import type {
  Graph,
  Image,
  ImageDescription,
  ImagePage,
  ImagePageSummary,
  ImageProvenance,
  ImageReleaseNotes,
  ImageSBOM,
  ImageScorecard,
  PaginationMetadata,
  WorkflowRun,
} from './models'

import DemoDataUrl from '../../public/assets/demo.json?url'

export interface Dump {
  tags: string[]
  pages: ImagePage[]
  resources: Record<string, Record<string, unknown>>
}

export class DemoApiClient implements ApiClient {
  #loaded: Promise<void>
  #dump: Dump

  constructor() {
    this.#loaded = fetch(DemoDataUrl)
      .then((x) => x.json())
      .then((x) => {
        this.#dump = x as unknown as Dump
      })

    this.#dump = {
      tags: [],
      pages: [],
      resources: {},
    }
  }

  async getTags(): Promise<string[]> {
    await this.#loaded
    return Promise.resolve(this.#dump.tags)
  }

  async getImages(options?: GetImagesOptions): Promise<ImagePage> {
    await this.#loaded

    const summary: ImagePageSummary = this.#dump.pages[0].summary

    const images: Image[] = []

    for (const page of this.#dump.pages) {
      for (const image of page.images) {
        if (options?.query) {
          const referenceMatches = image.reference.includes(options.query)
          const descriptionMatches = image.description
            ? image.description.includes(options.query)
            : false
          if (!referenceMatches && !descriptionMatches) {
            continue
          }
        }

        if (options?.tags && options.tags.length > 0) {
          const matchingTags = image.tags.reduce(
            (matchingTags, x) =>
              matchingTags + (options.tags?.includes(x) ? 1 : 0),
            0
          )

          const requiredMatches =
            options?.tagop === 'or' ? 1 : options.tags.length
          if (matchingTags < requiredMatches) {
            continue
          }
        }

        images.push(image)
      }
    }

    // The default sort / bump is the order the dump was created in - so we just
    // need to take care of the other modes
    if (options?.sort === 'reference') {
      images.sort((a, b) => a.reference.localeCompare(b.reference))
    }

    // TODO: Isn't this inverted for bump sort?
    if (options?.order === 'desc') {
      images.reverse()
    }

    const pageSize = options?.limit || 30
    const page = options?.page || 0

    const pagination: PaginationMetadata = {
      total: images.length,
      page: page + 1,
      size: pageSize,
    }

    return Promise.resolve({
      summary,
      images: images.slice(page * pageSize, (page + 1) * pageSize),
      pagination,
    })
  }

  async getImage(reference: string): Promise<Image | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]?.getImage as Image | null
    )
  }

  async getImageDescription(
    reference: string
  ): Promise<ImageDescription | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]
        ?.getImageDescription as ImageDescription | null
    )
  }

  async getImageReleaseNotes(
    reference: string
  ): Promise<ImageReleaseNotes | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]
        ?.getImageReleaseNotes as ImageReleaseNotes | null
    )
  }

  async getImageGraph(reference: string): Promise<Graph | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]?.getImageGraph as Graph | null
    )
  }

  async getImageScorecard(reference: string): Promise<ImageScorecard | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]
        ?.getImageScorecard as ImageScorecard | null
    )
  }

  async getImageProvenance(reference: string): Promise<ImageProvenance | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]
        ?.getImageProvenance as ImageProvenance | null
    )
  }

  async getImageSBOM(reference: string): Promise<ImageSBOM | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]?.getImageSBOM as ImageSBOM | null
    )
  }

  async getImageVulnerabilities(
    reference: string
  ): Promise<Vulnerability[] | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]?.getImageVulnerabilities as
        | Vulnerability[]
        | null
    )
  }

  async getLatestImageWorkflow(reference: string): Promise<WorkflowRun | null> {
    await this.#loaded
    return Promise.resolve(
      this.#dump.resources[reference]
        ?.getLatestImageWorkflow as WorkflowRun | null
    )
  }

  getLogoUrl(reference: string): string | undefined {
    // NOTE: This API is synchronous. As only the demo code requires
    // asynchronicity, don't wait for the promise.
    // In practice, the logo should not be retrieved before the rest of the
    // image data has been retrieved, meaning the demo data has loaded - so we
    // should be fine
    return this.#dump.resources[reference]?.getLogo as string | undefined
  }

  async scheduleImageScan(reference: string): Promise<void> {
    return Promise.resolve()
  }
}
