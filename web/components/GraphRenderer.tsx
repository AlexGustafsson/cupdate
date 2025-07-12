import { type ComponentType, type JSX, useEffect, useState } from 'react'
import { type Edge, EdgeRenderer, type EdgeRendererProps } from './EdgeRenderer'
import { Surface } from './Surface'

export type NodeProps<T> = {
  data: T
  className?: string
}

export interface Node<T> {
  id: string
  width?: number
  height?: number
  x?: number
  y?: number
  className?: string
  data: T
}

type GraphRendererProps<T> = {
  nodes: Node<T>[]
  edges: Edge[]
  bounds: { width: number; height: number }
  direction: EdgeRendererProps['direction']
  onNodeClick?: (node: Node<T>) => void
  onNodeHover?: (node: string | undefined) => void
  NodeElement: ComponentType<NodeProps<T>>
}

export function GraphRenderer<T>({
  nodes,
  edges,
  bounds: { width, height },
  direction,
  onNodeClick,
  onNodeHover,
  NodeElement,
}: GraphRendererProps<T>): JSX.Element {
  const [hoveredNode, setHoveredNode] = useState<string>()

  useEffect(() => {
    if (onNodeHover) {
      onNodeHover(hoveredNode)
    }
  }, [hoveredNode, onNodeHover])

  return (
    <div className="w-full h-full">
      <Surface>
        <div
          className="relative"
          role="tree"
          style={{ width: `${width}px`, height: `${height}px` }}
        >
          <EdgeRenderer edges={edges} direction={direction} />
          {nodes.map((node) => (
            // biome-ignore lint/a11y/useKeyWithClickEvents: Nodes cannot be focused using keyboard
            // biome-ignore lint/a11y/useFocusableInteractive: Nodes cannot be focused using keyboard
            <div
              role="treeitem"
              key={node.id}
              className="absolute"
              onClick={() => onNodeClick?.(node)}
              onPointerEnter={() => setHoveredNode(node.id)}
              onPointerLeave={() =>
                setHoveredNode((current) =>
                  current === node.id ? undefined : current
                )
              }
              style={{
                top: node.y === undefined ? undefined : `${node.y}px`,
                left: node.x === undefined ? undefined : `${node.x}px`,
                width: node.width === undefined ? undefined : `${node.width}px`,
                height:
                  node.height === undefined ? undefined : `${node.height}px`,
              }}
            >
              <NodeElement data={node.data} className={node.className} />
            </div>
          ))}
        </div>
      </Surface>
    </div>
  )
}
