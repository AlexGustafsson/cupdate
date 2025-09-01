export interface Purl {
  scheme: string
  type: string
  namespace?: string
  name: string
  version?: string
  qualifiers?: Record<string, string>
  subpath?: string
}

export function parsePurl(purl: string): Purl | null {
  try {
    const url = new URL(purl)
    if (url.protocol !== 'pkg:') {
      return null
    }

    const [path, version] = url.pathname.split('@')
    const parts = path.split('/').map(decodeURIComponent)
    const type = parts[0]
    const name = parts.length === 2 ? parts[1] : parts.slice(2).join('/')
    const namespace = parts.length > 2 ? parts[1] : undefined

    const result: Purl = {
      scheme: 'pkg',
      type,
      name,
    }

    if (namespace) {
      result.namespace = namespace
    }

    if (version) {
      result.version = version
    }

    if (url.search !== '') {
      result.qualifiers = Object.fromEntries(url.searchParams.entries())
    }

    if (url.hash !== '') {
      result.subpath = url.hash.substring(1)
    }

    return result
  } catch {
    return null
  }
}

export function purlLink(purl: Purl): string | null {
  switch (purl.type) {
    case 'apk':
      switch (purl.namespace) {
        case 'alpine':
          return purlLinkAlpine(purl)
      }
      break
    case 'deb':
      switch (purl.namespace) {
        case 'ubuntu':
          return purlLinkUbuntu(purl)
      }
      break
    case 'golang':
      return purlLinkGolang(purl)
  }

  return null
}

export function purlType(purl: Purl): string | null {
  switch (purl.type) {
    case 'apk':
      switch (purl.namespace) {
        case 'alpine':
          return 'an Alpine package'
      }
      break
    case 'deb':
      switch (purl.namespace) {
        case 'ubuntu':
          return 'an Ubuntu package'
      }
      break
    case 'golang':
      switch (purl.name) {
        case 'stdlib':
          return 'the Golang standard library'
        default:
          return 'a Golang package'
      }
  }

  return null
}

function purlLinkAlpine(purl: Purl): string {
  const url = new URL('https://pkgs.alpinelinux.org/packages')
  url.searchParams.append('name', purl.name)
  const osVersion = purl.qualifiers?.os_version
  if (osVersion) {
    url.searchParams.append(
      'branch',
      osVersion.startsWith('v') ? osVersion : `v${osVersion}`
    )
  }
  return url.toString()
}

function purlLinkUbuntu(purl: Purl): string {
  let url = `https://launchpad.net/ubuntu/+source/${purl.name}`
  if (purl.version) {
    url += `/${purl.version.replace(/\+.*$/, '')}`
  }
  return url
}

function purlLinkGolang(purl: Purl): string | null {
  if (purl.name === 'stdlib') {
    return 'https://pkg.go.dev/std'
  }

  if (purl.namespace) {
    return `https://pkg.go.dev/${purl.namespace}/${purl.name}`
  }

  return null
}
