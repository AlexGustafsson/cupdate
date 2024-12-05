export function formatRelativeTimeTo(date: Date): string {
  let diff = Date.now() - date.getTime()

  const prefix = diff < 0 ? 'in ' : ''
  const suffix = diff >= 0 ? ' ago' : ''

  diff = Math.floor(Math.abs(diff) / 1000)

  const seconds = diff
  if (seconds == 0) {
    return 'just now'
  } else if (seconds < 60) {
    return `${prefix}${Math.floor(seconds)} second${seconds > 1 ? 's' : ''}${suffix}`
  }

  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) {
    return `${prefix}${Math.floor(minutes)} minute${minutes > 1 ? 's' : ''}${suffix}`
  }

  const hours = Math.floor(minutes / 60)
  if (hours < 24) {
    return `${prefix}${Math.floor(hours)} hour${hours > 1 ? 's' : ''}${suffix}`
  }

  const days = Math.floor(hours / 24)
  if (days < 30) {
    return `${prefix}${Math.floor(days)} day${days > 1 ? 's' : ''}${suffix}`
  }

  const months = Math.floor(days / 30)
  if (months < 12) {
    return `${prefix}${Math.floor(months)} month${months > 1 ? 's' : ''}${suffix}`
  }

  const years = Math.floor(days / 365)
  return `${prefix}${Math.floor(years)} year${years > 1 ? 's' : ''}${suffix}`
}
