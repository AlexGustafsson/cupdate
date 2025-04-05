/** Tag holds information about a label / tag used to categorize images. */
export interface Tag {
  /** Name is the human-readable name of the tag. */
  name: string
  /** Description optionally holds a description for the tag. */
  description?: string
  /**
   * Color optionally holds a CSS color or light/dark color combo for the tag.
   */
  color?: string | { light: string; dark: string }
}

const palette: Record<string, string | { light: string; dark: string }> = {
  // Brand colors
  kubernetes: { light: '#316ce6', dark: '#1955cc' },
  docker: { light: '#1d63ed', dark: '#104dc6' },
  github: { light: '#8B57E8', dark: '#571ac1' },
  gitlab: { light: '#fc6d26', dark: '#c94503' },
  quay: { light: '#40B4E5', dark: '#177ba6' },

  // Emotions
  negative: {
    // Tailwind's red-500
    light: 'oklch(63.7% 0.237 25.331)',
    // Tailwind's red-700
    dark: 'oklch(50.5% 0.213 27.518)',
  },
  warning: {
    // Tailwind's amber-500
    light: 'oklch(82.8% 0.189 84.429)',
    // Tailwind's amber-700
    dark: 'oklch(55.5% 0.163 48.998)',
  },
  // Akin to a red color without being red
  maybeWarning: {
    // Tailwind's purple-500
    light: 'oklch(62.7% 0.265 303.9)',
    dark: 'oklch(49.6% 0.265 301.924)',
  },
  positive: {
    // Tailwind's green-500
    light: 'oklch(72.3% 0.219 149.579)',
    // Tailwind's green-700
    dark: 'oklch(52.7% 0.154 150.069)',
  },

  // Bumps - chosen to have a growing sense of urgency if the bump is higher,
  // but without using emotive colors (red, blue, green). That is, major is akin
  // to a red color without being red, minor blue and patch green.
  major: {
    // Tailwind's purple-500
    light: 'oklch(62.7% 0.265 303.9)',
    // Tailwind's purple-700
    dark: 'oklch(49.6% 0.265 301.924)',
  },
  minor: {
    // Tailwind's violet-500
    light: 'oklch(60.6% 0.25 292.717)',
    // Tailwind's violet-700
    dark: 'oklch(49.1% 0.27 292.581)',
  },
  patch: {
    // Tailwind's indigo-500
    light: 'oklch(58.5% 0.233 277.117)',
    // Tailwind's indigo-700
    dark: 'oklch(45.7% 0.24 277.023)',
  },

  // Generic
  blue: {
    // Tailwind's blue-500
    light: 'oklch(62.3% 0.214 259.815)',
    // Tailwind's blue-700
    dark: 'oklch(48.8% 0.243 264.376)',
  },
}

const KubernetesTags: Tag[] = [
  {
    name: 'pod',
    description: 'A kubernetes pod',
    color: palette.kubernetes,
  },
  {
    name: 'job',
    description: 'A kubernetes job',
    color: palette.kubernetes,
  },
  {
    name: 'cron job',
    description: 'A kubernetes cron job',
    color: palette.kubernetes,
  },
  {
    name: 'deployment',
    description: 'A kubernetes deployment',
    color: palette.kubernetes,
  },
  {
    name: 'replica set',
    description: 'A kubernetes replica set',
    color: palette.kubernetes,
  },
  {
    name: 'daemon set',
    description: 'A kubernetes daemon set',
    color: palette.kubernetes,
  },
  {
    name: 'stateful set',
    description: 'A kubernetes stateful set',
    color: palette.kubernetes,
  },
]

const DockerTags: Tag[] = [
  {
    name: 'service',
    description: 'A Docker Swarm service',
    color: palette.docker,
  },
  {
    name: 'task',
    description: 'A Docker Swarm task',
    color: palette.docker,
  },
  {
    name: 'docker',
    description: 'A docker image',
    color: palette.docker,
  },
]

/** Holds known / well-defined tags in their intended sort order. */
export const Tags: Tag[] = [
  // Vulnerability warning
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'vulnerable',
    description: 'Vulnerable images',
    color: palette.negative,
  },
  {
    name: 'vulnerability:critical',
    description: 'Critical severity vulnerabilities discovered',
    color: palette.negative,
  },
  {
    name: 'vulnerability:high',
    description: 'High severity vulnerabilities discovered',
    color: palette.negative,
  },

  // Outdated, failure status
  {
    name: 'outdated',
    description: 'Outdated images',
    color: palette.negative,
  },
  {
    name: 'failed',
    description: 'Failed images',
    color: palette.negative,
  },

  // Security warnings
  {
    name: 'vulnerability:medium',
    description: 'Medium severity vulnerabilities discovered',
    color: palette.warning,
  },
  {
    name: 'vulnerability:low',
    description: 'Medium severity vulnerabilities discovered',
    color: palette.warning,
  },
  {
    name: 'risk:high',
    description: 'High risk project',
    color: palette.negative,
  },
  {
    name: 'risk:medium',
    description: 'Medium risk project',
    color: palette.warning,
  },
  {
    name: 'vulnerability:unspecified',
    description: 'Unspecified severity vulnerabilities discovered',
    color: palette.maybeWarning,
  },

  // Bump
  {
    name: 'bump:major',
    description: 'Major update',
    color: palette.major,
  },
  {
    name: 'bump:minor',
    description: 'Minor update',
    color: palette.minor,
  },
  {
    name: 'bump:patch',
    description: 'Patch update',
    color: palette.patch,
  },
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'major',
    description: 'Major update',
    color: palette.major,
  },
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'minor',
    description: 'Minor update',
    color: palette.minor,
  },
  {
    // TODO: Remove in v1. Replaced with prefixed version
    name: 'patch',
    description: 'Patch update',
    color: palette.patch,
  },

  // Status
  {
    name: 'up-to-date',
    description: 'Up-to-date images',
    color: palette.positive,
  },

  // Security information
  {
    name: 'risk:low',
    description: 'Low risk project',
    color: palette.blue,
  },
  {
    name: 'attestation',
    description: 'This image contains attestations',
    color: palette.blue,
  },
  {
    name: 'sbom',
    description: 'This image contains an SBOM',
    color: palette.positive,
  },

  // VCS
  {
    name: 'ghcr',
    description: 'A ghcr image',
    color: palette.github,
  },
  {
    name: 'github',
    description: 'A github project',
    color: palette.github,
  },

  // Deployment information
  ...KubernetesTags,
  ...DockerTags,
  {
    name: 'namespace:*',
    description: 'A namespace',
    color: palette.blue,
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
  } else if (aSort < bSort) {
    return 1
  } else {
    return a.localeCompare(b)
  }
}
