import type { JSX } from 'react'
import React from 'react'
import { FluentBug16Regular } from '../../components/icons/fluent-bug-16-regular'
import type { Vulnerability } from '../../lib/osv/osv'
import { compareSeverity, normalizedSeverity } from '../../lib/osv/severity'
import { Card } from './Card'

export type VulnerabilitiesCardProps = {
  vulnerabilities: Vulnerability[]
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
  vulnerabilities: Vulnerability[]
): VulnerabilityCount {
  const counts: VulnerabilityCount = {
    critical: 0,
    high: 0,
    medium: 0,
    low: 0,
    unspecified: 0,
  }

  for (const vulnerability of vulnerabilities) {
    const severity = normalizedSeverity(vulnerability)

    if (Object.hasOwn(counts, severity)) {
      counts[severity as keyof VulnerabilityCount]++
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

  vulnerabilities = vulnerabilities.sort((a, b) =>
    compareSeverity(normalizedSeverity(a), normalizedSeverity(b))
  )

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
            </div>
          ),
        },
        {
          label: 'Details',
          content: (
            <div className="markdown-body">
              <dl>
                {vulnerabilities.map((x) => (
                  <React.Fragment key={x.id}>
                    <dt>
                      {x.id} ({normalizedSeverity(x)})
                    </dt>
                    <dd>
                      {x.summary && <p>{x.summary}</p>}
                      <ul>
                        {typeof x.database_specific?.url === 'string' && (
                          <li key={x.database_specific.url}>
                            <a
                              href={x.database_specific.url}
                              target="_blank"
                              rel="noreferrer"
                            >
                              {x.database_specific.url}
                            </a>
                          </li>
                        )}
                        {x.references?.map((x) => (
                          <li key={x.url}>
                            <a href={x.url} target="_blank" rel="noreferrer">
                              {x.url}
                            </a>
                          </li>
                        ))}
                      </ul>
                    </dd>
                  </React.Fragment>
                ))}
              </dl>
            </div>
          ),
        },
      ]}
    />
  )
}
