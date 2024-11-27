export interface Tag {
  name: string
  description?: string
  color?: string | { light: string; dark: string }
}

// SEE: https://designlibrary.sebgroup.com/patterns/pattern-darkmode
const palette = [
  ['#007AC7', '#2C9CD9'],
  ['#0092E1', '#01ADFF'],
  ['#40b0ed', '#6d58b8'],

  ['#318801', '#318801'],
  ['#45b401', '#75ba49'],
  ['#61cc18', '#61cc18'],

  ['#3f2587', '#4a328f'],
  ['#4f2c99', '#7e52cc'],
  ['#673ab5', '#ac91dc'],

  ['#f7a000', '#ebaa39'],
  ['#ffb400', '#f0be47'],
  ['#ffc502', '#ffe1b2'],

  ['#bb010b', '#9e2120'],
  ['#d81b1b', '#c82929'],
  ['#f0352a', '#f7706c'],
]

export const Tags: Tag[] = [
  {
    name: 'pod',
    description: 'A kubernetes pod',
    color: { light: palette[0][0], dark: palette[0][1] },
  },
  {
    name: 'job',
    description: 'A kubernetes job',
    color: { light: palette[1][0], dark: palette[1][1] },
  },
  {
    name: 'cron job',
    description: 'A kubernetes cron job',
    color: { light: palette[2][0], dark: palette[2][1] },
  },
  {
    name: 'deployment',
    description: 'A kubernetes deployment',
    color: { light: palette[7][0], dark: palette[7][1] },
  },
  {
    name: 'replica set',
    description: 'A kubernetes replica set',
    color: { light: palette[8][0], dark: palette[8][1] },
  },
  {
    name: 'daemon set',
    description: 'A kubernetes daemon set',
    color: { light: palette[9][0], dark: palette[9][1] },
  },
  {
    name: 'stateful set',
    description: 'A kubernetes stateful set',
    color: { light: palette[10][0], dark: palette[10][1] },
  },
  {
    name: 'docker',
    description: 'A docker image',
    color: { light: palette[7][0], dark: palette[7][1] },
  },
  {
    name: 'ghcr',
    description: 'A ghcr image',
    color: { light: palette[8][0], dark: palette[8][1] },
  },
  {
    name: 'github',
    description: 'A github project',
    color: { light: palette[0][0], dark: palette[0][1] },
  },
  {
    name: 'up-to-date',
    description: 'Up-to-date images',
    color: { light: palette[6][0], dark: palette[6][1] },
  },
  {
    name: 'outdated',
    description: 'Outdated images',
    color: { light: palette[13][0], dark: palette[13][1] },
  },
]

export const TagsByName: Record<string, Tag> = Object.fromEntries(
  Tags.map((x) => [x.name, x])
)
