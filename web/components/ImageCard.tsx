import type { JSX } from 'react'

import { NavLink } from 'react-router-dom'
import { TagsByName } from '../tags'
import { formatRelativeTimeTo } from '../time'
import { Badge } from './Badge'
import { ImageLogo } from './ImageLogo'
import { InfoTooltip } from './InfoTooltip'
import { FluentShieldError24Filled } from './icons/fluent-shield-error-24-filled'
import { FluentWarning16Filled } from './icons/fluent-warning-16-filled'
import { SimpleIconsOci } from './icons/simple-icons-oci'

export type ImageCardProps = {
  name: string
  currentVersion: string
  latestVersion?: string
  logo?: string
  vulnerabilities: number
  updated?: Date
  description?: string
  tags: string[]
}

export function ImageCard({
  name,
  logo,
  currentVersion,
  latestVersion,
  updated,
  vulnerabilities,
  description,
  tags,
  className,
}: ImageCardProps & { className?: string }): JSX.Element {
  return (
    <div
      className={`flex gap-x-4 p-4 md:p-6 bg-white dark:bg-[#1e1e1e] rounded-lg shadow ${className || ''}`}
    >
      <ImageLogo src={logo} width={48} height={48} />
      <div className="flex flex-col w-full">
        <div className="flex flex-col sm:flex-row w-full sm:items-center sm:justify-between">
          <div className="flex items-center">
            <p className="text-sm break-all font-semibold">{name}</p>
            {vulnerabilities > 0 && (
              <InfoTooltip
                icon={<FluentShieldError24Filled className="text-red-600" />}
              >
                {vulnerabilities} vulnerabilities reported.
              </InfoTooltip>
            )}
          </div>
          <div className="flex flex-row items-center gap-x-2">
            {latestVersion ? (
              currentVersion === latestVersion ? (
                <p className="text-green-600">{latestVersion}</p>
              ) : (
                <>
                  <p className="text-red-600 line-through">{currentVersion}</p>
                  <p className="text-green-600">{latestVersion}</p>
                </>
              )
            ) : (
              <>
                <p className="text-yellow-600">{currentVersion}</p>
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
        <p className="text-sm mt-2">{description}</p>
        <div className="flex flex-wrap gap-2 mt-4">
          {tags.map((x) => (
            <NavLink key={x} to={`/?tag=${encodeURIComponent(x)}`}>
              <Badge
                label={x}
                color={TagsByName[x]?.color}
                className="hover:opacity-90"
              />
            </NavLink>
          ))}
        </div>
      </div>
    </div>
  )
}
