export function version(reference: string): string {
  const parts = reference.split(':')
  if (parts.length === 1) {
    return 'latest'
  }
  return parts[1]
}
export function name(reference: string): string {
  const parts = reference.split(':')
  return parts[0]
}
