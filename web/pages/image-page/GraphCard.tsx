import { type JSX, useMemo } from 'react'
import type { Graph, GraphNode } from '../../api'
import { DependencyGraphNode } from '../../components/DependencyGraphNode'
import { GraphRenderer } from '../../components/GraphRenderer'
import { useGraphLayout } from '../../graph'

type GraphCardProps = {
  graph: Graph
}

export function GraphCard({ graph }: GraphCardProps): JSX.Element {
  const [formattedGraph, options] = useMemo(() => {
    return [
      {
        nodes: Object.entries(graph.nodes).map(([id, data]) => ({
          id,
          width: 250,
          height: 75,
          data,
        })),
        edges: Object.entries(graph.edges)
          .map(([from, adjacent]) =>
            Object.entries(adjacent)
              .filter(([_, isParent]) => isParent)
              .map(([to, _]) => ({
                id: `${from}->${to}`,
                // Reverse the tree order
                from: to,
                to: from,
              }))
          )
          .flat(2),
      },
      {
        'elk.algorithm': 'mrtree',
        'elk.spacing.nodeNode': '50',
      },
    ]
  }, [graph])

  const [nodes, edges, bounds] = useGraphLayout<GraphNode>(
    formattedGraph,
    options
  )

  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-2 shadow-sm h-[480px]">
      <GraphRenderer
        edges={edges}
        nodes={nodes}
        bounds={bounds}
        NodeElement={DependencyGraphNode}
      />
    </div>
  )
}
