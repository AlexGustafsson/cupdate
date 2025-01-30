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
}

type EdgeRendererParams = {
  edges: Edge[]
}

export function EdgeRenderer({ edges }: EdgeRendererParams): JSX.Element {
  const beziers = useMemo(() => {
    return edges.map((edge) => {
      const pathStart = `M ${edge.start.x},${edge.start.y}`
      const pathStartControl = `C ${edge.start.x},${edge.start.y + 30}`

      const pathEnd = `${edge.end.x},${edge.end.y}`
      const pathEndControl = `${edge.end.x},${edge.end.y - 30}`

      return [pathStart, pathStartControl, pathEndControl, pathEnd].join(' ')
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
