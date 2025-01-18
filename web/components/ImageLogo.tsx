import { type JSX, useCallback, useRef, useState } from 'react'
import { SimpleIconsOci } from './icons/simple-icons-oci'

import CupdateLogoURL from '../favicon.png'

export type ImageLogoProps = {
  reference: string
  src?: string
  width: number
  height: number
}

export function ImageLogo({
  reference,
  src,
  width,
  height,
}: ImageLogoProps): JSX.Element {
  const ref = useRef<HTMLImageElement>(null)
  const [isLoading, setIsLoading] = useState(src !== undefined)

  // Don't show the default artwork for Cupdate
  // SEE: https://github.com/opencontainers/image-spec/issues/1231
  if (!src && reference.includes('ghcr.io/alexgustafsson/cupdate')) {
    src = CupdateLogoURL
  }

  const onLoad = useCallback(() => {
    setIsLoading(false)
  }, [])

  return (
    <div
      className="shrink-0"
      style={{ width: `${width}px`, height: `${height}px` }}
    >
      {src && (
        <img
          ref={ref}
          loading="lazy"
          alt="logo"
          src={src}
          className={`w-full h-full rounded transition-opacity ${isLoading ? 'opacity-0' : 'opacity-1'}`}
          width={width}
          height={height}
          onLoad={onLoad}
        />
      )}
      {!src && (
        <div className="flex items-center justify-center w-full h-full rounded bg-blue-400 dark:bg-blue-700">
          <SimpleIconsOci className="w-2/3 h-2/3 text-white dark:text-[#dddddd]" />
        </div>
      )}
    </div>
  )
}
