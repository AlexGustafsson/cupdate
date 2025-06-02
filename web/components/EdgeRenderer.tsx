import { type JSX, useMemo } from 'react'

export interface Edge {
  id: string
  start: {
    x: number
    y: number
  }
  end: {
    x: number
    y: number
  }
  className?: string
}

export type EdgeRendererProps = {
  edges: Edge[]
  direction: 'top-down' | 'left-right'
}

export function EdgeRenderer({
  edges,
  direction,
}: EdgeRendererProps): JSX.Element {
  const beziers = useMemo(() => {
    return edges.map((edge) => {
      if (direction === 'left-right') {
        // Straight line
        if (Math.abs(edge.start.y - edge.end.y) < 5) {
          return [
            `M ${edge.start.x},${edge.start.y}`,
            `L ${edge.end.x},${edge.end.y}`,
          ].join(' ')
        }

        const direction = edge.end.y < edge.start.y ? 1 : -1

        const offset = 30
        const radius = 10

        return [
          `M ${edge.start.x},${edge.start.y}`,
          // Line right
          `L ${edge.end.x - offset - radius},${edge.start.y}`,
          // Curve 90 degree
          `Q ${edge.end.x - offset + radius},${edge.start.y} ${edge.end.x - offset + radius},${edge.start.y - direction * radius}`,
          // Line up
          `L ${edge.end.x - offset + radius},${edge.end.y + direction * radius}`,
          // Curve 90 degree
          `Q ${edge.end.x - offset + radius},${edge.end.y} ${edge.end.x - offset + 2 * radius},${edge.end.y}`,
          // Line right
          `L ${edge.end.x},${edge.end.y}`,
        ].join(' ')
      }

      if (direction === 'top-down') {
        // Straight line
        if (Math.abs(edge.start.x - edge.end.x) < 5) {
          return [
            `M ${edge.start.x},${edge.start.y}`,
            `L ${edge.end.x},${edge.end.y}`,
          ].join(' ')
        }

        const direction = edge.end.x < edge.start.x ? 1 : -1

        const offset = 30
        const radius = 10

        return [
          `M ${edge.start.x},${edge.start.y}`,
          // Line down
          `L ${edge.start.x},${edge.start.y + offset - radius}`,
          // Curve 90 degree
          `Q ${edge.start.x},${edge.start.y + offset} ${edge.start.x - direction * radius},${edge.start.y + offset}`,
          // Line sideways
          `L ${edge.end.x + direction * radius},${edge.start.y + offset}`,
          // Curve 90 degree
          `Q ${edge.end.x},${edge.end.y - offset + radius} ${edge.end.x},${edge.end.y - offset + 2 * radius}`,
          // // Line down
          `L ${edge.end.x},${edge.end.y}`,
        ].join(' ')
      }
    })
  }, [edges, direction])

  return (
    <svg role="img" aria-label="Graph edge" className="w-full h-full">
      <g>
        {beziers.map((bezier, i) => (
          <path
            key={bezier}
            d={bezier}
            className={`transition-all fill-none stroke-2 stroke-[#ebebeb] dark:stroke-[#333333] ${edges[i].className ?? ''}`}
          />
        ))}
      </g>
    </svg>
  )
}
