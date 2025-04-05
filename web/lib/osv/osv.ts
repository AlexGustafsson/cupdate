export interface Vulnerability {
  id: string
  modified: string
  affected?: Affected[]
  aliases?: string[]
  credits?: Credit[]
  database_specific?: Record<string, unknown>
  details?: string
  published?: string
  references?: Reference[]
  related?: string[]
  schema_version?: string
  severity?: Severity[]
  summary?: string
  withdrawn?: string
}

export interface Affected {
  database_specific?: Record<string, unknown>
  ecosystem_Specific?: Record<string, unknown>
  package?: AffectedPackage
  ranges?: AffectedRange[]
  severity?: Severity[]
  versions?: string[]
}

export interface AffectedPackage {
  ecosystem: string
  name: string
  purl?: string
}

export interface AffectedRange {
  type: string
  database_specific?: Record<string, unknown>
  events?: Event[]
  repo?: string
}

export interface Event {
  introduced?: string
  fixed?: string
  last_affected?: string
  limit?: string
}

export interface Credit {
  name: string
  contact?: string[]
  type?: string
}

export type ReferenceType =
  | 'ADVISORY'
  | 'ARTICLE'
  | 'DETECTION'
  | 'DISCUSSION'
  | 'REPORT'
  | 'FIX'
  | 'GIT'
  | 'INTRODUCED'
  | 'PACKAGE'
  | 'EVIDENCE'
  | 'WEB'

export interface Reference {
  type: ReferenceType
  url: string
}

export interface Severity {
  type: string
  score: string
}
