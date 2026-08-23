export type Color =
  | 'brand-kubernetes'
  | 'brand-docker'
  | 'brand-github'
  | 'brand-gitlab'
  | 'brand-quay'
  | 'negative'
  | 'warning'
  | 'maybe-warning'
  | 'positive'
  | 'major'
  | 'minor'
  | 'patch'
  | 'accent'

/** Tag holds information about a label / tag used to categorize images. */
export interface Tag {
  /** Name is the human-readable name of the tag. */
  name: string
  /** Description optionally holds a description for the tag. */
  description?: string
  /**
   * Color optionally holds a reference to a known color.
   */
  color?: Color
}

const KubernetesTags: Tag[] = [
  {
    name: 'pod',
    description: 'A kubernetes pod',
    color: 'brand-kubernetes',
  },
  {
    name: 'job',
    description: 'A kubernetes job',
    color: 'brand-kubernetes',
  },
  {
    name: 'cron job',
    description: 'A kubernetes cron job',
    color: 'brand-kubernetes',
  },
  {
    name: 'deployment',
    description: 'A kubernetes deployment',
    color: 'brand-kubernetes',
  },
  {
    name: 'replica set',
    description: 'A kubernetes replica set',
    color: 'brand-kubernetes',
  },
  {
    name: 'daemon set',
    description: 'A kubernetes daemon set',
    color: 'brand-kubernetes',
  },
  {
    name: 'stateful set',
    description: 'A kubernetes stateful set',
    color: 'brand-kubernetes',
  },
]

const DockerTags: Tag[] = [
  {
    name: 'service',
    description: 'A Docker Swarm service',
    color: 'brand-docker',
  },
  {
    name: 'task',
    description: 'A Docker Swarm task',
    color: 'brand-docker',
  },
  {
    name: 'brand-docker',
    description: 'A docker image',
    color: 'brand-docker',
  },
]

/** Holds known / well-defined tags in their intended sort order. */
export const Tags: Tag[] = [
  // Vulnerability warning
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'vulnerable',
    description: 'Vulnerable images',
    color: 'negative',
  },
  {
    name: 'vulnerability:critical',
    description: 'Critical severity vulnerabilities discovered',
    color: 'negative',
  },
  {
    name: 'vulnerability:high',
    description: 'High severity vulnerabilities discovered',
    color: 'negative',
  },

  // Outdated, failure status
  {
    name: 'outdated',
    description: 'Outdated images',
    color: 'negative',
  },
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'failed',
    description: 'Failed images',
    color: 'negative',
  },
  {
    name: 'status:failed',
    description: 'Failed images',
    color: 'negative',
  },
  {
    name: 'status:incomplete',
    description: 'Incomplete information',
    color: 'warning',
  },

  // Security warnings
  {
    name: 'vulnerability:medium',
    description: 'Medium severity vulnerabilities discovered',
    color: 'warning',
  },
  {
    name: 'vulnerability:low',
    description: 'Medium severity vulnerabilities discovered',
    color: 'warning',
  },
  {
    name: 'risk:high',
    description: 'High risk project',
    color: 'negative',
  },
  {
    name: 'risk:medium',
    description: 'Medium risk project',
    color: 'warning',
  },
  {
    name: 'vulnerability:unspecified',
    description: 'Unspecified severity vulnerabilities discovered',
    color: 'maybe-warning',
  },

  // Bump
  {
    name: 'bump:major',
    description: 'Major update',
    color: 'major',
  },
  {
    name: 'bump:minor',
    description: 'Minor update',
    color: 'minor',
  },
  {
    name: 'bump:patch',
    description: 'Patch update',
    color: 'patch',
  },
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'major',
    description: 'Major update',
    color: 'major',
  },
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'minor',
    description: 'Minor update',
    color: 'minor',
  },
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'patch',
    description: 'Patch update',
    color: 'patch',
  },

  // Status
  {
    name: 'up-to-date',
    description: 'Up-to-date images',
    color: 'positive',
  },
  {
    name: 'unprocessed',
    description: 'Unprocessed images',
    color: 'warning',
  },

  // Security information
  {
    name: 'risk:low',
    description: 'Low risk project',
    color: 'accent',
  },
  {
    name: 'attestation',
    description: 'This image contains attestations',
    color: 'accent',
  },
  {
    name: 'sbom',
    description: 'This image contains an SBOM',
    color: 'positive',
  },

  // VCS
  {
    name: 'ghcr',
    description: 'A ghcr image',
    color: 'brand-github',
  },
  {
    name: 'github',
    description: 'A github project',
    color: 'brand-github',
  },
  {
    name: 'gitlab',
    description: 'A gitlab project',
    color: 'brand-gitlab',
  },
  {
    name: 'quay',
    description: 'A quay project',
    color: 'brand-quay',
  },

  // Deployment information
  ...KubernetesTags,
  ...DockerTags,
  {
    name: 'namespace:*',
    description: 'A namespace',
    color: 'accent',
  },
]

/** Holds known / well-defined tags, mapped by their name. */
const TagsByName: Record<string, Tag> = Object.fromEntries(
  Tags.map((x) => [x.name, x])
)

export function tagByName(name: string): Tag | undefined {
  // Get tag by name
  let tag = TagsByName[name]

  // Fall back to a tag prefix
  if (!tag && name.includes(':')) {
    tag = TagsByName[`${name.substring(0, name.indexOf(':'))}:*`]
  }

  return tag
}

function tagSortValue(name: string, selected?: boolean): number {
  const tag = tagByName(name)

  const values: number[] = [
    // Prioritize selected tags
    selected ? 1 : 0,
    // Priority based on definition order
    tag ? Tags.length - Tags.indexOf(tag) : 0,
  ]

  let value = 0
  for (let i = 0; i < values.length; i++) {
    value |= values[i] << Math.floor(32 / (values.length - i))
  }

  return value
}

/** Sort tags lexically, putting prefixed tags last, selected tags first. */
export function compareTags(
  a: string,
  b: string,
  aSelected?: boolean,
  bSelected?: boolean
): number {
  const aSort = tagSortValue(a, aSelected)
  const bSort = tagSortValue(b, bSelected)

  if (aSort > bSort) {
    return -1
  }
  if (aSort < bSort) {
    return 1
  }
  return a.localeCompare(b)
}
