import type {
  Edge,
  Node,
  NodeTypes,
  OnEdgesChange,
  OnNodesChange,
} from '@xyflow/react'
import { getNodesBounds, useEdgesState, useNodesState } from '@xyflow/react'
import ELK from 'elkjs/lib/elk.bundled'
import { type ReactNode, useEffect } from 'react'

import type { Graph, GraphNode } from './api'
import { CustomGraphNode } from './components/CustomGraphNode'
import { SimpleIconsKubernetes } from './components/icons/simple-icons-kubernetes'
import { SimpleIconsOci } from './components/icons/simple-icons-oci'

export const nodeTypes: NodeTypes = {
  custom: CustomGraphNode,
}

export interface NodeType extends Node {
  data: {
    subtitle: string
    title: string
    label: ReactNode
  }
}

export interface EdgeType extends Edge {}

const titles: Record<string, Record<string, string | undefined> | undefined> = {
  oci: {
    image: 'Image',
  },
  kubernetes: {
    'core/v1/pod': 'Pod',
    'core/v1/namespace': 'Namespace',
    'core/v1/container': 'Container',
    'apps/v1/deployment': 'Deployment',
    'apps/v1/replicaset': 'Replica set',
    'batch/v1/job': 'Job',
    'batch/v1/cronjob': 'Cron job',
    'apps/v1/statefulset': 'Stateful set',
  },
}

function formatNode(id: string, node: GraphNode): NodeType {
  let label: ReactNode
  switch (node.domain) {
    case 'oci':
      label = <SimpleIconsOci className="text-blue-700" />
      break
    case 'kubernetes':
      label = <SimpleIconsKubernetes className="text-blue-400" />
      break
  }

  return {
    id,
    type: 'custom',
    data: {
      title: titles[node.domain]?.[node.type] || node.type,
      subtitle: node.name,
      label,
    },
    position: { x: 0, y: 0 },
  }
}

async function formatGraph(graph: Graph): Promise<[NodeType[], EdgeType[]]> {
  const nodes: NodeType[] = []
  const edges: EdgeType[] = []

  const elk = new ELK()

  const root = await elk.layout({
    id: 'root',
    layoutOptions: { 'elk.algorithm': 'mrtree', 'elk.spacing.nodeNode': '50' },
    children: Object.entries(graph.nodes).map(([id, _node]) => ({
      id,
      width: 250,
      height: 75,
    })),
    edges: Object.entries(graph.edges)
      .map(([from, adjacent]) =>
        Object.entries(adjacent)
          .filter(([_, isParent]) => isParent)
          .map(([to, _]) => ({
            id: `${from}->${to}`,
            sources: [to],
            targets: [from],
          }))
      )
      .flat(2),
  })

  for (const node of root.children || []) {
    const formatted = formatNode(node.id, graph.nodes[node.id])
    formatted.position.x = (node.x || 0) - (node.width || 0) / 2
    formatted.position.y = node.y || 0
    formatted.width = node.width
    formatted.height = node.height
    nodes.push(formatted)
  }

  const bounds = getNodesBounds(nodes)

  // The node for the image is (should be) on a "row" of its own. Always center
  // the image node on that row
  for (const node of nodes) {
    if (
      graph.nodes[node.id].domain === 'oci' &&
      graph.nodes[node.id].type === 'image'
    ) {
      // 42 is a magic number - an offset that makes centering look nice. I
      // think that there's something going on with the layout engine - that the
      // center of the graph isn't always the center of the image.
      // TODO: There are other good reasons to just write a simple viewer, this
      // is another reason why
      node.position.x = bounds.width / 2 - (node.width ?? 0) + 42
    }
  }

  // Pod nodes without names are typically templates, mark these as such
  for (const node of nodes) {
    if (
      graph.nodes[node.id].domain === 'kubernetes' &&
      graph.nodes[node.id].type === 'core/v1/pod' &&
      node.data.subtitle === ''
    ) {
      node.data.subtitle = '<template>'
    }
  }

  for (const edge of root.edges || []) {
    edges.push({
      id: edge.id,
      source: edge.sources[0],
      target: edge.targets[0],
    })
  }

  return [nodes, edges]
}

export function useNodesAndEdges(
  graph: Graph
): [
  [NodeType[], OnNodesChange<NodeType>],
  [EdgeType[], OnEdgesChange<EdgeType>],
] {
  const [nodes, setNodes, onNodesChange] = useNodesState<NodeType>([])
  const [edges, setEdges, onEdgesChange] = useEdgesState<EdgeType>([])

  useEffect(() => {
    formatGraph(graph).then(([nodes, edges]) => {
      setNodes(nodes)
      setEdges(edges)
    })
  }, [graph])

  return [
    [nodes, onNodesChange],
    [edges, onEdgesChange],
  ]
}
