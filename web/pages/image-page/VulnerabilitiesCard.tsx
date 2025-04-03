import type { JSX } from 'react'
import type { ImageVulnerability } from '../../api'
import { FluentBug16Regular } from '../../components/icons/fluent-bug-16-regular'
import { Card } from './Card'

export type VulnerabilitiesCardProps = {
  vulnerabilities: ImageVulnerability[]
}

function unique<T>(previousValue: T[], currentValue: T): T[] {
  if (previousValue.includes(currentValue)) {
    return previousValue
  }

  previousValue.push(currentValue)
  return previousValue
}

function flattened<T>(previousValue: T[], currentValue: T[]): T[] {
  previousValue.push(...currentValue)
  return previousValue
}

type VulnerabilityCount = {
  critical: number
  high: number
  medium: number
  low: number
  unspecified: number
}

function countVulnerabilities(
  vulnerabilities: ImageVulnerability[]
): VulnerabilityCount {
  const counts: VulnerabilityCount = {
    critical: 0,
    high: 0,
    medium: 0,
    low: 0,
    unspecified: 0,
  }

  for (const vulnerability of vulnerabilities) {
    if (Object.hasOwn(counts, vulnerability.severity)) {
      counts[vulnerability.severity as keyof VulnerabilityCount]++
    } else {
      counts.unspecified++
    }
  }

  return counts
}

export function VulnerabilitiesCard({
  vulnerabilities,
}: VulnerabilitiesCardProps): JSX.Element {
  const counts = countVulnerabilities(vulnerabilities)

  return (
    <Card
      persistenceKey="vulnerabilities"
      tabs={[
        {
          icon: <FluentBug16Regular />,
          label: 'Vulnerabilities',
          content: (
            <div className="markdown-body">
              <h2>Summary</h2>
              <ul>
                {counts.critical > 0 && <li>Critical: {counts.critical}</li>}
                {counts.high > 0 && <li>High: {counts.high}</li>}
                {counts.medium > 0 && <li>Medium: {counts.medium}</li>}
                {counts.low > 0 && <li>Low: {counts.low}</li>}
                {counts.unspecified > 0 && (
                  <li>Unspecified: {counts.unspecified}</li>
                )}
              </ul>

              <h2>Authorities</h2>
              <ul>
                {vulnerabilities
                  .map((x) => x.authority)
                  .reduce(unique<string>, [])
                  .map((x) => (
                    <li key={x}>{x}</li>
                  ))}
              </ul>

              <h2>Links</h2>
              <ul>
                {vulnerabilities
                  .map((x) => x.links)
                  .reduce(flattened<string>, [])
                  .reduce(unique<string>, [])
                  .map((x) => (
                    <li key={x}>
                      <a href={x} target="_blank" rel="noreferrer">
                        {x}
                      </a>
                    </li>
                  ))}
              </ul>
            </div>
          ),
        },
      ]}
    />
  )
}
