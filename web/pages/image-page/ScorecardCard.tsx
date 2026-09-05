import type { JSX } from 'react'
import { Card } from '../../components/Card'
import { FluentBoxSearch16Regular } from '../../components/icons/fluent-box-search-16-regular'
import type { ImageScorecard } from '../../lib/api/models'
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
          strokeWidth="3"
          strokeDasharray="50 100"
          strokeLinecap="round"
        />

        <circle
          cx="18"
          cy="18"
          r="16"
          fill="none"
          className="stroke-current"
          strokeWidth="1"
          strokeDasharray={`${percentage * 50} 100`}
          strokeLinecap="round"
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
  id: string
  scorecard: ImageScorecard
  collapsed: boolean
  onToggleCollapsed: () => void
}
export function ScorecardCard({
  id,
  scorecard,
  collapsed,
  onToggleCollapsed,
}: ScorecardCardProps): JSX.Element {
  const color =
    scorecard.score <= 2.5
      ? 'text-negative'
      : scorecard.score <= 5.0
        ? 'text-negative'
        : scorecard.score <= 7.5
          ? 'text-warning'
          : 'text-positive'

  return (
    <Card
      id={id}
      collapsed={collapsed}
      onToggleCollapsed={onToggleCollapsed}
      tabs={[
        {
          icon: <FluentBoxSearch16Regular />,
          label: 'Risk score',
          content: (
            <div className="markdown-body">
              <div className="flex justify-center">
                <Gauge
                  className={color}
                  percentage={scorecard.score / 10}
                  value={scorecard.score.toString()}
                  label={`${scorecard.risk} risk`}
                />
              </div>
              <p>
                The project associated with this image has been found to pose a{' '}
                <span className="font-semibold">{scorecard.risk} risk</span>,
                scoring{' '}
                <span className="font-semibold">
                  {scorecard.score}
                  /10
                </span>{' '}
                on{' '}
                <a
                  target="_blank"
                  rel="noreferrer"
                  href="https://scorecard.dev"
                >
                  Open Source Security Foundation's Scorecard
                </a>
                . The report was generated{' '}
                {formatRelativeTimeTo(new Date(scorecard.generatedAt))}. For
                more details, see the{' '}
                <a target="_blank" rel="noreferrer" href={scorecard.reportUrl}>
                  full report
                </a>
                .
              </p>
            </div>
          ),
        },
      ]}
    />
  )
}
