export interface Tag {
  name: string
  description?: string
  color?: string | { light: string; dark: string }
}

export const Tags: Tag[] = [
  {
    name: 'pod',
    description: 'A kubernetes pod',
    color: '#FFEDD5',
  },
  {
    name: 'job',
    description: 'A kubernetes job',
    color: '#DBEAFE',
  },
  {
    name: 'cron job',
    description: 'A kubernetes cron job',
    color: '#DBEAFE',
  },
  {
    name: 'deployment',
    description: 'A kubernetes deployment',
    color: '#DBEAFE',
  },
  {
    name: 'replica set',
    description: 'A kubernetes replica set',
    color: '#DBEAFE',
  },
  {
    name: 'daemon set',
    description: 'A kubernetes daemon set',
    color: '#DBEAFE',
  },
  {
    name: 'stateful set',
    description: 'A kubernetes stateful set',
    color: '#DBEAFE',
  },
  {
    name: 'docker',
    description: 'A docker image',
    color: '#FEE2E2',
  },
  {
    name: 'ghcr',
    description: 'A ghcr image',
    color: '#FEE2E2',
  },
  {
    name: 'github',
    description: 'A github project',
    color: '#FEE2E2',
  },
  {
    name: 'up-to-date',
    description: 'Up-to-date images',
    color: '#DCFCE7',
  },
  {
    name: 'outdated',
    description: 'Outdated images',
    color: '#FEE2E2',
  },
]

export const TagsByName: Record<string, Tag> = Object.fromEntries(
  Tags.map((x) => [x.name, x])
)
