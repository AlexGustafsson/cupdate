import ELK from 'elkjs/lib/elk.bundled'
import { type ReactNode, useEffect, useState } from 'react'

import type { Graph, GraphNode } from './api'
import { SimpleIconsDocker } from './components/icons/simple-icons-docker'
import { SimpleIconsKubernetes } from './components/icons/simple-icons-kubernetes'
import { SimpleIconsOci } from './components/icons/simple-icons-oci'

export interface Node {
  id: string
  data: {
    title: string
    subtitle: string
    label: ReactNode
  }
  width: number
  height: number
  position: Position
}

export interface Position {
  x: number
  y: number
}

export interface Edge {
  id: string
  start: Position
  end: Position
}

const titles: Record<string, Record<string, string | undefined> | undefined> = {
  oci: {
    image: 'Image',
  },
  kubernetes: {
    'core/v1/pod': 'Pod',
    'core/v1/namespace': 'Namespace',
    'core/v1/container': 'Container',
    'apps/v1/deployment': 'Deployment',
    'apps/v1/daemonset': 'Daemon set',
    'apps/v1/replicaset': 'Replica set',
    'batch/v1/job': 'Job',
    'batch/v1/cronjob': 'Cron job',
    'apps/v1/statefulset': 'Stateful set',
    unknown: '<unknown resource>',
  },
  docker: {
    container: 'Container',
    'swarm/task': 'Task',
    'swarm/service': 'Service',
    'swarm/namespace': 'Namespace',
    'compose/service': 'Service',
    'compose/project': 'Project',
  },
}

function formatNode(id: string, node: GraphNode): Node {
  let label: ReactNode
  switch (node.domain) {
    case 'oci':
      label = <SimpleIconsOci className="text-blue-400" />
      break
    case 'kubernetes':
      label = <SimpleIconsKubernetes className="text-blue-400" />
      break
    case 'docker':
      label = <SimpleIconsDocker className="text-blue-500" />
  }

  return {
    id,
    data: {
      title: titles[node.domain]?.[node.type] || node.type,
      subtitle: node.name,
      label,
    },
    width: 0,
    height: 0,
    position: { x: 0, y: 0 },
  }
}

async function formatGraph(
  graph: Graph
): Promise<[Node[], Edge[], { width: number; height: number }]> {
  const nodes: Node[] = []
  const edges: Edge[] = []

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
    formatted.position.x = node.x ?? 0
    formatted.position.y = node.y ?? 0
    formatted.width = node.width ?? 0
    formatted.height = node.height ?? 0
    nodes.push(formatted)
  }

  // The node for the image is (should be) on a "row" of its own. Always center
  // the image node on that row
  for (const node of nodes) {
    if (
      graph.nodes[node.id].domain === 'oci' &&
      graph.nodes[node.id].type === 'image'
    ) {
      node.position.x = (root.width ?? 0) / 2 - (node.width ?? 0) / 2
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
    const startNode = nodes.find((x) => x.id === edge.sources[0])
    const endNode = nodes.find((x) => x.id === edge.targets[0])
    if (!startNode || !endNode) {
      continue
    }

    const start = { ...startNode.position }
    const end = { ...endNode.position }

    start.x += startNode.width / 2
    start.y += startNode.height

    end.x += startNode.width / 2

    edges.push({
      id: edge.id,
      start,
      end,
    })
  }

  return [nodes, edges, { width: root.width || 0, height: root.height || 0 }]
}

export function useNodesAndEdges(
  graph: Graph
): [Node[], Edge[], { width: number; height: number }] {
  const [nodes, setNodes] = useState<Node[]>([])
  const [edges, setEdges] = useState<Edge[]>([])
  const [bounds, setBounds] = useState<{ width: number; height: number }>({
    width: 0,
    height: 0,
  })

  useEffect(() => {
    formatGraph(graph).then(([nodes, edges, bounds]) => {
      setNodes(nodes)
      setEdges(edges)
      setBounds(bounds)
    })
  }, [graph])

  return [nodes, edges, bounds]
}
