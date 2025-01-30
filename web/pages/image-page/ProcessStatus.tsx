import { type JSX, useCallback, useState } from 'react'
import { scheduleScan } from '../../api'
import { InfoTooltip } from '../../components/InfoTooltip'
import { FluentArrowSync16Regular } from '../../components/icons/fluent-arrow-sync-16-regular'
import { FluentWarning16Filled } from '../../components/icons/fluent-warning-16-filled'
import { formatRelativeTimeTo } from '../../time'

type ProcessStatusProps = {
  lastModified: string
  reference: string
}

export function ProcessStatus({
  lastModified,
  reference,
}: ProcessStatusProps): JSX.Element {
  const [status, setStatus] = useState<
    'idle' | 'in-flight' | 'successful' | 'failed'
  >('idle')

  const onSchedule = useCallback(() => {
    setStatus('in-flight')
    scheduleScan(reference)
      .then(() => setStatus('successful'))
      .catch(() => setStatus('failed'))
  }, [reference])

  return (
    <>
      {status !== 'successful' && (
        <p>
          Last processed{' '}
          <span title={new Date(lastModified).toLocaleString()}>
            {formatRelativeTimeTo(new Date(lastModified))}
          </span>
        </p>
      )}
      <p>{status === 'successful' && 'Image is scheduled for processing'}</p>
      <button
        type="button"
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
            Failed to schedule image. Cupdate is likely busy. Try again later.
          </InfoTooltip>
        )}
      </button>
    </>
  )
}
