import { Controls, ReactFlow } from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import { ReactNode } from 'react'
import { Navigate, useSearchParams } from 'react-router-dom'

import {
  Graph,
  useImage,
  useImageDescription,
  useImageGraph,
  useImageReleaseNotes,
  useTags,
} from '../api'
import { Badge } from '../components/Badge'
import { HTML } from '../components/HTML'
import { Markdown } from '../components/Markdown'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'
import { FluentLink24Filled } from '../components/icons/fluent-link-24-filled'
import { Quay } from '../components/icons/quay'
import { SimpleIconsDocker } from '../components/icons/simple-icons-docker'
import { SimpleIconsGit } from '../components/icons/simple-icons-git'
import { SimpleIconsGithub } from '../components/icons/simple-icons-github'
import { SimpleIconsGitlab } from '../components/icons/simple-icons-gitlab'
import { SimpleIconsOci } from '../components/icons/simple-icons-oci'
import { nodeTypes, useNodesAndEdges } from '../graph'
import { name, version } from '../oci'

const titles: Record<string, string | undefined> = {
  github: "Project's page on GitHub",
  ghcr: "Project's package page on GHCR",
  gitlab: "Project's page on GitLab",
  docker: "Project's page on Docker Hub",
  quay: "Project's page on Quay.io",
  git: "Project's git page",
  'oci-registry': "Project's OCI registry",
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
      icon = <SimpleIconsGithub className="text-black" />
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
    default:
      icon = <FluentLink24Filled />
  }

  return (
    <a title={title} href={url} target="_blank">
      {icon}
    </a>
  )
}

type GraphCardProps = {
  graph: Graph
}

export function GraphCard({ graph }: GraphCardProps): JSX.Element {
  const [[nodes, onNodesChange], [edges, _onEdgesChange]] =
    useNodesAndEdges(graph)

  return (
    <div className="rounded-lg bg-white dark:bg-[#121212] px-4 py-2 shadow h-[480px]">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        nodeTypes={nodeTypes}
        fitView
        edgesFocusable={false}
        nodesDraggable={true}
        nodesConnectable={false}
        nodesFocusable={false}
        draggable={true}
        panOnDrag={true}
        elementsSelectable={false}
        minZoom={0.001}
      >
        <Controls />
      </ReactFlow>
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

  const imageTags = tags.value.filter((x) => image.value!.tags.includes(x.name))

  return (
    <div className="flex flex-col items-center w-full py-[40px] px-[20px]">
      {/* Header */}
      {image.value.image ? (
        <img
          className="w-16 h-16 rounded"
          src={image.value.image}
          no-referrer
        />
      ) : (
        <div className="w-16 h-16 rounded bg-blue-500 flex items-center justify-center">
          <SimpleIconsOci className="text-white" />
        </div>
      )}
      {/* Image name */}
      <h1 className="text-2xl font-medium">{name(image.value.reference)}</h1>
      {/* Image version */}
      <div className="flex items-center">
        {!image.value.latestReference ||
        image.value.reference === image.value.latestReference ? (
          <p className="font-medium">{version(image.value.reference)}</p>
        ) : (
          <>
            <FluentChevronDown20Regular className="text-red-500" />
            <p className="font-medium text-red-500">
              {version(image.value.reference)}
            </p>
            <p className="font-medium ml-4 text-green-500">
              {image.value.latestReference
                ? version(image.value.latestReference)
                : 'unknown'}
            </p>
            <FluentChevronUp20Regular className="text-green-500" />
          </>
        )}
      </div>
      {/* Image description, if available */}
      {image.value.description && <p>{image.value.description}</p>}
      {/* Image tags */}
      <div className="flex mt-2 items-center">
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

      <main className="min-w-[200px] max-w-[980px] w-full box-border space-y-6 mt-6">
        {/* Release notes */}
        {releaseNotes.value?.html && (
          <div className="rounded-lg bg-white dark:bg-[#121212] px-4 py-6 shadow">
            <div className="markdown-body">
              <h1>{releaseNotes.value?.title}</h1>
              <HTML>{releaseNotes.value?.html}</HTML>
            </div>
          </div>
        )}

        {/* Description */}
        {(description.value?.html || description.value?.markdown) && (
          <div className="rounded-lg bg-white dark:bg-[#121212] px-4 py-6 shadow">
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

        <div className="flex justify-center">
          <p>
            Last updated {new Date(image.value?.lastModified).toLocaleString()}
          </p>
        </div>
      </main>
    </div>
  )
}
