import {
  Connection,
  Controls,
  MiniMap,
  NodeTypes,
  ReactFlow,
  addEdge,
  useEdgesState,
  useNodesState,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import { useCallback } from 'react'
import { useSearchParams } from 'react-router-dom'

import {
  useImage,
  useImageDescription,
  useImageReleaseNotes,
  useTags,
} from '../api'
import { Badge } from '../components/Badge'
import { HTML } from '../components/HTML'
import CustomNode from '../components/Node'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'
import { Quay } from '../components/icons/quay'
import { SimpleIconsDocker } from '../components/icons/simple-icons-docker'
import { SimpleIconsGit } from '../components/icons/simple-icons-git'
import { SimpleIconsGithub } from '../components/icons/simple-icons-github'
import { SimpleIconsGitlab } from '../components/icons/simple-icons-gitlab'

interface Tag {
  label: string
  color: string
}

const nodeTypes: NodeTypes = {
  custom: CustomNode,
}

const initNodes = [
  {
    id: '1',
    type: 'custom',
    data: { subtitle: 'default', title: 'Namespace', label: 'N' },
    position: { x: 0, y: 50 },
  },
  {
    id: '2',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Deployment', label: 'D' },
    position: { x: 0, y: 150 },
  },
  {
    id: '3',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Pod', label: 'P' },
    position: { x: 0, y: 250 },
  },
  {
    id: '4',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Container', label: 'C' },
    position: { x: 0, y: 350 },
  },

  {
    id: '8',
    type: 'custom',
    data: { subtitle: 'test', title: 'Namespace', label: 'N' },
    position: { x: 250, y: 150 },
  },
  {
    id: '6',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Pod', label: 'P' },
    position: { x: 250, y: 250 },
  },
  {
    id: '7',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Container', label: 'C' },
    position: { x: 250, y: 350 },
  },

  {
    id: '5',
    type: 'custom',
    data: { subtitle: 'home-assistant', title: 'Image', label: 'I' },
    position: { x: 0, y: 450 },
  },
]

const initEdges = [
  {
    id: 'e1',
    source: '1',
    target: '2',
  },
  {
    id: 'e2',
    source: '2',
    target: '3',
  },
  {
    id: 'e3',
    source: '3',
    target: '4',
  },
  {
    id: 'e4',
    source: '4',
    target: '5',
  },
  {
    id: 'e5',
    source: '6',
    target: '7',
  },
  {
    id: 'e6',
    source: '7',
    target: '5',
  },
  {
    id: 'e7',
    source: '8',
    target: '6',
  },
]

export function ImagePage(): JSX.Element {
  const [params, _] = useSearchParams()

  const imageName = params.get('name')
  const imageVersion = params.get('version')

  const image = useImage()
  const description = useImageDescription()
  const releaseNotes = useImageReleaseNotes()
  const tags = useTags()

  const [nodes, setNodes, onNodesChange] = useNodesState(initNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initEdges)

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
        <a
          title="Project's page on GitHub"
          href="https://github.com/home-assistant/core"
          target="_blank"
        >
          <SimpleIconsGithub className="text-black" />
        </a>
        <a
          title="Project's page on GitLab"
          href="https://gitlab.com/arm-research/smarter/smarter-device-manager"
          target="_blank"
        >
          <SimpleIconsGitlab className="text-orange-500" />
        </a>
        <a
          title="Project's page on Docker Hub"
          href="https://hub.docker.com/r/homeassistant/home-assistant"
          target="_blank"
        >
          <SimpleIconsDocker className="text-blue-500" />
        </a>
        <a
          title="Project's page on Quay"
          href="https://quay.io/repository/jetstack/cert-manager-webhook?tab=info"
          target="_blank"
        >
          <Quay className="text-blue-500" />
        </a>
        <a
          title="Project's source code"
          href="https://quay.io/repository/jetstack/cert-manager-webhook?tab=info"
          target="_blank"
        >
          <SimpleIconsGit className="text-orange-500" />
        </a>
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
