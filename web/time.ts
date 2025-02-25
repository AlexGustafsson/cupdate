/**
 * Returns a human-readable string representing the current relative time to
 * the specified date.
 */
export function formatRelativeTimeTo(date: Date): string {
  let diff = Date.now() - date.getTime()

  const prefix = diff < 0 ? 'in ' : ''
  const suffix = diff >= 0 ? ' ago' : ''

  diff = Math.floor(Math.abs(diff) / 1000)

  if (diff === 0) {
    return 'just now'
  }

  return prefix + formatDuration(diff) + suffix
}

/** Returns a human-readable string representing the duration. */
export function formatDuration(seconds: number): string {
  if (seconds < 1) {
    return 'less than a second'
  }

  if (seconds < 60) {
    return `${Math.floor(seconds)} second${seconds > 1 ? 's' : ''}`
  }

  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) {
    return `${Math.floor(minutes)} minute${minutes > 1 ? 's' : ''}`
  }

  const hours = Math.floor(minutes / 60)
  if (hours < 24) {
    return `${Math.floor(hours)} hour${hours > 1 ? 's' : ''}`
  }

  const days = Math.floor(hours / 24)
  if (days < 30) {
    return `${Math.floor(days)} day${days > 1 ? 's' : ''}`
  }

  const months = Math.floor(days / 30)
  if (months < 12) {
    return `${Math.floor(months)} month${months > 1 ? 's' : ''}`
  }

  const years = Math.floor(days / 365)
  return `${Math.floor(years)} year${years > 1 ? 's' : ''}`
}
