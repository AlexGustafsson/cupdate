import DOMPurify from 'dompurify'
import { type JSX, useMemo } from 'react'

export function HTML({
  children,
}: React.PropsWithChildren<Record<never, never>>): JSX.Element {
  if (typeof children !== 'string') {
    throw new Error('invalid HTML')
  }

  const purified = useMemo(() => {
    const purify = DOMPurify()

    purify.addHook('afterSanitizeElements', async (node) => {
      if (node instanceof HTMLElement) {
        // Set referrer-policy for elements with src
        switch (node.tagName.toLowerCase()) {
          case 'a':
          case 'area':
          case 'img':
          case 'video':
          case 'iframe':
          case 'script':
            node.setAttribute('referrer-policy', 'no-referrer')
        }
      }
    })

    return purify.sanitize(children, { ADD_ATTR: ['referrer-policy'] })
  }, [children])

  // biome-ignore lint/security/noDangerouslySetInnerHtml: the DOM is purified
  return <div dangerouslySetInnerHTML={{ __html: purified }} />
}
