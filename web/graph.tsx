import type {
  Edge,
  Node,
  NodeTypes,
  OnEdgesChange,
  OnNodesChange,
} from '@xyflow/react'
import { useEdgesState, useNodesState } from '@xyflow/react'
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
    'batch/v1/cron': 'Cron job',
    'apps/v1/statefulset': 'Stateful set',
  },
}

// For now, let's assume that the all nodes under the root (its parents) have
// at most one parent node and that all nodes are unique. This enables us to
// render the graph more easily than taking care of cycles, references and
// layout solves
function extractBranches(root: GraphNode): GraphNode[][] {
  const branches: GraphNode[][] = []

  // for (let i = 0; i < root.parents.length; i++) {
  //   const branch: GraphNode[] = []
  //   let current = root.parents[i]
  //   while (current) {
  //     branch.push(current)
  //     current = current.parents[0]
  //   }
  //   branches.push(branch)
  // }

  return branches
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

function formatGraph(graph: Graph): [NodeType[], EdgeType[]] {
  const nodes: NodeType[] = []
  const edges: EdgeType[] = []

  for (const [id, node] of Object.entries(graph.nodes)) {
    nodes.push(formatNode(id, node))

    for (const [to, isParent] of Object.entries(graph.edges[id])) {
      if (isParent) {
        edges.push({
          id: `${id}->${to}`,
          // NOTE: The graph / tree is built with images as roots, but in the UI
          // we wish to map them with the root at the bottom, i.e. invert the
          // tree
          source: to,
          target: id,
        })
      }
    }
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
    const [nodes, edges] = formatGraph(graph)
    setNodes(nodes)
    setEdges(edges)
  }, [graph])

  return [
    [nodes, onNodesChange],
    [edges, onEdgesChange],
  ]
}
