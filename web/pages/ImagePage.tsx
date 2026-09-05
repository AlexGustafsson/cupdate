import React, { type JSX, useCallback, useState } from 'react'
import { Link, Navigate, useSearchParams } from 'react-router'
import { Badge } from '../components/Badge'
import { Card } from '../components/Card'
import { DemoWarning } from '../components/DemoWarning'
import { HTML } from '../components/HTML'
import { ImageLogo } from '../components/ImageLogo'
import { InfoTooltip } from '../components/InfoTooltip'
import { FluentBookOpen16Regular } from '../components/icons/fluent-book-open-16-regular'
import { FluentBoxSearch16Regular } from '../components/icons/fluent-box-search-16-regular'
import { FluentBranch16Regular } from '../components/icons/fluent-branch-16-regular'
import { FluentBug16Regular } from '../components/icons/fluent-bug-16-regular'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'
import { FluentDocumentRibbon16Regular } from '../components/icons/fluent-document-ribbon-16-regular'
import { FluentFlow16Regular } from '../components/icons/fluent-flow-16-regular'
import { FluentLink16Regular } from '../components/icons/fluent-link-16-regular'
import { FluentWarning16Filled } from '../components/icons/fluent-warning-16-filled'
import { Markdown } from '../components/Markdown'
import { Toast } from '../components/Toast'
import { WordBreak } from '../components/WordBreak'
import { type Event, useEvents } from '../EventProvider'
import { useCollapseState } from '../hooks/useCollapseState'
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
import { ImageSkeleton } from './image-page/ImageSkeleton'
import { TableOfContents } from './image-page/TableOfContents'

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
        (e.type === 'imageUpdated' || e.type === 'imageProcessed') &&
        e.reference === reference
      ) {
        setIsUpdateAvailable(true)
      }
    },
    [reference]
  )

  const scrollIntoView = useCallback((id: string) => {
    const element = document.getElementById(id)
    element?.scrollIntoView({ block: 'start', behavior: 'smooth' })
  }, [])

  const [collapsed, toggleCollapsed] = useCollapseState()

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
      <div className="w-full flex flex-col items-center pt-2 pb-10 px-2">
        {/* First content / glance */}
        <div
          id="overview"
          className="flex flex-col w-full items-center scroll-mt-[64px]"
        >
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
                  icon={<FluentWarning16Filled className="text-warning" />}
                  className="ml-1"
                >
                  The latest version hasn't been identified. This could be due
                  to the image not being processed yet, not being available, the
                  registry not being supported, missing authentication or a
                  temporary issue.
                </InfoTooltip>
              </>
            ) : image.value.reference === image.value.latestReference ? (
              <p
                className="font-medium max-w-[15ch] truncate"
                title={fullVersion(image.value.reference)}
              >
                {formattedVersion(
                  image.value.reference,
                  image.value.annotations
                )}
              </p>
            ) : (
              <>
                <FluentChevronDown20Regular className="text-negative" />
                <p
                  className="font-medium text-negative max-w-[15ch] truncate"
                  title={fullVersion(image.value.reference)}
                >
                  {formattedVersion(
                    image.value.reference,
                    image.value.annotations
                  )}
                </p>
                <p
                  className="font-medium ml-4 text-positive max-w-[15ch] truncate"
                  title={fullVersion(image.value.latestReference)}
                >
                  {image.value.latestReference
                    ? formattedVersion(
                        image.value.latestReference,
                        image.value.latestAnnotations
                      )
                    : 'unknown'}
                </p>
                <FluentChevronUp20Regular className="text-positive" />
              </>
            )}
          </div>
          {/* Image release dates, if newer and available */}
          {image.value.latestCreated && (
            <p>
              Last updated{' '}
              <span
                title={new Date(image.value.latestCreated).toLocaleString()}
              >
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
        </div>
        <div className="w-full mt-6 grid grid-cols-1 lg:grid-cols-[1fr_max(200px,min(100%,800px))_1fr]">
          {/* Left column */}
          <div className="pr-6 hidden lg:block"></div>
          {/* Big center column */}
          <div className="flex flex-col w-full space-y-6">
            {/* Scorecard report */}
            {scorecard.value && (
              <ScorecardCard
                id="scorecard"
                scorecard={scorecard.value}
                collapsed={collapsed.has('scorecard')}
                onToggleCollapsed={() => toggleCollapsed('scorecard')}
              />
            )}

            {/* Vulnerability report */}
            {vulnerabilities.value && vulnerabilities.value.length > 0 && (
              <VulnerabilitiesCard
                id="vulnerabilities"
                vulnerabilities={vulnerabilities.value}
                collapsed={collapsed.has('vulnerabilities')}
                onToggleCollapsed={() => toggleCollapsed('vulnerabilities')}
              />
            )}

            {/* Release notes */}
            {releaseNotes.value?.html && (
              <Card
                id="release-notes"
                collapsed={collapsed.has('release-notes')}
                onToggleCollapsed={() => toggleCollapsed('release-notes')}
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
                id="description"
                collapsed={collapsed.has('description')}
                onToggleCollapsed={() => toggleCollapsed('description')}
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
              <LinksCard
                id="links"
                links={image.value.links}
                collapsed={collapsed.has('links')}
                onToggleCollapsed={() => toggleCollapsed('links')}
              />
            )}

            {/* Provenance report */}
            {provenance.value && (
              <ProvenanceCard
                id="provenance"
                provenance={provenance.value}
                collapsed={collapsed.has('provenance')}
                onToggleCollapsed={() => toggleCollapsed('provenance')}
              />
            )}

            {/* SBOM */}
            {sbom.value && (
              <SBOMCard
                id="sbom"
                sbom={sbom.value}
                collapsed={collapsed.has('sbom')}
                onToggleCollapsed={() => toggleCollapsed('sbom')}
              />
            )}

            {/* Graph */}
            {graph.value && (
              <GraphCard
                id="graph"
                graph={graph.value}
                collapsed={collapsed.has('graph')}
                onToggleCollapsed={() => toggleCollapsed('graph')}
              />
            )}

            {/* Workflow summary */}
            <WorkflowCard
              id="workflow"
              workflowRun={workflowRun.value}
              reference={image.value.reference}
              lastModified={image.value.lastModified}
              collapsed={collapsed.has('workflow')}
              onToggleCollapsed={() => toggleCollapsed('workflow')}
            />
          </div>
          {/* Right column */}
          <div className="pl-6 hidden lg:block">
            <div className="sticky top-[64px]">
              <TableOfContents
                onClick={scrollIntoView}
                onToggleCollapsed={toggleCollapsed}
                items={[
                  {
                    id: 'overview',
                    icon: <FluentBookOpen16Regular />,
                    label: 'Overview',
                  },
                  {
                    id: 'scorecard',
                    icon: <FluentBoxSearch16Regular />,
                    label: 'Risk score',
                    disabled: scorecard.value === null,
                    collapsed: collapsed.has('scorecard'),
                  },
                  {
                    id: 'vulnerabilities',
                    icon: <FluentBug16Regular />,
                    label: 'Vulnerablities',
                    disabled:
                      vulnerabilities.value === null ||
                      vulnerabilities.value.length === 0,
                    collapsed: collapsed.has('vulnerabilities'),
                  },
                  {
                    id: 'release-notes',
                    icon: <FluentBookOpen16Regular />,
                    label: 'Release notes',
                    disabled: releaseNotes.value === null,
                    collapsed: collapsed.has('release-notes'),
                  },
                  {
                    id: 'description',
                    icon: <FluentBookOpen16Regular />,
                    label: 'Description',
                    disabled: description.value === null,
                    collapsed: collapsed.has('description'),
                  },
                  {
                    id: 'links',
                    icon: <FluentLink16Regular />,
                    label: 'Links',
                    disabled: image.value === null,
                    collapsed: collapsed.has('links'),
                  },
                  {
                    id: 'provenance',
                    icon: <FluentDocumentRibbon16Regular />,
                    label: 'Provenance',
                    disabled: provenance.value === null,
                    collapsed: collapsed.has('provenance'),
                  },
                  {
                    id: 'sbom',
                    icon: <FluentBoxSearch16Regular />,
                    label: 'SBOM',
                    disabled: sbom.value === null,
                    collapsed: collapsed.has('sbom'),
                  },
                  {
                    id: 'graph',
                    icon: <FluentBranch16Regular />,
                    label: 'Graph',
                    disabled: graph.value === null,
                    collapsed: collapsed.has('graph'),
                  },
                  {
                    id: 'workflow',
                    icon: <FluentFlow16Regular />,
                    label: 'Workflow',
                    disabled: workflowRun.value === null,
                    collapsed: collapsed.has('workflow'),
                  },
                ]}
              />
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
