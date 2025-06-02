import type { JSX, ReactNode } from 'react'
import type { GraphNode } from '../lib/api/models'
import type { NodeProps } from './GraphRenderer'
import { InfoTooltip } from './InfoTooltip'
import { FluentTag16Regular } from './icons/fluent-tag-16-regular'
import { SimpleIconsDocker } from './icons/simple-icons-docker'
import { SimpleIconsKubernetes } from './icons/simple-icons-kubernetes'
import { SimpleIconsOci } from './icons/simple-icons-oci'

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

export function DependencyGraphNode({
  data,
}: NodeProps<GraphNode>): JSX.Element {
  let label: ReactNode
  switch (data.domain) {
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
    <div className="px-4 py-2 cursor-pointer hover:shadow-md transition-shadow rounded-md bg-white dark:bg-[#262626] border-2 border-[#ebebeb] dark:border-[#333333]">
      <div className="flex">
        <div className="rounded-full w-12 h-12 flex justify-center items-center bg-gray-100 dark:bg-[#363a3a] shrink-0">
          {label}
        </div>
        <div className="ml-2 grow min-w-0">
          <div className="flex items-center">
            <div className="text-lg font-bold truncate flex-grow">
              {titles[data.domain]?.[data.type] || data.type}
            </div>
            {data.labels && Object.keys(data.labels).length > 0 && (
              <InfoTooltip icon={<FluentTag16Regular />}>
                This node has labels that might affect how the image is
                processed.
              </InfoTooltip>
            )}
          </div>
          <div className="text-gray-500 truncate">{data.name}</div>
        </div>
      </div>
    </div>
  )
}
