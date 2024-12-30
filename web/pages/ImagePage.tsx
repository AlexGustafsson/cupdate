import type { JSX, ReactNode } from 'react'
import { Navigate, useSearchParams } from 'react-router-dom'

import {
  type Graph,
  useImage,
  useImageDescription,
  useImageGraph,
  useImageReleaseNotes,
  useScheduleScan,
  useTags,
} from '../api'
import { Badge } from '../components/Badge'
import { Graph as GraphRenderer } from '../components/Graph'
import { HTML } from '../components/HTML'
import { ImageLogo } from '../components/ImageLogo'
import { InfoTooltip } from '../components/InfoTooltip'
import { Markdown } from '../components/Markdown'
import { FluentArrowSync16Regular } from '../components/icons/fluent-arrow-sync-16-regular'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'
import { FluentLink24Filled } from '../components/icons/fluent-link-24-filled'
import { FluentShieldError24Filled } from '../components/icons/fluent-shield-error-24-filled'
import { FluentWarning16Filled } from '../components/icons/fluent-warning-16-filled'
import { Quay } from '../components/icons/quay'
import { SimpleIconsDocker } from '../components/icons/simple-icons-docker'
import { SimpleIconsGit } from '../components/icons/simple-icons-git'
import { SimpleIconsGithub } from '../components/icons/simple-icons-github'
import { SimpleIconsGitlab } from '../components/icons/simple-icons-gitlab'
import { SimpleIconsOci } from '../components/icons/simple-icons-oci'
import { useNodesAndEdges } from '../graph'
import { name, version } from '../oci'
import { formatRelativeTimeTo } from '../time'

const titles: Record<string, string | undefined> = {
  github: "Project's page on GitHub",
  ghcr: "Project's package page on GHCR",
  gitlab: "Project's page on GitLab",
  docker: "Project's page on Docker Hub",
  quay: "Project's page on Quay.io",
  git: "Project's git page",
  'oci-registry': "Project's OCI registry",
}

function unique<T>(previousValue: T[], currentValue: T): T[] {
  if (previousValue.includes(currentValue)) {
    return previousValue
  }

  previousValue.push(currentValue)
  return previousValue
}

function flattened<T>(previousValue: T[], currentValue: T[]): T[] {
  previousValue.push(...currentValue)
  return previousValue
}

export function ImageLink({
  type,
  url,
}: {
  type: string
  url: string
}): JSX.Element {
  const title = titles[type]
  let icon: ReactNode
  switch (type) {
    case 'github':
    case 'ghcr':
      icon = (
        <SimpleIconsGithub className="text-black dark:dark:text-[#dddddd]" />
      )
      break
    case 'gitlab':
      icon = <SimpleIconsGitlab className="text-orange-500" />
      break
    case 'docker':
      icon = <SimpleIconsDocker className="text-blue-500" />
      break
    case 'quay':
      icon = <Quay className="text-blue-700" />
      break
    case 'git':
      icon = <SimpleIconsGit className="text-orange-500" />
      break
    case 'svc':
      if (url.includes('github.com')) {
        return <ImageLink type="github" url={url} />
      } else if (url.includes('gitlab')) {
        return <ImageLink type="gitlab" url={url} />
      } else {
        return <ImageLink type="git" url={url} />
      }
    default:
      icon = <FluentLink24Filled />
  }

  return (
    <a title={title} href={url} target="_blank" rel="noreferrer">
      {icon}
    </a>
  )
}

type GraphCardProps = {
  graph: Graph
}

export function GraphCard({ graph }: GraphCardProps): JSX.Element {
  const [nodes, edges, bounds] = useNodesAndEdges(graph)

  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-2 shadow h-[480px]">
      <GraphRenderer edges={edges} nodes={nodes} bounds={bounds} />
    </div>
  )
}

export function ImagePage(): JSX.Element {
  const [params, _] = useSearchParams()

  const reference = params.get('reference')!

  const tags = useTags()
  const image = useImage(reference)
  const description = useImageDescription(reference)
  const releaseNotes = useImageReleaseNotes(reference)
  const graph = useImageGraph(reference)

  const scheduleScan = useScheduleScan()

  if (
    tags.status === 'idle' ||
    image.status === 'idle' ||
    description.status === 'idle' ||
    releaseNotes.status === 'idle' ||
    graph.status === 'idle'
  ) {
    return <></>
  }

  if (
    tags.status === 'rejected' ||
    image.status === 'rejected' ||
    description.status === 'rejected' ||
    releaseNotes.status === 'rejected' ||
    graph.status === 'rejected'
  ) {
    return <p>Error</p>
  }

  // Redirect if image was not found
  if (!image.value) {
    return <Navigate to="/" replace />
  }

  const imageTags = tags.value.filter((x) => image.value?.tags.includes(x.name))

  return (
    <div className="flex flex-col items-center w-full pt-6 pb-10 px-2">
      {/* Header */}
      <ImageLogo src={image.value.image} width={90} height={90} />
      {/* Image name */}
      <h1 className="text-2xl font-medium">
        {name(image.value.reference)}
        {image.value.vulnerabilities.length > 0 && (
          <InfoTooltip
            icon={<FluentShieldError24Filled className="text-red-600" />}
          >
            {image.value.vulnerabilities.length} vulnerabilities reported.
          </InfoTooltip>
        )}
      </h1>
      {/* Image version */}
      <div className="flex items-center">
        {!image.value.latestReference ? (
          <p className="font-medium">
            {version(image.value.reference)}{' '}
            <InfoTooltip
              icon={<FluentWarning16Filled className="text-yellow-600" />}
            >
              The latest version cannot be identified. This could be due to the
              image not being available, the registry not being supported,
              missing authentication or a temporary issue.
            </InfoTooltip>
          </p>
        ) : image.value.reference === image.value.latestReference ? (
          <p className="font-medium">{version(image.value.reference)}</p>
        ) : (
          <>
            <FluentChevronDown20Regular className="text-red-600" />
            <p className="font-medium text-red-600">
              {version(image.value.reference)}
            </p>
            <p className="font-medium ml-4 text-green-600">
              {image.value.latestReference
                ? version(image.value.latestReference)
                : 'unknown'}
            </p>
            <FluentChevronUp20Regular className="text-green-600" />
          </>
        )}
      </div>
      {/* Image description, if available */}
      {image.value.description && (
        <p className="mt-2">{image.value.description}</p>
      )}
      {/* Image tags */}
      <div className="flex mt-4 items-center gap-x-2">
        {imageTags.map((x) => (
          <Badge key={x.name} label={x.name} color={x.color} />
        ))}
      </div>
      {/* Links */}
      <div className="flex mt-2 space-x-4 items-center">
        {image.value.links.map((link) => (
          <ImageLink
            key={`${link.type}:${link.url}`}
            type={link.type}
            url={link.url}
          />
        ))}
      </div>

      <main className="min-w-[200px] max-w-[800px] w-full box-border space-y-6 mt-6">
        {/* Cupdate settings */}
        {image.value?.reference === 'ghcr.io/alexgustafsson/cupdate' && (
          <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow">
            <p>
              Cupdate version:{' '}
              {import.meta.env.VITE_CUPDATE_VERSION || 'development build'}.
            </p>
          </div>
        )}

        {/* Vulnerability report */}
        {image.value.vulnerabilities.length > 0 && (
          <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow">
            <div className="markdown-body">
              <h1>Vulnerabilities</h1>
              <ul>
                <li>
                  Critical:{' '}
                  {
                    image.value.vulnerabilities.filter(
                      (x) => x.severity === 'critical'
                    ).length
                  }
                </li>
                <li>
                  High:{' '}
                  {
                    image.value.vulnerabilities.filter(
                      (x) => x.severity === 'high'
                    ).length
                  }
                </li>
                <li>
                  Medium:{' '}
                  {
                    image.value.vulnerabilities.filter(
                      (x) => x.severity === 'medium'
                    ).length
                  }
                </li>
                <li>
                  Low:{' '}
                  {
                    image.value.vulnerabilities.filter(
                      (x) => x.severity === 'low'
                    ).length
                  }
                </li>
                <li>
                  Unspecified:{' '}
                  {
                    image.value.vulnerabilities.filter(
                      (x) => x.severity === 'unspecified'
                    ).length
                  }
                </li>
              </ul>

              <h2>Authorities</h2>
              <ul>
                {image.value.vulnerabilities
                  .map((x) => x.authority)
                  .reduce(unique<string>, [])
                  .map((x) => (
                    <li key={x}>{x}</li>
                  ))}
              </ul>

              <h2>Links</h2>
              <ul>
                {image.value.vulnerabilities
                  .map((x) => x.links)
                  .reduce(flattened<string>, [])
                  .reduce(unique<string>, [])
                  .map((x) => (
                    <li key={x}>
                      <a href={x}>{x}</a>
                    </li>
                  ))}
              </ul>
            </div>
          </div>
        )}

        {/* Release notes */}
        {releaseNotes.value?.html && (
          <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow">
            <div className="markdown-body">
              <h1>{releaseNotes.value?.title}</h1>
              <HTML>{releaseNotes.value?.html}</HTML>
            </div>
          </div>
        )}

        {/* Description */}
        {(description.value?.html || description.value?.markdown) && (
          <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow">
            <div className="markdown-body">
              {description.value.html && <HTML>{description.value.html}</HTML>}
              {description.value.markdown && (
                <Markdown>{description.value.markdown}</Markdown>
              )}
            </div>
          </div>
        )}

        {/* Graph */}
        {graph.value && <GraphCard graph={graph.value} />}

        <div className="flex justify-center gap-x-2 items-center">
          <p>
            Last updated{' '}
            <span title={new Date(image.value.lastModified).toLocaleString()}>
              {formatRelativeTimeTo(new Date(image.value.lastModified))}
            </span>
          </p>
          <button
            type="button"
            title="Schedule update"
            onClick={() =>
              image.value ? scheduleScan(image.value.reference) : undefined
            }
          >
            <FluentArrowSync16Regular />
          </button>
        </div>
      </main>
    </div>
  )
}
