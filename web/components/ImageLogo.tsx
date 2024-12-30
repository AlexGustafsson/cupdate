import { type JSX, useCallback, useRef, useState } from 'react'
import { SimpleIconsOci } from './icons/simple-icons-oci'

export type ImageLogoProps = {
  src?: string
  width: number
  height: number
}

export function ImageLogo({ src, width, height }: ImageLogoProps): JSX.Element {
  const ref = useRef<HTMLImageElement>(null)
  const [isLoading, setIsLoading] = useState(src !== undefined)

  const onLoad = useCallback(() => {
    setIsLoading(false)
  }, [])

  return (
    <div className={`w-[${width}px] h-[${height}px] shrink-0`}>
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
