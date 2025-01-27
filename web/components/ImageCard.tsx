import type { JSX } from 'react'

import { NavLink, useNavigate } from 'react-router-dom'
import { TagsByName, compareTags } from '../tags'
import { formatRelativeTimeTo } from '../time'
import { Badge } from './Badge'
import { ImageLogo } from './ImageLogo'
import { InfoTooltip } from './InfoTooltip'
import { FluentShieldError24Filled } from './icons/fluent-shield-error-24-filled'
import { FluentWarning16Filled } from './icons/fluent-warning-16-filled'
import { SimpleIconsOci } from './icons/simple-icons-oci'

export type ImageCardProps = {
  reference: string
  name: string
  currentVersion: string
  fullCurrentVersion: string
  latestVersion?: string
  fullLatestVersion?: string
  logo?: string
  vulnerabilities: number
  updated?: Date
  description?: string
  tags: string[]
  compact?: boolean
}

export function ImageCard({
  reference,
  name,
  logo,
  currentVersion,
  fullCurrentVersion,
  latestVersion,
  fullLatestVersion,
  updated,
  vulnerabilities,
  description,
  tags,
  compact,
  className,
}: ImageCardProps & { className?: string }): JSX.Element {
  const navigate = useNavigate()

  return (
    <div
      className={`flex gap-x-4 bg-white dark:bg-[#1e1e1e] rounded-lg shadow-sm ${compact ? 'p-3' : 'p-4 md:p-6'} ${className || ''}`}
    >
      <ImageLogo
        src={logo}
        reference={reference}
        className="w-[48px] h-[48px]"
      />
      <div className="flex flex-col w-full">
        <div
          className={`${compact ? '' : 'flex flex-col sm:flex-row w-full sm:items-center sm:justify-between'}`}
        >
          <div
            className={`${compact ? 'flex flex-col-reverse' : 'flex items-center'}`}
          >
            <p className="text-sm break-all line-clamp-2 font-semibold">
              {name}
            </p>
            {vulnerabilities > 0 && (
              <InfoTooltip
                icon={<FluentShieldError24Filled className="text-red-600" />}
              >
                {vulnerabilities} vulnerabilities reported.
              </InfoTooltip>
            )}
          </div>
          <div className="flex flex-row items-center gap-x-2">
            {/* Digests are formatted like <algo>:<digest>, such as sha256:<digest>. Show a maximum of 5 hex digits before truncating with ellipsis (hence 15ch) */}
            {latestVersion ? (
              fullCurrentVersion === fullLatestVersion ? (
                <p
                  className="text-green-600 max-w-[15ch] truncate"
                  title={fullLatestVersion}
                >
                  {latestVersion}
                </p>
              ) : (
                <>
                  <p
                    className="text-red-600 line-through max-w-[15ch] truncate"
                    title={fullCurrentVersion}
                  >
                    {currentVersion}
                  </p>
                  <p
                    className="text-green-600 max-w-[15ch] truncate"
                    title={fullLatestVersion}
                  >
                    {latestVersion}
                  </p>
                </>
              )
            ) : (
              <>
                <p
                  className="text-yellow-600 max-w-[15ch] truncate"
                  title={fullCurrentVersion}
                >
                  {currentVersion}
                </p>
                {!latestVersion && (
                  <InfoTooltip
                    icon={<FluentWarning16Filled className="text-yellow-600" />}
                  >
                    The latest version cannot be identified. This could be due
                    to the image not being available, the registry not being
                    supported, missing authentication or a temporary issue.
                  </InfoTooltip>
                )}
              </>
            )}
          </div>
        </div>
        {updated && (
          <p className="text-sm">Updated {formatRelativeTimeTo(updated)}</p>
        )}
        <p className={`text-sm mt-2 ${compact ? 'line-clamp-2' : ''}`}>
          {description}
        </p>
        {!compact && (
          <div className="flex flex-wrap gap-2 mt-4">
            {tags
              .toSorted((a, b) => compareTags(a, b))
              .map((x) => (
                <Badge
                  key={x}
                  label={x}
                  color={TagsByName[x]?.color}
                  className="hover:opacity-90"
                  // It's illegal to nest anchors in HTML, so unfortunately we need
                  // to use onClick here
                  onClick={(e) => {
                    e.metaKey || e.ctrlKey
                      ? openTab(`/?tag=${encodeURIComponent(x)}`)
                      : navigate(`/?tag=${encodeURIComponent(x)}`)
                    e.preventDefault()
                  }}
                />
              ))}
          </div>
        )}
      </div>
    </div>
  )
}

function openTab(target: string) {
  // Using window.open with target _blank creates a new window on Safari, macOS
  // so use this cross-platform solution instead
  const a = document.createElement('a')

  a.rel = 'noreferrer'
  a.target = '_blank'
  a.href = target
  a.click()
}
