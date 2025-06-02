import ELK, { type LayoutOptions } from 'elkjs/lib/elk.bundled'
import { useEffect, useState } from 'react'

export interface Node<T> {
  id: string
  width?: number
  height?: number
  x?: number
  y?: number
  data: T
}

export interface Edge {
  id: string
  start: {
    nodeId: string
    x: number
    y: number
  }
  end: {
    nodeId: string
    x: number
    y: number
  }
}

export interface Graph<T> {
  nodes: Node<T>[]
  edges: { from: string; to: string }[]
}

async function graphLayout<T>(
  graph: Graph<T>,
  layoutOptions?: LayoutOptions
): Promise<[Node<T>[], Edge[], { width: number; height: number }]> {
  const nodes: Node<T>[] = []
  const edges: Edge[] = []

  const elk = new ELK()

  const root = await elk.layout({
    id: 'root',
    layoutOptions,
    children: graph.nodes,
    edges: graph.edges.map((edge, i) => ({
      id: i.toString(),
      sources: [edge.from],
      targets: [edge.to],
    })),
  })

  // ELK can return nodes that are offset from the top left and bounds that are
  // larger than what's necessary. Recalculate offsets and bounds to return a
  // minimum rect, making things like centering of the graph easier.
  const leftOffset =
    root.children?.reduce(
      (leftOffset, x) => (x.x && x.x < leftOffset ? x.x : leftOffset),
      Number.MAX_VALUE
    ) || 0

  const topOffset =
    root.children?.reduce(
      (topOffset, x) => (x.y && x.y < topOffset ? x.y : topOffset),
      Number.MAX_VALUE
    ) || 0

  const width =
    root.children?.reduce(
      (width, x) =>
        x.x && x.width && x.x + x.width > width ? x.x + x.width : width,
      0
    ) || 0

  const height =
    root.children?.reduce(
      (height, x) =>
        x.y && x.height && x.y + x.height > height ? x.y + x.height : height,
      0
    ) || 0

  // Add nodes
  for (const node of root.children || []) {
    nodes.push({
      id: node.id,
      width: node.width,
      height: node.height,
      x: (node.x || 0) - leftOffset,
      y: (node.y || 0) - topOffset,
      data: graph.nodes.find((x) => x.id === node.id)!.data,
    })
  }

  // Add edges
  for (const edge of root.edges || []) {
    const startNode = nodes.find((x) => x.id === edge.sources[0])
    const endNode = nodes.find((x) => x.id === edge.targets[0])
    if (!startNode || !endNode) {
      continue
    }

    const start = { x: startNode.x || 0, y: startNode.y || 0 }
    const end = { x: endNode.x || 0, y: endNode.y || 0 }

    if (layoutOptions?.['elk.direction'] === 'RIGHT') {
      start.x += startNode.width || 0
      start.y += (startNode.height || 0) / 2

      end.y += (endNode.height || 0) / 2
    } else {
      start.x += (startNode.width || 0) / 2
      start.y += startNode.height || 0

      end.x += (startNode.width || 0) / 2
    }

    edges.push({
      id: edge.id,
      start: {
        nodeId: startNode.id,
        x: start.x,
        y: start.y,
      },
      end: {
        nodeId: endNode.id,
        x: end.x,
        y: end.y,
      },
    })
  }

  return [
    nodes,
    edges,
    { width: width - leftOffset, height: height - topOffset },
  ]
}

export function useGraphLayout<T>(
  graph: Graph<T>,
  layoutOptions?: LayoutOptions
): [Node<T>[], Edge[], { width: number; height: number }] {
  const [nodes, setNodes] = useState<Node<T>[]>([])
  const [edges, setEdges] = useState<Edge[]>([])
  const [bounds, setBounds] = useState<{ width: number; height: number }>({
    width: 0,
    height: 0,
  })

  useEffect(() => {
    graphLayout(graph, layoutOptions).then(([nodes, edges, bounds]) => {
      setNodes(nodes)
      setEdges(edges)
      setBounds(bounds)
    })
  }, [graph, layoutOptions])

  return [nodes, edges, bounds]
}
