import { type JSX, useMemo } from 'react'
import type { Edge, Node } from '../graph'
import { Surface } from './Surface'

export type GraphNode = Node

export function GraphNode({
  width,
  height,
  data: { subtitle, title, label },
}: GraphNode): JSX.Element {
  return (
    <div
      className="px-4 py-2 shadow-md rounded-md bg-white dark:bg-[#262626] border-2 border-[#ebebeb] dark:border-[#333333]"
      style={{ width: `${width}px`, height: `${height}px` }}
    >
      <div className="flex">
        <div className="rounded-full w-12 h-12 flex justify-center items-center bg-gray-100 dark:bg-[#363a3a] shrink-0">
          {label}
        </div>
        <div className="ml-2 grow min-w-0">
          <div className="text-lg font-bold truncate">{title}</div>
          <div className="text-gray-500 truncate">{subtitle}</div>
        </div>
      </div>
    </div>
  )
}

type EdgesParams = {
  edges: Edge[]
}

function Edges({ edges }: EdgesParams): JSX.Element {
  const beziers = useMemo(() => {
    return edges.map((edge) => {
      const pathStart = `M${edge.start.x},${edge.start.y}`
      const pathStartControl = `C${edge.start.x}, ${edge.start.y + 30}`

      const pathEnd = `${edge.end.x},${edge.end.y}`
      const pathEndControl = `${edge.end.x},${edge.end.y - 30}`

      return `${pathStart}, ${pathStartControl}, ${pathEndControl}, ${pathEnd}, `
    })
  }, [edges])

  return (
    <svg role="img" aria-label="Graph edge" className="w-full h-full">
      <g>
        {beziers.map((bezier) => (
          <path
            key={bezier}
            d={bezier}
            className="fill-none stroke-2 stroke-[#ebebeb] dark:stroke-[#333333]"
          />
        ))}
      </g>
    </svg>
  )
}

type GraphProps = {
  nodes: Node[]
  edges: Edge[]
  bounds: { width: number; height: number }
}

export function Graph({
  nodes,
  edges,
  bounds: { width, height },
}: GraphProps): JSX.Element {
  return (
    <div className="w-full h-full">
      <Surface>
        <div
          className="relative"
          style={{ width: `${width}px`, height: `${height}px` }}
        >
          <Edges edges={edges} />
          {nodes.map((node) => (
            <div
              key={node.id}
              className="absolute"
              style={{
                top: `${node.position.y}px`,
                left: `${node.position.x}px`,
              }}
            >
              <GraphNode {...node} />
            </div>
          ))}
        </div>
      </Surface>
    </div>
  )
}
