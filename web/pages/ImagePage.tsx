import React, { type JSX, useState } from 'react'
import { Link, Navigate, useSearchParams } from 'react-router-dom'
import { Badge } from '../components/Badge'
import { DemoWarning } from '../components/DemoWarning'
import { HTML } from '../components/HTML'
import { ImageLogo } from '../components/ImageLogo'
import { InfoTooltip } from '../components/InfoTooltip'
import { FluentBookOpen16Regular } from '../components/icons/fluent-book-open-16-regular'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'
import { FluentWarning16Filled } from '../components/icons/fluent-warning-16-filled'
import { Markdown } from '../components/Markdown'
import { Toast } from '../components/Toast'
import { WordBreak } from '../components/WordBreak'
import { type Event, useEvents } from '../EventProvider'
import {
  useImage,
  useImageDescription,
  useImageGraph,
  useImageProvenance,
  useImageReleaseNotes,
  useImageSBOM,
  useImageScorecard,
  useImageVulnerabilities,
  useLatestWorkflowRun,
  useTags,
} from '../lib/api/ApiProvider'
import { formattedVersion, fullVersion, name } from '../oci'
import { compareTags } from '../tags'
import { formatRelativeTimeTo } from '../time'
import { Card } from './image-page/Card'
import { ImageSkeleton } from './image-page/ImageSkeleton'

const GraphCard = React.lazy(() =>
  import('./image-page/GraphCard').then((module) => ({
    default: module.GraphCard,
  }))
)

const LinksCard = React.lazy(() =>
  import('./image-page/LinksCard').then((module) => ({
    default: module.LinksCard,
  }))
)

const ProvenanceCard = React.lazy(() =>
  import('./image-page/ProvenanceCard').then((module) => ({
    default: module.ProvenanceCard,
  }))
)

const SBOMCard = React.lazy(() =>
  import('./image-page/SBOMCard').then((module) => ({
    default: module.SBOMCard,
  }))
)

const ScorecardCard = React.lazy(() =>
  import('./image-page/ScorecardCard').then((module) => ({
    default: module.ScorecardCard,
  }))
)

const SettingsCard = React.lazy(() =>
  import('./image-page/SettingsCard').then((module) => ({
    default: module.SettingsCard,
  }))
)

const VulnerabilitiesCard = React.lazy(() =>
  import('./image-page/VulnerabilitiesCard').then((module) => ({
    default: module.VulnerabilitiesCard,
  }))
)

const WorkflowCard = React.lazy(() =>
  import('./image-page/WorkflowCard').then((module) => ({
    default: module.WorkflowCard,
  }))
)

export function ImagePage(): JSX.Element {
  const [params, _] = useSearchParams()

  const reference = params.get('reference')!

  const [tags, updateTags] = useTags()
  const [image, updateImage] = useImage(reference)
  const [description, updateDescription] = useImageDescription(reference)
  const [releaseNotes, updateReleaseNotes] = useImageReleaseNotes(reference)
  const [graph, updateGraph] = useImageGraph(reference)
  const [scorecard, updateScorecard] = useImageScorecard(reference)
  const [provenance, updateProvenance] = useImageProvenance(reference)
  const [sbom, updateSBOM] = useImageSBOM(reference)
  const [vulnerabilities, updateVulnerabilities] =
    useImageVulnerabilities(reference)
  const [workflowRun, updateWorkflowRun] = useLatestWorkflowRun(reference)

  const [isUpdateAvailable, setIsUpdateAvailable] = useState(false)

  useEvents(
    (e: Event) => {
      // All data but the workflow runs are covered by the image updated event,
      // include the image processed event to cover workflow runs changing
      if (
        e.reference === reference &&
        (e.type === 'imageUpdated' || e.type === 'imageProcessed')
      ) {
        setIsUpdateAvailable(true)
      }
    },
    [reference]
  )

  if (
    tags.status !== 'resolved' ||
    image.status !== 'resolved' ||
    description.status !== 'resolved' ||
    releaseNotes.status !== 'resolved' ||
    graph.status !== 'resolved' ||
    scorecard.status !== 'resolved' ||
    provenance.status !== 'resolved' ||
    sbom.status !== 'resolved' ||
    vulnerabilities.status !== 'resolved' ||
    workflowRun.status !== 'resolved'
  ) {
    return <ImageSkeleton />
  }

  // Redirect if image was not found
  if (!image.value) {
    return <Navigate to="/" replace />
  }

  const imageTags = tags.value
    .filter((x) => image.value?.tags.includes(x.name))
    .toSorted((a, b) => compareTags(a.name, b.name))

  const externalReleasesUrl = image.value?.links.find(
    (x) => x.type === 'github-releases'
  )?.url

  return (
    <>
      <div className="fixed bottom-[env(safe-area-inset-bottom))] flex justify-center w-full sm:w-auto sm:right-0 p-4 z-100">
        {isUpdateAvailable && (
          <Toast
            title="New data available"
            body="The image is updated. Update to show the latest data."
            secondaryAction="Dismiss"
            onSecondaryAction={() => setIsUpdateAvailable(false)}
            primaryAction="Update"
            onPrimaryAction={() => {
              setIsUpdateAvailable(false)
              updateTags()
              updateImage()
              updateDescription()
              updateReleaseNotes()
              updateGraph()
              updateScorecard()
              updateProvenance()
              updateSBOM()
              updateVulnerabilities()
              updateWorkflowRun()
            }}
          />
        )}
      </div>
      <div className="flex flex-col items-center w-full pt-2 pb-10 px-2">
        <DemoWarning />
        {/* Header */}
        <ImageLogo
          className="w-[90px] h-[90px] mt-4"
          reference={image.value.reference}
        />
        {/* Image name */}
        <h1 className="text-2xl font-medium">
          <WordBreak delimiter="/">{name(image.value.reference)}</WordBreak>
        </h1>
        {/* Image version */}
        {/* Digests are formatted like <algo>:<digest>, such as sha256:<digest>. Show a maximum of 5 hex digits before truncating with ellipsis (hence 15ch) */}
        <div className="flex items-center">
          {!image.value.latestReference ? (
            <>
              <p
                className="font-medium max-w-[15ch] truncate"
                title={fullVersion(image.value.reference)}
              >
                {formattedVersion(
                  image.value.reference,
                  image.value.annotations
                )}{' '}
              </p>
              <InfoTooltip
                icon={<FluentWarning16Filled className="text-yellow-600" />}
              >
                The latest version cannot be identified. This could be due to
                the image not being available, the registry not being supported,
                missing authentication or a temporary issue.
              </InfoTooltip>
            </>
          ) : image.value.reference === image.value.latestReference ? (
            <p
              className="font-medium max-w-[15ch] truncate"
              title={fullVersion(image.value.reference)}
            >
              {formattedVersion(image.value.reference, image.value.annotations)}
            </p>
          ) : (
            <>
              <FluentChevronDown20Regular className="text-red-600" />
              <p
                className="font-medium text-red-600 max-w-[15ch] truncate"
                title={fullVersion(image.value.reference)}
              >
                {formattedVersion(
                  image.value.reference,
                  image.value.annotations
                )}
              </p>
              <p
                className="font-medium ml-4 text-green-600 max-w-[15ch] truncate"
                title={fullVersion(image.value.latestReference)}
              >
                {image.value.latestReference
                  ? formattedVersion(
                      image.value.latestReference,
                      image.value.latestAnnotations
                    )
                  : 'unknown'}
              </p>
              <FluentChevronUp20Regular className="text-green-600" />
            </>
          )}
        </div>
        {/* Image release dates, if newer and available */}
        {image.value.latestCreated && (
          <p>
            Last updated{' '}
            <span title={new Date(image.value.latestCreated).toLocaleString()}>
              {formatRelativeTimeTo(new Date(image.value.latestCreated))}
            </span>
          </p>
        )}
        {/* Image description, if available */}
        {image.value.description && (
          <p className="mt-2">{image.value.description}</p>
        )}
        {/* Image tags */}
        <div className="flex mt-2 items-center gap-1 flex-wrap justify-center">
          {imageTags.map((x) => (
            <Link
              key={x.name}
              to={`/?tag=${encodeURIComponent(x.name)}`}
              className="group/link"
              tabIndex={0}
            >
              <Badge
                key={x.name}
                label={x.name}
                color={x.color}
                className="hover:opacity-90 group-focus/link:opacity-90"
              />
            </Link>
          ))}
        </div>

        <main className="min-w-[200px] max-w-[800px] w-full box-border space-y-6 mt-6">
          {/* Cupdate settings */}
          {image.value?.reference.startsWith(
            'ghcr.io/alexgustafsson/cupdate'
          ) && <SettingsCard />}

          {/* Scorecard report */}
          {scorecard.value && <ScorecardCard scorecard={scorecard.value} />}

          {/* Vulnerability report */}
          {vulnerabilities.value && vulnerabilities.value.length > 0 && (
            <VulnerabilitiesCard vulnerabilities={vulnerabilities.value} />
          )}

          {/* Release notes */}
          {releaseNotes.value?.html && (
            <Card
              persistenceKey="release-notes"
              tabs={[
                {
                  icon: <FluentBookOpen16Regular />,
                  label: 'Release notes',
                  action: externalReleasesUrl
                    ? {
                        type: 'external-link',
                        href: externalReleasesUrl,
                        title: `External releases page for ${name(image.value.reference)}`,
                      }
                    : undefined,
                  content: (
                    <div className="markdown-body">
                      {releaseNotes.value?.title && (
                        <h1>{releaseNotes.value?.title}</h1>
                      )}
                      <HTML>{releaseNotes.value?.html}</HTML>
                    </div>
                  ),
                },
              ]}
            />
          )}

          {/* Description */}
          {(description.value?.html || description.value?.markdown) && (
            <Card
              persistenceKey="description"
              tabs={[
                {
                  icon: <FluentBookOpen16Regular />,
                  label: 'Description',
                  content: (
                    <div className="markdown-body">
                      {description.value.html && (
                        <HTML>{description.value.html}</HTML>
                      )}
                      {description.value.markdown && (
                        <Markdown>{description.value.markdown}</Markdown>
                      )}
                    </div>
                  ),
                },
              ]}
            />
          )}

          {/* Links */}
          {image.value && image.value.links.length > 0 && (
            <LinksCard links={image.value.links} />
          )}

          {/* Provenance report */}
          {provenance.value && <ProvenanceCard provenance={provenance.value} />}

          {/* SBOM */}
          {sbom.value && <SBOMCard sbom={sbom.value} />}

          {/* Graph */}
          {graph.value && <GraphCard graph={graph.value} />}

          {/* Workflow summary */}
          <WorkflowCard
            workflowRun={workflowRun.value}
            reference={image.value.reference}
            lastModified={image.value.lastModified}
          />
        </main>
      </div>
    </>
  )
}
