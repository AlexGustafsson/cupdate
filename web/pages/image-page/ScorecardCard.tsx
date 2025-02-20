import type { JSX } from 'react'
import type { ImageScorecard } from '../../api'
import { formatRelativeTimeTo } from '../../time'

type GaugeProps = {
  percentage: number
  value: string
  label: string
  className: string
}

function Gauge({
  percentage,
  value,
  label,
  className,
}: GaugeProps): JSX.Element {
  return (
    <div className="relative w-32 h-32">
      <svg
        className={`size-full rotate-180 ${className}`}
        role="img"
        aria-label="icon"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 36 36"
      >
        <circle
          cx="18"
          cy="18"
          r="16"
          fill="none"
          className="stroke-current opacity-25"
          stroke-width="3"
          stroke-dasharray="50 100"
          stroke-linecap="round"
        />

        <circle
          cx="18"
          cy="18"
          r="16"
          fill="none"
          className="stroke-current"
          stroke-width="1"
          stroke-dasharray={`${percentage * 50} 100`}
          stroke-linecap="round"
        />
      </svg>

      <div className="absolute top-9 start-1/2 transform -translate-x-1/2 text-center">
        <span className="text-2xl font-bold">{value}</span>
        <span className="text-xs block">{label}</span>
      </div>
    </div>
  )
}

export type ScorecardCardProps = {
  scorecard: ImageScorecard
}
export function ScorecardCard({ scorecard }: ScorecardCardProps): JSX.Element {
  const color =
    scorecard.score <= 2.5
      ? 'text-red-400'
      : scorecard.score <= 5.0
        ? 'text-red-400'
        : scorecard.score <= 7.5
          ? 'text-orange-400'
          : 'text-yellow-400'

  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow">
      <div className="markdown-body">
        <h1>Risk score</h1>
        <div className="flex justify-center">
          <Gauge
            className={color}
            percentage={scorecard.score / 10}
            value={scorecard.score.toString()}
            label={scorecard.risk}
          />
        </div>
        <p>
          The project associated with this image has been found to pose a{' '}
          <span className="font-semibold">{scorecard.risk}</span> risk, scoring{' '}
          <span className="font-semibold">
            {scorecard.score}
            /10
          </span>{' '}
          on{' '}
          <a target="_blank" rel="noreferrer" href="https://scorecard.dev">
            Open Source Security Foundation's Scorecard
          </a>
          . The report was generated{' '}
          {formatRelativeTimeTo(new Date(scorecard.generatedAt))}. For more
          details, see the{' '}
          <a target="_blank" rel="noreferrer" href={scorecard.reportUrl}>
            full report
          </a>
          .
        </p>
      </div>
    </div>
  )
}
