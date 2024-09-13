import DOMPurify from 'dompurify'

export function HTML({ children }: React.PropsWithChildren<{}>): JSX.Element {
  if (typeof children !== 'string') {
    throw new Error('invalid HTML')
  }

  const purified = DOMPurify.sanitize(children)

  return <div dangerouslySetInnerHTML={{ __html: purified }} />
}
