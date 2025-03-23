import type { JSX } from 'react'
import type { ImageProvenance, ProvenanceBuildInfo } from '../../api'
import { FluentDocumentRibbon16Regular } from '../../components/icons/fluent-document-ribbon-16-regular'
import { SimpleIconsDocker } from '../../components/icons/simple-icons-docker'
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
          Image manifest:{' '}
          <code className="truncate">{buildInfo.imageDigest}</code>
        </li>
        <li>
          Build started: <code>{buildInfo.buildStartedOn}</code>
        </li>
        <li>
          Build finished: <code>{buildInfo.buildFinishedOn}</code>
        </li>
        <li>
          Architecture:{' '}
          <code>
            {[
              buildInfo.operatingSystem,
              buildInfo.architecture,
              buildInfo.architectureVariant,
            ]
              .filter((x) => x !== undefined)
              .join('/')}
          </code>
        </li>
        {buildInfo.source && (
          <li>
            {[buildInfo.source, buildInfo.sourceRevision]
              .filter((x) => x !== undefined)
              .join('@')}
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
          label: 'Attestations',
          content: (
            <div className="markdown-body">
              <p>
                Some images include <i>attestations</i> - means of asserting
                facts about an image's build process or contents. These details
                are helpful for users to understand what an image contains and
                where it's from, as well as for automated systems to identify
                vulnerabilities contained in an image.
              </p>
              <p>There are two main forms of attestations:</p>
              <ul>
                <li>
                  provenance attestations - details about an image's build
                  process. May include Dockerfiles, build timestamps and version
                  control metadata.
                </li>
                <li>
                  Software Bill of Materials (SBOM) attestations - details about
                  an image's software contents
                </li>
              </ul>
              <p>
                More information can be found here:{' '}
                <a href="https://docs.docker.com/build/metadata/attestations/">
                  https://docs.docker.com/build/metadata/attestations/
                </a>
                .
              </p>
            </div>
          ),
        },
        ...provenance.buildInfo.map((buildInfo) => ({
          icon: <SimpleIconsDocker className="w-[16px] h-[16px]" />,
          label:
            provenance.buildInfo.length > 0 && buildInfo.architecture
              ? `Provenance (${buildInfo.architecture})`
              : 'Provenance',
          content: <BuildInfo buildInfo={buildInfo} />,
        })),
      ]}
    />
  )
}
