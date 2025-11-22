import type { Vulnerability } from '../osv/osv'
import type {
  Graph,
  Image,
  ImageDescription,
  ImagePage,
  ImageProvenance,
  ImageReleaseNotes,
  ImageSBOM,
  ImageScorecard,
  WorkflowRun,
} from './models'

// TODO: Provider for API calls. That way we can build with a mock API provider
// and optionally build the API with mock data. This way we can create a build
// with everything needed to run locally.
// TODO: There's a github pages action to publish a page. That way this can be
// build and pushed on push to main
export interface ApiClient {
  getTags(): Promise<string[]>
  getImages(options?: GetImagesOptions): Promise<ImagePage>
  getImage(reference: string): Promise<Image | null>
  getImageDescription(reference: string): Promise<ImageDescription | null>
  getImageReleaseNotes(reference: string): Promise<ImageReleaseNotes | null>
  getImageGraph(reference: string): Promise<Graph | null>
  getImageScorecard(reference: string): Promise<ImageScorecard | null>
  getImageProvenance(reference: string): Promise<ImageProvenance | null>
  getImageSBOM(reference: string): Promise<ImageSBOM | null>
  getImageVulnerabilities(reference: string): Promise<Vulnerability[] | null>
  getLatestImageWorkflow(reference: string): Promise<WorkflowRun | null>
  getLogoUrl(reference: string): string | undefined
  scheduleImageScan(reference: string): Promise<void>
  pollImages(): Promise<void>
}

export interface GetImagesOptions {
  tags?: string[]
  tagop?: 'and' | 'or'
  sort?: 'reference' | 'bump'
  order?: 'asc' | 'desc'
  page?: number
  limit?: number
  query?: string
}
