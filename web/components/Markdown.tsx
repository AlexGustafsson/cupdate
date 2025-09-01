import { marked } from 'marked'
import { type JSX, useMemo } from 'react'
import { HTML } from './HTML'

export function Markdown({
  children,
}: React.PropsWithChildren<Record<never, never>>): JSX.Element {
  if (typeof children !== 'string') {
    throw new Error('invalid HTML')
  }

  const html = useMemo(
    () => marked.parse(children, { async: false }),
    [children]
  )

  if (typeof html !== 'string') {
    console.log(html)
  }

  return <HTML>{html}</HTML>
}
