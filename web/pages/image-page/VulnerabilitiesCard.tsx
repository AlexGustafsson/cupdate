import type { JSX } from 'react'
import React from 'react'
import { FluentBug16Regular } from '../../components/icons/fluent-bug-16-regular'
import { Markdown } from '../../components/Markdown'
import type { AffectedPackage, Vulnerability } from '../../lib/osv/osv'
import { compareSeverity, normalizedSeverity } from '../../lib/osv/severity'
import { parsePurl, purlLink, purlType } from '../../lib/purl'
import { formatRelativeTimeTo } from '../../time'
import { Card } from './Card'

export type VulnerabilitiesCardProps = {
  vulnerabilities: Vulnerability[]
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

function AffectedPackageDetails({
  package: pkg,
}: {
  package: AffectedPackage
}): JSX.Element {
  const purl = pkg.purl ? parsePurl(pkg.purl) : undefined
  const href = purl ? purlLink(purl) : undefined
  const type = purl ? purlType(purl) : undefined

  const ref = [pkg.name, purl?.version].filter((x) => x).join(' ')

  return (
    <>
      {href ? (
        <a href={href} target="_blank" rel="noreferrer">
          {ref}
        </a>
      ) : (
        ref
      )}
      {type && <> ({type})</>}
    </>
  )
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
              {vulnerabilities.map((x) => (
                <React.Fragment key={x.id}>
                  <h2>
                    {x.id} ({normalizedSeverity(x)})
                  </h2>
                  {x.published && (
                    <p>
                      Published {formatRelativeTimeTo(new Date(x.published))}.
                    </p>
                  )}
                  {x.modified !== x.published && (
                    <p>
                      Modified {formatRelativeTimeTo(new Date(x.modified))}.
                    </p>
                  )}
                  {x.summary && (
                    <>
                      <h3>Summary</h3>
                      <Markdown>{x.summary}</Markdown>
                    </>
                  )}
                  {x.affected && x.affected.length > 0 && (
                    <>
                      <h3>Affected packages</h3>
                      <ul>
                        {x.affected
                          ?.filter((x) => x.package)
                          .map((x) => (
                            <li key={x.package!.purl || x.package!.name}>
                              <AffectedPackageDetails package={x.package!} />
                            </li>
                          ))}
                      </ul>
                    </>
                  )}
                  <h3>References</h3>
                  <ul>
                    {/* Database-specific references */}
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
                    {/* Included references */}
                    {x.references?.map((x) => (
                      <li key={x.url}>
                        <a href={x.url} target="_blank" rel="noreferrer">
                          {x.url}
                        </a>
                      </li>
                    ))}
                    {/* NIST */}
                    {(x.id.startsWith('CVE-') || x.id.startsWith('cve-')) && (
                      <li>
                        <a
                          href={`https://nvd.nist.gov/vuln/detail/${x.id}`}
                          target="_blank"
                          rel="noreferrer"
                        >
                          {`https://nvd.nist.gov/vuln/detail/${x.id}`}
                        </a>
                      </li>
                    )}
                    {/* OSV */}
                    <li>
                      <a
                        href={`https://osv.dev/vulnerability/${x.id}`}
                        target="_blank"
                        rel="noreferrer"
                      >
                        {`https://osv.dev/vulnerability/${x.id}`}
                      </a>
                    </li>
                  </ul>
                </React.Fragment>
              ))}
            </div>
          ),
        },
      ]}
    />
  )
}
