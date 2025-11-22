import { type JSX, useCallback, useState } from 'react'
import { InfoTooltip } from '../../components/InfoTooltip'
import { FluentArrowSync16Regular } from '../../components/icons/fluent-arrow-sync-16-regular'
import { FluentWarning16Filled } from '../../components/icons/fluent-warning-16-filled'
import { useEvents } from '../../EventProvider'
import { useScheduleScan } from '../../lib/api/ApiProvider'
import { formatRelativeTimeTo } from '../../time'

type ProcessStatusProps = {
  lastModified: string
  reference: string
}

export function ProcessStatus({
  lastModified: initialLastModified,
  reference,
}: ProcessStatusProps): JSX.Element {
  const scheduleScan = useScheduleScan()

  const [status, setStatus] = useState<
    'idle' | 'in-flight' | 'successful' | 'failed'
  >('idle')

  // Get the time from the image once, then rely on events to update it
  const [lastModified, setLastModified] = useState(
    new Date(initialLastModified)
  )

  const onSchedule = useCallback(() => {
    setStatus('in-flight')
    scheduleScan(reference)
      .then(() => setStatus('successful'))
      .catch(() => setStatus('failed'))
  }, [scheduleScan, reference])

  useEvents(
    (e) => {
      if (e.type === 'imageProcessed' && e.reference === reference) {
        // TODO: Use time from event rather then the current time
        setLastModified(new Date())

        // If we successfully queued the image for processing, clear the state
        // when the reference was processed
        if (status === 'successful') {
          setStatus('idle')
        }
      }
    },
    [reference, status]
  )

  return (
    <div className="flex justify-center">
      {status !== 'successful' && (
        <p>
          Last processed{' '}
          <span title={lastModified.toLocaleString()}>
            {formatRelativeTimeTo(lastModified)}
          </span>
        </p>
      )}
      <p>{status === 'successful' && 'Image is scheduled for processing'}</p>
      <button
        type="button"
        className="cursor-pointer"
        title={status === 'idle' ? 'Schedule update' : ''}
        onClick={onSchedule}
        disabled={status !== 'idle'}
      >
        {(status === 'idle' || status === 'in-flight') && (
          <FluentArrowSync16Regular
            className={`ml-1 hover:opacity-90 active:opacity-80 disabled:opacity-70 ${status === 'in-flight' ? 'animate-spin' : ''}`}
          />
        )}
        {status === 'failed' && (
          <InfoTooltip icon={<FluentWarning16Filled />}>
            Failed to schedule image. Try again later.
          </InfoTooltip>
        )}
      </button>
    </div>
  )
}
