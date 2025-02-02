import {
  type JSX,
  type ReactNode,
  useCallback,
  useMemo,
  useRef,
  useState,
} from 'react'
import type { JobRun, StepRun, WorkflowRun } from '../../api'
import { GraphRenderer } from '../../components/GraphRenderer'
import { FluentCheckmarkCircle20Filled } from '../../components/icons/fluent-checkmark-circle-20-filled'
import { FluentCheckmarkCircle20Regular } from '../../components/icons/fluent-checkmark-circle-20-regular'
import { FluentDismissCircle20Filled } from '../../components/icons/fluent-dismiss-circle-20-filled'
import { useGraphLayout } from '../../graph'
import { formatDuration, formatRelativeTimeTo } from '../../time'

function Job({
  data,
}: {
  data: JobRun
}): JSX.Element {
  let label: ReactNode
  let status: string
  switch (data.result) {
    case 'succeeded':
      label = (
        <FluentCheckmarkCircle20Filled className="text-green-500 dark:text-green-700" />
      )
      status = `Succeded in ${formatDuration(data.duration)}`
      break
    case 'skipped':
      label = <FluentCheckmarkCircle20Regular className="text-gray-500" />
      status = 'Skipped'
      break
    case 'failed':
      label = (
        <FluentDismissCircle20Filled className="text-red-500 dark:text-red-700" />
      )
      status = `Failed after ${formatDuration(data.duration)}`
  }

  return (
    <div className="px-4 py-2 cursor-pointer shadow-md hover:shadow-lg rounded-md bg-white dark:bg-[#262626] border-2 border-[#ebebeb] dark:border-[#333333]">
      <div className="flex">
        <div
          className={`rounded-full w-12 h-12 flex justify-center items-center ${data.result === 'succeeded' ? 'bg-green-400/20 dark:bg-green-800/20' : data.result === 'skipped' ? 'bg-gray-100 dark:bg-[#363a3a]' : 'bg-red-400/20 dark:bg-red-800/20'} shrink-0`}
        >
          {label}
        </div>
        <div className="ml-2 grow min-w-0">
          <div className="text-lg font-bold truncate">{data.jobName}</div>
          <div className="text-gray-500 truncate">{status}</div>
        </div>
      </div>
    </div>
  )
}

export type WorkflowRunCardProps = {
  workflowRun: WorkflowRun
}

type StepRunListItemProps = {
  stepRun: StepRun
}

function StepRunListItem({ stepRun }: StepRunListItemProps): JSX.Element {
  let icon: ReactNode
  switch (stepRun.result) {
    case 'succeeded':
      icon = (
        <FluentCheckmarkCircle20Filled className="text-green-500 dark:text-green-700" />
      )
      break
    case 'skipped':
      icon = <FluentCheckmarkCircle20Regular className="text-gray-500" />
      break
    case 'failed':
      icon = (
        <FluentDismissCircle20Filled className="text-red-500 dark:text-red-700" />
      )
  }

  return (
    <>
      <div className="flex w-full gap-x-2 items-center">
        {icon}
        <p
          className={`flex-grow truncate text-sm ${stepRun.result === 'skipped' ? 'opacity-60' : ''}`}
        >
          {stepRun.stepName}
        </p>
        <p className="truncate text-sm opacity-60">
          {stepRun.duration ? formatDuration(stepRun.duration) : ''}
        </p>
      </div>
      {stepRun.error && (
        <pre className="p-4 rounded-lg bg-[#fafafa] dark:bg-[#121212] text-xs">
          <code>{stepRun.error}</code>
        </pre>
      )}
    </>
  )
}

type JobRunDialogProps = {
  ref: React.RefObject<HTMLDialogElement | null>
  traceId: string | undefined
  jobRun: JobRun | undefined
}

function JobRunDialog({
  ref,
  traceId,
  jobRun,
}: JobRunDialogProps): JSX.Element {
  let status: string
  switch (jobRun?.result) {
    case 'succeeded':
      status = `Succeeded ${formatRelativeTimeTo(new Date(jobRun.started))} after ${formatDuration(jobRun.duration)}`
      break
    case 'skipped':
      status = 'Skipped'
      break
    case 'failed':
      status = `Failed ${formatRelativeTimeTo(new Date(jobRun.started))} after ${formatDuration(jobRun.duration)}`
      break
    default:
      status = ''
  }

  traceId = 'fe39288ccaa21e41c764bc74466a48ad'

  return (
    // biome-ignore lint/a11y/useKeyWithClickEvents: The dialog element already handles ESC
    <dialog
      ref={ref}
      className="backdrop:bg-black/20 backdrop:backdrop-blur-sm bg-transparent m-auto"
      onClick={(e) => e.target === ref.current && ref.current.close()}
    >
      <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow w-[90vw] max-w-[800px] max-h-[80vh] overflow-y-scroll">
        <p className="font-bold">{jobRun?.jobName}</p>
        <p className="text-sm opacity-60">{status}</p>
        <div className="mt-4 flex flex-col gap-y-4 overflow-y-scroll">
          {jobRun?.steps
            .filter((x) => x.stepName)
            .map((x, i) => (
              <StepRunListItem key={i.toString()} stepRun={x} />
            ))}
        </div>
        {traceId && (
          <p className="text-sm opacity-60 text-center mt-4">
            Trace id: {traceId}
          </p>
        )}
      </div>
    </dialog>
  )
}

export function WorkflowCard({
  workflowRun,
}: WorkflowRunCardProps): JSX.Element {
  const [formattedGraph, options] = useMemo(() => {
    return [
      {
        nodes: workflowRun.jobs.map((data, i) => ({
          id: data.jobId || i.toString(),
          width: 350,
          height: 75,
          data,
        })),
        edges: workflowRun.jobs.flatMap((job, i) =>
          job.dependsOn.map((dependency) => ({
            // Reverse order
            from: dependency,
            to: job.jobId || i.toString(),
          }))
        ),
      },
      {
        'elk.algorithm': 'mrtree',
        'elk.spacing.nodeNode': '50',
        'elk.direction': 'RIGHT',
      },
    ]
  }, [workflowRun])

  const [nodes, edges, bounds] = useGraphLayout<JobRun>(formattedGraph, options)

  const [jobRun, setJobRun] = useState<JobRun>()
  const dialogRef = useRef<HTMLDialogElement>(null)

  const showJobRun = useCallback((jobRun: JobRun) => {
    setJobRun(jobRun)
    dialogRef.current?.showModal()
  }, [])

  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-2 shadow-sm h-[480px]">
      <JobRunDialog
        ref={dialogRef}
        traceId={workflowRun.traceId}
        jobRun={jobRun}
      />
      <GraphRenderer
        edges={edges}
        nodes={nodes}
        bounds={bounds}
        onNodeClick={(node) => showJobRun(node.data)}
        NodeElement={Job}
      />
    </div>
  )
}
