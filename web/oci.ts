/**
 * Regex is taken from the distribution reference, which is under an Apache 2.0
 * license.
 * @see {@link https://github.com/distribution/reference/blob/8c942b0459dfdcc5b6685581dd0a5a470f615bff/regexp.go#L34}
 */
const ReferenceRegexp =
  /^((?:(?:(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])(?:\.(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]))*|\[(?:[a-fA-F0-9:]+)\])(?::[0-9]+)?\/)?[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*(?:\/[a-z0-9]+(?:(?:[._]|__|[-]+)[a-z0-9]+)*)*)(?::([\w][\w.-]{0,127}))?(?:@([A-Za-z][A-Za-z0-9]*(?:[-_+.][A-Za-z][A-Za-z0-9]*)*[:][0-9A-Fa-f]{32,}))?$/

/** Parse an OCI reference. Returns null for invalid references.  */
export function parse(reference: string): {
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

/**
 * Returns a string representing the reference's version in a way users would
 * normally associate with an image.
 */
export function version(reference: string): string {
  const result = parse(reference)
  if (!result) {
    return ''
  }

  // For now, if both a tag and a digest is specified - show only the tag
  if (result.tag.length > 0) {
    return result.tag
  }

  if (result.digest.length > 0) {
    return result.digest
  }

  return 'latest'
}

/**
 * Returns a string representing the reference's version, including both tag and
 * digest, if available.
 */
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
  }

  if (result.tag.length > 0) {
    return result.tag
  }

  return 'latest'
}

/** Name returns the name of the reference. I.e. its registry and path combo. */
export function name(reference: string): string {
  const result = parse(reference)
  if (!result) {
    return ''
  }

  return result.name
}

/**
 * Returns a string representing the reference's version in a way users would
 * normally associate with an image. If that ends up being a digest or the
 * 'latest' tag, any known version annotation will be included in the string.
 */
export function formattedVersion(
  reference: string,
  annotations?: Record<string, string>
): string {
  let versionString = version(reference)
  if (
    (versionString === 'latest' || versionString.startsWith('sha256:')) &&
    annotations
  ) {
    const versionAnnotationString = versionAnnotation(annotations)
    if (versionAnnotationString && versionAnnotationString !== versionString) {
      versionString = `${versionString} (${versionAnnotationString})`
    }
  }

  return versionString
}

export function versionAnnotation(
  annotations: Record<string, string>
): string | undefined {
  return (
    annotations['org.opencontainers.image.version'] ||
    annotations['org.label-schema.version'] ||
    annotations.version
  )
}
