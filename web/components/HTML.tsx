import DOMPurify from 'dompurify'
import type { JSX } from 'react'

export function HTML({
  children,
}: React.PropsWithChildren<Record<never, never>>): JSX.Element {
  if (typeof children !== 'string') {
    throw new Error('invalid HTML')
  }

  const purified = DOMPurify.sanitize(children)

  // biome-ignore lint/security/noDangerouslySetInnerHtml: the DOM is purified
  return <div dangerouslySetInnerHTML={{ __html: purified }} />
}
