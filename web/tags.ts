export interface Tag {
  name: string
  description?: string
  color?: string | { light: string; dark: string }
}

const palette: Record<string, string | { light: string; dark: string }> = {
  // Brands
  kubernetes: { light: '#316ce6', dark: '#1955cc' },
  docker: { light: '#1d63ed', dark: '#104dc6' },
  github: { light: '#8B57E8', dark: '#571ac1' },
  gitlab: { light: '#fc6d26', dark: '#c94503' },
  quay: { light: '#40B4E5', dark: '#177ba6' },

  // Generic
  // SEE: https://spectrum.adobe.com/page/badge/
  grey: { light: '#6d6d6d', dark: '#545454' },
  purple1: { light: '#893de7', dark: '#6f1ad5' },
  purble2: { light: '#b622b7', dark: '#8e1a8e' },
  green1: { light: '#318801', dark: '#246501' },
  green2: { light: '#007772', dark: '#004d49' },
  red1: { light: '#c82269', dark: '#9d1b53' },
  red2: { light: '#d31510', dark: '#8e0f0b' },
  blue1: { light: '#0265dc', dark: '#024597' },
  blue2: { light: '#5258e4', dark: '#1c21b0' },
  yellow: { light: '#e8c600', dark: '#998200' },
}

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
]

export const TagsByName: Record<string, Tag> = Object.fromEntries(
  Tags.map((x) => [x.name, x])
)