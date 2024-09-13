import type {
  Edge,
  Node,
  NodeTypes,
  OnEdgesChange,
  OnNodesChange,
} from '@xyflow/react'
import { useEdgesState, useNodesState } from '@xyflow/react'
import { type ReactNode, useEffect } from 'react'

import type { Graph, GraphNode, Result } from './api'
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
  },
}

// For now, let's assume that the all nodes under the root (its parents) have
// at most one parent node and that all nodes are unique. This enables us to
// render the graph more easily than taking care of cycles, references and
// layout solves
function extractBranches(root: GraphNode): GraphNode[][] {
  const branches: GraphNode[][] = []

  for (let i = 0; i < root.parents.length; i++) {
    const branch: GraphNode[] = []
    let current = root.parents[i]
    while (current) {
      branch.push(current)
      current = current.parents[0]
    }
    branches.push(branch)
  }

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

function formatGraph(root: GraphNode): [NodeType[], EdgeType[]] {
  const branches = extractBranches(root)

  const rootNode = formatNode('root', root)
  const nodes: NodeType[] = [rootNode]
  const edges: EdgeType[] = []

  const offsetX =
    -(branches.length * 150 + Math.max(branches.length - 1, 0) * 100) / 2 + 75
  const offsetY = -100
  for (let i = 0; i < branches.length; i++) {
    let previous = rootNode
    for (let j = 0; j < branches[i].length; j++) {
      const x = offsetX + i * 150 + i * 100
      const y = offsetY - j * 100

      const node = formatNode(`$${i}$${y}$`, branches[i][j])
      node.position = { x, y }
      nodes.push(node)

      edges.push({
        id: `${node.id} -> ${previous.id}`,
        source: node.id,
        target: previous.id,
      })

      previous = node
    }
  }

  return [nodes, edges]
}

export function useNodesAndEdges(
  graph: Result<Graph>
): [
  [NodeType[], OnNodesChange<NodeType>],
  [EdgeType[], OnEdgesChange<EdgeType>],
] {
  const [nodes, setNodes, onNodesChange] = useNodesState<NodeType>([])
  const [edges, setEdges, onEdgesChange] = useEdgesState<EdgeType>([])

  useEffect(() => {
    if (graph.status !== 'resolved') {
      return
    }

    const [nodes, edges] = formatGraph(graph.value.root)
    setNodes(nodes)
    setEdges(edges)
  }, [graph])

  return [
    [nodes, onNodesChange],
    [edges, onEdgesChange],
  ]
}
