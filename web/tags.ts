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

  // Generic colors
  // Originally based on https://spectrum.adobe.com/page/badge/
  grey: { light: '#6d6d6d', dark: '#545454' },
  purple1: { light: '#893de7', dark: '#6f1ad5' },
  purble2: { light: '#b622b7', dark: '#8e1a8e' },
  green1: { light: '#44cb01', dark: '#339801' },
  green2: { light: '#007772', dark: '#004d49' },
  red1: { light: '#de3b82', dark: '#9d1b53' },
  red2: { light: '#f14741', dark: '#8e0f0b' },
  blue1: { light: '#0265dc', dark: '#024597' },
  blue2: { light: '#5258e4', dark: '#1c21b0' },
  yellow: { light: '#e8c600', dark: '#998200' },
}

/** Holds known / well-defined tags. */
export const Tags: Tag[] = [
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
  {
    name: 'up-to-date',
    description: 'Up-to-date images',
    color: palette.green1,
  },
  {
    name: 'outdated',
    description: 'Outdated images',
    color: palette.red2,
  },
  {
    name: 'failed',
    description: 'Failed images',
    color: palette.red2,
  },
  {
    name: 'vulnerable',
    description: 'Vulnerable images',
    color: palette.red2,
  },
  {
    name: 'patch',
    description: 'Patch update',
    color: palette.blue1,
  },
  {
    name: 'minor',
    description: 'Minor update',
    color: palette.yellow,
  },
  {
    name: 'major',
    description: 'Major update',
    color: palette.red1,
  },
]

/** Holds known / well-defined tags, mapped by their name. */
export const TagsByName: Record<string, Tag> = Object.fromEntries(
  Tags.map((x) => [x.name, x])
)

/** Sort tags lexically, putting prefixed tags last, selected tags first. */
export function compareTags(
  a: string,
  b: string,
  aSelected?: boolean,
  bSelected?: boolean
): number {
  // Prioritize selected tags
  if (aSelected === true && bSelected === false) {
    return -1
  } else if (aSelected === false && bSelected === true) {
    return 1
  }

  const aString = typeof a === 'string' ? a : a
  const bString = typeof b === 'string' ? b : b

  if (aString.includes(':') === bString.includes(':')) {
    return aString.localeCompare(bString)
  } else if (aString.includes(':')) {
    return 1
  } else {
    return -1
  }
}
