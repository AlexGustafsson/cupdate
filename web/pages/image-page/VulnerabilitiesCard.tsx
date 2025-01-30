import type { JSX } from 'react'
import type { Image } from '../../api'

export type VulnerabilitiesCardProps = {
  image: Image
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

export function VulnerabilitiesCard({
  image,
}: VulnerabilitiesCardProps): JSX.Element {
  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow">
      <div className="markdown-body">
        <h1>Vulnerabilities</h1>
        <ul>
          <li>
            Critical:{' '}
            {
              image.vulnerabilities.filter((x) => x.severity === 'critical')
                .length
            }
          </li>
          <li>
            High:{' '}
            {image.vulnerabilities.filter((x) => x.severity === 'high').length}
          </li>
          <li>
            Medium:{' '}
            {
              image.vulnerabilities.filter((x) => x.severity === 'medium')
                .length
            }
          </li>
          <li>
            Low:{' '}
            {image.vulnerabilities.filter((x) => x.severity === 'low').length}
          </li>
          <li>
            Unspecified:{' '}
            {
              image.vulnerabilities.filter((x) => x.severity === 'unspecified')
                .length
            }
          </li>
        </ul>

        <h2>Authorities</h2>
        <ul>
          {image.vulnerabilities
            .map((x) => x.authority)
            .reduce(unique<string>, [])
            .map((x) => (
              <li key={x}>{x}</li>
            ))}
        </ul>

        <h2>Links</h2>
        <ul>
          {image.vulnerabilities
            .map((x) => x.links)
            .reduce(flattened<string>, [])
            .reduce(unique<string>, [])
            .map((x) => (
              <li key={x}>
                <a href={x}>{x}</a>
              </li>
            ))}
        </ul>
      </div>
    </div>
  )
}
