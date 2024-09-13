import {
  Connection,
  Controls,
  Edge,
  MiniMap,
  Node,
  NodeTypes,
  OnEdgesChange,
  OnNodesChange,
  ReactFlow,
  addEdge,
  useEdgesState,
  useNodesState,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import { ReactNode, useCallback, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'

import {
  Graph,
  Result,
  useImage,
  useImageDescription,
  useImageGraph,
  useImageReleaseNotes,
  useTags,
} from '../api'
import { Badge } from '../components/Badge'
import { HTML } from '../components/HTML'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'
import { Quay } from '../components/icons/quay'
import { SimpleIconsDocker } from '../components/icons/simple-icons-docker'
import { SimpleIconsGit } from '../components/icons/simple-icons-git'
import { SimpleIconsGithub } from '../components/icons/simple-icons-github'
import { SimpleIconsGitlab } from '../components/icons/simple-icons-gitlab'
import { SimpleIconsOci } from '../components/icons/simple-icons-oci'
import { nodeTypes, useNodesAndEdges } from '../graph'

const titles: Record<string, string | undefined> = {
  github: "Project's page on GitHub",
  gitlab: "Project's page on GitLab",
  docker: "Project's page on Docker Hub",
  quay: "Project's page on Quay.io",
  git: "Project's git page",
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
  }

  return (
    <a title={title} href={url} target="_blank">
      {icon}
    </a>
  )
}

export function ImagePage(): JSX.Element {
  const [params, _] = useSearchParams()

  const imageName = params.get('name')
  const imageVersion = params.get('version')

  const image = useImage()
  const description = useImageDescription()
  const releaseNotes = useImageReleaseNotes()
  const graph = useImageGraph()
  const tags = useTags()

  const [[nodes, onNodesChange], [edges, _onEdgesChange]] =
    useNodesAndEdges(graph)

  if (
    image.status !== 'resolved' ||
    description.status !== 'resolved' ||
    releaseNotes.status !== 'resolved' ||
    tags.status !== 'resolved'
  ) {
    return <p>Loading</p>
  }

  const imageTags = tags.value.filter((x) => image.value.tags.includes(x.name))

  return (
    <div className="flex flex-col items-center w-full py-[40px] px-[20px]">
      {/* Header */}
      {image.value.image && (
        <img className="w-16 rounded" src={image.value.image} />
      )}
      <h1 className="text-2xl font-medium">{image.value.name}</h1>
      <div className="flex items-center">
        <FluentChevronDown20Regular className="text-red-500" />
        <p className="font-medium text-red-500">{image.value.currentVersion}</p>
        <p className="font-medium ml-4 text-green-500">
          {image.value.latestVersion}
        </p>
        <FluentChevronUp20Regular className="text-green-500" />
      </div>
      <div className="flex mt-2 items-center">
        {imageTags.map((x) => (
          <Badge key={x.name} label={x.name} color={x.color} />
        ))}
      </div>

      {/* Release notes */}
      <div className="flex mt-2 space-x-4 items-center">
        {image.value.links.map((link) => (
          <ImageLink
            key={`${link.type}:${link.url}`}
            type={link.type}
            url={link.url}
          />
        ))}
      </div>

      <main className="min-w-[200px] max-w-[980px] box-border space-y-6 mt-6">
        <div className="rounded-lg bg-white px-4 py-6 shadow">
          <div className="markdown-body">
            <HTML>{description.value?.html}</HTML>
          </div>
        </div>

        <div className="rounded-lg bg-white px-4 py-6 shadow">
          <div className="markdown-body">
            <h1>{releaseNotes.value?.title}</h1>
            <HTML>{releaseNotes.value?.html}</HTML>
          </div>
        </div>

        {/* Graph */}
        <div className="rounded-lg bg-white px-4 py-2 shadow h-[480px]">
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
          >
            <Controls />
          </ReactFlow>
        </div>
      </main>
    </div>
  )
}
