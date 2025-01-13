function parse(reference: string): {
  name: string
  tag: string
  digest: string
} {
  let digest = ''
  const digestDelimiter = reference.indexOf('@')
  if (digestDelimiter >= 0) {
    digest = reference.substring(digestDelimiter + 1)
    reference = reference.substring(0, digestDelimiter)
  }

  let tag = ''
  const tagDelimiter = reference.indexOf(':')
  if (tagDelimiter >= 0) {
    digest = reference.substring(tagDelimiter + 1)
    tag = reference.substring(0, tagDelimiter)
  }

  return {
    name: reference,
    tag,
    digest,
  }
}

export function version(reference: string): string {
  const { tag, digest } = parse(reference)
  if (digest.length > 0) {
    return digest
  } else if (tag.length > 0) {
    return tag
  } else {
    return 'latest'
  }
}

export function name(reference: string): string {
  const { name } = parse(reference)
  return name
}
