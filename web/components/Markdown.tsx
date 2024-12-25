import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { type JSX, useMemo } from 'react'

export function Markdown({
  children,
}: React.PropsWithChildren<Record<never, never>>): JSX.Element {
  if (typeof children !== 'string') {
    throw new Error('invalid HTML')
  }

  const purified = useMemo(() => {
    try {
      const dom = marked.parse(children, { async: false })
      return DOMPurify.sanitize(dom)
    } catch (error) {
      console.log('Failed to render markdown', error)
      return ''
    }
  }, [children])

  // biome-ignore lint/security/noDangerouslySetInnerHtml: the DOM is purified
  return <div dangerouslySetInnerHTML={{ __html: purified }} />
}
