import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { useMemo } from 'react'

export function Markdown({
  children,
}: React.PropsWithChildren<{}>): JSX.Element {
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

  return <div dangerouslySetInnerHTML={{ __html: purified }} />
}
