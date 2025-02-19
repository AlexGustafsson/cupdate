import {
  type JSX,
  type ReactNode,
  useCallback,
  useMemo,
  useRef,
  useState,
} from 'react'
import type { Graph, GraphNode } from '../../api'
import { DependencyGraphNode } from '../../components/DependencyGraphNode'
import { GraphRenderer } from '../../components/GraphRenderer'
import { SimpleIconsDocker } from '../../components/icons/simple-icons-docker'
import { SimpleIconsKubernetes } from '../../components/icons/simple-icons-kubernetes'
import { SimpleIconsOci } from '../../components/icons/simple-icons-oci'
import { useGraphLayout } from '../../graph'

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

type GraphNodeProps = {
  ref: React.RefObject<HTMLDialogElement | null>
  graphNode: GraphNode | undefined
}

function GraphNodeDialog({ ref, graphNode }: GraphNodeProps): JSX.Element {
  let label: ReactNode
  switch (graphNode?.domain) {
    case 'oci':
      label = <SimpleIconsOci className="text-blue-400" />
      break
    case 'kubernetes':
      label = <SimpleIconsKubernetes className="text-blue-400" />
      break
    case 'docker':
      label = <SimpleIconsDocker className="text-blue-500" />
  }

  return (
    // biome-ignore lint/a11y/useKeyWithClickEvents: The dialog element already handles ESC
    <dialog
      ref={ref}
      className="backdrop:bg-black/20 backdrop:backdrop-blur-sm bg-transparent m-auto"
      onClick={(e) => e.target === ref.current && ref.current.close()}
    >
      <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow w-[90vw] max-w-[800px] max-h-[80vh] overflow-y-scroll">
        <p className="font-bold">
          {graphNode
            ? titles[graphNode.domain]?.[graphNode.type] || graphNode.type
            : ''}
        </p>
        <p className="text-sm opacity-60 break-all">{graphNode?.name}</p>
        {graphNode?.labels && (
          <>
            <p className="mt-2">Labels</p>
            <ul className="text-sm">
              {graphNode?.labels &&
                Object.entries(graphNode.labels).map(([k, v]) => (
                  <li key={`${k}`}>
                    <code>
                      {k}: {v}
                    </code>
                  </li>
                ))}
            </ul>
          </>
        )}
      </div>
    </dialog>
  )
}

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

  const [graphNode, setGraphNode] = useState<GraphNode>()
  const dialogRef = useRef<HTMLDialogElement>(null)

  const showGraphNode = useCallback((graphNode: GraphNode) => {
    setGraphNode(graphNode)
    dialogRef.current?.showModal()
  }, [])

  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-2 shadow-sm h-[480px]">
      <GraphNodeDialog ref={dialogRef} graphNode={graphNode} />
      <GraphRenderer
        edges={edges}
        nodes={nodes}
        bounds={bounds}
        onNodeClick={(node) => showGraphNode(node.data)}
        NodeElement={DependencyGraphNode}
      />
    </div>
  )
}
