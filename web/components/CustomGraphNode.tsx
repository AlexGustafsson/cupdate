import { Handle, type NodeProps, Position } from '@xyflow/react'
import type { JSX, ReactNode } from 'react'

export function CustomGraphNode(
  node: NodeProps & {
    data: {
      subtitle: string
      title: string
      label: ReactNode
    }
    type: 'custom'
  }
): JSX.Element {
  return (
    <div
      className="px-4 py-2 shadow-md rounded-md bg-white dark:bg-[#262626] border-2 border-[#ebebeb] dark:border-[#333333]"
      style={{ width: node.width, height: node.height }}
    >
      <div className="flex">
        <div className="rounded-full w-12 h-12 flex justify-center items-center bg-gray-100 dark:bg-[#363a3a] shrink-0">
          {node.data.label}
        </div>
        <div className="ml-2 grow min-w-0">
          <div className="text-lg font-bold truncate">{node.data.title}</div>
          <div className="text-gray-500 truncate">{node.data.subtitle}</div>
        </div>
      </div>

      <Handle
        type="target"
        position={Position.Top}
        className="opacity-0 pointer-events-none !cursor-grab"
      />
      <Handle
        type="source"
        position={Position.Bottom}
        className="opacity-0 pointer-events-none !cursor-grab"
      />
    </div>
  )
}
