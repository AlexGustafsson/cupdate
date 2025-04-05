import type { Vulnerability } from './osv'

export type NormalizedSeverity =
  | 'critical'
  | 'high'
  | 'medium'
  | 'low'
  | 'unspecified'

export function normalizedSeverity(
  vulnerability: Vulnerability
): NormalizedSeverity {
  if (vulnerability.database_specific) {
    switch (vulnerability.database_specific.severity) {
      case 'CRITICAL':
        return 'critical'
      case 'HIGH':
        return 'high'
      case 'MODERATE':
      case 'MEDIUM':
        return 'medium'
      case 'LOW':
        return 'low'
    }
  }

  // TODO: Get the severities from packages / parse CVSS score?
  return 'unspecified'
}

export function compareSeverity(
  a: NormalizedSeverity,
  b: NormalizedSeverity
): number {
  if (a === b) {
    return 0
  }

  const order: NormalizedSeverity[] = [
    'critical',
    'high',
    'medium',
    'low',
    'unspecified',
  ]

  let orderA = order.indexOf(a)
  if (orderA === -1) {
    orderA = order.length
  }

  let orderB = order.indexOf(b)
  if (orderB === -1) {
    orderB = order.length
  }

  return orderA - orderB
}
