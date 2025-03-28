import type { JSX } from 'react'
import type { ImageProvenance, ProvenanceBuildInfo } from '../../api'
import { FluentDocumentRibbon16Regular } from '../../components/icons/fluent-document-ribbon-16-regular'
import { Card } from './Card'

type BuildInfoProps = {
  buildInfo: ProvenanceBuildInfo
}

function BuildInfo({ buildInfo }: BuildInfoProps): JSX.Element {
  return (
    <div className="markdown-body">
      {buildInfo.dockerfile && (
        <pre className="max-h-100">
          <code className="language-dockerfile">{buildInfo.dockerfile}</code>
        </pre>
      )}
      <ul>
        <li>
          <div className="flex items-center">
            <p className="m-0 flex-shrink-0">Image manifest: </p>
            <code className="truncate block">{buildInfo.imageDigest}</code>
          </div>
        </li>
        <li>
          <div className="flex items-center">
            <p className="m-0 flex-shrink-0">Build started: </p>
            <code className="truncate block">{buildInfo.buildStartedOn}</code>
          </div>
        </li>
        <li>
          <div className="flex items-center">
            <p className="m-0 flex-shrink-0">Build finished: </p>
            <code className="truncate block">{buildInfo.buildFinishedOn}</code>
          </div>
        </li>
        <li>
          <div className="flex items-center">
            <p className="m-0 flex-shrink-0">Architecture: </p>
            <code className="truncate block">
              {[
                buildInfo.operatingSystem,
                buildInfo.architecture,
                buildInfo.architectureVariant,
              ]
                .filter((x) => x !== undefined)
                .join('/')}
            </code>
          </div>
        </li>
        {buildInfo.source && (
          <li>
            <div className="flex items-center">
              <p className="m-0 flex-shrink-0">Source: </p>
              <code className="truncate block">{buildInfo.source}</code>
            </div>
          </li>
        )}
        {buildInfo.sourceRevision && (
          <li>
            <div className="flex items-center">
              <p className="m-0 flex-shrink-0">Source revision: </p>
              <code className="truncate block">{buildInfo.sourceRevision}</code>
            </div>
          </li>
        )}
      </ul>
    </div>
  )
}

export type ProvenanceCardProps = {
  provenance: ImageProvenance
}

export function ProvenanceCard({
  provenance,
}: ProvenanceCardProps): JSX.Element {
  return (
    <Card
      persistenceKey="provenance"
      tabs={[
        {
          icon: <FluentDocumentRibbon16Regular />,
          label: 'Provenance',
          content: (
            <div className="markdown-body">
              <p>
                Some images include <i>provenance attestations</i> - means of
                asserting facts about an image's build process. These details
                are helpful for users to understand where an image is from and
                how it was built.
              </p>
              <p>
                More information can be found here:{' '}
                <a href="https://docs.docker.com/build/metadata/attestations/slsa-provenance/">
                  https://docs.docker.com/build/metadata/attestations/slsa-provenance/
                </a>
                .
              </p>
            </div>
          ),
        },
        ...provenance.buildInfo.map((buildInfo) => ({
          label:
            provenance.buildInfo.length > 0 && buildInfo.architecture
              ? `${[
                  buildInfo.operatingSystem,
                  buildInfo.architecture,
                  buildInfo.architectureVariant,
                ]
                  .filter((x) => x !== undefined)
                  .join('/')}`
              : 'Build info',
          content: <BuildInfo buildInfo={buildInfo} />,
        })),
      ]}
    />
  )
}
