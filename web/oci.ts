/**
 * Regex is taken from the distribution reference, which is under an Apache 2.0
 * license.
 * @see {@link https://github.com/distribution/reference/blob/8c942b0459dfdcc5b6685581dd0a5a470f615bff/regexp.go#L34}
 */
const ReferenceRegexp =
  /^((?:(?:(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])(?:\.(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]))*|\[(?:[a-fA-F0-9:]+)\])(?::[0-9]+)?\/)?[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*(?:\/[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*)*)(?::([\w][\w.-]{0,127}))?(?:@([A-Za-z][A-Za-z0-9]*(?:[-_+.][A-Za-z][A-Za-z0-9]*)*[:][0-9A-Fa-f]{32,}))?$/

function parse(reference: string): {
  name: string
  tag: string
  digest: string
} | null {
  const match = ReferenceRegexp.exec(reference)
  if (!match) {
    return null
  }

  return {
    name: match[1] || '',
    tag: match[2] || '',
    digest: match[3] || '',
  }
}

export function version(reference: string): string {
  const result = parse(reference)
  if (!result) {
    return ''
  }

  // For now, if both a tag and a digest is specified - show only the tag
  if (result.tag.length > 0) {
    return result.tag
  } else if (result.digest.length > 0) {
    return result.digest
  } else {
    return 'latest'
  }
}

export function fullVersion(reference: string): string {
  const result = parse(reference)
  if (!result) {
    return ''
  }

  if (result.digest.length > 0) {
    if (result.tag.length > 0) {
      return `${result.tag}@${result.digest}`
    }

    return result.digest
  } else if (result.tag.length > 0) {
    return result.tag
  } else {
    return 'latest'
  }
}

export function name(reference: string): string {
  const result = parse(reference)
  if (!result) {
    return ''
  }

  return result.name
}
