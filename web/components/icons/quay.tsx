import type { SVGProps } from 'react'

export function Quay(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      role="img"
      aria-label="icon"
      xmlns="http://www.w3.org/2000/svg"
      width="24px"
      height="24px"
      viewBox="0 0 8.467 8.467"
      {...props}
    >
      <path
        d="m6.626.343 1.84 3.89-1.84 3.89H5.061l1.838-3.89L5.061.344Z"
        fill="currentColor"
      />
      <path
        d="m5.06 8.123-1.839-3.89 1.84-3.89h1.565l-1.838 3.89 1.838 3.89Z"
        fill="currentColor"
      />
      <path
        d="M4.233 2.094 3.406.344H1.84L3.45 3.75ZM3.45 4.717 1.84 8.123h1.565l.827-1.75Z"
        fill="currentColor"
      />
      <path
        d="M1.84 8.123.002 4.233l1.84-3.89h1.565l-1.84 3.89 1.84 3.89Z"
        fill="currentColor"
      />
    </svg>
  )
}
