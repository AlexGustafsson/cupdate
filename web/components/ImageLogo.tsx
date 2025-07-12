import { type JSX, useCallback, useRef, useState } from 'react'
import { SimpleIconsOci } from './icons/simple-icons-oci'

import { useApiClient } from '../lib/api/ApiProvider'
import CupdateLogoURL from '../public/assets/icon.png'

export type ImageLogoProps = {
  reference: string
  src?: string
  className?: string
}

export function ImageLogo({
  reference,
  src,
  className,
}: ImageLogoProps): JSX.Element {
  const ref = useRef<HTMLImageElement>(null)
  const [isLoading, setIsLoading] = useState(src !== undefined)
  const [isError, setIsError] = useState(false)
  const apiClient = useApiClient()

  // Don't show the default artwork for Cupdate
  // SEE: https://github.com/opencontainers/image-spec/issues/1231
  if (!src && reference.includes('ghcr.io/alexgustafsson/cupdate')) {
    src = CupdateLogoURL
  }

  // Try to get the logo from the API (from disk, via user's config)
  if (!src) {
    src = apiClient.getLogoUrl(reference)
  }

  const onLoad = useCallback(() => {
    setIsLoading(false)
  }, [])

  const onError = useCallback(() => {
    setIsError(true)
  }, [])

  return (
    <div className={`shrink-0 ${className}`}>
      {src && !isError && (
        <img
          ref={ref}
          loading="lazy"
          alt="logo"
          src={src}
          className={`w-full h-full rounded-sm transition-opacity ${isLoading ? 'opacity-0' : 'opacity-100'}`}
          onLoad={onLoad}
          onError={onError}
        />
      )}
      {(!src || isError) && (
        <div className="flex items-center justify-center w-full h-full rounded-sm bg-blue-400 dark:bg-blue-700">
          <SimpleIconsOci className="w-2/3 h-2/3 text-white dark:text-[#dddddd]" />
        </div>
      )}
    </div>
  )
}
