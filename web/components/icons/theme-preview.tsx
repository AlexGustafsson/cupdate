import type { SVGProps } from 'react'

export function ThemePreview(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      role="img"
      aria-label="icon"
      xmlns="http://www.w3.org/2000/svg"
      width="175"
      height="98"
      viewBox="0 0 46.302 25.929"
      {...props}
    >
      {/* Background */}
      <path d="M0 0h46.302v25.929H0z" fill="var(--surface-0-bg)" />
      {/* Title */}
      <path d="M19.778 1.342h6.747v1.386h-6.747z" fill="var(--surface-0-fg)" />
      {/* RSS button */}
      <path
        d="M42.349 1.342h3.088v1.386h-3.088z"
        fill="var(--color-brand-rss)"
      />
      {/* Summaries */}
      <path d="M4.99 5.56h5.065v1.358H4.99Z" fill="var(--surface-0-fg)" />
      <path d="M12.378 5.538h4.324v1.359h-4.324z" fill="var(--surface-0-fg)" />
      <path d="M20.619 5.68h3.583v1.359H20.62z" fill="var(--surface-0-fg)" />
      <path d="M28.945 5.715h5.065v1.359h-5.065z" fill="var(--surface-0-fg)" />
      <path d="M36.247 5.715h3.88v1.359h-3.88z" fill="var(--surface-0-fg)" />
      <path d="M4.824 7.993h1.853v1.36H4.824Z" fill="var(--color-positive)" />
      <path
        d="M12.95 7.993h1.852v1.36H12.95zM20.712 7.993h1.853v1.36h-1.853z"
        fill="var(--color-negative)"
      />
      <path
        d="M29.09 7.993h1.854v1.36H29.09zM36.388 7.993h1.853v1.36h-1.853Z"
        fill="var(--color-negative)"
      />
      {/* Inputs */}
      <path
        d="M4.3 11.363h37.702v1.057H4.3z"
        fill="var(--surface-1-bg)"
        stroke="var(--surface-1-stroke)"
        strokeWidth="0.2"
      />
      <path
        d="M4.237 13.28h11.859v1.184H4.237z"
        fill="var(--surface-1-bg)"
        stroke="var(--surface-1-stroke)"
        strokeWidth="0.2"
      />
      <path
        d="M18.046 13.28h11.772v1.185H18.046z"
        fill="var(--surface-1-bg)"
        stroke="var(--surface-1-stroke)"
        strokeWidth="0.2"
      />
      <path
        d="M31.598 13.275h10.393v1.195H31.598z"
        fill="var(--surface-1-bg)"
        stroke="var(--surface-1-stroke)"
        strokeWidth="0.2"
      />
      {/* Image card */}
      <rect
        width="37.31"
        height="5.436"
        x="4.496"
        y="15.644"
        rx=".529"
        ry=".529"
        fill="var(--surface-1-bg)"
      />
      {/* Logo */}
      <rect
        width="3.175"
        height="3.175"
        x="5.674"
        y="16.775"
        rx=".529"
        ry=".529"
        fill="var(--color-accent)"
      />
      {/* Name / description */}
      <path d="M10.492 16.405h9.142v1.112h-9.142z" fill="var(--surface-1-fg)" />
      <path d="M10.504 17.986h19.89v1.112h-19.89z" fill="var(--surface-1-fg)" />
      {/* Tags */}
      <rect
        width="2.965"
        height="1.112"
        x="10.535"
        y="19.467"
        rx=".265"
        ry=".265"
        fill="var(--color-negative)"
      />
      <rect
        width="2.965"
        height="1.112"
        x="14.364"
        y="19.467"
        rx=".265"
        fill="var(--color-warning)"
      />
      <rect
        width="2.965"
        height="1.112"
        x="18.018"
        y="19.467"
        rx=".265"
        ry=".265"
        fill="var(--color-major)"
      />
      <rect
        width="2.965"
        height="1.112"
        x="21.689"
        y="19.467"
        rx=".265"
        ry=".265"
        fill="var(--color-brand-kubernetes)"
      />
      <rect
        width="2.965"
        height="1.112"
        x="25.375"
        y="19.467"
        rx=".265"
        ry=".265"
        fill="var(--color-brand-github)"
      />
      {/* Old version */}
      <path
        d="M33.942 16.45h2.965v1.112h-2.965z"
        fill="var(--color-negative)"
      />
      {/* New version */}
      <path
        d="M38.075 16.45h2.965v1.112h-2.965z"
        fill="var(--color-positive)"
      />
    </svg>
  )
}
