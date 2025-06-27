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
  buildArguments?: Record<string, string>
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

export interface WebPushServerKey {
  key: string
}

export interface WebPushSubscription {
  endpoint: string
  keys: {
    p256dh: string
    auth: string
  }
}
