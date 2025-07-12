import { describe, expect } from 'vitest'
import { compareSeverity, type NormalizedSeverity } from './severity'

describe('NormalizedSeverity', (it) => {
  it('is comparable', () => {
    const expected: string[] = [
      'critical',
      'high',
      'medium',
      'low',
      'unspecified',
      'unknown or unsupported',
    ]

    // Property-based testing
    for (let i = 0; i < 100; i++) {
      const actual = [...expected]
        .map((x) => ({ value: x, sort: Math.random() }))
        .sort((a, b) => a.sort - b.sort)
        .map((x) => x.value)
        .sort((a, b) =>
          compareSeverity(a as NormalizedSeverity, b as NormalizedSeverity)
        )

      expect(actual).toEqual(expected)
    }
  })
})
