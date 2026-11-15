import * as CycloneDX from './cyclonedx'
import * as SPDX2 from './spdx2'

export interface SBOM {}

// TODO: Idea is that we can normalize SBOMs to essentially be an array of root
// entries (think frontend, cupdate, osv-scanner) and a map of id: package.
// A package would be a name, description and included packages. Goal is to list
// vulnerabilities by tree, e.g. foobar@123 has vulnerability, attribute to it
// being used by dompurify, being used by frontend. In order to help user
// understanding of probable vulnerabilities
export function normalize(sbom: string): SBOM | null {
  let x: {
    spdxVersion?: string
    bomFormat?: string
    specVersion?: string
  }
  try {
    x = JSON.parse(sbom)
  } catch {
    return null
  }

  switch (x.spdxVersion) {
    case 'SPDX-2.3':
      return normalizeSPDX2(x as SPDX2.Document)
  }

  switch (x.bomFormat) {
    case 'CycloneDX':
      switch (x.specVersion) {
        case '1.6':
          return normalizeCycloneDX(x as CycloneDX.Document)
      }
  }

  return null
}

function normalizeSPDX2(sbom: SPDX2.Document): SBOM | null {
  return null
}

function normalizeCycloneDX(sbom: CycloneDX.Document): SBOM | null {
  return null
}
