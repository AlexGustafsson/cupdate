import type { JSX } from 'react'

export type WordBreakProps = {
  delimiter: string
  children: string
}

/**
 * WordBreak a text at a specific string / character.
 * Useful to word break at specific delimiters like "/" where CSS options would
 * otherwise cause breaking at unwanted characters.
 */
export function WordBreak({
  children,
  delimiter,
}: WordBreakProps): JSX.Element {
  const parts = children.split(delimiter)

  return (
    <>
      {parts.map((part, i) => (
        <>
          {part}
          {i < parts.length - 1 && (
            <>
              /
              <wbr />
            </>
          )}
        </>
      ))}
    </>
  )
}
