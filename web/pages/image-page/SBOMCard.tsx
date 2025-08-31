import type { JSX } from 'react'
import { FluentBoxSearch16Regular } from '../../components/icons/fluent-box-search-16-regular'
import type { ImageSBOM, SBOM } from '../../lib/api/models'
import { Card } from './Card'

type SBOMTabProps = {
  sbom: SBOM
}

function SBOMTab({ sbom }: SBOMTabProps): JSX.Element {
  return (
    <div className="markdown-body">
      <pre className="max-h-100">
        <code className="language-dockerfile">{sbom.sbom}</code>
      </pre>
    </div>
  )
}

export type SBOMCardProps = {
  sbom: ImageSBOM
}

export function SBOMCard({ sbom }: SBOMCardProps): JSX.Element {
  return (
    <Card
      persistenceKey="sbom"
      tabs={[
        {
          icon: <FluentBoxSearch16Regular />,
          label: 'SBOM',
          content: (
            <div className="markdown-body">
              <p>
                Some images include{' '}
                <i>Software Bill of Materials (SBOM) attestations</i> - means of
                asserting facts about an image's contents. These details are
                helpful for users to understand what software is running on
                their machines, as well as for services like Cupdate to
                automatically identify vulnerabilities.
              </p>
              <p>
                More information can be found here:{' '}
                <a
                  href="https://docs.docker.com/build/metadata/attestations/sbom/"
                  target="_blank"
                  rel="noreferrer"
                >
                  https://docs.docker.com/build/metadata/attestations/sbom/
                </a>
                .
              </p>
            </div>
          ),
        },
        ...sbom.sbom.map((sbom) => ({
          label:
            sbom.sbom.length > 0 && sbom.architecture
              ? `${[
                  sbom.operatingSystem,
                  sbom.architecture,
                  sbom.architectureVariant,
                ]
                  .filter((x) => x !== undefined)
                  .join('/')}`
              : 'SBOM',
          content: <SBOMTab sbom={sbom} />,
        })),
      ]}
    />
  )
}
