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
import { FluentBranch16Regular } from '../../components/icons/fluent-branch-16-regular'
import { SimpleIconsDocker } from '../../components/icons/simple-icons-docker'
import { SimpleIconsKubernetes } from '../../components/icons/simple-icons-kubernetes'
import { SimpleIconsOci } from '../../components/icons/simple-icons-oci'
import { useGraphLayout } from '../../graph'
import { parse } from '../../oci'
import { Card } from './Card'

const titles: Record<string, Record<string, string | undefined> | undefined> = {
  oci: {
    image: 'Image',
  },
  kubernetes: {
    'core/v1/node': 'Node',
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
    host: 'Host',
  },
}

const internalLabels: Record<string, string> = {
  'host-architecture': 'Architecture',
  'host-operating-system': 'Operating system',
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

  const oci =
    graphNode?.domain === 'oci' && graphNode.type === 'image'
      ? parse(graphNode.name)
      : null

  return (
    // biome-ignore lint/a11y/useKeyWithClickEvents: The dialog element already handles ESC
    <dialog
      ref={ref}
      className="starting:backdrop:bg-black/0 backdrop:bg-black/20 backdrop:backdrop-blur-sm bg-transparent m-auto backdrop:transition-all"
      onClick={(e) => e.target === ref.current && ref.current.close()}
    >
      <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow w-[90vw] max-w-[800px] max-h-[80vh] overflow-y-scroll">
        <p className="font-bold">
          {graphNode
            ? titles[graphNode.domain]?.[graphNode.type] || graphNode.type
            : ''}
        </p>
        {graphNode?.domain !== 'oci' && (
          <p className="text-sm opacity-60 break-all">{graphNode?.name}</p>
        )}
        {oci && (
          <>
            <p className="mt-2">Details</p>
            <ul className="text-sm">
              <li>
                <code>
                  Registry:{' '}
                  {oci.name.substring(0, oci.name.indexOf('/')) || 'docker.io'}
                </code>
              </li>
              <li>
                <code>
                  Name:{' '}
                  {oci.name.includes('.')
                    ? oci.name.substring(oci.name.indexOf('/') + 1)
                    : oci.name}
                </code>
              </li>
              <li>
                <code>Tag: {oci.tag}</code>
              </li>
              <li>
                <code className="break-all">Digest: {oci.digest}</code>
              </li>
            </ul>
          </>
        )}
        {graphNode?.internalLabels && (
          <>
            <p className="mt-2">Details</p>
            <ul className="text-sm">
              {graphNode?.internalLabels &&
                Object.entries(graphNode.internalLabels).map(([k, v]) => (
                  <li key={`${k}`}>
                    <code>
                      {internalLabels[k] || k}: {v}
                    </code>
                  </li>
                ))}
            </ul>
          </>
        )}
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
    <Card
      persistenceKey="graph"
      tabs={[
        {
          icon: <FluentBranch16Regular />,
          label: 'Graph',
          content: (
            <div className="h-[480px]">
              <GraphNodeDialog ref={dialogRef} graphNode={graphNode} />
              <GraphRenderer
                edges={edges}
                nodes={nodes}
                bounds={bounds}
                direction="top-down"
                onNodeClick={(node) => showGraphNode(node.data)}
                NodeElement={DependencyGraphNode}
              />
            </div>
          ),
        },
      ]}
    />
  )
}
