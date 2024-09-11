import { Handle, NodeProps, Position } from '@xyflow/react'
import { memo } from 'react'

function CustomNode({
  data,
}: NodeProps & {
  data: any
  type: any
}): JSX.Element {
  return (
    <div className="px-4 py-2 shadow-md rounded-md bg-white border-2 border-stone-400">
      <div className="flex">
        <div className="rounded-full w-12 h-12 flex justify-center items-center bg-gray-100">
          {data.label}
        </div>
        <div className="ml-2">
          <div className="text-lg font-bold">{data.title}</div>
          <div className="text-gray-500">{data.subtitle}</div>
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

export default memo(CustomNode)
