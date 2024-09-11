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

import { Badge } from '../components/Badge'
import CustomNode from '../components/Node'
import { FluentChevronDown20Regular } from '../components/icons/fluent-chevron-down-20-regular'
import { FluentChevronUp20Regular } from '../components/icons/fluent-chevron-up-20-regular'

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

  const [nodes, setNodes, onNodesChange] = useNodesState(initNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initEdges)

  const tags: Tag[] = [
    { label: 'k8s', color: 'bg-blue-100' },
    { label: 'docker', color: 'bg-red-100' },
    { label: 'Pod', color: 'bg-orange-100' },
    { label: 'Job', color: 'bg-blue-100' },
    { label: 'ghcr', color: 'bg-blue-100' },
    { label: 'up-to-date', color: 'bg-green-100' },
    { label: 'outdated', color: 'bg-red-100' },
  ]

  const onConnect = useCallback(
    (connection: Connection) => setEdges((eds) => addEdge(connection, eds)),
    []
  )

  return (
    <div className="flex flex-col items-center w-full py-[40px] px-[20px]">
      <h1 className="text-2xl font-medium">{imageName}</h1>
      <div className="flex items-center">
        <FluentChevronDown20Regular />
        <p className="font-medium">{imageVersion}</p>
        <p className="font-medium ml-4">{imageVersion}</p>
        <FluentChevronUp20Regular />
      </div>
      <div className="flex mt-2 items-center">
        {tags.map((x) => (
          <Badge label={x.label} color={x.color} />
        ))}
      </div>
      <div className="rounded-lg bg-white px-4 py-2 shadow w-[780px] h-[480px] mt-6">
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
    </div>
  )
}
