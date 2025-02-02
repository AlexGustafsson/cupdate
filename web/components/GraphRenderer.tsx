import type { ComponentType, JSX } from 'react'
import { type Edge, EdgeRenderer } from './EdgeRenderer'
import { Surface } from './Surface'

export type NodeProps<T> = {
  data: T
}

export interface Node<T> {
  id: string
  width?: number
  height?: number
  x?: number
  y?: number
  data: T
}

type GraphRendererProps<T> = {
  nodes: Node<T>[]
  edges: Edge[]
  bounds: { width: number; height: number }
  onNodeClick?: (node: Node<T>) => void
  NodeElement: ComponentType<NodeProps<T>>
}

export function GraphRenderer<T>({
  nodes,
  edges,
  bounds: { width, height },
  onNodeClick,
  NodeElement,
}: GraphRendererProps<T>): JSX.Element {
  return (
    <div className="w-full h-full">
      <Surface>
        <div
          className="relative"
          style={{ width: `${width}px`, height: `${height}px` }}
        >
          <EdgeRenderer edges={edges} />
          {nodes.map((node) => (
            // biome-ignore lint/a11y/useKeyWithClickEvents: Nodes cannot be focused using keyboard
            <div
              key={node.id}
              className="absolute"
              onClick={() => onNodeClick?.(node)}
              style={{
                top: node.y === undefined ? undefined : `${node.y}px`,
                left: node.x === undefined ? undefined : `${node.x}px`,
                width: node.width === undefined ? undefined : `${node.width}px`,
                height:
                  node.height === undefined ? undefined : `${node.height}px`,
              }}
            >
              <NodeElement data={node.data} />
            </div>
          ))}
        </div>
      </Surface>
    </div>
  )
}
