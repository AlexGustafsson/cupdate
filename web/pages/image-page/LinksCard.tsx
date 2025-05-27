import type { JSX, ReactNode } from 'react'
import { FluentBookQuestionMark24Filled } from '../../components/icons/fluent-book-question-mark-24-filled'
import { FluentLink16Regular } from '../../components/icons/fluent-link-16-regular'
import { FluentLink24Filled } from '../../components/icons/fluent-link-24-filled'
import { Quay } from '../../components/icons/quay'
import { SimpleIconsDocker } from '../../components/icons/simple-icons-docker'
import { SimpleIconsGit } from '../../components/icons/simple-icons-git'
import { SimpleIconsGithub } from '../../components/icons/simple-icons-github'
import { SimpleIconsGitlab } from '../../components/icons/simple-icons-gitlab'
import type { ImageLink } from '../../lib/api/models'
import { Card } from './Card'

const titles: Record<string, string | undefined> = {
  github: 'GitHub project page',
  'github-releases': 'GitHub releases',
  ghcr: 'GHCR package page',
  gitlab: 'GitLab projectp age',
  docker: 'Docker Hub project page',
  quay: 'Quay project page',
  git: 'Git project page',
  docs: 'Documentation',
  'oci-registry': "Project's OCI registry",
}

function Link({ type, url }: { type: string; url: string }): JSX.Element {
  const { hostname } = new URL(url)
  const title = titles[type] || url
  let icon: ReactNode
  switch (type) {
    case 'github':
    case 'github-releases':
    case 'ghcr':
      icon = (
        <SimpleIconsGithub className="text-black dark:text-[#dddddd] shrink-0" />
      )
      break
    case 'gitlab':
      icon = <SimpleIconsGitlab className="text-orange-500 shrink-0" />
      break
    case 'docker':
      icon = <SimpleIconsDocker className="text-blue-500 shrink-0" />
      break
    case 'quay':
      icon = <Quay className="text-blue-700 shrink-0" />
      break
    case 'git':
      icon = <SimpleIconsGit className="text-orange-500 shrink-0" />
      break
    case 'svc':
      switch (hostname) {
        case 'github.com':
          return <Link type="github" url={url} />
        case 'gitlab.com':
          return <Link type="gitlab" url={url} />
        default:
          return <Link type="git" url={url} />
      }
    case 'docs':
      icon = (
        <FluentBookQuestionMark24Filled className="text-black dark:text-[#dddddd] shrink-0" />
      )
      break
    default:
      icon = (
        <FluentLink24Filled className="text-black dark:text-[#dddddd] shrink-0" />
      )
  }

  return (
    <a href={url}>
      <div className="flex items-center gap-x-2">
        {icon} <p className="m-0 flex-shrink-0">{title}</p>
      </div>
    </a>
  )
}

export type LinksCardProps = {
  links: ImageLink[]
}

export function LinksCard({ links }: LinksCardProps): JSX.Element {
  return (
    <Card
      persistenceKey="links"
      tabs={[
        {
          icon: <FluentLink16Regular />,
          label: 'Links',
          content: (
            <div className="markdown-body">
              <ul className="p-0">
                {links.map((x) => (
                  <li key={x.url} className="list-none w-min">
                    <Link type={x.type} url={x.url} />
                  </li>
                ))}
              </ul>
            </div>
          ),
        },
      ]}
    />
  )
}
